package sum

import (
	"crypto/sha256"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func Sha256() (sum []byte, err error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(exe)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
