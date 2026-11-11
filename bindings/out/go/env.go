package Env;

import "unsafe"


// go:wasmimport Env log-info
func _Env_LogInfo()
func LogInfo(msg string){
	_r1 := (*uint8)(unsafe.StringData(msg))
	_r2 := uint32(len(msg))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env log-warn
func _Env_LogWarn()
func LogWarn(msg string){
	_r1 := (*uint8)(unsafe.StringData(msg))
	_r2 := uint32(len(msg))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env log-error
func _Env_LogError()
func LogError(msg string){
	_r1 := (*uint8)(unsafe.StringData(msg))
	_r2 := uint32(len(msg))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env reg-add-block
func _Env_RegAddBlock()
func RegAddBlock(block BlockType)string{
}
// go:wasmimport Env reg-get-block
func _Env_RegGetBlock()
func RegGetBlock(namespaceid string)BlockType{
	_r1 := (*uint8)(unsafe.StringData(namespaceid))
	_r2 := uint32(len(namespaceid))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env append-lang
func _Env_AppendLang()
func AppendLang(name string){
	_r1 := (*uint8)(unsafe.StringData(name))
	_r2 := uint32(len(name))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env config-set-int
func _Env_ConfigSetInt()
func ConfigSetInt(name string, value int64){
	_r1 := (*uint8)(unsafe.StringData(name))
	_r2 := uint32(len(name))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env config-set-float
func _Env_ConfigSetFloat()
func ConfigSetFloat(name string, value float64){
	_r1 := (*uint8)(unsafe.StringData(name))
	_r2 := uint32(len(name))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env config-set-string
func _Env_ConfigSetString()
func ConfigSetString(name string, value string){
	_r1 := (*uint8)(unsafe.StringData(name))
	_r2 := uint32(len(name))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
	_r4 := (*uint8)(unsafe.StringData(value))
	_r5 := uint32(len(value))
	var _r6 wasm_string;
	_r6.ptr = _r4
	_r6.len = _r5
}
// go:wasmimport Env config-get-int
func _Env_ConfigGetInt()
func ConfigGetInt(name string)int64{
	_r1 := (*uint8)(unsafe.StringData(name))
	_r2 := uint32(len(name))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env config-get-float
func _Env_ConfigGetFloat()
func ConfigGetFloat(name string)float64{
	_r1 := (*uint8)(unsafe.StringData(name))
	_r2 := uint32(len(name))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}
// go:wasmimport Env config-get-string
func _Env_ConfigGetString()
func ConfigGetString(name string)string{
	_r1 := (*uint8)(unsafe.StringData(name))
	_r2 := uint32(len(name))
	var _r3 wasm_string;
	_r3.ptr = _r1
	_r3.len = _r2
}