# Design overview

1. No pointers are exposed to `guest` or `host`.
1. No pointers are exposed on the interface's surface, they are only internal.
1. No unions exposed to guest or host.

## Types

1. Primitives
	+ Sent directly.
1. Structs/Plexes (Product types)
	+ Either flattened, or sent via pointer copy.
1. Unions (Sum types)
	+ Cast to a byte array, and then back. Casting is recursive.

+ Every type is known at compile-time. Their sizes and layout are fixed.


### Syntax
<!-- ENUM = 'enum', PRIM_INT, '{', {Name, ['(', [FIELD, {',', FIELD}], ')']}, '}'; -->

```EBNF
FIELD = Name, TYPE;

PRIM_INT = 'u8'|'u16'|'u32'|'u64'|'i8'|'i16'|'i32'|'i64'|'usize';
PRIMITIVE = PRIM_INT|'f32'|'f64'|'bool'|'string';
STRUCT = 'struct', '{', {FIELD}, '}';
UNION = 'union', '{', {FIELD}, '}';
ENUM = 'enum', PRIM_INT, '{', {Name}, '}';
ARRAY = '[', Integer, ']', TYPE;
SLICE = '[', ']', TYPE;
POINTER = '*', ['mut'], TYPE;
FUNCTION = 'fn', '(', [FIELD, {',', FIELD}], ')', ['(', [FIELD], ')'];
IDENT = Name, ['<', TYPE, {',', TYPE}, '>'];

TYPE = PRIMITIVE | STRUCT | UNION | ENUM | ARRAY | SLICE | POINTER | FUNCTION | IDENT | '(', TYPE, ')';

DECL = Name, ['<' Name, {',', Name}, '>'], TYPE;
MODULE = 'module', Name, '{', { DECL }, '}';
SOURCE = { MODULE }, Eof;
```

## It's purpose

The purpose of the generator is to convert the abi interface declarations into two files per module. One file contains the bindings functions for the `host`,
```go
func MyInterface(handler func(mod api.Module, name string) ()) (string, interface{}) { ... }
```
these generate the binding that gets registered in the `host`. The other file contains the `guest` side implementation, that can be imported for `guest` module development.
```go
// go:wasmimport env my-interface
func wasm_my_interface(t0 i32, t1 i32)
func MyInterface(name string) () { ... }
```

Both the `guest` and `host` get extra utilities, that are language dependant, but still emmitted by the generator, such as:
+ Union type for golang, with access helpers.
+ Exported `wasm_alloc` function for the `guest`. Since at the boundary everything is borrowed/copied, there's no need to export `wasm_free`.

## Guest (Calls the host)
Receives domain specific data, returns domain specific data.
1. Preprocess data (like how go and c handle strings so differently, first have to handle that).
1. Serialize data. Flatten, or pass pointers.
1. Call wasm function with args.

1. Deserialize data.
1. Postprocess data.
1. Clear out allocations by the `host` and `guest`.
1. Return data (this is never a pointer).

Or: `domain -> wire -> call -> read -> free -> return`

## Host (Handles calls from guest)
Receives wasm compatible data, returns wasm compatible data.
1. Copy underlying data from pointers from `guest`.
1. Deserialize data.
1. Postprocess data.
1. Call handler.

1. Allocate space inside `guest` for pointers.
1. Preprocess data (like how go and c handle strings so differently, first have to handle that).
1. Serialize data.
1. Return data (this can be a pointer).

Or: `wire -> copy -> decode -> call -> alloc -> write -> return`

## Operations & Instructions

The simplest way to generate this kind of system from an interface, especially since data only moves in one direction, is SSA (Static Single Asignment). This way behaviours can be split into individual instructions, and each language only has to implement the templates for these, without having to worry about generating expressions. However the more specific an instruction, the better, because each language has it's own way of handing types.

To convert arrays, loops are also necessary.

Guest calls an api function binding. Input args of the guest are converted from the types of the guest domain, to types compatible with the interface (null terminated string -> slice string, etc), then these get lowered/encoded/serialized. Then the host function is called. The results are decoded/deserialized/raised back to interface types, after which they are converted to compatible types for the domain, and interface values are freed if needed. Domain values are then returned.

### Instructions

1. `Primitive` -> `Primitive`
	1. `string` -> `Pointer, Length`
1. `Enum` -> `typeof(Enum.Tag)` (WIP)
1. `Struct` -> `Pointer(Struct)`
1. `Union` -> `Pointer(Union)` (WIP)
1. `Array` -> `Pointer(Array)`
1. `Slice` -> `Pointer, Length`
1. `Pointer` -> `Primitive`
1. `Function` -> `Invalid`
1. `Identifier` -> `Body`

#### Casting

1. `string`
```go
// Guest write
var adr = int32(unsafe.Pointer(&arg[0]));
var sln = int32(len(arg));
// Host read
var raw = m.Read(adr, sln)
var arg = string(raw)

// Host write
var sln = int32(len(arg))
var adr = m.Alloc(sln)
m.Write(adr, []byte(arg))
m.Write(out_adr, adr)
m.Write(out_sln, sln)
// Guest read
var out_adr int32 // pointer
var out_sln int32 // length
_api_func(
	int32(uintptr(unsafe.Pointer(&out_adr))),
	int32(uintptr(unsafe.Pointer(&out_sln))),
)
var arg = unsafe.String(
	(*byte)(unsafe.Pointer(uintptr(out_adr))),
	int(out_sln)
)
```
1. `Struct` `Array` `Union` 
```go
// Guest write
var arg MyStruct
var res wasm_MyStruct
res.
```
1. `slice`
```go
var adr = uintptr(unsafe.Pointer(&arg[0]));
var len = uintptr(len(arg));
```
1. Iterate
```go
for i := range arg {
	e := arg[i]
}
```

## The solution

Guest has two type declarations for everything. In one everything is in a local type, that the guest module can work with. In the other every type has an explicit fixed format, that is shared between guest and host.


# FUCK

- I need an ABI.
- Convert language specific types into an agnostic universal format.
- Golang templates create the actual abi.
- The model encodes what the abi contains.
	- Does it also contain instructions for conversion?
	- What do these instructions look like?
	- It has to generate an SSA (Single Static Assignment) IR (Intermediate Representation).
	- The template takes the IR and emits language specific code.
- Is it a good idea to have the instructions, types, type conversions data driven? 

## TODO

- [ ] Standardize the ABI format.

- [ ] Design the IR.
	- [ ] What is a function?
	- [ ] Do functions have blocks?
	- [ ] What are instructions?
	- [ ] What instructions are there?

## The Design

I need to make a custom ABI interface generator for WASM. It's for a gameengine so developers can declare api functions for their systems, on one side is the host which is the game written in golang, on the other is the wasm guest. Both the guest and host get dencoding/decoding functions. Data gets copied from and to the guest, so the guest has to free/allocate its own memory, and the host can only allocate in the guest not free.


Iterate the input params of the function. Recursively unpack and convert every type into the universal format, while emitting conversion instructions for the templating system, and defining the layout of the transmission. For this we need a datastructure for every type, to encode the decoded types, and instructions to encode and decode it.

The IR generator then looks up the types in the datastructure and emits the correct instructions and type declarations for the chosen language.