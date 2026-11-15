//go:build wasip1
// module Env
package generated

import (
	"encoding/binary"
	"unsafe"
	"math"
	"sync"
)

var (
	allocMu sync.Mutex
	allocs  = make(map[uint32][]byte)
)

//go:wasmexport malloc
func wasm_malloc(n uint32) uint32 {
	if n == 0 {
		return 0
	}
	buf := make([]byte, n)
	ptr := uint32(uintptr(unsafe.Pointer(unsafe.SliceData(buf))))
	allocMu.Lock()
	allocs[ptr] = buf
	allocMu.Unlock()
	return ptr
}

func wasm_free(ptr uint32) {
	if ptr == 0 {
		return
	}
	allocMu.Lock()
	delete(allocs, ptr)
	allocMu.Unlock()
}

type BlockKind uint8
const E_Solid = 0
const E_Air = 1
const E_Transparent = 2
const E_Liquid = 3
type ColorU8 struct{ R uint8; G uint8; B uint8; A uint8 }
type BlockType struct{ Name string; Color ColorU8; Kind BlockKind }


//go:wasmimport env log-info
func _Env_LogInfo(r3 uint32, r4 uint32)
func LogInfo(msg string)(){
	r1 := make([]byte,0,5)
	
	r7 := []byte(msg)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	_Env_LogInfo(r3, r4)
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
}
//go:wasmimport env log-warn
func _Env_LogWarn(r3 uint32, r4 uint32)
func LogWarn(msg string)(){
	r1 := make([]byte,0,5)
	
	r7 := []byte(msg)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	_Env_LogWarn(r3, r4)
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
}
//go:wasmimport env log-error
func _Env_LogError(r3 uint32, r4 uint32)
func LogError(msg string)(){
	r1 := make([]byte,0,5)
	
	r7 := []byte(msg)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	_Env_LogError(r3, r4)
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
}
//go:wasmimport env reg-add-block
func _Env_RegAddBlock(r3 uint32, r4 uint32, r5 uint32, r6 uint32)
func RegAddBlock(block BlockType)(string){
	r1 := make([]byte,0,10)
	
	r7 := block.Name
	r8 := []byte(r7)
	r10 := uint32(len(r8))
	r12 := r10
	r1 = binary.LittleEndian.AppendUint32(r1, r12)
	var r9 uint32 = 0
	for r9 < r10 {
	r11 := r8[r9]
	r1 = append(r1, r11)
	r9 += 1
}
	r13 := block.Color
	r14 := r13.R
	r1 = append(r1, r14)
	r15 := r13.G
	r1 = append(r1, r15)
	r16 := r13.B
	r1 = append(r1, r16)
	r17 := r13.A
	r1 = append(r1, r17)
	r18 := block.Kind
	r19 := uint8(r18)
	r1 = append(r1, r19)
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	var r5 uint32
	var r6 uint32
	_Env_RegAddBlock(r3, r4, uint32(uintptr(unsafe.Pointer(&r5))), uint32(uintptr(unsafe.Pointer(&r6))))
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
	r2 := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(r5))), int(r6))
	
	r24 := binary.LittleEndian.Uint32(r2)
		r2 = r2[4:]
	r22 := r24
	r20 := make([]uint8,r22)
	var r21 uint32 = 0
	for r21 < r22 {
	r23 := r2[0]
		r2 = r2[1:]
	r20[r21] = r23
	r21 += 1
}
	namespaceid := string(r20)
	if len(r2) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r2[0]))))
		}
	return namespaceid
}
//go:wasmimport env reg-get-block
func _Env_RegGetBlock(r3 uint32, r4 uint32, r5 uint32, r6 uint32)
func RegGetBlock(namespaceid string)(BlockType, string){
	r1 := make([]byte,0,5)
	
	r7 := []byte(namespaceid)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	var r5 uint32
	var r6 uint32
	_Env_RegGetBlock(r3, r4, uint32(uintptr(unsafe.Pointer(&r5))), uint32(uintptr(unsafe.Pointer(&r6))))
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
	r2 := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(r5))), int(r6))
	
	var block BlockType
	r17 := binary.LittleEndian.Uint32(r2)
		r2 = r2[4:]
	r15 := r17
	r13 := make([]uint8,r15)
	var r14 uint32 = 0
	for r14 < r15 {
	r16 := r2[0]
		r2 = r2[1:]
	r13[r14] = r16
	r14 += 1
}
	r12 := string(r13)
	block.Name = r12
	var r18 ColorU8
	r19 := r2[0]
		r2 = r2[1:]
	r18.R = r19
	r20 := r2[0]
		r2 = r2[1:]
	r18.G = r20
	r21 := r2[0]
		r2 = r2[1:]
	r18.B = r21
	r22 := r2[0]
		r2 = r2[1:]
	r18.A = r22
	block.Color = r18
	r24 := r2[0]
		r2 = r2[1:]
	r23 := (BlockKind)(r24)
	block.Kind = r23
	r29 := binary.LittleEndian.Uint32(r2)
		r2 = r2[4:]
	r27 := r29
	r25 := make([]uint8,r27)
	var r26 uint32 = 0
	for r26 < r27 {
	r28 := r2[0]
		r2 = r2[1:]
	r25[r26] = r28
	r26 += 1
}
	err := string(r25)
	if len(r2) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r2[0]))))
		}
	return block, err
}
//go:wasmimport env append-lang
func _Env_AppendLang(r3 uint32, r4 uint32, r5 uint32, r6 uint32)
func AppendLang(name string)(string){
	r1 := make([]byte,0,5)
	
	r7 := []byte(name)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	var r5 uint32
	var r6 uint32
	_Env_AppendLang(r3, r4, uint32(uintptr(unsafe.Pointer(&r5))), uint32(uintptr(unsafe.Pointer(&r6))))
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
	r2 := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(r5))), int(r6))
	
	r16 := binary.LittleEndian.Uint32(r2)
		r2 = r2[4:]
	r14 := r16
	r12 := make([]uint8,r14)
	var r13 uint32 = 0
	for r13 < r14 {
	r15 := r2[0]
		r2 = r2[1:]
	r12[r13] = r15
	r13 += 1
}
	err := string(r12)
	if len(r2) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r2[0]))))
		}
	return err
}
//go:wasmimport env config-set-int
func _Env_ConfigSetInt(r3 uint32, r4 uint32)
func ConfigSetInt(name string, value int64)(){
	r1 := make([]byte,0,13)
	
	r7 := []byte(name)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	r1 = binary.LittleEndian.AppendUint64(r1, uint64(value))
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	_Env_ConfigSetInt(r3, r4)
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
}
//go:wasmimport env config-set-float
func _Env_ConfigSetFloat(r3 uint32, r4 uint32)
func ConfigSetFloat(name string, value float64)(){
	r1 := make([]byte,0,13)
	
	r7 := []byte(name)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	r1 = binary.LittleEndian.AppendUint64(r1, math.Float64bits(value))
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	_Env_ConfigSetFloat(r3, r4)
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
}
//go:wasmimport env config-set-string
func _Env_ConfigSetString(r3 uint32, r4 uint32)
func ConfigSetString(name string, value string)(){
	r1 := make([]byte,0,10)
	
	r7 := []byte(name)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	r12 := []byte(value)
	r14 := uint32(len(r12))
	r16 := r14
	r1 = binary.LittleEndian.AppendUint32(r1, r16)
	var r13 uint32 = 0
	for r13 < r14 {
	r15 := r12[r13]
	r1 = append(r1, r15)
	r13 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	_Env_ConfigSetString(r3, r4)
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
}
//go:wasmimport env config-get-int
func _Env_ConfigGetInt(r3 uint32, r4 uint32, r5 uint32, r6 uint32)
func ConfigGetInt(name string)(int64){
	r1 := make([]byte,0,5)
	
	r7 := []byte(name)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	var r5 uint32
	var r6 uint32
	_Env_ConfigGetInt(r3, r4, uint32(uintptr(unsafe.Pointer(&r5))), uint32(uintptr(unsafe.Pointer(&r6))))
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
	r2 := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(r5))), int(r6))
	
	value := int64(binary.LittleEndian.Uint64(r2))
		r2 = r2[8:]
	if len(r2) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r2[0]))))
		}
	return value
}
//go:wasmimport env config-get-float
func _Env_ConfigGetFloat(r3 uint32, r4 uint32, r5 uint32, r6 uint32)
func ConfigGetFloat(name string)(float64){
	r1 := make([]byte,0,5)
	
	r7 := []byte(name)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	var r5 uint32
	var r6 uint32
	_Env_ConfigGetFloat(r3, r4, uint32(uintptr(unsafe.Pointer(&r5))), uint32(uintptr(unsafe.Pointer(&r6))))
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
	r2 := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(r5))), int(r6))
	
	value := math.Float64frombits(binary.LittleEndian.Uint64(r2))
		r2 = r2[8:]
	if len(r2) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r2[0]))))
		}
	return value
}
//go:wasmimport env config-get-string
func _Env_ConfigGetString(r3 uint32, r4 uint32, r5 uint32, r6 uint32)
func ConfigGetString(name string)(string){
	r1 := make([]byte,0,5)
	
	r7 := []byte(name)
	r9 := uint32(len(r7))
	r11 := r9
	r1 = binary.LittleEndian.AppendUint32(r1, r11)
	var r8 uint32 = 0
	for r8 < r9 {
	r10 := r7[r8]
	r1 = append(r1, r10)
	r8 += 1
}
	
	var r3 uint32 = 0
		if len(r1) > 0 {
			r3 = uint32(uintptr(unsafe.Pointer(&r1[0])))
		}
	r4 := uint32(len(r1))
	var r5 uint32
	var r6 uint32
	_Env_ConfigGetString(r3, r4, uint32(uintptr(unsafe.Pointer(&r5))), uint32(uintptr(unsafe.Pointer(&r6))))
	if len(r1) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r1[0]))))
		}
	r2 := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(r5))), int(r6))
	
	r16 := binary.LittleEndian.Uint32(r2)
		r2 = r2[4:]
	r14 := r16
	r12 := make([]uint8,r14)
	var r13 uint32 = 0
	for r13 < r14 {
	r15 := r2[0]
		r2 = r2[1:]
	r12[r13] = r15
	r13 += 1
}
	value := string(r12)
	if len(r2) > 0 {
			wasm_free(uint32(uintptr(unsafe.Pointer(&r2[0]))))
		}
	return value
}