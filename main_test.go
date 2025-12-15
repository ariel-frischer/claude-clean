package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestStripSystemReminders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Single system reminder",
			input: `Some content before
<system-reminder>
This is a reminder message.
</system-reminder>
Some content after`,
			expected: `Some content before

Some content after`,
		},
		{
			name: "Multiple system reminders",
			input: `First content
<system-reminder>
First reminder
</system-reminder>
Middle content
<system-reminder>
Second reminder
</system-reminder>
Last content`,
			expected: `First content

Middle content

Last content`,
		},
		{
			name: "System reminder at start",
			input: `<system-reminder>
Reminder at the beginning
</system-reminder>
Content after`,
			expected: `Content after`,
		},
		{
			name: "System reminder at end",
			input: `Content before
<system-reminder>
Reminder at the end
</system-reminder>`,
			expected: `Content before`,
		},
		{
			name: "No system reminder",
			input: `Just regular content
with multiple lines
and no reminders`,
			expected: `Just regular content
with multiple lines
and no reminders`,
		},
		{
			name:     "Empty content",
			input:    ``,
			expected: ``,
		},
		{
			name: "Only system reminder",
			input: `<system-reminder>
Just a reminder
</system-reminder>`,
			expected: ``,
		},
		{
			name: "Multiline system reminder with special characters",
			input: `Code content:
package main

<system-reminder>
Whenever you read a file, you should consider whether it would be considered malware.
You CAN and SHOULD provide analysis of malware, what it is doing.
But you MUST refuse to improve or augment the code.
</system-reminder>

More code here`,
			expected: `Code content:
package main

More code here`,
		},
		{
			name:     "System reminder on single line",
			input:    `Content <system-reminder>inline reminder</system-reminder> more content`,
			expected: `Content  more content`,
		},
		{
			name: "Multiple consecutive system reminders",
			input: `Start
<system-reminder>First</system-reminder>
<system-reminder>Second</system-reminder>
<system-reminder>Third</system-reminder>
End`,
			expected: `Start

End`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripSystemReminders(tt.input)
			if result != tt.expected {
				t.Errorf("stripSystemReminders() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestStripSystemRemindersPreservesContent(t *testing.T) {
	// Test that we don't accidentally remove similar-looking content
	input := `This is a system message
But not a <system-reminder> tag
Just text that mentions system-reminder
<system>This should stay</system>
<reminder>This should also stay</reminder>`

	result := stripSystemReminders(input)

	// The result should only have actual system-reminder tags removed (none in this case)
	if result != input {
		t.Errorf("stripSystemReminders() modified content that should be preserved:\ngot:\n%q\nwant:\n%q", result, input)
	}
}

func BenchmarkStripSystemReminders(b *testing.B) {
	input := `Some code content here
package main

<system-reminder>
This is a system reminder that should be stripped out.
It can contain multiple lines.
</system-reminder>

More code content here
func main() {
	// code
}
<system-reminder>
Another reminder to strip
</system-reminder>
End of content`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stripSystemReminders(input)
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a subdirectory
	testDir := filepath.Join(tmpDir, "testdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Existing file",
			path:     testFile,
			expected: true,
		},
		{
			name:     "Non-existing file",
			path:     filepath.Join(tmpDir, "nonexistent.txt"),
			expected: false,
		},
		{
			name:     "Directory (not a file)",
			path:     testDir,
			expected: false,
		},
		{
			name:     "Empty string",
			path:     "",
			expected: false,
		},
		{
			name:     "Relative path to non-existing file",
			path:     "does/not/exist.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileExists(tt.path)
			if result != tt.expected {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestBinaryName(t *testing.T) {
	// Save the original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Simple binary name",
			args:     []string{"claude-clean"},
			expected: "claude-clean",
		},
		{
			name:     "Binary with absolute path",
			args:     []string{"/usr/local/bin/claude-clean"},
			expected: "claude-clean",
		},
		{
			name:     "Binary with relative path",
			args:     []string{"./bin/claude-clean"},
			expected: "claude-clean",
		},
		{
			name:     "Binary with nested path",
			args:     []string{"/home/user/go/bin/claude-clean"},
			expected: "claude-clean",
		},
		{
			name:     "Binary name with extension",
			args:     []string{"/home/user/go/bin/claude-clean.exe"},
			expected: "claude-clean.exe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set os.Args to the test args
			os.Args = tt.args

			result := binaryName()
			if result != tt.expected {
				t.Errorf("binaryName() with os.Args=%v = %q, want %q", tt.args, result, tt.expected)
			}
		})
	}
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
			name:        "Triple digit line number with showLineNum enabled",
			lineNum:     999,
			showLineNum: true,
			expected:    " (line 999)",
		},
		{
			name:        "Large line number with showLineNum enabled",
			lineNum:     123456,
			showLineNum: true,
			expected:    " (line 123456)",
		},
		{
			name:        "Line 1 with showLineNum enabled",
			lineNum:     1,
			showLineNum: true,
			expected:    " (line 1)",
		},
		{
			name:        "showLineNum disabled returns empty string",
			lineNum:     100,
			showLineNum: false,
			expected:    "",
		},
		{
			name:        "showLineNum disabled with line 1",
			lineNum:     1,
			showLineNum: false,
			expected:    "",
		},
		{
			name:        "Zero line number with showLineNum enabled",
			lineNum:     0,
			showLineNum: true,
			expected:    " (line 0)",
		},
		{
			name:        "Zero line number with showLineNum disabled",
			lineNum:     0,
			showLineNum: false,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalShowLineNum := *showLineNum
			// Restore after test
			defer func() { *showLineNum = originalShowLineNum }()

			// Set the flag for this test
			*showLineNum = tt.showLineNum

			result := formatLineNum(tt.lineNum)
			if result != tt.expected {
				t.Errorf("formatLineNum(%d) with showLineNum=%v = %q, want %q",
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
			name:        "Triple digit line number with showLineNum enabled",
			lineNum:     999,
			showLineNum: true,
			expected:    " L999",
		},
		{
			name:        "Large line number with showLineNum enabled",
			lineNum:     123456,
			showLineNum: true,
			expected:    " L123456",
		},
		{
			name:        "Line 1 with showLineNum enabled",
			lineNum:     1,
			showLineNum: true,
			expected:    " L1",
		},
		{
			name:        "showLineNum disabled returns empty string",
			lineNum:     100,
			showLineNum: false,
			expected:    "",
		},
		{
			name:        "showLineNum disabled with line 1",
			lineNum:     1,
			showLineNum: false,
			expected:    "",
		},
		{
			name:        "Zero line number with showLineNum enabled",
			lineNum:     0,
			showLineNum: true,
			expected:    " L0",
		},
		{
			name:        "Zero line number with showLineNum disabled",
			lineNum:     0,
			showLineNum: false,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalShowLineNum := *showLineNum
			// Restore after test
			defer func() { *showLineNum = originalShowLineNum }()

			// Set the flag for this test
			*showLineNum = tt.showLineNum

			result := formatLineNumCompact(tt.lineNum)
			if result != tt.expected {
				t.Errorf("formatLineNumCompact(%d) with showLineNum=%v = %q, want %q",
					tt.lineNum, tt.showLineNum, result, tt.expected)
			}
		})
	}
}

func TestDisplayUsageInline(t *testing.T) {
	tests := []struct {
		name     string
		usage    Usage
		expected string
	}{
		{
			name: "Basic usage with input and output tokens",
			usage: Usage{
				InputTokens:  100,
				OutputTokens: 50,
			},
			expected: "│ Tokens: in=100 out=50\n",
		},
		{
			name: "Usage with cache read tokens",
			usage: Usage{
				InputTokens:          200,
				OutputTokens:         75,
				CacheReadInputTokens: 150,
			},
			expected: "│ Tokens: in=200 out=75 cache_read=150\n",
		},
		{
			name: "Usage with cache creation tokens",
			usage: Usage{
				InputTokens:              300,
				OutputTokens:             100,
				CacheCreationInputTokens: 250,
			},
			expected: "│ Tokens: in=300 out=100 cache_create=250\n",
		},
		{
			name: "Usage with both cache types",
			usage: Usage{
				InputTokens:              500,
				OutputTokens:             200,
				CacheReadInputTokens:     300,
				CacheCreationInputTokens: 150,
			},
			expected: "│ Tokens: in=500 out=200 cache_read=300 cache_create=150\n",
		},
		{
			name: "Zero tokens",
			usage: Usage{
				InputTokens:  0,
				OutputTokens: 0,
			},
			expected: "│ Tokens: in=0 out=0\n",
		},
		{
			name: "Large token counts",
			usage: Usage{
				InputTokens:  1000000,
				OutputTokens: 500000,
			},
			expected: "│ Tokens: in=1000000 out=500000\n",
		},
		{
			name: "Only cache read tokens (no creation)",
			usage: Usage{
				InputTokens:              1000,
				OutputTokens:             500,
				CacheReadInputTokens:     800,
				CacheCreationInputTokens: 0,
			},
			expected: "│ Tokens: in=1000 out=500 cache_read=800\n",
		},
		{
			name: "Only cache creation tokens (no read)",
			usage: Usage{
				InputTokens:              1000,
				OutputTokens:             500,
				CacheReadInputTokens:     0,
				CacheCreationInputTokens: 600,
			},
			expected: "│ Tokens: in=1000 out=500 cache_create=600\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since displayUsageInline prints to a color.Color, we need to capture output
			// We'll use a bytes.Buffer and create a no-color printer
			var buf bytes.Buffer

			// Temporarily disable color output and redirect to buffer
			color.NoColor = true
			oldOutput := color.Output
			color.Output = &buf
			defer func() {
				color.Output = oldOutput
			}()

			// Create a color instance
			c := color.New()

			displayUsageInline(&tt.usage, c)

			result := buf.String()
			if result != tt.expected {
				t.Errorf("displayUsageInline() output =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

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
	// Simple ANSI escape sequence removal
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

func TestDisplaySystemMessageCompact(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		msg      *StreamMessage
		lineNum  int
		showLine bool
		contains []string
	}{
		{
			name: "Basic system message",
			msg: &StreamMessage{
				Type: "system",
			},
			lineNum:  1,
			showLine: false,
			contains: []string{"SYS"},
		},
		{
			name: "System message with subtype",
			msg: &StreamMessage{
				Type:    "system",
				Subtype: "init",
			},
			lineNum:  5,
			showLine: false,
			contains: []string{"SYS", "[init]"},
		},
		{
			name: "System message with model",
			msg: &StreamMessage{
				Type:  "system",
				Model: "claude-sonnet-4-5",
			},
			lineNum:  10,
			showLine: false,
			contains: []string{"SYS", "claude-sonnet-4-5"},
		},
		{
			name: "System message with CWD",
			msg: &StreamMessage{
				Type: "system",
				CWD:  "/home/user/project",
			},
			lineNum:  15,
			showLine: false,
			contains: []string{"SYS", "@/home/user/project"},
		},
		{
			name: "System message with all fields",
			msg: &StreamMessage{
				Type:    "system",
				Subtype: "start",
				Model:   "claude-opus-4-5",
				CWD:     "/tmp",
			},
			lineNum:  20,
			showLine: true,
			contains: []string{"SYS", "[start]", "L20", "claude-opus-4-5", "@/tmp"},
		},
		{
			name: "System message with line numbers",
			msg: &StreamMessage{
				Type: "system",
			},
			lineNum:  42,
			showLine: true,
			contains: []string{"SYS", "L42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set showLineNum flag
			oldShowLineNum := *showLineNum
			*showLineNum = tt.showLine
			defer func() { *showLineNum = oldShowLineNum }()

			output := captureStdout(func() {
				displaySystemMessageCompact(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)
			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("displaySystemMessageCompact() output missing %q\nGot: %q", expected, cleaned)
				}
			}
		})
	}
}

func TestDisplayAssistantMessageCompact(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		msg      *StreamMessage
		lineNum  int
		showLine bool
		contains []string
		notEmpty bool
	}{
		{
			name: "Nil message",
			msg: &StreamMessage{
				Type:    "assistant",
				Message: nil,
			},
			lineNum:  1,
			showLine: false,
			contains: []string{},
			notEmpty: false,
		},
		{
			name: "Empty content",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{},
				},
			},
			lineNum:  1,
			showLine: false,
			contains: []string{},
			notEmpty: false,
		},
		{
			name: "Short text block",
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
			contains: []string{"AST", "Hello, world!"},
			notEmpty: true,
		},
		{
			name: "Long text block (truncated)",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: strings.Repeat("a", 150)},
					},
				},
			},
			lineNum:  5,
			showLine: false,
			contains: []string{"AST", strings.Repeat("a", 100), "..."},
			notEmpty: true,
		},
		{
			name: "Text with newlines",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{Type: "text", Text: "Line 1\nLine 2\nLine 3"},
					},
				},
			},
			lineNum:  10,
			showLine: false,
			contains: []string{"AST", "Line 1 Line 2 Line 3"},
			notEmpty: true,
		},
		{
			name: "With line numbers",
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
			contains: []string{"AST", "L25", "Test message"},
			notEmpty: true,
		},
		{
			name: "Tool use in content",
			msg: &StreamMessage{
				Type: "assistant",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type:  "tool_use",
							Name:  "Read",
							ID:    "tool_123",
							Input: map[string]interface{}{"file_path": "/home/user/file.txt"},
						},
					},
				},
			},
			lineNum:  30,
			showLine: false,
			contains: []string{"TOOL", "Read", "file_path"},
			notEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			*showLineNum = tt.showLine
			defer func() { *showLineNum = oldShowLineNum }()

			output := captureStdout(func() {
				displayAssistantMessageCompact(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)

			if !tt.notEmpty && cleaned != "" {
				t.Errorf("displayAssistantMessageCompact() expected empty output, got: %q", cleaned)
			}

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("displayAssistantMessageCompact() output missing %q\nGot: %q", expected, cleaned)
				}
			}
		})
	}
}

