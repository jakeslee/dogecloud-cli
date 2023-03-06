#!/bin/sh

version=$(git describe --long --dirty --abbrev=6 --tags)
flags="-X github.com/jakeslee/dogecloud-cli/cmd.BuildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X github.com/jakeslee/dogecloud-cli/cmd.Version=$version"

export CGO_ENABLED=0

echo "building version: $version"
go build -ldflags "$flags" -o ./output/dogecloud-cli main.go
