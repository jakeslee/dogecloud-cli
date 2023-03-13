#!/bin/sh

version=$(git describe --long --dirty --abbrev=6 --tags)
flags="-X github.com/jakeslee/dogecloud-cli/cmd.BuildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X github.com/jakeslee/dogecloud-cli/cmd.Version=$version"

export CGO_ENABLED=0

echo "building version: $version"

for os in "linux" "darwin" ; do
  for arch in "arm64" "amd64" ; do
    echo "building for $os/$arch"
    GOOS="$os" GOARCH="$arch" go build -ldflags "$flags" -o "./output/dogecloud-cli-$os-$arch" main.go
    chmod +x "./output/dogecloud-cli-$os-$arch"
  done
done
