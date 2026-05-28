package src

import (
	"embed"
	"io/fs"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

//go:embed all:go
var all embed.FS
var FS = must(fs.Sub(all, "go"))
