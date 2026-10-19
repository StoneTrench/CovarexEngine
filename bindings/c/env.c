#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

typedef struct { const char* ptr; uint32_t len; } wasm_string;
typedef struct { void* ptr; uint32_t len; } wasm_slice;

typedef struct Color {
	uint8_t R;
	uint8_t G;
	uint8_t B;
	uint8_t A;
} Color;
typedef struct BlockType {
	wasm_string Name;
	Color Color;
	bool IsSolid;
} BlockType;
__attribute__((import_module("env"), import_name("log-info")))
extern void env_log_info(uint32_t msg_str,uint32_t msg_len);
void log_info(wasm_string msg){
uint32_t msg_str;uint32_t msg_len;msg_str=(uint32_t)((uint32_t)(uintptr_t)msg.ptr);msg_len=(uint32_t)(msg.len);env_log_info(msg_str,msg_len);}
__attribute__((import_module("env"), import_name("log-warn")))
extern void env_log_warn(uint32_t msg_str,uint32_t msg_len);
void log_warn(wasm_string msg){
uint32_t msg_str;uint32_t msg_len;msg_str=(uint32_t)((uint32_t)(uintptr_t)msg.ptr);msg_len=(uint32_t)(msg.len);env_log_warn(msg_str,msg_len);}
__attribute__((import_module("env"), import_name("log-error")))
extern void env_log_error(uint32_t msg_str,uint32_t msg_len);
void log_error(wasm_string msg){
uint32_t msg_str;uint32_t msg_len;msg_str=(uint32_t)((uint32_t)(uintptr_t)msg.ptr);msg_len=(uint32_t)(msg.len);env_log_error(msg_str,msg_len);}
__attribute__((import_module("env"), import_name("reg-add-block")))
extern void env_reg_add_block(uint32_t block_Name_f_str,uint32_t block_Name_f_len,uint32_t block_IsSolid_f_b,uint32_t block_Color_f_R_f_w,uint32_t block_Color_f_G_f_w,uint32_t block_Color_f_B_f_w,uint32_t block_Color_f_A_f_w,uint32_t namespaceid_str,uint32_t namespaceid_len);
wasm_string reg_add_block(BlockType block){
wasm_string block_Name_f;block_Name_f=(wasm_string)((uint32_t)(uintptr_t)&block);uint32_t namespaceid_str;uint32_t namespaceid_len;env_reg_add_block(block_Name_f_str,block_Name_f_len,block_IsSolid_f_b,block_Color_f_R_f_w,block_Color_f_G_f_w,block_Color_f_B_f_w,block_Color_f_A_f_w,(uint32_t)(uintptr_t)&namespaceid_str,(uint32_t)(uintptr_t)&namespaceid_len);wasm_string namespaceid;namespaceid=(wasm_string)((wasm_string){ (const char*)(uintptr_t)namespaceid_str, namespaceid_len });wasm_reset_malloc();return namespaceid;}
__attribute__((import_module("env"), import_name("reg-get-block")))
extern void env_reg_get_block(uint32_t namespaceid_str,uint32_t namespaceid_len,uint32_t block_Name_f_str,uint32_t block_Name_f_len,uint32_t block_IsSolid_f_b,uint32_t block_Color_f_R_f_w,uint32_t block_Color_f_G_f_w,uint32_t block_Color_f_B_f_w,uint32_t block_Color_f_A_f_w);
BlockType reg_get_block(wasm_string namespaceid){
uint32_t namespaceid_str;uint32_t namespaceid_len;namespaceid_str=(uint32_t)((uint32_t)(uintptr_t)namespaceid.ptr);namespaceid_len=(uint32_t)(namespaceid.len);uint32_t block_Name_f_str;uint32_t block_Name_f_len;uint32_t block_IsSolid_f_b;uint32_t block_Color_f_R_f_w;uint32_t block_Color_f_G_f_w;uint32_t block_Color_f_B_f_w;uint32_t block_Color_f_A_f_w;env_reg_get_block(namespaceid_str,namespaceid_len,(uint32_t)(uintptr_t)&block_Name_f_str,(uint32_t)(uintptr_t)&block_Name_f_len,(uint32_t)(uintptr_t)&block_IsSolid_f_b,(uint32_t)(uintptr_t)&block_Color_f_R_f_w,(uint32_t)(uintptr_t)&block_Color_f_G_f_w,(uint32_t)(uintptr_t)&block_Color_f_B_f_w,(uint32_t)(uintptr_t)&block_Color_f_A_f_w);BlockType block;void* block;block=*(void**)(uintptr_t)block_Name_f;wasm_reset_malloc();return block;}
__attribute__((import_module("env"), import_name("append-lang")))
extern void env_append_lang(uint32_t name_str,uint32_t name_len,uint32_t err_ptr);
wasm_string* append_lang(wasm_string name){
uint32_t name_str;uint32_t name_len;name_str=(uint32_t)((uint32_t)(uintptr_t)name.ptr);name_len=(uint32_t)(name.len);uint32_t err_ptr;env_append_lang(name_str,name_len,(uint32_t)(uintptr_t)&err_ptr);wasm_string* err;err=(wasm_string*)(uintptr_t)err_ptr;wasm_reset_malloc();return err;}
__attribute__((import_module("env"), import_name("config-set-int")))
extern void env_config_set_int(uint32_t name_str,uint32_t name_len,int64_t value_w);
void config_set_int(wasm_string name,int64_t value){
uint32_t name_str;uint32_t name_len;name_str=(uint32_t)((uint32_t)(uintptr_t)name.ptr);name_len=(uint32_t)(name.len);int64_t value_w;value_w=(int64_t)(value);env_config_set_int(name_str,name_len,value_w);}
__attribute__((import_module("env"), import_name("config-set-float")))
extern void env_config_set_float(uint32_t name_str,uint32_t name_len,double value_w);
void config_set_float(wasm_string name,double value){
uint32_t name_str;uint32_t name_len;name_str=(uint32_t)((uint32_t)(uintptr_t)name.ptr);name_len=(uint32_t)(name.len);double value_w;value_w=(double)(value);env_config_set_float(name_str,name_len,value_w);}
__attribute__((import_module("env"), import_name("config-set-string")))
extern void env_config_set_string(uint32_t name_str,uint32_t name_len,uint32_t value_str,uint32_t value_len);
void config_set_string(wasm_string name,wasm_string value){
uint32_t name_str;uint32_t name_len;name_str=(uint32_t)((uint32_t)(uintptr_t)name.ptr);name_len=(uint32_t)(name.len);uint32_t value_str;uint32_t value_len;value_str=(uint32_t)((uint32_t)(uintptr_t)value.ptr);value_len=(uint32_t)(value.len);env_config_set_string(name_str,name_len,value_str,value_len);}
__attribute__((import_module("env"), import_name("config-get-int")))
extern void env_config_get_int(uint32_t name_str,uint32_t name_len,uint32_t value_w);
int64_t config_get_int(wasm_string name){
uint32_t name_str;uint32_t name_len;name_str=(uint32_t)((uint32_t)(uintptr_t)name.ptr);name_len=(uint32_t)(name.len);uint32_t value_w;env_config_get_int(name_str,name_len,(uint32_t)(uintptr_t)&value_w);int64_t value;value=(int64_t)(value_w);return value;}
__attribute__((import_module("env"), import_name("config-get-float")))
extern void env_config_get_float(uint32_t name_str,uint32_t name_len,uint32_t value_w);
double config_get_float(wasm_string name){
uint32_t name_str;uint32_t name_len;name_str=(uint32_t)((uint32_t)(uintptr_t)name.ptr);name_len=(uint32_t)(name.len);uint32_t value_w;env_config_get_float(name_str,name_len,(uint32_t)(uintptr_t)&value_w);double value;value=(double)(value_w);return value;}
__attribute__((import_module("env"), import_name("config-get-string")))
extern void env_config_get_string(uint32_t name_str,uint32_t name_len,uint32_t value_str,uint32_t value_len);
wasm_string config_get_string(wasm_string name){
uint32_t name_str;uint32_t name_len;name_str=(uint32_t)((uint32_t)(uintptr_t)name.ptr);name_len=(uint32_t)(name.len);uint32_t value_str;uint32_t value_len;env_config_get_string(name_str,name_len,(uint32_t)(uintptr_t)&value_str,(uint32_t)(uintptr_t)&value_len);wasm_string value;value=(wasm_string)((wasm_string){ (const char*)(uintptr_t)value_str, value_len });wasm_reset_malloc();return value;}
