#!/bin/bash

go build -ldflags="-X github.com/lachlan2k/phatcrack/api/internal/version.version=$(git describe --tags)" -o phatcrack-api main.go