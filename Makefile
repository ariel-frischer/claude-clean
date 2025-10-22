.PHONY: help build install setup-alias run run-verbose test clean fmt vet deps all build-release release

# Binary name
BINARY_NAME=claude-clean

# Default target
help:
	@echo "Claude Clean Output - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  install       - Install to ~/.local/bin and optionally setup 'cclean' alias"
	@echo "  setup-alias   - Setup shell alias for 'cclean' command"
	@echo "  run           - Run with sample mock data"
	@echo "  run-verbose   - Run with verbose output on sample data"
	@echo "  test          - Run tests"
	@echo "  clean         - Remove built binaries"
	@echo "  fmt           - Format code with gofmt"
	@echo "  vet           - Run go vet"
	@echo "  deps          - Download dependencies"
	@echo "  all           - Format, vet, and build"
	@echo ""
	@echo "Release targets:"
	@echo "  build-release - Build binaries for all platforms (Linux, macOS, Windows)"
	@echo "  release       - Create a GitHub release (requires version, e.g., make release VERSION=v0.1.0)"

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
		echo "Or manually add ONE of these to your ~/.bashrc or ~/.zshrc:"; \
		echo ""; \
		echo "Option 1 - OAuth (Claude Pro/Team plan, FREE):"; \
		echo "  # Setting ANTHROPIC_API_KEY=\"\" forces OAuth, ignoring any API key"; \
		echo "cclean() {"; \
		echo "  ANTHROPIC_API_KEY=\"\" claude -p \"\$$*\" --verbose --output-format stream-json | claude-clean"; \
		echo "}"; \
		echo ""; \
		echo "Option 2 - API Key (pay-per-use):"; \
		echo "  # Uses configured API key, you'll be billed per request"; \
		echo "cclean() {"; \
		echo "  claude -p \"\$$*\" --verbose --output-format stream-json | claude-clean"; \
		echo "}"; \
	fi

# Setup shell alias (can be run separately)
setup-alias:
	@echo "Setting up 'cclean' alias..."
	@echo ""
	@echo "Choose authentication method:"
	@echo "  [1] OAuth (use your Claude Pro/Team plan) - FREE, no API costs"
	@echo "      Sets ANTHROPIC_API_KEY=\"\" to force OAuth and use your plan"
	@echo ""
	@echo "  [2] API Key (pay-per-use) - charges to your Anthropic API account"
	@echo "      Uses configured API key, billed per request"
	@echo ""
	@read -p "Enter choice [1/2]: " -n 1 -r AUTH_CHOICE; \
	echo; \
	if [ "$$AUTH_CHOICE" = "1" ]; then \
		ALIAS_CMD='ANTHROPIC_API_KEY="" claude -p "$$*" --verbose --output-format stream-json | claude-clean'; \
		echo "✓ Using OAuth (sets ANTHROPIC_API_KEY=\"\" to force plan usage)"; \
	else \
		ALIAS_CMD='claude -p "$$*" --verbose --output-format stream-json | claude-clean'; \
		echo "✓ Using API Key (pay-per-use billing)"; \
	fi; \
	if [ -n "$$ZSH_VERSION" ] || [ -f ~/.zshrc ]; then \
		if grep -q "cclean()" ~/.zshrc 2>/dev/null; then \
			echo "✓ Alias 'cclean' already exists in ~/.zshrc - not adding again"; \
		else \
			echo "" >> ~/.zshrc; \
			echo "# Claude Clean alias" >> ~/.zshrc; \
			echo "cclean() {" >> ~/.zshrc; \
			echo "  $$ALIAS_CMD" >> ~/.zshrc; \
			echo "}" >> ~/.zshrc; \
			echo "✓ Added alias to ~/.zshrc"; \
			echo "Run: source ~/.zshrc"; \
		fi \
	elif [ -f ~/.bashrc ]; then \
		if grep -q "cclean()" ~/.bashrc 2>/dev/null; then \
			echo "✓ Alias 'cclean' already exists in ~/.bashrc - not adding again"; \
		else \
			echo "" >> ~/.bashrc; \
			echo "# Claude Clean alias" >> ~/.bashrc; \
			echo "cclean() {" >> ~/.bashrc; \
			echo "  $$ALIAS_CMD" >> ~/.bashrc; \
			echo "}" >> ~/.bashrc; \
			echo "✓ Added alias to ~/.bashrc"; \
			echo "Run: source ~/.bashrc"; \
		fi \
	else \
		echo "Could not detect shell config file."; \
		echo "Please manually add the alias to your shell config."; \
	fi

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
		./$(BINARY_NAME) -v mocks/claude-stream-json-simple.jsonl; \
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

# Build binaries for all platforms
build-release:
	@./scripts/build-binaries.sh $(VERSION)

# Create a GitHub release
release: build-release
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required"; \
		echo "Usage: make release VERSION=v0.1.0 [NOTES='Release notes']"; \
		exit 1; \
	fi
	@./scripts/create-release.sh $(VERSION) "$(NOTES)"
