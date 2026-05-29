package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

// go help mod download
type Module struct {
	Path     string // module path
	Query    string // version query corresponding to this version
	Version  string // module version
	Error    string // error loading module
	Info     string // absolute path to cached .info file
	GoMod    string // absolute path to cached .mod file
	Zip      string // absolute path to cached .zip file
	Dir      string // absolute path to cached source root directory
	Sum      string // checksum for path, version (as in go.sum)
	GoModSum string // checksum for go.mod (as in go.sum)
	Origin   any    // provenance of module
	Reuse    bool   // reuse of old module info is safe
}

func GoModDownload(module string) (*Module, error) {
	cmd := exec.Command("go", "mod", "download", "-json", module)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(stdout)
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("%v start: %w", err)
	}

	var m *Module
	err = dec.Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("%v decode JSON: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("%v wait: %w\nstderr:\n%s", cmd, err, stderr)
	}

	return m, nil
}

func ToolchainFS(v string, goos string, goarch string) (fs.FS, error) {
	vTrimmed := strings.TrimPrefix(v, "v")
	m, err := GoModDownload("golang.org/toolchain@v0.0.1-go" + vTrimmed + "." + goos + "-" + goarch)
	if err != nil {
		return nil, err
	}
	return os.DirFS(m.Dir), nil
}
