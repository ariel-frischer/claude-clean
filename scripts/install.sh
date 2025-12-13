#!/bin/bash
set -e

# Read BINARY_NAME from Makefile (single source of truth)
BINARY_NAME=$(grep '^BINARY_NAME=' Makefile | cut -d'=' -f2)
INSTALL_DIR="$HOME/.local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_warning() { echo -e "${YELLOW}!${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go first:"
    echo "  Arch: sudo pacman -S go"
    echo "  macOS: brew install go"
    echo "  Ubuntu: sudo apt install golang-go"
    exit 1
fi

# Build the binary
echo "Building $BINARY_NAME..."
go build -o "$BINARY_NAME" .
print_success "Built $BINARY_NAME"

# Install to ~/.local/bin
echo ""
echo "Installing to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"
print_success "Installed to $INSTALL_DIR/$BINARY_NAME"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    print_warning "$INSTALL_DIR is not in your PATH"
    echo "  Add this to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "Usage:"
echo "  $BINARY_NAME \"your prompt\"           Run Claude with clean output"
echo "  $BINARY_NAME -oauth \"your prompt\"    Use OAuth instead of API key"
echo "  $BINARY_NAME log.jsonl               Parse existing JSON log"
echo "  cat log.jsonl | $BINARY_NAME         Parse from stdin"
echo ""
echo "Run '$BINARY_NAME -h' for all options."
