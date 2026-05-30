package main

import (
	"go.jcbhmr.com.internal/go-sea/cmd/go/internal/all"
	linux_amd64 "go.jcbhmr.com.internal/go-sea/cmd/go/internal/linux-amd64"
	"go.jcbhmr.com.internal/go-sea/internal/fscomb"
)

var fsys = fscomb.Merge(linux_amd64.FS, all.FS)
