package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/StoneTrench/go-mc-clone/bindings/generator/langs"
	. "github.com/StoneTrench/go-mc-clone/bindings/generator/types"
)

func main() {
	LANG := GOLANG

	language := "go"

	root := "bindings/template"
	entries, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(root, e.Name())
		name, decl, err := P_parse(path)
		if err != nil {
			panic(err)
		}
		bguest := &strings.Builder{}
		bhost := &strings.Builder{}

		bhost.WriteString("package host\n")
		bhost.WriteString("import \"github.com/tetratelabs/wazero/api\"\n")
		bhost.WriteString("import \"context\"\n")
		LANG.WriteHeader(bguest, name)

		for _, d := range decl {
			switch e := d.Type.(type) {
			case TypeFunc:
				DeclHostFunc(bhost, d.Name, e)
			default:
				GOLANG.WriteGuestDecl(bhost, name, d)
			}
			LANG.WriteGuestDecl(bguest, name, d)
		}

		if LANG.WriteFooter != nil {
			LANG.WriteFooter(bguest, name)
		}

		file, err := os.Create(fmt.Sprintf("bindings/%[2]s/%[1]s.%[2]s", name, language))
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(bguest.String())
		if err != nil {
			panic(err)
		}
		file.Close()

		file, err = os.Create(fmt.Sprintf("bindings/host/%s.go", name))
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(bhost.String())
		if err != nil {
			panic(err)
		}
		file.Close()
	}

}
