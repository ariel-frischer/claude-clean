package display

import (
	"fmt"
	"strings"

	"github.com/ariel-frischer/claude-clean/internal/parser"
)

func displayMessageDefault(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	switch msg.Type {
	case "system":
		displaySystemMessage(msg, lineNum, cfg)
	case "assistant":
		displayAssistantMessage(msg, lineNum, cfg)
	case "user":
		displayUserMessage(msg, lineNum, cfg)
	case "result":
		displayResultMessage(msg, lineNum, cfg)
	default:
		Gray.Printf("│ [Line %d] Unknown message type: %s\n", lineNum, msg.Type)
	}
}

func displaySystemMessage(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	BoldCyan.Print("┌─ ")
	BoldCyan.Print("SYSTEM")
	if msg.Subtype != "" {
		Cyan.Printf(" [%s]", msg.Subtype)
	}
	Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

	if msg.CWD != "" {
		Cyan.Printf("│ Working Directory: %s\n", msg.CWD)
	}
	if msg.Model != "" {
		Cyan.Printf("│ Model: %s\n", msg.Model)
	}
	if msg.ClaudeCodeVersion != "" {
		Cyan.Printf("│ Claude Code: v%s\n", msg.ClaudeCodeVersion)
	}
	if len(msg.Tools) > 0 {
		Cyan.Printf("│ Tools: %d available\n", len(msg.Tools))
	}

	Cyan.Println("└─")
}

func displayAssistantMessage(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.Message == nil {
		return
	}

	content := msg.Message.Content
	if len(content) == 0 {
		return
	}

	// Group consecutive text blocks together
	var textBlocks []string
	var toolUses []parser.ContentBlock

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
		BoldGreen.Print("┌─ ")
		BoldGreen.Print("ASSISTANT")
		Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		for _, text := range textBlocks {
			Green.Print("│ ")
			White.Println(text)
		}

		if cfg.Verbose && msg.Message.Usage != nil {
			DisplayUsage(msg.Message.Usage)
		}
		Green.Println("└─")
	}

	// Display tool uses
	for _, tool := range toolUses {
		displayToolUse(&tool, lineNum, cfg)
	}
}

func displayToolUse(tool *parser.ContentBlock, lineNum int, cfg *Config) {
	BoldYellow.Print("┌─ ")
	BoldYellow.Printf("TOOL: %s", tool.Name)
	Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

	if cfg.Verbose {
		Yellow.Printf("│ ID: %s\n", tool.ID)
	}

	if tool.Input != nil {
		Yellow.Println("│ Input:")
		for key, value := range tool.Input {
			// Pretty print the value
			Yellow.Printf("│   %s: ", key)

			switch v := value.(type) {
			case string:
				// Show more context for strings - first 200 + last 100 chars
				if len(v) > 300 {
					White.Printf("%s ... (%d chars omitted) ... %s\n",
						v[:200], len(v)-300, v[len(v)-100:])
				} else {
					White.Println(v)
				}
			case []interface{}:
				// Special handling for todos array in TodoWrite tool
				if tool.Name == "TodoWrite" && key == "todos" {
					White.Println()
					DisplayTodos(v)
				} else {
					White.Printf("[%d items]\n", len(v))
				}
			case map[string]interface{}:
				White.Println("{...}")
			default:
				White.Printf("%v\n", v)
			}
		}
	}

	Yellow.Println("└─")
}

func displayUserMessage(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.Message == nil {
		return
	}

	content := msg.Message.Content
	if len(content) == 0 {
		return
	}

	for _, block := range content {
		if block.Type == "tool_result" {
			displayToolResult(&block, lineNum, cfg)
		}
	}
}

