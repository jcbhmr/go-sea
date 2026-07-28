//go:build unix || plan9

package crossexec

import "syscall"

func Exec(argv0 string, argv []string, envv []string) (err error) {
	return syscall.Exec(argv0, argv, envv)
}
