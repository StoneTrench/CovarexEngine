package langs

import (
	"fmt"
	"strings"

	. "github.com/StoneTrench/go-mc-clone/bindings/generator/types"
)

var CLANG = LangBuilder{
	WriteHeader: c_WriteHeader,
	WriteFooter: nil,

	WriteGuestDecl: c_WriteGuestDecl,

	WriteGuestEncode: c_WriteGuestEncode,
	WriteGuestDecode: c_WriteGuestDecode,

	StringType: c_StringType,
}

func c_WriteHeader(b *strings.Builder, moduleName string) {
	Writes(b, "#include <stdint.h>\n")
	Writes(b, "#include <stdbool.h>\n")
	Writes(b, "#include <stddef.h>\n\n")
	Writes(b, "typedef struct { const char* ptr; uint32_t len; } wasm_string;\n")
	Writes(b, "typedef struct { void* ptr; uint32_t len; } wasm_slice;\n\n")
}
func c_Ident(name string) string { return strings.ReplaceAll(name, "-", "_") }
func c_writeFunc(b *strings.Builder, moduleName string, name string, fn TypeFunc) {
	bind_name := c_Ident(name)
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

	Writes(b, c_StringWasmImport(moduleName, name, api_name, TypeFunc{Input: wasm_params}))
	Writes(b, c_StringFnHeader(bind_name, fn))
	Writes(b, "{\n")
	for _, p := range params {
		c_WriteGuestEncode(b, p)
	}
	for _, p := range wasm_output_vars {
		Writes(b, c_StringDeclare(p))
	}
	Writes(b, c_StringCall(api_name, wasm_input_vars, wasm_output_vars))
	for _, p := range rets {
		c_WriteGuestDecode(b, p)
	}
	if uses_malloc {
		Writes(b, c_StringResetMalloc())
	}
	if len(rets) > 0 {
		Writes(b, c_StringReturn(rets))
	}
	Writes(b, "}\n")
}
func c_WriteGuestDecl(b *strings.Builder, moduleName string, p TypeParam) {
	switch e := p.Type.(type) {
	case TypeFunc:
		c_writeFunc(b, moduleName, p.Name, e)
	case TypeStruct:
		Writef(b, "typedef struct %s {\n", p.Name)
		for _, f := range e.Fields {
			Writef(b, "\t%s %s;\n", c_StringType(f.Type), f.Name)
		}
		Writef(b, "} %s;\n", p.Name)
	default:
		Writef(b, "typedef %s %s;\n", c_StringType(p.Type), p.Name)
	}
}
func c_WriteGuestEncode(b *strings.Builder, p TypeParam) {
	pr, _ := ParamToABI(p)
	if len(pr) == 0 {
		return
	}
	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writes(b, c_StringDeclare(pr[0]))
			Writes(b, c_StringDeclare(pr[1]))
			Writes(b, c_StringAssign(pr[0], "(uint32_t)(uintptr_t)%s.ptr", name))
			Writes(b, c_StringAssign(pr[1], "%s.len", name))
			return
		case TypeP_bool:
			Writes(b, c_StringDeclare(pr[0]))
			Writes(b, c_StringAssign(pr[0], "%s?1:0", name))
			return
		}
		Writes(b, c_StringDeclare(pr[0]))
		Writes(b, c_StringAssign(pr[0], "%s", name))
		return
	case TypeArray:
		Writes(b, c_StringDeclare(pr[0]))
		Writes(b, c_StringAssign(pr[0], "(uint32_t)(uintptr_t)&%s[0]", name))
		return
	case TypeSlice:
		Writes(b, c_StringDeclare(pr[0]))
		Writes(b, c_StringDeclare(pr[1]))
		Writes(b, c_StringAssign(pr[0], "(uint32_t)(uintptr_t)%s.ptr", name))
		Writes(b, c_StringAssign(pr[1], "%s.len", name))
		return
	case TypeFunc:
		Writes(b, c_StringDeclare(pr[0]))
		Writes(b, c_StringAssign(pr[0], "(uint32_t)(uintptr_t)&%s", name))
		return
	case TypePointer:
		Writes(b, c_StringDeclare(pr[0]))
		Writes(b, c_StringAssign(pr[0], "(uint32_t)(uintptr_t)%s", name))
		return
	case TypeStruct:
		Writes(b, c_StringDeclare(pr[0]))
		Writes(b, c_StringAssign(pr[0], "(uint32_t)(uintptr_t)&%s", name))
		return
	case TypeIdent:
		p.Type = e.Type
		c_WriteGuestEncode(b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type: %v", p.Type))
	}
}

