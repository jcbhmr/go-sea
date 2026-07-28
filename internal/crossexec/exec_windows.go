package crossexec

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
)

func Exec(argv0 string, argv []string, envv []string) (err error) {
	cmd := &exec.Cmd{
		Path:   argv0,
		Args:   argv,
		Env:    envv,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	runtime.LockOSThread()
	signal.Ignore(os.Interrupt)
	err = cmd.Run()
	if err != nil {
		if _, ok := errors.AsType[*exec.ExitError](err); ok {
			// continue
		} else {
			log.Fatal(err)
		}
	}
	os.Exit(cmd.ProcessState.ExitCode())
	panic("exited")
}
