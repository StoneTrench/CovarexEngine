package langs

import (
	"fmt"
	"strings"

	. "github.com/StoneTrench/go-mc-clone/bindings/generator/types"
)

var GOLANG = LangBuilder{
	WriteHeader: go_WriteHeader,
	WriteFooter: nil,

	WriteGuestDecl: go_WriteGuestDecl,

	WriteGuestEncode: go_WriteGuestEncode,
	WriteGuestDecode: go_WriteGuestDecode,

	StringType: go_StringType,
}

type GoLangBuilder struct{}

func go_WriteHeader(b *strings.Builder, moduleName string) {
	Writes(b, "package api\n")
	b.WriteString("import (\n")
	b.WriteString("\"unsafe\"")
	b.WriteString("\n")
	b.WriteString(")\n")
}

func go_writeFunc(b *strings.Builder, moduleName string, name string, fn TypeFunc) {
	bind_name := HyphenCaseToPascalCase(name)
	api_name := fmt.Sprintf("%s_%s", moduleName, bind_name)

	params := fn.Input
	rets := fn.Output

	uses_malloc := false
	wasm_params := []TypeParam{}
	wasm_input_vars := []TypeParam{}
	wasm_output_vars := []TypeParam{}

	for _, p := range rets {
		uses_malloc = uses_malloc || GetABIAllocation(p)
	}

	for _, p := range RepParamToABI(params) {
		wasm_input_vars = append(wasm_input_vars, p)
		wasm_params = append(wasm_params, p)
	}

	for _, e := range RepParamToABI(rets) {
		e.Type = TypeP_usize
		wasm_output_vars = append(wasm_output_vars, e)
		wasm_params = append(wasm_params, e)
	}

	Writes(
		b,
		go_StringWasmImport(
			moduleName, name,
			api_name,
			TypeFunc{
				Input: wasm_params,
			},
		),
	)
	Writes(b, go_StringFnHeader(bind_name, fn))
	Writes(b, "{\n")
	{
		Writes(b, "// Guest Encode\n")
		for _, p := range params {
			go_WriteGuestEncode(b, p)
			Writes(b, "\n")
		}
		Writes(b, "// Declare Outputs\n")
		if len(rets) > 0 {
			for _, p := range wasm_output_vars {
				Writes(b, go_StringDeclare(p))
				Writes(b, "\n")
			}
		}
		Writes(b, "// Call\n")
		Writes(b, go_StringCall(api_name, wasm_input_vars, wasm_output_vars))
		Writes(b, "\n")
		Writes(b, "// Guest Decode\n")
		if len(rets) > 0 {
			for _, p := range rets {
				go_WriteGuestDecode(b, p)
				Writes(b, "\n")
			}
		}
		Writes(b, "// Finish\n")
		if uses_malloc {
			Writes(b, go_StringResetMalloc())
			Writes(b, "\n")
		}
		if len(rets) > 0 {
			Writes(b, go_StringReturn(rets))
			Writes(b, "\n")
		}
	}
	Writes(b, "}\n")
}

