package game

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"path"

	host "github.com/StoneTrench/CovarexEngine/bindings/out/go/host"
	"github.com/StoneTrench/CovarexEngine/src/engine"
	"github.com/StoneTrench/CovarexEngine/src/engine/core/log"
	"github.com/StoneTrench/CovarexEngine/src/engine/platform"
	"github.com/StoneTrench/CovarexEngine/src/engine/resources"
	"github.com/StoneTrench/CovarexEngine/src/engine/systems/modding"
	"github.com/StoneTrench/CovarexEngine/src/game/world"
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
	_, err = RegistryBlocks.ApplyIdMap([]string{"me.stonetrench.core:air", "me.stonetrench.core:stone", "me.stonetrench.core:dirt", "me.stonetrench.core:grass", "me.stonetrench.core:water"})
	if err != nil {
		log.Error(err)
	}

	return nil
}

func load_api(m *modding.HostModuleBuilder) {
	m.Add(host.HostAPILogInfo(func(m api.Module, msg string) {
		slog.Info(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))
	m.Add(host.HostAPILogWarn(func(m api.Module, msg string) {
		slog.Warn(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))
	m.Add(host.HostAPILogError(func(m api.Module, msg string) {
		slog.Error(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))

	m.Add(host.HostAPIRegAddBlock(func(m api.Module, block host.BlockType) (namespaceid string) {
		IsSolid := false
		if block.Kind == host.E_Solid {
			IsSolid = true
		}
		return RegistryBlocks.RegisterPanic(m.Name(), block.Name, world.BlockType{
			Color:   color.RGBA(block.Color),
			IsSolid: IsSolid,
		})
	}))
	m.Add(host.HostAPIRegGetBlock(func(m api.Module, namespaceid string) (host.BlockType, string) {
		b, err := RegistryBlocks.GetByIdentifier(namespaceid)
		wblock := host.BlockType{}
		if err != nil {
			return wblock, err.Error()
		}

		wblock.Name = namespaceid
		wblock.Color = host.ColorU8(b.Color)
		var Kind host.BlockKind
		if b.IsSolid {
			Kind = host.E_Solid
		} else {
			Kind = host.E_Air
		}
		wblock.Kind = Kind
		return wblock, ""
	}))
	m.Add(host.HostAPIAppendLang(func(m api.Module, name string) string {
		mod := modding.GetMod(m)
		if mod == nil {
			return "failed to append language resource: cant find this mod"
		}
		if len(name) != 5 {
			return "failed to append language resource: name must be 5 characters long"
		}
		err := resources.LangAppend(path.Join(mod.Folder, "assets", "lang", name+".json"))
		if err != nil {
			return fmt.Errorf("failed to append language resource: %w", err).Error()
		}
		return ""
	}))
}
