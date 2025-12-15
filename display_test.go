package main

import (
	"strings"
	"testing"

	"github.com/fatih/color"
)

// TestDisplayAssistantMessage tests the default style displayAssistantMessage function
func TestDisplayAssistantMessage(t *testing.T) {
	// Disable colors for testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		msg              *StreamMessage
		lineNum          int
		showLineNum      bool
		verboseMode      bool
		expectedIncludes []string
		expectedExcludes []string
	}{
		{
			name: "Nil message",
			msg: &StreamMessage{
				Type:    "assistant",
				Message: nil,
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{},
		},
		{
			name: "Empty content",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{},
				},
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{},
		},
		{
			name: "Simple text message",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type: "text",
							Text: "Hello, world!",
						},
					},
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"ASSISTANT",
				"Hello, world!",
				"└─",
			},
		},
		{
			name: "Multiple text blocks",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type: "text",
							Text: "First block",
						},
						{
							Type: "text",
							Text: "Second block",
						},
					},
				},
			},
			lineNum:     2,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"ASSISTANT",
				"First block",
				"Second block",
			},
		},
		{
			name: "Text with usage in verbose mode",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type: "text",
							Text: "Response text",
						},
					},
					Usage: &Usage{
						InputTokens:  100,
						OutputTokens: 50,
					},
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: true,
			expectedIncludes: []string{
				"ASSISTANT",
				"Response text",
				"Tokens: in=100 out=50",
			},
		},
		{
			name: "Text without usage in non-verbose mode",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type: "text",
							Text: "Response text",
						},
					},
					Usage: &Usage{
						InputTokens:  100,
						OutputTokens: 50,
					},
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"ASSISTANT",
				"Response text",
			},
			expectedExcludes: []string{
				"Tokens:",
			},
		},
		{
			name: "With line number",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type: "text",
							Text: "Test message",
						},
					},
				},
			},
			lineNum:     99,
			showLineNum: true,
			verboseMode: false,
			expectedIncludes: []string{
				"ASSISTANT",
				"(line 99)",
				"Test message",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origShowLineNum := *showLineNum
			origVerbose := *verbose
			*showLineNum = tt.showLineNum
			*verbose = tt.verboseMode
			defer func() {
				*showLineNum = origShowLineNum
				*verbose = origVerbose
			}()

			output := captureStdout(func() {
				displayAssistantMessage(tt.msg, tt.lineNum)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displayAssistantMessage() output missing %q\nGot:\n%s", expected, output)
				}
			}

			for _, excluded := range tt.expectedExcludes {
				if strings.Contains(output, excluded) {
					t.Errorf("displayAssistantMessage() output should not contain %q\nGot:\n%s", excluded, output)
				}
			}
		})
	}
}

// TestDisplayUserMessage tests the default style displayUserMessage function
func TestDisplayUserMessage(t *testing.T) {
	// Disable colors for testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		msg              *StreamMessage
		lineNum          int
		showLineNum      bool
		verboseMode      bool
		expectedIncludes []string
	}{
		{
			name: "Nil message",
			msg: &StreamMessage{
				Type:    "user",
				Message: nil,
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{},
		},
		{
			name: "Empty content",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{},
				},
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{},
		},
		{
			name: "Tool result success",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type:      "tool_result",
							ToolUseID: "tool_123",
							Content:   "Command executed successfully",
							IsError:   false,
						},
					},
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT",
				"Command executed successfully",
			},
		},
		{
			name: "Tool result error",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type:      "tool_result",
							ToolUseID: "tool_456",
							Content:   "Error: command failed",
							IsError:   true,
						},
					},
				},
			},
			lineNum:     2,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT ERROR",
				"Error: command failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origShowLineNum := *showLineNum
			origVerbose := *verbose
			*showLineNum = tt.showLineNum
			*verbose = tt.verboseMode
			defer func() {
				*showLineNum = origShowLineNum
				*verbose = origVerbose
			}()

			output := captureStdout(func() {
				displayUserMessage(tt.msg, tt.lineNum)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displayUserMessage() output missing %q\nGot:\n%s", expected, output)
				}
			}
		})
	}
}

