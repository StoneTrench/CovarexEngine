package types

import (
	"fmt"
	"strings"
)

type TypeField struct {
	Name string
	Type TypeRef
}
type TypeParam struct {
	TypeField
	Give bool // Given also known as Ownership change or the opposite of Borrow.
}

type IType interface {
	itype()
	Kind() string
}
type TypePrimitive uint32
type TypeStruct struct {
	Fields []TypeField
}
type TypeEnum struct {
	Type   TypeRef
	Fields []string
}
type TypeArray struct {
	Count int
	Type  TypeRef
}
type TypeSlice struct {
	Type TypeRef
}
type TypePointer struct {
	Type TypeRef
	Mut  bool // Mutable also known as the opposite of Constant (Constant gets propagated)
}
type TypeFunc struct {
	Input  []TypeParam
	Output []TypeParam
}
type TypeIdent struct {
	Type TypeRef
	Name string
}

// These functions are for the interface.

func (TypePrimitive) itype() {}
func (TypeStruct) itype()    {}
func (TypeEnum) itype()      {}
func (TypeArray) itype()     {}
func (TypeSlice) itype()     {}
func (TypePointer) itype()   {}
func (TypeFunc) itype()      {}
func (TypeIdent) itype()     {}

type TypeKind string

// These Kind functions are here for the go templating to be able to figure out what kind to emit.

func (t TypePrimitive) Kind() string {
	switch t {
	case TypeP_Invalid:
		return "invalid"
	case TypeP_u8:
		return "u8"
	case TypeP_u16:
		return "u16"
	case TypeP_u32:
		return "u32"
	case TypeP_u64:
		return "u64"
	case TypeP_i8:
		return "i8"
	case TypeP_i16:
		return "i16"
	case TypeP_i32:
		return "i32"
	case TypeP_i64:
		return "i64"
	case TypeP_f32:
		return "f32"
	case TypeP_f64:
		return "f64"
	case TypeP_bool:
		return "bool"
	case TypeP_string:
		return "string"
	case TypeP_usize:
		return "usize"
	}
	return "unknown"
}
func (TypeStruct) Kind() string  { return "struct" }
func (TypeEnum) Kind() string    { return "enum" }
func (TypeArray) Kind() string   { return "array" }
func (TypeSlice) Kind() string   { return "slice" }
func (TypePointer) Kind() string { return "pointer" }
func (TypeFunc) Kind() string    { return "function" }
func (TypeIdent) Kind() string   { return "identifier" }

const (
	TypeP_Invalid TypePrimitive = iota
	TypeP_u8
	TypeP_u16
	TypeP_u32
	TypeP_u64

	TypeP_i8
	TypeP_i16
	TypeP_i32
	TypeP_i64

	TypeP_f32
	TypeP_f64

	TypeP_bool

	TypeP_string
	TypeP_usize // pointer
)

func Type_bodyof_string(c *TypeCollection) TypeRef {
	return c.Ref(TypeIdent{Name: "wasm_string", Type: Type_bodyof_slice(c, c.Ref(TypeP_u8), false)})
}
func Type_bodyof_slice(c *TypeCollection, t TypeRef, is_mut bool) TypeRef {
	return c.Ref(TypeStruct{
		Fields: []TypeField{
			{Name: "ptr", Type: c.Ref(TypePointer{Type: t, Mut: is_mut})},
			{Name: "len", Type: c.Ref(TypeP_usize)},
		},
	})
}

func keyof(t TypeRef) string {
	return fmt.Sprintf("%%%d", t)
}

// String converts a type into a string without recursion.
func String(p IType) string {
	switch e := p.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_u8:
			return "u8"
		case TypeP_u16:
			return "u16"
		case TypeP_u32:
			return "u32"
		case TypeP_u64:
			return "u64"

		case TypeP_i8:
			return "i8"
		case TypeP_i16:
			return "i16"
		case TypeP_i32:
			return "i32"
		case TypeP_i64:
			return "i64"

		case TypeP_f32:
			return "f32"
		case TypeP_f64:
			return "f64"

		case TypeP_bool:
			return "bool"
		case TypeP_string:
			return "string"
		case TypeP_usize:
			return "usize"
		}

		return fmt.Sprintf("invalid primitive: %d", e)
	case TypeStruct:
		s := strings.Builder{}
		s.WriteString("struct{")
		for _, f := range e.Fields {
			s.WriteString(f.Name)
			s.WriteString(" ")
			s.WriteString(keyof(f.Type))
			s.WriteString(",")
		}
		s.WriteString("}")
		return s.String()
	case TypeEnum:
		s := strings.Builder{}
		s.WriteString("enum[")
		s.WriteString(keyof(e.Type))
		s.WriteString("]{")
		for _, f := range e.Fields {
			s.WriteString(f)
			s.WriteString(",")
		}
		s.WriteString("}")
		return s.String()
	case TypeArray:
		return fmt.Sprintf("[%d]%s", e.Count, keyof(e.Type))
	case TypeSlice:
		return fmt.Sprintf("[]%s", keyof(e.Type))
	case TypePointer:
		is_mut := ""
		if e.Mut {
			is_mut = "mut "
		}
		return fmt.Sprintf("*%s%s", is_mut, keyof(e.Type))
	case TypeFunc:
		s := strings.Builder{}
		s.WriteString("fn(")
		for _, f := range e.Input {
			s.WriteString(f.Name)
			s.WriteString(" ")
			s.WriteString(keyof(f.Type))
			s.WriteString(",")
		}
		s.WriteString(")(")
		for _, f := range e.Output {
			s.WriteString(f.Name)
			s.WriteString(" ")
			s.WriteString(keyof(f.Type))
			s.WriteString(",")
		}
		s.WriteString(")")
		return s.String()
	case TypeIdent:
		return fmt.Sprintf("%s:%s", e.Name, keyof(e.Type))
	default:
		return fmt.Sprintf("invalid type kind: %v", p)
	}
}

