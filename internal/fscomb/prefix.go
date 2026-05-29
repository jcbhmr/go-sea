package fscomb

import (
	"io/fs"
	"path"
)

type prefixFS struct {
	fsys   fs.FS
	prefix string
}

func Prefix(fsys fs.FS, prefix string) fs.FS {
	return prefixFS{fsys: fsys, prefix: prefix}
}
func (fsys prefixFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}
	return fsys.fsys.Open(path.Clean(path.Join(fsys.prefix, name)))
}
