//go:build gc

package main

import (
	"context"
	"errors"
	"io"
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

func copyOverwriteFS(dir string, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fpath, err := filepath.Localize(path)
		if err != nil {
			return err
		}
		newPath := filepath.Join(dir, fpath)

		switch d.Type() {
		case fs.ModeDir:
			return os.MkdirAll(newPath, 0777)
		case fs.ModeSymlink:
			target, err := fs.ReadLink(fsys, path)
			if err != nil {
				return err
			}
			os.Remove(target)
			return os.Symlink(target, newPath)
		case 0:
			r, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()
			info, err := r.Stat()
			if err != nil {
				return err
			}
			w, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY, 0666|info.Mode()&0777)
			if err != nil {
				return err
			}

			if _, err := io.Copy(w, r); err != nil {
				w.Close()
				return &fs.PathError{Op: "Copy", Path: newPath, Err: err}
			}
			return w.Close()
		default:
			return &fs.PathError{Op: "CopyFS", Path: path, Err: fs.ErrInvalid}
		}
	})
}

func main() {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	appCacheDir := filepath.Join(userCacheDir, "gosfx", "1.26.5")

	if _, err := os.Stat(filepath.Join(appCacheDir, ".unpacked-success")); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatal(err)
		}

		err := os.Mkdir(appCacheDir, 0o777)
		if err == nil {
			err := copyOverwriteFS(appCacheDir, nil)
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
	signal.Ignore()
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