func c_WriteGuestDecode(b *strings.Builder, p TypeParam) {
	pr, _ := ParamToABI(p)
	if len(pr) == 0 {
		return
	}
	Writes(b, c_StringDeclare(p))
	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writes(b, c_StringAssign(p, "(wasm_string){ (const char*)(uintptr_t)%s, %s }", pr[0].Name, pr[1].Name))
			return
		case TypeP_bool:
			Writes(b, c_StringAssign(p, "%s!=0", pr[0].Name))
			return
		}
		Writes(b, c_StringAssign(p, "%s", pr[0].Name))
		return
	case TypeArray:
		Writef(b, "%s=*(%s*)(uintptr_t)%s;", name, c_StringType(p.Type), pr[0].Name)
		return
	case TypeSlice:
		Writes(b, c_StringAssign(p, "(wasm_slice){ (void*)(uintptr_t)%s, %s }", pr[0].Name, pr[1].Name))
		return
	case TypeFunc:
		Writef(b, "%s=*(%s*)(uintptr_t)%s;", name, c_StringType(p.Type), pr[0].Name)
		return
	case TypePointer:
		Writef(b, "%s=(%s)(uintptr_t)%s;", name, c_StringType(p.Type), pr[0].Name)
		return
	case TypeStruct:
		Writef(b, "%s=*(%s*)(uintptr_t)%s;", name, c_StringType(p.Type), pr[0].Name)
		return
	case TypeIdent:
		p.Type = e.Type
		c_WriteGuestDecode(b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type: %v", p.Type))
	}
}

func c_StringType(t IType) string {
	switch e := t.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_void:
			return "void"
		case TypeP_u8:
			return "uint8_t"
		case TypeP_u16:
			return "uint16_t"
		case TypeP_u32:
			return "uint32_t"
		case TypeP_u64:
			return "uint64_t"
		case TypeP_i8:
			return "int8_t"
		case TypeP_i16:
			return "int16_t"
		case TypeP_i32:
			return "int32_t"
		case TypeP_i64:
			return "int64_t"
		case TypeP_f32:
			return "float"
		case TypeP_f64:
			return "double"
		case TypeP_string:
			return "wasm_string"
		case TypeP_bool:
			return "bool"
		case TypeP_error:
			return "int32_t"
		}
	case TypeArray:
		return fmt.Sprintf("%s[%d]", c_StringType(e.Type), e.Count)
	case TypeSlice:
		return "wasm_slice"
	case TypeFunc:
		return "void*"
	case TypePointer:
		return fmt.Sprintf("%s*", c_StringType(e.Type))
	case TypeStruct:
		return "void*"
	case TypeIdent:
		return e.Name
	}
	panic(fmt.Sprintf("invalid type: %v", t))
}

func c_StringDeclare(p TypeParam) string {
	return fmt.Sprintf("%s %s;", c_StringType(p.Type), p.Name)
}
func c_StringAssign(dst TypeParam, format string, args ...any) string {
	return fmt.Sprintf("%s=(%s)(%s);", dst.Name, c_StringType(dst.Type), fmt.Sprintf(format, args...))
}
func c_StringCall(name string, args []TypeParam, rets []TypeParam) string {
	str := strings.Builder{}
	for i, e := range args {
		str.WriteString(e.Name)
		if i < len(args)-1 || len(rets) > 0 {
			str.WriteString(",")
		}
	}
	for i, e := range rets {
		Writef(&str, "(uint32_t)(uintptr_t)&%s", e.Name)
		if i < len(rets)-1 {
			str.WriteString(",")
		}
	}
	return fmt.Sprintf("%s(%s);", name, str.String())
}
func c_StringReturn(p []TypeParam) string {
	if len(p) == 0 {
		return "return;"
	}
	if len(p) == 1 {
		return fmt.Sprintf("return %s;", p[0].Name)
	}
	panic("multi-return not supported in C")
}
func c_StringFnHeader(name string, fn TypeFunc) string {
	params := strings.Builder{}
	for i, e := range fn.Input {
		Writef(&params, "%s %s", c_StringType(e.Type), e.Name)
		if i < len(fn.Input)-1 {
			params.WriteString(",")
		}
	}
	if len(fn.Input) == 0 {
		params.WriteString("void")
	}
	rets := c_StringType(TypeP_void)
	if len(fn.Output) == 1 {
		rets = c_StringType(fn.Output[0].Type)
	} else if len(fn.Output) > 1 {
		panic("multi-return not supported in C")
	}
	return fmt.Sprintf("%s %s(%s)", rets, name, params.String())
}
func c_StringWasmImport(module string, name string, bind_name string, fn TypeFunc) string {
	params := strings.Builder{}
	for i, e := range fn.Input {
		Writef(&params, "%s %s", c_StringType(e.Type), e.Name)
		if i < len(fn.Input)-1 {
			params.WriteString(",")
		}
	}
	if len(fn.Input) == 0 {
		params.WriteString("void")
	}
	return fmt.Sprintf(
		"__attribute__((import_module(\"%s\"), import_name(\"%s\")))\nextern void %s(%s);\n",
		module, name, bind_name, params.String(),
	)
}
func c_StringResetMalloc() string {
	return "wasm_reset_malloc();"
}
