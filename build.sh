#!/usr/bin/env bash

version=$(git describe --long --dirty --abbrev=10 --tags)
lf="-X github.com/intel/oneapi-cli/cmd.version=${version}"

GOOS=linux GOARCH=amd64 go build -ldflags "$lf" -o linux/bin/oneapi-cli
GOOS=windows GOARCH=amd64 go build -ldflags "$lf" -o win/bin/oneapi-cli.exe
GOOS=darwin GOARCH=amd64 go build -ldflags "$lf" -o osx/bin/oneapi-cli