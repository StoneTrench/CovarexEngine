package intermediate_representation

import (
	"fmt"

	types "github.com/StoneTrench/go-mc-clone/bindings/generator/model"
)

type fnBuilder struct {
	TypColl       *TColl
	NameC         *NameCount
	Func          *FuncDecl
	Block         []Instr
	Buffer        Name
	BufferMinSize int
}

func (b *fnBuilder) Ref(t types.IType) TYPE_REF {
	return b.TypColl.Ref(t)
}
func (b *fnBuilder) Body(t TYPE_REF) (types.IType, error) {
	return b.TypColl.Body(t)
}
func (b *fnBuilder) Next(t TYPE_REF) Name {
	return b.NameC.Next(t)
}
func (b *fnBuilder) NextT(t types.IType) Name {
	return b.Next(b.Ref(t))
}

func (b *fnBuilder) Write(i Instr) *Instr {
	i.Func = b.Func
	i.Index = len(b.Block)
	b.Block = append(b.Block, i)
	return &b.Block[len(b.Block)-1]
}
func (b *fnBuilder) WriteBuf(mode ir_writer_mode, kind string, input Name) *Instr {
	b.BufferMinSize += sizeof(b.TypColl, input.Type)
	switch mode {
	case WRITE_TO_BUFFER:
		return b.Write(Instr{Kind: kind, Input: []Name{b.Buffer, input}, Output: []Name{}})
	case READ_FROM_BUFFER:
		return b.Write(Instr{Kind: kind, Input: []Name{b.Buffer}, Output: []Name{input}})
	}
	return nil
}
func (b *fnBuilder) Write_0x1(kind string, output Name) *Instr {
	return b.Write(Instr{Kind: kind, Input: []Name{}, Output: []Name{output}})
}
func (b *fnBuilder) Write_1x0(kind string, input Name) *Instr {
	return b.Write(Instr{Kind: kind, Input: []Name{input}, Output: []Name{}})
}
func (b *fnBuilder) Write_1x1(kind string, input Name, output Name) *Instr {
	return b.Write(Instr{Kind: kind, Input: []Name{input}, Output: []Name{output}})
}
func (b *fnBuilder) Write_2x0(kind string, input_1 Name, input_2 Name) *Instr {
	return b.Write(Instr{Kind: kind, Input: []Name{input_1, input_2}, Output: []Name{}})
}
func (b *fnBuilder) Write_2x1(kind string, input_1 Name, input_2 Name, output Name) *Instr {
	return b.Write(Instr{Kind: kind, Input: []Name{input_1, input_2}, Output: []Name{output}})
}
func (b *fnBuilder) Write_3x0(kind string, input_1 Name, input_2 Name, input_3 Name) *Instr {
	return b.Write(Instr{Kind: kind, Input: []Name{input_1, input_2, input_3}, Output: []Name{}})
}
func (b *fnBuilder) Write_nop(message string) *Instr {
	return b.Write(Instr{Kind: op_nop_comment, Text: message})
	// return nil
}

func (b *fnBuilder) MakeBuffer() *Instr {
	return b.Write_0x1(op_buffer_make, b.Buffer)
}
func (b *fnBuilder) MakeBufferFromRaw(p, l Name) *Instr {
	return b.Write_2x1(op_buffer_make_from_raw, p, l, b.Buffer)
}
func (b *fnBuilder) FreeBuffer() *Instr {
	return b.Write_1x0(op_buffer_free, b.Buffer)
}
func (b *fnBuilder) BufferPointer(out Name) *Instr {
	return b.Write_1x1(op_buffer_pointer, b.Buffer, out)
}
func (b *fnBuilder) BufferLength(out Name) *Instr {
	return b.Write_1x1(op_buffer_length, b.Buffer, out)
}
func (b *fnBuilder) WriteMemoryU32(addr, val Name) *Instr {
	return b.Write_2x0(op_host_write_memory_u32, addr, val)
}
func (b *fnBuilder) NopComment(text string) *Instr {
	return b.Write_nop(text)
}

func (b *fnBuilder) Serialize(m ir_writer_mode, names []Name) []error {
	err_list := []error{}
	for _, v := range names {
		if err := bWrite_name(b, m, v); err != nil {
			err_list = append(err_list, fmt.Errorf("failed to serialize (%s): %w", v.Name, err))
			continue
		}
	}
	return err_list
}

func sizeof(c *TColl, p TYPE_REF) int {
	body, err := c.Body(p)
	if err != nil {
		return 0
	}

	if p, ok := body.(types.TypePrimitive); ok {
		switch p {
		case types.TypeP_u8:
			return 1
		case types.TypeP_u16:
			return 2
		case types.TypeP_u32:
			return 4
		case types.TypeP_u64:
			return 8

		case types.TypeP_i8:
			return 1
		case types.TypeP_i16:
			return 2
		case types.TypeP_i32:
			return 4
		case types.TypeP_i64:
			return 8

		case types.TypeP_f32:
			return 4
		case types.TypeP_f64:
			return 8

		case types.TypeP_usize:
			return 4
		}
	}
	return 0
}
