#!/bin/sh

set -eux

name="$1"
version="$2"
target="$3"
base="${target##*${name}-}"

export CGO_ENABLED=0
export GOARCH="${base#*-}"
export GOOS="${base%-*}"

exec go build -v -a -tags netgo -o "$target" -ldflags "-X main.version=${version}" ./
