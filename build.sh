#!/usr/bin/env bash
export CGO_ENABLED=0

export CGO_ENABLED=0
version=$(git describe --long --dirty --abbrev=10 --tags)
lf="-X github.com/intel/oneapi-cli/cmd.version=${version}"


GOOS=linux GOARCH=amd64 go build -ldflags "$lf" -trimpath -mod=readonly -gcflags="all=-spectre=all -N -l" -asmflags="all=-spectre=all" -ldflags="all=-s -w" -o linux/bin/oneapi-cli
GOOS=windows GOARCH=amd64 go build -ldflags "$lf" -trimpath -mod=readonly -gcflags="all=-spectre=all -N -l" -asmflags="all=-spectre=all" -ldflags="all=-s -w"  -o win/bin/oneapi-cli.exe
GOOS=darwin GOARCH=amd64 go build -ldflags "$lf" -trimpath -mod=readonly -gcflags="all=-spectre=all -N -l" -asmflags="all=-spectre=all" -ldflags="all=-s -w" -o osx/bin/oneapi-cli