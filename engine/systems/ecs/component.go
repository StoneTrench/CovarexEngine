package ecs

import (
	"fmt"
	"sync"
)

type csize_t uint64

const (
	__ERR_INVALID_SIZE  = "invalid component size: expected %d got %d bytes"
	__ERR_INVALID_EID   = "invalid entity id (probably forgot to resize): index outside the bounds of the array (%d)[%d]"
	__ERR_CAPACITY      = "component set capacity exceeded"
	__ERR_INVALID_STATE = "invalid state: there are no components, and yet this entity %d claims to have a component from this set with component id %d"
)

// ComponentSet is a byte array, that stores the components as raw bytes, to allow for new component types at runtime.
type ComponentSet struct {
	mu sync.RWMutex

	// CSize is the size of a single component
	CSize csize_t

	// 0 is a sentinel value, meaning there is no component attached
	EntityToComponent []id_t
	// there is no sentinel value
	ComponentToEntity []id_t

	// Components is the raw byte array representing the components
	Components []byte
}

// TryResize resizes the set when a new entity is added.
func (s *ComponentSet) TryResize(e EntityRef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := e.id + 1
	if n > id_t(len(s.EntityToComponent)) {
		s.EntityToComponent = append(s.EntityToComponent, make([]id_t, n)...)
	}
}

func (s *ComponentSet) impl_GetComponentUnsafe(e EntityRef) ([]byte, error) {
	if int(e.id) >= len(s.EntityToComponent) {
		return nil, fmt.Errorf(__ERR_INVALID_EID, len(s.EntityToComponent), e.id)
	}
	cid := csize_t(s.EntityToComponent[e.id])
	if cid == 0 {
		return nil, nil
	}

	index := cid - 1
	cptr := index * s.CSize
	return s.Components[cptr:(cptr + s.CSize)], nil
}

// SetComponent adds a new component to an entity, or overwrites an existing one on the entity.
func (s *ComponentSet) SetComponent(e EntityRef, src []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(src) != int(s.CSize) {
		return fmt.Errorf(__ERR_INVALID_SIZE, s.CSize, len(src))
	}
	if int(e.id) >= len(s.EntityToComponent) {
		return fmt.Errorf(__ERR_INVALID_EID, len(s.EntityToComponent), e.id)
	}

	component_id := s.EntityToComponent[e.id]
	if component_id == 0 {
		count := len(s.Components) / int(s.CSize)
		if count >= int(^id_t(0)) {
			return fmt.Errorf(__ERR_CAPACITY)
		}

		s.Components = append(s.Components, src...)
		s.ComponentToEntity = append(s.ComponentToEntity, e.id)

		s.EntityToComponent[e.id] = id_t(count) + 1
	} else {
		// keeps the same component id
		// same mapping too
		// just replaces it with the new value

		index := uint32(component_id - 1)
		component_ptr := index * uint32(s.CSize)
		copy(s.Components[component_ptr:(component_ptr+uint32(s.CSize))], src)
	}
	return nil
}
// CopyComponent copies the raw bytes of a component to the destination buffer.
func (s *ComponentSet) CopyComponent(e EntityRef, dst []byte) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(dst) != int(s.CSize) {
		return false, fmt.Errorf(__ERR_INVALID_SIZE, s.CSize, len(dst))
	}

	component, err := s.impl_GetComponentUnsafe(e)
	if err != nil {
		return false, err
	}
	copy(dst, component)
	return true, nil
}
// CopyComponentInternal copies the raw bytes of a component from one entity to the other inside the ComponentSet.
func (s *ComponentSet) CopyComponentInternal(dst EntityRef, src EntityRef) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c_dst, err := s.impl_GetComponentUnsafe(dst)
	if err != nil {
		return false, err
	}
	c_src, err := s.impl_GetComponentUnsafe(src)
	if err != nil {
		return false, err
	}

	copy(c_dst, c_src)
	return true, nil
}
// RemoveComponent removes a component from the entity.
func (s *ComponentSet) RemoveComponent(e EntityRef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if int(e.id) >= len(s.EntityToComponent) {
		return // invalid id
	}

	cid_remove := csize_t(s.EntityToComponent[e.id]) // get the component id
	cid_last := csize_t(len(s.Components)) / s.CSize // get the last component's id

	if cid_remove == 0 {
		return // entity doesnt have a component
	}
	if cid_last == 0 {
		panic(fmt.Errorf(__ERR_INVALID_STATE, e.id, cid_remove))
		// return // there are no components, therefore entity shouldn't have had a component to begin with, which technically makes this an invalid state
	}

	c_index_remove := cid_remove - 1
	c_index_last := cid_last - 1

	cptr_remove := c_index_remove * s.CSize
	cptr_last := c_index_last * s.CSize

	s.EntityToComponent[e.id] = 0 // set the component id to the sentinel value

	// move the last component to the target component
	if cid_last != cid_remove {
		copy(s.Components[cptr_remove:(cptr_remove+s.CSize)], s.Components[cptr_last:(cptr_last+s.CSize)])
		c := s.ComponentToEntity[c_index_last]
		s.EntityToComponent[c] = id_t(cid_remove)
		s.ComponentToEntity[c_index_remove] = c
	}
	// shrink the slice by 1
	s.Components = s.Components[:cptr_last]
	s.ComponentToEntity = s.ComponentToEntity[:c_index_last]
}
