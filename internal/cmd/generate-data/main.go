package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.jcbhmr.com.internal/go-sea/internal/fscomb"
	"go.jcbhmr.com.internal/go-sea/internal/golang/versionfile"
	"go.jcbhmr.com.internal/go-sea/internal/platform"
	"go.jcbhmr.com.internal/go-sea/internal/tgzfs"
)

var godir string

func init() {
	flag.StringVar(&godir, "godir", "", "copy of https://go.googlesource.com/go")
}

func main() {
	flag.Parse()

	data, err := os.ReadFile(filepath.Join(godir, "VERSION"))
	if err != nil {
		log.Fatal(err)
	}

	versionFile, err := versionfile.Parse(filepath.Join(godir, "VERSION"), data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("VERSION: %s", versionFile.Version.Version)

	err = os.RemoveAll("data")
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll("data", 0o777)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("data/.gitignore", []byte("*\n"), 0o666)
	if err != nil {
		log.Fatal(err)
	}

	replaces := []string{}

	for _, p := range platform.List {
		if !platform.FirstClass(p.GOOS, p.GOARCH) {
			continue
		}
		// skip 32-bit systems
		if p.GOARCH == "386" || p.GOARCH == "arm" {
			continue
		}

		err := os.MkdirAll("data/"+p.GOOS+"-"+p.GOARCH, 0o777)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Building %q for %s/%s...", godir, p.GOOS, p.GOARCH)
		cmd := exec.Command("go", "tool", "build-go", "-dir", godir, "-goos", p.GOOS, "-goarch", p.GOARCH)
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		b := &bytes.Buffer{}
		fmt.Fprintf(b, "module %s\n", "go.jcbhmr.com.internal/go-sea/internal/data/"+p.GOOS+"-"+p.GOARCH)
		fmt.Fprintln(b)
		fmt.Fprintf(b, "go %s\n", "1.26.0")
		fmt.Fprintln(b)
		fmt.Fprintf(b, "toolchain %s\n", "go1.26.3")
		err = os.WriteFile("data/"+p.GOOS+"-"+p.GOARCH+"/go.mod", b.Bytes(), 0o666)
		if err != nil {
			log.Fatal(err)
		}

		var fsCloser interface {
			fs.FS
			io.Closer
		}
		if p.GOOS == "windows" {
			var err error
			fsCloser, err = zip.OpenReader(filepath.Join(godir, "pkg/distpack/"+versionFile.Version.Version+"."+p.GOOS+"-"+p.GOARCH+".zip"))
			if err != nil {
				log.Fatal(err)
			}
			defer fsCloser.Close()
		} else {
			var err error
			fsCloser, err = tgzfs.Open(filepath.Join(godir, "pkg/distpack/"+versionFile.Version.Version+"."+p.GOOS+"-"+p.GOARCH+".tar.gz"))
			if err != nil {
				log.Fatal(err)
			}
			defer fsCloser.Close()
		}

		fsys, err := fs.Sub(fsCloser, "go")
		if err != nil {
			panic(err)
		}
		fsys = fscomb.Filter(fsys, func(name string) bool {
			if name == "." {
				return true
			}
			return (name == "bin" || strings.HasPrefix(name, "bin/")) || (name == "pkg" || name == "pkg/tool" || name == "pkg/tool/"+p.GOOS+"_"+p.GOARCH || strings.HasPrefix(name, "pkg/tool/"+p.GOOS+"_"+p.GOARCH+"/"))
		})

		log.Printf("Copying bin/ and pkg/tool/%s_%s/ to %q", p.GOOS, p.GOARCH, "data/"+p.GOOS+"-"+p.GOARCH)
		err = os.CopyFS("data/"+p.GOOS+"-"+p.GOARCH, fsys)
		if err != nil {
			log.Fatal(err)
		}

		b = &bytes.Buffer{}
		fmt.Fprintf(b, "package %s_%s\n", p.GOOS, p.GOARCH)
		fmt.Fprintln(b)
		fmt.Fprintf(b, "import (\n\t%s\n)\n", strings.Join([]string{`"embed"`}, "\n\t"))
		fmt.Fprintln(b)
		fmt.Fprintf(b, "//go:embed all:bin all:pkg\n")
		fmt.Fprintf(b, "var FS embed.FS\n")
		err = os.WriteFile("data/"+p.GOOS+"-"+p.GOARCH+"/embed.go", b.Bytes(), 0o666)
		if err != nil {
			log.Fatal(err)
		}

		replaces = append(replaces, "go.jcbhmr.com.internal/go-sea/internal/data/"+p.GOOS+"-"+p.GOARCH+" => "+"./"+p.GOOS+"-"+p.GOARCH)
	}

	log.Printf("Building %q for %s/%s...", godir, "linux", "amd64")
	cmd := exec.Command("go", "tool", "build-go", "-dir", godir, "-goos", "linux", "-goarch", "amd64")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	b := &bytes.Buffer{}
	fmt.Fprintf(b, "module %s\n", "go.jcbhmr.com.internal/go-sea/internal/data")
	fmt.Fprintln(b)
	fmt.Fprintf(b, "go %s\n", "1.26.0")
	fmt.Fprintln(b)
	fmt.Fprintf(b, "toolchain %s\n", "go1.26.3")
	fmt.Fprintln(b)
	fmt.Fprintf(b, "replace (\n\t%s\n)\n", strings.Join(replaces, "\n\t"))
	err = os.WriteFile("data/go.mod", b.Bytes(), 0o666)
	if err != nil {
		log.Fatal(err)
	}

	b = &bytes.Buffer{}
	fmt.Fprintf(b, "package %s\n", "data")
	fmt.Fprintln(b)
	fmt.Fprintf(b, "import (\n\t%s\n)\n", strings.Join([]string{`_ "embed"`, `"archive/zip"`}, "\n\t"))
	fmt.Fprintln(b)
	fmt.Fprintf(b, "//go:embed any.zip\n")
	fmt.Fprintf(b, "var any_zip []byte\n")
	fmt.Fprintln(b)
	fmt.Fprintf(b, "var FS = func() *zip.Reader {\n\t\n}()\n")
	err = os.WriteFile("data/embed.go", b.Bytes(), 0o666)
	if err != nil {
		log.Fatal(err)
	}
}
