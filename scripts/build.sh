#!/usr/bin/bash

buildPath="./cmd/sshcm"

GOOS=linux GOARCH=amd64 go build -o sshcm_linux_amd64 $buildPath
GOOS=linux GOARCH=arm64 go build -o sshcm_linux_arm64 $buildPath
GOOS=windows GOARCH=amd64 go build -o sshcm_windows_amd64 $buildPath