package main

import (
	"errors"
	"io"
	"sync"
)

type readerAt struct {
	r  io.ReadSeeker
	mu sync.Mutex
}

func newReaderAt(r io.ReadSeeker) *readerAt {
	return &readerAt{r: r}
}

func (ra *readerAt) Read(p []byte) (n int, err error) {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	return ra.r.Read(p)
}

func (ra *readerAt) Seek(offset int64, whence int) (int64, error) {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	return ra.r.Seek(offset, whence)
}

func (ra *readerAt) ReadAt(p []byte, off int64) (n int, err error) {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	old, err := ra.r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	defer func() {
		_, err2 := ra.r.Seek(old, io.SeekStart)
		err = errors.Join(err, err2)
	}()

	_, err = ra.r.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return ra.r.Read(p)
}
