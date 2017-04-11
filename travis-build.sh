#!/usr/bin/env bash

export GOOS=linux
export GOARCH=amd64

echo "Building linux binary: rbacq with env variables:"
env | grep GO
go build -ldflags='-s -w' -v -o $TRAVIS_BUILD_DIR/rbacq

export GOOS=windows
export GOARCH=amd64

echo "Building windows binary: rbacq.exe with env variables:"
env | grep GO
go build -ldflags='-s -w' -v -o $TRAVIS_BUILD_DIR/rbacq.exe
