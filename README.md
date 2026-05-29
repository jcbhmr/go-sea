# Go as a Single Executable Application

<table>
<tr align=center><th>Before<th>After
<tr align=center valign=top><td align=left>

```
/usr/local/go
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

<td align=left>

```
~/.local/bin
├── go
└── gofmt
```

</table>

## Installation

- **Windows x86-64:** [`go.windows-amd64.zip`](#TODO)
- **Linux x86-64:** [`go.linux-amd64.zip`](#TODO)
- **Linux AArch64:** [`go.linux-arm64.zip`](#TODO)
- **macOS AArch64:** [`go.darwin-arm64.zip`](#TODO)

<dl>
<div>
<dt>Linux, macOS, and other *nix-like systems
<dd>

```sh
unzip ./go.*.zip -d ~/.local/bin
```

<details><summary>Add <code>~/.local/bin</code> to your <code>PATH</code></summary>

```sh
echo 'PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
```

</details>

</div>
<div>
<dt>Windows
<dd>

```powershell
Expand-Archive -Path ./go.*.zip -DestinationPath "$Env:LocalAppData/Programs"
```

<details><summary>Add <code>%LocalAppData%</code> to your <code>PATH</code></summary>

```powershell
$binDir = "$env:LocalAppData\Programs"
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
$currentPathList = $currentPath -split ";"
if ($currentPathList -notcontains $binDir) {
    $newPath = if ([string]::IsNullOrWhiteSpace($currentPath)) {
        $binDir
    } else {
        "$currentPath;$binDir"
    }
    [System.Environment]::SetEnvironmentVariable("Path", $newPath, [System.EnvironmentVariableTarget]::User)
    Write-Host "$binDir added successfully to User PATH."
} else {
    Write-Host "$binDir already exists in User PATH."
}
```

</details>

</div>
</dl>

## Usage

Use it like you would the official `go` binary!

## Development

```sh
go build ./cmd/go
go build cmd/gofmt
```

### How it works

These binaries contain `//go:embed`-ed data that is unpacked and cached when run.
