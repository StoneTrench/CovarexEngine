package modding

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/StoneTrench/CovarexEngine/engine/platform"
)

type DependencySeverity uint8

const (
	DependencySeverity_Info DependencySeverity = iota
	DependencySeverity_Warning
	DependencySeverity_Error
)

type DependencyIssue struct {
	ModuleId string
	TargetId string
	Op       string
	Message  string
	Severity DependencySeverity
}

func (e *DependencyIssue) Error() string {
	res := e.Message

	if len(e.Op) > 0 {
		res = fmt.Sprintf("%s: %s", e.Op, res)
	}

	if e.ModuleId != e.TargetId {
		res = fmt.Sprintf("%s has a problem with %s: %s", e.ModuleId, e.TargetId, res)
	} else {
		res = fmt.Sprintf("%s has a problem: %s", e.ModuleId, res)
	}

	return res
}

type HeaderCollection struct {
	Headers []GuestModuleHeader
	IdMap   map[string]uint
}

func (h *HeaderCollection) Init() *HeaderCollection {
	h.Headers = make([]GuestModuleHeader, 0)
	h.IdMap = make(map[string]uint)

	return h
}

const (
	// error if exists
	checkfor_NOT_EXISTS uint8 = iota
	// error if not exists
	checkfor_EXISTS
)

func validate_map(coll HeaderCollection, issues *[]DependencyIssue, cnst map[string]semver.Constraints, mod_id string, op string, severity DependencySeverity, checkfor uint8) {
	for other_id, constraint := range cnst {
		other, other_exists := coll.IdMap[other_id]
		var err error

		switch checkfor {
		case checkfor_NOT_EXISTS:
			if other_exists {
				head := coll.Headers[other]
				if ok := constraint.Check(head.Version); ok {
					err = fmt.Errorf("incompatible module")
				}
			}
		case checkfor_EXISTS:
			if other_exists {
				head := coll.Headers[other]
				if ok, errs := constraint.Validate(head.Version); !ok {
					err = errors.Join(errs...)
				}
			} else {
				err = fmt.Errorf("missing module")
			}
		}

		if err != nil {
			*issues = append(*issues, DependencyIssue{
				ModuleId: mod_id,
				TargetId: other_id,
				Op:       op,
				Severity: severity,
				Message:  err.Error(),
			})
		}
	}
}

func single_mod_check_issues(coll HeaderCollection, issues *[]DependencyIssue, h *GuestModuleHeader) {
	game_version := platform.GetVersion()

	// check game version
	if ok, errs := h.Compat.GameVersion.Validate(game_version); !ok {
		*issues = append(*issues, DependencyIssue{
			ModuleId: h.ModuleId,
			TargetId: h.ModuleId,
			Op:       "check_game_version",
			Severity: DependencySeverity_Error,
			Message:  errors.Join(errs...).Error(),
		})
	}

	// check exists
	validate_map(
		coll, issues, h.Compat.Depends,
		h.ModuleId, "check_depends", DependencySeverity_Error, checkfor_EXISTS,
	)
	validate_map(
		coll, issues, h.Compat.Recommends,
		h.ModuleId, "check_recommends", DependencySeverity_Warning, checkfor_EXISTS,
	)
	validate_map(
		coll, issues, h.Compat.Suggests,
		h.ModuleId, "check_suggests", DependencySeverity_Info, checkfor_EXISTS,
	)
	validate_map(
		coll, issues, h.Compat.Conflicts,
		h.ModuleId, "check_conflicts", DependencySeverity_Warning, checkfor_NOT_EXISTS,
	)
	validate_map(
		coll, issues, h.Compat.Breaks,
		h.ModuleId, "check_breaks", DependencySeverity_Error, checkfor_NOT_EXISTS,
	)

	// check order
	// something something topological sort
	mod_index := coll.IdMap[h.ModuleId]
	for _, other_id := range h.Compat.LoadBefore {
		if other_index, exists := coll.IdMap[other_id]; exists {
			if other_index > mod_index {
				*issues = append(*issues, DependencyIssue{
					ModuleId: h.ModuleId,
					TargetId: other_id,
					Op:       "check_load_before",
					Severity: DependencySeverity_Error,
					Message:  "target should load before the module",
				})
			}
		}
	}
	for _, other_id := range h.Compat.LoadAfter {
		if other_index, exists := coll.IdMap[other_id]; exists {
			if other_index < mod_index {
				*issues = append(*issues, DependencyIssue{
					ModuleId: h.ModuleId,
					TargetId: other_id,
					Op:       "check_load_after",
					Severity: DependencySeverity_Error,
					Message:  "target should load after the module",
				})
			}
		}
	}
}

// Analyze module list, and return list of issues between mods
func CheckForHeaderIssues(coll HeaderCollection) []DependencyIssue {
	issues := []DependencyIssue{}

	for i := range coll.Headers {
		h := &coll.Headers[i]
		single_mod_check_issues(coll, &issues, h)
	}

	return issues
}

// Sort the module list
// func AutoSortModules(coll *HeaderCollection) []DependencyIssue {

// }