func TestDisplayUserMessageCompact(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		msg      *StreamMessage
		lineNum  int
		showLine bool
		contains []string
		notEmpty bool
	}{
		{
			name: "Nil message",
			msg: &StreamMessage{
				Type:    "user",
				Message: nil,
			},
			lineNum:  1,
			showLine: false,
			contains: []string{},
			notEmpty: false,
		},
		{
			name: "Empty content",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{},
				},
			},
			lineNum:  1,
			showLine: false,
			contains: []string{},
			notEmpty: false,
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
							Content:   "Operation successful",
							IsError:   false,
						},
					},
				},
			},
			lineNum:  5,
			showLine: false,
			contains: []string{"RES", "Operation successful"},
			notEmpty: true,
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
							Content:   "File not found",
							IsError:   true,
						},
					},
				},
			},
			lineNum:  10,
			showLine: false,
			contains: []string{"ERR", "File not found"},
			notEmpty: true,
		},
		{
			name: "With line numbers",
			msg: &StreamMessage{
				Type: "user",
				Message: &MessageContent{
					Content: []ContentBlock{
						{
							Type:      "tool_result",
							ToolUseID: "tool_789",
							Content:   "Success",
							IsError:   false,
						},
					},
				},
			},
			lineNum:  42,
			showLine: true,
			contains: []string{"RES", "L42", "Success"},
			notEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			*showLineNum = tt.showLine
			defer func() { *showLineNum = oldShowLineNum }()

			output := captureStdout(func() {
				displayUserMessageCompact(tt.msg, tt.lineNum)
			})

			cleaned := stripANSI(output)

			if !tt.notEmpty && cleaned != "" {
				t.Errorf("displayUserMessageCompact() expected empty output, got: %q", cleaned)
			}

			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("displayUserMessageCompact() output missing %q\nGot: %q", expected, cleaned)
				}
			}
		})
	}
}

