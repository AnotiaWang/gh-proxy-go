#!/bin/bash

# 定义输出目录
OUTPUT_DIR="build"
mkdir -p ${OUTPUT_DIR}

# Linux amd64
echo "Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/gh-proxy-go-linux-amd64 main.go
echo "Output: ${OUTPUT_DIR}/gh-proxy-go-linux-amd64"

# Linux arm64
echo "Building for Linux arm64..."
GOOS=linux GOARCH=arm64 go build -o ${OUTPUT_DIR}/gh-proxy-go-linux-arm64 main.go
echo "Output: ${OUTPUT_DIR}/gh-proxy-go-linux-arm64"

# Linux ARMv7
echo "Building for Linux ARMv7..."
GOOS=linux GOARCH=arm GOARM=7 go build -o ${OUTPUT_DIR}/gh-proxy-go-linux-armv7 main.go
echo "Output: ${OUTPUT_DIR}/gh-proxy-go-linux-armv7"

# Windows amd64
echo "Building for Windows amd64..."
GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/gh-proxy-go-windows-amd64.exe main.go
echo "Output: ${OUTPUT_DIR}/gh-proxy-go-windows-amd64.exe"

# macOS amd64
echo "Building for macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/gh-proxy-go-macos-amd64 main.go
echo "Output: ${OUTPUT_DIR}/gh-proxy-go-macos-amd64"

# macOS arm64
echo "Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -o ${OUTPUT_DIR}/gh-proxy-go-macos-arm64 main.go
echo "Output: ${OUTPUT_DIR}/gh-proxy-go-macos-arm64"

echo "Build completed."
