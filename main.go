package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Message types from Claude stream
type StreamMessage struct {
	Type             string          `json:"type"`
	Subtype          string          `json:"subtype,omitempty"`
	Message          *MessageContent `json:"message,omitempty"`
	SessionID        string          `json:"session_id,omitempty"`
	ParentToolUseID  string          `json:"parent_tool_use_id,omitempty"`
	CWD              string          `json:"cwd,omitempty"`
	Tools            []string        `json:"tools,omitempty"`
	Model            string          `json:"model,omitempty"`
	ClaudeCodeVersion string         `json:"claude_code_version,omitempty"`
}

type MessageContent struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Model        string         `json:"model"`
	Content      []ContentBlock `json:"content"`
	StopReason   *string        `json:"stop_reason"`
	StopSequence *string        `json:"stop_sequence"`
	Usage        *Usage         `json:"usage"`
}

type ContentBlock struct {
	Type       string                 `json:"type"`
	Text       string                 `json:"text,omitempty"`
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Input      map[string]interface{} `json:"input,omitempty"`
	ToolUseID  string                 `json:"tool_use_id,omitempty"`
	Content    interface{}            `json:"content,omitempty"`
	IsError    bool                   `json:"is_error,omitempty"`
}

type Usage struct {
	InputTokens               int                  `json:"input_tokens"`
	OutputTokens              int                  `json:"output_tokens"`
	CacheCreationInputTokens  int                  `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens      int                  `json:"cache_read_input_tokens,omitempty"`
	CacheCreation             *CacheCreationDetail `json:"cache_creation,omitempty"`
	ServiceTier               string               `json:"service_tier,omitempty"`
}

type CacheCreationDetail struct {
	Ephemeral5mInputTokens  int `json:"ephemeral_5m_input_tokens"`
	Ephemeral1hInputTokens  int `json:"ephemeral_1h_input_tokens"`
}

// Color definitions
var (
	boldCyan    = color.New(color.FgCyan, color.Bold)
	boldGreen   = color.New(color.FgGreen, color.Bold)
	boldYellow  = color.New(color.FgYellow, color.Bold)
	boldRed     = color.New(color.FgRed, color.Bold)
	boldMagenta = color.New(color.FgMagenta, color.Bold)
	cyan        = color.New(color.FgCyan)
	green       = color.New(color.FgGreen)
	yellow      = color.New(color.FgYellow)
	red         = color.New(color.FgRed)
	gray        = color.New(color.FgHiBlack)
	white       = color.New(color.FgWhite)
)

func main() {
	var reader io.Reader

	// Check if we have a file argument or should read from stdin
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse the JSON
		var msg StreamMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			gray.Printf("│ [Line %d] Invalid JSON: %v\n", lineNum, err)
			continue
		}

		// Format and display the message
		displayMessage(&msg, lineNum)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func displayMessage(msg *StreamMessage, lineNum int) {
	switch msg.Type {
	case "system":
		displaySystemMessage(msg, lineNum)
	case "assistant":
		displayAssistantMessage(msg, lineNum)
	case "user":
		displayUserMessage(msg, lineNum)
	default:
		gray.Printf("│ [Line %d] Unknown message type: %s\n", lineNum, msg.Type)
	}
}

func displaySystemMessage(msg *StreamMessage, lineNum int) {
	boldCyan.Print("┌─ ")
	boldCyan.Print("SYSTEM")
	if msg.Subtype != "" {
		cyan.Printf(" [%s]", msg.Subtype)
	}
	cyan.Printf(" (line %d)\n", lineNum)

	if msg.CWD != "" {
		cyan.Printf("│ Working Directory: %s\n", msg.CWD)
	}
	if msg.Model != "" {
		cyan.Printf("│ Model: %s\n", msg.Model)
	}
	if msg.ClaudeCodeVersion != "" {
		cyan.Printf("│ Claude Code: v%s\n", msg.ClaudeCodeVersion)
	}
	if len(msg.Tools) > 0 {
		cyan.Printf("│ Tools: %d available\n", len(msg.Tools))
	}

	cyan.Println("└─")
}

func displayAssistantMessage(msg *StreamMessage, lineNum int) {
	if msg.Message == nil {
		return
	}

	content := msg.Message.Content
	if len(content) == 0 {
		return
	}

	// Group consecutive text blocks together
	var textBlocks []string
	var toolUses []ContentBlock

	for _, block := range content {
		switch block.Type {
		case "text":
			if block.Text != "" {
				textBlocks = append(textBlocks, block.Text)
			}
		case "tool_use":
			toolUses = append(toolUses, block)
		}
	}

	// Display text blocks
	if len(textBlocks) > 0 {
		boldGreen.Print("┌─ ")
		boldGreen.Print("ASSISTANT")
		gray.Printf(" (line %d)\n", lineNum)

		for _, text := range textBlocks {
			green.Print("│ ")
			white.Println(text)
		}

		if msg.Message.Usage != nil {
			displayUsage(msg.Message.Usage)
		}
		green.Println("└─")
	}

	// Display tool uses
	for _, tool := range toolUses {
		displayToolUse(&tool, lineNum)
	}
}

func displayToolUse(tool *ContentBlock, lineNum int) {
	boldYellow.Print("┌─ ")
	boldYellow.Printf("TOOL: %s", tool.Name)
	gray.Printf(" (line %d)\n", lineNum)
	yellow.Printf("│ ID: %s\n", tool.ID)

	if tool.Input != nil {
		yellow.Println("│ Input:")
		for key, value := range tool.Input {
			// Pretty print the value
			yellow.Printf("│   %s: ", key)

			switch v := value.(type) {
			case string:
				// Truncate long strings
				if len(v) > 100 {
					white.Printf("%s...\n", v[:100])
				} else {
					white.Println(v)
				}
			case []interface{}:
				white.Printf("[%d items]\n", len(v))
			case map[string]interface{}:
				white.Println("{...}")
			default:
				white.Printf("%v\n", v)
			}
		}
	}

	yellow.Println("└─")
}

func displayUserMessage(msg *StreamMessage, lineNum int) {
	if msg.Message == nil {
		return
	}

	content := msg.Message.Content
	if len(content) == 0 {
		return
	}

	for _, block := range content {
		if block.Type == "tool_result" {
			displayToolResult(&block, lineNum)
		}
	}
}

func displayToolResult(block *ContentBlock, lineNum int) {
	if block.IsError {
		boldRed.Print("┌─ ")
		boldRed.Print("TOOL RESULT ERROR")
		gray.Printf(" (line %d)\n", lineNum)
		red.Printf("│ Tool ID: %s\n", block.ToolUseID)

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		red.Print("│ ")
		white.Println(contentStr)
		red.Println("└─")
	} else {
		boldMagenta.Print("┌─ ")
		boldMagenta.Print("TOOL RESULT")
		gray.Printf(" (line %d)\n", lineNum)
		gray.Printf("│ Tool ID: %s\n", block.ToolUseID)

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		if contentStr == "" {
			gray.Println("│ (no output)")
		} else {
			// Truncate very long output
			lines := strings.Split(contentStr, "\n")
			maxLines := 10

			for i, line := range lines {
				if i >= maxLines {
					gray.Printf("│ ... (%d more lines)\n", len(lines)-maxLines)
					break
				}
				gray.Print("│ ")
				white.Println(line)
			}
		}

		gray.Println("└─")
	}
}

func displayUsage(usage *Usage) {
	gray.Print("│ ")
	gray.Printf("Tokens: in=%d out=%d", usage.InputTokens, usage.OutputTokens)

	if usage.CacheReadInputTokens > 0 {
		gray.Printf(" cache_read=%d", usage.CacheReadInputTokens)
	}
	if usage.CacheCreationInputTokens > 0 {
		gray.Printf(" cache_create=%d", usage.CacheCreationInputTokens)
	}

	gray.Println()
}
