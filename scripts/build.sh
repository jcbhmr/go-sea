#!/usr/bin/env bash
set -Eeuo pipefail
go generate ./...
CGO_ENABLED=0 go build ./cmd/go
CGO_ENABLED=0 go build cmd/gofmt
