package list

import (
	"encoding/json"
	"os/exec"
)

type Module struct {
	Path     string
	Versions []string
}

func Versions(package_ string) (m *Module, err error) {
	cmd := exec.Command("go", "list", "-json", "-m", "-versions", package_)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &m)
	return m, err
}
