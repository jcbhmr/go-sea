package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var goos string
var goarch string

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	gosw := goSourceWorkspace{Dir: must(filepath.Abs("./go"))}
	goswVersion, err := gosw.Version()
	if err != nil {
		log.Fatalf("Go source workspace %q Version(): %v", gosw.Dir, err)
	}

	log.Printf("Host:       %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Printf("Target:     %s/%s", goos, goarch)
	log.Printf("Go version: %s")

	log.Printf("Running 'go generate' to pre-build non-Go-source dependencies...")
	cmd := exec.Command("go", "generate", "-x", "./...")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	log.Printf("> %v", cmd)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Running 'go build ./cmd/go' to produce './.out/linux-amd64/go' SEA binary...")
	cmd = exec.Command("go", "build", "./cmd/go")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	log.Printf("> %v", cmd)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Copied 'gofmt' from Go build output to './gofmt'. It is already a Single Executable Application.")
}
