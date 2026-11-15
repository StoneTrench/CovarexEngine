//go:build !js

package filesystem

import (
	"io"
	"os"
)

func impl_open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}
func impl_create(name string) (io.ReadWriteCloser, error) {
	return os.Create(name)
}
func impl_writeFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}
func impl_readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
func impl_readDir(name string) ([]DirEntry, error) {
	return os.ReadDir(name)
}
