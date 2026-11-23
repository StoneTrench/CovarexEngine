package platform

import (
	"bytes"
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

// OpenFile_DEF_PERM is the default perm value, for Create.
const OpenFile_DEF_PERM fs.FileMode = 0666

// MkdirAll_DEF_PERM is the default perm value, for MkdirAll.
const MkdirAll_DEF_PERM fs.FileMode = 0755

var ErrNoFileSystem = errors.New("has no OS filesystem access")

// Init initializes a virtual filesystem.
//
// It's a wrapper around the FS hierarchical file system.
// It has an array of read-only file systems, that it iterates 0 to n.
// And a flag that should depend on the target system.
//
// The purpose of this is to easily embed files into the game, so every system that needs file access goes through this.
// Also for WASM targets, especially on a webpage, you don't have os disk access to the disk. Usually, if you can't write but still try to, the function will return an error.
func (vfs *VFS) Init(embeds []fs.FS, can_write bool) *VFS {
	vfs.read_only = embeds
	vfs.can_write = can_write
	return vfs
}

// AppendFilesystem adds an FS to the end of the read-only array.
func (vfs *VFS) AppendFilesystem(f fs.FS) {
	vfs.read_only = append(vfs.read_only, f)
}

// HasWriteAccess returns true if the project has access to os.Create and os.Open, so it can write and read.
//
// This should probably be called HasDiskAccess.
func (vfs *VFS) HasWriteAccess() bool {
	return vfs.can_write
}

// Open opens the named file. It's the same as fs.Open or os.Open.
//
// If it has no disk access, it returns a PathError.
func (vfs *VFS) Open(name string) (io.ReadCloser, error) {
	for _, f := range vfs.read_only {
		res, err := f.Open(name)
		if err == nil {
			return res, err
		}
	}
	if !vfs.can_write {
		return nil, &fs.PathError{Op: "open", Path: name, Err: ErrNoFileSystem}
	}

	return os.Open(name)
}

// ReadDir first checks the read-only filesystems then it tries to ReadDir on the actual disk, if that fails, it returns a PathError.
func (vfs *VFS) ReadDir(name string) ([]fs.DirEntry, error) {
	for _, f := range vfs.read_only {
		res, err := fs.ReadDir(f, name)
		if err == nil {
			return res, err
		}
	}
	if !vfs.can_write {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: ErrNoFileSystem}
	}

	return os.ReadDir(name)
}

// OpenFile tries to call os.OpenFile, and fails if the VFS has no file system access, at which point it returns a PathError.
func (vfs *VFS) OpenFile(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
	if !vfs.can_write {
		return nil, &fs.PathError{Op: "openfile", Path: name, Err: ErrNoFileSystem}
	}
	return os.OpenFile(name, flag, perm)
}

// MkdirAll tries to call os.Mkdirall, and fails if the VFS has no file system access, at which point it returns a PathError.
func (vfs *VFS) MkdirAll(name string, perm fs.FileMode) error {
	if !vfs.can_write {
		return &fs.PathError{Op: "mkdirall", Path: name, Err: ErrNoFileSystem}
	}
	return os.MkdirAll(name, perm)
}

// Remove tries to remove the named file, and fails if the VFS has no file system access, at which point it returns a PathError.
func (vfs *VFS) Remove(name string) error {
	if !vfs.can_write {
		return &fs.PathError{Op: "remove", Path: name, Err: ErrNoFileSystem}
	}
	return os.Remove(name)
}

// Stat returns the FileInfo for the named file, or a PathError if it has no system access.
func (vfs *VFS) Stat(name string) (fs.FileInfo, error) {
	for _, f := range vfs.read_only {
		res, err := fs.Stat(f, name)
		if err == nil {
			return res, err
		}
	}
	if !vfs.can_write {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: ErrNoFileSystem}
	}
	return os.Stat(name)
}

// Create just calls vfs.OpenFile with the appropriate flags and perms for read and write.
func (vfs *VFS) Create(name string) (io.ReadWriteCloser, error) {
	return vfs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, OpenFile_DEF_PERM)
}

// Exists returns true if the named file exists. Uses vfs.Open.
func (vfs *VFS) Exists(name string) (bool, error) {
	file, err := vfs.Open(name)
	if file != nil {
		file.Close()
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// WriteFile uses vfs.OpenFile to open a named file for write, and io.Copy to write the bytes into it. Returns a PathError if there is no disk access.
func (vfs *VFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	file, err := vfs.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, werr := io.Copy(file, bytes.NewReader(data))
	cerr := file.Close()
	if werr == nil {
		werr = cerr
	}
	if werr != nil {
		return werr
	}
	return nil
}

// ReadFile uses vfs.Open to open a named file, and io.ReadAll to return the bytes in the file.
func (vfs *VFS) ReadFile(name string) ([]byte, error) {
	file, err := vfs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
