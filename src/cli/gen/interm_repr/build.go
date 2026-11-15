package intermediate_representation

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/StoneTrench/go-mc-clone/bindings/generator/model"
)

// NormalizedIdentifier converts an .abi identifier to an NormalizedIdentifier identifier:
//   - '-' removed, next rune capitalized
//   - '_' preserved, next rune capitalized
//   - first rune forced uppercase
func NormalizedIdentifier(s string) string {
	var b strings.Builder
	capNext := true
	for _, r := range s {
		switch {
		case r == '-':
			capNext = true
		case r == '_':
			b.WriteRune(r)
			capNext = true
		case r >= 'a' && r <= 'z':
			if capNext {
				r -= 'a' - 'A'
			}
			b.WriteRune(r)
			capNext = false
		default:
			b.WriteRune(r)
			capNext = false
		}
	}
	return b.String()
}
func NameToTypeParam(n Name) TypeParam {
	return TypeParam{
		TypeField: TypeField{Name: n.Name, Type: n.Type},
	}
}

func BuildAgnosticModule(c *TColl, name string, fields []TypeField) (*Module, error) {
	mod := new(Module)
	mod.WasmName = name
	mod.Name = NormalizedIdentifier(name)

	err_list := []error{}

	for _, field := range fields {
		body, err := c.Body(field.Type)
		if err != nil {
			err_list = append(err_list, err)
			continue
		}
		if type_fn, ok := body.(TypeFunc); ok {
			fun, err := decl_func_guest(c, mod, field.Name, &type_fn)
			if err != nil {
				err_list = append(err_list, err)
				continue
			}
			mod.Guest = append(mod.Guest, *fun)

			fun, err = decl_func_host(c, mod, field.Name, &type_fn)
			if err != nil {
				err_list = append(err_list, err)
				continue
			}
			mod.Host = append(mod.Host, *fun)
		} else {
			mod.Types = append(mod.Types, Name{Name: NormalizedIdentifier(field.Name), Type: field.Type})
		}
	}

	for i, t := range c.Types {
		{
			ident, ok := t.(TypeIdent)
			if ok {
				ref := TypeRef(i + 1)
				ident.Name = NormalizedIdentifier(ident.Name)
				err := c.Update(ref, ident)
				if err != nil {
					err_list = append(err_list, err)
				}
			}
		}
		// This has to ignore TypeIdentUnknown, DO NOT NORMALIZE THAT
	}

	if len(err_list) > 0 {
		return nil, errors.Join(err_list...)
	}

	return mod, nil
}
func decl_func_guest(c *TColl, mod *Module, name string, fn *TypeFunc) (*FuncDecl, error) {
	TYPE_WASM_POINTER := c.Ref(TypeP_u32)
	TYPE_BUFFER := c.Ref(TypeP_buffer)

	result := new(FuncDecl)
	result.Module = mod
	result.WasmName = name
	result.Name = NormalizedIdentifier(name)

	b := new(fnBuilder)
	b.TypColl = c
	b.NameC = new(NameCount)
	b.Block = []Instr{}
	b.Func = result

	N_BUFFER_INPUT := b.Next(TYPE_BUFFER)
	N_BUFFER_OUTPUT := b.Next(TYPE_BUFFER)

	N_BUFFER_INPUT_PTR := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_INPUT_LEN := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_OUTPUT_PTR := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_OUTPUT_LEN := b.Next(TYPE_WASM_POINTER)

	HAS_INPUT := len(fn.Input) > 0
	HAS_OUTPUT := len(fn.Output) > 0

	N_FN_INPUT := []Name{}
	N_FN_OUTPUT := []Name{}

	err_list := []error{}

	for _, e := range fn.Input {
		N_FN_INPUT = append(N_FN_INPUT, Name{Name: e.Name, Type: e.Type})
	}
	TYPE_OUTPUT_STRUCT := TypeStruct{}
	for _, e := range fn.Output {
		N_FN_OUTPUT = append(N_FN_OUTPUT, Name{Name: e.Name, Type: e.Type})
		TYPE_OUTPUT_STRUCT.Fields = append(TYPE_OUTPUT_STRUCT.Fields, e.TypeField)
	}
	if len(N_FN_OUTPUT) > 1 {
		result.OutputStruct = b.Ref(TYPE_OUTPUT_STRUCT)
	}

	result.Input = N_FN_INPUT
	result.Output = N_FN_OUTPUT

	if HAS_INPUT {
		b.Func.WasmParams = append(b.Func.WasmParams, N_BUFFER_INPUT_PTR, N_BUFFER_INPUT_LEN)
		b.Buffer = N_BUFFER_INPUT
		instr_buf := b.MakeBuffer()
		b.NopComment("encoding input")
		err_list = append(err_list, b.Serialize(WRITE_TO_BUFFER, N_FN_INPUT)...)
		b.NopComment("send")
		b.BufferPointer(N_BUFFER_INPUT_PTR)
		b.BufferLength(N_BUFFER_INPUT_LEN)
		b.Block[instr_buf.Index].Integer = b.BufferMinSize
		b.BufferMinSize = 0
	}

	if HAS_OUTPUT {
		b.Write_0x1(op_declare_unknown, N_BUFFER_OUTPUT_PTR)
		b.Write_0x1(op_declare_unknown, N_BUFFER_OUTPUT_LEN)
	}

	kind, args := op_call_empty, []Name{}
	switch {
	case HAS_INPUT && HAS_OUTPUT:
		kind = op_call_IO
		args = []Name{N_BUFFER_INPUT_PTR, N_BUFFER_INPUT_LEN, N_BUFFER_OUTPUT_PTR, N_BUFFER_OUTPUT_LEN}
	case HAS_INPUT:
		kind = op_call_I
		args = []Name{N_BUFFER_INPUT_PTR, N_BUFFER_INPUT_LEN}
	case HAS_OUTPUT:
		kind = op_call_O
		args = []Name{N_BUFFER_OUTPUT_PTR, N_BUFFER_OUTPUT_LEN}
	}
	b.Write(Instr{Kind: kind, Text: name, Input: args})

	if HAS_INPUT {
		b.FreeBuffer()
	}

	if HAS_OUTPUT {
		b.Func.WasmParams = append(b.Func.WasmParams, N_BUFFER_OUTPUT_PTR, N_BUFFER_OUTPUT_LEN)
		b.Buffer = N_BUFFER_OUTPUT
		b.MakeBufferFromRaw(N_BUFFER_OUTPUT_PTR, N_BUFFER_OUTPUT_LEN)
		b.NopComment("decoding output")
		err_list = append(err_list, b.Serialize(READ_FROM_BUFFER, N_FN_OUTPUT)...)
		b.FreeBuffer()

		old_block := b.Block
		b.Block = []Instr{}

		if result.OutputStruct != TypeRef_Invalid {
			plex := b.Next(result.OutputStruct)
			b.Write_0x1(op_plex_make, plex)
			for _, f := range N_FN_OUTPUT {
				b.Write_2x0(op_plex_index[READ_FROM_BUFFER], plex, f).Text = f.Name
			}
			b.Write(Instr{Kind: op_return, Input: []Name{plex}})
		}

		new_block := b.Block
		b.Block = old_block

		b.Write(Instr{Kind: op_return, Input: N_FN_OUTPUT, Block: new_block})
	}
	result.Block = b.Block

	if len(err_list) > 0 {
		return nil, errors.Join(err_list...)
	}

	return result, nil
}
func decl_func_host(c *TColl, mod *Module, name string, fn *TypeFunc) (*FuncDecl, error) {
	TYPE_WASM_POINTER := c.Ref(TypeP_u32)
	TYPE_BUFFER := c.Ref(TypeP_buffer)

	result := new(FuncDecl)
	result.Module = mod
	result.WasmName = name
	result.Name = "HostAPI" + NormalizedIdentifier(name)

	b := new(fnBuilder)
	b.TypColl = c
	b.NameC = new(NameCount)
	b.Block = []Instr{}
	b.Func = result

	N_BUFFER_INPUT := b.Next(TYPE_BUFFER)
	N_BUFFER_OUTPUT := b.Next(TYPE_BUFFER)

	N_BUFFER_INPUT_PTR := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_INPUT_LEN := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_OUTPUT_PTR := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_OUTPUT_LEN := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_OUTPUT_PTR_PTR := b.Next(TYPE_WASM_POINTER)
	N_BUFFER_OUTPUT_LEN_PTR := b.Next(TYPE_WASM_POINTER)

	HAS_INPUT := len(fn.Input) > 0
	HAS_OUTPUT := len(fn.Output) > 0

	N_FN_INPUT := []Name{}
	N_FN_OUTPUT := []Name{}

	err_list := []error{}

	for _, e := range fn.Input {
		N_FN_INPUT = append(N_FN_INPUT, Name{Name: e.Name, Type: e.Type})
	}
	for _, e := range fn.Output {
		N_FN_OUTPUT = append(N_FN_OUTPUT, Name{Name: e.Name, Type: e.Type})
	}

	N_API_MODULE := Name{Name: "mod", Type: b.Ref(TypeAnonymous{Name: "api.Module"})}

	N_CALLBACK := b.NextT(TypeFunc{
		Input:  append([]TypeParam{NameToTypeParam(N_API_MODULE)}, fn.Input...),
		Output: fn.Output,
	})
	N_HOST_BINDING := b.NextT(TypeAnonymous{Name: "HostBinding"})
	N_BINDING_INPUT := []Name{
		{Name: "ctx", Type: b.Ref(TypeAnonymous{Name: "context.Context"})},
		N_API_MODULE,
	}

	result.Input = []Name{N_CALLBACK}
	result.Output = []Name{N_HOST_BINDING}

	if HAS_INPUT {
		N_BINDING_INPUT = append(N_BINDING_INPUT, N_BUFFER_INPUT_PTR, N_BUFFER_INPUT_LEN)
		b.Buffer = N_BUFFER_INPUT
		b.MakeBufferFromRaw(N_BUFFER_INPUT_PTR, N_BUFFER_INPUT_LEN)
		b.NopComment("decoding input")
		err_list = append(err_list, b.Serialize(READ_FROM_BUFFER, N_FN_INPUT)...)
	}

	b.NopComment("send")
	b.Write(Instr{
		Kind:   op_host_call,
		Input:  append([]Name{N_CALLBACK}, N_FN_INPUT...),
		Output: N_FN_OUTPUT,
	})

	if HAS_OUTPUT {
		N_BINDING_INPUT = append(N_BINDING_INPUT, N_BUFFER_OUTPUT_PTR_PTR, N_BUFFER_OUTPUT_LEN_PTR)
		b.Buffer = N_BUFFER_OUTPUT
		instr_buf := b.MakeBuffer() // makes it in the host
		b.NopComment("encoding output")
		err_list = append(err_list, b.Serialize(WRITE_TO_BUFFER, N_FN_OUTPUT)...)
		b.BufferPointer(N_BUFFER_OUTPUT_PTR) // copies to guest
		b.BufferLength(N_BUFFER_OUTPUT_LEN)
		b.WriteMemoryU32(N_BUFFER_OUTPUT_PTR_PTR, N_BUFFER_OUTPUT_PTR)
		b.WriteMemoryU32(N_BUFFER_OUTPUT_LEN_PTR, N_BUFFER_OUTPUT_LEN)
		b.FreeBuffer()
		b.Block[instr_buf.Index].Integer = b.BufferMinSize
		b.BufferMinSize = 0
	}

	new_block := b.Block
	b.Block = []Instr{}

	b.Write(Instr{Kind: op_host_make_binding, Input: N_BINDING_INPUT, Output: result.Output, Block: new_block})
	b.Write_1x0(op_return, N_HOST_BINDING)

	if len(err_list) > 0 {
		return nil, errors.Join(err_list...)
	}
	result.Block = b.Block

	return result, nil
}

