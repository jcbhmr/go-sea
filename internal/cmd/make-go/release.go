package main

import (
	"fmt"
	"net/url"
)

func platformURL(v string, goos string, goarch string) *url.URL {
	archiveExt := ".tar.gz"
	if goos == "windows" {
		archiveExt = ".zip"
	}
	// "arm" is a special case for some reason.
	if goarch == "arm" {
		goarch = "armv6l"
	}
	u, err := url.Parse(fmt.Sprintf("https://go.dev/dl/%s.%s-%s%s", v, goos, goarch, archiveExt))
	if err != nil {
		panic(err)
	}
	return u
}

func srcURL(v string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("https://go.dev/dl/%s.src.tar.gz", v))
	if err != nil {
		panic(err)
	}
	return u
}
