#!/bin/bash

VERSION_STR=$(git describe --tags)

echo "Building version ${VERSION_STR}"

go build -buildvcs=false -ldflags="-X github.com/lachlan2k/phatcrack/api/internal/version.version=${VERSION_STR}" -o phatcrack-api main.go