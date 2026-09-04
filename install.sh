#!/usr/bin/env bash
set -e

echo -e "\n========================================================"
echo -e "  🚀 Antigravity Account Switcher - Installer"
echo -e "========================================================\n"

if ! command -v go &> /dev/null; then
    echo "❌ Go toolchain not found in PATH!"
    echo "Please install Go 1.24+ from https://go.dev/dl/"
    exit 1
fi

echo "[1/3] Compiling antigravity-account-switcher..."
mkdir -p ~/.local/bin
CGO_ENABLED=0 go build -o ~/.local/bin/antigravity-account-switcher ./cmd/antigravity-account-switcher
chmod +x ~/.local/bin/antigravity-account-switcher
echo "  ✓ Installed to ~/.local/bin/antigravity-account-switcher"

echo "[2/3] Setting default port to 1831..."
~/.local/bin/antigravity-account-switcher config set port 1831 > /dev/null 2>&1 || true
echo "  ✓ Default port configured to 1831"

echo "[3/3] Creating desktop integration (optional)..."
~/.local/bin/antigravity-account-switcher install-desktop > /dev/null 2>&1 || true

echo -e "\n========================================================"
echo -e "  🎉 Setup Complete!"
echo -e "========================================================"
echo -e "\nNext Steps:"
echo -e "  1. Add your Google account:"
echo -e "     antigravity-account-switcher add-account\n"
echo -e "  2. Launch Antigravity 2.0 under supervision:"
echo -e "     antigravity-account-switcher launch\n"