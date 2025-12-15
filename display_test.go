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
