# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Claude Clean (`cclean`) is a terminal parser that transforms Claude Code's raw stream-json output into beautiful, human-readable terminal output with colors and box-drawing characters. It supports multiple output styles (default/compact/minimal/plain) and can process live streams, files, or stdin.

## Development Commands

```bash
make build    # Compile to bin/cclean
make test     # Run tests with verbose output
make fmt      # Format Go code (required before commits)
make vet      # Static analysis
make all      # Format + vet + build (recommended before commits)
make run      # Build and run with sample data from mocks/
```

Single-letter shortcuts: `make b` (build), `make t` (test), `make f` (fmt), `make a` (all)

Run a single test:
```bash
go test -v -run TestFunctionName ./...
```

## Architecture

The codebase is organized into a pipeline architecture in a single-package structure:

**main.go** (~1500 lines) - Core processing pipeline:
1. **Input routing**: File, stdin, or live Claude execution via `runClaude()`
2. **JSONL parsing**: Line-by-line JSON deserialization into `StreamMessage` structs
3. **Duplicate detection**: Buffers assistant messages to skip duplicates in result messages
4. **Display routing**: `displayMessage()` routes to style-specific formatters

**types.go** - JSON schema definitions matching Claude Code's stream-json format:
- `StreamMessage` - Top-level wrapper (type: system/assistant/user/result)
- `MessageContent` - Content container with `[]ContentBlock`
- `ContentBlock` - Individual content pieces (text/tool_use/tool_result)

**Display functions** (48 total, in main.go) - Each message type has 4 style variants:
- `displaySystemMessage()`, `displaySystemMessageCompact()`, `displaySystemMessageMinimal()`, `displaySystemMessagePlain()`
- Same pattern for Assistant, User/ToolResult, and Result messages

## Key Implementation Details

- **10MB max buffer** for handling large outputs (line ~168)
- **System reminder stripping**: Regex removes `<system-reminder>` tags in non-verbose mode
- **Large output truncation**: Shows first 20 + last 20 lines with middle summary
- **Todo formatting**: Special visual indicators for TodoWrite tool (✓/→/○)

## Testing

Test files are organized by output style:
- `main_test.go` - Core function tests (stripSystemReminders, etc.)
- `display_test.go` - Default style tests
- `display_minimal_test.go` - Minimal style tests
- `display_plain_test.go` - Plain style tests

Sample test data: `mocks/claude-stream-json-simple.jsonl`
