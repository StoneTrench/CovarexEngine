package resources

import (
	"fmt"
	"iter"
	"log/slog"
	"slices"
	"sync"

	"golang.org/x/exp/constraints"
)

/*

1. Register resources
1. Freeze resources, Increment version counter (Permanent until reset)
1. Can query by identifier
	1. Apply id map, Validate ids
	1. Can query by id
	1. Invalidate ids, Increment version counter
	1. Repeat
1. Reset
1. Register resources
1. Repeat

There's two phases of the registry's lifespan:
+ Registering resources
+ Querying for resources, updating ids, saving id mappings
The numerical ids of resources, can be assigned after the registration phase is over.

Q: What happens when a new mod is added?
A: Id mappings that already exist, get correctly assigned, then the registry handles the new resources, by just assigning new, unique Ids.

Q: What happens when a mod is removed?
A: In the save file we store the mapping of resources, and a list of mods used. The game will issue a warning of missing mods.

Q: What happens when a mod updates names?
A: Don't.

Q: What happens when mods are loaded in a different order?
A: Order doesn't matter, because resources can be loaded in any order, as their numerical ids are handled after registration.

Q: How specifically does id mapping work?
A:

*/

type RegistryErr[T any, I constraints.Integer] struct {
	RegName    string
	RegFrozen  bool
	RegValid   bool
	RegVersion uint32
	Ops        string
	Err        error
}

func (e *RegistryErr[T, I]) Error() string {
	return fmt.Sprintf("registry (name:%q,frozen:%v,valid:%v,version:%d): operation %s: %s",
		e.RegName, e.RegFrozen, e.RegValid, e.RegVersion,
		e.Ops, e.Err.Error(),
	)
}

func new_error[T any, I constraints.Integer](reg *Registry[T, I], op string, err error) error {
	return &RegistryErr[T, I]{
		RegName:    reg.name,
		RegFrozen:  reg.is_frozen,
		RegValid:   reg.is_id_valid,
		RegVersion: reg.version,
		Ops:        op,
		Err:        err,
	}
}
func new_errorf[T any, I constraints.Integer](reg *Registry[T, I], op string, format string, args ...any) error {
	return new_error(reg, op, fmt.Errorf(format, args...))
}

const (
	REG_ERR_FROZEN     = "is already frozen, registration phase over"
	REG_ERR_NOT_FROZEN = "is not yet frozen, registration phase still underway"

	REG_ERR_INVALIDATED     = "entries have already been invalidated, waiting for new id map"
	REG_ERR_NOT_INVALIDATED = "entries have not been invalidated, cannot apply new id map"

	REG_ERR_EXIST_IDENT     = "entry with identifier %q already exists"
	REG_ERR_NOT_EXIST_IDENT = "entry with identifier %q does not exist"
	REG_ERR_NOT_EXIST_ID    = "entry with id %d does not exist"

	REG_ERR_INVALID_NAME = "entry name %q contains invalid characters or has a length of 0"
	REG_ERR_INVALID_MOD  = "mod name %q contains invalid characters or has a length of 0"

	REG_ERR_DUPLICATE_MAPPING = "the identifier %q appears more than once in the mapping"
)

type Registry[T any, I constraints.Integer] struct {
	mu sync.RWMutex

	name string

	resources map[string]T

	is_frozen   bool
	is_id_valid bool // when you leave a game, ids get invalidated, when you join a game, the ids get loaded

	version uint32 // registry_version is incremented every time the the registry gets invalidated, when it is frozen and when a new id mapping is applied.

	id2ident_arr []string
	ident2id_map map[string]I
}

// Init doesnt lock the registry for reading or writing.
func (r *Registry[T, I]) Init(name string) *Registry[T, I] {
	r.name = name
	r.resources = make(map[string]T)
	r.is_frozen = false
	r.is_id_valid = false
	r.version = 0
	r.id2ident_arr = make([]string, 0)
	r.ident2id_map = make(map[string]I)
	return r
}

// Reset calls init, but also locks the registry for writing.
func (r *Registry[T, I]) Reset() *Registry[T, I] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Init(r.name)
}