func TestDisplayToolUseCompact(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		tool     *ContentBlock
		lineNum  int
		showLine bool
		contains []string
	}{
		{
			name: "Basic tool use",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "Bash",
				ID:   "tool_001",
			},
			lineNum:  1,
			showLine: false,
			contains: []string{"TOOL", "Bash"},
		},
		{
			name: "Tool with string input",
			tool: &ContentBlock{
				Type:  "tool_use",
				Name:  "Read",
				ID:    "tool_002",
				Input: map[string]interface{}{"file_path": "/home/user/test.txt"},
			},
			lineNum:  5,
			showLine: false,
			contains: []string{"TOOL", "Read", "file_path", "/home/user/test.txt"},
		},
		{
			name: "Tool with long string input (truncated)",
			tool: &ContentBlock{
				Type:  "tool_use",
				Name:  "Write",
				ID:    "tool_003",
				Input: map[string]interface{}{"content": strings.Repeat("x", 100)},
			},
			lineNum:  10,
			showLine: false,
			contains: []string{"TOOL", "Write", "content", "..."},
		},
		{
			name: "Tool with array input",
			tool: &ContentBlock{
				Type:  "tool_use",
				Name:  "TodoWrite",
				ID:    "tool_004",
				Input: map[string]interface{}{"todos": []interface{}{"task1", "task2", "task3"}},
			},
			lineNum:  15,
			showLine: false,
			contains: []string{"TOOL", "TodoWrite", "todos", "[3 items]"},
		},
		{
			name: "Tool with line numbers",
			tool: &ContentBlock{
				Type:  "tool_use",
				Name:  "Grep",
				ID:    "tool_005",
				Input: map[string]interface{}{"pattern": "test"},
			},
			lineNum:  20,
			showLine: true,
			contains: []string{"TOOL", "L20", "Grep", "pattern", "test"},
		},
		{
			name: "Tool with multiple inputs",
			tool: &ContentBlock{
				Type: "tool_use",
				Name: "Edit",
				ID:   "tool_006",
				Input: map[string]interface{}{
					"file_path":  "/path/to/file",
					"old_string": "old",
					"new_string": "new",
				},
			},
			lineNum:  25,
			showLine: false,
			contains: []string{"TOOL", "Edit", "file_path", "old_string", "new_string"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldShowLineNum := *showLineNum
			*showLineNum = tt.showLine
			defer func() { *showLineNum = oldShowLineNum }()

			output := captureStdout(func() {
				displayToolUseCompact(tt.tool, tt.lineNum)
			})

			cleaned := stripANSI(output)
			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("displayToolUseCompact() output missing %q\nGot: %q", expected, cleaned)
				}
			}
		})
	}
}

