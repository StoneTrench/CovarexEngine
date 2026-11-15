package modding

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/Masterminds/semver/v3"
	"github.com/StoneTrench/CovarexEngine/src/engine/platform/filesystem"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type GuestModuleHeaderJson struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`

	Homepage string `json:"homepage"`

	ApiVersion        int               `json:"api_version"`
	Dependencies      map[string]string `json:"dependencies"`
	Incompatibilities map[string]string `json:"incompatibilities"`

	Environment map[string]string `json:"environment"`
}

func read_mod_header(folder string, out *GuestModuleHeaderJson) error {
	file_mod_json, err := filesystem.ReadFile(path.Join(folder, "mod.json"))
	if err != nil {
		return fmt.Errorf("failed to read mod.json: %w", err)
	}

	var head GuestModuleHeaderJson
	err = json.Unmarshal(file_mod_json, &head)
	if err != nil {
		return fmt.Errorf("failed to parse mod.json: %w", err)
	}
	*out = head

	return nil
}
func init_mod_wasm(ctx context.Context, file_path string, name string, rt wazero.Runtime) (api.Module, error) {
	bin, err := filesystem.ReadFile(file_path)
	if err != nil {
		return nil, fmt.Errorf("failed to open wasm file: %w", err)
	}

	comp, err := rt.CompileModule(ctx, bin)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}
	defer comp.Close(ctx)

	cnf := wazero.NewModuleConfig().
		WithStartFunctions("_initialize"). // this has to be called here manually for some reason
		WithName(name)

	inst, err := rt.InstantiateModule(ctx, comp, cnf)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	return inst, nil
}

type GuestModule struct {
	header GuestModuleHeaderJson

	Folder string

	Version      semver.Version
	Dependencies map[string]semver.Constraints

	Instance api.Module
}

func (m GuestModule) GetName() string {
	return m.header.Name
}
func (m GuestModule) GetVersion() *semver.Version {
	return &m.Version
}
func (m GuestModule) GetDependencies() map[string]semver.Constraints {
	return m.Dependencies
}
func (m GuestModule) GetAuthor() string {
	return m.header.Author
}
func (m GuestModule) GetDescription() string {
	return m.header.Description
}
func (m GuestModule) GetHomepage() string {
	return m.header.Homepage
}
func (m GuestModule) IsApiVersionValid(api_version int) error {
	if m.header.ApiVersion != api_version {
		return fmt.Errorf("API version of engine and mod do not match (%d != %d)", api_version, m.header.ApiVersion)
	}
	return nil
}

func (m *GuestModule) Init(folder string) *GuestModule {
	m.Folder = folder
	return m
}
func (m *GuestModule) Close(ctx context.Context) {
	if m.Instance != nil {
		m.Instance.Close(ctx)
	}
}
func (m *GuestModule) LoadHeader() error {
	head := &m.header

	err := read_mod_header(m.Folder, head)
	if err != nil {
		return err
	}

	vers, err := semver.NewVersion(head.Version)
	if err != nil {
		return fmt.Errorf("failed to parse mod version: %w", err)
	}
	m.Version = *vers

	m.Dependencies = make(map[string]semver.Constraints)
	for k, v := range head.Dependencies {
		vers, err := semver.NewConstraint(v)
		if err != nil {
			return fmt.Errorf("failed to parse mod dependency (%s: %s): %w", k, v, err)
		}
		m.Dependencies[k] = *vers
	}

	return nil
}
func (m *GuestModule) LoadWasm(ctx context.Context, rt wazero.Runtime) error {
	inst, err := init_mod_wasm(ctx, path.Join(m.Folder, "main.wasm"), m.GetName(), rt)
	if err != nil {
		return fmt.Errorf("failed to load mod (%s %s): %w", m.GetName(), m.Version.String(), err)
	}
	m.Instance = inst
	return nil
}