func go_WriteGuestDecl(b *strings.Builder, moduleName string, p TypeParam) {
	switch e := p.Type.(type) {
	case TypeFunc:
		go_writeFunc(b, moduleName, p.Name, e)
	default:
		Writef(b, "type %s %s\n", p.Name, go_StringType(p.Type))
	}
}
func go_WriteGuestEncode(b *strings.Builder, p TypeParam) {
	pr, _ := ParamToABI(p)

	if len(pr) == 0 {
		return
	}

	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writef(b, "var %s = (%s)(uintptr(unsafe.Pointer(unsafe.StringData(%s))));", pr[0].Name, go_StringType(pr[0].Type), name)
			Writes(b, "\n")
			Writef(b, "var %s = (%s)(len(%s));", pr[1].Name, go_StringType(pr[1].Type), name)
			Writes(b, "\n")
			return
		case TypeP_bool:
			Writes(b, go_StringDeclare(pr[0]))
			Writes(b, "\n")
			Writef(b, "if %s{", name)
			Writef(b, "%s = 1;", pr[0].Name)
			Writef(b, "}else{")
			Writef(b, "%s = 0;", pr[0].Name)
			Writef(b, "};")
			Writes(b, "\n")
			return
		}

		Writef(b, "var %s = (%s)(%s);", pr[0].Name, go_StringType(pr[0].Type), name)
		return
	case TypeArray:
		Writes(b, go_StringDeclare(pr[0]))
		Writes(b, "\n")
		Writef(b, "if len(%s)>0{", name)
		Writef(b, "%s = (%s)(uintptr(unsafe.Pointer(&%s[0])));", pr[0].Name, go_StringType(pr[0].Type), name)
		Writef(b, "};")
		Writes(b, "\n")
		return
	case TypeSlice:
		Writes(b, go_StringDeclare(pr[0]))
		Writes(b, "\n")
		Writef(b, "var %s = (%s)(len(%s));", pr[1].Name, go_StringType(pr[1].Type), name)
		Writef(b, "if %s>0{", pr[0].Name)
		Writef(b, "%s = (%s)(uintptr(unsafe.Pointer(&%s[0])));", pr[0].Name, go_StringType(pr[0].Type), name)
		Writef(b, "};")
		Writes(b, "\n")
		return
	case TypeFunc:
		Writef(b, "var %s = (%s)(uintptr(unsafe.Pointer(&%s)));", pr[0].Name, go_StringType(pr[0].Type), name)
		Writes(b, "\n")
		return
	case TypePointer:
		Writef(b, "var %s = (%s)(uintptr(unsafe.Pointer(%s)));", pr[0].Name, go_StringType(pr[0].Type), name)
		Writes(b, "\n")
		return
	case TypeStruct:
		for i, f := range e.Fields {
			Writef(b, "var %s = (%s)(%s.%s);", pr[i].Name, go_StringType(pr[i].Type), name, f.Name)
			go_WriteGuestEncode(b, pr[i])
		}
		return
	case TypeIdent:
		p.Type = e.Type
		go_WriteGuestEncode(b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}
func go_WriteGuestDecode(b *strings.Builder, p TypeParam) {
	pr, _ := ParamToABI(p)

	if len(pr) == 0 {
		return
	}

	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writef(b, "var %s = (%s)(unsafe.String((*byte)(unsafe.Pointer(uintptr(%s))),int(%s)));", p.Name, go_StringType(p.Type), pr[0].Name, pr[1].Name)
			Writes(b, "\n")
			return
		case TypeP_bool:
			Writef(b, "var %s = (%s)(%s!=0);", p.Name, go_StringType(p.Type), pr[0].Name)
			Writes(b, "\n")
			return
		}

		Writef(b, "var %s = (%s)(%s);", p.Name, go_StringType(p.Type), pr[0].Name)
		Writes(b, "\n")
		return
	case TypeArray:
		Writef(b, "var %s=*(*%s)(unsafe.Pointer(uintptr(%s)));", name, go_StringType(p.Type), pr[0].Name)
		Writes(b, "\n")
		return
	case TypeSlice:
		Writef(b, "var %s=*(*%s)(unsafe.Slice((*%s)(unsafe.Pointer(uintptr(%s)),int(%s))));", name, go_StringType(p.Type), go_StringType(e.Type), pr[0].Name, pr[1].Name)
		Writes(b, "\n")
		return
	case TypeFunc:
		Writef(b, "var %s=*(*%s)(unsafe.Pointer(uintptr(%s)));", name, go_StringType(p.Type), pr[0].Name)
		Writes(b, "\n")
		return
	case TypePointer:
		Writef(b, "var %s=*(*%s)(unsafe.Pointer(uintptr(%s)));", name, go_StringType(p.Type), pr[0].Name)
		Writes(b, "\n")
		return
	case TypeStruct:
		Writef(b, "var %s %s;", p.Name, go_StringType(p.Type))
		Writes(b, "\n")
		for i, f := range e.Fields {
			go_WriteGuestDecode(b, pr[i])
			Writef(b, "%s.%s=%s;", name, f.Name, pr[i].Name)
			Writes(b, "\n")
		}
		return
	case TypeIdent:
		p.Type = e.Type
		go_WriteGuestDecode(b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}

