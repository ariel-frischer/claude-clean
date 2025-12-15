package main

import (
	"strings"
	"testing"

	"github.com/fatih/color"
)

// ============================================================================
// PLAIN STYLE DISPLAY FUNCTION TESTS
// ============================================================================

// TestDisplaySystemMessagePlain tests displaySystemMessagePlain function
func TestDisplaySystemMessagePlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name        string
		msg         *StreamMessage
		lineNum     int
		showLine    bool
		contains    []string
		notContains []string
	}{
		{
			name: "Basic system message",
			msg: &StreamMessage{
				Type:    "system",
				Subtype: "init",
				CWD:     "/home/user/project",
				Model:   "claude-opus-4",
			},
			lineNum:  1,
			showLine: false,
			contains: []string{
				"SYSTEM",
				"[init]",
				"Working Directory: /home/user/project",
				"Model: claude-opus-4",
			},
		},
		{
			name: "System message with line number",
			msg: &StreamMessage{
				Type:  "system",
				Model: "claude-sonnet-4",
			},
			lineNum:  42,
			showLine: true,
			contains: []string{
				"SYSTEM",
				"(line 42)",
				"Model: claude-sonnet-4",
			},
		},
		{
			name: "System message with all fields",
			msg: &StreamMessage{
				Type:              "system",
				Subtype:           "start",
				CWD:               "/test/path",
				Model:             "claude-opus-4",
				ClaudeCodeVersion: "1.2.3",
				Tools:             []string{"Bash", "Read", "Write"},
			},
			lineNum:  5,
			showLine: false,
			contains: []string{
				"SYSTEM",
				"[start]",
				"Working Directory: /test/path",
				"Model: claude-opus-4",
				"Claude Code: v1.2.3",
				"Tools: 3 available",
			},
		},
		{
			name: "System message with no subtype",
			msg: &StreamMessage{
				Type:  "system",
				Model: "claude-haiku-4",
			},
			lineNum:     10,
			showLine:    false,
			contains:    []string{"SYSTEM", "Model: claude-haiku-4"},
			notContains: []string{"[]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			*showLineNum = tt.showLine
			defer func() { *showLineNum = oldShowLineNum }()

			output := captureStdout(func() {
				displaySystemMessagePlain(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)
			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}

			for _, notExp := range tt.notContains {
				if strings.Contains(cleaned, notExp) {
					t.Errorf("Expected output NOT to contain %q, but it did.\nGot:\n%s", notExp, cleaned)
				}
			}
		})
	}
}

// TestDisplayAssistantMessagePlain tests displayAssistantMessagePlain function
func TestDisplayAssistantMessagePlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name        string
		msg         *StreamMessage
		lineNum     int
		showLine    bool
		verbose     bool
		contains    []string
		notContains []string
		notEmpty    bool
	}{
		{
			name: "Simple text message",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "Hello, world!"},
					},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"ASSISTANT", "Hello, world!"},
			notEmpty: true,
		},
		{
			name: "Message with line number",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "Test message"},
					},
				},
			},
			lineNum:  25,
			showLine: true,
			verbose:  false,
			contains: []string{"ASSISTANT", "(line 25)", "Test message"},
			notEmpty: true,
		},
		{
			name: "Multiple text blocks",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "First block"},
						{Type: "text", Text: "Second block"},
					},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"ASSISTANT", "First block", "Second block"},
			notEmpty: true,
		},
		{
			name: "Message with usage info in verbose mode",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "Response text"},
					},
					Usage: &Usage{
						InputTokens:  100,
						OutputTokens: 50,
					},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  true,
			contains: []string{"ASSISTANT", "Response text", "Tokens: in=100 out=50"},
			notEmpty: true,
		},
		{
			name: "Message with cache usage in verbose mode",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "Cached response"},
					},
					Usage: &Usage{
						InputTokens:              200,
						OutputTokens:             75,
						CacheReadInputTokens:     150,
						CacheCreationInputTokens: 50,
					},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  true,
			contains: []string{"ASSISTANT", "Cached response", "Tokens: in=200 out=75", "cache_read=150", "cache_create=50"},
			notEmpty: true,
		},
		{
			name: "Empty message",
			msg: &StreamMessage{
				Type:    "assistant",
				Message: nil,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{},
			notEmpty: false,
		},
		{
			name: "Message with empty content",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{},
			notEmpty: false,
		},
		{
			name: "Usage not shown in non-verbose mode",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "Response"},
					},
					Usage: &Usage{
						InputTokens:  100,
						OutputTokens: 50,
					},
				},
			},
			lineNum:     1,
			showLine:    false,
			verbose:     false,
			contains:    []string{"ASSISTANT", "Response"},
			notContains: []string{"Tokens:"},
			notEmpty:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			oldVerbose := *verbose
			*showLineNum = tt.showLine
			*verbose = tt.verbose
			defer func() {
				*showLineNum = oldShowLineNum
				*verbose = oldVerbose
			}()

			output := captureStdout(func() {
				displayAssistantMessagePlain(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)

			if !tt.notEmpty && cleaned != "" {
				t.Errorf("Expected empty output, got: %q", cleaned)
			}

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}

			for _, notExp := range tt.notContains {
				if strings.Contains(cleaned, notExp) {
					t.Errorf("Expected output NOT to contain %q, but it did.\nGot:\n%s", notExp, cleaned)
				}
			}
		})
	}
}

