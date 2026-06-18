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
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("os.Executable() err: %v", err)
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		log.Fatalf("filepath.EvalSymlinks(%#v) err: %v", executable, err)
	}

	dir := executable + ".UNPACKED"
	// os.Mkdir is atomic on all systems.
	err = os.Mkdir(dir, 0o777)
	if err == nil {
		err := os.CopyFS(dir, fsys)
		if err != nil {
			log.Fatalf("os.CopyFS(%#v, fsys) err: %v", dir, err)
		}
	} else if !errors.Is(err, fs.ErrExist) {
		log.Fatalf("os.Mkdir(%#v, 0o777) err: %v", dir, err)
	}

	argv0 := filepath.Join(dir, "bin", "go")
	if runtime.GOOS == "windows" {
		argv0 += ".exe"
	}
	err = crossexec.Exec(argv0, os.Args, os.Environ())
	log.Fatalf("crossexec.Exec(%#v, os.Args, os.Environ()) err: %v", argv0, err)
}
