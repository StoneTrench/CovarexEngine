package api
import (
"unsafe"
)
type Color struct{R uint8;G uint8;B uint8;A uint8}
type BlockType struct{Name string;Color Color;IsSolid bool}
type WError [1]string
type OptError *WError
//go:wasmimport env log-info
func env_LogInfo(msg_str uint32,msg_len uint32);
func LogInfo(msg string){
// Guest Encode
var msg_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(msg))));
var msg_len = (uint32)(len(msg));

// Declare Outputs
// Call
env_LogInfo(msg_str,msg_len);
// Guest Decode
// Finish
}
//go:wasmimport env log-warn
func env_LogWarn(msg_str uint32,msg_len uint32);
func LogWarn(msg string){
// Guest Encode
var msg_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(msg))));
var msg_len = (uint32)(len(msg));

// Declare Outputs
// Call
env_LogWarn(msg_str,msg_len);
// Guest Decode
// Finish
}
//go:wasmimport env log-error
func env_LogError(msg_str uint32,msg_len uint32);
func LogError(msg string){
// Guest Encode
var msg_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(msg))));
var msg_len = (uint32)(len(msg));

// Declare Outputs
// Call
env_LogError(msg_str,msg_len);
// Guest Decode
// Finish
}
//go:wasmimport env reg-add-block
func env_RegAddBlock(block_Name_f_str uint32,block_Name_f_len uint32,block_IsSolid_f_b uint32,block_Color_f_R_f_w uint32,block_Color_f_G_f_w uint32,block_Color_f_B_f_w uint32,block_Color_f_A_f_w uint32,namespaceid_str uint32,namespaceid_len uint32);
func RegAddBlock(block BlockType)(string){
// Guest Encode
var block_Name_f = (string)(block.Name);var block_Name_f_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(block_Name_f))));
var block_Name_f_len = (uint32)(len(block_Name_f));
var block_Color_f = (Color)(block.Color);var block_Color_f_R_f = (uint8)(block_Color_f.R);var block_Color_f_R_f_w = (uint32)(block_Color_f_R_f);var block_Color_f_G_f = (uint8)(block_Color_f.G);var block_Color_f_G_f_w = (uint32)(block_Color_f_G_f);var block_Color_f_B_f = (uint8)(block_Color_f.B);var block_Color_f_B_f_w = (uint32)(block_Color_f_B_f);var block_Color_f_A_f = (uint8)(block_Color_f.A);var block_Color_f_A_f_w = (uint32)(block_Color_f_A_f);var block_IsSolid_f = (bool)(block.IsSolid);var block_IsSolid_f_b uint32;
if block_IsSolid_f{block_IsSolid_f_b = 1;}else{block_IsSolid_f_b = 0;};

// Declare Outputs
var namespaceid_str uint32;
var namespaceid_len uint32;
// Call
env_RegAddBlock(block_Name_f_str,block_Name_f_len,block_IsSolid_f_b,block_Color_f_R_f_w,block_Color_f_G_f_w,block_Color_f_B_f_w,block_Color_f_A_f_w,(uint32)(uintptr(unsafe.Pointer(&namespaceid_str))),(uint32)(uintptr(unsafe.Pointer(&namespaceid_len))));
// Guest Decode
var namespaceid = (string)(unsafe.String((*byte)(unsafe.Pointer(uintptr(namespaceid_str))),int(namespaceid_len)));

// Finish
uResetMalloc();
return namespaceid;
}
//go:wasmimport env reg-get-block
func env_RegGetBlock(namespaceid_str uint32,namespaceid_len uint32,err_ptr uint32,block_Name_f_str uint32,block_Name_f_len uint32,block_IsSolid_f_b uint32,block_Color_f_R_f_w uint32,block_Color_f_G_f_w uint32,block_Color_f_B_f_w uint32,block_Color_f_A_f_w uint32);
func RegGetBlock(namespaceid string)(BlockType,OptError){
// Guest Encode
var namespaceid_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(namespaceid))));
var namespaceid_len = (uint32)(len(namespaceid));

