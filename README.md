# Claude Clean

[![Pipeline](https://gitlab.com/ariel-frischer/claude-clean/badges/main/pipeline.svg)](https://gitlab.com/ariel-frischer/claude-clean/-/pipelines)
[![Release](https://img.shields.io/gitlab/v/release/ariel-frischer/claude-clean)](https://gitlab.com/ariel-frischer/claude-clean/-/releases)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A beautiful terminal parser for Claude Code's streaming JSON output. Transforms the raw stream-json log format into clean, colorful, and readable terminal output.

## What is this for?

When you run Claude Code with verbose stream-json output, you get raw JSON logs that look like this:

```bash
claude -p "what is 2+2" --verbose --output-format stream-json
```

This produces a stream of JSON objects that are difficult to read in real-time. **Claude Clean** parses this stream and transforms it into beautiful, human-readable terminal output with colors, boxes, and proper formatting.

### Before (Raw JSON stream):

```json
{"type":"system","subtype":"init","cwd":"/home/user/project","session_id":"6170607e-7232-407c-82c3-7fc983d60064","tools":["Task","Bash","Glob","Grep","ExitPlanMode","Read","Edit","Write","NotebookEdit","WebFetch","TodoWrite","WebSearch","BashOutput","KillShell","Skill","SlashCommand","mcp__perplexity-mcp__perplexity_ask","mcp__perplexity-mcp__perplexity_reason","mcp__sequential-thinking__sequentialthinking"],"mcp_servers":[{"name":"perplexity-mcp","status":"connected"},{"name":"sequential-thinking","status":"connected"}],"model":"claude-sonnet-4-5-20250929","permissionMode":"bypassPermissions","slash_commands":["commit-commands:clean_gone","commit-commands:commit-push-pr","commit-commands:commit","feature-dev:feature-dev","compact","context","cost","init","output-style:new","pr-comments","release-notes","todos","review","security-review"],"apiKeySource":"none","claude_code_version":"2.0.25","output_style":"default","agents":["general-purpose","statusline-setup","output-style-setup","Explore","feature-dev:code-architect","feature-dev:code-explorer","feature-dev:code-reviewer","thoughts-analyzer","codebase-analyzer","base-template-generator","codebase-pattern-finder","web-search-researcher","codebase-locator","thoughts-locator"],"skills":[],"uuid":"268b2d2e-e6c1-4bd9-9e0d-fc6eef0e1779"}
{"type":"assistant","message":{"model":"claude-sonnet-4-5-20250929","id":"msg_01Rws28Xg2tBY3A5fNdrk6Mf","type":"message","role":"assistant","content":[{"type":"text","text":"I'll run a comprehensive diagnostic using all the requested tools."}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":3,"cache_creation_input_tokens":3117,"cache_read_input_tokens":14723,"cache_creation":{"ephemeral_5m_input_tokens":3117,"ephemeral_1h_input_tokens":0},"output_tokens":1,"service_tier":"standard"}},"parent_tool_use_id":null,"session_id":"6170607e-7232-407c-82c3-7fc983d60064","uuid":"68dda894-f1c3-4c75-b134-b364dd5d580d"}
{"type":"assistant","message":{"model":"claude-sonnet-4-5-20250929","id":"msg_01Rws28Xg2tBY3A5fNdrk6Mf","type":"message","role":"assistant","content":[{"type":"tool_use","id":"toolu_01VdNvyRGtzZvniXJGQQjvEP","name":"Glob","input":{"pattern":"**/*.go"}}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":3,"cache_creation_input_tokens":3117,"cache_read_input_tokens":14723,"cache_creation":{"ephemeral_5m_input_tokens":3117,"ephemeral_1h_input_tokens":0},"output_tokens":1,"service_tier":"standard"}},"parent_tool_use_id":null,"session_id":"6170607e-7232-407c-82c3-7fc983d60064","uuid":"e59b712a-7489-4166-80b9-13978fd2a6ea"}
```

### After (With Claude Clean):

```
┌─ SYSTEM [init]
│ Working Directory: /home/user/project
│ Model: claude-sonnet-4-5-20250929
│ Claude Code: v2.0.25
│ Tools: 19 available
└─

┌─ ASSISTANT
│ I'll run a comprehensive diagnostic using all the requested tools.
└─

┌─ TOOL: Glob
│ Input:
│   pattern: **/*.go
└─
```

Much better! 🎨

This is the type of output this tool handles best: verbose stream-json logs from Claude Code.

**Bonus:** Using the `-p` flag with `--output-format stream-json` bypasses Claude's interactive UI/JS, which can help avoid segfault issues some users experience (e.g., with bun).

## Features

- **Multiple output styles**: Choose between default, compact, minimal, or plain formats
- Beautiful color-coded output with box drawing characters (default style)
- Supports stdin/stdout for piping
- Can read from log files
- Real-time streaming support
- Parses all Claude message types:
  - System initialization messages
  - Assistant text responses
  - Tool use calls (Bash, Read, Write, etc.)
  - Tool results (success and errors)
  - Token usage statistics

## Installation

### Option 1: One-line installer (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/claude-clean/main/install.sh | sh
```

This automatically detects your OS/architecture, downloads the correct binary, and installs it to `/usr/local/bin` (or `~/.local/bin` if no sudo access).

### Option 2: Manual download

Download from the [releases page](https://github.com/ariel-frischer/claude-clean/releases):

| Platform | Download |
|----------|----------|
| Linux (x86_64) | [claude-clean-linux-amd64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/claude-clean-linux-amd64) |
| Linux (ARM64) | [claude-clean-linux-arm64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/claude-clean-linux-arm64) |
| macOS (Intel) | [claude-clean-darwin-amd64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/claude-clean-darwin-amd64) |
| macOS (Apple Silicon) | [claude-clean-darwin-arm64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/claude-clean-darwin-arm64) |
| Windows | [claude-clean-windows-amd64.exe](https://github.com/ariel-frischer/claude-clean/releases/latest/download/claude-clean-windows-amd64.exe) |

After downloading, make it executable and move to your PATH:
```bash
chmod +x claude-clean-*
sudo mv claude-clean-* /usr/local/bin/claude-clean
```

### Option 3: Build from source

```bash
go build -o claude-clean
```

Or use the Makefile:
```bash
make install  # Builds and installs to ~/.local/bin with optional alias setup
```

## Usage

### Live streaming with Claude Code (recommended):
```bash
# Uses your configured API key (pay-per-use)
claude -p "what is 2+2" --verbose --output-format stream-json | claude-clean

# Or use OAuth with your Claude Pro/Team plan (FREE, no API costs)
ANTHROPIC_API_KEY="" claude -p "what is 2+2" --verbose --output-format stream-json | claude-clean
```

### From a log file:
```bash
claude-clean log.jsonl
```

### From stdin (pipe):
```bash
cat log.jsonl | claude-clean
```

### Using different output styles:
```bash
# Default style (boxed format with colors)
claude-clean log.jsonl

# Compact style (minimal single-line format)
claude-clean -s compact log.jsonl

# Minimal style (simple indented format, no boxes)
claude-clean -s minimal log.jsonl

# Plain style (no colors, no boxes)
claude-clean -s plain log.jsonl

# Verbose output with compact style
claude-clean -v -s compact log.jsonl
```

## Output Styles

The tool supports four different output styles, each optimized for different use cases:

### Default Style (Boxed Format)
The default style uses box drawing characters and colors for maximum readability:
- **SYSTEM** (Cyan): System initialization, configuration
- **ASSISTANT** (Green): Claude's text responses
- **TOOL** (Yellow): Tool invocations (Bash, Read, Write, etc.)
- **TOOL RESULT** (Gray/Magenta): Tool execution results
- **TOOL RESULT ERROR** (Red): Tool execution errors

Each message includes:
- Message type and details
- Content (truncated for readability)
- Line numbers (with `-l` flag)
- Token usage statistics (with `-V` flag)

### Compact Style
Minimal single-line format for quick scanning:
```
SYS[init] claude-sonnet-4-5-20250929 @/home/user/project
AST I'll help you with that task
TOOL Bash {description: "List files", command: "ls -la"}
RES total 24 drwxr-xr-x 3 user user 4096 ...
OK turns=2 1.23s $0.0012 in=150 out=45
```

### Minimal Style
Simple indented format without boxes, but with colors:
```
SYSTEM [init]
  Working Directory: /home/user/project
  Model: claude-sonnet-4-5-20250929

ASSISTANT
  I'll help you with that task

TOOL: Bash
  Input:
    command: ls -la

TOOL RESULT
  total 24
  drwxr-xr-x 3 user user 4096 ...
```

### Plain Style
No colors, no boxes - just plain text for logging or piping:
```
SYSTEM [init]
  Working Directory: /home/user/project
  Model: claude-sonnet-4-5-20250929

ASSISTANT
  I'll help you with that task
```

## Example Output

```
┌─ SYSTEM [init]
│ Working Directory: /home/user/project
│ Model: claude-sonnet-4-5-20250929
│ Claude Code: v2.0.25
│ Tools: 19 available
└─

┌─ ASSISTANT
│ I'll help you set up this Go project...
└─

┌─ TOOL: Bash
│ Input:
│   description: List repository contents
│   command: ls -la
└─

┌─ TOOL RESULT
│ total 24
│ drwxr-xr-x 3 user user 4096 Oct 21 12:00 .
│ ...
└─
```

## Dependencies

- [github.com/fatih/color](https://github.com/fatih/color) - Terminal color output
- [github.com/mattn/go-colorable](https://github.com/mattn/go-colorable) - Cross-platform colored terminal output (indirect)
- [github.com/mattn/go-isatty](https://github.com/mattn/go-isatty) - TTY detection (indirect)
- [golang.org/x/sys](https://golang.org/x/sys) - System call support (indirect)

## License

MIT
