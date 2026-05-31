package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

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

func copyFSIgnoreExists(dir string, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fpath, err := filepath.Localize(path)
		if err != nil {
			return err
		}
		newPath := filepath.Join(dir, fpath)

		if _, err := os.Stat(newPath); err == nil {
			return nil
		}

		switch d.Type() {
		case fs.ModeDir:
			return os.MkdirAll(newPath, 0o777)
		case fs.ModeSymlink:
			target, err := fs.ReadLink(fsys, path)
			if err != nil {
				return err
			}
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
			w, err := os.OpenFile(newPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666|info.Mode()&0o777)
			if err != nil {
				return err
			}

			if _, err := io.Copy(w, r); err != nil {
				w.Close()
				return &os.PathError{Op: "Copy", Path: newPath, Err: err}
			}
			return w.Close()
		default:
			return &os.PathError{Op: "CopyFS_ignoreExists", Path: path, Err: fs.ErrInvalid}
		}
	})
}
