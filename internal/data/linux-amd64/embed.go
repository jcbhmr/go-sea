package linux_amd64

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"io/fs"
)

//go:generate env -C ../../ go run ./internal/cmd/generate-platform -goos linux -goarch amd64 -o ./internal/linux-amd64/linux-amd64.zip
//go:embed linux-amd64.zip
var linux_amd64 []byte
var FS = func() fs.FS {
	br := bytes.NewReader(linux_amd64)
	zr, err := zip.NewReader(br, br.Size())
	if err != nil {
		panic(err)
	}
	return zr
}()
