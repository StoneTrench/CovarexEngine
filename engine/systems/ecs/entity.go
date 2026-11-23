package ecs

import (
	"sync"
)

// id_t is the entity_id and component_id. They're the same so we can have as many components as entities.
type id_t uint32

// get_t is for entity validation (because ids can get reassigned), and if it doesnt match the latest generation, then the entity was destroyed, so this reference is invalid.
type gen_t uint16

// EntityRef is a reference to an entity. Keeps track of the generation index from when the entity was created to check freshness.
type EntityRef struct {
	// the id of the entity
	id id_t

	// gen counts the number of times that this entity id has been reused
	gen gen_t
}

// EntityWorld is the object that keeps track of entity ids, their generation counter, and the components.
//
// Use it to create component sets and new entities.
type EntityWorld struct {
	mu         sync.RWMutex
	OpenSet    []id_t
	Generation []gen_t
	NextId     id_t

	ComponentSetNames map[string]int
	ComponentSets     []*ComponentSet
}

// Init sets every field inside EntityWorld to a new default value.
func (s *EntityWorld) Init() *EntityWorld {
	s.mu = sync.RWMutex{}
	s.OpenSet = []id_t{}
	s.Generation = []gen_t{}
	s.NextId = 0
	s.ComponentSetNames = make(map[string]int)
	s.ComponentSets = []*ComponentSet{}

	return s
}

// CreateEntity creates a new entity, either from the open set or by making a new id, and returns the reference to it.
func (s *EntityWorld) CreateEntity() EntityRef {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.OpenSet) == 0 {
		res := EntityRef{
			id:  s.NextId,
			gen: 0,
		}
		s.NextId++
		s.Generation = append(s.Generation, 0)
		return res
	}

	last_index := len(s.OpenSet) - 1
	id := s.OpenSet[last_index]
	s.OpenSet = s.OpenSet[:last_index]
	gens := s.Generation[id]

	return EntityRef{
		id:  id,
		gen: gens,
	}
}

// DestroyEntity destroys the entity, increments the global generation counter for it, and adds the id to the open set to be recycled.
func (s *EntityWorld) DestroyEntity(e EntityRef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Generation) <= int(e.id) {
		return
	}
	if s.Generation[e.id] != e.gen {
		return // because the entity has already been destroyed
	}
	for _, c := range s.ComponentSets {
		c.RemoveComponent(e)
	}
	s.Generation[e.id] += 1
	s.OpenSet = append(s.OpenSet, e.id)
}

// IsFresh returns true if this entity still exists, and wasn't deleted.
func (s *EntityWorld) IsFresh(e EntityRef) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.Generation) <= int(e.id) {
		return false
	}
	return s.Generation[e.id] == e.gen
}

// CreateComponentSet creates a component set, and automatically sizes it to fit every currently existing entity.
func (s *EntityWorld) CreateComponentSet(name string, size csize_t) {
	s.mu.Lock()
	defer s.mu.Unlock()

	set := &ComponentSet{
		mu:                sync.RWMutex{},
		CSize:             size,
		EntityToComponent: []id_t{},
		ComponentToEntity: []id_t{},
		Components:        []byte{},
	}
	if s.NextId > 0 {
		id := s.NextId - 1
		gen := s.Generation[id]
		set.TryResize(EntityRef{id: id, gen: gen})
	}

	id := len(s.ComponentSets)
	s.ComponentSets = append(s.ComponentSets, set)
	s.ComponentSetNames[name] = id
}
