#!/usr/bin/env bash

set -euo pipefail

# Docker Hub username
DOCKER_USER="gaiety"
IMAGE_NAME="cli-proxy-api-plus"
TAG="latest"
FULL_IMAGE_NAME="${DOCKER_USER}/${IMAGE_NAME}:${TAG}"

echo "Gathering build information..."
VERSION="dev-cline"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo 'none')"
BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

echo "Building Docker image: ${FULL_IMAGE_NAME}"
echo "  Version: ${VERSION}"
echo "  Commit:  ${COMMIT}"
echo "  Date:    ${BUILD_DATE}"
echo

docker build -t "${FULL_IMAGE_NAME}" \
  --build-arg VERSION="${VERSION}" \
  --build-arg COMMIT="${COMMIT}" \
  --build-arg BUILD_DATE="${BUILD_DATE}" \
  .

echo
echo "Build complete!"
echo "To push this image to Docker Hub, make sure you are logged in (docker login) and run:"
echo "docker push ${FULL_IMAGE_NAME}"