func (r *Registry[T, I]) GetVersion() uint32 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.version
}

func (r *Registry[T, I]) Register(mod string, name string, resource T) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.is_frozen {
		return "", new_errorf(r, "register", REG_ERR_FROZEN)
	}
	if IsNameInvalid(mod) {
		return "", new_errorf(r, "register", REG_ERR_INVALID_MOD, mod)
	}
	if IsNameInvalid(name) {
		return "", new_errorf(r, "register", REG_ERR_INVALID_NAME, name)
	}
	ident := fmt.Sprintf("%s:%s", mod, name)
	if _, exists := r.resources[ident]; exists {
		return "", new_errorf(r, "register", REG_ERR_EXIST_IDENT, ident)
	}
	r.resources[ident] = resource
	return ident, nil
}
func (r *Registry[T, I]) RegisterPanic(mod string, name string, resource T) string {
	n, err := r.Register(mod, name, resource)
	if err != nil {
		panic(err)
	}
	return n
}

func (r *Registry[T, I]) Freeze() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.is_frozen {
		return new_errorf(r, "freeze", REG_ERR_FROZEN)
	}
	r.is_frozen = true
	r.version++
	return nil
}
func (r *Registry[T, I]) InvalidateIds() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.is_frozen {
		return new_errorf(r, "invalidate_ids", REG_ERR_NOT_FROZEN)
	}
	if !r.is_id_valid {
		return new_errorf(r, "invalidate_ids", REG_ERR_INVALIDATED)
	}
	r.is_id_valid = false
	r.version++
	return nil
}

/*
ApplyIdMap applies the given m id->identifier map, and sets the state to valid.

Let R be the set of things that have been registered.

Let M be the set of things in the mapping expects to have been registered.

Then R*M' is the set of entries that were registered but have no mapping (new entries).

Then R'*M is the set of entries that were not registered but have a mapping (missing entries).
*/
func (r *Registry[T, I]) ApplyIdMap(m []string) (new_map, new_entries, missing_entries []string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.is_frozen {
		return nil, nil, nil, new_errorf(r, "apply_id_map", REG_ERR_NOT_FROZEN)
	}
	if r.is_id_valid {
		return nil, nil, nil, new_errorf(r, "apply_id_map", REG_ERR_NOT_INVALIDATED)
	}
	r.id2ident_arr = make([]string, len(m))
	copy(r.id2ident_arr, m)
	r.ident2id_map = make(map[string]I)

	// CONJ -> AND -> INTERSECTION -> *
	// DISJ -> OR -> UNION -> +

	resources_size := len(r.resources)

	registered_entries := make([]string, 0, resources_size)
	for ident := range r.resources {
		registered_entries = append(registered_entries, ident)
	}
	slices.Sort(registered_entries)

	// have not been registered but have a mapping
	missing_entries = make([]string, 0, resources_size)
	for id, ident := range r.id2ident_arr {
		if _, exists := r.ident2id_map[ident]; exists {
			return nil, nil, nil, new_errorf(r, "apply_id_map", REG_ERR_DUPLICATE_MAPPING, ident)
		}

		r.ident2id_map[ident] = I(id)

		if _, exists := r.resources[ident]; !exists {
			missing_entries = append(missing_entries, ident)
		}
	}

	// have been registered but have no mapping
	new_entries = make([]string, 0, resources_size)
	for _, ident := range registered_entries {

		if _, exists := r.ident2id_map[ident]; !exists {
			new_entries = append(new_entries, ident)
			r.ident2id_map[ident] = I(len(r.id2ident_arr))
			r.id2ident_arr = append(r.id2ident_arr, ident)
		}
	}

	if len(missing_entries) > 0 {
		slog.Warn(new_errorf(r, "apply_id_map", "missing registry entries when loading id map: %v", missing_entries).Error())
	}
	if len(new_entries) > 0 {
		slog.Info(new_errorf(r, "apply_id_map", "new registry entries when loading id map: %v", new_entries).Error())
	}

	r.is_id_valid = true
	r.version++
	new_map = make([]string, len(r.id2ident_arr))
	copy(new_map, r.id2ident_arr)
	return new_map, new_entries, missing_entries, nil
}

