package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	model "github.com/StoneTrench/go-mc-clone/bindings/generator/agnostic"
	types "github.com/StoneTrench/go-mc-clone/bindings/generator/model"
	"github.com/StoneTrench/go-mc-clone/bindings/generator/parser"
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

	var path_tmpl_common string
	var path_tmpl_guest string
	var path_tmpl_host string
	var tmpl_counter int
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && strings.HasSuffix(name, ".tmpl") && !strings.HasPrefix(name, ".") {
			switch name[0] {
			case 'c':
				path_tmpl_common = filepath.Join(tmplDir, e.Name())
			case 'g':
				path_tmpl_guest = filepath.Join(tmplDir, e.Name())
			case 'h':
				path_tmpl_host = filepath.Join(tmplDir, e.Name())
			default:
				panic(fmt.Errorf("invalid template: valid names are, in lowercase 'common.tmpl', 'guest.tmpl', 'host.tmpl' (only the first letter, and the suffix is checked. '.' prefix makes a file hidden)"))
			}
			tmpl_counter++
		}
	}
	if len(path_tmpl_guest) == 0 && len(path_tmpl_host) == 0 {
		panic(fmt.Errorf("there has to be a at least a host.tmpl or guest.tmpl"))
	}
	if tmpl_counter == 0 {
		panic(fmt.Errorf("target %q has no .tmpl files", *lang))
	}

	funcs_map := template.FuncMap{
		"GetTypeBody": func(t types.TypeRef) types.IType {
			body, err := coll.Body(t)
			if err != nil {
				return types.TypeP_Invalid
			}
			return body
		},
	}

	var tmpl_guest *template.Template
	var tmpl_host *template.Template

	if len(path_tmpl_guest) > 0 {
		tmpl_guest = template.New("GuestTemplate").Funcs(funcs_map)

		if len(path_tmpl_common) > 0 {
			tmpl_guest, err = tmpl_guest.ParseFiles(path_tmpl_common)
			if err != nil {
				panic(err)
			}
		}
		tmpl_guest, err = tmpl_guest.ParseFiles(path_tmpl_guest)
		if err != nil {
			panic(err)
		}
	}
	if len(path_tmpl_host) > 0 {
		tmpl_host = template.New("HostTemplate").Funcs(funcs_map)

		if len(path_tmpl_common) > 0 {
			tmpl_host, err = tmpl_host.ParseFiles(path_tmpl_common)
			if err != nil {
				panic(err)
			}
		}
		tmpl_host, err = tmpl_host.ParseFiles(path_tmpl_host)
		if err != nil {
			panic(err)
		}
	}

	for _, m := range modules {
		mod, err := model.BuildAgnosticModule(coll, m.Name, m.Types)
		if err != nil {
			panic(err)
		}
		var dest string

		if tmpl_guest != nil {
			dest = filepath.Join(*out, *lang, "guest", m.Name+"."+*lang)
			render(tmpl_guest, "main", mod, dest)
		}

		if tmpl_host != nil {
			dest = filepath.Join(*out, *lang, "host", m.Name+"."+*lang)
			render(tmpl_host, "main", mod, dest)
		}
	}
	// for id, value := range coll.Types {
	// 	fmt.Printf("type %%%d : %s\n", id+1, types.String(value))
	// }
}

func render(t *template.Template, name string, data *model.Module, out string) {
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// fmt.Print(f, data.StringPadding(0, "  "))
	if err := t.ExecuteTemplate(f, name, data); err != nil {
		panic(err)
	}
}
