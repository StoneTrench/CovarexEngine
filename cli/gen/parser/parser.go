package parser

import (
	"errors"
	"fmt"
	"os"

	. "github.com/StoneTrench/CovarexEngine/cli/gen/types"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type t_SourceFile struct {
	Modules []*t_Module `parser:"@@*"`
}
type t_Module struct {
	Name  string    `parser:"'module' @Ident"`
	Types []*t_Decl `parser:"'{' @@* '}'"`
}
type t_Decl struct {
	Name     string   `parser:"@Ident"`
	Generics []string `parser:"('<' @Ident (',' @Ident)* '>')?"`
	Type     *t_Type  `parser:"@@"`
}
type t_Param struct {
	Name string  `parser:"@Ident"`
	Type *t_Type `parser:"@@"`
}
type t_Field struct {
	Name string  `parser:"@Ident"`
	Type *t_Type `parser:"@@"`
}

type t_Type struct {
	Paren   *t_Type    `parser:"'(' @@ ')'"`
	Struct  *t_Struct  `parser:"| @@"`
	Union   *t_Union   `parser:"| @@"`
	Enum    *t_Enum    `parser:"| @@"`
	Slice   *t_Slice   `parser:"| @@"`
	Func    *t_Func    `parser:"| @@"`
	Pointer *t_Pointer `parser:"| @@"`
	Array   *t_Array   `parser:"| @@"`
	Ident   *t_Ident   `parser:"| @@"`
}

type t_Ident struct {
	Name     string    `parser:"@Ident"`
	Generics []*t_Type `parser:"('<' @@ (',' @@)* '>')?"`
}
type t_Struct struct {
	Fields []*t_Field `parser:"'struct' '{' @@* '}'"`
}
type t_Union struct {
	Fields []*t_Field `parser:"'union' '{' @@* '}'"`
}
type t_Enum struct {
	Tag    string   `parser:"'enum' @('u8'|'u16'|'u32'|'u64'|'i8'|'i16'|'i32'|'i64'|'usize')"`
	Fields []string `parser:"'{' @Ident* '}'"`
}
type t_Array struct {
	Count int64   `parser:"'[' @Int ']'"`
	Type  *t_Type `parser:"@@"`
}
type t_Slice struct {
	Type *t_Type `parser:"'[' ']' @@"`
}
type t_Pointer struct {
	Mut  bool    `parser:"'*' @'mut'?"`
	Type *t_Type `parser:"@@"`
}
type t_Func struct {
	Input []*t_Param `parser:"'fn' '(' (@@ (',' @@)*)? ')'"`
	// Output []*P_Param `parser:"'(' (@@ (',' @@)*)? ')'"`
	Output []*t_Param `parser:"('(' (@@ (',' @@)*)? ')')?"`
}

type Module struct {
	Name  string
	Types []TypeField
}
type module_decl struct {
	Raw  *t_Decl
	Type TypeRef
}
type fnLookup = func(name string) (module_decl, error)

func convert_ident(c *TypeCollection, LUF fnLookup, node *t_Ident) (TypeRef, error) {
	if len(node.Generics) == 0 {
		if p, err := convert_primitive(node.Name); err == nil {
			return c.Ref(p), nil
		}
	}

	decl, err := LUF(node.Name)
	if err != nil {
		return TypeRef_Invalid, fmt.Errorf("error occured while looking for declaration: %w", err)
	}
	if len(node.Generics) != len(decl.Raw.Generics) {
		return TypeRef_Invalid, fmt.Errorf("incorrect number of generics: expected %d, got %d", len(decl.Raw.Generics), len(node.Generics))
	}

	if decl.Type != TypeRef_Invalid {
		return c.Ref(TypeIdent{
			Type: decl.Type,
			Name: decl.Raw.Name,
		}), nil
	}

	new_LUF := LUF
	if len(node.Generics) > 0 {
		local_generics_map := make(map[string]*t_Type)

		for i := range node.Generics {
			local_generics_map[decl.Raw.Generics[i]] = node.Generics[i]
		}

		new_LUF = func(name string) (module_decl, error) {
			if t, ok := local_generics_map[name]; ok {
				return module_decl{
					Raw: &t_Decl{
						Name:     "",
						Generics: nil,
						Type:     t,
					},
					Type: TypeRef_Invalid,
				}, nil
			}
			// if name is not inside generics, then go up one level
			return LUF(name)
		}
	}

	ref, err := convert_type(c, new_LUF, decl.Raw.Type)
	if err != nil {
		return TypeRef_Invalid, err
	}
	if len(decl.Raw.Generics) == 0 {
		decl.Type = ref
	}

	if len(decl.Raw.Name) == 0 {
		return ref, nil
	} else if len(decl.Raw.Generics) == 0 {
		return c.Ref(TypeIdent{
			Type: ref,
			Name: decl.Raw.Name,
		}), nil
	} else {
		return ref, nil
		// return c.Ref(TypeIdent{
		// 	Type: ref,
		// 	Name: fmt.Sprintf("%s-%d", decl.Raw.Name, ref),
		// }), nil
	}
}
func convert_primitive(name string) (TypePrimitive, error) {
	switch name {
	case "u8":
		return TypeP_u8, nil
	case "u16":
		return TypeP_u16, nil
	case "u32":
		return TypeP_u32, nil
	case "u64":
		return TypeP_u64, nil

	case "i8":
		return TypeP_i8, nil
	case "i16":
		return TypeP_i16, nil
	case "i32":
		return TypeP_i32, nil
	case "i64":
		return TypeP_i64, nil

	case "f32":
		return TypeP_f32, nil
	case "f64":
		return TypeP_f64, nil

	case "bool":
		return TypeP_bool, nil
	case "string":
		return TypeP_string, nil
	case "usize":
		return TypeP_usize, nil
	}

	return TypeP_Invalid, fmt.Errorf("unknown primitive: %s", name)
}
func convert_struct(c *TypeCollection, LUF fnLookup, t *t_Struct) (TypeRef, error) {
	fields, err := convert_fields(c, LUF, t.Fields)
	if err != nil {
		return TypeRef_Invalid, err
	}
	return c.Ref(TypeStruct{Fields: fields}), nil
}
func convert_enum(c *TypeCollection, _ fnLookup, t *t_Enum) (TypeRef, error) {
	tag, err := convert_primitive(t.Tag)
	if err != nil {
		return TypeRef_Invalid, err
	}
	body := TypeEnum{}
	body.Type = c.Ref(tag)
	body.Fields = t.Fields
	return c.Ref(body), nil
}
func convert_array(c *TypeCollection, LUF fnLookup, t *t_Array) (TypeRef, error) {
	e, err := convert_type(c, LUF, t.Type)
	if err != nil {
		return TypeRef_Invalid, err
	}
	return c.Ref(TypeArray{
		Count: int(t.Count),
		Type:  e,
	}), nil
}
func convert_slice(c *TypeCollection, LUF fnLookup, t *t_Slice) (TypeRef, error) {
	e, err := convert_type(c, LUF, t.Type)
	if err != nil {
		return TypeRef_Invalid, err
	}
	return c.Ref(TypeSlice{
		Type: e,
	}), nil
}
func convert_pointer(c *TypeCollection, LUF fnLookup, t *t_Pointer) (TypeRef, error) {
	e, err := convert_type(c, LUF, t.Type)
	if err != nil {
		return TypeRef_Invalid, err
	}
	return c.Ref(TypePointer{
		Type: e,
		Mut:  t.Mut,
	}), nil
}
func convert_func(c *TypeCollection, LUF fnLookup, t *t_Func) (TypeRef, error) {
	input, err := convert_params(c, LUF, t.Input)
	if err != nil {
		return TypeRef_Invalid, err
	}
	output, err := convert_params(c, LUF, t.Output)
	if err != nil {
		return TypeRef_Invalid, err
	}
	return c.Ref(TypeFunc{
		Input:  input,
		Output: output,
	}), nil
}
func convert_type(c *TypeCollection, LUF fnLookup, t *t_Type) (TypeRef, error) {
	if t == nil {
		return TypeRef_Invalid, nil
	} else if t.Paren != nil {
		return convert_type(c, LUF, t.Paren)
	} else if t.Ident != nil {
		return convert_ident(c, LUF, t.Ident)
	} else if t.Struct != nil {
		return convert_struct(c, LUF, t.Struct)
	} else if t.Enum != nil {
		return convert_enum(c, LUF, t.Enum)
	} else if t.Array != nil {
		return convert_array(c, LUF, t.Array)
	} else if t.Slice != nil {
		return convert_slice(c, LUF, t.Slice)
	} else if t.Pointer != nil {
		return convert_pointer(c, LUF, t.Pointer)
	} else if t.Func != nil {
		return convert_func(c, LUF, t.Func)
	}
	return TypeRef_Invalid, fmt.Errorf("failed to parse type, everything was nil")
}

func convert_module(c *TypeCollection, LUF fnLookup, t *t_Module) (*Module, error) {
	decl := make(map[string]int)
	res_decl := []module_decl{}

	new_LUF := func(name string) (module_decl, error) {
		if i, ok := decl[name]; ok {
			return res_decl[i], nil
		}
		return LUF(name)
	}

	for _, t := range t.Types {
		decl[t.Name] = len(res_decl)
		res_decl = append(res_decl, module_decl{
			Raw:  t,
			Type: TypeRef_Invalid,
		})
	}

	res_types := []TypeField{}
	// res_funcs := []TypeField{}
	err_list := []error{}
	for _, d := range res_decl {
		if len(d.Raw.Generics) == 0 {
			t := d.Type
			if t == TypeRef_Invalid {
				var err error
				t, err = convert_type(c, new_LUF, d.Raw.Type)
				if err != nil {
					err_list = append(err_list, err)
					continue
				}
			}
			r := TypeField{
				Name: d.Raw.Name,
				Type: t,
			}
			res_types = append(res_types, r)
			// switch d.Raw.Kind {
			// case DECL_FN:
			// 	res_funcs = append(res_funcs, r)
			// case DECL_TYPE:
			// 	res_types = append(res_types, r)
			// default:
			// 	err_list = append(err_list, fmt.Errorf("invalid kind of declaration: %s", d.Raw.Kind))
			// 	continue
			// }
		}
	}

	if len(err_list) > 0 {
		return nil, errors.Join(err_list...)
	}

	return &Module{
		Name:  t.Name,
		Types: res_types,
		// Funcs: res_funcs,
	}, nil
}

func convert_param(c *TypeCollection, LUF fnLookup, t *t_Param) (*TypeParam, error) {
	if t == nil {
		return nil, nil
	}
	e, err := convert_type(c, LUF, t.Type)
	if err != nil {
		return nil, err
	}
	return &TypeParam{
		TypeField: TypeField{Name: t.Name, Type: e},
		Give:      false,
	}, nil
}
func convert_params(c *TypeCollection, LUF fnLookup, t []*t_Param) ([]TypeParam, error) {
	res := make([]TypeParam, len(t))
	for i, p := range t {
		r, err := convert_param(c, LUF, p)
		if err != nil {
			return nil, err
		}
		res[i] = *r
	}
	return res, nil
}
func convert_fields(c *TypeCollection, LUF fnLookup, t []*t_Field) ([]TypeField, error) {
	res := make([]TypeField, len(t))
	for i, p := range t {
		e, err := convert_type(c, LUF, p.Type)
		if err != nil {
			return nil, err
		}
		res[i] = TypeField{
			Name: p.Name,
			Type: e,
		}
	}
	return res, nil
}

var LEXER = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comment", Pattern: `//[^\n]*`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_\-]*`},
	{Name: "Punct", Pattern: `[{}()\[\],*<>]`}, // added '<' '>'
	{Name: "Whitespace", Pattern: `[ \t\r\n;]+`},
	{Name: "Int", Pattern: `[0-9][0-9]*`},
})
var PARSER = participle.MustBuild[t_SourceFile](
	participle.Lexer(LEXER),
	participle.Elide("Comment", "Whitespace"),
)

func ParseSource(c *TypeCollection, name string) ([]Module, error) {
	bits, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	return parse_bytes(c, name, bits)
}
func parse_bytes(c *TypeCollection, name string, bits []byte) ([]Module, error) {
	src, err := PARSER.ParseBytes(name, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	LUF := func(name string) (module_decl, error) {
		var n module_decl
		return n, fmt.Errorf("declaration %s doesn't exist", name)
	}
	result := make([]Module, len(src.Modules))
	for i, mod := range src.Modules {
		m, err := convert_module(c, LUF, mod)
		if err != nil {
			return nil, err
		}
		result[i] = *m
	}
	return result, nil
}
