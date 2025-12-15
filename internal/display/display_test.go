package display

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ariel-frischer/claude-clean/internal/parser"
	"github.com/fatih/color"
)

// captureStdout captures stdout during function execution
func captureStdout(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set color package output to our pipe
	color.Output = w

	// Capture output in a goroutine
	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		io.Copy(&buf, r)
		close(done)
	}()

	// Run the function
	f()

	// Close writer and restore stdout
	w.Close()
	os.Stdout = oldStdout
	color.Output = oldStdout

	// Wait for reader to finish
	<-done

	return buf.String()
}

// stripANSI removes ANSI escape codes from a string for easier testing
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			i++ // Skip the '['
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

func TestFormatLineNum(t *testing.T) {
	tests := []struct {
		name        string
		lineNum     int
		showLineNum bool
		expected    string
	}{
		{
			name:        "Single digit line number with showLineNum enabled",
			lineNum:     5,
			showLineNum: true,
			expected:    " (line 5)",
		},
		{
			name:        "Double digit line number with showLineNum enabled",
			lineNum:     42,
			showLineNum: true,
			expected:    " (line 42)",
		},
		{
			name:        "showLineNum disabled returns empty string",
			lineNum:     100,
			showLineNum: false,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatLineNum(tt.lineNum, tt.showLineNum)
			if result != tt.expected {
				t.Errorf("FormatLineNum(%d, %v) = %q, want %q",
					tt.lineNum, tt.showLineNum, result, tt.expected)
			}
		})
	}
}

func TestFormatLineNumCompact(t *testing.T) {
	tests := []struct {
		name        string
		lineNum     int
		showLineNum bool
		expected    string
	}{
		{
			name:        "Single digit line number with showLineNum enabled",
			lineNum:     5,
			showLineNum: true,
			expected:    " L5",
		},
		{
			name:        "Double digit line number with showLineNum enabled",
			lineNum:     42,
			showLineNum: true,
			expected:    " L42",
		},
		{
			name:        "showLineNum disabled returns empty string",
			lineNum:     100,
			showLineNum: false,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatLineNumCompact(tt.lineNum, tt.showLineNum)
			if result != tt.expected {
				t.Errorf("FormatLineNumCompact(%d, %v) = %q, want %q",
					tt.lineNum, tt.showLineNum, result, tt.expected)
			}
		})
	}
}