// TestDisplayUserMessagePlain tests displayUserMessagePlain function
func TestDisplayUserMessagePlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		msg      *StreamMessage
		lineNum  int
		showLine bool
		verbose  bool
		contains []string
		notEmpty bool
	}{
		{
			name: "Tool result message",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type:      "tool_result",
							ToolUseID: "tool123",
							Content:   "Tool output here",
							IsError:   false,
						},
					},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL RESULT", "Tool output here"},
			notEmpty: true,
		},
		{
			name: "Empty user message",
			msg: &StreamMessage{
				Type:    "user",
				Message: nil,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{},
			notEmpty: false,
		},
		{
			name: "Empty content array",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{},
			notEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			oldVerbose := *verbose
			*showLineNum = tt.showLine
			*verbose = tt.verbose
			defer func() {
				*showLineNum = oldShowLineNum
				*verbose = oldVerbose
			}()

			output := captureStdout(func() {
				displayUserMessagePlain(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)

			if !tt.notEmpty && cleaned != "" {
				t.Errorf("Expected empty output, got: %q", cleaned)
			}

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}
		})
	}
}

// TestDisplayToolUsePlain tests displayToolUsePlain function
func TestDisplayToolUsePlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name        string
		tool        *ContentBlock
		lineNum     int
		showLine    bool
		verbose     bool
		contains    []string
		notContains []string
	}{
		{
			name: "Simple tool use",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "Bash",
				ID:   "tool123",
				Input: map[string]interface{}{
					"command": "ls -la",
				},
			},
			lineNum:     1,
			showLine:    false,
			verbose:     false,
			contains:    []string{"TOOL: Bash", "Input:", "command:", "ls -la"},
			notContains: []string{"ID:"},
		},
		{
			name: "Tool use with line number",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "Read",
				ID:   "tool456",
				Input: map[string]interface{}{
					"file_path": "/home/user/test.txt",
				},
			},
			lineNum:  15,
			showLine: true,
			verbose:  false,
			contains: []string{"TOOL: Read", "(line 15)", "Input:", "file_path:", "/home/user/test.txt"},
		},
		{
			name: "Tool use with ID in verbose mode",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "Write",
				ID:   "tool789",
				Input: map[string]interface{}{
					"file_path": "/tmp/output.txt",
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  true,
			contains: []string{"TOOL: Write", "ID: tool789", "Input:", "file_path:", "/tmp/output.txt"},
		},
		{
			name: "Tool use with long string truncation",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "Edit",
				ID:   "tool999",
				Input: map[string]interface{}{
					"old_string": strings.Repeat("a", 400),
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL: Edit", "Input:", "old_string:", "chars omitted"},
		},
		{
			name: "Tool use with array input",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "TestTool",
				ID:   "tool111",
				Input: map[string]interface{}{
					"items": []interface{}{"item1", "item2", "item3"},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL: TestTool", "Input:", "items:", "[3 items]"},
		},
		{
			name: "Tool use with map input",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "ComplexTool",
				ID:   "tool222",
				Input: map[string]interface{}{
					"config": map[string]interface{}{"key": "value"},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL: ComplexTool", "Input:", "config:", "{...}"},
		},
		{
			name: "TodoWrite tool with todos",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "TodoWrite",
				ID:   "tool333",
				Input: map[string]interface{}{
					"todos": []interface{}{
						map[string]interface{}{
							"content": "First task",
							"status":  "completed",
						},
						map[string]interface{}{
							"content": "Second task",
							"status":  "in_progress",
						},
					},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL: TodoWrite", "Input:", "todos:", "First task", "Second task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			oldVerbose := *verbose
			*showLineNum = tt.showLine
			*verbose = tt.verbose
			defer func() {
				*showLineNum = oldShowLineNum
				*verbose = oldVerbose
			}()

			output := captureStdout(func() {
				displayToolUsePlain(tt.tool, tt.lineNum)
			})

			cleaned := stripANSI(output)

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}

			for _, notExp := range tt.notContains {
				if strings.Contains(cleaned, notExp) {
					t.Errorf("Expected output NOT to contain %q, but it did.\nGot:\n%s", notExp, cleaned)
				}
			}
		})
	}
}

// TestDisplayToolResultPlain tests displayToolResultPlain function
func TestDisplayToolResultPlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name        string
		block       *ContentBlock
		lineNum     int
		showLine    bool
		verbose     bool
		contains    []string
		notContains []string
	}{
		{
			name: "Successful tool result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool123",
				Content:   "Command executed successfully",
				IsError:   false,
			},
			lineNum:     1,
			showLine:    false,
			verbose:     false,
			contains:    []string{"TOOL RESULT", "Command executed successfully"},
			notContains: []string{"ERROR", "Tool ID:"},
		},
		{
			name: "Error tool result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool456",
				Content:   "Command failed: permission denied",
				IsError:   true,
			},
			lineNum:     1,
			showLine:    false,
			verbose:     false,
			contains:    []string{"TOOL RESULT ERROR", "Command failed: permission denied"},
			notContains: []string{"Tool ID:"},
		},
		{
			name: "Tool result with line number",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool789",
				Content:   "Output text",
				IsError:   false,
			},
			lineNum:  30,
			showLine: true,
			verbose:  false,
			contains: []string{"TOOL RESULT", "(line 30)", "Output text"},
		},
		{
			name: "Tool result with ID in verbose mode",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool999",
				Content:   "Verbose output",
				IsError:   false,
			},
			lineNum:  1,
			showLine: false,
			verbose:  true,
			contains: []string{"TOOL RESULT", "Tool ID: tool999", "Verbose output"},
		},
		{
			name: "Empty tool result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool111",
				Content:   "",
				IsError:   false,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL RESULT", "(no output)"},
		},
		{
			name: "Tool result with system reminder stripped",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool222",
				Content:   "Output here\n<system-reminder>This should be removed</system-reminder>\nMore output",
				IsError:   false,
			},
			lineNum:     1,
			showLine:    false,
			verbose:     false,
			contains:    []string{"TOOL RESULT", "Output here", "More output"},
			notContains: []string{"system-reminder", "This should be removed"},
		},
		{
			name: "Tool result with many lines (truncated)",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool333",
				Content:   generateTestLines(50),
				IsError:   false,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL RESULT", "Line 1", "Line 20", "more lines", "Line 31", "Line 50"},
		},
		{
			name: "Tool result with few lines (not truncated)",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool444",
				Content:   "Line 1\nLine 2\nLine 3",
				IsError:   false,
			},
			lineNum:     1,
			showLine:    false,
			verbose:     false,
			contains:    []string{"TOOL RESULT", "Line 1", "Line 2", "Line 3"},
			notContains: []string{"more lines"},
		},
		{
			name: "Tool result with non-string content",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool555",
				Content:   123,
				IsError:   false,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"TOOL RESULT", "123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			oldVerbose := *verbose
			*showLineNum = tt.showLine
			*verbose = tt.verbose
			defer func() {
				*showLineNum = oldShowLineNum
				*verbose = oldVerbose
			}()

			output := captureStdout(func() {
				displayToolResultPlain(tt.block, tt.lineNum)
			})

			cleaned := stripANSI(output)

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}

			for _, notExp := range tt.notContains {
				if strings.Contains(cleaned, notExp) {
					t.Errorf("Expected output NOT to contain %q, but it did.\nGot:\n%s", notExp, cleaned)
				}
			}
		})
	}
}