// Declare Outputs
var err_ptr uint32;
var block_Name_f_str uint32;
var block_Name_f_len uint32;
var block_IsSolid_f_b uint32;
var block_Color_f_R_f_w uint32;
var block_Color_f_G_f_w uint32;
var block_Color_f_B_f_w uint32;
var block_Color_f_A_f_w uint32;
// Call
env_RegGetBlock(namespaceid_str,namespaceid_len,(uint32)(uintptr(unsafe.Pointer(&err_ptr))),(uint32)(uintptr(unsafe.Pointer(&block_Name_f_str))),(uint32)(uintptr(unsafe.Pointer(&block_Name_f_len))),(uint32)(uintptr(unsafe.Pointer(&block_IsSolid_f_b))),(uint32)(uintptr(unsafe.Pointer(&block_Color_f_R_f_w))),(uint32)(uintptr(unsafe.Pointer(&block_Color_f_G_f_w))),(uint32)(uintptr(unsafe.Pointer(&block_Color_f_B_f_w))),(uint32)(uintptr(unsafe.Pointer(&block_Color_f_A_f_w))));
// Guest Decode
var block struct{Name string;Color Color;IsSolid bool};
var block_Name_f = (string)(unsafe.String((*byte)(unsafe.Pointer(uintptr(block_Name_f_str))),int(block_Name_f_len)));
block.Name=block_Name_f;
var block_Color_f struct{R uint8;G uint8;B uint8;A uint8};
var block_Color_f_R_f = (uint8)(block_Color_f_R_f_w);
block_Color_f.R=block_Color_f_R_f;
var block_Color_f_G_f = (uint8)(block_Color_f_G_f_w);
block_Color_f.G=block_Color_f_G_f;
var block_Color_f_B_f = (uint8)(block_Color_f_B_f_w);
block_Color_f.B=block_Color_f_B_f;
var block_Color_f_A_f = (uint8)(block_Color_f_A_f_w);
block_Color_f.A=block_Color_f_A_f;
block.Color=block_Color_f;
var block_IsSolid_f = (bool)(block_IsSolid_f_b!=0);
block.IsSolid=block_IsSolid_f;

var err=*(**WError)(unsafe.Pointer(uintptr(err_ptr)));

// Finish
uResetMalloc();
return block,err;
}
//go:wasmimport env append-lang
func env_AppendLang(name_str uint32,name_len uint32,err_ptr uint32);
func AppendLang(name string)(OptError){
// Guest Encode
var name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))));
var name_len = (uint32)(len(name));

// Declare Outputs
var err_ptr uint32;
// Call
env_AppendLang(name_str,name_len,(uint32)(uintptr(unsafe.Pointer(&err_ptr))));
// Guest Decode
var err=*(**WError)(unsafe.Pointer(uintptr(err_ptr)));

// Finish
uResetMalloc();
return err;
}
//go:wasmimport env config-set-int
func env_ConfigSetInt(name_str uint32,name_len uint32,value_w int64);
func ConfigSetInt(name string,value int64){
// Guest Encode
var name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))));
var name_len = (uint32)(len(name));

var value_w = (int64)(value);
// Declare Outputs
// Call
env_ConfigSetInt(name_str,name_len,value_w);
// Guest Decode
// Finish
}
//go:wasmimport env config-set-float
func env_ConfigSetFloat(name_str uint32,name_len uint32,value_w float64);
func ConfigSetFloat(name string,value float64){
// Guest Encode
var name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))));
var name_len = (uint32)(len(name));

var value_w = (float64)(value);
// Declare Outputs
// Call
env_ConfigSetFloat(name_str,name_len,value_w);
// Guest Decode
// Finish
}
//go:wasmimport env config-set-string
func env_ConfigSetString(name_str uint32,name_len uint32,value_str uint32,value_len uint32);
func ConfigSetString(name string,value string){
// Guest Encode
var name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))));
var name_len = (uint32)(len(name));

var value_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(value))));
var value_len = (uint32)(len(value));

// Declare Outputs
// Call
env_ConfigSetString(name_str,name_len,value_str,value_len);
// Guest Decode
// Finish
}
//go:wasmimport env config-get-int
func env_ConfigGetInt(name_str uint32,name_len uint32,value_w uint32);
func ConfigGetInt(name string)(int64){
// Guest Encode
var name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))));
var name_len = (uint32)(len(name));

// Declare Outputs
var value_w uint32;
// Call
env_ConfigGetInt(name_str,name_len,(uint32)(uintptr(unsafe.Pointer(&value_w))));
// Guest Decode
var value = (int64)(value_w);

// Finish
return value;
}
//go:wasmimport env config-get-float
func env_ConfigGetFloat(name_str uint32,name_len uint32,value_w uint32);
func ConfigGetFloat(name string)(float64){
// Guest Encode
var name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))));
var name_len = (uint32)(len(name));

// Declare Outputs
var value_w uint32;
// Call
env_ConfigGetFloat(name_str,name_len,(uint32)(uintptr(unsafe.Pointer(&value_w))));
// Guest Decode
var value = (float64)(value_w);

// Finish
return value;
}
//go:wasmimport env config-get-string
func env_ConfigGetString(name_str uint32,name_len uint32,value_str uint32,value_len uint32);
func ConfigGetString(name string)(string){
// Guest Encode
var name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))));
var name_len = (uint32)(len(name));

// Declare Outputs
var value_str uint32;
var value_len uint32;
// Call
env_ConfigGetString(name_str,name_len,(uint32)(uintptr(unsafe.Pointer(&value_str))),(uint32)(uintptr(unsafe.Pointer(&value_len))));
// Guest Decode
var value = (string)(unsafe.String((*byte)(unsafe.Pointer(uintptr(value_str))),int(value_len)));

// Finish
uResetMalloc();
return value;
}
