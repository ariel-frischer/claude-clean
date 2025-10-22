# Claude Clean Output

A beautiful terminal parser for Claude Code's streaming JSON output. Transforms the raw stream-json log format into clean, colorful, and readable terminal output.

## Features

- Beautiful color-coded output with box drawing characters
- Supports stdin/stdout for piping
- Can read from log files
- Parses all Claude message types:
  - System initialization messages
  - Assistant text responses
  - Tool use calls (Bash, Read, Write, etc.)
  - Tool results (success and errors)
  - Token usage statistics

## Installation

```bash
go build -o claude-clean-output
```

## Usage

### From a log file:
```bash
./claude-clean-output mocks/claude-stream-json-log.log
```

### From stdin (pipe):
```bash
cat mocks/claude-stream-json-log.log | ./claude-clean-output
```

### Live streaming:
```bash
# If you have Claude Code outputting to a stream
your-claude-command | ./claude-clean-output
```

## Output Format

The parser formats different message types with distinct colors and styles:

- **SYSTEM** (Cyan): System initialization, configuration
- **ASSISTANT** (Green): Claude's text responses
- **TOOL** (Yellow): Tool invocations (Bash, Read, Write, etc.)
- **TOOL RESULT** (Gray/Magenta): Tool execution results
- **TOOL RESULT ERROR** (Red): Tool execution errors

Each message includes:
- Line number from the source file
- Message type and details
- Content (truncated for readability)
- Token usage statistics (for assistant messages)

## Example Output

```
┌─ SYSTEM [init] (line 1)
│ Working Directory: /home/user/project
│ Model: claude-sonnet-4-5-20250929
│ Claude Code: v2.0.25
│ Tools: 19 available
└─

┌─ ASSISTANT (line 2)
│ I'll help you set up this Go project...
│ Tokens: in=2 out=3 cache_create=17639
└─

┌─ TOOL: Bash (line 4)
│ ID: toolu_011FM74Ft8CR2ji5uPhY2k6d
│ Input:
│   description: List repository contents
│   command: ls -la
└─

┌─ TOOL RESULT (line 8)
│ Tool ID: toolu_011FM74Ft8CR2ji5uPhY2k6d
│ total 24
│ drwxr-xr-x 3 user user 4096 Oct 21 12:00 .
│ ...
└─
```

## Dependencies

- [github.com/fatih/color](https://github.com/fatih/color) - Terminal color output

## License

MIT
