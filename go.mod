module go.jcbhmr.com.internal/go-sea

go 1.26.0

toolchain go1.26.3

require (
	github.com/dschmidt/go-layerfs v0.2.0 // indirect
	github.com/mattn/goveralls v0.0.12 // indirect
	go.jcbhmr.com/crossexec v1.1.1 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/tools v0.8.0 // indirect
)

replace (
	go.jcbhmr.com.internal/go-sea/internal/data => ./internal/data
	go.jcbhmr.com.internal/go-sea/internal/data/linux-amd64 => ./internal/data/linux-amd64
)
