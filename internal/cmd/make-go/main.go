/*
build-go builds a copy of [The Go Programming Language] (Go) source tree.

[The Go Programming Language]: https://go.dev/
*/
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/nlepage/go-tarfs"
	"github.com/timpalpant/gzran"
	"go.jcbhmr.com.internal/go-sea/internal/golang/versionfile"
	"go.jcbhmr.com.internal/go-sea/internal/httputil"
)

var goos string
var goarch string
var restore bool
var distpack bool

func init() {
	flag.StringVar(&goos, "goos", runtime.GOOS, "target GOOS")
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "target GOARCH")
	flag.BoolVar(&restore, "restore", false, "restore instead")
	flag.BoolVar(&distpack, "distpack", false, "")
}

func main() {
	err := mainValue()
	if err != nil {
		log.Fatal(err)
	}
}

func mainValue() (err error) {
	flag.Parse()

	name := "VERSION"
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	version, err := versionfile.Parse(name, data)
	if err != nil {
		return err
	}

	if restore {
		if version == nil {
			return fmt.Errorf("no file %q", name)
		}

		u := platformURL(version.Version.Version, goos, goarch)
		name := filepath.Join(os.TempDir(), path.Base(u.Path))
		exists, err := regularExists(name)
		if err != nil {
			return err
		}
		if !exists {
			err = httputil.Download(name, u.String())
			if err != nil {
				return err
			}
		}

		f, err := os.Open(name)
		if err != nil {
			return err
		}

		gr, err := gzran.NewReader(f)
		if err != nil {
			return err
		}

		gra := newReaderAt(gr)

		fsys, err := tarfs.New(gra)
		if err != nil {
			return err
		}

		fsys, err = fs.Sub(fsys, "go")
		if err != nil {
			panic(err)
		}

		err = copyFSIgnoreExists(".", fsys)
		if err != nil {
			return err
		}
	} else {

	}
	return nil
}
