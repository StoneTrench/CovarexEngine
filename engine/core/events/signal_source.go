package events

import (
	"errors"
	"sync"
)

type SubId uint16
type Subscriber[T any] func(e T) error

type listener_t[T any] struct {
	id SubId
	fn Subscriber[T]
}

type SignalSource[T any] struct {
	mu       sync.RWMutex
	next_id  SubId
	handlers []listener_t[T]
}

func NewSource[T any]() *SignalSource[T] {
	return &SignalSource[T]{
		mu:       sync.RWMutex{},
		next_id:  0,
		handlers: make([]listener_t[T], 0),
	}
}

func (e *SignalSource[T]) handle_event(event T) error {
	err_list := []error{}
	for i := range e.handlers {
		err := e.handlers[i].fn(event)
		if err != nil {
			err_list = append(err_list, err)
		}
	}
	return errors.Join(err_list...)
}

func (e *SignalSource[T]) EmitImmediate(event T) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.handle_event(event)
}

func (e *SignalSource[T]) EmitDeferred(event T) {
	go e.EmitImmediate(event)
}

func (e *SignalSource[T]) Subscribe(fn Subscriber[T]) SubId {
	e.mu.Lock()
	defer e.mu.Unlock()
	id := e.next_id
	e.next_id++
	e.handlers = append(e.handlers, listener_t[T]{
		id: id,
		fn: fn,
	})
	return id
}

func (e *SignalSource[T]) Unsubscribe(id SubId) {
	e.mu.Lock()
	defer e.mu.Unlock()
	last := len(e.handlers) - 1
	for i := range e.handlers {
		if e.handlers[i].id == id {
			e.handlers[i] = e.handlers[last]
			e.handlers = e.handlers[:last]
			return
		}
	}
}
