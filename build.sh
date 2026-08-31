#!/usr/bin/env bash

#====================================================
# SingUI Multi-Platform Cross-Build Script
#====================================================

set -e

VERSION="1.0.0"
OUTPUT_DIR="build"

echo "=== Building SingUI Frontend Assets ==="
cd frontend
npm install
npm run build
cd ..

rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

PLATFORMS=(
    "linux/amd64/singbox-ui-linux-amd64"
    "linux/arm64/singbox-ui-linux-arm64"
    "linux/arm/singbox-ui-linux-armv7"
    "darwin/amd64/singbox-ui-darwin-amd64"
    "darwin/arm64/singbox-ui-darwin-arm64"
    "windows/amd64/singbox-ui-windows-amd64.exe"
)

echo "=== Compiling SingUI Standalone Binaries ==="
for PLATFORM in "${PLATFORMS[@]}"; do
    IFS="/" read -r OS ARCH OUTPUT <<< "$PLATFORM"
    echo "Building for $OS/$ARCH -> $OUTPUT_DIR/$OUTPUT"
    
    GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.Version=$VERSION'" -o "$OUTPUT_DIR/$OUTPUT" main.go
    
    # Tar/Zip release packaging
    cd "$OUTPUT_DIR"
    if [ "$OS" = "windows" ]; then
        zip "${OUTPUT%.exe}.zip" "$OUTPUT"
    else
        tar -czf "${OUTPUT}.tar.gz" "$OUTPUT"
    fi
    cd ..
done

echo "=== Build Complete! Packages available in $OUTPUT_DIR/ ==="
