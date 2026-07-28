package main

import (
	"app/internal/httputil"
	"archive/zip"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type AdderFS interface {
	AddFS(fsys fs.FS) error
}

func Restore(dst AdderFS, goversion string, goos string, goarch string) (err error) {
	src := fmt.Sprintf("https://go.dev/dl/mod/golang.org/toolchain/@v/v0.0.1-%s.%s-%s", goversion, goos, goarch)
	f, err := httputil.Download("", src)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
		err = errors.Join(err, os.Remove(f.Name()))
	}()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(f, info.Size())
	if err != nil {
		return err
	}

	binFS, err := fs.Sub(zipReader, "bin")
	if err != nil {
		return err
	}

	toolsFS, err := fs.Sub(zipReader, fmt.Sprintf("pkg/tools/%s_%s", goos, goarch))
	if err != nil {
		return err
	}

	return errors.Join(dst.AddFS(binFS), dst.AddFS(toolsFS))
}
