#!/bin/bash
set -e

# wit-bindgen go --out-dir ./bindings/go ./bindings/wit
# rm -rf ./bindings/go/go.mod
# rm -rf ./bindings/go/wit_exports.go
# mv ./bindings/go/engine_bindings_api/* ./bindings/go/
# rm -rf ./bindings/go/engine_bindings_api/

# cd ./bindings/generator
# componentize-go.exe -d="../wit" -w=main bindings
# go mod tidy
# cd ../..
# wit-bindgen markdown --out-dir ./docs/api ./bindings/wit

go run src/cli/gen/main.go -l go -s bindings/env.abi -o bindings/out/ -t bindings/target
# go run main.go -l c -s ./env.abi -o ../out/
go run src/cli/gen/main.go -h