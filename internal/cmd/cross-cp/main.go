package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
)

var r bool

func init() {
	flag.BoolVar(&r, "r", false, "recursive")
}

func copyFile(dst string, src string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = errors.Join(err, srcFile.Close())
	}()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcFileInfo.Mode().Perm())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = errors.Join(err, dstFile.Close())
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	dst := flag.Args()[0]
	src := flag.Args()[1]

	if r {
		err := os.CopyFS(dst, os.DirFS(src))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := copyFile(dst, src)
		if err != nil {
			log.Fatal(err)
		}
	}
}
