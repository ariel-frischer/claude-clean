#!/bin/bash
set -e

BINARY_NAME="claude-clean"
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
print_success "Binary installed to $INSTALL_DIR/$BINARY_NAME"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_warning "$INSTALL_DIR is not in your PATH"
    echo "  Add this to your shell config:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
fi

# Ask about alias setup
echo ""
echo "Would you like to set up the 'cclean' alias? (Recommended)"
echo "This adds a shell function for quick Claude prompts with clean output."
echo ""
read -p "Setup alias? [y/N] " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Skipped alias setup. You can run this script again later."
    exit 0
fi

# Choose authentication method
echo ""
echo "Choose authentication method:"
echo "  [1] OAuth (use your Claude Pro/Team plan) - FREE, no API costs"
echo "  [2] API Key (pay-per-use) - charges to your Anthropic API account"
echo ""
read -p "Enter choice [1/2]: " -n 1 -r AUTH_CHOICE
echo ""

if [ "$AUTH_CHOICE" = "1" ]; then
    CLAUDE_CMD='ANTHROPIC_API_KEY="" claude -p "$*" --verbose --output-format stream-json | claude-clean'
    print_success "Using OAuth (sets ANTHROPIC_API_KEY=\"\" to force plan usage)"
else
    CLAUDE_CMD='claude -p "$*" --verbose --output-format stream-json | claude-clean'
    print_success "Using API Key (pay-per-use billing)"
fi

# The cclean function with --help handler
CCLEAN_FUNC='# Claude Clean alias
cclean() {
  if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    echo "cclean - Quick Claude prompts with clean terminal output"
    echo ""
    echo "Usage: cclean <prompt>"
    echo ""
    echo "This is a convenience wrapper that pipes Claude output through claude-clean."
    echo "All arguments are passed as a prompt to Claude - this command has no flags."
    echo ""
    echo "Examples:"
    echo "  cclean what is 2+2"
    echo "  cclean explain the difference between TCP and UDP"
    echo ""
    echo "Related commands:"
    echo "  claude-clean --help   Parser options (styles, verbose mode)"
    echo "  claude --help         Claude Code CLI options"
    return 0
  fi
  '"$CLAUDE_CMD"'
}'

# Detect shell config file
if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
    SHELL_RC="$HOME/.zshrc"
elif [ -f "$HOME/.bashrc" ]; then
    SHELL_RC="$HOME/.bashrc"
else
    print_error "Could not detect shell config file."
    echo "Please manually add this to your shell config:"
    echo ""
    echo "$CCLEAN_FUNC"
    exit 1
fi

# Check if alias already exists
if grep -q "cclean()" "$SHELL_RC" 2>/dev/null; then
    print_warning "Alias 'cclean' already exists in $SHELL_RC"
    read -p "Replace it? [y/N] " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Remove existing cclean function (handles multi-line)
        sed -i '/# Claude Clean alias/,/^}/d' "$SHELL_RC"
        echo "" >> "$SHELL_RC"
        echo "$CCLEAN_FUNC" >> "$SHELL_RC"
        print_success "Replaced alias in $SHELL_RC"
    else
        echo "Kept existing alias."
        exit 0
    fi
else
    echo "" >> "$SHELL_RC"
    echo "$CCLEAN_FUNC" >> "$SHELL_RC"
    print_success "Added alias to $SHELL_RC"
fi

echo ""
echo "Run this to activate:"
echo "  source $SHELL_RC"
echo ""
echo "Then try:"
echo "  cclean --help"
echo "  cclean what is 2+2"
