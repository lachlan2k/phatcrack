#!/bin/bash

cd api

echo "API Checks"
nilaway ./...
golangci-lint run ./... --disable errcheck
gocritic check ./...

echo "Agent Checks"
cd ../agent

nilaway -exclude-pkgs 'github.com/lachlan2k/phatcrack/agent/internal/installer' ./...
golangci-lint run ./... --disable errcheck
gocritic check -disable exitAfterDefer ./...