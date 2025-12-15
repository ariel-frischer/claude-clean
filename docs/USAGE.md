# USAGE

How to use cclean to format Claude Code's streaming JSON output.

## Basic Usage

### Pipe from Claude Code (Recommended)

```bash
claude -p "your prompt" --verbose --output-format stream-json | cclean
```

### Read from File

```bash
cclean logfile.jsonl
```

### Read from Stdin

```bash
cat logfile.jsonl | cclean
```

## Shell Alias Setup

For convenience, add a shell function to your config:

### Bash/Zsh (`~/.bashrc` or `~/.zshrc`)

```bash
# Using OAuth (Claude Pro/Team plan - FREE)
cc() {
  ANTHROPIC_API_KEY="" claude -p "$*" --verbose --output-format stream-json | cclean
}

# Or using API key (pay-per-use)
cc() {
  claude -p "$*" --verbose --output-format stream-json | cclean
}
```

### Fish (`~/.config/fish/config.fish`)

```fish
function cc
  env ANTHROPIC_API_KEY="" claude -p $argv --verbose --output-format stream-json | cclean
end
```

Then use simply:

```bash
cc "what is 2+2"
cc "help me debug this error"
```

## Command Line Flags

| Flag | Description |
|------|-------------|
| `-s <style>` | Output style: `default`, `compact`, `minimal`, `plain` |
| `-v` | Verbose mode (more details) |
| `-V` | Very verbose (includes token stats) |
| `-l` | Show line numbers |
| `--version` | Show version info |
| `-h`, `--help` | Show help |

## Output Styles

### Default (Boxed)

Beautiful box-drawn output with colors:

```
┌─ SYSTEM [init]
│ Working Directory: /home/user/project
│ Model: claude-sonnet-4-5-20250929
└─

┌─ ASSISTANT
│ I'll help you with that task...
└─

┌─ TOOL: Bash
│ Input:
│   command: ls -la
└─
```

### Compact

Single-line format for quick scanning:

```
SYS[init] claude-sonnet-4-5-20250929 @/home/user/project
AST I'll help you with that task
TOOL Bash {command: "ls -la"}
RES total 24 drwxr-xr-x ...
OK turns=2 1.23s $0.0012
```

Use with:
```bash
cclean -s compact logfile.jsonl
```

### Minimal

Simple indented format, colors but no boxes:

```
SYSTEM [init]
  Working Directory: /home/user/project
  Model: claude-sonnet-4-5-20250929

ASSISTANT
  I'll help you with that task...

TOOL: Bash
  Input:
    command: ls -la
```

Use with:
```bash
cclean -s minimal logfile.jsonl
```

### Plain

No colors, no boxes - ideal for logging or piping to files:

```
SYSTEM [init]
  Working Directory: /home/user/project

ASSISTANT
  I'll help you with that task...
```

Use with:
```bash
cclean -s plain logfile.jsonl > output.txt
```

## Message Types

cclean parses and formats these message types:

| Type | Color | Description |
|------|-------|-------------|
| SYSTEM | Cyan | Initialization, config, session info |
| ASSISTANT | Green | Claude's text responses |
| TOOL | Yellow | Tool invocations (Bash, Read, Write, etc.) |
| TOOL RESULT | Gray | Successful tool execution results |
| TOOL RESULT ERROR | Red | Failed tool executions |
| RESULT | Magenta | Final result/summary |

## Examples

### Basic prompt

```bash
cc "what is 2+2"
```

### Complex task with verbose output

```bash
claude -p "refactor main.go to use interfaces" \
  --verbose --output-format stream-json | cclean -V
```

### Save formatted output

```bash
claude -p "explain this code" --verbose --output-format stream-json \
  | cclean -s plain > explanation.txt
```

### View with line numbers

```bash
cclean -l -s minimal logfile.jsonl
```

### Real-time streaming

The tool handles real-time streaming naturally - output appears as Claude generates it:

```bash
claude -p "write a long story" --verbose --output-format stream-json | cclean
```

## Why Use This?

1. **Readability** - Raw JSON streams are hard to follow; this makes them beautiful
2. **Real-time** - See output as it streams, not just at the end
3. **Debugging** - Easily see tool calls, inputs, and results
4. **Avoid Segfaults** - Using `-p` with `--output-format stream-json` bypasses Claude's interactive UI, avoiding potential segfault issues
