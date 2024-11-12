#!/bin/bash

# 변수 설정
PROJECT_NAME="lazy-mzstudy-go"
DIR_NAME="bin"
BUILD_PATH=${DIR_NAME}/${PROJECT_NAME}

# 빌드 디렉토리 생성
mkdir -p ${DIR_NAME}

# Windows 64-bit 바이너리 빌드
echo "Building Windows 64-bit binary..."
GOOS=windows GOARCH=amd64 go build -o ${BUILD_PATH}_windows_amd64.exe

# Linux 64-bit 바이너리 빌드
echo "Building Linux 64-bit binary..."
GOOS=linux GOARCH=amd64 go build -o ${BUILD_PATH}_linux_amd64

# macOS Silicon (ARM64) 바이너리 빌드
echo "Building macOS Silicon (ARM64) binary..."
GOOS=darwin GOARCH=arm64 go build -o ${BUILD_PATH}_macos_arm64

echo "All binaries built successfully!"
