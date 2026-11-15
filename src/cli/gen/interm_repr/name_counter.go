package intermediate_representation

import "fmt"

type NameCount uint

func (n *NameCount) CurrString() string {
	return fmt.Sprintf("r%d", *n)
}
func (n *NameCount) NextString() string {
	(*n) = (*n) + 1
	return n.CurrString()
}
func (n *NameCount) Next(t TYPE_REF) Name {
	return Name{Name: n.NextString(), Type: t}
}
