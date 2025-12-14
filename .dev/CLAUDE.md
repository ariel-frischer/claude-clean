# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Claude Clean Output is a terminal parser that transforms Claude Code's raw stream-json output into beautiful, human-readable terminal output with colors and formatting. It reads JSONL (JSON Lines) input either from stdin or files and displays formatted output in real-time.

## Common Commands

### Build and Run
```bash
make build          # Build the binary
make run            # Run with sample mock data (mocks/claude-stream-json-simple.log)
make run-verbose    # Run with verbose output (-v flag)
```

### Development
```bash
make all            # Format, vet, and build (run before commits)
make fmt            # Format code with gofmt
make vet            # Run go vet for suspicious code
make test           # Run tests
```

### Manual Execution
```bash
# From a file
./claude-clean-output mocks/your-test-file.log

# With verbose output (shows tool IDs, token usage)
./claude-clean-output -v mocks/your-test-file.log

# Different output styles
./claude-clean-output -s compact mocks/your-test-file.log
./claude-clean-output -s minimal mocks/your-test-file.log
./claude-clean-output -s plain mocks/your-test-file.log

# From stdin (for piping)
cat mocks/your-test-file.log | ./claude-clean-output

# Live with Claude Code
claude-code -p "your prompt" --verbose --output-format stream-json | ./claude-clean-output
```

## Code Architecture

### Single-File Design
The entire application is in `main.go` (~600 lines). This is intentional for simplicity and ease of maintenance.

### Message Processing Pipeline

1. **Input Stage**: Reads JSONL input line-by-line from stdin or file (main loop at main.go:165)
2. **Parsing Stage**: Unmarshals each JSON line into `StreamMessage` struct (main.go:176)
3. **Duplicate Detection**: Buffers assistant messages to avoid showing duplicates that appear in result messages (main.go:182-203)
4. **Display Routing**: Routes messages to appropriate display functions based on type and style (main.go:226-238)

### Key Data Structures

- **StreamMessage**: Top-level struct for all message types from Claude Code stream-json format (main.go:16)
  - Handles `system`, `assistant`, `user`, and `result` message types
  - Contains nested structures for message content, tool usage, and token statistics

- **MessageContent**: Assistant message details including content blocks and usage stats (main.go:38)

- **ContentBlock**: Represents individual content items - can be text, tool_use, or tool_result (main.go:49)

- **Usage**: Token usage statistics with cache information (main.go:60)

### Display System

The display system has a multi-style architecture:

1. **Style Selection** (main.go:92-99): Four output styles (default, compact, minimal, plain)
2. **Display Router** (main.go:226-238): Routes to style-specific formatters
3. **Message Type Handlers** (main.go:240-574):
   - `displaySystemMessage`: Cyan boxes showing initialization and config
   - `displayAssistantMessage`: Green boxes for Claude's responses
   - `displayToolUse`: Yellow boxes for tool invocations (Bash, Read, Edit, etc.)
   - `displayToolResult`: Magenta/Red boxes for tool execution results
   - `displayResultMessage`: Blue boxes for final session summary

### Special Features

- **Large Line Handling** (main.go:102-107): Scanner buffer set to 10MB to handle very large JSON lines from tool results or file contents. If a line exceeds 10MB (extremely rare), shows a detailed warning and exits gracefully with exit code 2.
- **Duplicate Detection** (main.go:182-203): Buffers assistant messages and compares with result messages to avoid showing the same content twice
- **Smart Truncation**:
  - Tool input strings: Shows first 200 + last 100 chars (main.go:343-350)
  - Tool result output: Shows first 20 + last 20 lines (main.go:459-485)
- **TodoWrite Special Handling** (main.go:352-398): Pretty-prints todo lists with status icons (✓, →, ○)
- **Verbose Mode** (main.go:103): Shows tool IDs, token usage details, and per-model statistics

### Color System

Colors are defined globally using `github.com/fatih/color` (main.go:75-89):
- System messages: Cyan
- Assistant text: Green
- Tool invocations: Yellow
- Tool results: Magenta (success) / Red (errors)
- Final results: Blue
- Metadata: Gray

## Adding New Features

### Supporting New Message Types
1. Add fields to `StreamMessage` struct (main.go:16)
2. Create a `displayXxxMessage()` function following existing patterns
3. Add case to switch statement in `displayMessageDefault()` (main.go:240)
4. Define colors if needed in the global color variables (main.go:75)

### Adding New Output Styles
1. Add style constant to OutputStyle enum (main.go:92)
2. Create `displayMessageXxx()` function for the new style
3. Add case to `displayMessage()` router (main.go:226)
4. Update help text to document the new style (main.go:121)

### Testing Changes
- Add sample JSON data to `mocks/` directory
- Run `make run` or `make run-verbose` to test parsing
- Test with different styles using `-s` flag
- Verify colors and formatting work correctly