func go_StringType(t IType) string {
	switch e := t.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_void:
			return "struct{}"
		case TypeP_u8:
			return "uint8"
		case TypeP_u16:
			return "uint16"
		case TypeP_u32:
			return "uint32"
		case TypeP_u64:
			return "uint64"
		case TypeP_i8:
			return "int8"
		case TypeP_i16:
			return "int16"
		case TypeP_i32:
			return "int32"
		case TypeP_i64:
			return "int64"
		case TypeP_f32:
			return "float32"
		case TypeP_f64:
			return "float64"
		case TypeP_string:
			return "string"
		case TypeP_bool:
			return "bool"
		case TypeP_error:
			return "error"
		}
	case TypeArray:
		return fmt.Sprintf("[%d]%s", e.Count, go_StringType(e.Type))
	case TypeSlice:
		return fmt.Sprintf("[]%s", go_StringType(e.Type))
	case TypeFunc:
		input := make([]string, len(e.Input))
		output := make([]string, len(e.Output))

		for i, val := range e.Input {
			input[i] = fmt.Sprintf("%s %s", val.Name, go_StringType(val.Type))
		}
		for i, val := range e.Output {
			output[i] = fmt.Sprintf("%s %s", val.Name, go_StringType(val.Type))
		}

		return fmt.Sprintf("func(%s)(%s)", strings.Join(input, ","), strings.Join(output, ","))
	case TypePointer:
		return fmt.Sprintf("*%s", go_StringType(e.Type))
	case TypeStruct:
		fields := make([]string, len(e.Fields))
		for i, val := range e.Fields {
			fields[i] = fmt.Sprintf("%s %s", val.Name, go_StringType(val.Type))
		}
		return fmt.Sprintf("struct{%s}", strings.Join(fields, ";"))
	case TypeIdent:
		return e.Name
	}
	panic(fmt.Sprintf("invalid type type: %v", t))
}

func go_StringDeclare(p TypeParam) string {
	return fmt.Sprintf("var %s %s;", p.Name, go_StringType(p.Type))
}
func go_StringCall(name string, args []TypeParam, rets []TypeParam) string {
	str := strings.Builder{}
	for i, e := range args {
		str.WriteString(e.Name)
		if i < len(args)-1 {
			str.WriteString(",")
		}
	}
	if len(rets) > 0 {
		str.WriteString(",")
		for i, e := range rets {
			Writef(&str, "(%s)(uintptr(unsafe.Pointer(&%s)))", go_StringType(e.Type), e.Name)
			if i < len(rets)-1 {
				str.WriteString(",")
			}
		}
	}
	return fmt.Sprintf("%s(%s);", name, str.String())
}
func go_StringReturn(p []TypeParam) string {
	args := strings.Builder{}
	for i, e := range p {
		args.WriteString(e.Name)
		if i < len(p)-1 {
			args.WriteString(",")
		}
	}
	return fmt.Sprintf("return %s;", args.String())
}
func go_StringFnHeader(name string, fn TypeFunc) string {
	params := strings.Builder{}
	for i, e := range fn.Input {
		Writef(&params, "%s %s", e.Name, go_StringType(e.Type))
		if i < len(fn.Input)-1 {
			params.WriteString(",")
		}
	}
	rets := strings.Builder{}
	if len(fn.Output) > 0 {
		rets.WriteString("(")
		for i, e := range fn.Output {
			Writef(&rets, "%s", go_StringType(e.Type))
			if i < len(fn.Output)-1 {
				rets.WriteString(",")
			}
		}
		rets.WriteString(")")
	}
	return fmt.Sprintf("func %s(%s)%s", name, params.String(), rets.String())
}
func go_StringWasmImport(module string, name string, bind_name string, fn TypeFunc) string {
	params := strings.Builder{}
	for i, e := range fn.Input {
		Writef(&params, "%s %s", e.Name, go_StringType(e.Type))
		if i < len(fn.Input)-1 {
			params.WriteString(",")
		}
	}
	rets := strings.Builder{}
	if len(fn.Output) > 0 {
		rets.WriteString("(")
		for i, e := range fn.Output {
			Writef(&rets, "%s", go_StringType(e.Type))
			if i < len(fn.Output)-1 {
				rets.WriteString(",")
			}
		}
		rets.WriteString(")")
	}
	return fmt.Sprintf("//go:wasmimport %s %s\nfunc %s(%s)%s;\n", module, name, bind_name, params.String(), rets.String())
}
func go_StringResetMalloc() string {
	return "uResetMalloc();"
}
