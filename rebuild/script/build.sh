#!/bin/bash
set -e
APP_NAME="warp"
VERSION=$(cat ../VERSION.md)
BUILD_DIR="build"
echo "Building $APP_NAME version $VERSION..."
mkdir -p $BUILD_DIR
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/${APP_NAME}-linux-amd64 ./cmd/warp
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/${APP_NAME}-darwin-amd64 ./cmd/warp
echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o $BUILD_DIR/${APP_NAME}-darwin-arm64 ./cmd/warp
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/${APP_NAME}-windows-amd64.exe ./cmd/warp
echo "Build complete. Artifacts are in $BUILD_DIR/"