func (r *Registry[T, I]) impl_GetByIdent(ident string) (T, error) {
	if !r.is_frozen {
		var NA T
		return NA, new_errorf(r, "get_by_ident", "identifier %q: %s", ident, REG_ERR_NOT_FROZEN)
	}
	res, exists := r.resources[ident]
	if !exists {
		return res, new_errorf(r, "get_by_ident", REG_ERR_NOT_EXIST_IDENT, ident)
	}
	return res, nil
}
func (r *Registry[T, I]) impl_GetId(ident string) (I, error) {
	var NA I
	if !r.is_frozen {
		return NA, new_errorf(r, "get_id", "identifier %q: %s", ident, REG_ERR_NOT_FROZEN)
	}
	if !r.is_id_valid {
		return NA, new_errorf(r, "get_id", "identifier %q: %s", ident, REG_ERR_INVALIDATED)
	}
	res, exists := r.ident2id_map[ident]
	if !exists {
		return res, new_errorf(r, "get_id", REG_ERR_NOT_EXIST_IDENT, ident)
	}
	return res, nil
}
func (r *Registry[T, I]) impl_GetIdent(id I) (string, error) {
	var NA string
	if !r.is_frozen {
		return NA, new_errorf(r, "get_ident", "id %d: %s", id, REG_ERR_NOT_FROZEN)
	}
	if !r.is_id_valid {
		return NA, new_errorf(r, "get_ident", "id %d: %s", id, REG_ERR_INVALIDATED)
	}
	if 0 > id || id >= I(len(r.id2ident_arr)) {
		return NA, new_errorf(r, "get_ident", REG_ERR_NOT_EXIST_ID, id)
	}
	return r.id2ident_arr[id], nil
}
func (r *Registry[T, I]) impl_GetById(id I) (T, error) {
	var NA T
	ident, err := r.impl_GetIdent(id)
	if err != nil {
		return NA, err
	}
	return r.impl_GetByIdent(ident)
}

func (r *Registry[T, I]) GetByIdentifier(ident string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.impl_GetByIdent(ident)
}
func (r *Registry[T, I]) GetId(ident string) (I, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.impl_GetId(ident)
}
func (r *Registry[T, I]) GetIdentifier(id I) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.impl_GetIdent(id)
}
func (r *Registry[T, I]) GetById(id I) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.impl_GetById(id)
}
func (r *Registry[T, I]) GetIdPanic(ident string) I {
	res, err := r.GetId(ident)
	if err != nil {
		panic(err)
	}
	return res
}

// Creates a snapshot of the registry.
func (r *Registry[T, I]) Iterate() iter.Seq2[string, T] {
	r.mu.RLock()

	snapshot := make([]struct {
		key   string
		value T
	}, len(r.resources))
	i := 0
	for k, v := range r.resources {
		snapshot[i].key = k
		snapshot[i].value = v
		i++
	}

	r.mu.RUnlock()

	return func(yield func(string, T) bool) {
		for _, pair := range snapshot {
			if !yield(pair.key, pair.value) {
				return
			}
		}
	}
}

//	func (r *Registry[T, I]) RegisterAlias(mod string, name string, alias string) (string, error) {
//		r.mu.Lock()
//		defer r.mu.Unlock()
//		if r.is_frozen {
//			return "", new_errorf(r, "registeralias", REG_ERR_FROZEN)
//		}
//		original_ident := fmt.Sprintf("%s:%s", mod, name)
//		id, exists := r.ident2id_map[original_ident]
//		if !exists {
//			return "", new_errorf(r, "registeralias", REG_ERR_NOT_EXIST_IDENT, original_ident)
//		}
//		alias_ident := fmt.Sprintf("%s:%s", mod, alias)
//		if _, exists := r.ident2id_map[alias_ident]; exists {
//			return "", new_errorf(r, "registeralias", REG_ERR_EXIST_IDENT, alias_ident)
//		}
//		r.ident2id_map[alias_ident] = id
//		return alias_ident, nil
//	}
