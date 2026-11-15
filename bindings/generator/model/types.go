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
type TypeAnonymous struct {
	Name string
}

// These functions are for the interface.

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
	case TypeP_buffer:
		return "buffer"
	}
	return "unknown"
}
func (TypeStruct) Kind() string    { return "struct" }
func (TypeEnum) Kind() string      { return "enum" }
func (TypeArray) Kind() string     { return "array" }
func (TypeSlice) Kind() string     { return "slice" }
func (TypePointer) Kind() string   { return "pointer" }
func (TypeFunc) Kind() string      { return "func" }
func (TypeIdent) Kind() string     { return "ident" }
func (TypeAnonymous) Kind() string { return "ident_unknown" }

const (
	TypeP_Invalid TypePrimitive = iota
	TypeP_buffer

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

func Keyof(t TypeRef) string {
	return fmt.Sprintf("%%%d", t)
}

// String converts a type into a string without recursion.
func String(p IType) string {
	switch e := p.(type) {
	case TypePrimitive:
		return e.Kind()
	case TypeStruct:
		s := strings.Builder{}
		s.WriteString("struct{")
		for _, f := range e.Fields {
			s.WriteString(f.Name)
			s.WriteString(" ")
			s.WriteString(Keyof(f.Type))
			s.WriteString(",")
		}
		s.WriteString("}")
		return s.String()
	case TypeEnum:
		s := strings.Builder{}
		s.WriteString("enum[")
		s.WriteString(Keyof(e.Type))
		s.WriteString("]{")
		for _, f := range e.Fields {
			s.WriteString(f)
			s.WriteString(",")
		}
		s.WriteString("}")
		return s.String()
	case TypeArray:
		return fmt.Sprintf("[%d]%s", e.Count, Keyof(e.Type))
	case TypeSlice:
		return fmt.Sprintf("[]%s", Keyof(e.Type))
	case TypePointer:
		is_mut := ""
		if e.Mut {
			is_mut = "mut "
		}
		return fmt.Sprintf("*%s%s", is_mut, Keyof(e.Type))
	case TypeFunc:
		s := strings.Builder{}
		s.WriteString("fn(")
		for _, f := range e.Input {
			s.WriteString(f.Name)
			s.WriteString(" ")
			s.WriteString(Keyof(f.Type))
			s.WriteString(",")
		}
		s.WriteString(")(")
		for _, f := range e.Output {
			s.WriteString(f.Name)
			s.WriteString(" ")
			s.WriteString(Keyof(f.Type))
			s.WriteString(",")
		}
		s.WriteString(")")
		return s.String()
	case TypeIdent:
		return fmt.Sprintf("%s:%s", e.Name, Keyof(e.Type))
	case TypeAnonymous:
		return e.Name
	default:
		return fmt.Sprintf("invalid type kind (%s): %v", p.Kind(), p)
	}
}
