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

// Body returns the IType of the r TypeRef.
func (c *TypeCollection) Body(r TypeRef) (IType, error) {
	if r == TypeRef_Invalid {
		return nil, fmt.Errorf("failed to get body: invalid typeref")
	}
	id := int(r - 1)
	if len(c.Types) <= id {
		return nil, fmt.Errorf("failed to get body: index [%d] was greater than size of types list [%d]", id, len(c.Types))
	}
	return c.Types[id], nil
}

func (c *TypeCollection) RebuildUniqueTable() {
	c.Unique = make(map[string]TypeRef)
	for id, t := range c.Types {
		key := String(t)
		c.Unique[key] = TypeRef(id + 1)
	}
}

func (c *TypeCollection) Reset() {
	c.Types = make([]IType, 0)
	c.Unique = make(map[string]TypeRef)
}

func (c *TypeCollection) Update(r TypeRef, t IType) error {
	old_t, err := c.Body(r)
	if err != nil {
		return err
	}
	old_key := String(old_t)
	new_key := String(t)
	delete(c.Unique, old_key)
	c.Unique[new_key] = r
	return nil
}
