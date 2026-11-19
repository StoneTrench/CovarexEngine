package platform

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
)

// VFS is a virtual filesystem abstraction to allow for wasm build targets.
type VFS struct {
	read_only []fs.FS
	can_write bool // Has access to os.Create and os.Open, so it can write and read
}

const OpenFile_DEF_PERM fs.FileMode = 0666
const MkdirAll_DEF_PERM fs.FileMode = 0755
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

func (o *VFS) Open(name string) (fs.File, error) {
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
func (o *VFS) OpenFile(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
	if !o.can_write {
		return nil, &fs.PathError{Op: "openfile", Path: name, Err: errors.New(__ERR_NO_ACC)}
	}
	return os.OpenFile(name, flag, perm)
}
func (o *VFS) MkdirAll(name string, perm fs.FileMode) error {
	if !o.can_write {
		return &fs.PathError{Op: "mkdirall", Path: name, Err: errors.New(__ERR_NO_ACC)}
	}
	return os.MkdirAll(name, perm)
}
func (o *VFS) Remove(name string) error {
	if !o.can_write {
		return &fs.PathError{Op: "remove", Path: name, Err: errors.New(__ERR_NO_ACC)}
	}
	return os.Remove(name)
}
func (o *VFS) Stat(name string) (fs.FileInfo, error) {
	for _, f := range o.read_only {
		res, err := fs.Stat(f, name)
		if err == nil {
			return res, err
		}
	}
	if !o.can_write {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: errors.New(__ERR_NO_ACC)}
	}
	return os.Stat(name)
}

func (o *VFS) Create(name string) (io.ReadWriteCloser, error) {
	return o.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, OpenFile_DEF_PERM)
}
func (o *VFS) Exists(name string) (bool, error) {
	file, err := o.Open(name)
	if file != nil {
		file.Close()
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (o *VFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	f, err := o.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, werr := io.Copy(f, bytes.NewReader(data))
	cerr := f.Close()
	if werr == nil {
		werr = cerr
	}
	if werr != nil {
		return werr
	}
	return nil
}
func (o *VFS) ReadFile(name string) ([]byte, error) {
	f, err := o.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func (o *VFS) TreeFS() string {
	nodes := make(map[string]bool)
	for _, f := range o.read_only {
		_ = fs.WalkDir(f, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil || p == "." {
				return nil // swallow per-entry errors; this is a debug view
			}
			if _, seen := nodes[p]; !seen {
				nodes[p] = d.IsDir()
			}
			return nil
		})
	}
	if len(nodes) == 0 {
		return ".\n"
	}

	// Bucket children by parent so we can recurse without a second walk.
	children := make(map[string][]string, len(nodes))
	for p := range nodes {
		parent := path.Dir(p)
		children[parent] = append(children[parent], p)
	}
	for _, kids := range children {
		sort.Strings(kids)
	}

	var b strings.Builder
	b.WriteString(".\n")

	var walk func(parent, prefix string)
	walk = func(parent, prefix string) {
		kids := children[parent]
		for i, k := range kids {
			last := i == len(kids)-1
			branch, cont := "├── ", "│   "
			if last {
				branch, cont = "└── ", "    "
			}
			name := path.Base(k)
			if nodes[k] {
				name += "/"
			}
			b.WriteString(prefix)
			b.WriteString(branch)
			b.WriteString(name)
			b.WriteString("\n")
			if nodes[k] {
				walk(k, prefix+cont)
			}
		}
	}
	walk(".", "")
	return b.String()
}
