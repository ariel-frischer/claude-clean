# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Removed

- Removed automatic `claude` process spawning when no arguments provided; cclean now reads from stdin by default

## [1.0.0] - 2025-12-12

### Added

- Initial public release
- Parse Claude Code stream-json output into human-readable terminal format
- Color-coded output for different message types (system, assistant, tool use, tool result)
- Four output styles: default, compact, minimal, plain
- Verbose mode with tool IDs and token usage statistics
- Smart truncation for large tool inputs and outputs
- Special formatting for TodoWrite tool with status icons
- Support for reading from stdin or file input
- 10MB buffer for handling large JSON lines
- OAuth support for Claude prompts via `--oauth` flag

### Documentation

- README with installation and usage instructions
- CONTRIBUTING guide for contributors
- CLAUDE.md with architecture details for AI-assisted development
