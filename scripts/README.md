# Release Scripts

Scripts for building and releasing claude-clean-output binaries.

## Usage

### Build binaries for all platforms

```bash
make build-release VERSION=v0.1.0
```

This creates a `release/` directory with binaries for:
- Linux (amd64, arm64)
- macOS (amd64/Intel, arm64/Apple Silicon)
- Windows (amd64)

Plus a `SHA256SUMS` file for verification.

### Create a GitHub release

```bash
# Basic release
make release VERSION=v0.1.0

# With custom release notes
make release VERSION=v0.1.0 NOTES="Initial release with cool features"
```

This will:
1. Build binaries for all platforms
2. Create a git tag (if it doesn't exist)
3. Push the tag to GitHub
4. Create a GitHub release with all binaries attached
5. Include the SHA256SUMS file

**Requirements:**
- [GitHub CLI](https://cli.github.com/) (`gh`) must be installed
- You must be authenticated (`gh auth login`)
- You must have push access to the repository

## Manual Usage

You can also run the scripts directly:

```bash
# Build binaries
./scripts/build-binaries.sh v0.1.0

# Create release (after building)
./scripts/create-release.sh v0.1.0 "Release notes here"
```

## What gets released

Each release includes:
- `claude-clean-linux-amd64` - Linux 64-bit
- `claude-clean-linux-arm64` - Linux ARM 64-bit
- `claude-clean-darwin-amd64` - macOS Intel
- `claude-clean-darwin-arm64` - macOS Apple Silicon
- `claude-clean-windows-amd64.exe` - Windows 64-bit
- `SHA256SUMS` - Checksums for verification

## Binary optimization

Binaries are built with `-ldflags="-s -w"` to strip debug information and reduce file size.
