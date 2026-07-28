# Go as a single executable application

🚚 `go` and `gofmt` as standalone binaries

<div align=center>
<table>
<tr><th>Before<th>After
<tr valign=top><td>

```
~/go
├── api/
│   └── ...
├── bin/
│   ├── go
│   └── gofmt
├── doc/
│   └── ...
├── lib/
│   └── ...
├── misc/
│   └── ...
├── pkg/
│   └── ...
├── src/
│   └── ...
├── test/
│   └── ...
├── codereview.cfg
├── CONTRIBUTING.md
├── go.env
├── LICENSE
├── PATENTS
├── README.md
├── SECURITY.md
└── VERSION
```

<td>

```
~/.local/bin
├── go
└── gofmt
```

</table>
</div>

😎 It's cool to be a single binary

## Installation

![sh](https://img.shields.io/badge/sh-4EAA25?style=for-the-badge&logo=GNU+Bash&logoColor=FFFFFF)
![PowerShell](https://img.shields.io/badge/PowerShell-4C84E8?style=for-the-badge)
![ZIP](https://img.shields.io/badge/ZIP-E05E00?style=for-the-badge)
![tar](https://img.shields.io/badge/tar-A42E2B?style=for-the-badge)

You might also be looking for [the official Go installation guide](https://go.dev/doc/install) instead. 😉

This module is not installable via `go get`/`go install`. If you could use `go get`/`go install` you'd already have Go installed! Instead, this project releases precompiled `.tar.gz` and `.zip` archives with the `go` and `gofmt` binaries in their root. These binaries are entirely standalone; you can run `go fmt` without `gofmt` being present.

1. Navigate to https://github.com/jcbhmr/go-sea/releases.
2. Choose the version you wish to install. That's probably the one marked "Latest".
3. Download the platform-specific `.tar.gz` or `.zip` archive for your platform.
4. If you're on Linux, macOS, WSL, or another \*nix-like operating system:
    1. Make sure `~/.local/bin` exists to unpack to.
    2. Make sure `~/.local/bin` is in your `PATH`.
    3. Unpack the archive you just downloaded to `~/.local/bin`.
5. If you're on Windows:
    1. Make sure `%LOCALAPPDATA%\Programs` exists to unpack to.
    2. Make sure `%LOCALAPPDATA%\Programs` is in your `PATH`.
    3. Unpack the archive you just downloaded to `%LOCALAPPDATA%\Programs`
  
<sup>TODO: Add install script? Contributions welcome! ❤️</sup>
  
## Usage

![Linux](https://img.shields.io/badge/Linux-222222?style=for-the-badge&logo=Linux&logoColor=FCC624)
![macOS](https://img.shields.io/badge/macOS-000000?style=for-the-badge&logo=macOS&logoColor=FFFFFF)
![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge)

Use it like you would the official `go` binary!

## Development

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)

```sh
go generate ./...
go build ./cmd/go ./cmd/gofmt
```

### How it works

These binaries contain `//go:embed`-ed data that is unpacked and cached when run.
