package data

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"io/fs"

	"go.jcbhmr.com.internal/go-sea/internal/fscomb"
)

//go:generate go run ./cmd/generate-all
//go:embed all.zip
var all_zip []byte

var all = func() fs.FS {
	br := bytes.NewReader(all_zip)
	zr, err := zip.NewReader(br, br.Size())
	if err != nil {
		panic(err)
	}
	return zr
}()

var FS = fscomb.Merge(all, platform)