func TestDisplayUsageInline(t *testing.T) {
	tests := []struct {
		name     string
		usage    parser.Usage
		expected string
	}{
		{
			name: "Basic usage with input and output tokens",
			usage: parser.Usage{
				InputTokens:  100,
				OutputTokens: 50,
			},
			expected: "│ Tokens: in=100 out=50\n",
		},
		{
			name: "Usage with cache read tokens",
			usage: parser.Usage{
				InputTokens:          200,
				OutputTokens:         75,
				CacheReadInputTokens: 150,
			},
			expected: "│ Tokens: in=200 out=75 cache_read=150\n",
		},
		{
			name: "Usage with cache creation tokens",
			usage: parser.Usage{
				InputTokens:              300,
				OutputTokens:             100,
				CacheCreationInputTokens: 250,
			},
			expected: "│ Tokens: in=300 out=100 cache_create=250\n",
		},
		{
			name: "Usage with both cache types",
			usage: parser.Usage{
				InputTokens:              500,
				OutputTokens:             200,
				CacheReadInputTokens:     300,
				CacheCreationInputTokens: 150,
			},
			expected: "│ Tokens: in=500 out=200 cache_read=300 cache_create=150\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			color.NoColor = true
			oldOutput := color.Output
			color.Output = &buf
			defer func() {
				color.Output = oldOutput
			}()

			c := color.New()
			DisplayUsageInline(&tt.usage, c)

			result := buf.String()
			if result != tt.expected {
				t.Errorf("DisplayUsageInline() output =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

// TestDisplaySystemMessageDefault tests the default style displaySystemMessage function
func TestDisplaySystemMessageDefault(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		msg              *parser.StreamMessage
		lineNum          int
		showLineNum      bool
		expectedIncludes []string
	}{
		{
			name: "Basic system message",
			msg: &parser.StreamMessage{
				Type: "system",
			},
			lineNum:     1,
			showLineNum: false,
			expectedIncludes: []string{
				"SYSTEM",
				"└─",
			},
		},
		{
			name: "System message with subtype",
			msg: &parser.StreamMessage{
				Type:    "system",
				Subtype: "init",
			},
			lineNum:     1,
			showLineNum: false,
			expectedIncludes: []string{
				"SYSTEM",
				"[init]",
			},
		},
		{
			name: "System message with all fields",
			msg: &parser.StreamMessage{
				Type:              "system",
				Subtype:           "startup",
				CWD:               "/home/user/project",
				Model:             "claude-opus-4-5",
				ClaudeCodeVersion: "1.2.3",
				Tools:             []string{"bash", "read", "write"},
			},
			lineNum:     5,
			showLineNum: true,
			expectedIncludes: []string{
				"SYSTEM",
				"[startup]",
				"(line 5)",
				"Working Directory: /home/user/project",
				"Model: claude-opus-4-5",
				"Claude Code: v1.2.3",
				"Tools: 3 available",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Style:       StyleDefault,
				Verbose:     false,
				ShowLineNum: tt.showLineNum,
			}

			output := captureStdout(func() {
				displaySystemMessage(tt.msg, tt.lineNum, cfg)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displaySystemMessage() output missing %q\nGot:\n%s", expected, output)
				}
			}
		})
	}
}

// TestDisplayAssistantMessageDefault tests the default style displayAssistantMessage function
func TestDisplayAssistantMessageDefault(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		msg              *parser.StreamMessage
		lineNum          int
		showLineNum      bool
		verboseMode      bool
		expectedIncludes []string
		expectedExcludes []string
	}{
		{
			name: "Nil message",
			msg: &parser.StreamMessage{
				Type:    "assistant",
				Message: nil,
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{},
		},
		{
			name: "Simple text message",
			msg: &parser.StreamMessage{
				Type: "assistant",
				Message: &parser.MessageContent{
					Content: []parser.ContentBlock{
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
			name: "Text with usage in verbose mode",
			msg: &parser.StreamMessage{
				Type: "assistant",
				Message: &parser.MessageContent{
					Content: []parser.ContentBlock{
						{
							Type: "text",
							Text: "Response text",
						},
					},
					Usage: &parser.Usage{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Style:       StyleDefault,
				Verbose:     tt.verboseMode,
				ShowLineNum: tt.showLineNum,
			}

			output := captureStdout(func() {
				displayAssistantMessage(tt.msg, tt.lineNum, cfg)
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

// TestDisplayTodos tests the DisplayTodos function
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				DisplayTodos(tt.todos)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("DisplayTodos() output missing %q\nGot:\n%s", expected, output)
				}
			}

			for _, excluded := range tt.expectedExcludes {
				if strings.Contains(output, excluded) {
					t.Errorf("DisplayTodos() output should not contain %q\nGot:\n%s", excluded, output)
				}
			}
		})
	}
}

// TestDisplayResultMessageDefault tests the default style displayResultMessage function
func TestDisplayResultMessageDefault(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		msg              *parser.StreamMessage
		lineNum          int
		showLineNum      bool
		verboseMode      bool
		expectedIncludes []string
		expectedExcludes []string
	}{
		{
			name: "Success result",
			msg: &parser.StreamMessage{
				Type:         "result",
				IsError:      false,
				NumTurns:     5,
				DurationMS:   15000,
				TotalCostUSD: 0.0025,
				Usage: &parser.Usage{
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
			msg: &parser.StreamMessage{
				Type:    "result",
				IsError: true,
			},
			lineNum:          1,
			showLineNum:      false,
			verboseMode:      false,
			expectedIncludes: []string{"RESULT: ERROR"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Style:       StyleDefault,
				Verbose:     tt.verboseMode,
				ShowLineNum: tt.showLineNum,
			}

			output := captureStdout(func() {
				displayResultMessage(tt.msg, tt.lineNum, cfg)
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
