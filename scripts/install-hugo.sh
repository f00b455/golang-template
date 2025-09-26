#!/bin/bash

# Install Hugo binary for the project
# This script installs a specific version of Hugo

set -e

# Configuration - allow override via environment variable
HUGO_VERSION="${HUGO_VERSION:-0.140.2}"
HUGO_BASE_URL="${HUGO_BASE_URL:-https://github.com/gohugoio/hugo/releases/download}"
BIN_DIR="${BIN_DIR:-$(pwd)/bin}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case ${ARCH} in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: ${ARCH}"
        exit 1
        ;;
esac

# Set download URL based on OS
case ${OS} in
    linux)
        HUGO_FILE="hugo_extended_${HUGO_VERSION}_Linux-${ARCH}.tar.gz"
        ;;
    darwin)
        HUGO_FILE="hugo_extended_${HUGO_VERSION}_Darwin-universal.tar.gz"
        ;;
    *)
        echo "Unsupported OS: ${OS}"
        exit 1
        ;;
esac

DOWNLOAD_URL="${HUGO_BASE_URL}/v${HUGO_VERSION}/${HUGO_FILE}"

# Create bin directory if it doesn't exist
mkdir -p "${BIN_DIR}"

# Download and extract Hugo
echo "Installing Hugo v${HUGO_VERSION}..."
echo "Downloading from: ${DOWNLOAD_URL}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf ${TMP_DIR}' EXIT

# Download Hugo
if ! curl -L -o "${TMP_DIR}/${HUGO_FILE}" "${DOWNLOAD_URL}"; then
    echo "Failed to download Hugo from ${DOWNLOAD_URL}"
    exit 1
fi

# Extract Hugo
if ! tar -xzf "${TMP_DIR}/${HUGO_FILE}" -C "${TMP_DIR}"; then
    echo "Failed to extract Hugo archive"
    exit 1
fi

# Move Hugo binary to bin directory
if ! mv "${TMP_DIR}/hugo" "${BIN_DIR}/hugo"; then
    echo "Failed to move Hugo binary to ${BIN_DIR}"
    exit 1
fi

# Make it executable
if ! chmod +x "${BIN_DIR}/hugo"; then
    echo "Failed to make Hugo binary executable"
    exit 1
fi

# Verify installation
if "${BIN_DIR}/hugo" version > /dev/null 2>&1; then
    echo "Hugo successfully installed!"
    "${BIN_DIR}/hugo" version
else
    echo "Hugo installation failed!"
    exit 1
fi

echo "Hugo installed to: ${BIN_DIR}/hugo"
echo "Add ${BIN_DIR} to your PATH or use ./bin/hugo to run Hugo"