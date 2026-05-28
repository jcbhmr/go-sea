package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.jcbhmr.com/crossexec"
)

func main() {
	flag.Parse()

	argv := []string{}
	envv := os.Environ()
	state := 0
	for _, a := range flag.Args() {
		if state == 0 {
			if strings.Contains(a, "=") {
				envv = append(envv, a)
			} else {
				state = 1
			}
		}
		if state == 1 {
			argv = append(argv, a)
		}
	}

	argv0, err := exec.LookPath(argv[0])
	if err != nil {
		log.Fatal(err)
	}
	crossexec.Exec(argv0, argv, envv)
}
