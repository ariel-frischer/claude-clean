.PHONY: help build install setup-alias run run-verbose test clean fmt vet deps all

# Binary name
BINARY_NAME=claude-clean-output

# Default target
help:
	@echo "Claude Clean Output - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  install     - Install to ~/.local/bin and optionally setup 'cclean' alias"
	@echo "  setup-alias - Setup shell alias for 'cclean' command"
	@echo "  run         - Run with sample mock data"
	@echo "  run-verbose - Run with verbose output on sample data"
	@echo "  test        - Run tests"
	@echo "  clean       - Remove built binaries"
	@echo "  fmt         - Format code with gofmt"
	@echo "  vet         - Run go vet"
	@echo "  deps        - Download dependencies"
	@echo "  all         - Format, vet, and build"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# Install to ~/.local/bin with optional alias setup
install: build
	@echo "Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp $(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	@chmod +x ~/.local/bin/$(BINARY_NAME)
	@echo "✓ Binary installed to ~/.local/bin/$(BINARY_NAME)"
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH:"
	@echo "  export PATH=\"\$$HOME/.local/bin:\$$PATH\""
	@echo ""
	@echo "Would you like to set up the 'cclean' alias? (Recommended)"
	@echo "This will add a shell function to your config file."
	@echo ""
	@read -p "Setup alias? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(MAKE) setup-alias; \
	else \
		echo "Skipped alias setup. You can run 'make setup-alias' later."; \
		echo ""; \
		echo "Or manually add this to your ~/.bashrc or ~/.zshrc:"; \
		echo ""; \
		echo "cclean() {"; \
		echo "  claude-code -p \"\$$*\" --verbose --output-format stream-json | claude-clean-output"; \
		echo "}"; \
	fi

# Setup shell alias (can be run separately)
setup-alias:
	@echo "Setting up 'cclean' alias..."
	@if [ -n "$$ZSH_VERSION" ] || [ -f ~/.zshrc ]; then \
		echo "" >> ~/.zshrc; \
		echo "# Claude Clean Output alias" >> ~/.zshrc; \
		echo "cclean() {" >> ~/.zshrc; \
		echo "  claude-code -p \"\$$*\" --verbose --output-format stream-json | claude-clean-output" >> ~/.zshrc; \
		echo "}" >> ~/.zshrc; \
		echo "✓ Added alias to ~/.zshrc"; \
		echo "Run: source ~/.zshrc"; \
	elif [ -f ~/.bashrc ]; then \
		echo "" >> ~/.bashrc; \
		echo "# Claude Clean Output alias" >> ~/.bashrc; \
		echo "cclean() {" >> ~/.bashrc; \
		echo "  claude-code -p \"\$$*\" --verbose --output-format stream-json | claude-clean-output" >> ~/.bashrc; \
		echo "}" >> ~/.bashrc; \
		echo "✓ Added alias to ~/.bashrc"; \
		echo "Run: source ~/.bashrc"; \
	else \
		echo "Could not detect shell config file."; \
		echo "Please manually add the alias to your shell config."; \
	fi

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
