package modding

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/StoneTrench/CovarexEngine/engine/resources"
)

type json_SchemaVersion struct {
	SchemaVersion *int `json:"schemaVersion"`
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

const LATEST_SCHEMA = 0

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
func ParseJsonHeader(file_json []byte, out *GuestModuleHeader) error {
	var schema json_SchemaVersion
	err := json.Unmarshal(file_json, &schema)
	if err != nil {
		return fmt.Errorf("failed to get schema version from header: %w", err)
	}

	if schema.SchemaVersion == nil {
		i := LATEST_SCHEMA
		schema.SchemaVersion = &i
	}

	var head json_GMH_0
	switch *schema.SchemaVersion {
	case 0:
		err := json.Unmarshal(file_json, &head)
		if err != nil {
			return fmt.Errorf("failed to parse header: %w", err)
		}
	default:
		return fmt.Errorf("invalid schema version %d", schema.SchemaVersion)
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