// TestDisplayToolUse tests the default style displayToolUse function
func TestDisplayToolUse(t *testing.T) {
	// Disable colors for testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		block            *ContentBlock
		lineNum          int
		showLineNum      bool
		verboseMode      bool
		expectedIncludes []string
		expectedExcludes []string
	}{
		{
			name: "Simple tool use",
			block: &ContentBlock{
				Type: "tool_use",
				Name: "Bash",
				ID:   "tool_abc",
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL: Bash",
				"└─",
			},
			expectedExcludes: []string{
				"ID:",
			},
		},
		{
			name: "Tool use with verbose ID",
			block: &ContentBlock{
				Type: "tool_use",
				Name: "Read",
				ID:   "tool_xyz",
			},
			lineNum:     5,
			showLineNum: false,
			verboseMode: true,
			expectedIncludes: []string{
				"TOOL: Read",
				"ID: tool_xyz",
			},
		},
		{
			name: "Tool use with string input",
			block: &ContentBlock{
				Type: "tool_use",
				Name: "Bash",
				ID:   "tool_123",
				Input: map[string]interface{}{
					"command": "ls -la",
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL: Bash",
				"Input:",
				"command:",
				"ls -la",
			},
		},
		{
			name: "Tool use with long string input",
			block: &ContentBlock{
				Type: "tool_use",
				Name: "Write",
				ID:   "tool_456",
				Input: map[string]interface{}{
					"file_path": "/home/user/file.txt",
					"content":   strings.Repeat("a", 400),
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL: Write",
				"file_path:",
				"/home/user/file.txt",
				"content:",
				"chars omitted",
			},
		},
		{
			name: "Tool use with array input",
			block: &ContentBlock{
				Type: "tool_use",
				Name: "SomeTool",
				ID:   "tool_789",
				Input: map[string]interface{}{
					"items": []interface{}{"item1", "item2", "item3"},
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL: SomeTool",
				"items:",
				"[3 items]",
			},
		},
		{
			name: "Tool use with line number",
			block: &ContentBlock{
				Type: "tool_use",
				Name: "Edit",
				ID:   "tool_edit",
			},
			lineNum:     42,
			showLineNum: true,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL: Edit",
				"(line 42)",
			},
		},
		{
			name: "TodoWrite tool with todos",
			block: &ContentBlock{
				Type: "tool_use",
				Name: "TodoWrite",
				ID:   "tool_todo",
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
						map[string]interface{}{
							"content": "Third task",
							"status":  "pending",
						},
					},
				},
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL: TodoWrite",
				"Input:",
				"todos:",
				"First task",
				"Second task",
				"Third task",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origShowLineNum := *showLineNum
			origVerbose := *verbose
			*showLineNum = tt.showLineNum
			*verbose = tt.verboseMode
			defer func() {
				*showLineNum = origShowLineNum
				*verbose = origVerbose
			}()

			output := captureStdout(func() {
				displayToolUse(tt.block, tt.lineNum)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displayToolUse() output missing %q\nGot:\n%s", expected, output)
				}
			}

			for _, excluded := range tt.expectedExcludes {
				if strings.Contains(output, excluded) {
					t.Errorf("displayToolUse() output should not contain %q\nGot:\n%s", excluded, output)
				}
			}
		})
	}
}

