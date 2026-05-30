package golang

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

func Make(dir string, goos string, goarch string, extraEnv []string) error {
	makeName := "./make.bash"
	if runtime.GOOS == "windows" {
		makeName = "./make.bat"
	}
	cmd := exec.Command(makeName)
	cmd.Dir = filepath.Join(dir, "src")
	cmd.Env = append(cmd.Environ(), "GOOS="+goos, "GOARCH="+goarch)
	cmd.Env = append(cmd.Env, extraEnv...)
	return cmd.Run()
}

func MakeDistpack(dir string, goos string, goarch string, extraEnv []string) error {
	makeName := "./make.bash"
	if runtime.GOOS == "windows" {
		makeName = "./make.bat"
	}
	cmd := exec.Command(makeName, "-distpack")
	cmd.Dir = filepath.Join(dir, "src")
	cmd.Env = append(cmd.Environ(), "GOOS="+goos, "GOARCH="+goarch)
	cmd.Env = append(cmd.Env, extraEnv...)
	return cmd.Run()
}
