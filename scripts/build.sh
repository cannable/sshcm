#!/usr/bin/bash

buildPath="./cmd/sshcm"
outPath="./build"

# Create the build output directory
mkdir -p $outPath

# Clean up old builds
if [ -d outPath ]; then
    rm -f "${outPath}/*"
fi

GOOS=linux GOARCH=amd64 go build -o "${outPath}/sshcm_linux_amd64" $buildPath
GOOS=linux GOARCH=arm64 go build -o "${outPath}/sshcm_linux_arm64" $buildPath
GOOS=windows GOARCH=amd64 go build -o "${outPath}/sshcm_windows_amd64.exe" $buildPath
