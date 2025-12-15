# Release Process

This project supports releasing to both GitHub and GitLab with automatic platform detection.

## Overview

Releases are triggered by git tags matching the `v*` pattern (e.g., `v1.0.0`). The release system:

1. Detects your git remote (GitHub or GitLab)
2. Creates and pushes a git tag
3. Builds cross-platform binaries via GoReleaser or GitHub Actions
4. Publishes the release with artifacts and checksums

## GitHub Releases

### Automated (CI)

When a tag is pushed, GitHub Actions (`.github/workflows/release.yml`) automatically:

- Runs tests
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64/Apple Silicon)
  - Windows (amd64)
- Generates SHA256 checksums
- Creates a GitHub Release with auto-generated release notes

### Manual (Local)

Use the Makefile targets for local releases:

```bash
# Auto-detects GitHub from remote
make patch   # Bump v0.0.X
make minor   # Bump v0.X.0
make major   # Bump vX.0.0

# Or specify exact version
make release VERSION=v1.2.3

# Force GitHub platform
make patch PLATFORM=github
```

Local releases use GoReleaser (`.goreleaser.yaml`) which:
- Builds optimized binaries with `-s -w` ldflags
- Embeds version info via ldflags
- Creates tar.gz archives (zip for Windows)
- Generates checksums
- Publishes to GitHub Releases

## GitLab Releases

### Manual (Local)

```bash
# Auto-detects GitLab from remote
make patch

# Force GitLab platform
make patch PLATFORM=gitlab
```

The Makefile handles token management:
- GitHub: Uses `gh auth token` and unsets `GITLAB_TOKEN`
- GitLab: Uses `GITLAB_TOKEN` env var and unsets `GITHUB_TOKEN`

## Platform Detection

The Makefile auto-detects your platform from the git remote URL:

```makefile
REMOTE_URL := $(shell git remote get-url origin)
DETECTED_PLATFORM := $(shell echo $(REMOTE_URL) | grep -q github && echo github || ...)
```

Override with `PLATFORM=github` or `PLATFORM=gitlab` if needed.

## Release Artifacts

Each release includes:

| File | Description |
|------|-------------|
| `cclean-linux-amd64` | Linux x86_64 binary |
| `cclean-linux-arm64` | Linux ARM64 binary |
| `cclean-darwin-amd64` | macOS Intel binary |
| `cclean-darwin-arm64` | macOS Apple Silicon binary |
| `cclean-windows-amd64.exe` | Windows x86_64 binary |
| `SHA256SUMS` / `checksums.txt` | SHA256 checksums |

## Workflow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  make patch │ ──▶ │  git tag    │ ──▶ │  git push   │
└─────────────┘     │  -a vX.Y.Z  │     │  origin tag │
                    └─────────────┘     └─────────────┘
                                               │
                    ┌──────────────────────────┘
                    ▼
         ┌─────────────────────┐
         │  GitHub Actions OR  │
         │  GoReleaser (local) │
         └─────────────────────┘
                    │
                    ▼
         ┌─────────────────────┐
         │  Build binaries     │
         │  Generate checksums │
         │  Create release     │
         └─────────────────────┘
```

## Prerequisites

- **GitHub**: `gh` CLI authenticated (`gh auth login`)
- **GitLab**: `GITLAB_TOKEN` environment variable set
- **Local releases**: GoReleaser installed (`go install github.com/goreleaser/goreleaser@latest`)
