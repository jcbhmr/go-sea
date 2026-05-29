package fscomb_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"go.jcbhmr.com.internal/go-sea/internal/fscomb"
)

var FS = fstest.MapFS{
	"data.txt":         &fstest.MapFile{Data: []byte("Hi there! It's me, Alan Turing!\n"), Mode: 0o666},
	"folder/info.json": &fstest.MapFile{Data: []byte(`{"look":"some json"}`), Mode: 0o666},
}

func TestFilter_skipFolder(t *testing.T) {
	fsys := fscomb.Filter(FS, func(name string) bool {
		r := name == "." || (name != "folder" && !strings.HasPrefix(name, "folder/"))
		// slog.Info("filter", "name", name, "r", r)
		return r
	})
	err := fstest.TestFS(fsys, "folder/info.json")
	if err == nil {
		t.Fatalf("expected err for %s", "folder/info.json")
	}
	err = fstest.TestFS(fsys, "data.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilter_onlyFolder(t *testing.T) {
	fsys := fscomb.Filter(FS, func(name string) bool {
		return name == "." || name == "folder" || strings.HasPrefix(name, "folder/")
	})
	err := fstest.TestFS(fsys, "folder/info.json")
	if err != nil {
		t.Fatal(err)
	}
	err = fstest.TestFS(fsys, "data.txt")
	if err == nil {
		t.Fatalf("expected err for %s", "data.txt")
	}
}
