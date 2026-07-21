package osutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyFSExcl(dir string, fsys fs.FS) error {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	dirDir, dirBase := filepath.Split(dirAbs)
	if dirDir == dirAbs {
		return fmt.Errorf("%v has no parent", dirAbs)
	}

	dirSib, err := os.MkdirTemp(dirDir, dirBase)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dirSib)

	err = os.CopyFS(dirSib, fsys)
	if err != nil {
		return err
	}

	return os.Rename(dirSib, dirAbs)
}

func IsDir(name string) (bool, error) {
	info, err := os.Lstat(name)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func IsRegular(name string) (bool, error) {
	info, err := os.Lstat(name)
	if err != nil {
		return false, err
	}

	return info.Mode().IsRegular(), nil
}
