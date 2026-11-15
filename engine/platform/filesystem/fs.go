package filesystem

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

type DirEntry = fs.DirEntry

func Open(name string) (io.ReadCloser, error) {
	return impl_open(name)
}
func Create(name string) (io.ReadWriteCloser, error) {
	return impl_create(name)
}
func WriteFile(name string, data []byte) error {
	return impl_writeFile(name, data)
}
func ReadFile(name string) ([]byte, error) {
	return impl_readFile(name)
}
func ReadDir(name string) ([]fs.DirEntry, error) {
	return impl_readDir(name)
}
func OpenOrCreate(name string) (io.ReadCloser, error) {
	file, err := Open(name)
	if errors.Is(err, os.ErrNotExist) {
		err = file.Close()
		if err != nil {
			return nil, err
		}
		file, err = Create(name)
	}
	return file, err
}
