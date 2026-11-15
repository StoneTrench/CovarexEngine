package modding

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero"
)

/*
HostModuleBuilder registers functions inside the runtime.

When a GuestInstance calls a host function, the go function is called.

There's not really a need to make more than one Runtime.
*/
type HostModuleBuilder struct {
	rt     wazero.Runtime
	module wazero.HostModuleBuilder
}

type HostBinding = struct {
	Name   string
	Method any
}

func (g *HostModuleBuilder) Init(ctx context.Context, rt wazero.Runtime, name string) *HostModuleBuilder {
	g.rt = rt
	g.module = g.rt.NewHostModuleBuilder(name)
	return g
}
func (g *HostModuleBuilder) Add(b HostBinding) *HostModuleBuilder {
	g.module.NewFunctionBuilder().WithFunc(b.Method).Export(b.Name)
	return g
}
func (g *HostModuleBuilder) Instantiate(ctx context.Context) error {
	if g.module == nil {
		return fmt.Errorf("host module is already compiled")
	}
	comp, err := g.module.Compile(ctx)
	if err != nil {
		return fmt.Errorf("failed to compile host module: %w", err)
	}
	defer comp.Close(ctx)
	_, err = g.rt.InstantiateModule(ctx, comp, wazero.NewModuleConfig())
	if err != nil {
		return fmt.Errorf("failed to instantiate host module: %w", err)
	}

	return nil
}
