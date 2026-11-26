# CovarexEngine

This previously was a game, but the game and the engine has been split into two repositories.

## Features

Platform:

- Embedded read-only filesystem wrapper.
- Pull version info from the github repository during build and embed that into the project.

Core:

- Event system: Signals that you can subscribe to.
- Logger system that uses slog in the background.
- Math: `Vec2`, `Vec3`, `Quat`, `Mat4x4` are all generic.

Resources:

- Configuration module.
- Language translation module.
- Resource registry module.

System:

- Entity component system (without the system for now, still working on that).
- Mod dependency issue analyzer.
- WASM ABI binding generator compatible host module builder for Wazero.
- UDP based protocol with reliable and split packet support.
- Worker pool for asyncronous job execution.

## Planned features

- [ ] Custom rendering pipeline
- [ ] Runtime texture/image generation
- [ ] Networking

## Usage

### Build

To build, run:

```sh
go run github.com/StoneTrench/CovarexEngine/cli/build
```

Release:

```sh
go run github.com/StoneTrench/CovarexEngine/cli/build --os linux --arch amd64 --release
```

Flags: `--os`, `--arch`, `--name`, `--out`, `--dist`, `--release`.
Output: `_out/`.
Release tarballs: `_dist/`.

### ABI Bindings

Refer to [this](./cli/gen/README.md).

```sh
go run github.com/StoneTrench/CovarexEngine/cli/gen -l go -s bindings/env.abi -o bindings/
```

Generates `bindings/go/{host,guest}/` from `bindings/env.abi`.