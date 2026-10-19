package host

import (
	"fmt"
	"unsafe"

	"github.com/tetratelabs/wazero/api"
)

type usize = uint32

const Sizeof_usize = 32
const MAX_usize = (1 << Sizeof_usize) - 1

type HostBinding struct {
	Name string
	Fn   interface{}
}
type FnAllocator func(size usize) (offset usize)

func uPanic(m api.Module, err error) {
	panic(fmt.Errorf("error occured in (%s): %w", m.Name(), err))
}
func uPanicf(m api.Module, format string, args ...any) {
	uPanic(m, fmt.Errorf(format, args...))
}
func uPanics(m api.Module, msg string) {
	uPanicf(m, "%s", msg)
}

func uReadString(m api.Module, offset, length usize) string {
	if length == 0 {
		return ""
	}
	raw, ok := m.Memory().Read(offset, length)
	if !ok {
		uPanics(m, "out of bounds")
	}
	return string(raw)
}
func uReadPlex[T any](m api.Module, offset usize) T {
	var v T
	size := usize(unsafe.Sizeof(v))
	if size == 0 {
		return v
	}
	raw, ok := m.Memory().Read(offset, size)
	if !ok {
		uPanics(m, "out of bounds")
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&v)), size), raw)
	return v
}
func uReadSlice[T any](m api.Module, offset, length usize) []T {
	if length == 0 {
		return nil
	}
	var zero T
	elemSize := usize(unsafe.Sizeof(zero))
	size := length * elemSize
	raw, ok := m.Memory().Read(offset, size)
	if !ok {
		uPanics(m, "out of bounds")
	}
	out := make([]T, length)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&out[0])), size), raw)
	return out
}
func uReadPointer[T any](m api.Module, offset usize) *T {
	if offset == 0 {
		return nil
	}
	v := new(T)
	size := usize(unsafe.Sizeof(*v))
	raw, ok := m.Memory().Read(offset, size)
	if !ok {
		uPanics(m, "out of bounds")
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(v)), size), raw)
	return v
}

func uWriteString(m api.Module, alloc FnAllocator, str string, ptr_ptr, ptr_len usize) {
	if len(str) > MAX_usize {
		uPanics(m, "string too large")
	}
	length := usize(len(str))
	if length == 0 {
		uWriteAnyFixed[usize](m, 0, ptr_ptr)
		uWriteAnyFixed[usize](m, 0, ptr_len)
		return
	}
	offset := alloc(length)
	if offset == 0 {
		uPanics(m, "alloc returned null pointer")
	}
	ok := m.Memory().Write(offset, ([]byte)(str))
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[usize](m, offset, ptr_ptr)
	uWriteAnyFixed[usize](m, length, ptr_len)
}
func uWritePlex[T any](m api.Module, alloc FnAllocator, v T, ptr_ptr usize) {
	size := usize(unsafe.Sizeof(v))
	if size == 0 {
		uWriteAnyFixed[usize](m, 0, ptr_ptr)
		return
	}
	offset := alloc(size)
	if offset == 0 {
		uPanics(m, "alloc returned null pointer")
	}
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&v)), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[usize](m, offset, ptr_ptr)
}
func uWriteSlice[T any](m api.Module, alloc FnAllocator, v []T, ptr_ptr, ptr_len usize) {
	if len(v) == 0 {
		uWriteAnyFixed[usize](m, 0, ptr_ptr)
		uWriteAnyFixed[usize](m, 0, ptr_len)
		return
	}
	var zero T
	elemSize := usize(unsafe.Sizeof(zero))
	length := usize(len(v))
	size := length * elemSize
	offset := alloc(size)
	if offset == 0 {
		uPanics(m, "alloc returned null pointer")
	}
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&v[0])), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[usize](m, offset, ptr_ptr)
	uWriteAnyFixed[usize](m, length, ptr_len)
}
func uWritePointer[T any](m api.Module, alloc FnAllocator, v *T, ptr_ptr usize) {
	if v == nil {
		uWriteAnyFixed[usize](m, 0, ptr_ptr)
		return
	}
	size := usize(unsafe.Sizeof(*v))
	offset := alloc(size)
	if offset == 0 {
		uPanics(m, "alloc returned null pointer")
	}
	raw := unsafe.Slice((*byte)(unsafe.Pointer(v)), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[usize](m, offset, ptr_ptr)
}
func uWriteAnyFixed[T any](m api.Module, v T, offset usize) {
	size := usize(unsafe.Sizeof(v))
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&v)), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
}
