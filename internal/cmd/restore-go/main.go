package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"go.jcbhmr.com.internal/go-sea/internal/platform"
)

var dir string
var goos string
var goarch string
var go386 string
var goamd64 string
var goarm string
var goarm64 string
var gomips string
var gomips64 string
var goppc64 string
var goriscv64 string
var gowasm string

func init() {
	flag.StringVar(&dir, "dir", mustGetwd(), "copy of https://go.googlesource.com/go")
	flag.StringVar(&goos, "goos", runtime.GOOS, "target GOOS")
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "target GOARCH")
	flag.StringVar(&go386, "go386", "", "GOARCH=386 target GO386")
	flag.StringVar(&goamd64, "goamd64", "", "GOARCH=amd64 target GOAMD64")
	flag.StringVar(&goarm, "goarm", "", "GOARCH=arm target GOARM")
	flag.StringVar(&goarm64, "goarm64", "", "GOARCH=arm64 target GOARM64")
	flag.StringVar(&gomips, "gomips", "", "GOARCH=mips{,le} target GOMIPS")
	flag.StringVar(&gomips64, "gomips64", "", "GOARCH=mips{,le} target GOMIPS64")
	flag.StringVar(&goppc64, "goppc64", "", "GOARCH=ppc64{,le} target GOPPC64")
	flag.StringVar(&goriscv64, "goriscv64", "", "GOARCH=riscv64 target GORISCV64")
	flag.StringVar(&gowasm, "gowasm", "", "GOARCH=wasm target GOWASM")
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func parse() {
	flag.Parse()
	if !slices.Contains(platform.List, platform.OSArch{GOOS: goos, GOARCH: goarch}) {
		log.Fatalf("invalid platform %q", goos+"/"+goarch)
	}
}

func main() {
	parse()

	envv := os.Environ()
	envv = append(envv, "GOOS="+goos, "GOARCH="+goarch, "GO"+strings.ToUpper(goarch)+"="+gogoarch)

	dst := filepath.Join(mustUserCacheDir())
}
