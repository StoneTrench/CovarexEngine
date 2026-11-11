package main

import (
	"fmt"

	. "github.com/StoneTrench/go-mc-clone/bindings/generator/types"
	"github.com/StoneTrench/go-mc-clone/src/engine/platform/filesystem"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type P_Module struct {
	Name  string    `parser:"'module' @Ident"`
	Types []*P_Decl `parser:"@@*"`
}

type P_Decl struct {
	Name string  `parser:"'decl' @Ident"`
	Type *P_Type `parser:"@@"`
}
type P_Param struct {
	Name string  `parser:"@Ident"`
	Give bool    `parser:"'give'?"`
	Type *P_Type `parser:"@@"`
}
type P_Field struct {
	Name string  `parser:"@Ident"`
	Type *P_Type `parser:"@@"`
}

type P_Type struct {
	Struct  *P_Struct  `parser:"@@"`
	Array   *P_Array   `parser:"| @@"`
	Slice   *P_Slice   `parser:"| @@"`
	Pointer *P_Pointer `parser:"| @@"`
	Func    *P_Func    `parser:"| @@"`
	Ident   *P_Ident   `parser:"| @@"`
}

type P_Ident struct {
	Name string `parser:"@Ident"`
}
type P_Struct struct {
	Fields []*P_Field `parser:"'struct' '{' @@* '}'"`
}
type P_Array struct {
	Count int64   `parser:"'[' @Int ']'"`
	Type  *P_Type `parser:"@@"`
}
type P_Slice struct {
	Type *P_Type `parser:"'[' ']' @@"`
}
type P_Pointer struct {
	_    struct{} `parser:"'*'"`
	Mut  bool     `parser:"@'mut'?"`
	Type *P_Type  `parser:"@@"`
}
type P_Func struct {
	Input  []*P_Param `parser:"'(' (@@ (',' @@)*)? ')'"`
	Output []*P_Param `parser:"'(' (@@ (',' @@)*)? ')'"`
}

var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comment", Pattern: `//[^\n]*`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_\-]*`},
	{Name: "Punct", Pattern: `[{}()\[\],*]`},
	{Name: "Whitespace", Pattern: `[ \t\r\n;]+`},
	{Name: "Int", Pattern: `[0-9][0-9]*`},
})

var parser = participle.MustBuild[P_Module](
	participle.Lexer(lexerDef),
	participle.Elide("Comment", "Whitespace"),
)

func P_parse(name string) (string, []TypeField, error) {
	bits, err := filesystem.ReadFile(name)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse file: %w", err)
	}
	module, err := parser.ParseBytes(name, bits)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse file: %w", err)
	}
	name, decl := P_collect_decl(module)
	return name, decl, nil
}

func P_collect_decl(t *P_Module) (string, []TypeField) {
	decl := make(map[string]int)
	res_decl := []TypeField{}

	LUF := func(name string) IType {
		i, exists := decl[name]
		if !exists {
			return nil
		}
		return res_decl[i].Type
	}

	for _, t := range t.Types {
		res_t := P_convert_type(LUF, t.Type)
		decl[t.Name] = len(res_decl)
		res_decl = append(res_decl, TypeField{Name: t.Name, Type: res_t})
	}

	return t.Name, res_decl
}

func P_convert_type(decl func(name string) IType, t *P_Type) IType {
	if t.Ident != nil {
		return P_convert_ident(decl, t.Ident)
	} else if t.Struct != nil {
		return TypeStruct{
			Fields: P_convert_field(decl, t.Struct.Fields),
		}
	} else if t.Array != nil {
		return TypeArray{
			Count: int(t.Array.Count),
			Type:  P_convert_type(decl, t.Array.Type),
		}
	} else if t.Slice != nil {
		return TypeSlice{
			Type: P_convert_type(decl, t.Slice.Type),
		}
	} else if t.Pointer != nil {
		return TypePointer{
			Type: P_convert_type(decl, t.Pointer.Type),
			Mut:  t.Pointer.Mut,
		}
	} else if t.Func != nil {
		return TypeFunc{
			Input:  P_convert_param(decl, t.Func.Input),
			Output: P_convert_param(decl, t.Func.Output),
		}
	}
	panic("failed to parse type, everything was nil")
}

func P_convert_field(decl func(name string) IType, t []*P_Field) []TypeField {
	res := make([]TypeField, len(t))
	for i, p := range t {
		res[i] = TypeField{
			Name: p.Name,
			Type: P_convert_type(decl, p.Type),
		}
	}
	return res
}

func P_convert_param(decl func(name string) IType, t []*P_Param) []TypeParam {
	res := make([]TypeParam, len(t))
	for i, p := range t {
		res[i] = TypeParam{
			TypeField: TypeField{
				Name: p.Name,
				Type: P_convert_type(decl, p.Type),
			},
			Give: p.Give,
		}
	}
	return res
}

func P_convert_ident(decl func(name string) IType, t *P_Ident) IType {
	switch t.Name {
	case "void":
		return TypeP_void

	case "u8":
		return TypeP_u8
	case "u16":
		return TypeP_u16
	case "u32":
		return TypeP_u32
	case "u64":
		return TypeP_u64

	case "i8":
		return TypeP_i8
	case "i16":
		return TypeP_i16
	case "i32":
		return TypeP_i32
	case "i64":
		return TypeP_i64

	case "f32":
		return TypeP_f32
	case "f64":
		return TypeP_f64

	case "bool":
		return TypeP_bool
	case "string":
		return TypeP_string
	case "error":
		return TypeP_error
	case "size":
		return TypeP_usize
	}

	if d := decl(t.Name); d != nil {
		return TypeIdent{
			Type: d,
			Name: t.Name,
		}
	}

	panic("unknown identifier: " + t.Name)
}
