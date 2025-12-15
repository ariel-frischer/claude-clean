package display

import (
	"fmt"
	"strings"

	"github.com/ariel-frischer/claude-clean/internal/parser"
)

func displayMessageMinimal(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	switch msg.Type {
	case "system":
		displaySystemMessageMinimal(msg, lineNum, cfg)
	case "assistant":
		displayAssistantMessageMinimal(msg, lineNum, cfg)
	case "user":
		displayUserMessageMinimal(msg, lineNum, cfg)
	case "result":
		displayResultMessageMinimal(msg, lineNum, cfg)
	}
}

func displaySystemMessageMinimal(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	BoldCyan.Printf("SYSTEM")
	if msg.Subtype != "" {
		Cyan.Printf(" [%s]", msg.Subtype)
	}
	Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

	if msg.CWD != "" {
		Cyan.Printf("  Working Directory: %s\n", msg.CWD)
	}
	if msg.Model != "" {
		Cyan.Printf("  Model: %s\n", msg.Model)
	}
	if msg.ClaudeCodeVersion != "" {
		Cyan.Printf("  Claude Code: v%s\n", msg.ClaudeCodeVersion)
	}
	if len(msg.Tools) > 0 {
		Cyan.Printf("  Tools: %d available\n", len(msg.Tools))
	}
	fmt.Println()
}

func displayAssistantMessageMinimal(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	// Group consecutive text blocks
	var textBlocks []string
	var toolUses []parser.ContentBlock

	for _, block := range msg.Message.Content {
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
		BoldGreen.Printf("ASSISTANT")
		Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		for _, text := range textBlocks {
			White.Printf("  %s\n", text)
		}

		if cfg.Verbose && msg.Message.Usage != nil {
			Gray.Printf("  Tokens: in=%d out=%d", msg.Message.Usage.InputTokens, msg.Message.Usage.OutputTokens)
			if msg.Message.Usage.CacheReadInputTokens > 0 {
				Gray.Printf(" cache_read=%d", msg.Message.Usage.CacheReadInputTokens)
			}
			if msg.Message.Usage.CacheCreationInputTokens > 0 {
				Gray.Printf(" cache_create=%d", msg.Message.Usage.CacheCreationInputTokens)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Display tool uses
	for _, tool := range toolUses {
		displayToolUseMinimal(&tool, lineNum, cfg)
	}
}

func displayToolUseMinimal(tool *parser.ContentBlock, lineNum int, cfg *Config) {
	BoldYellow.Printf("TOOL: %s", tool.Name)
	Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

	if cfg.Verbose {
		Yellow.Printf("  ID: %s\n", tool.ID)
	}

	if tool.Input != nil {
		Yellow.Println("  Input:")
		for key, value := range tool.Input {
			Yellow.Printf("    %s: ", key)

			switch v := value.(type) {
			case string:
				if len(v) > 300 {
					White.Printf("%s ... (%d chars omitted) ... %s\n", v[:200], len(v)-300, v[len(v)-100:])
				} else {
					White.Println(v)
				}
			case []interface{}:
				if tool.Name == "TodoWrite" && key == "todos" {
					White.Println()
					DisplayTodosMinimal(v)
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
	fmt.Println()
}

func displayUserMessageMinimal(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		if block.Type == "tool_result" {
			displayToolResultMinimal(&block, lineNum, cfg)
		}
	}
}

func displayToolResultMinimal(block *parser.ContentBlock, lineNum int, cfg *Config) {
	if block.IsError {
		BoldRed.Printf("TOOL RESULT ERROR")
		Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		if cfg.Verbose {
			Red.Printf("  Tool ID: %s\n", block.ToolUseID)
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

		White.Printf("  %s\n", contentStr)
	} else {
		BoldMagenta.Printf("TOOL RESULT")
		Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		if cfg.Verbose {
			Gray.Printf("  Tool ID: %s\n", block.ToolUseID)
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
			Gray.Println("  (no output)")
		} else {
			lines := strings.Split(contentStr, "\n")
			firstLines := parser.FirstLines
			lastLines := parser.LastLines
			totalLines := len(lines)

			if totalLines <= firstLines+lastLines {
				for _, line := range lines {
					White.Printf("  %s\n", line)
				}
			} else {
				for i := 0; i < firstLines; i++ {
					White.Printf("  %s\n", lines[i])
				}
				Gray.Printf("  ... (%d more lines) ...\n", totalLines-firstLines-lastLines)
				for i := totalLines - lastLines; i < totalLines; i++ {
					White.Printf("  %s\n", lines[i])
				}
			}
		}
	}
	fmt.Println()
}

func displayResultMessageMinimal(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.IsError {
		BoldRed.Printf("RESULT: ERROR")
	} else {
		BoldBlue.Printf("RESULT: SUCCESS")
	}
	Gray.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

	if msg.NumTurns > 0 {
		Blue.Printf("  Turns: %d\n", msg.NumTurns)
	}
	if msg.DurationMS > 0 {
		Blue.Printf("  Duration: %.2fs", float64(msg.DurationMS)/1000.0)
		if msg.DurationAPIMS > 0 {
			Blue.Printf(" (API: %.2fs)", float64(msg.DurationAPIMS)/1000.0)
		}
		Blue.Println()
	}
	if msg.TotalCostUSD > 0 {
		Blue.Printf("  Cost: $%.4f\n", msg.TotalCostUSD)
	}

	if msg.Usage != nil {
		Blue.Printf("  Tokens: in=%d out=%d", msg.Usage.InputTokens, msg.Usage.OutputTokens)
		if msg.Usage.CacheReadInputTokens > 0 {
			Blue.Printf(" cache_read=%d", msg.Usage.CacheReadInputTokens)
		}
		if msg.Usage.CacheCreationInputTokens > 0 {
			Blue.Printf(" cache_create=%d", msg.Usage.CacheCreationInputTokens)
		}
		Blue.Println()
	}

	if cfg.Verbose && msg.ModelUsage != nil && len(msg.ModelUsage) > 0 {
		Blue.Println()
		Blue.Println("  Model Usage:")
		for model, usageData := range msg.ModelUsage {
			Blue.Printf("    %s:\n", model)
			if usageMap, ok := usageData.(map[string]interface{}); ok {
				if inputTokens, ok := usageMap["inputTokens"].(float64); ok {
					Blue.Printf("      Input: %.0f tokens\n", inputTokens)
				}
				if outputTokens, ok := usageMap["outputTokens"].(float64); ok {
					Blue.Printf("      Output: %.0f tokens\n", outputTokens)
				}
				if cost, ok := usageMap["costUSD"].(float64); ok {
					Blue.Printf("      Cost: $%.4f\n", cost)
				}
			}
		}
	}

	if len(msg.PermissionDenials) > 0 {
		fmt.Println()
		Red.Printf("  Permission Denials: %d\n", len(msg.PermissionDenials))
		if cfg.Verbose {
			for i, denial := range msg.PermissionDenials {
				Red.Printf("    [%d] %v\n", i+1, denial)
			}
		}
	}

	if msg.Result != "" {
		fmt.Println()
		lines := strings.Split(msg.Result, "\n")
		for _, line := range lines {
			White.Printf("  %s\n", line)
		}
	}

	fmt.Println()
}
