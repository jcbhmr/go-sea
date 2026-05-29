package fscomb

import (
	"io"
	"io/fs"
	"path"
	"slices"
)

type filterFS struct {
	fs     fs.FS
	filter func(name string) bool
}

func Filter(fsys fs.FS, filter func(name string) bool) fs.FS {
	return filterFS{fs: fsys, filter: filter}
}
func (fsys filterFS) Open(name string) (f fs.File, err error) {
	var clean string
	// defer func() {
	// 	slog.Info("open", "fsys", fsys, "name", name, "clean", clean, "f?", f != nil, "err", err)
	// }()
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}
	clean = path.Clean(name)
	if !fsys.filter(clean) {
		return nil, fs.ErrNotExist
	}
	f, err = fsys.fs.Open(clean)
	if err != nil {
		return nil, err
	}
	if rdf, ok := f.(fs.ReadDirFile); ok {
		return filterReadDirFile{rdf: rdf, dir: clean, filter: fsys.filter}, nil
	}
	return f, nil
}

type filterReadDirFile struct {
	rdf    fs.ReadDirFile
	dir    string
	filter func(name string) bool
}

func (f filterReadDirFile) Stat() (fs.FileInfo, error) {
	return f.rdf.Stat()
}
func (f filterReadDirFile) Read(b []byte) (int, error) {
	return f.rdf.Read(b)
}
func (f filterReadDirFile) Close() error {
	return f.rdf.Close()
}
func (f filterReadDirFile) ReadDir(n int) (entries []fs.DirEntry, err error) {
	// defer func() {
	// 	entriesNames := make([]string, len(entries))
	// 	for i, v := range entries {
	// 		entriesNames[i] = v.Name()
	// 	}
	// 	entriesPaths := make([]string, len(entries))
	// 	for i, v := range entries {
	// 		entriesPaths[i] = path.Clean(path.Join(f.dir, v.Name()))
	// 	}
	// 	slog.Info("readdir", "f.dir", f.dir, "n", n, "entries names", entriesNames, "entries paths", entriesPaths, "err", err)
	// }()
	entries = []fs.DirEntry{}
	remaining := n
	for {
		some, err := f.rdf.ReadDir(remaining)
		some = slices.DeleteFunc(some, func(e fs.DirEntry) bool {
			p := path.Clean(path.Join(f.dir, e.Name()))
			return !f.filter(p)
		})
		entries = append(entries, some...)
		remaining -= len(some)
		if err != nil {
			return entries, err
		}
		if n <= 0 {
			return entries, nil
		}
		if remaining <= 0 {
			return entries, io.EOF
		}
	}
}
