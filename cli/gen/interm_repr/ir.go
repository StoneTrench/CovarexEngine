package interm_repr

import (
	"fmt"
	"strings"

	"github.com/StoneTrench/CovarexEngine/cli/gen/types"
)

type TColl = types.TypeCollection

/*
These are the types that the template will process to create the binding sourcecode in the target language.

Assume that types are all packed (no padding in struct).
*/

type Name struct {
	Name string
	Type types.TypeRef
}
type Module struct {
	WasmName string
	Name     string
	Guest    []FuncDecl
	Host     []FuncDecl
	Types    []Name
}
type FuncDecl struct {
	Module        *Module
	WasmName      string
	Name          string
	DeclaredNames []Name
	WasmParams    []Name
	OutputStruct  types.TypeRef
	Input         []Name
	Output        []Name
	Block         []Instr
}
type Instr struct {
	Func  *FuncDecl
	Index int
	Kind  string

	Integer int
	Text    string

	Input  []Name
	Output []Name
	Block  []Instr
}

func (t Module) StringPadding(depth int, pad string) string {
	s := strings.Builder{}

	has_padding := len(pad) > 0
	pad_string := strings.Repeat(pad, depth)

	if has_padding {
		s.WriteString(pad_string)
	}

	s.WriteString(t.Name)
	s.WriteString("{")
	if has_padding {
		s.WriteString("\n")
	}
	for _, e := range t.Guest {
		s.WriteString(e.StringPadding(depth+1, pad))
		if has_padding {
			s.WriteString("\n")
		}
	}
	s.WriteString(pad_string)
	s.WriteString("}")
	// s.WriteString("{")
	// if has_padding {
	// 	s.WriteString("\n")
	// }
	// for _, e := range t.Host {
	// 	s.WriteString(e.StringPadding(depth+1, pad))
	// 	if has_padding {
	// 		s.WriteString("\n")
	// 	}
	// }
	// s.WriteString(pad_string)
	// s.WriteString("}")

	return s.String()
}
func (t Module) String() string {
	return t.StringPadding(0, "")
}

func (t FuncDecl) StringPadding(depth int, pad string) string {
	s := strings.Builder{}

	has_padding := len(pad) > 0
	pad_string := strings.Repeat(pad, depth)

	if has_padding {
		s.WriteString(pad_string)
	}

	s.WriteString(t.Name)
	s.WriteString("(")
	arr := make([]string, len(t.Input))
	for i, e := range t.Input {
		arr[i] = fmt.Sprintf("%s:%s", e.Name, types.Keyof(e.Type))
	}
	s.WriteString(strings.Join(arr, ", "))
	s.WriteString(")")
	if t.Output != nil {
		s.WriteString("(")
		arr := make([]string, len(t.Output))
		for i, e := range t.Output {
			arr[i] = fmt.Sprintf("%s:%s", e.Name, types.Keyof(e.Type))
		}
		s.WriteString(strings.Join(arr, ", "))
		s.WriteString(")")
	}
	s.WriteString("{")
	if has_padding {
		s.WriteString("\n")
	}
	for _, e := range t.Block {
		s.WriteString(e.StringPadding(depth+1, pad))
		s.WriteString(";")
		if has_padding {
			s.WriteString("\n")
		}
	}
	s.WriteString(pad_string)
	s.WriteString("}")

	return s.String()
}
func (t FuncDecl) String() string {
	return t.StringPadding(0, "")
}

func (t Instr) StringPadding(depth int, pad string) string {
	s := new(strings.Builder)

	has_padding := len(pad) > 0
	pad_string := strings.Repeat(pad, depth)

	if has_padding {
		s.WriteString(pad_string)
	}

	if len(t.Output) > 0 {
		arr := make([]string, len(t.Output))
		for i, e := range t.Output {
			arr[i] = fmt.Sprintf("%s:%s", e.Name, types.Keyof(e.Type))
		}
		s.WriteString(strings.Join(arr, ", "))
		s.WriteString(" = ")
	}
	s.WriteString(t.Kind)
	fmt.Fprintf(s, "[%d,%q]", t.Integer, t.Text)
	s.WriteString("(")
	if len(t.Input) > 0 {
		arr := make([]string, len(t.Input))
		for i, e := range t.Input {
			arr[i] = fmt.Sprintf("%s:%s", e.Name, types.Keyof(e.Type))
		}
		s.WriteString(strings.Join(arr, ", "))
	}
	s.WriteString(")")

	if len(t.Block) > 0 {
		s.WriteString("{")
		if has_padding {
			s.WriteString("\n")
		}
		for _, e := range t.Block {
			s.WriteString(e.StringPadding(depth+1, pad))
			s.WriteString(";")
			if has_padding {
				s.WriteString("\n")
			}
		}

		s.WriteString(pad_string)
		s.WriteString("}")
	}

	return s.String()
}
func (t Instr) String() string {
	return t.StringPadding(0, "")
}

