# Claude Clean Output

A beautiful terminal parser for Claude Code's streaming JSON output. Transforms the raw stream-json log format into clean, colorful, and readable terminal output.

## What is this for?

When you run Claude Code with verbose stream-json output, you get raw JSON logs that look like this:

```bash
claude-code -p "what is 2+2" --verbose --output-format stream-json
```

This produces a stream of JSON objects that are difficult to read in real-time. **Claude Clean Output** parses this stream and transforms it into beautiful, human-readable terminal output with colors, boxes, and proper formatting.

### Before (Raw JSON stream):

```json
{"type":"system","subtype":"init","cwd":"/home/ari/repos/claude-code-clean-output","session_id":"6170607e-7232-407c-82c3-7fc983d60064","tools":["Task","Bash","Glob","Grep","ExitPlanMode","Read","Edit","Write","NotebookEdit","WebFetch","TodoWrite","WebSearch","BashOutput","KillShell","Skill","SlashCommand","mcp__perplexity-mcp__perplexity_ask","mcp__perplexity-mcp__perplexity_reason","mcp__sequential-thinking__sequentialthinking"],"mcp_servers":[{"name":"perplexity-mcp","status":"connected"},{"name":"sequential-thinking","status":"connected"}],"model":"claude-sonnet-4-5-20250929","permissionMode":"bypassPermissions","slash_commands":["commit-commands:clean_gone","commit-commands:commit-push-pr","commit-commands:commit","feature-dev:feature-dev","compact","context","cost","init","output-style:new","pr-comments","release-notes","todos","review","security-review"],"apiKeySource":"none","claude_code_version":"2.0.25","output_style":"default","agents":["general-purpose","statusline-setup","output-style-setup","Explore","feature-dev:code-architect","feature-dev:code-explorer","feature-dev:code-reviewer","thoughts-analyzer","codebase-analyzer","base-template-generator","codebase-pattern-finder","web-search-researcher","codebase-locator","thoughts-locator"],"skills":[],"uuid":"268b2d2e-e6c1-4bd9-9e0d-fc6eef0e1779"}
{"type":"assistant","message":{"model":"claude-sonnet-4-5-20250929","id":"msg_01Rws28Xg2tBY3A5fNdrk6Mf","type":"message","role":"assistant","content":[{"type":"text","text":"I'll run a comprehensive diagnostic using all the requested tools."}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":3,"cache_creation_input_tokens":3117,"cache_read_input_tokens":14723,"cache_creation":{"ephemeral_5m_input_tokens":3117,"ephemeral_1h_input_tokens":0},"output_tokens":1,"service_tier":"standard"}},"parent_tool_use_id":null,"session_id":"6170607e-7232-407c-82c3-7fc983d60064","uuid":"68dda894-f1c3-4c75-b134-b364dd5d580d"}
{"type":"assistant","message":{"model":"claude-sonnet-4-5-20250929","id":"msg_01Rws28Xg2tBY3A5fNdrk6Mf","type":"message","role":"assistant","content":[{"type":"tool_use","id":"toolu_01VdNvyRGtzZvniXJGQQjvEP","name":"Glob","input":{"pattern":"**/*.go"}}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":3,"cache_creation_input_tokens":3117,"cache_read_input_tokens":14723,"cache_creation":{"ephemeral_5m_input_tokens":3117,"ephemeral_1h_input_tokens":0},"output_tokens":1,"service_tier":"standard"}},"parent_tool_use_id":null,"session_id":"6170607e-7232-407c-82c3-7fc983d60064","uuid":"e59b712a-7489-4166-80b9-13978fd2a6ea"}
```

### After (With Claude Clean Output):

```
┌─ SYSTEM [init] (line 1)
│ Working Directory: /home/ari/repos/claude-code-clean-output
│ Model: claude-sonnet-4-5-20250929
│ Claude Code: v2.0.25
│ Tools: 19 available
└─

┌─ ASSISTANT (line 2)
│ I'll run a comprehensive diagnostic using all the requested tools.
└─

┌─ TOOL: Glob (line 3)
│ Input:
│   pattern: **/*.go
└─
```

Much better! 🎨

This is the type of output this tool handles best: verbose stream-json logs from Claude Code.

## Features

- Beautiful color-coded output with box drawing characters
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

```bash
go build -o claude-clean-output
```

## Usage

### Live streaming with Claude Code (recommended):
```bash
claude-code -p "what is 2+2" --verbose --output-format stream-json | ./claude-clean-output
```

### From a log file:
```bash
./claude-clean-output mocks/claude-stream-json-log.log
```

### From stdin (pipe):
```bash
cat mocks/claude-stream-json-log.log | ./claude-clean-output
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
