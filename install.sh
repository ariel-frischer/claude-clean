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

    # Build asset pattern to match
    if [ "$OS" = "windows" ]; then
        ASSET_PATTERN="${BINARY_NAME}_.*_${OS}_${ARCH}\.zip"
    else
        ASSET_PATTERN="${BINARY_NAME}_.*_${OS}_${ARCH}\.tar\.gz"
    fi

    # Get download URL from GitHub API
    RELEASE_JSON=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest")
    DOWNLOAD_URL=$(echo "$RELEASE_JSON" | grep -o "\"browser_download_url\": \"[^\"]*${OS}_${ARCH}[^\"]*\"" | grep -E "\.(tar\.gz|zip)\"" | head -1 | sed 's/"browser_download_url": "//' | sed 's/"$//')
    VERSION=$(echo "$RELEASE_JSON" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//')

    if [ -z "$DOWNLOAD_URL" ]; then
        error "Failed to find download URL for ${OS}/${ARCH}"
    fi

    ARCHIVE_FILE=$(basename "$DOWNLOAD_URL")

    # Determine install location
    INSTALL_DIR=$(get_install_dir)
    INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"

    info "Download URL: ${DOWNLOAD_URL}"
    info "Install path: ${INSTALL_PATH}"
    echo

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    # Download archive
    info "Downloading ${BINARY_NAME} v${VERSION}..."
    ARCHIVE_PATH="${TMP_DIR}/${ARCHIVE_FILE}"
    if ! download "$DOWNLOAD_URL" "$ARCHIVE_PATH"; then
        error "Failed to download ${BINARY_NAME}. Check if the release exists at ${GITHUB_URL}/releases"
    fi

    # Extract archive
    info "Extracting..."
    if [ "$OS" = "windows" ]; then
        unzip -q "$ARCHIVE_PATH" -d "$TMP_DIR"
    else
        tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"
    fi

    # Find and install binary
    EXTRACTED_BINARY="${TMP_DIR}/${BINARY_NAME}"
    if [ ! -f "$EXTRACTED_BINARY" ]; then
        error "Binary not found in archive"
    fi

    chmod +x "$EXTRACTED_BINARY"

    if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ ! -w /usr/local/bin ]; then
        info "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "$EXTRACTED_BINARY" "$INSTALL_PATH"
    else
        mv "$EXTRACTED_BINARY" "$INSTALL_PATH"
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
