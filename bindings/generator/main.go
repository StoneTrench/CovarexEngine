package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/StoneTrench/go-mc-clone/bindings/generator/model"
	"github.com/StoneTrench/go-mc-clone/bindings/generator/parser"
	"github.com/StoneTrench/go-mc-clone/bindings/generator/types"
)

func main() {
	var (
		src  = flag.String("src", "/env.abi", "input .abi file")
		lang = flag.String("lang", "go", "target language directory target/<lang>/*.tmpl")
		out  = flag.String("out", "generated/", "output directory")
	)
	flag.Parse()

	coll := new(types.TypeCollection).Init()
	modules, err := parser.ParseSource(coll, *src)
	if err != nil {
		panic(err)
	}

	tmplDir := filepath.Join("target", *lang)
	entries, err := os.ReadDir(tmplDir)
	if err != nil {
		panic(fmt.Errorf("target %q: %w", *lang, err))
	}

	var tmplFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".tmpl") && !strings.HasPrefix(e.Name(), "_") {
			tmplFiles = append(tmplFiles, filepath.Join(tmplDir, e.Name()))
		}
	}
	if len(tmplFiles) == 0 {
		panic(fmt.Errorf("target %q has no .tmpl files", *lang))
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"GetTypeBody": func(t types.TypeRef) types.IType {
			body, err := coll.Body(t)
			if err != nil {
				return types.TypeP_Invalid
			}
			return body
		},
	}).ParseFiles(tmplFiles...)
	if err != nil {
		panic(err)
	}

	outDir := filepath.Join(*out, *lang)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}

	for _, m := range modules {
		mod, err := model.Build(coll, m.Name, m.Decl)
		if err != nil {
			panic(err)
		}
		dest := filepath.Join(outDir, m.Name+"."+*lang)
		render(tmpl, "main", mod, dest)
	}
}

func render(t *template.Template, name string, data any, out string) {
	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := t.ExecuteTemplate(f, name, data); err != nil {
		panic(err)
	}
}
