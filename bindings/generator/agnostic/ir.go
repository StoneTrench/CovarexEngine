package model

import types "github.com/StoneTrench/go-mc-clone/bindings/generator/model"

type TypeRef = types.TypeRef

/*
These are the types that the template will process to create the binding sourcecode in the target language.

Assume that types are all packed (no padding in struct).
*/

type Module struct {
	Name  string
	Funcs []FuncDecl
}

type FuncDecl struct {
	Module  *Module
	Name    string
	Params  []FuncParam
	RetType TypeRef
	Block   []Instr
}

type FuncParam struct {
	Name string
	Type TypeRef
}

type InstrParam struct {
	Name string
	Type TypeRef
}

type Instr struct {
	Kind string

	Count  int // for the fixed loop
	Input  []InstrParam
	Output []InstrParam
	Block  []Instr
}

type InstrKind string

const (
	Instr_Invalid InstrKind = "Invalid"

	Instr_LoopFixed   InstrKind = "LoopFixed"
	Instr_LoopDynamic InstrKind = "LoopDynamic"

	Instr_
)
