# CovarexEngine

This previously was a game, but the game and the engine has been split into two repositories.

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

### Bindings

Refer to [this](./cli/gen/README.md).

```sh
go run github.com/StoneTrench/CovarexEngine/cli/gen -l go -s bindings/env.abi -o bindings/
```

Generates `bindings/go/{host,guest}/` from `bindings/env.abi`.