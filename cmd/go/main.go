//go:build (aix && ppc64) || (darwin && (amd64 || arm64)) || (dragonfly && amd64) || (freebsd && (386 || amd64 || arm || arm64)) || (illumos && amd64) || (linux && (386 || amd64 || arm64 || arm || loong64 || mips || mips64 || mips64le || mipsle || ppc64 || ppc64le || riscv64 || s390x)) || (netbsd && (386 || amd64 || arm || arm64)) || (openbsd && (386 || amd64 || arm || arm64 || ppc64 || riscv64)) || (plan9 && (386 || amd64 || arm)) || (solaris && amd64) || (windows && (386 || amd64 || arm64))

package main

import (
	"app/data"
	"app/data/binfs"
	"app/internal/crossexec"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/dschmidt/go-layerfs"
	"github.com/josharian/filterfs"
	"github.com/unstoppablemango/ihfs/prefixfs"
)

var exeSuffix = func() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	} else {
		return ""
	}
}()

var GoFS = layerfs.New(
	data.Assets,
	prefixfs.New(filterfs.ExcludePaths(binfs.BinFS, "asm"+exeSuffix, "cgo"+exeSuffix, "compile"+exeSuffix, "cover"+exeSuffix, "fix"+exeSuffix, "link"+exeSuffix, "preprofile"+exeSuffix, "vet"+exeSuffix), "bin"),
	prefixfs.New(filterfs.ExcludePaths(binfs.BinFS, "go"+exeSuffix, "gofmt"+exeSuffix), "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH),
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	appCacheDir := filepath.Join(userCacheDir, "go-sea", "1.26.5")

	_, err = os.Stat(appCacheDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = nil

			err = os.CopyFS(appCacheDir, GoFS)
			if err != nil {
				return
			}
			err = errors.Join(
				os.Chmod(filepath.Join(appCacheDir, "bin/go"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "bin/gofmt"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/asm"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/cgo"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/compile"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/cover"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/fix"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/link"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/preprofile"+exeSuffix), 0o777),
				os.Chmod(filepath.Join(appCacheDir, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH+"/vet"+exeSuffix), 0o777),
			)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	err = crossexec.Exec(filepath.Join(appCacheDir, "bin", "go"+exeSuffix), os.Args, os.Environ())
}