func pad_struct(offset, alignment int) int {
	remainder := offset % alignment
	if remainder > 0 {
		return alignment - remainder
	} else {
		return 0
	}
}

// Sizeof returns the size of the type.
func SizeofAlignof(c *TypeCollection, r TypeRef, sizeof_ptr int) (int, int, error) {
	body, err := c.Body(r)
	if err != nil {
		return 0, 1, err
	}

	switch e := body.(type) {
	case TypePrimitive:
		{
			switch e {
			case TypeP_Invalid:
				return 0, 1, nil
			case TypeP_u8:
				return 1, 1, nil
			case TypeP_u16:
				return 2, 2, nil
			case TypeP_u32:
				return 4, 4, nil
			case TypeP_u64:
				return 8, sizeof_ptr, nil

			case TypeP_i8:
				return 1, 1, nil
			case TypeP_i16:
				return 2, 2, nil
			case TypeP_i32:
				return 4, 4, nil
			case TypeP_i64:
				return 8, sizeof_ptr, nil

			case TypeP_f32:
				return 4, 4, nil
			case TypeP_f64:
				return 8, sizeof_ptr, nil

			case TypeP_bool:
				return 1, 1, nil
			case TypeP_string:
				return SizeofAlignof(c, Type_bodyof_string(c), sizeof_ptr)
			case TypeP_usize:
				return sizeof_ptr, sizeof_ptr, nil
			}

			return 0, 1, nil
		}
	case TypeStruct:
		s := 0
		a := 1
		for _, f := range e.Fields {
			sf, af, err := SizeofAlignof(c, f.Type, sizeof_ptr)
			if err != nil {
				return 0, 1, err
			}
			if a < af {
				a = af
			}
			s += pad_struct(s, af)
			s += sf
		}
		s += pad_struct(s, a)
		return s, a, nil
	case TypeEnum:
		return SizeofAlignof(c, e.Type, sizeof_ptr)
	case TypeArray:
		sf, af, err := SizeofAlignof(c, e.Type, sizeof_ptr)
		return sf * e.Count, af, err
	case TypeSlice:
		return SizeofAlignof(c, Type_bodyof_slice(c, e.Type, true), sizeof_ptr)
	case TypePointer:
		return sizeof_ptr, sizeof_ptr, nil
	case TypeFunc:
		return sizeof_ptr, sizeof_ptr, nil
	case TypeIdent:
		return SizeofAlignof(c, e.Type, sizeof_ptr)
	}
	return 0, 1, fmt.Errorf("unknown kind: %v", body.Kind())
}

func _convert_fully_TypeToExplicitType(c *TypeCollection, r TypeRef, convert_fully bool) (TypeRef, error) {
	if convert_fully {
		return TypeToExplicitType(c, r, convert_fully)
	} else {
		return r, nil
	}
}

// TypeToExplicitType converts the type into a universal format, where everything is explicit and mostly ABI compatible.
// These however don't get directly sent across the abi using a function call.
func TypeToExplicitType(c *TypeCollection, r TypeRef, convert_fully bool) (TypeRef, error) {
	body, err := c.Body(r)
	if err != nil {
		return 0, err
	}

	switch e := body.(type) {
	case TypePrimitive:
		switch e {
		// case TypeP_u8, TypeP_u16, TypeP_u32, TypeP_i8, TypeP_i16, TypeP_i32, TypeP_bool, TypeP_usize:
		// 	return c.AddGetRef(TypeP_i32), nil
		// case TypeP_u64, TypeP_i64:
		// 	return c.AddGetRef(TypeP_i64), nil

		// case TypeP_f32:
		// 	return c.AddGetRef(TypeP_f32), nil
		// case TypeP_f64:
		// 	return c.AddGetRef(TypeP_f64), nil

		case TypeP_bool:
			return c.Ref(TypeP_u8), nil
		case TypeP_usize:
			return c.Ref(TypeP_i32), nil

		case TypeP_string:
			return _convert_fully_TypeToExplicitType(c, Type_bodyof_string(c), convert_fully)

		default:
			return r, nil
		}
	case TypeStruct:
		s := TypeStruct{}
		for _, f := range e.Fields {
			sf, err := TypeToExplicitType(c, f.Type, convert_fully)
			if err != nil {
				return 0, err
			}

			s.Fields = append(s.Fields, TypeField{
				Name: f.Name,
				Type: sf,
			})
		}
		return c.Ref(s), nil
	case TypeEnum:
		return _convert_fully_TypeToExplicitType(c, e.Type, convert_fully)
	case TypeArray:
		sf, err := TypeToExplicitType(c, e.Type, convert_fully)
		if err != nil {
			return 0, err
		}
		s := TypeArray{Count: e.Count, Type: sf}
		return c.Ref(s), nil
	case TypeSlice:
		return _convert_fully_TypeToExplicitType(c, Type_bodyof_slice(c, e.Type, true), convert_fully)
	case TypePointer:
		return _convert_fully_TypeToExplicitType(c, c.Ref(TypeP_usize), convert_fully)
	case TypeFunc:
		return _convert_fully_TypeToExplicitType(c, c.Ref(TypeP_usize), convert_fully)
	case TypeIdent:
		sf, err := TypeToExplicitType(c, e.Type, convert_fully)
		if err != nil {
			return 0, err
		}
		s := TypeIdent{Name: e.Name, Type: sf}
		return c.Ref(s), nil
	default:
		return 0, fmt.Errorf("invalid type kind: %s", body.Kind())
	}
}
