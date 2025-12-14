# CONTRIBUTING

Guidelines for contributing to claude-clean.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/claude-clean.git
   cd claude-clean
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/ariel-frischer/claude-clean.git
   ```
4. Install dependencies:
   ```bash
   make deps
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

- Write clear, concise code
- Follow existing code style
- Add tests for new functionality

### 3. Test Your Changes

```bash
# Format code (required)
make fmt

# Run linter
make vet

# Run tests
make test

# Build and test locally
make build
./cclean mocks/claude-stream-json-simple.jsonl
```

### 4. Commit

Write clear commit messages:

```bash
git commit -m "feat: add support for new message type"
git commit -m "fix: handle empty tool results"
git commit -m "docs: update usage examples"
```

Commit message prefixes:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code refactoring
- `test:` - Tests
- `chore:` - Maintenance

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Code Style

- Run `make fmt` before committing
- Run `make vet` to catch issues
- Keep functions focused and small
- Add comments for non-obvious logic

## Testing

- Add test cases for new functionality
- Ensure existing tests pass: `make test`
- Test with real Claude output when possible

### Creating Test Data

```bash
# Generate test JSONL file
claude -p "test prompt" --verbose --output-format stream-json > mocks/new-test.jsonl
```

## Pull Request Guidelines

1. **Title**: Clear, descriptive title
2. **Description**: Explain what and why
3. **Tests**: Include tests for changes
4. **Docs**: Update docs if needed
5. **Single Focus**: One feature/fix per PR

## CI Requirements

PRs must pass:
- Format check (`gofmt`)
- Linter (`go vet`)
- Build
- Tests

## Reporting Issues

When reporting bugs:

1. **Search existing issues** first
2. **Include**:
   - Go version (`go version`)
   - OS and architecture
   - Steps to reproduce
   - Expected vs actual behavior
   - Sample JSONL input (if applicable)

## Feature Requests

- Open an issue with `[Feature]` prefix
- Describe the use case
- Explain why it's useful

## Questions?

Open a GitHub issue or discussion.
