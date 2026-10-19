package main

import (
	"fmt"
	"strings"

	. "github.com/StoneTrench/go-mc-clone/bindings/generator/langs"
	. "github.com/StoneTrench/go-mc-clone/bindings/generator/types"
)

type ConfigHostWriter struct {
	Golang        LangBuilder
	NameApiModule string
	NameFnAlloc   string
}

// var Type_AllocFn = TypeFunc{Input: []Param{{"size", TypeP_u32}}, Output: []Param{{"offset", TypeP_WasmPointer}}}

func DeclHostFunc(b *strings.Builder, name string, fn TypeFunc) {
	bind_name := HyphenCaseToPascalCase(name)
	config := &ConfigHostWriter{
		Golang:        GOLANG,
		NameApiModule: "mod",
		NameFnAlloc:   "alloc",
	}
	const fn_name = "fn"
	state := config.Golang
	uses_malloc := false

	wasm_params := []TypeParam{
		{Name: "ctx", Type: TypeIdent{Type: nil, Name: "context.Context"}},
		{Name: config.NameApiModule, Type: TypeIdent{Type: nil, Name: "api.Module"}},
	}

	params := fn.Input
	rets := fn.Output

	for _, p := range rets {
		uses_malloc = uses_malloc || GetABIAllocation(p.Type)
	}

	wasm_params = append(wasm_params, RepParamToABI(params)...)
	for _, p := range RepParamToABI(rets) {
		p.Type = TypeP_usize
		wasm_params = append(wasm_params, p)
	}

	Writef(b, "func %s(", bind_name)
	header_inputs := []TypeParam{
		{Name: fn_name, Type: TypeFunc{
			Input:  append([]TypeParam{{Name: "m", Type: TypeIdent{Type: nil, Name: "api.Module"}}}, params...),
			Output: rets,
		}},
	}
	for i, e := range header_inputs {
		Writes(b, e.Name)
		Writes(b, " ")
		Writes(b, state.StringType(e.Type))
		if i < len(header_inputs)-1 {
			Writes(b, ",")
		}
	}
	Writes(b, ")(w HostBinding)")

	Writes(b, "{\n")
	{
		Writef(b, "return HostBinding{\"%s\",", name)
		{
			Writes(b, state.StringType(TypeFunc{Input: wasm_params}))
			Writes(b, "{\n")
			{
				if uses_malloc {
					Writef(b, "_malloc:=%s.ExportedFunction(\"malloc\");", config.NameApiModule)
					Writef(b, "%s:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));if err!=nil{panic(err)};return uint32(res[0])};", config.NameFnAlloc)
				}
				for _, p := range params {
					WriteGoHostDecode(config, b, p)
				}
				if len(rets) > 0 {
					for i, p := range rets {
						Writes(b, p.Name)
						if i < len(rets)-1 {
							Writes(b, ",")
						}
					}
					Writes(b, ":=")
				}
				Writef(b, "%s(%s,", fn_name, config.NameApiModule)
				for i, p := range params {
					Writes(b, p.Name)
					if i < len(params)-1 {
						Writes(b, ",")
					}
				}
				Writes(b, ");")
				if len(rets) > 0 {
					for _, p := range rets {
						WriteGoHostEncode(config, b, p)
					}
				}
			}
			Writes(b, "},\n")
		}
		Writes(b, "}\n")
	}
	Writes(b, "}\n")
}
func WriteGoHostEncode(config *ConfigHostWriter, b *strings.Builder, p TypeParam) {
	g := config.Golang
	m := config.NameApiModule
	alloc := config.NameFnAlloc
	pr, _ := ParamToABI(p)

	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writef(b, "uWriteString(%s,%s,%s,%s,%s);", m, alloc, name, pr[0].Name, pr[1].Name)
			Writes(b, "\n")
			return
		case TypeP_bool:
			Writef(b, "if %s{", name)
			Writef(b, "uWriteAnyFixed[%s](%s,%s,%s)", g.StringType(pr[0].Type), m, "1", pr[0].Name)
			Writes(b, "}else{")
			Writef(b, "uWriteAnyFixed[%s](%s,%s,%s)", g.StringType(pr[0].Type), m, "0", pr[0].Name)
			Writes(b, "};")
			Writes(b, "\n")
			return
		}

		Writef(b, "uWriteAnyFixed(%s,%s,%s);", m, name, pr[0].Name)
		Writes(b, "\n")
		return
	case TypeArray:
		Writef(b, "uWritePlex(%s,%s,%s,%s);", m, alloc, name, pr[0].Name)
		Writes(b, "\n")
		return
	case TypeSlice:
		Writef(b, "uWriteSlice(%s,%s,%s,%s,%s);", m, alloc, name, pr[0].Name, pr[1].Name)
		Writes(b, "\n")
		return
	case TypeFunc:
		panic("can't pass function through wasm abi")
	case TypePointer:
		Writef(b, "uWritePointer(%s,%s,%s,%s);", m, alloc, name, pr[0].Name)
		Writes(b, "\n")
		return
	case TypeStruct:
		for i, f := range e.Fields {
			Writef(b, "var %s = (%s)(%s.%s);", pr[i].Name, g.StringType(pr[i].Type), name, f.Name)
			Writes(b, "\n")
			WriteGoHostEncode(config, b, pr[i])
		}
		return
	case TypeIdent:
		p.Type = e.Type
		WriteGoHostEncode(config, b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}
func WriteGoHostDecode(config *ConfigHostWriter, b *strings.Builder, p TypeParam) {
	g := config.Golang
	m := config.NameApiModule
	pr, _ := ParamToABI(p)

	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writef(b, "var %s = uReadString(%s, %s, %s);", p.Name, m, pr[0].Name, pr[1].Name)
			Writes(b, "\n")
			return
		case TypeP_bool:
			Writef(b, "var %s = (%s)(%s!=0);", p.Name, g.StringType(p.Type), pr[0].Name)
			Writes(b, "\n")
			return
		}

		Writef(b, "var %s = (%s)(%s);", p.Name, g.StringType(p.Type), pr[0].Name)
		Writes(b, "\n")
		return
	case TypeArray:
		Writef(b, "var %s = uReadPlex[%s](%s, %s);", p.Name, g.StringType(p.Type), m, pr[0].Name)
		Writes(b, "\n")
		return
	case TypeSlice:
		Writef(b, "var %s = uReadSlice[%s](%s, %s, %s);", p.Name, g.StringType(e.Type), m, pr[0].Name, pr[1].Name)
		Writes(b, "\n")
		return
	case TypeFunc:
		panic("can't pass function through wasm abi")
	case TypePointer:
		Writef(b, "var %s = uReadPointer[%s](%s, %s);", p.Name, g.StringType(p.Type), m, pr[0].Name)
		Writes(b, "\n")
		return
	case TypeStruct:
		Writef(b, "var %s %s;", p.Name, g.StringType(p.Type))
		Writes(b, "\n")
		for i, f := range e.Fields {
			WriteGoHostDecode(config, b, pr[i])
			Writef(b, "%s.%s=%s;", p.Name, f.Name, pr[i].Name)
			Writes(b, "\n")
		}
		return
	case TypeIdent:
		p.Type = e.Type
		WriteGoHostDecode(config, b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}
