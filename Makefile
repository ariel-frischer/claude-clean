.PHONY: help build install run run-verbose test clean fmt vet deps all snapshot release
.PHONY: b i r rv t c f v d a s

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
	@echo "Release targets:"
	@echo "  snapshot (s)  - Build snapshot release locally (no publish)"
	@echo "  release       - Create a GitHub release (requires git tag, uses goreleaser)"

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

# Create a GitHub release (requires GITHUB_TOKEN and git tag)
release:
	goreleaser release --clean

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
