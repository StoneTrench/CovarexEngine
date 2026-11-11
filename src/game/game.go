package game

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"path"

	"github.com/StoneTrench/go-mc-clone/bindings/host"
	"github.com/StoneTrench/go-mc-clone/src/engine"
	"github.com/StoneTrench/go-mc-clone/src/engine/core/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/platform"
	"github.com/StoneTrench/go-mc-clone/src/engine/resources"
	"github.com/StoneTrench/go-mc-clone/src/engine/systems/modding"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func Init() error {
	var err error = nil
	err = engine.Init("logs/")
	if err != nil {
		return fmt.Errorf("failed to initialize engine: %w", err)
	}

	log.Info(platform.GetFormattedApplicationLabel())

	ctx := context.Background()
	RT := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.Instantiate(ctx, RT)
	api_builder := new(modding.HostModuleBuilder).Init(ctx, RT, "env")
	load_api(api_builder)
	err = api_builder.Instantiate(ctx)
	if err != nil {
		log.Error(err)
	}
	err = modding.LoadMods(ctx, RT, "mods/")
	if err != nil {
		log.Error(err)
	}
	RegistryBlocks.Freeze()
	_, err = RegistryBlocks.ApplyIdMap([]string{"core:air", "core:stone", "core:dirt", "core:grass", "core:water"})
	if err != nil {
		log.Error(err)
	}

	return nil
}

func load_api(m *modding.HostModuleBuilder) {
	m.Add(host.LogInfo(func(m api.Module, msg string) {
		slog.Info(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))
	m.Add(host.LogWarn(func(m api.Module, msg string) {
		slog.Warn(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))
	m.Add(host.LogError(func(m api.Module, msg string) {
		slog.Error(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))

	m.Add(host.RegAddBlock(func(m api.Module, block host.BlockType) (namespaceid string) {
		return RegistryBlocks.RegisterPanic(m.Name(), block.Name, world.BlockType{
			Color:   color.RGBA(block.Color),
			IsSolid: block.IsSolid,
		})
	}))
	m.Add(host.RegGetBlock(func(m api.Module, namespaceid string) (host.BlockType, host.OptError) {
		b, err := RegistryBlocks.GetByIdentifier(namespaceid)
		wblock := host.BlockType{
			Name:    namespaceid,
			Color:   host.Color(b.Color),
			IsSolid: b.IsSolid,
		}
		if err != nil {
			return wblock, &host.WError{err.Error()}
		}
		return wblock, nil
	}))
	m.Add(host.AppendLang(func(m api.Module, name string) host.OptError {
		mod := modding.GetMod(m)
		if mod == nil {
			return &host.WError{"failed to append language resource: cant find this mod"}
		}
		if len(name) != 5 {
			return &host.WError{"failed to append language resource: name must be 5 characters long"}
		}
		err := resources.LangAppend(path.Join(mod.Folder, "assets", "lang", name+".json"))
		if err != nil {
			return &host.WError{fmt.Errorf("failed to append language resource: %w", err).Error()}
		}
		return nil
	}))
}
