//go:build (darwin && (amd64 || arm64)) || (linux && (386 || amd64 || arm || arm64)) || (windows && (386 || amd64))

package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"go.jcbhmr.com/crossexec"
)

func main() {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	cacheDir := filepath.Join(userCacheDir, "go-sea")

	err = os.MkdirAll(cacheDir, 0o777)
	if err != nil {
		log.Fatalf("mkdir all path=%v perm=%o: %v", cacheDir, 0o777, err)
	}

	dir := filepath.Join(cacheDir, "go1.26.3")
	err = os.Mkdir(dir, 0o777)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			// continue
		} else {
			log.Fatalf("mkdir path=%v perm=%o: %v", dir, 0o777, err)
		}
	} else {
		err = os.CopyFS(dir, fsys)
		if err != nil {
			log.Fatalf("copy fs dir=%v fsys=%v: %v", dir, fsys, err)
		}
	}

	argv0 := filepath.Join(dir, "bin", "go")
	if runtime.GOOS == "windows" {
		argv0 += ".exe"
	}
	crossexec.Exec(argv0, os.Args, os.Environ())
}
