package main

import (
	"archive/zip"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.jcbhmr.com.internal/go-sea/internal/fscomb"
)

func main() {
	modDir, err := filepath.Abs("../../../../")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("./scripts/build-go.sh")
	cmd.Dir = modDir
	cmd.Env = append(cmd.Environ(), "GOOS=linux", "GOARCH=amd64")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	log.Printf("+ %v", cmd)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(modDir, "go/go/VERSION"))
	if err != nil {
		log.Fatal(err)
	}
	var goversion string
	for l := range strings.Lines(string(data)) {
		goversion = strings.TrimRight(l, "\r\n")
		break
	}

	zr, err := zip.OpenReader(filepath.Join(modDir, "go/go/pkg/distpack", "v0.0.1-"+goversion+".linux-amd64.zip"))
	if err != nil {
		log.Fatal(err)
	}
	fsys, err := fs.Sub(zr, "golang.org/toolchain@v0.0.1-"+goversion+".linux-amd64")
	if err != nil {
		panic(err)
	}

	fsys = fscomb.Filter(fsys, func(name string) bool {
		if name == "." {
			return true
		}
		return !(name == "pkg/tool/linux_amd64" || strings.HasPrefix(name, "pkg/tool/linux_amd64/")) && !(name == "bin" || strings.HasPrefix(name, "bin/"))
	})

	f, err := os.Create("all.zip")
	if err != nil {
		log.Fatal(err)
	}
	zw := zip.NewWriter(f)
	err = zw.AddFS(fsys)
	if err != nil {
		log.Fatal(err)
	}
	zw.Close()
	f.Close()
}
