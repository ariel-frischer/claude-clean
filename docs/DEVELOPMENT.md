# DEVELOPMENT

Local development setup and workflow for claude-clean.

## Prerequisites

- Go 1.24+
- Git
- GoReleaser (for local releases)

## Quick Start

```bash
# Clone the repo
git clone https://github.com/ariel-frischer/claude-clean.git
cd claude-clean

# Install dependencies
make deps

# Build
make build

# Run tests
make test
```

## Makefile Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `make build` | `make b` | Build the `cclean` binary |
| `make install` | `make i` | Install to `~/.local/bin` with alias setup |
| `make run` | `make r` | Run with sample mock data |
| `make run-verbose` | `make rv` | Run with verbose output |
| `make test` | `make t` | Run tests with `-v` |
| `make clean` | `make c` | Remove built binaries |
| `make fmt` | `make f` | Format code with `gofmt` |
| `make vet` | `make v` | Run `go vet` |
| `make deps` | `make d` | Download and tidy dependencies |
| `make all` | `make a` | Format, vet, and build |
| `make snapshot` | `make s` | Build snapshot release locally |

## Project Structure

```
claude-clean/
├── main.go              # Entry point and core logic
├── Makefile             # Build automation
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── install.sh           # One-line installer script
├── .goreleaser.yaml     # GoReleaser config
├── bin/                 # Build output (gitignored)
├── mocks/               # Sample JSONL test data
├── scripts/             # Helper scripts
├── docs/                # Documentation
└── .github/
    └── workflows/
        ├── ci.yml       # CI pipeline
        └── release.yml  # Release pipeline
```

## Git Hooks

Install git hooks after cloning:

```bash
make dev-setup
```

### Available Hooks

| Hook | Purpose |
|------|---------|
| `pre-rebase` | Backs up `.dev/` to `.git/.dev-backup/` before rebase |
| `post-rewrite` | Restores `.dev/` from backup after rebase on `dev` branch |
| `post-merge` | Auto-cleans `.dev/` when merging to `main` |
| `pre-merge-commit` | Prevents merging `main` into `dev` (use rebase instead) |

### Branch Workflow

This project uses two main branches:

- **`main`** - Stable release branch (public, no `.dev/` files)
- **`dev`** - Development branch (has `.dev/` files with internal docs/scripts)

| Action | Allowed |
|--------|---------|
| Merge `dev` → `main` | ✅ Yes |
| Rebase `dev` onto `main` | ✅ Yes (preferred) |
| Merge `main` → `dev` | ❌ No (use rebase) |

#### Syncing dev with main

```bash
git checkout dev
git rebase main              # pre-rebase backs up .dev/, post-rewrite restores it
git push origin dev --force-with-lease
```

## Development Workflow

### 1. Make Changes

```bash
# Create a feature branch
git checkout -b feature/my-feature

# Make your changes...

# Format code
make fmt

# Run linter
make vet

# Run tests
make test
```

### 2. Test Locally

```bash
# Build and run with sample data
make run

# Or test with real Claude output
claude -p "test prompt" --verbose --output-format stream-json | ./bin/cclean
```

### 3. Test Different Output Styles

```bash
./bin/cclean -s default mocks/claude-stream-json-simple.jsonl
./bin/cclean -s compact mocks/claude-stream-json-simple.jsonl
./bin/cclean -s minimal mocks/claude-stream-json-simple.jsonl
./bin/cclean -s plain mocks/claude-stream-json-simple.jsonl
```

## CI Pipeline

The CI pipeline (`.github/workflows/ci.yml`) runs on pushes to `main`/`dev` and PRs:

1. **Install dependencies** - `go mod download`
2. **Verify dependencies** - `go mod verify`
3. **Format check** - Fails if code isn't formatted
4. **Vet** - Static analysis with `go vet`
5. **Build** - Compile the project
6. **Test** - Run tests with race detection and coverage

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/fatih/color` | Terminal color output |
| `golang.org/x/term` | Terminal size detection |

## Build Flags

The binary is built with ldflags for optimization:

```bash
-ldflags="-s -w"  # Strip debug info, reduce binary size
```

GoReleaser also embeds version info:

```bash
-X main.version={{.Version}}
-X main.commit={{.ShortCommit}}
-X main.date={{.Date}}
```

## Testing Mock Data

Sample JSONL files in `mocks/` directory can be used for testing:

```bash
# Run with mock data
./bin/cclean mocks/claude-stream-json-simple.jsonl

# Create new mock data
claude -p "your prompt" --verbose --output-format stream-json > mocks/new-test.jsonl
```
