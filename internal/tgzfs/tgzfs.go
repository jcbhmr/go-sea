package tgzfs

import (
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nlepage/go-tarfs"
)

func mustUserCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	return dir
}

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

type tgzFS struct {
	TgzName string
	TarFile *os.File
	fs.FS
}

func (fsys *tgzFS) Close() error {
	return fsys.TarFile.Close()
}

func Open(name string) (fsys interface {
	fs.FS
	io.Closer
}, err error) {
	tgzFile, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, tgzFile.Close())
	}()

	h := sha256.New()

	_, err = io.Copy(h, tgzFile)
	if err != nil {
		return nil, err
	}

	sum := h.Sum(nil)
	hex := fmt.Sprintf("%x", sum)

	_, err = tgzFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	tarFileName := filepath.Join(mustUserCacheDir(), hex+".tar")
	exists, err := regularExists(tarFileName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := func() (err error) {
			gr, err := gzip.NewReader(tgzFile)
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, gr.Close())
			}()

			fw, err := os.Create(tarFileName)
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, fw.Close())
			}()

			_, err = io.Copy(fw, gr)
			return err
		}()
		if err != nil {
			return nil, err
		}
	}

	tarFile, err := os.Open(tarFileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tarFile.Close())
		}
	}()

	tarFS, err := tarfs.New(tarFile)
	if err != nil {
		return nil, err
	}

	return &tgzFS{TgzName: tgzFile.Name(), TarFile: tarFile, FS: tarFS}, nil
}
