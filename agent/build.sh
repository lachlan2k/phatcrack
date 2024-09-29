#!/bin/bash

VERSION_STR=$(git describe --tags)

echo "Building version ${VERSION_STR}"

GOOS=${1:-linux}
FILENAME="phatcrack-agent"
if [ "$GOOS" = "windows" ]; then
    FILENAME="phatcrack-agent.exe"
fi

GOOS=$GOOS CGO_ENABLED=0 go build -ldflags="-X github.com/lachlan2k/phatcrack/agent/internal/version.version=${VERSION_STR}" -o $FILENAME main.go