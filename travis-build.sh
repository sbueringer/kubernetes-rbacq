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

echo "Downloading upx"
cd $TRAVIS_BUILD_DIR
curl -L -O https://github.com/upx/upx/releases/download/v3.93/upx-3.93-amd64_linux.tar.xz
tar xvf upx-3.93-amd64_linux.tar.xz

echo "Using upx on rbacq"
upx-3.93-amd64_linux/upx $TRAVIS_BUILD_DIR/rbacq

echo "Using upx on rbacq.exe"
upx-3.93-amd64_linux/upx $TRAVIS_BUILD_DIR/rbacq.exe
