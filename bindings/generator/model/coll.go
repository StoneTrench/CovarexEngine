package types

import "fmt"

type TypeRef uint32

const (
	TypeRef_Invalid TypeRef = iota
)

type TypeCollection struct {
	Types  []IType
	Unique map[string]TypeRef
}

func (c *TypeCollection) Init() *TypeCollection {
	c.Types = make([]IType, 0)
	c.Unique = make(map[string]TypeRef)
	return c
}
func (c *TypeCollection) Ref(t IType) TypeRef {
	key := String(t)
	if id, ok := c.Unique[key]; ok {
		return id
	}
	id := TypeRef(len(c.Types)) + 1
	c.Types = append(c.Types, t)
	c.Unique[key] = id
	return id
}
func (c *TypeCollection) Body(r TypeRef) (IType, error) {
	if r == TypeRef_Invalid {
		return nil, fmt.Errorf("invalid typeref")
	}
	return c.Types[r-1], nil
}
