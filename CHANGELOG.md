# Changelog

All notable changes to claude-clean will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Watch mode for monitoring Claude log files in real-time
- Project constitution and worktree setup script

## [0.2.0] - 2025-12-20

### Added

- Multiple display styles (compact, default, minimal, plain)
- Message parsing with system reminder stripping
- Unit tests for parser and display functions

### Changed

- Restructured into public parser/ and display/ packages for library use
- Simplified README manual download instructions

### Fixed

- Release workflow build output paths
- Version ldflags case (Version not version)
- goreleaser configuration and CI workflow

### Removed

- Obsolete autospec scripts and project restructuring specs

## [0.1.4] - 2025-12-15

### Fixed

- Use browser_download_url from API instead of constructing URL
- Download tar.gz archives instead of standalone binaries

## [0.1.3] - 2025-12-15

### Added

- Uninstall target in Makefile
- Uninstall instructions in documentation

### Changed

- Simplified installation by integrating steps into Makefile

### Fixed

- goreleaser archives syntax and configuration

## [0.1.2] - 2025-12-14

### Added

- Compact, default, minimal, and plain display styles
- Display function tests

### Changed

- Renamed binary from claude-clean to cclean
- Updated all documentation references for new binary name

## [0.1.1] - 2025-12-14

### Added

- Git hooks (pre-merge-commit, pre-rebase, post-rewrite)
- CONTRIBUTING.md and development documentation

### Changed

- Formatting improvements for command-line flags

### Removed

- Outdated CONTRIBUTING.md and setup-hooks.sh

## [0.1.0] - 2025-12-14

### Added

- Initial release of claude-clean (cclean)
- JSONL stream parser for Claude Code stream-json output
- Multiple output styles (default, compact, minimal, plain)
- System reminder stripping from assistant messages
- Duplicate assistant message detection
- 10MB max buffer for large outputs
- Large line handling with graceful exit
- Line number display option
- Version flag
- goreleaser configuration for cross-platform builds
- One-line installer script
- OAuth support for Claude prompts

[Unreleased]: https://gitlab.com/ariel-frischer/claude-clean/-/compare/v0.2.0...HEAD
[0.2.0]: https://gitlab.com/ariel-frischer/claude-clean/-/compare/v0.1.4...v0.2.0
[0.1.4]: https://gitlab.com/ariel-frischer/claude-clean/-/compare/v0.1.3...v0.1.4
[0.1.3]: https://gitlab.com/ariel-frischer/claude-clean/-/compare/v0.1.2...v0.1.3
[0.1.2]: https://gitlab.com/ariel-frischer/claude-clean/-/compare/v0.1.1...v0.1.2
[0.1.1]: https://gitlab.com/ariel-frischer/claude-clean/-/compare/v0.1.0...v0.1.1