// TestDisplayTodosPlain tests displayTodosPlain function
func TestDisplayTodosPlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name        string
		todos       []interface{}
		contains    []string
		notContains []string
	}{
		{
			name: "Mixed status todos",
			todos: []interface{}{
				map[string]interface{}{"content": "Completed task", "status": "completed"},
				map[string]interface{}{"content": "In progress task", "status": "in_progress"},
				map[string]interface{}{"content": "Pending task", "status": "pending"},
			},
			contains: []string{
				"[✓] Completed task",
				"[→] In progress task",
				"[○] Pending task",
			},
		},
		{
			name: "Unknown status",
			todos: []interface{}{
				map[string]interface{}{"content": "Unknown status task", "status": "unknown"},
			},
			contains: []string{"[-] Unknown status task"},
		},
		{
			name: "Empty todos",
			todos: []interface{}{},
		},
		{
			name: "Invalid todo format",
			todos: []interface{}{
				"not a map",
				123,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				displayTodosPlain(tt.todos)
			})

			cleaned := stripANSI(output)

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}

			for _, notExp := range tt.notContains {
				if strings.Contains(cleaned, notExp) {
					t.Errorf("Expected output NOT to contain %q, but it did.\nGot:\n%s", notExp, cleaned)
				}
			}
		})
	}
}

