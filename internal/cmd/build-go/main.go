/*
build-go builds a copy of [The Go Programming Language] (Go) source tree.

[The Go Programming Language]: https://go.dev/
*/
package main

import (
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"go.jcbhmr.com.internal/go-sea/internal/golang"
	"go.jcbhmr.com.internal/go-sea/internal/golang/release"
	"go.jcbhmr.com.internal/go-sea/internal/golang/versionfile"
	"go.jcbhmr.com.internal/go-sea/internal/httputil"
)

var dir string
var goos string
var goarch string
var fromsrc bool

func init() {
	flag.StringVar(&dir, "dir", "", "copy of https://go.googlesource.com/go")
	flag.StringVar(&goos, "goos", runtime.GOOS, "target GOOS")
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "target GOARCH")
	flag.BoolVar(&fromsrc, "fromsrc", false, "force building from source")
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func regularExists(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}
	return fi.Mode().IsRegular(), nil
}

func main() {
	flag.Parse()

	initCwd := mustGetwd()
	_ = initCwd

	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Changed directory to %q", mustGetwd())

	versionFile := func() *versionfile.File {
		data, err := os.ReadFile("VERSION")
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			} else {
				log.Fatal(err)
			}
		}

		versionFile, err := versionfile.Parse("VERSION", data)
		if err != nil {
			log.Fatal(err)
		}

		return versionFile
	}()

	{
		v := "<not present>"
		if versionFile.Version != nil {
			v = versionFile.Version.Version
		}
		log.Printf("VERSION file: %s", v)
	}

	if fromsrc {
		log.Printf("Running './make.{bash,bat,rc} -distpack' for %s/%s...", goos, goarch)
		err := golang.MakeDistpack(".", goos, goarch, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if versionFile == nil {
			log.Fatal("prebuilt release: no VERSION file")
		}

		err := os.MkdirAll("./pkg/distpack", 0o777)
		if err != nil {
			log.Fatal(err)
		}

		u := release.SrcURL(versionFile.Version.Version)
		name := filepath.Join("./pkg/distpack", path.Base(u.Path))
		exists, err := regularExists(name)
		if err != nil {
			log.Fatal(err)
		}
		if exists {
			log.Printf("%q already exists", name)
		} else {
			log.Printf("Downloading %q to %q...", u, name)
			err := httputil.Download(name, u.String())
			if err != nil {
				log.Fatal(err)
			}
		}

		u = release.URL(versionFile.Version.Version, goos, goarch)
		name = filepath.Join("./pkg/distpack", path.Base(u.Path))
		exists, err = regularExists(name)
		if err != nil {
			log.Fatal(err)
		}
		if exists {
			log.Printf("%q already exists", name)
		} else {
			log.Printf("Downloading %q to %q...", u, name)
			err := httputil.Download(name, u.String())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
