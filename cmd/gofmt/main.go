//go:build (aix && ppc64) || (darwin && (amd64 || arm64)) || (dragonfly && amd64) || (freebsd && (386 || amd64 || arm || arm64)) || (illumos && amd64) || (linux && (386 || amd64 || arm64 || arm || loong64 || mips || mips64 || mips64le || mipsle || ppc64 || ppc64le || riscv64 || s390x)) || (netbsd && (386 || amd64 || arm || arm64)) || (openbsd && (386 || amd64 || arm || arm64 || ppc64 || riscv64)) || (plan9 && (386 || amd64 || arm)) || (solaris && amd64) || (windows && (386 || amd64 || arm64))

/*
go cache directory is ${cache}/go-sea/${version}/go/.
gofmt cache directory is ${cache}/go-sea/${version}/gofmt/.

	if go_cache_directory successfully unpacked:
		exec ${go_cache_directory}/go-sea/${version}/go/bin/gofmt(.exe)
	
	unpack embedded gofmt to ${gofmt_cache_directory}
*/
package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"
)

func pollStat(ctx context.Context, name string, d time.Duration) (fs.FileInfo, error) {
	for {
		info, err := os.Stat(name)
		if err == nil {
			return info, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(d):
		}
	}
}

func main() {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	appCacheDir := filepath.Join(userCacheDir, "go-sea", "1.26.3")

	if _, err := os.Stat(filepath.Join(appCacheDir, ".unpacked-success")); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatal(err)
		}

		if parent := filepath.Dir(appCacheDir); parent != appCacheDir {
			err := os.MkdirAll(parent, 0o777)
			if err != nil {
				log.Fatal(err)
			}
		}
		err := os.Mkdir(appCacheDir, 0o777)
		if err == nil {
			err := copyOverwriteFS(appCacheDir, os.DirFS("/usr/local/go"))
			if err != nil {
				os.RemoveAll(appCacheDir)
				log.Fatal(err)
			}
		} else {
			if !errors.Is(err, fs.ErrExist) {
				log.Fatal(err)
			}

			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err := pollStat(ctx, filepath.Join(appCacheDir, ".unpacked-success"), 50*time.Millisecond)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}

		err = os.WriteFile(filepath.Join(appCacheDir, ".unpacked-success"), nil, 0o666)
		if err != nil {
			log.Fatal(err)
		}
	}

	var exeSuffix string
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}
	cmd := &exec.Cmd{
		Path:   filepath.Join(appCacheDir, "bin", "go"+exeSuffix),
		Args:   os.Args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	runtime.LockOSThread()
	signal.Ignore(os.Interrupt)
	err = cmd.Run()
	if err != nil {
		if _, ok := errors.AsType[*exec.ExitError](err); ok {
			// continue
		} else {
			log.Fatal(err)
		}
	}
	os.Exit(cmd.ProcessState.ExitCode())
}
