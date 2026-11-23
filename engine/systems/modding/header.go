package modding

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/StoneTrench/CovarexEngine/engine/resources"
)

type HeaderCollection struct {
	Headers []*GuestModuleHeader
	IdMap   map[string]int
}

func (h *HeaderCollection) Init() *HeaderCollection {
	h.Headers = make([]*GuestModuleHeader, 0)
	h.IdMap = make(map[string]int)

	return h
}
func (h *HeaderCollection) Add(gst *GuestModuleHeader) {
	id := len(h.Headers)
	h.Headers = append(h.Headers, gst)
	h.IdMap[gst.ModuleId] = id
}

type json_GMH_0 struct {
	ModuleId   string `json:"moduleid"`
	ModVersion string `json:"version"`

	Links map[string]string `json:"links"`

	Compat struct {
		GameVersion string `json:"gameversion"`

		Depends    map[string]string `json:"depends"`
		Recommends map[string]string `json:"recommends"`
		Suggests   map[string]string `json:"suggests"`
		Conflicts  map[string]string `json:"conflicts"`
		Breaks     map[string]string `json:"breaks"`

		LoadAfter  []string `json:"loadafter"`
		LoadBefore []string `json:"loadbefore"`
	} `json:"compat"`
}
type GuestModuleHeader struct {
	// ModuleId should never change for a module
	ModuleId string
	Version  *semver.Version

	// Links
	Links map[string]string

	Compat struct {
		// This is a constraint for game version

		GameVersion semver.Constraints

		// These 5 here check if the module exists

		Depends    map[string]semver.Constraints // error when missing
		Recommends map[string]semver.Constraints // warning when missing
		Suggests   map[string]semver.Constraints // info when missing
		Conflicts  map[string]semver.Constraints // warning when exists
		Breaks     map[string]semver.Constraints // error when exists

		// These two determine their order

		LoadAfter  []string
		LoadBefore []string
	}
}

func parse_deps(out *map[string]semver.Constraints, inp *map[string]string) error {
	*out = make(map[string]semver.Constraints)

	for mod_id, vrsn := range *inp {
		if err := resources.IsNameInvalidError(mod_id); err != nil {
			return err
		}
		v, err := semver.NewConstraint(vrsn)
		if err != nil {
			return err
		}
		(*out)[mod_id] = *v
	}
	return nil
}
func ParseJsonHeaderExtra[T any](file_json []byte, out *GuestModuleHeader, out_extra *T) error {
	if out == nil {
		return fmt.Errorf("argument out was nil")
	}

	if out_extra != nil {
		err := json.Unmarshal(file_json, out_extra)
		if err != nil {
			return fmt.Errorf("failed to parse extra header: %w", err)
		}
	}

	var head json_GMH_0
	err := json.Unmarshal(file_json, &head)
	if err != nil {
		return fmt.Errorf("failed to parse header: %w", err)
	}

	*out = GuestModuleHeader{}
	if err := resources.IsNameInvalidError(head.ModuleId); err != nil {
		return fmt.Errorf("failed to parse modid: %w", err)
	}
	out.ModuleId = head.ModuleId
	out.Version, err = semver.NewVersion(head.ModVersion)
	if err != nil {
		return fmt.Errorf("failed to parse version: %w", err)
	}

	out.Links = head.Links

	game_version, err := semver.NewConstraint(head.Compat.GameVersion)
	if err != nil {
		return fmt.Errorf("failed to parse game version: %w", err)
	}
	out.Compat.GameVersion = *game_version
	if err := parse_deps(&out.Compat.Depends, &head.Compat.Depends); err != nil {
		return fmt.Errorf("failed to parse depends list: %w", err)
	}
	if err := parse_deps(&out.Compat.Recommends, &head.Compat.Recommends); err != nil {
		return fmt.Errorf("failed to parse recommends list: %w", err)
	}
	if err := parse_deps(&out.Compat.Suggests, &head.Compat.Suggests); err != nil {
		return fmt.Errorf("failed to parse suggests list: %w", err)
	}
	if err := parse_deps(&out.Compat.Conflicts, &head.Compat.Conflicts); err != nil {
		return fmt.Errorf("failed to parse conflicts list: %w", err)
	}
	if err := parse_deps(&out.Compat.Breaks, &head.Compat.Breaks); err != nil {
		return fmt.Errorf("failed to parse breaks list: %w", err)
	}

	out.Compat.LoadAfter = head.Compat.LoadAfter
	out.Compat.LoadBefore = head.Compat.LoadBefore

	return nil
}
func ParseJsonHeader(file_json []byte, out *GuestModuleHeader) error {
	return ParseJsonHeaderExtra[any](file_json, out, nil)
}