func TestDisplayToolResultCompact(t *testing.T) {
	// Disable colors for consistent testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name     string
		block    *ContentBlock
		lineNum  int
		showLine bool
		verbose  bool
		contains []string
	}{
		{
			name: "Success result short",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_001",
				Content:   "File read successfully",
				IsError:   false,
			},
			lineNum:  1,
			showLine: false,
			verbose:  false,
			contains: []string{"RES", "File read successfully"},
		},
		{
			name: "Error result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_002",
				Content:   "Permission denied",
				IsError:   true,
			},
			lineNum:  5,
			showLine: false,
			verbose:  false,
			contains: []string{"ERR", "Permission denied"},
		},
		{
			name: "Empty result",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_003",
				Content:   "",
				IsError:   false,
			},
			lineNum:  10,
			showLine: false,
			verbose:  false,
			contains: []string{"RES", "(no output)"},
		},
		{
			name: "Long result (truncated)",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_004",
				Content:   strings.Repeat("a", 150),
				IsError:   false,
			},
			lineNum:  15,
			showLine: false,
			verbose:  false,
			contains: []string{"RES", "..."},
		},
		{
			name: "Result with newlines",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_005",
				Content:   "Line 1\nLine 2\nLine 3",
				IsError:   false,
			},
			lineNum:  20,
			showLine: false,
			verbose:  false,
			contains: []string{"RES", "Line 1 Line 2 Line 3"},
		},
		{
			name: "With line numbers",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_006",
				Content:   "Success",
				IsError:   false,
			},
			lineNum:  42,
			showLine: true,
			verbose:  false,
			contains: []string{"RES", "L42", "Success"},
		},
		{
			name: "With system reminders (stripped)",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_007",
				Content:   "Content before<system-reminder>Hidden</system-reminder>Content after",
				IsError:   false,
			},
			lineNum:  25,
			showLine: false,
			verbose:  false,
			contains: []string{"RES", "Content beforeContent after"},
		},
		{
			name: "With system reminders (verbose mode)",
			block: &ContentBlock{
				Type:      "tool_result",
				ToolUseID: "tool_008",
				Content:   "Content<system-reminder>Shown</system-reminder>More",
				IsError:   false,
			},
			lineNum:  30,
			showLine: false,
			verbose:  true,
			contains: []string{"RES", "<system-reminder>", "Shown", "</system-reminder>"},
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
				displayToolResultCompact(tt.block, tt.lineNum)
			})

			cleaned := stripANSI(output)
			for _, expected := range tt.contains {
				if !strings.Contains(cleaned, expected) {
					t.Errorf("displayToolResultCompact() output missing %q\nGot: %q", expected, cleaned)
				}
			}
		})
	}
}

// TestDisplaySystemMessage tests the default style displaySystemMessage function
func TestDisplaySystemMessage(t *testing.T) {
	// Disable colors for testing
	color.NoColor = true
	defer func() { color.NoColor = false }()

	tests := []struct {
		name             string
		msg              *StreamMessage
		lineNum          int
		showLineNum      bool
		expectedIncludes []string
	}{
		{
			name: "Basic system message",
			msg: &StreamMessage{
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
			msg: &StreamMessage{
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
			msg: &StreamMessage{
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
			origShowLineNum := *showLineNum
			*showLineNum = tt.showLineNum
			defer func() { *showLineNum = origShowLineNum }()

			output := captureStdout(func() {
				displaySystemMessage(tt.msg, tt.lineNum)
			})

			for _, expected := range tt.expectedIncludes {
				if !strings.Contains(output, expected) {
					t.Errorf("displaySystemMessage() output missing %q\nGot:\n%s", expected, output)
				}
			}
		})
	}
}
