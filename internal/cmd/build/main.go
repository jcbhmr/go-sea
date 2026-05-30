package main

import (
	"flag"
	"log"
	"os/exec"
)

var distpack bool

func init() {
	flag.BoolVar(&distpack, "distpack", false, "build release dist pack")
}

func main() {
	flag.Parse()

	cmd := exec.Command("go", "generate", "./...")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	{
		cmd := exec.Command("go", "build", "./cmd/go")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		cmd = exec.Command("go", "build", "cmd/gofmt")
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
