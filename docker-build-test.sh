#!/usr/bin/env bash

# Local Docker build test script
set -e

IMAGE_NAME="singui:latest"

echo "Building Docker image: $IMAGE_NAME..."
docker build -t "$IMAGE_NAME" .

echo "Docker image built successfully! Running test container on port 2096..."
docker run --rm -p 2096:2096 "$IMAGE_NAME"
