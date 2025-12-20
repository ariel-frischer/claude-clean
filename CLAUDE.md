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

The codebase follows standard Go project layout with public packages:

```
claude-clean/
├── cmd/cclean/           # Application entry point
│   └── main.go           # CLI parsing, input routing
├── parser/               # Public: JSONL parsing and types
│   ├── parser.go         # StripSystemReminders, constants
│   ├── types.go          # StreamMessage, ContentBlock, Usage
│   └── parser_test.go
├── display/              # Public: Output formatting
│   ├── display.go        # Common utilities, color definitions, Config
│   ├── default.go        # Default style (box-drawing)
│   ├── compact.go        # Single-line summaries
│   ├── minimal.go        # No box-drawing
│   ├── plain.go          # No colors
│   └── display_test.go
├── mocks/                # Test data
└── bin/                  # Build output
```

## Using as a Library

Other Go modules can import the parser and display packages:

```go
import (
    "github.com/ariel-frischer/claude-clean/parser"
    "github.com/ariel-frischer/claude-clean/display"
)

// Use types
msg := parser.StreamMessage{...}
cfg := &display.Config{Style: display.StyleDefault}
display.DisplayMessage(&msg, 1, cfg)
```

**cmd/cclean/main.go** - Application entry point:
1. **CLI parsing**: Flag handling for style, verbose, line numbers
2. **Input routing**: File or stdin (defaults to stdin when no args)
3. **JSONL parsing**: Line-by-line JSON deserialization into `StreamMessage` structs
4. **Duplicate detection**: Buffers assistant messages to skip duplicates in result messages

**internal/parser/** - Data types and parsing utilities:
- `types.go`: JSON schema definitions matching Claude Code's stream-json format
  - `StreamMessage` - Top-level wrapper (type: system/assistant/user/result)
  - `MessageContent` - Content container with `[]ContentBlock`
  - `ContentBlock` - Individual content pieces (text/tool_use/tool_result)
- `parser.go`: Utility functions like `StripSystemReminders`

**internal/display/** - Output formatting:
- `display.go`: Common utilities, color definitions, `Config` struct, `DisplayMessage()` router
- Style-specific files: Each message type has 4 style variants across files

## Key Implementation Details

- **10MB max buffer** for handling large outputs (`parser.MaxBufferCapacity`)
- **System reminder stripping**: Regex removes `<system-reminder>` tags in non-verbose mode
- **Large output truncation**: Shows first 20 + last 20 lines with middle summary
- **Todo formatting**: Special visual indicators for TodoWrite tool (✓/→/○)

## Testing

Tests are organized by package:
- `internal/parser/parser_test.go` - Parser tests (stripSystemReminders, etc.)
- `internal/display/display_test.go` - Display function tests

Sample test data: `mocks/claude-stream-json-simple.jsonl`