type ir_writer_mode uint8

const (
	WRITE_TO_BUFFER ir_writer_mode = iota
	READ_FROM_BUFFER
)

func bWrite_primitive_write(b *fnBuilder, input Name, e TypePrimitive) error {
	m := WRITE_TO_BUFFER

	switch e {
	case TypeP_Invalid:
		return fmt.Errorf("primitive was invalid")
	case TypeP_buffer:
		return fmt.Errorf("cannot put a buffer primitive into a buffer primitive")

	case TypeP_bool:
		e = TypeP_u8
		output := b.NextT(e)
		b.Write_1x1(op_bool_conv[m], input, output)
		return bWrite_primitive_write(b, output, e)
	case TypeP_usize:
		e = TypeP_u32
		output := b.NextT(e)
		b.Write_1x1(op_usize_conv[m], input, output)
		return bWrite_primitive_write(b, output, e)

	case TypeP_string:
		output := b.NextT(TypeSlice{Type: b.Ref(TypeP_u8)})
		b.Write_1x1(op_string_conv[m], input, output)
		return bWrite_name(b, m, output)
	}

	op, ok := op_buffer_io[e]
	if !ok {
		return fmt.Errorf("missing buffer instruction for (%s)", e.Kind())
	}
	b.WriteBuf(m, op[m], input)
	return nil
}
func bWrite_primitive_read(b *fnBuilder, output Name, e TypePrimitive) error {
	m := READ_FROM_BUFFER

	switch e {
	case TypeP_Invalid:
		return fmt.Errorf("primitive was invalid")
	case TypeP_buffer:
		return fmt.Errorf("cannot put a buffer primitive into a buffer primitive")

	case TypeP_bool:
		e = TypeP_u8
		input := b.NextT(e)
		err := bWrite_primitive_read(b, input, e)
		if err != nil {
			return err
		}
		b.Write_1x1(op_bool_conv[m], input, output)
		return nil
	case TypeP_usize:
		e = TypeP_u32
		input := b.NextT(e)
		err := bWrite_primitive_read(b, input, e)
		if err != nil {
			return err
		}
		b.Write_1x1(op_usize_conv[m], input, output)
		return nil

	case TypeP_string:
		input := b.NextT(TypeSlice{Type: b.Ref(TypeP_u8)})
		err := bWrite_name(b, m, input)
		if err != nil {
			return err
		}
		b.Write_1x1(op_string_conv[m], input, output)
		return nil
	}

	op, ok := op_buffer_io[e]
	if !ok {
		return fmt.Errorf("missing buffer instruction for (%s)", e.Kind())
	}
	b.WriteBuf(m, op[m], output)
	return nil
}
func bWrite_primitive(b *fnBuilder, m ir_writer_mode, name Name, e TypePrimitive) error {
	switch m {
	case WRITE_TO_BUFFER:
		return bWrite_primitive_write(b, name, e)
	case READ_FROM_BUFFER:
		return bWrite_primitive_read(b, name, e)
	}
	return nil
}
func bWrite_struct(b *fnBuilder, m ir_writer_mode, name Name, e TypeStruct) error {
	if m == READ_FROM_BUFFER {
		b.Write_0x1(op_plex_make, name)
	}
	for _, f := range e.Fields {
		output := b.Next(f.Type)
		if m == WRITE_TO_BUFFER {
			b.Write_1x1(op_plex_index[m], name, output).Text = f.Name
		}
		err := bWrite_name(b, m, output)
		if err != nil {
			return fmt.Errorf("failure in struct field (%s:%s): %w", f.Name, Keyof(f.Type), err)
		}
		if m == READ_FROM_BUFFER {
			b.Write_2x0(op_plex_index[m], name, output).Text = f.Name
		}
	}
	return nil
}
func bWrite_enum(b *fnBuilder, m ir_writer_mode, name Name, e TypeEnum) error {
	enum_type, err := b.Body(e.Type)
	if err != nil {
		return fmt.Errorf("failed to find body of enum: %w", err)
	}
	prim, ok := enum_type.(TypePrimitive)
	if !ok {
		return fmt.Errorf("expected primitive in enum, got: %s", enum_type.Kind())
	}
	op, ok := op_enum_conv[prim]
	if !ok {
		return fmt.Errorf("expected integer in enum, got: %s", prim.Kind())
	}

	output := b.Next(e.Type)
	switch m {
	case WRITE_TO_BUFFER:
		b.Write_1x1(op[m], name, output)
		return bWrite_name(b, m, output)
	case READ_FROM_BUFFER:
		err := bWrite_name(b, m, output)
		if err != nil {
			return err
		}
		b.Write_1x1(op[m], output, name)
	}

	return nil
}
func bWrite_array(b *fnBuilder, m ir_writer_mode, name Name, e TypeArray) error {
	index := b.NextT(TypeP_usize)
	element := b.Next(e.Type)

	b.Write_0x1(op_counter_init, index)

	if m == READ_FROM_BUFFER {
		b.Write_0x1(op_array_make, name)
	}

	old_block := b.Block
	b.Block = []Instr{}

	switch m {
	case WRITE_TO_BUFFER:
		b.Write_2x1(op_array_index[m], name, index, element)
		err := bWrite_name(b, m, element)
		if err != nil {
			return err
		}
	case READ_FROM_BUFFER:
		err := bWrite_name(b, m, element)
		if err != nil {
			return err
		}
		b.Write_3x0(op_array_index[m], name, index, element)
	}

	b.Write_1x0(op_counter_increment, index)

	loop := Instr{Kind: op_loop_fixed, Integer: e.Count, Block: b.Block}
	b.Block = old_block
	b.Write(loop)
	return nil
}
func bWrite_slice(b *fnBuilder, m ir_writer_mode, name Name, e TypeSlice) error {
	index := b.NextT(TypeP_usize)
	length := b.NextT(TypeP_usize)
	element := b.Next(e.Type)

	switch m {
	case WRITE_TO_BUFFER:
		b.Write_1x1(op_slice_length, name, length)
		err := bWrite_name(b, m, length)
		if err != nil {
			return err
		}
	case READ_FROM_BUFFER:
		err := bWrite_name(b, m, length)
		if err != nil {
			return err
		}
		b.Write_1x1(op_slice_make, length, name)
	}
	b.Write_0x1(op_counter_init, index)

	old_block := b.Block
	b.Block = []Instr{}

	switch m {
	case WRITE_TO_BUFFER:
		b.Write_2x1(op_slice_index[m], name, index, element)
		err := bWrite_name(b, m, element)
		if err != nil {
			return err
		}
	case READ_FROM_BUFFER:
		err := bWrite_name(b, m, element)
		if err != nil {
			return err
		}
		b.Write_3x0(op_slice_index[m], name, index, element)
	}

	b.Write_1x0(op_counter_increment, index)

	loop := Instr{Kind: op_loop_until, Input: []Name{index, length}, Block: b.Block}
	b.Block = old_block
	b.Write(loop)
	return nil
}
func bWrite_name(b *fnBuilder, mode ir_writer_mode, input Name) error {
	body, err := b.Body(input.Type)
	if err != nil {
		return err
	}

	for {
		if ident, ok := body.(TypeIdent); ok {
			body, err = b.Body(ident.Type)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	switch e := body.(type) {
	case TypePrimitive:
		return bWrite_primitive(b, mode, input, e)
	case TypeStruct:
		return bWrite_struct(b, mode, input, e)
	case TypeEnum:
		return bWrite_enum(b, mode, input, e)
	case TypeArray:
		return bWrite_array(b, mode, input, e)
	case TypeSlice:
		return bWrite_slice(b, mode, input, e)
	case TypePointer:
		return fmt.Errorf("sending pointers across the wire is unsafe")
	case TypeIdent:
		return fmt.Errorf("this should be unreachable, if it was reached, identifier clensing didnt work")
	case TypeAnonymous:
		return fmt.Errorf("cannot write an anonymous identifier to the buffer")
	}

	return fmt.Errorf("unknown kind (%s)", body.Kind())
}
