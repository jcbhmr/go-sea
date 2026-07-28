package httputil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(dst string, srcURL string) (f *os.File, err error) {
	response, err := http.Get(srcURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, response.Body.Close())
	}()

	if got, want := response.StatusCode, http.StatusOK; got != want {
		return nil, fmt.Errorf("%s got %d, want %d", response.Request.URL, got, want)
	}

	if dst == "" {
		f, err = os.CreateTemp("", "")
	} else {
		f, err = os.Create(dst)
	}
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	return f, err
}