// TestDisplayToolResult tests the default style displayToolResult function
func TestDisplayToolResult(t *testing.T) {
	// Disable colors for testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		block            *ContentBlock
		lineNum          int
		showLineNum      bool
		verboseMode      bool
		expectedIncludes []string
		expectedExcludes []string
	}{
		{
			name: "Successful tool result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_123",
				Content:   "Success output",
				IsError:   false,
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT",
				"Success output",
				"└─",
			},
			expectedExcludes: []string{
				"ERROR",
				"Tool ID:",
			},
		},
		{
			name: "Error tool result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_456",
				Content:   "Error occurred",
				IsError:   true,
			},
			lineNum:     2,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT ERROR",
				"Error occurred",
			},
			expectedExcludes: []string{
				"Tool ID:",
			},
		},
		{
			name: "Tool result with verbose ID",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_verbose_123",
				Content:   "Output text",
				IsError:   false,
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: true,
			expectedIncludes: []string{
				"TOOL RESULT",
				"Tool ID: tool_verbose_123",
				"Output text",
			},
		},
		{
			name: "Empty tool result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_empty",
				Content:   "",
				IsError:   false,
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT",
				"(no output)",
			},
		},
		{
			name: "Tool result with line number",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_789",
				Content:   "Some output",
				IsError:   false,
			},
			lineNum:     55,
			showLineNum: true,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT",
				"(line 55)",
				"Some output",
			},
		},
		{
			name: "Tool result with long output",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_long",
				Content:   strings.Repeat("line\n", 50),
				IsError:   false,
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT",
				"more lines",
			},
		},
		{
			name: "Tool result with system reminder (non-verbose)",
			block: &ContentBlock{
				Type:    "tool_result",
				Content: "Content before\n<system-reminder>This should be stripped</system-reminder>\nContent after",
				IsError: false,
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"Content before",
				"Content after",
			},
			expectedExcludes: []string{
				"system-reminder",
				"This should be stripped",
			},
		},
		{
			name: "Tool result with non-string content",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_num",
				Content:   123,
				IsError:   false,
			},
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"TOOL RESULT",
				"123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origShowLineNum := *showLineNum
			origVerbose := *verbose
			*showLineNum = tt.showLineNum
			*verbose = tt.verboseMode
			defer func() {
				*showLineNum = origShowLineNum
				*verbose = origVerbose
			}()

			output := captureStdout(func() {
				displayToolResult(tt.block, tt.lineNum)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displayToolResult() output missing %q\nGot:\n%s", expected, output)
				}
			}

			for _, excluded := range tt.expectedExcludes {
				if strings.Contains(output, excluded) {
					t.Errorf("displayToolResult() output should not contain %q\nGot:\n%s", excluded, output)
				}
			}
		})
	}
}

