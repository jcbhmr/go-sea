#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

temp_dir=$(mktemp -d)
pushd "$temp_dir"

cleanup() {
    trap - SIGINT SIGTERM ERR EXIT
    cd /
    rm -rf "$temp_dir"
}
trap cleanup SIGINT SIGTERM ERR EXIT

url=$(curl -fsSLo /dev/null -w '%{url_effective}\n' https://github.com/tinygo-org/tinygo/releases/latest)
tag="${url##*/}"
version="${tag#v}"
curl -fsSLo tinygo.deb "https://github.com/tinygo-org/tinygo/releases/download/$tag/tinygo_${version}_amd64.deb"
sudo dpkg -i tinygo.deb