func displayToolResult(block *parser.ContentBlock, lineNum int, cfg *Config) {
	if block.IsError {
		BoldRed.Print("┌─ ")
		BoldRed.Print("TOOL RESULT ERROR")
		Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		if cfg.Verbose {
			Red.Printf("│ Tool ID: %s\n", block.ToolUseID)
		}

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		// Strip system reminders in non-verbose mode
		if !cfg.Verbose {
			contentStr = parser.StripSystemReminders(contentStr)
		}

		Red.Print("│ ")
		White.Println(contentStr)
		Red.Println("└─")
	} else {
		BoldMagenta.Print("┌─ ")
		BoldMagenta.Print("TOOL RESULT")
		Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		if cfg.Verbose {
			Gray.Printf("│ Tool ID: %s\n", block.ToolUseID)
		}

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		// Strip system reminders in non-verbose mode
		if !cfg.Verbose {
			contentStr = parser.StripSystemReminders(contentStr)
		}

		if contentStr == "" {
			Gray.Println("│ (no output)")
		} else {
			// Show first 20 + last 20 lines for long output
			lines := strings.Split(contentStr, "\n")
			firstLines := parser.FirstLines
			lastLines := parser.LastLines
			totalLines := len(lines)

			if totalLines <= firstLines+lastLines {
				// Show all lines if content is short enough
				for _, line := range lines {
					Gray.Print("│ ")
					White.Println(line)
				}
			} else {
				// Show first 20 lines
				for i := 0; i < firstLines; i++ {
					Gray.Print("│ ")
					White.Println(lines[i])
				}

				// Show summary of middle content
				Gray.Printf("│ ... (%d more lines) ...\n", totalLines-firstLines-lastLines)

				// Show last 20 lines
				for i := totalLines - lastLines; i < totalLines; i++ {
					Gray.Print("│ ")
					White.Println(lines[i])
				}
			}
		}

		Gray.Println("└─")
	}
}

func displayResultMessage(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.IsError {
		BoldRed.Print("┌─ ")
		BoldRed.Print("RESULT: ERROR")
	} else {
		BoldBlue.Print("┌─ ")
		BoldBlue.Print("RESULT: SUCCESS")
	}
	Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

	// Show summary stats
	if msg.NumTurns > 0 {
		Blue.Printf("│ Turns: %d\n", msg.NumTurns)
	}
	if msg.DurationMS > 0 {
		Blue.Printf("│ Duration: %.2fs", float64(msg.DurationMS)/1000.0)
		if msg.DurationAPIMS > 0 {
			Blue.Printf(" (API: %.2fs)", float64(msg.DurationAPIMS)/1000.0)
		}
		Blue.Println()
	}
	if msg.TotalCostUSD > 0 {
		Blue.Printf("│ Cost: $%.4f\n", msg.TotalCostUSD)
	}

	// Show detailed token usage
	if msg.Usage != nil {
		Blue.Println("│")
		Blue.Print("│ ")
		Blue.Printf("Tokens: in=%d out=%d", msg.Usage.InputTokens, msg.Usage.OutputTokens)
		if msg.Usage.CacheReadInputTokens > 0 {
			Blue.Printf(" cache_read=%d", msg.Usage.CacheReadInputTokens)
		}
		if msg.Usage.CacheCreationInputTokens > 0 {
			Blue.Printf(" cache_create=%d", msg.Usage.CacheCreationInputTokens)
		}
		Blue.Println()
	}

	// Show per-model usage in verbose mode
	if cfg.Verbose && msg.ModelUsage != nil && len(msg.ModelUsage) > 0 {
		Blue.Println("│")
		Blue.Println("│ Model Usage:")
		for model, usageData := range msg.ModelUsage {
			Blue.Printf("│   %s:\n", model)
			if usageMap, ok := usageData.(map[string]interface{}); ok {
				if inputTokens, ok := usageMap["inputTokens"].(float64); ok {
					Blue.Printf("│     Input: %.0f tokens\n", inputTokens)
				}
				if outputTokens, ok := usageMap["outputTokens"].(float64); ok {
					Blue.Printf("│     Output: %.0f tokens\n", outputTokens)
				}
				if cost, ok := usageMap["costUSD"].(float64); ok {
					Blue.Printf("│     Cost: $%.4f\n", cost)
				}
			}
		}
	}

	// Show permission denials if present
	if len(msg.PermissionDenials) > 0 {
		Blue.Println("│")
		Red.Printf("│ Permission Denials: %d\n", len(msg.PermissionDenials))
		if cfg.Verbose {
			for i, denial := range msg.PermissionDenials {
				Red.Printf("│   [%d] %v\n", i+1, denial)
			}
		}
	}

	// Show result content if present
	if msg.Result != "" {
		Blue.Println("│")
		// Split result into lines and display
		lines := strings.Split(msg.Result, "\n")
		for _, line := range lines {
			Blue.Print("│ ")
			White.Println(line)
		}
	}

	Blue.Println("└─")
}