// TestDisplayResultMessagePlain tests displayResultMessagePlain function
func TestDisplayResultMessagePlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name        string
		msg         *StreamMessage
		lineNum     int
		showLine    bool
		verbose     bool
		contains    []string
		notContains []string
	}{
		{
			name: "Success result",
			msg: &StreamMessage{
				Type:         "result",
				IsError:      false,
				NumTurns:     5,
				DurationMS:   15000,
				TotalCostUSD: 0.0025,
				Usage: &Usage{
					InputTokens:  1000,
					OutputTokens: 500,
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{
				"RESULT: SUCCESS",
				"Turns: 5",
				"Duration: 15.00s",
				"Cost: $0.0025",
				"Tokens: in=1000 out=500",
			},
			notContains: []string{"ERROR"},
		},
		{
			name: "Error result",
			msg: &StreamMessage{
				Type:    "result",
				IsError: true,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"RESULT: ERROR"},
		},
		{
			name: "Result with line number",
			msg: &StreamMessage{
				Type:    "result",
				IsError: false,
			},
			lineNum:  47,
			showLine: true,
			verbose:  false,
			contains: []string{"RESULT: SUCCESS", "(line 47)"},
		},
		{
			name: "Result with API duration",
			msg: &StreamMessage{
				Type:          "result",
				IsError:       false,
				DurationMS:    10000,
				DurationAPIMS: 8000,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"Duration: 10.00s", "(API: 8.00s)"},
		},
		{
			name: "Result with cache tokens",
			msg: &StreamMessage{
				Type:    "result",
				IsError: false,
				Usage: &Usage{
					InputTokens:              2000,
					OutputTokens:             1000,
					CacheReadInputTokens:     1500,
					CacheCreationInputTokens: 500,
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{
				"Tokens: in=2000 out=1000",
				"cache_read=1500",
				"cache_create=500",
			},
		},
		{
			name: "Result with result text",
			msg: &StreamMessage{
				Type:    "result",
				IsError: false,
				Result:  "Task completed successfully",
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"Task completed successfully"},
		},
		{
			name: "Result with permission denials",
			msg: &StreamMessage{
				Type:              "result",
				IsError:           false,
				PermissionDenials: []interface{}{"Denied action 1", "Denied action 2"},
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"Permission Denials: 2"},
		},
		{
			name: "Result with model usage in verbose mode",
			msg: &StreamMessage{
				Type:    "result",
				IsError: false,
				ModelUsage: map[string]interface{}{
					"claude-opus-4": map[string]interface{}{
						"inputTokens":  float64(1000),
						"outputTokens": float64(500),
						"costUSD":      float64(0.05),
					},
				},
			},
			lineNum:  1,
			showLine: false,
			verbose:  true,
			contains: []string{
				"Model Usage:",
				"claude-opus-4:",
				"Input: 1000 tokens",
				"Output: 500 tokens",
				"Cost: $0.0500",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			oldVerbose := *verbose
			*showLineNum = tt.showLine
			*verbose = tt.verbose
			defer func() {
				*showLineNum = oldShowLineNum
				*verbose = oldVerbose
			}()

			output := captureStdout(func() {
				displayResultMessagePlain(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}

			for _, notExp := range tt.notContains {
				if strings.Contains(cleaned, notExp) {
					t.Errorf("Expected output NOT to contain %q, but it did.\nGot:\n%s", notExp, cleaned)
				}
			}
		})
	}
}

// TestDisplayMessagePlain tests displayMessagePlain routing function
func TestDisplayMessagePlain(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		msg      *StreamMessage
		lineNum  int
		contains []string
	}{
		{
			name: "Routes system message",
			msg: &StreamMessage{
				Type:  "system",
				Model: "claude-opus-4",
			},
			lineNum:  1,
			contains: []string{"SYSTEM", "Model: claude-opus-4"},
		},
		{
			name: "Routes assistant message",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "Hello!"},
					},
				},
			},
			lineNum:  2,
			contains: []string{"ASSISTANT", "Hello!"},
		},
		{
			name: "Routes user message",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "tool_result", ToolUseID: "t1", Content: "Result"},
					},
				},
			},
			lineNum:  3,
			contains: []string{"TOOL RESULT", "Result"},
		},
		{
			name: "Routes result message",
			msg: &StreamMessage{
				Type:    "result",
				IsError: false,
			},
			lineNum:  4,
			contains: []string{"RESULT: SUCCESS"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				displayMessagePlain(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)
			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot:\n%s", expected, cleaned)
				}
			}
		})
	}
}

// Helper function to generate numbered test lines
func generateTestLines(count int) string {
	var sb strings.Builder
	for i := 1; i <= count; i++ {
		if i > 1 {
			sb.WriteString("\n")
		}
		sb.WriteString("Line ")
		sb.WriteString(itoa(i))
	}
	return sb.String()
}

// Simple int to string conversion
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
