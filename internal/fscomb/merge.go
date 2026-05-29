package fscomb

import (
	"io/fs"

	"github.com/dschmidt/go-layerfs"
)

func Merge(fsyss ...fs.FS) fs.FS {
	return layerfs.New(fsyss...)
}
