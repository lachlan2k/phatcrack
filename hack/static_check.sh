#!/bin/bash

cd api

nilaway ./...
golangci-lint run ./...

cd ../agent

nilaway ./...
golangci-lint run ./...