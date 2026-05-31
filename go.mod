module go.jcbhmr.com.internal/go-sea

go 1.26.0

toolchain go1.26.3

tool (
	go.jcbhmr.com.internal/go-sea/internal/cmd/build
	go.jcbhmr.com.internal/go-sea/internal/cmd/makew
)

require (
	github.com/dschmidt/go-layerfs v0.2.0
	github.com/nlepage/go-tarfs v1.2.1
	go.jcbhmr.com/crossexec v1.1.1
)

require (
	github.com/mattn/goveralls v0.0.12 // indirect
	github.com/timpalpant/gzran v0.0.0-20201127163450-7b631e56f57b // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/tools v0.8.0 // indirect
)

replace (
	go.jcbhmr.com.internal/go-sea/internal/data => ./internal/data
	go.jcbhmr.com.internal/go-sea/internal/data/linux-amd64 => ./internal/data/linux-amd64
)

ignore ./go
