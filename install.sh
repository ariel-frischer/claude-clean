#!/bin/sh
# Claude Clean Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/ariel-frischer/claude-clean/main/install.sh | sh

set -e

REPO="ariel-frischer/claude-clean"
BINARY_NAME="cclean"
GITHUB_URL="https://github.com/${REPO}"

# Colors (disabled if not a tty)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    CYAN='\033[0;36m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    CYAN=''
    NC=''
fi

info() { printf "${CYAN}%s${NC}\n" "$1"; }
success() { printf "${GREEN}%s${NC}\n" "$1"; }
warn() { printf "${YELLOW}%s${NC}\n" "$1"; }
error() { printf "${RED}%s${NC}\n" "$1" >&2; exit 1; }

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        CYGWIN*|MINGW*|MSYS*) echo "windows" ;;
        *) error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        armv7l) echo "arm" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get download command (curl or wget)
get_downloader() {
    if command -v curl > /dev/null 2>&1; then
        echo "curl -fsSL"
    elif command -v wget > /dev/null 2>&1; then
        echo "wget -qO-"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download file to path
download() {
    url="$1"
    dest="$2"
    if command -v curl > /dev/null 2>&1; then
        curl -fsSL -o "$dest" "$url"
    else
        wget -qO "$dest" "$url"
    fi
}

# Get install directory
get_install_dir() {
    # Try /usr/local/bin first (requires sudo), then ~/.local/bin
    if [ -w /usr/local/bin ]; then
        echo "/usr/local/bin"
    elif [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        echo "$HOME/.local/bin"
    else
        error "Cannot find writable install directory. Please create ~/.local/bin or run with sudo."
    fi
}

# Check if directory is in PATH
check_path() {
    dir="$1"
    case ":$PATH:" in
        *":$dir:"*) return 0 ;;
        *) return 1 ;;
    esac
}

main() {
    info "Installing ${BINARY_NAME}..."
    echo

    OS=$(detect_os)
    ARCH=$(detect_arch)

    info "Detected: ${OS}/${ARCH}"

    # Build download URL
    if [ "$OS" = "windows" ]; then
        BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}.exe"
    else
        BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}"
    fi

    DOWNLOAD_URL="${GITHUB_URL}/releases/latest/download/${BINARY_FILE}"

    # Determine install location
    INSTALL_DIR=$(get_install_dir)
    INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"

    info "Download URL: ${DOWNLOAD_URL}"
    info "Install path: ${INSTALL_PATH}"
    echo

    # Create temp file
    TMP_FILE=$(mktemp)
    trap 'rm -f "$TMP_FILE"' EXIT

    # Download binary
    info "Downloading ${BINARY_NAME}..."
    if ! download "$DOWNLOAD_URL" "$TMP_FILE"; then
        error "Failed to download ${BINARY_NAME}. Check if the release exists at ${GITHUB_URL}/releases"
    fi

    # Make executable and move to install location
    chmod +x "$TMP_FILE"

    if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ ! -w /usr/local/bin ]; then
        info "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "$TMP_FILE" "$INSTALL_PATH"
    else
        mv "$TMP_FILE" "$INSTALL_PATH"
    fi

    # Verify installation
    if [ -x "$INSTALL_PATH" ]; then
        success "Successfully installed ${BINARY_NAME} to ${INSTALL_PATH}"

        # Check version
        if VERSION=$("$INSTALL_PATH" --version 2>/dev/null); then
            info "Version: ${VERSION}"
        fi
    else
        error "Installation failed"
    fi

    # Warn if not in PATH
    if ! check_path "$INSTALL_DIR"; then
        echo
        warn "Note: ${INSTALL_DIR} is not in your PATH."
        warn "Add this to your shell config (~/.bashrc, ~/.zshrc, etc.):"
        echo
        echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
        echo
    fi

    echo
    success "Installation complete!"
    echo
    info "Usage: cclean [options] [file.jsonl]"
    echo
    info "Examples:"
    echo "  claude -p \"your prompt\" --verbose --output-format stream-json | cclean"
    echo "  cclean log.jsonl"
    echo
    info "For more info: ${GITHUB_URL}"
}

main "$@"