var op_host_call = "Call"                       // Call the host function (HOST ONLY)
var op_host_make_binding = "MakeHostBinding"    // Declare a function in a function (HOST ONLY)
var op_host_write_memory_u32 = "WriteMemoryU32" // Writes to the offset inside the wasm module (HOST ONLY)

var op_return = "Return" // Return statement

var op_call_IO = "CallWireIO"       // Call with an input and output buffer
var op_call_I = "CallWireI"         // Call with an input buffer only
var op_call_O = "CallWireO"         // Call with an output buffer only
var op_call_empty = "CallWireEmpty" // Call without any transmission

var op_declare_unknown = "DeclareUnknown" // Declare but do not initialize a variable with a type

var op_counter_init = "MakeCounter"           // New counter
var op_counter_increment = "IncrementCounter" // Add 1 to counter

var op_loop_fixed = "LoopFixed"       // Loop a fixed known n times
var op_loop_until = "LoopUntil"       // Loop based on a value
var op_if_not_zero = "IfNotZeroBlock" // Single if statement that executes if input is not 0

var op_buffer_make = "MakeBuffer"                 // Make buffer.
var op_buffer_make_from_raw = "MakeBufferFromRaw" // Make buffer from pointer and length.
var op_buffer_free = "FreeBuffer"                 // Destroy buffer.

var op_buffer_pointer = "BufferPointer" // Get 0/start pointer of the buffer. On the host this copies the buffer into the guest.
var op_buffer_length = "BufferLength"   // Get length of buffer.

var op_buffer_io = map[types.TypePrimitive][2]string{
	types.TypeP_u8:  {"WriteBufU8", "ReadBufU8"},
	types.TypeP_u16: {"WriteBufU16", "ReadBufU16"},
	types.TypeP_u32: {"WriteBufU32", "ReadBufU32"},
	types.TypeP_u64: {"WriteBufU64", "ReadBufU64"},
	types.TypeP_i8:  {"WriteBufI8", "ReadBufI8"},
	types.TypeP_i16: {"WriteBufI16", "ReadBufI16"},
	types.TypeP_i32: {"WriteBufI32", "ReadBufI32"},
	types.TypeP_i64: {"WriteBufI64", "ReadBufI64"},

	types.TypeP_f32: {"WriteBufF32", "ReadBufF32"},
	types.TypeP_f64: {"WriteBufF64", "ReadBufF64"},
}
var op_string_conv = [2]string{"StringToSlice", "SliceToString"}
var op_pointer_deref_conv = [2]string{"PointerToValue", "ValueToPointer"} // deref and ref
var op_pointer_usize_conv = [2]string{"PointerToUsize", "UsizeToPointer"}
var op_bool_conv = [2]string{"BoolToU8", "U8ToBool"}
var op_usize_u32_conv = [2]string{"UsizeToU32", "U32ToUsize"}

var op_enum_conv = map[types.TypePrimitive][2]string{
	types.TypeP_u8:  {"EnumToU8", "U8ToEnum"},
	types.TypeP_u16: {"EnumToU16", "U16ToEnum"},
	types.TypeP_u32: {"EnumToU32", "U32ToEnum"},
	types.TypeP_u64: {"EnumToU64", "U64ToEnum"},
	types.TypeP_i8:  {"EnumToI8", "I8ToEnum"},
	types.TypeP_i16: {"EnumToI16", "I16ToEnum"},
	types.TypeP_i32: {"EnumToI32", "I32ToEnum"},
	types.TypeP_i64: {"EnumToI64", "I64ToEnum"},
}

var op_plex_make = "MakeStruct" // New struct
var op_array_make = "MakeArray" // New array
var op_slice_make = "MakeSlice" // New slice
var op_slice_free = "FreeSlice" // Destroy slice.

var op_move_output = "MoveOutput" // Create an output variable from a block
var op_move = "MoveInput"         // Assign to the output variable

var op_plex_index = [2]string{"ReadStruct", "WriteStruct"} // read field from struct
var op_array_index = [2]string{"ReadArray", "WriteArray"}  // read element from fixed array
var op_slice_index = [2]string{"ReadSlice", "WriteSlice"}  // read element from slice
var op_slice_length = "SliceLength"                        // Get length of slice

var op_nop_comment = "Nop" // Its a comment
