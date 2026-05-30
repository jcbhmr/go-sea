package all

import (
	"archive/zip"
	"bytes"
	_ "embed"
)

//go:generate go run ./internal/cmd/generate
//go:embed all.zip
var data []byte
var FS = func() *zip.Reader {
	br := bytes.NewReader(data)
	zr, err := zip.NewReader(br, br.Size())
	if err != nil {
		panic(err)
	}
	return zr
}()
