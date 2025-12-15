# Missing Unit Tests

## Overview

Current test coverage is estimated at **<5%**. Only `stripSystemReminders()` has proper test coverage. This document lists all functions and areas that need unit tests.

---

## High Priority - Core Functions

### `main()` - Entry Point Logic
- [ ] Test command-line flag parsing (`-v`, `--version`, `-h`, `--help`, `--style`, `--oauth`, `-n`)
- [ ] Test file path detection vs prompt string detection
- [ ] Test stdin handling for piped input
- [ ] Test invalid style flag values
- [ ] Test version/help output exits correctly
- [ ] Test error handling when no input provided

### `fileExists(path string) bool`
- [ ] Test with existing file
- [ ] Test with non-existing file
- [ ] Test with directory path
- [ ] Test with permission-denied path
- [ ] Test with empty string

### `runClaude(prompt string) error`
- [ ] Test successful command execution
- [ ] Test OAuth environment variable is set when `useOAuth` is true
- [ ] Test error when claude binary not found
- [ ] Test stdout pipe creation
- [ ] Test large output handling (10MB buffer limit)

### `binaryName() string`
- [ ] Test returns correct executable name
- [ ] Test handles path with directories

### `printVersion()`
- [ ] Test version output format
- [ ] Test version variables are populated

### `printHelp()`
- [ ] Test help output contains usage info
- [ ] Test help output contains all flags
- [ ] Test help output contains examples

---

## High Priority - Display Functions

### Default Style

#### `displayMessage(msg *StreamMessage, lineNum int)`
- [ ] Test routing to correct display function based on message type
- [ ] Test nil message handling

#### `displayMessageDefault(msg *StreamMessage, lineNum int)`
- [ ] Test routing based on msg.Type values: "system", "assistant", "user", "result"
- [ ] Test fallback for unknown types

#### `displaySystemMessage(msg *StreamMessage, lineNum int)`
- [ ] Test output format with model and version
- [ ] Test with missing Model field
- [ ] Test with missing ClaudeCodeVersion field

#### `displayAssistantMessage(msg *StreamMessage, lineNum int)`
- [ ] Test text content display
- [ ] Test tool_use content detection and routing
- [ ] Test multiple content blocks
- [ ] Test empty content handling

#### `displayToolUse(msg *StreamMessage, lineNum int)`
- [ ] Test tool name and ID display
- [ ] Test input parameters display (JSON formatting)
- [ ] Test TodoWrite special handling routes to `displayTodos()`
- [ ] Test with nil/empty Input

#### `displayTodos(msg *StreamMessage, lineNum int)`
- [ ] Test todo item display format
- [ ] Test status icons (pending, in_progress, completed)
- [ ] Test empty todos array
- [ ] Test invalid todo input format

#### `displayUserMessage(msg *StreamMessage, lineNum int)`
- [ ] Test text content display
- [ ] Test tool_result content routing
- [ ] Test system reminder stripping

#### `displayToolResult(msg *StreamMessage, lineNum int)`
- [ ] Test success result display
- [ ] Test error result display (isError: true)
- [ ] Test truncated content handling

#### `displayResultMessage(msg *StreamMessage, lineNum int)`
- [ ] Test result display with cost and duration
- [ ] Test error result display
- [ ] Test routes to `displayUsage()`

#### `displayUsage(usage Usage)`
- [ ] Test token counts display
- [ ] Test cache statistics display
- [ ] Test zero values handling

---

### Compact Style

#### `displayMessageCompact(msg *StreamMessage, lineNum int)`
- [ ] Test routing logic mirrors default but with compact format

#### `displaySystemMessageCompact(msg *StreamMessage, lineNum int)`
- [ ] Test compact format output

#### `displayAssistantMessageCompact(msg *StreamMessage, lineNum int)`
- [ ] Test compact text display
- [ ] Test tool_use handling

#### `displayToolUseCompact(msg *StreamMessage, lineNum int)`
- [ ] Test compact tool display format

#### `displayUserMessageCompact(msg *StreamMessage, lineNum int)`
- [ ] Test compact user message

#### `displayToolResultCompact(msg *StreamMessage, lineNum int)`
- [ ] Test compact tool result

#### `displayResultMessageCompact(msg *StreamMessage, lineNum int)`
- [ ] Test compact result with inline usage

---

### Minimal Style

#### `displayMessageMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal routing

#### `displaySystemMessageMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal system message (no output expected)

#### `displayAssistantMessageMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal assistant output

#### `displayToolUseMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal tool display

#### `displayTodosMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal todos display

#### `displayUserMessageMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal user message

#### `displayToolResultMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal tool result

#### `displayResultMessageMinimal(msg *StreamMessage, lineNum int)`
- [ ] Test minimal result display

---

### Plain Style

#### `displayMessagePlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain routing

#### `displaySystemMessagePlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain system message

#### `displayAssistantMessagePlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain assistant output (no colors)

#### `displayToolUsePlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain tool display (no box drawing)

#### `displayTodosPlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain todos display

#### `displayUserMessagePlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain user message

#### `displayToolResultPlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain tool result

#### `displayResultMessagePlain(msg *StreamMessage, lineNum int)`
- [ ] Test plain result display

---

## Medium Priority - Utility Functions

### `formatLineNum(lineNum int) string`
- [ ] Test single digit formatting
- [ ] Test multi-digit formatting
- [ ] Test when `showLineNum` is false (returns empty)

### `formatLineNumCompact(lineNum int) string`
- [ ] Test compact line number format
- [ ] Test when `showLineNum` is false

### `displayUsageInline(usage Usage) string`
- [ ] Test inline format output
- [ ] Test with cache tokens
- [ ] Test with zero values

---

## Medium Priority - Type Validation

### StreamMessage struct
- [ ] Test JSON unmarshal with all field types
- [ ] Test JSON unmarshal with missing fields (defaults)
- [ ] Test JSON unmarshal with extra fields (ignored)
- [ ] Test JSON unmarshal with invalid types

### MessageContent struct
- [ ] Test nested Content array parsing
- [ ] Test Usage field parsing

### ContentBlock struct
- [ ] Test all block types: text, tool_use, tool_result
- [ ] Test Input field as map[string]interface{}

### Usage struct
- [ ] Test CacheCreation nested object parsing

---

## Low Priority - Edge Cases

### Error Handling
- [ ] Test graceful handling of malformed JSON input
- [ ] Test handling of very long lines (>10MB)
- [ ] Test handling of binary/non-UTF8 input
- [ ] Test handling of interrupted stdin

### Output Formatting
- [ ] Test color codes are applied correctly
- [ ] Test box drawing characters render properly
- [ ] Test multiline content indentation
- [ ] Test system reminder stripping in various positions

---

## Testing Strategy Notes

1. **Display functions** should capture stdout and compare against expected output (snapshot testing approach)
2. **Main function** may need refactoring to inject dependencies (stdin, stdout, args)
3. **runClaude** will need mocking of exec.Command
4. Consider using `testify` or similar for assertions
5. Add more mock JSONL files in `/mocks` for different scenarios

---

## Current Test Coverage

| File | Functions | Tested | Coverage |
|------|-----------|--------|----------|
| main.go | 47 | 1 | ~2% |
| types.go | 0 (structs only) | N/A | N/A |

**Total estimated coverage: <5%**
