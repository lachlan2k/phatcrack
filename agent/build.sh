#!/bin/bash

VERSION_STR=$(git describe --tags)

echo "Building version ${VERSION_STR}"

CGO_ENABLED=0 go build -ldflags="-X github.com/lachlan2k/phatcrack/agent/internal/version.version=${VERSION_STR}" -o phatcrack-agent main.go