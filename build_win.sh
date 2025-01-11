#!/bin/bash

BUILD_DATE="$(date +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git rev-parse HEAD)"
VERSION="$(git describe --tags --abbrev=0 | tr -d '\n')"
INFO="$VERSION
$BUILD_DATE
$GIT_COMMIT"

GOARCH=amd64 CGO_ENABLED=0 go build -o scms.exe -ldflags="-X 'main.VersionInfo=$INFO'"