#!/bin/bash

BUILD_DIR=$(dirname "$0")/build
mkdir -p "$BUILD_DIR"

go get -d -v ./...
go install -v ./...
go build -o "$BUILD_DIR"/backup
upx "$BUILD_DIR"/backup

"$BUILD_DIR"/backup -p 8088