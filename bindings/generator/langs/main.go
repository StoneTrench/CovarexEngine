package langs

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	. "github.com/StoneTrench/go-mc-clone/bindings/generator/types"
)

type Instruction uint8
const (
	Instr_cast Instruction = iota
	Instr_unwrap // flatten but only once non-recursively
	Instr_pointer
)

/*
LangBuilder ...

If function type: Encode input args to ABI compatible types.
After which, decode output values from ABI, and return them.

If other type: Create type declaration.
*/
type LangBuilder struct {
	WriteHeader func(b *strings.Builder, moduleName string)
	WriteFooter func(b *strings.Builder, moduleName string)

	// If function, write function declaration with WASM import.
	// If type, write type declaration.
	WriteGuestDecl func(b *strings.Builder, moduleName string, p TypeParam)

	// WriteGuestEncode Data -> ABI.
	WriteGuestEncode func(b *strings.Builder, p TypeParam)
	// WriteGuestEncode Data <- ABI.
	WriteGuestDecode func(b *strings.Builder, p TypeParam)

	// Converts the type to a string.
	StringType func(p IType) string
}

// ToUpperFirst capitalizes the first letter.
func ToUpperFirst(s string) string {
	if s == "" {
		return s
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// HyphenCaseToPascalCase converts hyphen case `my-ident` to pascal case `MyIdent`.
func HyphenCaseToPascalCase(s string) string {
	res := strings.Builder{}

	for _, s := range strings.Split(s, "-") {
		res.WriteString(ToUpperFirst(strings.ToLower(s)))
	}

	return res.String()
}

// Writef is a self explanatory helper function.
func Writef(b *strings.Builder, format string, args ...any) {
	b.WriteString(fmt.Sprintf(format, args...))
}

// Writes is a self explanatory helper function.
func Writes(b *strings.Builder, msg string) {
	b.WriteString(msg)
}
