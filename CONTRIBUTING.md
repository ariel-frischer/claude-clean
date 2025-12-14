# Contributing to Claude Clean Output

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the project.

## Quick Start

The project includes a Makefile with common commands to make development easier:

```bash
make help    # Show all available commands
```

## Development Commands

### Building

```bash
make build   # Build the binary
```

This compiles the Go code and creates the `claude-clean-output` executable in the current directory.

### Running

```bash
make run           # Run with sample mock data
make run-verbose   # Run with verbose output (-v flag)
```

These commands build the binary and run it with sample data from `mocks/claude-stream-json-simple.log` if available.

You can also run the binary directly:

```bash
# From a file
./claude-clean-output mocks/your-test-file.log

# With verbose output
./claude-clean-output -v mocks/your-test-file.log

# From stdin
cat mocks/your-test-file.log | ./claude-clean-output

# Live with Claude Code
claude-code -p "your prompt" --verbose --output-format stream-json | ./claude-clean-output
```

### Testing

```bash
make test    # Run all tests
```

### Code Quality

```bash
make fmt     # Format code with gofmt
make vet     # Run go vet for suspicious code
make all     # Format, vet, and build (recommended before commits)
```

### Dependencies

```bash
make deps    # Download and tidy dependencies
```

### Cleaning

```bash
make clean   # Remove built binaries
```

## Project Structure

```
.
├── main.go              # Main parser implementation
├── types.go             # Data structures
├── go.mod               # Go module definition
├── Makefile             # Build and development commands
├── README.md            # User-facing documentation
├── CONTRIBUTING.md      # This file
├── mocks/               # Test data files
├── tests/               # Unit tests
└── scripts/
    ├── hooks/           # Git hooks (install with setup-hooks.sh)
    ├── setup-hooks.sh   # Install git hooks
    ├── build-binaries.sh
    └── install.sh
```

## Development Workflow

1. **Fork and Clone**: Fork the repository and clone your fork
2. **Setup Hooks**: Run `./scripts/setup-hooks.sh` to install git hooks
3. **Create a Branch**: Create a feature branch for your changes
4. **Make Changes**: Implement your feature or fix
5. **Format and Check**: Run `make all` to format, vet, and build
6. **Test**: Ensure `make test` passes (and add tests for new features)
7. **Commit**: Write clear, descriptive commit messages
8. **Push**: Push to your fork
9. **Pull Request**: Open a PR with a clear description of changes

## Code Style

- Follow standard Go conventions and idioms
- Run `go fmt` (or `make fmt`) before committing
- Run `go vet` (or `make vet`) to catch common issues
- Keep functions focused and well-named
- Add comments for complex logic

## Testing

While the project doesn't have unit tests yet, you can test manually:

1. Create test files in the `mocks/` directory with Claude Code output
2. Run `make run` or `make run-verbose` to see the parsed output
3. Verify the output is correctly formatted and colored

When adding new features, please include test data that exercises the new functionality.

## Adding New Message Types

If you need to handle new Claude Code message types:

1. Update the structs in `main.go` to include new fields
2. Add a display function (e.g., `displayNewMessageType()`)
3. Update the switch statement in `displayMessage()`
4. Add color definitions if needed
5. Create sample test data in `mocks/`

## Dependencies

The project uses minimal dependencies:

- `github.com/fatih/color` - Terminal color output

To add a new dependency:

```bash
go get github.com/example/package
make deps  # Tidy up
```

## Commit Guidelines

- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Reference issues if applicable (e.g., "Fix #123: Handle empty messages")
- Keep commits focused on a single change

## Getting Help

- Check the README.md for user documentation
- Review existing code for patterns and conventions
- Open an issue for questions or discussions
- Reach out to maintainers if you need guidance

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT).
