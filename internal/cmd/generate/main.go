package main

import (
	"app/gocli/list"
	"app/gocli/mod"
	"archive/zip"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing/fstest"

	"github.com/josharian/filterfs"
)

type Platform struct {
	OS   string
	Arch string
}

const goversion = "go1.26.5"

var toolchainVersionRE = regexp.MustCompile(`^v0\.0\.1-(go1\..*?)\.([a-z0-9]+)-([a-z0-9]+)$`)

func platforms(goversion string) ([]Platform, error) {
	m, err := list.Versions("golang.org/toolchain")
	if err != nil {
		return nil, err
	}

	platforms := []Platform{}
	for _, v := range m.Versions {
		match := toolchainVersionRE.FindStringSubmatch(v)
		if match == nil {
			continue
		}
		if match[1] != goversion {
			continue
		}
		platforms = append(platforms, Platform{
			OS:   match[2],
			Arch: match[3],
		})
	}
	return platforms, nil
}

func assetsZip(dst string, goversion string) error {
	platforms, err := platforms(goversion)
	if err != nil {
		return err
	}
	if len(platforms) < 1 {
		return fmt.Errorf("no platforms for %s", goversion)
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
		if err != nil {
			err = errors.Join(err, os.Remove(f.Name()))
		}
	}()

	zipWriter := zip.NewWriter(f)
	defer func() {
		err = errors.Join(err, zipWriter.Close())
	}()

	module := fmt.Sprintf("golang.org/toolchain@v0.0.1-%s.%s-%s", goversion, "linux", "amd64")
	m, err := mod.Download(module)
	if err != nil {
		return err
	}

	dirFS := os.DirFS(m.Dir)

	err = zipWriter.AddFS(filterfs.ExcludePaths(dirFS, "bin", "pkg/tool", "CONTRIBUTING.md", "codereview.cfg", "SECURITY.md"))
	if err != nil {
		return err
	}

	mods := fstest.MapFS{}
	err = fs.WalkDir(dirFS, ".", func(path2 string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "_go.mod" {
			return nil
		}
		data, err := fs.ReadFile(dirFS, path2)
		if err != nil {
			return err
		}
		info, err := fs.Stat(dirFS, path2)
		if err != nil {
			return err
		}
		mods[path.Join(path.Dir(path2), "go.mod")] = &fstest.MapFile{
			Data:    data,
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
		}
		return nil
	})
	if err != nil {
		return err
	}
	return zipWriter.AddFS(mods)
}

func binDir(dst string, goversion string, goos string, goarch string) (err error) {
	module := fmt.Sprintf("golang.org/toolchain@v0.0.1-%s.%s-%s", goversion, goos, goarch)
	m, err := mod.Download(module)
	if err != nil {
		return err
	}

	dirFS := os.DirFS(m.Dir)

	binFS, err := fs.Sub(dirFS, "bin")
	if err != nil {
		return err
	}

	toolFS, err := fs.Sub(dirFS, fmt.Sprintf("pkg/tool/%s_%s", goos, goarch))
	if err != nil {
		return err
	}

	err = os.RemoveAll(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, os.RemoveAll(dst))
		}
	}()
	
	err = errors.Join(os.CopyFS(dst, binFS), os.CopyFS(dst, toolFS))
	if err != nil {
		return err
	}

	return filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		return os.Chmod(path, 0o777)
	})
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = os.MkdirAll("data", 0o777)
	if err != nil {
		return
	}

	err = assetsZip("data/assets.zip", goversion)
	if err != nil {
		return
	}
	log.Printf("Created data/assets.zip")

	err = os.WriteFile("data/assets.go", []byte(`package data

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
`), 0o666)
	if err != nil {
		return
	}

	platforms, err := platforms(goversion)
	if err != nil {
		return
	}
	for _, p := range platforms {
		err = binDir(filepath.Join("data", fmt.Sprintf("%s-%s", p.OS, p.Arch)), goversion, p.OS, p.Arch)
		if err != nil {
			return
		}
		log.Printf("Created data/%s-%s/", p.OS, p.Arch)

		exeSuffix := ""
		if p.OS == "windows" {
			exeSuffix = ".exe"
		}
		err = os.WriteFile(filepath.Join("data", fmt.Sprintf("%s-%s", p.OS, p.Arch), "all.go"), fmt.Appendf(nil, `package %s_%s

import (
	"embed"
)

//go:embed asm%[3]s cgo%[3]s compile%[3]s cover%[3]s fix%[3]s go%[3]s gofmt%[3]s link%[3]s preprofile%[3]s vet%[3]s
var All embed.FS
`, p.OS, p.Arch, exeSuffix), 0o666)
		if err != nil {
			return
		}

		err = os.MkdirAll("data/binfs", 0o777)
		if err != nil {
			return
		}

		err = os.WriteFile(filepath.Join("data/binfs", fmt.Sprintf("binfs_%s_%s.go", p.OS, p.Arch)), fmt.Appendf(nil, `package binfs

import (
	"app/data/%[1]s-%[2]s"
)

var BinFS = %[1]s_%[2]s.All
`, p.OS, p.Arch), 0o666)
		if err != nil {
			return
		}
	}
}
