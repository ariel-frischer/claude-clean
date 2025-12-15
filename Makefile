.PHONY: help build install uninstall run run-verbose test clean fmt vet deps all snapshot release patch minor major dev-setup
.PHONY: b i u r rv t c f v d a s p

# Binary name and output directory
BINARY_NAME=cclean
BUILD_DIR=bin

# Default target
help:
	@echo "Claude Clean Output - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build   (b)   - Build the binary"
	@echo "  install (i)   - Install cclean to ~/.local/bin"
	@echo "  uninstall (u) - Uninstall cclean from the system"
	@echo "  run     (r)   - Run with sample mock data"
	@echo "  run-verbose (rv) - Run with verbose output on sample data"
	@echo "  test    (t)   - Run tests"
	@echo "  clean   (c)   - Remove built binaries"
	@echo "  fmt     (f)   - Format code with gofmt"
	@echo "  vet     (v)   - Run go vet"
	@echo "  deps    (d)   - Download dependencies"
	@echo "  all     (a)   - Format, vet, and build"
	@echo "  dev-setup     - Install git hooks for development"
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
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cclean

# Install to ~/.local/bin (build from source)
install: build
	@mkdir -p $(HOME)/.local/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@chmod +x $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to ~/.local/bin/"

# Uninstall cclean from the system
uninstall:
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) --uninstall; \
	elif command -v $(BINARY_NAME) > /dev/null 2>&1; then \
		$(BINARY_NAME) --uninstall; \
	else \
		echo "Removing cclean from known locations..."; \
		rm -f $(HOME)/.local/bin/$(BINARY_NAME); \
		if [ -f /usr/local/bin/$(BINARY_NAME) ]; then \
			sudo rm -f /usr/local/bin/$(BINARY_NAME); \
		fi; \
		echo "Done."; \
	fi

# Run with sample data
run: build
	@echo "Running $(BINARY_NAME) with sample data..."
	@if [ -f mocks/claude-stream-json-simple.jsonl ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) mocks/claude-stream-json-simple.jsonl; \
	else \
		echo "No sample data found. Create mocks/claude-stream-json-simple.jsonl or pipe data to ./$(BUILD_DIR)/$(BINARY_NAME)"; \
	fi

# Run with verbose output
run-verbose: build
	@echo "Running $(BINARY_NAME) with verbose output..."
	@if [ -f mocks/claude-stream-json-simple.jsonl ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) -V mocks/claude-stream-json-simple.jsonl; \
	else \
		echo "No sample data found. Create mocks/claude-stream-json-simple.jsonl or pipe data to ./$(BUILD_DIR)/$(BINARY_NAME)"; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean built binaries
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
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

# Install git hooks
dev-setup:
	@./scripts/setup-hooks.sh

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
u: uninstall
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
