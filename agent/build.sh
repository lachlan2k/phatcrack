#!/bin/bash

CGO_ENABLED=0 go build -ldflags="-X github.com/lachlan2k/phatcrack/agent/internal/version.version=$(git describe --tags)" -o phatcrack-agent main.go