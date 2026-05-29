package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.jcbhmr.com.internal/go-sea/internal/platform"
)

type GoSource struct {
	root *os.Root
}

func GoSourceFrom(dir string) (*GoSource, error) {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	root, err := os.OpenRoot(dirAbs)
	if err != nil {
		return nil, err
	}
	return &GoSource{root: root}, nil
}

func (w *GoSource) GOROOT() string {
	return w.root.Name()
}

func (w *GoSource) FS() fs.FS {
	return w.root.FS()
}

type VERSIONFile struct {
	Version string // ex: "go1.26.0"
	Time    time.Time
}

// ParseVERSIONFile parses a VERSION file.
// The first line of the file is the Go version.
// Additional lines are 'key value' pairs setting other data.
// The only valid key at the moment is 'time', which sets the modification time for file archives.
func ParseVERSIONFile(name string, data []byte) (*VERSIONFile, error) {
	var version string
	var t time.Time
	var err error
	version, rest, _ := strings.Cut(string(data), "\n")
	for line := range strings.SplitSeq(rest, "\n") {
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		default:
			return nil, fmt.Errorf("%s: unexpected line: %s", name, line)
		case "time":
			if len(f) != 2 {
				return nil, fmt.Errorf("%s: unexpected time line: %s", name, line)
			}
			t, err = time.ParseInLocation(time.RFC3339, f[1], time.UTC)
			if err != nil {
				return nil, fmt.Errorf("%s: bad time: %s", name, err)
			}
		}
	}
	return &VERSIONFile{Version: version, Time: t}, nil
}

// ReadVERSION reads the VERSION file.
// The first line of the file is the Go version.
// Additional lines are 'key value' pairs setting other data.
// The only valid key at the moment is 'time', which sets the modification time for file archives.
func (w *GoSource) ReadVERSION() (*VERSIONFile, error) {
	name := filepath.Join(w.GOROOT(), "VERSION")
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return ParseVERSIONFile(name, data)
}

func (w *GoSource) Make(goos string, goarch string, extraEnv []string) (fs.FS, error) {
	
}
