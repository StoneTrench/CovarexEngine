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

type Input struct {
	Src  string
	Lang string
	Out  string
}

func main() {
	var (
		src  = flag.String("s", "env.abi", "input `.abi` source file")
		lang = flag.String("l", "go", "target language directory `target/<lang>/*.tmpl`")
		out  = flag.String("o", "out", "output directory")
	)
	flag.Parse()

	err := cli(&Input{
		Src:  *src,
		Lang: *lang,
		Out:  *out,
	})

	if err != nil {
		panic(fmt.Errorf("%w\n%s", err, flag.ErrHelp))
	}
}

func cli(input *Input) error {
	t_coll := new(types.TypeCollection).Init()
	funcs_map := template.FuncMap{
		"GetTypeBody": func(t types.TypeRef) types.IType {
			body, err := t_coll.Body(t)
			if err != nil {
				return types.TypeP_Invalid
			}
			return body
		},
	}

	modules, err := parser.ParseSource(t_coll, input.Src)
	if err != nil {
		return err
	}
	path_tmpl_common, path_tmpl_guest, path_tmpl_host, err := discover_tmpl(input.Lang)
	if err != nil {
		return fmt.Errorf("no template avaliable for target %q: %w", input.Lang, err)
	}
	tmpl_guest, tmpl_host, err := make_templates(path_tmpl_common, path_tmpl_guest, path_tmpl_host, funcs_map)
	if err != nil {
		return fmt.Errorf("failed to parse go templates for target %q: %w", input.Lang, err)
	}
	err = emit_modules(modules, t_coll, tmpl_guest, tmpl_host, input.Out, input.Lang)
	if err != nil {
		return fmt.Errorf("failed to emit modules for target %q: %w", input.Lang, err)
	}

	return nil
}

func make_templates(path_tmpl_common, path_tmpl_guest, path_tmpl_host string, funcs_map template.FuncMap) (tmpl_guest, tmpl_host *template.Template, err error) {
	if len(path_tmpl_guest) > 0 {
		tmpl_guest = template.New("GuestTemplate").Funcs(funcs_map)

		if len(path_tmpl_common) > 0 {
			tmpl_guest, err = tmpl_guest.ParseFiles(path_tmpl_common)
			if err != nil {
				return nil, nil, err
			}
		}
		tmpl_guest, err = tmpl_guest.ParseFiles(path_tmpl_guest)
		if err != nil {
			return nil, nil, err
		}
	}
	if len(path_tmpl_host) > 0 {
		tmpl_host = template.New("HostTemplate").Funcs(funcs_map)

		if len(path_tmpl_common) > 0 {
			tmpl_host, err = tmpl_host.ParseFiles(path_tmpl_common)
			if err != nil {
				return nil, nil, err
			}
		}
		tmpl_host, err = tmpl_host.ParseFiles(path_tmpl_host)
		if err != nil {
			return nil, nil, err
		}
	}
	return tmpl_guest, tmpl_host, nil
}
func emit_modules(modules []parser.Module, t_coll *types.TypeCollection, tmpl_guest, tmpl_host *template.Template, out, lang string) error {
	for _, m := range modules {
		mod, err := model.BuildAgnosticModule(t_coll, m.Name, m.Types)
		if err != nil {
			return err
		}
		var dest string

		if tmpl_guest != nil {
			dest = filepath.Join(out, lang, "guest", m.Name+"."+lang)
			err := emit_template(tmpl_guest, "main", mod, dest)
			if err != nil {
				return err
			}
		}

		if tmpl_host != nil {
			dest = filepath.Join(out, lang, "host", m.Name+"."+lang)
			err := emit_template(tmpl_host, "main", mod, dest)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func emit_template(t *template.Template, name string, data *model.Module, out string) error {
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	// fmt.Print(f, data.StringPadding(0, "  "))
	if err := t.ExecuteTemplate(f, name, data); err != nil {
		return err
	}
	return nil
}
func discover_tmpl(lang string) (path_tmpl_common, path_tmpl_guest, path_tmpl_host string, err error) {
	tmplDir := filepath.Join("target", lang)
	entries, err := os.ReadDir(tmplDir)
	if err != nil {
		return path_tmpl_common, path_tmpl_guest, path_tmpl_host,
			err
	}

	var tmpl_counter int
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && strings.HasSuffix(name, ".tmpl") && !strings.HasPrefix(name, ".") {
			switch name {
			case "common.tmpl":
				path_tmpl_common = filepath.Join(tmplDir, e.Name())
			case "guest.tmpl":
				path_tmpl_guest = filepath.Join(tmplDir, e.Name())
			case "host.tmpl":
				path_tmpl_host = filepath.Join(tmplDir, e.Name())
			default:
				return path_tmpl_common, path_tmpl_guest, path_tmpl_host,
					fmt.Errorf("valid names are, in lowercase 'common.tmpl', 'guest.tmpl', 'host.tmpl' ('.' prefix makes a file hidden)")
			}
			tmpl_counter++
		}
		if tmpl_counter == 3 {
			break
		}
	}
	if len(path_tmpl_guest) == 0 && len(path_tmpl_host) == 0 {
		return path_tmpl_common, path_tmpl_guest, path_tmpl_host,
			fmt.Errorf("there has to be a at least a host.tmpl or guest.tmpl")
	}
	if tmpl_counter == 0 {
		return path_tmpl_common, path_tmpl_guest, path_tmpl_host,
			fmt.Errorf("target %q has no .tmpl files", lang)
	}
	return path_tmpl_common, path_tmpl_guest, path_tmpl_host,
		nil
}
