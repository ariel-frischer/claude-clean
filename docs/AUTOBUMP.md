# Autobump: Makefile Version Management

The Makefile provides automatic semantic version bumping for releases.

## Quick Reference

| Command | Action | Example |
|---------|--------|---------|
| `make patch` or `make p` | Bump patch version | v1.2.3 → v1.2.4 |
| `make minor` | Bump minor version | v1.2.3 → v1.3.0 |
| `make major` | Bump major version | v1.2.3 → v2.0.0 |
| `make release VERSION=vX.Y.Z` | Set exact version | → vX.Y.Z |
| `make snapshot` or `make s` | Local build (no publish) | - |

## How It Works

### Version Detection

The current version is extracted from the latest git tag:

```makefile
CURRENT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
```

If no tags exist, defaults to `v0.0.0`.

### Version Parsing

The version is split into semantic components:

```makefile
MAJOR := $(shell echo $(CURRENT_VERSION) | sed 's/v//' | cut -d. -f1)
MINOR := $(shell echo $(CURRENT_VERSION) | sed 's/v//' | cut -d. -f2)
PATCH := $(shell echo $(CURRENT_VERSION) | sed 's/v//' | cut -d. -f3)
```

### Auto-Increment Logic

Each bump target calculates the new version using shell arithmetic:

```makefile
# Patch: v1.2.3 → v1.2.4
patch:
    @$(MAKE) release VERSION=v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))

# Minor: v1.2.3 → v1.3.0
minor:
    @$(MAKE) release VERSION=v$(MAJOR).$(shell echo $$(($(MINOR)+1))).0

# Major: v1.2.3 → v2.0.0
major:
    @$(MAKE) release VERSION=v$(shell echo $$(($(MAJOR)+1))).0.0
```

## Usage Examples

```bash
# Check current version
git describe --tags --abbrev=0
# v0.2.1

# Patch release (most common)
make patch
# Creates v0.2.2

# Minor release (new features)
make minor
# Creates v0.3.0

# Major release (breaking changes)
make major
# Creates v1.0.0

# Specific version
make release VERSION=v1.0.0-beta.1
```

## Platform Override

Combine with platform selection:

```bash
# Force GitHub
make patch PLATFORM=github

# Force GitLab
make minor PLATFORM=gitlab
```

## Snapshot Builds

Test the release process without publishing:

```bash
make snapshot
# or
make s
```

This runs GoReleaser with `--snapshot --clean` to build all artifacts locally in `dist/`.

## What Happens During a Release

1. **Version Calculation** - New version computed from current tag
2. **Tag Creation** - Annotated tag created: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
3. **Tag Push** - Tag pushed to origin: `git push origin vX.Y.Z`
4. **Build & Publish** - GoReleaser builds and publishes (or CI takes over)

## Abbreviations

The Makefile supports short aliases:

```bash
make p   # patch
make s   # snapshot
```

## Checking Version Before Release

```bash
# See current version
git describe --tags --abbrev=0

# Preview what patch would create
make -n patch  # dry-run
```

## Troubleshooting

**"v0.0.0" showing as current version**
- No git tags exist yet
- Run `make release VERSION=v0.1.0` to create first tag

**Version not incrementing correctly**
- Ensure tags are fetched: `git fetch --tags`
- Check tag format: must be `vX.Y.Z`
