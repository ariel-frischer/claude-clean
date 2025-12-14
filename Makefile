.PHONY: help build install run run-verbose test clean fmt vet deps all snapshot release patch minor major
.PHONY: b i r rv t c f v d a s p

# Binary name
BINARY_NAME=cclean

# Default target
help:
	@echo "Claude Clean Output - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build   (b)   - Build the binary"
	@echo "  install (i)   - Install cclean to ~/.local/bin"
	@echo "  run     (r)   - Run with sample mock data"
	@echo "  run-verbose (rv) - Run with verbose output on sample data"
	@echo "  test    (t)   - Run tests"
	@echo "  clean   (c)   - Remove built binaries"
	@echo "  fmt     (f)   - Format code with gofmt"
	@echo "  vet     (v)   - Run go vet"
	@echo "  deps    (d)   - Download dependencies"
	@echo "  all     (a)   - Format, vet, and build"
	@echo ""
	@echo "Release targets (auto-detects github/gitlab from remote, override with PLATFORM=github|gitlab):"
	@echo "  snapshot (s)  - Build snapshot release locally (no publish)"
	@echo "  patch   (p)   - Release patch bump (v0.0.X)"
	@echo "  minor         - Release minor bump (v0.X.0)"
	@echo "  major         - Release major bump (vX.0.0)"
	@echo "  release       - Release specific version (make release VERSION=v1.0.0)"
	@echo ""
	@echo "Examples:"
	@echo "  make p                    # auto-detect platform"
	@echo "  make p PLATFORM=github    # force GitHub release"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# Install to ~/.local/bin with optional alias setup
install:
	@./scripts/install.sh

# Run with sample data
run: build
	@echo "Running $(BINARY_NAME) with sample data..."
	@if [ -f mocks/claude-stream-json-simple.jsonl ]; then \
		./$(BINARY_NAME) mocks/claude-stream-json-simple.jsonl; \
	else \
		echo "No sample data found. Create mocks/claude-stream-json-simple.jsonl or pipe data to ./$(BINARY_NAME)"; \
	fi

# Run with verbose output
run-verbose: build
	@echo "Running $(BINARY_NAME) with verbose output..."
	@if [ -f mocks/claude-stream-json-simple.jsonl ]; then \
		./$(BINARY_NAME) -V mocks/claude-stream-json-simple.jsonl; \
	else \
		echo "No sample data found. Create mocks/claude-stream-json-simple.jsonl or pipe data to ./$(BINARY_NAME)"; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean built binaries
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	go clean

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build everything
all: fmt vet build

# Build snapshot release locally (no publish)
snapshot:
	goreleaser release --snapshot --clean

# Get current version from git tags (defaults to v0.0.0 if no tags)
CURRENT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
MAJOR := $(shell echo $(CURRENT_VERSION) | sed 's/v//' | cut -d. -f1)
MINOR := $(shell echo $(CURRENT_VERSION) | sed 's/v//' | cut -d. -f2)
PATCH := $(shell echo $(CURRENT_VERSION) | sed 's/v//' | cut -d. -f3)

# Platform detection (override with PLATFORM=github or PLATFORM=gitlab)
REMOTE_URL := $(shell git remote get-url origin 2>/dev/null)
DETECTED_PLATFORM := $(shell echo $(REMOTE_URL) | grep -q github && echo github || (echo $(REMOTE_URL) | grep -q gitlab && echo gitlab || echo unknown))
PLATFORM ?= $(DETECTED_PLATFORM)

# Create a release: make release VERSION=v1.0.0 [PLATFORM=github|gitlab]
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=v1.0.0"; \
		echo "  or use: make patch | make minor | make major"; \
		echo "  override platform: PLATFORM=github or PLATFORM=gitlab"; \
		exit 1; \
	fi
	@echo "Releasing $(VERSION) to $(PLATFORM)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
ifeq ($(PLATFORM),github)
	unset GITLAB_TOKEN && GITHUB_TOKEN=$$(gh auth token) goreleaser release --clean
else ifeq ($(PLATFORM),gitlab)
	unset GITHUB_TOKEN && goreleaser release --clean
else
	@echo "Error: Unknown platform '$(PLATFORM)'. Use PLATFORM=github or PLATFORM=gitlab"
	@exit 1
endif

# Auto-bump releases
patch:
	@$(MAKE) release VERSION=v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))

minor:
	@$(MAKE) release VERSION=v$(MAJOR).$(shell echo $$(($(MINOR)+1))).0

major:
	@$(MAKE) release VERSION=v$(shell echo $$(($(MAJOR)+1))).0.0

# Abbreviations
b: build
i: install
r: run
rv: run-verbose
t: test
c: clean
f: fmt
v: vet
d: deps
a: all
s: snapshot
p: patch
