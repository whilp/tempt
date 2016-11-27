#!/bin/sh

name="$1"
base="${2##*/${name}-}"

export CGO_ENABLED=0
export GOARCH="${base#*-}"
export GOOS="${base%-*}"

exec go build -v -a -tags netgo -o "$2" -ldflags "-X main.version=${version}" ./