// TestDisplayResultMessage tests the default style displayResultMessage function
func TestDisplayResultMessage(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		msg              *StreamMessage
		lineNum          int
		showLineNum      bool
		verboseMode      bool
		expectedIncludes []string
		expectedExcludes []string
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
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
				"RESULT: SUCCESS",
				"Turns: 5",
				"Duration: 15.00s",
				"Cost: $0.0025",
				"Tokens: in=1000 out=500",
			},
			expectedExcludes: []string{"ERROR"},
		},
		{
			name: "Error result",
			msg: &StreamMessage{
				Type:    "result",
				IsError: true,
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{"RESULT: ERROR"},
		},
		{
			name: "Result with line number",
			msg: &StreamMessage{
				Type:    "result",
				IsError: false,
			},
			lineNum:          47,
			showLineNum:      true,
			verboseMode:      false,
			expectedIncludes: []string{"RESULT: SUCCESS", "(line 47)"},
		},
		{
			name: "Result with API duration",
			msg: &StreamMessage{
				Type:          "result",
				IsError:       false,
				DurationMS:    10000,
				DurationAPIMS: 8000,
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{"Duration: 10.00s", "(API: 8.00s)"},
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
			lineNum:     1,
			showLineNum: false,
			verboseMode: false,
			expectedIncludes: []string{
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
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{"Task completed successfully"},
		},
		{
			name: "Result with permission denials",
			msg: &StreamMessage{
				Type:              "result",
				IsError:           false,
				PermissionDenials: []interface{}{"Denied action 1", "Denied action 2"},
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{"Permission Denials: 2"},
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
			lineNum:     1,
			showLineNum: false,
			verboseMode: true,
			expectedIncludes: []string{
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
			origShowLineNum := *showLineNum
			origVerbose := *verbose
			*showLineNum = tt.showLineNum
			*verbose = tt.verboseMode
			defer func() {
				*showLineNum = origShowLineNum
				*verbose = origVerbose
			}()

			output := captureStdout(func() {
				displayResultMessage(tt.msg, tt.lineNum)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displayResultMessage() output missing %q\nGot:\n%s", expected, output)
				}
			}

			for _, excluded := range tt.expectedExcludes {
				if strings.Contains(output, excluded) {
					t.Errorf("displayResultMessage() output should not contain %q\nGot:\n%s", excluded, output)
				}
			}
		})
	}
}

// TestDisplayTodos tests the default style displayTodos function
func TestDisplayTodos(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		todos            []interface{}
		expectedIncludes []string
		expectedExcludes []string
	}{
		{
			name: "Mixed status todos",
			todos: []interface{}{
				map[string]interface{}{"content": "Completed task", "status": "completed"},
				map[string]interface{}{"content": "In progress task", "status": "in_progress"},
				map[string]interface{}{"content": "Pending task", "status": "pending"},
			},
			expectedIncludes: []string{
				"Completed task",
				"In progress task",
				"Pending task",
			},
		},
		{
			name: "Unknown status",
			todos: []interface{}{
				map[string]interface{}{"content": "Unknown status task", "status": "unknown"},
			},
			expectedIncludes: []string{"Unknown status task"},
		},
		{
			name:             "Empty todos",
			todos:            []interface{}{},
			expectedIncludes: []string{},
		},
		{
			name: "Invalid todo format (skipped)",
			todos: []interface{}{
				"not a map",
				123,
			},
			expectedIncludes: []string{},
		},
		{
			name: "Todo without status defaults to dash",
			todos: []interface{}{
				map[string]interface{}{"content": "No status task"},
			},
			expectedIncludes: []string{"No status task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				displayTodos(tt.todos)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displayTodos() output missing %q\nGot:\n%s", expected, output)
				}
			}

			for _, excluded := range tt.expectedExcludes {
				if strings.Contains(output, excluded) {
					t.Errorf("displayTodos() output should not contain %q\nGot:\n%s", excluded, output)
				}
			}
		})
	}
}

// TestDisplayUsage tests the default style displayUsage function
func TestDisplayUsage(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		usage            *Usage
		expectedIncludes []string
	}{
		{
			name: "Basic usage with input and output tokens",
			usage: &Usage{
				InputTokens:  100,
				OutputTokens: 50,
			},
			expectedIncludes: []string{"Tokens: in=100 out=50"},
		},
		{
			name: "Usage with cache read tokens",
			usage: &Usage{
				InputTokens:          200,
				OutputTokens:         75,
				CacheReadInputTokens: 150,
			},
			expectedIncludes: []string{"Tokens: in=200 out=75", "cache_read=150"},
		},
		{
			name: "Usage with cache creation tokens",
			usage: &Usage{
				InputTokens:              300,
				OutputTokens:             100,
				CacheCreationInputTokens: 250,
			},
			expectedIncludes: []string{"Tokens: in=300 out=100", "cache_create=250"},
		},
		{
			name: "Usage with both cache types",
			usage: &Usage{
				InputTokens:              500,
				OutputTokens:             200,
				CacheReadInputTokens:     300,
				CacheCreationInputTokens: 150,
			},
			expectedIncludes: []string{
				"Tokens: in=500 out=200",
				"cache_read=300",
				"cache_create=150",
			},
		},
		{
			name: "Zero tokens",
			usage: &Usage{
				InputTokens:  0,
				OutputTokens: 0,
			},
			expectedIncludes: []string{"Tokens: in=0 out=0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				displayUsage(tt.usage)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displayUsage() output missing %q\nGot:\n%s", expected, output)
				}
			}
		})
	}
}
