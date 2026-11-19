package platform

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

// VFS is a virtual filesystem abstraction to allow for wasm build targets.
type VFS struct {
	read_only []fs.FS
	can_write bool // Has access to os.Create and os.Open, so it can write and read
}

const __ERR_NO_ACC = "has no OS filesystem access"

func (o *VFS) Init(embeds []fs.FS, can_write bool) *VFS {
	o.read_only = embeds
	o.can_write = can_write
	return o
}
func (o *VFS) AppendFilesystem(f fs.FS) {
	o.read_only = append(o.read_only, f)
}
func (o *VFS) HasWriteAccess() bool {
	return o.can_write
}

func (o *VFS) OpenFS(name string) (fs.File, error) {
	for _, f := range o.read_only {
		res, err := f.Open(name)
		if err == nil {
			return res, err
		}
	}
	if !o.can_write {
		return nil, &fs.PathError{Op: "open", Path: name, Err: errors.New(__ERR_NO_ACC)}
	}

	return os.Open(name)
}
func (o *VFS) ReadDir(name string) ([]fs.DirEntry, error) {
	for _, f := range o.read_only {
		res, err := fs.ReadDir(f, name)
		if err == nil {
			return res, err
		}
	}
	if !o.can_write {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: errors.New(__ERR_NO_ACC)}
	}

	return os.ReadDir(name)
}
func (o *VFS) Open(name string) (io.ReadCloser, error) {
	return o.OpenFS(name)
}
func (o *VFS) Create(name string) (io.ReadWriteCloser, error) {
	if !o.can_write {
		return nil, &fs.PathError{Op: "create", Path: name, Err: errors.New(__ERR_NO_ACC)}
	}
	return os.Create(name)
}

func (o *VFS) OpenOrCreate(name string) (io.ReadCloser, error) {
	file, err := o.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		err = file.Close()
		if err != nil {
			return nil, err
		}
		file, err = o.Create(name)
	}
	return file, err
}
func (o *VFS) Exists(name string) bool {
	_, err := o.Open(name)
	return !errors.Is(err, os.ErrNotExist)
}
