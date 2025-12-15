# Unified Install Script

The `install.sh` script provides a single-command installation that works across all supported platforms.

## One-Line Install

```bash
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/claude-clean/main/install.sh | sh
```

Or with wget:
```bash
wget -qO- https://raw.githubusercontent.com/ariel-frischer/claude-clean/main/install.sh | sh
```

## How It Works

### Platform Detection

The script automatically detects OS and architecture:

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Detect OS      │ ──▶ │  Detect Arch    │ ──▶ │  Build URL      │
│  (uname -s)     │     │  (uname -m)     │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

#### OS Detection (`detect_os`)

| `uname -s` | Mapped OS |
|------------|-----------|
| `Linux*` | `linux` |
| `Darwin*` | `darwin` |
| `CYGWIN*`, `MINGW*`, `MSYS*` | `windows` |

#### Architecture Detection (`detect_arch`)

| `uname -m` | Mapped Arch |
|------------|-------------|
| `x86_64`, `amd64` | `amd64` |
| `aarch64`, `arm64` | `arm64` |
| `armv7l` | `arm` |

### Download URL Construction

The binary name is constructed dynamically:

```bash
# Linux/macOS
BINARY_FILE="cclean-${OS}-${ARCH}"
# Example: cclean-linux-amd64

# Windows
BINARY_FILE="cclean-${OS}-${ARCH}.exe"
# Example: cclean-windows-amd64.exe
```

Final URL:
```
https://github.com/ariel-frischer/claude-clean/releases/latest/download/${BINARY_FILE}
```

### Install Location Selection

The script chooses the install directory in order of preference:

1. `/usr/local/bin` - if writable (system-wide)
2. `~/.local/bin` - user directory (created if needed)

```bash
get_install_dir() {
    if [ -w /usr/local/bin ]; then
        echo "/usr/local/bin"
    elif [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin"; then
        echo "$HOME/.local/bin"
    else
        error "Cannot find writable install directory"
    fi
}
```

### Download Method

Supports both `curl` and `wget`:

```bash
get_downloader() {
    if command -v curl > /dev/null; then
        echo "curl -fsSL"
    elif command -v wget > /dev/null; then
        echo "wget -qO-"
    fi
}
```

## Features

### Cross-Platform Support

| Platform | Architectures |
|----------|---------------|
| Linux | amd64, arm64, arm |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Windows | amd64 |

### Colored Output

When running in a terminal (TTY), the script uses colored output:

- Cyan: Info messages
- Green: Success messages
- Yellow: Warnings
- Red: Errors

Colors are disabled when piped or in non-interactive contexts.

### PATH Warning

If the install directory isn't in `$PATH`, the script warns and provides the fix:

```bash
Note: /home/user/.local/bin is not in your PATH.
Add this to your shell config:

  export PATH="$PATH:/home/user/.local/bin"
```

### Version Verification

After installation, the script runs `--version` to confirm success:

```bash
Successfully installed cclean to /home/user/.local/bin/cclean
Version: cclean v0.2.1
```

### Usage After Install

```bash
# Pipe from Claude Code
claude -p "your prompt" --verbose --output-format stream-json | cclean

# Or read from file
cclean logfile.jsonl
```

## Security Considerations

- Uses `set -e` for fail-fast behavior
- Downloads to a temp file first
- Uses `trap` to clean up temp file on exit
- Verifies binary is executable after install
- Supports sudo escalation only when needed

## Troubleshooting

**"Unsupported operating system"**
- Running on an unsupported OS
- Check `uname -s` output

**"Unsupported architecture"**
- Running on an unsupported CPU architecture
- Check `uname -m` output

**"Neither curl nor wget found"**
- Install curl: `apt install curl` or `brew install curl`

**"Failed to download"**
- Check internet connection
- Verify release exists at GitHub releases page

**"Cannot find writable install directory"**
- Create `~/.local/bin`: `mkdir -p ~/.local/bin`
- Or run with sudo for `/usr/local/bin`
