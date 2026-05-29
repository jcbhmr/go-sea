package main

import (
	"archive/zip"
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"strings"

	"go.jcbhmr.com.internal/go-sea/internal/fscomb"
	"go.jcbhmr.com.internal/go-sea/internal/platform"
)

const goVersion = "1.26.3"

func AllFS(v string) (fs.FS, error) {
	fsys, err := ToolchainFS(v, "linux", "amd64")
	if err != nil {
		return nil, err
	}
	return fscomb.Filter(fsys, func(name string) bool {
		return name == "." || (name != "bin" && !strings.HasPrefix(name, "bin/") && name != "pkg" && !strings.HasPrefix(name, "pkg/"))
	}), nil
}

func PlatformFS(v string, goos string, goarch string) (fs.FS, error) {
	fsys, err := ToolchainFS(v, "linux", "amd64")
	if err != nil {
		return nil, err
	}
	return fscomb.Filter(fsys, func(name string) bool {
		return name == "." || (name == "bin" || !strings.HasPrefix(name, "bin/")) || (name == "pkg" || !strings.HasPrefix(name, "pkg/"))
	}), nil
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()
	flag.Parse()

	err = func() (err error) {
		fsys, err := AllFS(goVersion)
		if err != nil {
			return
		}

		f, err := os.Create("data/all.zip")
		if err != nil {
			return
		}
		defer func() {
			err = errors.Join(err, f.Close())
		}()

		zw := zip.NewWriter(f)
		defer func() {
			err = errors.Join(err, zw.Close())
		}()
		return zw.AddFS(fsys)
	}()

	for _, p := range platform.List {
		if platform.FirstClass(p.GOOS, p.GOARCH) {
			err = func() (err error) {
				fsys, err := PlatformFS(goVersion, p.GOOS, p.GOARCH)
				if err != nil {
					return err
				}

				return os.CopyFS("data/"+p.GOOS+"-"+p.GOARCH, fsys)
			}()
		}
	}
}
