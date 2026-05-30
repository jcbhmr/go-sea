package httputil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(dst string, url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		return fmt.Errorf("%s got %d, want %d", resp.Request.URL, got, want)
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
