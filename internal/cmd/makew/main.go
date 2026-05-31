/*
build-go builds a copy of [The Go Programming Language] (Go) source tree.

[The Go Programming Language]: https://go.dev/
*/
package main

import (
	"errors"
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
	"go.jcbhmr.com.internal/go-sea/internal/versionfile"
	"go.jcbhmr.com.internal/go-sea/internal/httputil"
)

var chdir string
var goos string
var goarch string
var restore bool
var distpack bool

func init() {
	flag.StringVar(&chdir, "chdir", "", "chdir to this directory first")
	flag.StringVar(&goos, "goos", runtime.GOOS, "target GOOS")
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "target GOARCH")
	flag.BoolVar(&restore, "restore", false, "restore instead")
	flag.BoolVar(&distpack, "distpack", false, "")
}

func main() {
	err := main2()
	if err != nil {
		log.Fatal(err)
	}
}

func main2() (err error) {
	flag.Parse()

	if chdir != "" {
		err = os.Chdir(chdir)
		if err != nil {
			return err
		}
	}

	var version *versionfile.File
	name := "../VERSION"
	data, err := os.ReadFile(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// continue
		} else {
			return err
		}
	} else {
		var err error
		version, err = versionfile.Parse(name, data)
		if err != nil {
			return err
		}
	}

	if restore {
		if version == nil {
			return fmt.Errorf("no file %q", name)
		}

		platformURLValue := platformURL(version.Version.Version, goos, goarch)
		platformTempName := filepath.Join(os.TempDir(), path.Base(platformURLValue.Path))
		exists, err := regularExists(platformTempName)
		if err != nil {
			return err
		}
		if !exists {
			err = httputil.Download(platformTempName, platformURLValue.String())
			if err != nil {
				return err
			}
		}

		err = os.Rename("../bin/go", "../bin/go.bak")
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// continue
			} else {
				return err
			}
		} else {
			defer func() {
				exists, err2 := regularExists("../bin/go")
				if err2 != nil {
					err = errors.Join(err, err2)
					return
				}
				if exists {
					err = errors.Join(err, os.Remove("../bin/go.bak"))
				} else {
					err = errors.Join(err, os.Rename("../bin/go.bak", "../bin/go"))
				}
			}()
		}

		err = os.Rename("../bin/gofmt", "../bin/gofmt.bak")
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// continue
			} else {
				return err
			}
		} else {
			defer func() {
				exists, err2 := regularExists("../bin/gofmt")
				if err2 != nil {
					err = errors.Join(err, err2)
					return
				}
				if exists {
					err = errors.Join(err, os.Remove("../bin/gofmt.bak"))
				} else {
					err = errors.Join(err, os.Rename("../bin/gofmt.bak", "../bin/gofmt"))
				}
			}()
		}

		err = func() (err error) {
			f, err := os.Open(platformTempName)
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, f.Close())
			}()

			gr, err := gzran.NewReader(f)
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, gr.Close())
			}()

			gra := newReaderAt(gr)

			fsys, err := tarfs.New(gra)
			if err != nil {
				return err
			}

			fsys, err = fs.Sub(fsys, "go")
			if err != nil {
				panic(err)
			}

			return copyFSIgnoreExists("..", fsys)
		}()
		if err != nil {
			return err
		}

		if goos != runtime.GOOS || goarch != runtime.GOARCH {
			err := os.MkdirAll("../bin/"+goos+"_"+goarch, 0o777)
			if err != nil {
				return err
			}

			err = os.Rename("../bin/go", "../bin/"+goos+"_"+goarch+"/go")
			if err != nil {
				return err
			}
			err = os.Rename("../bin/gofmt", "../bin/"+goos+"_"+goarch+"/gofmt")
			if err != nil {
				return err
			}
		}

		if distpack {
			err = os.MkdirAll("../pkg/distpack", 0o777)
			if err != nil {
				return err
			}

			u := srcURL(version.Version.Version)
			name := filepath.Join("../pkg/distpack", path.Base(u.Path))
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

			name = filepath.Join("../pkg/distpack", path.Base(platformURLValue.Path))
			exists, err = regularExists(name)
			if err != nil {
				return err
			}
			if !exists {
				err = copyFile(name, platformTempName)
				if err != nil {
					return err
				}
			}

			u = toolchainInfoURL(version.Version.Version, goos, goarch)
			name = filepath.Join("../pkg/distpack", path.Base(u.Path))
			exists, err = regularExists(name)
			if err != nil {
				return err
			}
			if !exists {
				err = httputil.Download(name, u.String())
				if err != nil {
					return err
				}
			}

			u = toolchainModURL(version.Version.Version, goos, goarch)
			name = filepath.Join("../pkg/distpack", path.Base(u.Path))
			exists, err = regularExists(name)
			if err != nil {
				return err
			}
			if !exists {
				err = httputil.Download(name, u.String())
				if err != nil {
					return err
				}
			}

			u = toolchainZipURL(version.Version.Version, goos, goarch)
			name = filepath.Join("../pkg/distpack", path.Base(u.Path))
			exists, err = regularExists(name)
			if err != nil {
				return err
			}
			if !exists {
				err = httputil.Download(name, u.String())
				if err != nil {
					return err
				}
			}
		}
	} else {
		panic("todo")
	}
	return nil
}
