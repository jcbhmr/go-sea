package linux_amd64

import "embed"

//go:generate go run ./internal/cmd/generate
//go:embed all:bin all:pkg
var FS embed.FS
