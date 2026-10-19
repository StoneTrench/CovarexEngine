package types

import "fmt"

type TypeField struct {
	Name string
	Type IType
}

type TypeParam struct {
	TypeField
	Give bool // also known as Own or the opposite of Borrow.
}

type IType interface{ itype() }
type TypePrimitive uint32
type TypeStruct struct {
	Fields []TypeField
}
type TypeArray struct {
	Count int
	Type  IType
}
type TypeSlice struct {
	Type IType
}
type TypePointer struct {
	Type IType
	Mut  bool // Mutable also known as the opposite of Constant (Constant gets propagated)
}
type TypeFunc struct {
	Input  []TypeParam
	Output []TypeParam
}
type TypeIdent struct {
	Type IType
	Name string
}

func (TypePrimitive) itype() {}
func (TypeStruct) itype()    {}
func (TypeArray) itype()     {}
func (TypeSlice) itype()     {}
func (TypePointer) itype()   {}
func (TypeFunc) itype()      {}
func (TypeIdent) itype()     {}

const (
	TypeP_void TypePrimitive = iota

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
	TypeP_error
	TypeP_usize // pointer
)

var Type_bodyof_string = Type_bodyof_slice(TypeP_u8, false)
var Type_bodyof_error = TypePointer{Type: Type_bodyof_string, Mut: false}

func Type_bodyof_slice(t IType, is_mut bool) TypeStruct {
	return TypeStruct{
		Fields: []TypeField{
			{Name: "Ptr", Type: TypePointer{Type: t, Mut: is_mut}},
			{Name: "Len", Type: TypeP_usize},
		},
	}
}

// GetABIAllocation whether we should alloc/free or do nothing (aka it just copies/writes without alloc/free).
func GetABIAllocation(p IType) bool {
	switch e := p.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			return GetABIAllocation(Type_bodyof_string)
		case TypeP_error:
			return GetABIAllocation(Type_bodyof_error)
		case TypeP_usize:
			return false // this is correct, because here, the other side just copies.
		}

		return false
	case TypeArray:
		return true
	case TypeSlice:
		return true
	case TypeFunc:
		panic("can't pass function through wasm abi")
	case TypePointer:
		return true
	case TypeStruct:
		return true
	case TypeIdent:
		return GetABIAllocation(e.Type)
	default:
		panic(fmt.Sprintf("invalid type type: %v", p))
	}
}

// GetCountFlat gets the number of primitives this type will be flattened into.
func GetCountFlat(p IType) int {
	switch e := p.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			return 2
		case TypeP_error:
			return 1
		}

		return 1
	case TypeArray:
		return e.Count * GetCountFlat(e.Type)
	case TypeSlice:
		return 2
	case TypeFunc:
		return 1
	case TypePointer:
		return 1
	case TypeStruct:
		count := 0
		for _, f := range e.Fields {
			count += GetCountFlat(f.Type)
		}
		return count
	case TypeIdent:
		return GetCountFlat(e.Type)
	default:
		panic(fmt.Sprintf("invalid type type: %v", p))
	}
}

// // ParamToABI expands the parameter given, to the fundemental types it is made up of.
// // Everything becomes a number after enough iterations.
// // Only expands one step.
// func ParamToABI(p TypeParam) ([]TypeParam, bool) {
// 	name := p.Name
// 	switch e := p.Type.(type) {
// 	case TypePrimitive:
// 		p.Name = fmt.Sprintf("%s_w", name)
// 		switch e {
// 		case TypeP_string:
// 			return []TypeParam{
// 				{fmt.Sprintf("%s_str", name), TypeP_usize},
// 				{fmt.Sprintf("%s_len", name), TypeP_u32},
// 			}, true
// 		case TypeP_bool:
// 			return []TypeParam{{fmt.Sprintf("%s_b", name), TypeP_u32}}, true
// 		case TypeP_u8, TypeP_u16, TypeP_u32:
// 			p.Type = TypeP_u32
// 			return []TypeParam{p}, true
// 		case TypeP_i8, TypeP_i16, TypeP_i32:
// 			p.Type = TypeP_i32
// 			return []TypeParam{p}, true
// 		case TypeP_u64:
// 			p.Type = TypeP_u64
// 			return []TypeParam{p}, true
// 		case TypeP_i64:
// 			p.Type = TypeP_i64
// 			return []TypeParam{p}, true
// 		case TypeP_error:
// 			return []TypeParam{}, true
// 		}

// 		return []TypeParam{p}, true
// 	case TypeArray:
// 		return []TypeParam{
// 			{fmt.Sprintf("%s_ptr", name), TypeP_usize},
// 		}, true
// 	case TypeSlice:
// 		return []TypeParam{
// 			{fmt.Sprintf("%s_ptr", name), TypeP_usize},
// 			{fmt.Sprintf("%s_len", name), TypeP_u32},
// 		}, true
// 	case TypeFunc:
// 		return []TypeParam{
// 			{fmt.Sprintf("%s_fn", name), TypeP_usize},
// 		}, true
// 	case TypePointer:
// 		return []TypeParam{
// 			{fmt.Sprintf("%s_ptr", name), TypeP_usize},
// 		}, true
// 	case TypeStruct:
// 		res := []TypeParam{}
// 		for _, f := range e.Fields {
// 			f.Name = fmt.Sprintf("%s_%s_f", name, f.Name)
// 			res = append(res, f)
// 		}
// 		return res, false
// 	case TypeIdent:
// 		p.Type = e.Type
// 		return ParamToABI(p)
// 	default:
// 		panic(fmt.Sprintf("invalid type type: %v", p.Type))
// 	}
// }

// // RepParamToABI calls ParamToABI until all the types have been fully expanded to be ABI compatible.
// func RepParamToABI(p []TypeParam) []TypeParam {
// 	to_be_lowered := p
// 	result := []TypeParam{}

// 	for len(to_be_lowered) > 0 {
// 		next_to_be_lowered := []TypeParam{}
// 		for _, p := range to_be_lowered {
// 			r, is_low := ParamToABI(p)
// 			if is_low {
// 				result = append(result, r...)
// 			} else {
// 				next_to_be_lowered = append(next_to_be_lowered, r...)
// 			}
// 		}

// 		to_be_lowered = next_to_be_lowered
// 	}

// 	return result
// }
