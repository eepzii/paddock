#!/bin/sh
set -e

GITHUB_USER="eepzii"
GITHUB_REPO="paddock"
BINARY_NAME="paddock"


OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)  OS_TYPE="Linux" ;;
    Darwin) OS_TYPE="Darwin" ;;
    *)      echo "OS $OS not supported"; exit 1 ;;
esac

case "$ARCH" in
    x86_64) ARCH_TYPE="x86_64" ;;
    aarch64|arm64) ARCH_TYPE="arm64" ;;
    *)      echo "Architecture $ARCH not supported"; exit 1 ;;
esac

LATEST_URL="https://api.github.com/repos/$GITHUB_USER/$GITHUB_REPO/releases/latest"
VERSION=$(curl -s $LATEST_URL | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Error: Could not find latest release version. (Have you released it yet?)"
    exit 1
fi

ASSET_NAME="${BINARY_NAME}_${OS_TYPE}_${ARCH_TYPE}.tar.gz"
DOWNLOAD_URL="https://github.com/$GITHUB_USER/$GITHUB_REPO/releases/download/$VERSION/$ASSET_NAME"

TEMP_DIR=$(mktemp -d)
echo "Downloading $DOWNLOAD_URL..."

if ! curl -sL "$DOWNLOAD_URL" -o "$TEMP_DIR/$ASSET_NAME"; then
    echo "Error: Download failed. Check if the release asset exists."
    exit 1
fi

tar -xzf "$TEMP_DIR/$ASSET_NAME" -C "$TEMP_DIR"

INSTALL_DIR="/usr/local/bin"

echo "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Requires sudo permissions to write to $INSTALL_DIR"
    sudo mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

rm -rf "$TEMP_DIR"

echo "Success! Run '$BINARY_NAME --help' to get started."