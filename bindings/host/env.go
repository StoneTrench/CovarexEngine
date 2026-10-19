package host
import "github.com/tetratelabs/wazero/api"
import "context"
type Color struct{R uint8;G uint8;B uint8;A uint8}
type BlockType struct{Name string;Color Color;IsSolid bool}
type WError [1]string
type OptError *WError
func LogInfo(fn func(m api.Module,msg string)())(w HostBinding){
return HostBinding{"log-info",func(ctx context.Context,mod api.Module,msg_str uint32,msg_len uint32)(){
var msg = uReadString(mod, msg_str, msg_len);
fn(mod,msg);},
}
}
func LogWarn(fn func(m api.Module,msg string)())(w HostBinding){
return HostBinding{"log-warn",func(ctx context.Context,mod api.Module,msg_str uint32,msg_len uint32)(){
var msg = uReadString(mod, msg_str, msg_len);
fn(mod,msg);},
}
}
func LogError(fn func(m api.Module,msg string)())(w HostBinding){
return HostBinding{"log-error",func(ctx context.Context,mod api.Module,msg_str uint32,msg_len uint32)(){
var msg = uReadString(mod, msg_str, msg_len);
fn(mod,msg);},
}
}
func RegAddBlock(fn func(m api.Module,block BlockType)(namespaceid string))(w HostBinding){
return HostBinding{"reg-add-block",func(ctx context.Context,mod api.Module,block_Name_f_str uint32,block_Name_f_len uint32,block_IsSolid_f_b uint32,block_Color_f_R_f_w uint32,block_Color_f_G_f_w uint32,block_Color_f_B_f_w uint32,block_Color_f_A_f_w uint32,namespaceid_str uint32,namespaceid_len uint32)(){
_malloc:=mod.ExportedFunction("malloc");alloc:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));if err!=nil{panic(err)};return uint32(res[0])};var block struct{Name string;Color Color;IsSolid bool};
var block_Name_f = uReadString(mod, block_Name_f_str, block_Name_f_len);
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
namespaceid:=fn(mod,block);uWriteString(mod,alloc,namespaceid,namespaceid_str,namespaceid_len);
},
}
}
func RegGetBlock(fn func(m api.Module,namespaceid string)(block BlockType,err OptError))(w HostBinding){
return HostBinding{"reg-get-block",func(ctx context.Context,mod api.Module,namespaceid_str uint32,namespaceid_len uint32,err_ptr uint32,block_Name_f_str uint32,block_Name_f_len uint32,block_IsSolid_f_b uint32,block_Color_f_R_f_w uint32,block_Color_f_G_f_w uint32,block_Color_f_B_f_w uint32,block_Color_f_A_f_w uint32)(){
_malloc:=mod.ExportedFunction("malloc");alloc:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));if err!=nil{panic(err)};return uint32(res[0])};var namespaceid = uReadString(mod, namespaceid_str, namespaceid_len);
block,err:=fn(mod,namespaceid);var block_Name_f = (string)(block.Name);
uWriteString(mod,alloc,block_Name_f,block_Name_f_str,block_Name_f_len);
var block_Color_f = (Color)(block.Color);
var block_Color_f_R_f = (uint8)(block_Color_f.R);
uWriteAnyFixed(mod,block_Color_f_R_f,block_Color_f_R_f_w);
var block_Color_f_G_f = (uint8)(block_Color_f.G);
uWriteAnyFixed(mod,block_Color_f_G_f,block_Color_f_G_f_w);
var block_Color_f_B_f = (uint8)(block_Color_f.B);
uWriteAnyFixed(mod,block_Color_f_B_f,block_Color_f_B_f_w);
var block_Color_f_A_f = (uint8)(block_Color_f.A);
uWriteAnyFixed(mod,block_Color_f_A_f,block_Color_f_A_f_w);
var block_IsSolid_f = (bool)(block.IsSolid);
if block_IsSolid_f{uWriteAnyFixed[uint32](mod,1,block_IsSolid_f_b)}else{uWriteAnyFixed[uint32](mod,0,block_IsSolid_f_b)};
uWritePointer(mod,alloc,err,err_ptr);
},
}
}
func AppendLang(fn func(m api.Module,name string)(err OptError))(w HostBinding){
return HostBinding{"append-lang",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,err_ptr uint32)(){
_malloc:=mod.ExportedFunction("malloc");alloc:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));if err!=nil{panic(err)};return uint32(res[0])};var name = uReadString(mod, name_str, name_len);
err:=fn(mod,name);uWritePointer(mod,alloc,err,err_ptr);
},
}
}
func ConfigSetInt(fn func(m api.Module,name string,value int64)())(w HostBinding){
return HostBinding{"config-set-int",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w int64)(){
var name = uReadString(mod, name_str, name_len);
var value = (int64)(value_w);
fn(mod,name,value);},
}
}
func ConfigSetFloat(fn func(m api.Module,name string,value float64)())(w HostBinding){
return HostBinding{"config-set-float",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w float64)(){
var name = uReadString(mod, name_str, name_len);
var value = (float64)(value_w);
fn(mod,name,value);},
}
}
func ConfigSetString(fn func(m api.Module,name string,value string)())(w HostBinding){
return HostBinding{"config-set-string",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_str uint32,value_len uint32)(){
var name = uReadString(mod, name_str, name_len);
var value = uReadString(mod, value_str, value_len);
fn(mod,name,value);},
}
}
func ConfigGetInt(fn func(m api.Module,name string)(value int64))(w HostBinding){
return HostBinding{"config-get-int",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w uint32)(){
var name = uReadString(mod, name_str, name_len);
value:=fn(mod,name);uWriteAnyFixed(mod,value,value_w);
},
}
}
func ConfigGetFloat(fn func(m api.Module,name string)(value float64))(w HostBinding){
return HostBinding{"config-get-float",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w uint32)(){
var name = uReadString(mod, name_str, name_len);
value:=fn(mod,name);uWriteAnyFixed(mod,value,value_w);
},
}
}
func ConfigGetString(fn func(m api.Module,name string)(value string))(w HostBinding){
return HostBinding{"config-get-string",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_str uint32,value_len uint32)(){
_malloc:=mod.ExportedFunction("malloc");alloc:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));if err!=nil{panic(err)};return uint32(res[0])};var name = uReadString(mod, name_str, name_len);
value:=fn(mod,name);uWriteString(mod,alloc,value,value_str,value_len);
},
}
}
