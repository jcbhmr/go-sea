package data

import (
	"archive/zip"
	"bytes"
	_ "embed"
)

//go:embed assets.zip
var assetsBytes []byte
var Assets = func() *zip.Reader {
	r := bytes.NewReader(assetsBytes)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		panic(err)
	}
	return zr
}()
