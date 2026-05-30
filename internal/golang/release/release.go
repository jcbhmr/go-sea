package release

import (
	"fmt"
	"net/url"
)

func URL(v string, goos string, goarch string) *url.URL {
	archiveExt := ".tar.gz"
	if goos == "windows" {
		archiveExt = ".zip"
	}
	u, err := url.Parse(fmt.Sprintf("https://go.dev/dl/%s.%s-%s%s", v, goos, goarch, archiveExt))
	if err != nil {
		panic(err)
	}
	return u
}

func SrcURL(v string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("https://go.dev/dl/%s.src.tar.gz", v))
	if err != nil {
		panic(err)
	}
	return u
}
