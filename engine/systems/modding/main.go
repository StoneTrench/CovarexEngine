package modding

import (
	"context"
	"errors"
	"fmt"
	"path"

	fs "github.com/StoneTrench/CovarexEngine/engine/platform/filesystem"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func load_mod_headers(folder string, entries []fs.DirEntry) ([]GuestModule, error) {
	mods := []GuestModule{}
	err_list := []error{}
	for _, e := range entries {
		if e.IsDir() {
			mod := GuestModule{}
			mod.Init(path.Join(folder, e.Name()))

			err := mod.LoadHeader()
			if err != nil {
				err_list = append(err_list, fmt.Errorf("failed to load a mod header: %w", err))
			} else {
				mods = append(mods, mod)
			}
		}
	}
	return mods, errors.Join(err_list...)
}
func query_dependencies(mod_list []GuestModule, m GuestModule) ([]GuestModule, []error) {
	deps := m.GetDependencies()
	dep_count := len(deps)
	res := make([]GuestModule, dep_count)

	err_list := []error{}

	dep_index := 0
	for k, v := range deps {
		var found = false

		for _, dep := range mod_list {
			if dep.header.Name == k {
				if !v.Check(dep.GetVersion()) {
					err := fmt.Errorf(
						"wrong version of mod: (%s %s) expects (%s %s), got (%s %s) instead",
						m.GetName(),
						m.GetVersion().String(),
						k,
						v.String(),
						dep.GetName(),
						dep.GetVersion().String(),
					)

					err_list = append(err_list, err)
					continue
				}

				res[dep_index] = dep
				dep_index++
				found = true
			}
		}

		if !found {
			err := fmt.Errorf(
				"missing dependency: (%s %s) requires mod (%s %s)",
				m.GetName(),
				m.GetVersion().String(),
				k,
				v.String(),
			)
			err_list = append(err_list, err)
			continue
		}
	}

	return res, err_list
}
func sort_by_dependency_view(mod_list []GuestModule, mod GuestModule, visited map[string]struct{}, sorted *[]GuestModule) []error {
	if _, exists := visited[mod.header.Name]; exists {
		has := false
		for _, m := range *sorted {
			if mod.GetName() == m.GetName() &&
				mod.GetVersion().Equal(m.GetVersion()) {
				has = true
				break
			}
		}
		if !has {
			return []error{fmt.Errorf("cyclic dependency found in relation to mod (%s %s)", mod.GetName(), mod.GetVersion().String())}
		}
	}

	err_list := []error{}

	deps, err := query_dependencies(mod_list, mod)
	if err != nil {
		err_list = append(err_list, err...)
	}

	for _, e := range deps {
		err := sort_by_dependency_view(mod_list, e, visited, sorted)
		if err != nil {
			err_list = append(err_list, err...)
		}
	}

	*sorted = append((*sorted), mod)

	return err_list
}
func sort_by_dependency(mod_list []GuestModule) ([]GuestModule, error) {
	sorted := []GuestModule{}
	visited := make(map[string]struct{})

	err_list := []error{}

	for _, v := range mod_list {
		err := sort_by_dependency_view(mod_list, v, visited, &sorted)
		if err != nil {
			err_list = append(err_list, err...)
		}
	}

	return sorted, errors.Join(err_list...)
}

var LoadedModsList []GuestModule
var name_to_mod map[string]int

func LoadMods(ctx context.Context, rt wazero.Runtime, folder string) error {
	LoadedModsList = nil
	entries, err := fs.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("failed to load mods: %w", err)
	}

	// Load mod headers
	mods, err := load_mod_headers(folder, entries)
	if err != nil {
		return err
	}

	// Sort mods based on dependency
	// Skip mods with missing dependencies and log errors
	mods, err = sort_by_dependency(mods)
	if err != nil {
		return err
	}

	// create id map
	var error_arr []error
	name_to_mod = make(map[string]int)
	LoadedModsList = mods
	for i := range mods {
		m := &mods[i]
		id := m.GetName()
		if _, exists := name_to_mod[id]; exists {
			error_arr = append(error_arr, fmt.Errorf("duplicate mod id: %s", id))
		}
		name_to_mod[id] = i
	}

	// Load guest modules
	for i := range mods {
		m := &mods[i]
		err := m.LoadWasm(ctx, rt)
		if err != nil {
			error_arr = append(error_arr, err)
			continue
		}
	}

	// Finish
	return errors.Join(error_arr...)
}
func GetMod(m api.Module) *GuestModule {
	if LoadedModsList == nil {
		return nil
	}
	i, exists := name_to_mod[m.Name()]
	if !exists {
		return nil
	}
	return &LoadedModsList[i]
}
