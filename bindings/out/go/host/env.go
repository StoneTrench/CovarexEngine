// module Env
package generated

import (
	"encoding/binary"
	"math"
	"context"
	"github.com/tetratelabs/wazero/api"
	"fmt"
)

type HostBinding struct {
	Name string
	Method any
}

type BlockKind uint8
const E_Solid = 0
const E_Air = 1
const E_Transparent = 2
const E_Liquid = 3
type ColorU8 struct{ R uint8; G uint8; B uint8; A uint8 }
type BlockType struct{ Name string; Color ColorU8; Kind BlockKind }

func HostAPILogInfo(r9 func(mod api.Module, msg string)())(HostBinding){
	r10 := HostBinding {
			Name: "log-info",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	msg := string(r11)
	
	r9(mod ,msg)
},
		}
	return r10
}
func HostAPILogWarn(r9 func(mod api.Module, msg string)())(HostBinding){
	r10 := HostBinding {
			Name: "log-warn",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	msg := string(r11)
	
	r9(mod ,msg)
},
		}
	return r10
}
func HostAPILogError(r9 func(mod api.Module, msg string)())(HostBinding){
	r10 := HostBinding {
			Name: "log-error",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	msg := string(r11)
	
	r9(mod ,msg)
},
		}
	return r10
}
func HostAPIRegAddBlock(r9 func(mod api.Module, block BlockType)(namespaceid string))(HostBinding){
	r10 := HostBinding {
			Name: "reg-add-block",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32, r7 uint32, r8 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	var block BlockType
	r16 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r14 := r16
	r12 := make([]uint8,r14)
	var r13 uint32 = 0
	for r13 < r14 {
	r15 := r1[0]
		r1 = r1[1:]
	r12[r13] = r15
	r13 += 1
}
	r11 := string(r12)
	block.Name = r11
	var r17 ColorU8
	r18 := r1[0]
		r1 = r1[1:]
	r17.R = r18
	r19 := r1[0]
		r1 = r1[1:]
	r17.G = r19
	r20 := r1[0]
		r1 = r1[1:]
	r17.B = r20
	r21 := r1[0]
		r1 = r1[1:]
	r17.A = r21
	block.Color = r17
	r23 := r1[0]
		r1 = r1[1:]
	r22 := (BlockKind)(r23)
	block.Kind = r22
	
	namespaceid:=r9(mod ,block)
	r2 := make([]byte,0,15)
	
	r24 := []byte(namespaceid)
	r26 := uint32(len(r24))
	r28 := r26
	r2 = binary.LittleEndian.AppendUint32(r2, r28)
	var r25 uint32 = 0
	for r25 < r26 {
	r27 := r24[r25]
	r2 = append(r2, r27)
	r25 += 1
}
	var r5 uint32
		__len :=len(r2)
		if __len > 0 {
			__res, __err := mod.ExportedFunction("malloc").Call(ctx, uint64(__len))
			if __err != nil { panic(__err) }
			r5 = uint32(__res[0])
			mod.Memory().Write(r5, r2)
		}
	r6 := uint32(len(r2))
	mod.Memory().WriteUint32Le(r7, r5)
	mod.Memory().WriteUint32Le(r8, r6)
	// Free buffer r2
},
		}
	return r10
}
func HostAPIRegGetBlock(r9 func(mod api.Module, namespaceid string)(block BlockType, err string))(HostBinding){
	r10 := HostBinding {
			Name: "reg-get-block",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32, r7 uint32, r8 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	namespaceid := string(r11)
	
	block, err:=r9(mod ,namespaceid)
	r2 := make([]byte,0,20)
	
	r16 := block.Name
	r17 := []byte(r16)
	r19 := uint32(len(r17))
	r21 := r19
	r2 = binary.LittleEndian.AppendUint32(r2, r21)
	var r18 uint32 = 0
	for r18 < r19 {
	r20 := r17[r18]
	r2 = append(r2, r20)
	r18 += 1
}
	r22 := block.Color
	r23 := r22.R
	r2 = append(r2, r23)
	r24 := r22.G
	r2 = append(r2, r24)
	r25 := r22.B
	r2 = append(r2, r25)
	r26 := r22.A
	r2 = append(r2, r26)
	r27 := block.Kind
	r28 := uint8(r27)
	r2 = append(r2, r28)
	r29 := []byte(err)
	r31 := uint32(len(r29))
	r33 := r31
	r2 = binary.LittleEndian.AppendUint32(r2, r33)
	var r30 uint32 = 0
	for r30 < r31 {
	r32 := r29[r30]
	r2 = append(r2, r32)
	r30 += 1
}
	var r5 uint32
		__len :=len(r2)
		if __len > 0 {
			__res, __err := mod.ExportedFunction("malloc").Call(ctx, uint64(__len))
			if __err != nil { panic(__err) }
			r5 = uint32(__res[0])
			mod.Memory().Write(r5, r2)
		}
	r6 := uint32(len(r2))
	mod.Memory().WriteUint32Le(r7, r5)
	mod.Memory().WriteUint32Le(r8, r6)
	// Free buffer r2
},
		}
	return r10
}
func HostAPIAppendLang(r9 func(mod api.Module, name string)(err string))(HostBinding){
	r10 := HostBinding {
			Name: "append-lang",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32, r7 uint32, r8 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	name := string(r11)
	
	err:=r9(mod ,name)
	r2 := make([]byte,0,10)
	
	r16 := []byte(err)
	r18 := uint32(len(r16))
	r20 := r18
	r2 = binary.LittleEndian.AppendUint32(r2, r20)
	var r17 uint32 = 0
	for r17 < r18 {
	r19 := r16[r17]
	r2 = append(r2, r19)
	r17 += 1
}
	var r5 uint32
		__len :=len(r2)
		if __len > 0 {
			__res, __err := mod.ExportedFunction("malloc").Call(ctx, uint64(__len))
			if __err != nil { panic(__err) }
			r5 = uint32(__res[0])
			mod.Memory().Write(r5, r2)
		}
	r6 := uint32(len(r2))
	mod.Memory().WriteUint32Le(r7, r5)
	mod.Memory().WriteUint32Le(r8, r6)
	// Free buffer r2
},
		}
	return r10
}
func HostAPIConfigSetInt(r9 func(mod api.Module, name string, value int64)())(HostBinding){
	r10 := HostBinding {
			Name: "config-set-int",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	name := string(r11)
	value := int64(binary.LittleEndian.Uint64(r1))
		r1 = r1[8:]
	
	r9(mod ,name,value)
},
		}
	return r10
}
func HostAPIConfigSetFloat(r9 func(mod api.Module, name string, value float64)())(HostBinding){
	r10 := HostBinding {
			Name: "config-set-float",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	name := string(r11)
	value := math.Float64frombits(binary.LittleEndian.Uint64(r1))
		r1 = r1[8:]
	
	r9(mod ,name,value)
},
		}
	return r10
}
func HostAPIConfigSetString(r9 func(mod api.Module, name string, value string)())(HostBinding){
	r10 := HostBinding {
			Name: "config-set-string",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	name := string(r11)
	r20 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r18 := r20
	r16 := make([]uint8,r18)
	var r17 uint32 = 0
	for r17 < r18 {
	r19 := r1[0]
		r1 = r1[1:]
	r16[r17] = r19
	r17 += 1
}
	value := string(r16)
	
	r9(mod ,name,value)
},
		}
	return r10
}
func HostAPIConfigGetInt(r9 func(mod api.Module, name string)(value int64))(HostBinding){
	r10 := HostBinding {
			Name: "config-get-int",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32, r7 uint32, r8 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	name := string(r11)
	
	value:=r9(mod ,name)
	r2 := make([]byte,0,13)
	
	r2 = binary.LittleEndian.AppendUint64(r2, uint64(value))
	var r5 uint32
		__len :=len(r2)
		if __len > 0 {
			__res, __err := mod.ExportedFunction("malloc").Call(ctx, uint64(__len))
			if __err != nil { panic(__err) }
			r5 = uint32(__res[0])
			mod.Memory().Write(r5, r2)
		}
	r6 := uint32(len(r2))
	mod.Memory().WriteUint32Le(r7, r5)
	mod.Memory().WriteUint32Le(r8, r6)
	// Free buffer r2
},
		}
	return r10
}
func HostAPIConfigGetFloat(r9 func(mod api.Module, name string)(value float64))(HostBinding){
	r10 := HostBinding {
			Name: "config-get-float",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32, r7 uint32, r8 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	name := string(r11)
	
	value:=r9(mod ,name)
	r2 := make([]byte,0,13)
	
	r2 = binary.LittleEndian.AppendUint64(r2, math.Float64bits(value))
	var r5 uint32
		__len :=len(r2)
		if __len > 0 {
			__res, __err := mod.ExportedFunction("malloc").Call(ctx, uint64(__len))
			if __err != nil { panic(__err) }
			r5 = uint32(__res[0])
			mod.Memory().Write(r5, r2)
		}
	r6 := uint32(len(r2))
	mod.Memory().WriteUint32Le(r7, r5)
	mod.Memory().WriteUint32Le(r8, r6)
	// Free buffer r2
},
		}
	return r10
}
func HostAPIConfigGetString(r9 func(mod api.Module, name string)(value string))(HostBinding){
	r10 := HostBinding {
			Name: "config-get-string",
			Method: func(ctx context.Context, mod api.Module, r3 uint32, r4 uint32, r7 uint32, r8 uint32) {
	r1, __ok := mod.Memory().Read(r3, r4)		
		if !__ok { panic(fmt.Errorf("failed to copy buffer from guest [%d %d]", r3, r4)) }
	
	r15 := binary.LittleEndian.Uint32(r1)
		r1 = r1[4:]
	r13 := r15
	r11 := make([]uint8,r13)
	var r12 uint32 = 0
	for r12 < r13 {
	r14 := r1[0]
		r1 = r1[1:]
	r11[r12] = r14
	r12 += 1
}
	name := string(r11)
	
	value:=r9(mod ,name)
	r2 := make([]byte,0,10)
	
	r16 := []byte(value)
	r18 := uint32(len(r16))
	r20 := r18
	r2 = binary.LittleEndian.AppendUint32(r2, r20)
	var r17 uint32 = 0
	for r17 < r18 {
	r19 := r16[r17]
	r2 = append(r2, r19)
	r17 += 1
}
	var r5 uint32
		__len :=len(r2)
		if __len > 0 {
			__res, __err := mod.ExportedFunction("malloc").Call(ctx, uint64(__len))
			if __err != nil { panic(__err) }
			r5 = uint32(__res[0])
			mod.Memory().Write(r5, r2)
		}
	r6 := uint32(len(r2))
	mod.Memory().WriteUint32Le(r7, r5)
	mod.Memory().WriteUint32Le(r8, r6)
	// Free buffer r2
},
		}
	return r10
}
