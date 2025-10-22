.PHONY: help build install run test clean fmt vet deps all

# Binary name
BINARY_NAME=claude-clean-output

# Default target
help:
	@echo "Claude Clean Output - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  install    - Install to GOPATH/bin"
	@echo "  run        - Run with sample mock data"
	@echo "  run-verbose- Run with verbose output on sample data"
	@echo "  test       - Run tests"
	@echo "  clean      - Remove built binaries"
	@echo "  fmt        - Format code with gofmt"
	@echo "  vet        - Run go vet"
	@echo "  deps       - Download dependencies"
	@echo "  all        - Format, vet, and build"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# Install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install

# Run with sample data
run: build
	@echo "Running $(BINARY_NAME) with sample data..."
	@if [ -f mocks/claude-stream-json-simple.log ]; then \
		./$(BINARY_NAME) mocks/claude-stream-json-simple.log; \
	else \
		echo "No sample data found. Create mocks/claude-stream-json-simple.log or pipe data to ./$(BINARY_NAME)"; \
	fi

# Run with verbose output
run-verbose: build
	@echo "Running $(BINARY_NAME) with verbose output..."
	@if [ -f mocks/claude-stream-json-simple.log ]; then \
		./$(BINARY_NAME) -v mocks/claude-stream-json-simple.log; \
	else \
		echo "No sample data found. Create mocks/claude-stream-json-simple.log or pipe data to ./$(BINARY_NAME)"; \
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
