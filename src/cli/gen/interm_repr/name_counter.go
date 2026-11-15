package interm_repr

import (
	"fmt"

	"github.com/StoneTrench/go-mc-clone/src/cli/gen/types"
)

type NameCount uint

func (n *NameCount) CurrString() string {
	return fmt.Sprintf("r%d", *n)
}
func (n *NameCount) NextString() string {
	(*n) = (*n) + 1
	return n.CurrString()
}
func (n *NameCount) Next(t types.TypeRef) Name {
	return Name{Name: n.NextString(), Type: t}
}
