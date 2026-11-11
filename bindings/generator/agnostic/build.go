package model

import (
	"fmt"
	"strings"

	types "github.com/StoneTrench/go-mc-clone/bindings/generator/model"
)

type NameCounter uint

func (n *NameCounter) Next() string {
	(*n) = (*n) + 1
	return n.Curr()
}
func (n *NameCounter) Curr() string {
	return fmt.Sprintf("_r%d", *n)
}

type Translation struct {
	Body    []Instr
	OutName string
	OutType types.TypeRef
}

func BuildAgnosticModule(c *types.TypeCollection, name string, fields []types.TypeField) (*Module, error) {
	mod := new(Module)

	mod.Name = exported(name)

	for _, field := range fields {
		body, err := c.Body(field.Type)
		if err != nil {
			return nil, err
		}
		if type_fn, ok := body.(types.TypeFunc); ok {
			fun, err := decl_func(c, mod, field.Name, &type_fn)
			if err != nil {
				return nil, err
			}
			mod.Funcs = append(mod.Funcs, *fun)
		}
	}

	return mod, nil
}
func decl_func(c *types.TypeCollection, mod *Module, name string, fn *types.TypeFunc) (*FuncDecl, error) {
	res := new(FuncDecl)
	res.Module = mod
	res.Name = exported(name)

	names := new(NameCounter)

	for _, param := range fn.Input {
		res.Params = append(res.Params, FuncParam{
			Name: param.Name,
			Type: param.Type,
		})

		// res.Body = append(res.Body, Instr{Kind: "DeclareMove", Inp1: param.Name, Out1: names.Next()})
		instrs, err := translate(c, names, &Translation{
			Body:    []Instr{},
			OutName: param.Name,
			OutType: param.Type,
		})
		if err != nil {
			return nil, err
		}
		if instrs != nil {
			res.Block = append(res.Block, instrs.Body...)
		}
	}
	if len(fn.Output) > 0 {
		res.RetType = fn.Output[0].Type
	}

	return res, nil
}
func translate(c *types.TypeCollection, n *NameCounter, r *Translation) (*Translation, error) {
	ARG_BODY, err := c.Body(r.OutType)
	if err != nil {
		return nil, err
	}

	ARG_NAME := r.OutName

	switch e := ARG_BODY.(type) {
	case types.TypePrimitive:
		switch e {
		case types.TypeP_Invalid:
			return nil, fmt.Errorf("invalid primitive type")
		case types.TypeP_string:
			F_ptr := n.Next()
			F_len := n.Next()
			OUT_NAME := n.Next()
			OUT_TYPE := types.Type_bodyof_string(c)
			OUT := &Translation{
				Body: append(r.Body,
					Instr{Kind: "UnpackString", Output: []InstrParam{
						{ Name: F_ptr, Type: c.Ref(Types.TypeP) }
					} Out1: F_ptr, Out2: F_len, Inp1: ARG_NAME},
					Instr{Kind: "Declare", Out1: OUT_NAME, Typ1: OUT_TYPE},
					Instr{Kind: "SetField", Out1: OUT_NAME, Out2: "ptr", Inp1: F_ptr},
					Instr{Kind: "SetField", Out1: OUT_NAME, Out2: "len", Inp1: F_len},
				),
				OutName: OUT_NAME,
				OutType: OUT_TYPE,
			}
			return translate(c, n, OUT)
		}
	case types.TypeStruct:
		OUT_NAME := n.Next()
		OUT_INSTR := []Instr{}
		OUT_Struct := types.TypeStruct{}

		for _, f := range e.Fields {
			field, err := translate(c, n, &Translation{
				Body:    nil,
				OutName: f.Name,
				OutType: f.Type,
			})
			if err != nil {
				return nil, err
			}
			if field != nil {
				OUT_INSTR = append(OUT_INSTR, field.Body...)
				OUT_INSTR = append(OUT_INSTR, Instr{
					Kind: "SetField", Out1: OUT_NAME, Out2: f.Name, Inp1: field.OutName,
				})
				OUT_Struct.Fields = append(OUT_Struct.Fields, types.TypeField{
					Name: f.Name,
					Type: field.OutType,
				})
			}
		}

		OUT_TYPE := c.Ref(OUT_Struct)

		return &Translation{
			Body: append(append(r.Body, Instr{Kind: "Declare", Out1: OUT_NAME, Typ1: OUT_TYPE}),
				OUT_INSTR...,
			),
			OutName: OUT_NAME,
			OutType: OUT_TYPE,
		}, nil
	}
	return r, nil
}

// exported converts an .abi identifier to an exported identifier:
//   - '-' removed, next rune capitalized
//   - '_' preserved, next rune capitalized
//   - first rune forced uppercase
func exported(s string) string {
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
