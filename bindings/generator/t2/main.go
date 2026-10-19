package t2

import (
	"fmt"
	"strings"

	"github.com/StoneTrench/go-mc-clone/bindings/generator/types"
)

type Builder strings.Builder

func (w *Builder) Write(s string) {
	(*strings.Builder)(w).WriteString(s)
}
func (w *Builder) Writef(format string, args ...any) {
	fmt.Fprintf((*strings.Builder)(w), format, args...)
}
func (w *Builder) String() string {
	return (*strings.Builder)(w).String()
}


func GuestABIWrite(w *Builder, name_t string, type_t types.IType) int {
	switch e := type_t.(type) {
	case types.TypePrimitive:
		
		return 1
	}

	return 0
}
func GuestABIRead(w *Builder, name_t string, type_t types.IType) int {

}
func HostABIWrite(w *Builder, name_t string, type_t types.IType) int {

}
func HostABIRead(w *Builder, name_t string, type_t types.IType) int {

}
