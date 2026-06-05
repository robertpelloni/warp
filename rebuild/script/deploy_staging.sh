#!/bin/bash
set -e
APP_NAME="warp"
VERSION=$(cat ../VERSION.md)
BUILD_DIR="build"
STAGING_DIR="staging"
echo "Preparing staging deployment for $APP_NAME version $VERSION..."
if [ ! -d "$BUILD_DIR" ]; then echo "Error: Build directory not found. Run build.sh first."; exit 1; fi
mkdir -p $STAGING_DIR
echo "Creating staging bundles..."
tar -czf $STAGING_DIR/${APP_NAME}-v${VERSION}-staging.tar.gz -C $BUILD_DIR .
echo "Staging deployment prepared in $STAGING_DIR/"
