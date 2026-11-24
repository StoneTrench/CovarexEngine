# Bindings generator

Reads the `.abi` abi declaration file, and uses the embedded templates to generate bindings for the target language.

Function cannot be sent over the wire. Generics evaluate into their raw body, so its best to wrap them in a named type declaration.

## Usage

```sh
gen -l go -s ./env.abi -o ./out/
```

## Behaviour

Steps:
1. Parse the input `.abi` file
1. Collect function and type declarations, and evaluate generics
1. Build IR (intermediate representation)
1. Apply the chosen templates to the IR
1. Write the generated bindings

## Example

Here's the abi syntax:

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

And here's an example. You can define multiple modules and they will be emitted as multiple files:

```rs
module env {
	BlockKind enum u8 {
		Solid
		Air
		Transparent
		Liquid
	}

	color<T> struct { R T G T B T A T }
	ColorU8 color<u8>
	BlockType struct {
		Name string
		Color ColorU8
		Kind BlockKind
	}

	log-info fn(msg string)
	log-warn fn(msg string)
	log-error fn(msg string)

	reg-add-block fn(block BlockType)(namespaceid string)
	reg-get-block fn(namespaceid string)(block BlockType, err *string)
	append-lang fn(name string)(err *string)

	config-set-int fn(name string, value i64)
	config-set-float fn(name string, value f64)
	config-set-string fn(name string, value string)

	config-get-int fn(name string)(value i64)
	config-get-float fn(name string)(value f64)
	config-get-string fn(name string)(value string)
}
```