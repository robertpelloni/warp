#!/bin/bash
set -e

APP_NAME="warp"
VERSION=$(cat ../VERSION.md)
BUILD_DIR="build"
STAGING_DIR="staging"

echo "Preparing staging deployment for $APP_NAME version $VERSION..."

if [ ! -d "$BUILD_DIR" ]; then
    echo "Error: Build directory not found. Run build.sh first."
    exit 1
fi

mkdir -p $STAGING_DIR

# Verification
REQUIRED_BINARIES=(
    "warp-linux-amd64"
    "warp-darwin-amd64"
    "warp-darwin-arm64"
    "warp-windows-amd64.exe"
)

for bin in "${REQUIRED_BINARIES[@]}"; do
    if [ ! -f "$BUILD_DIR/$bin" ]; then
        echo "Error: Required binary $bin missing from build directory."
        exit 1
    fi
done

# Create deployment bundles
echo "Creating staging bundles..."
tar -czf $STAGING_DIR/${APP_NAME}-v${VERSION}-staging.tar.gz -C $BUILD_DIR .

echo "Staging deployment prepared in $STAGING_DIR/"
ls -l $STAGING_DIR/
