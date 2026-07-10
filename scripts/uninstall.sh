#!/bin/sh
# SIIG Uninstall Script (macOS / Linux)
# Usage: curl -fsSL https://get.siig.dev/uninstall.sh | sh

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="siig"
BINARY_PATH="${INSTALL_DIR}/${BINARY_NAME}"

if [ ! -f "$BINARY_PATH" ]; then
  echo "SIIG is not installed at ${BINARY_PATH}"
  echo "Nothing to uninstall."
  exit 0
fi

echo "→ Removing SIIG from ${BINARY_PATH}..."

if [ -w "$INSTALL_DIR" ]; then
  rm -f "$BINARY_PATH"
else
  echo "→ Requires sudo to remove from ${INSTALL_DIR}"
  sudo rm -f "$BINARY_PATH"
fi

echo ""
echo "✓ SIIG uninstalled successfully."
