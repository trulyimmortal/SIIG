#!/bin/sh
# SIIG Install Script (macOS / Linux)
# Usage: curl -fsSL https://get.siig.dev/install.sh | sh

set -e

REPO="yv3000/siig"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="siig"

# detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux) PLATFORM="linux" ;;
  darwin) PLATFORM="darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

ASSET_NAME="siig-${PLATFORM}-${ARCH}"

echo "→ Detecting system: ${PLATFORM}/${ARCH}"

# get latest release tag
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Error: Could not determine latest release."
  exit 1
fi

echo "→ Latest version: ${LATEST}"

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/${ASSET_NAME}"
CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${LATEST}/checksums.txt"

# create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "→ Downloading ${ASSET_NAME}..."
curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${BINARY_NAME}"

# verify checksum
echo "→ Verifying checksum..."
EXPECTED=$(curl -fsSL "$CHECKSUMS_URL" | grep "$ASSET_NAME" | awk '{print $1}')
if [ -n "$EXPECTED" ]; then
  if command -v sha256sum > /dev/null 2>&1; then
    ACTUAL=$(sha256sum "${TMP_DIR}/${BINARY_NAME}" | awk '{print $1}')
  else
    ACTUAL=$(shasum -a 256 "${TMP_DIR}/${BINARY_NAME}" | awk '{print $1}')
  fi

  if [ "$EXPECTED" != "$ACTUAL" ]; then
    echo "Error: Checksum mismatch!"
    echo "  Expected: $EXPECTED"
    echo "  Actual:   $ACTUAL"
    echo "Aborting installation."
    exit 1
  fi
  echo "→ Checksum verified ✓"
else
  echo "⚠ Warning: Could not verify checksum (checksums.txt not found)"
fi

# install binary
chmod +x "${TMP_DIR}/${BINARY_NAME}"

if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
  echo "→ Requires sudo to install to ${INSTALL_DIR}"
  sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo ""
echo "✓ SIIG ${LATEST} installed successfully!"
echo "  Location: ${INSTALL_DIR}/${BINARY_NAME}"
echo ""
echo "  Run 'siig' in any project directory to get started."
