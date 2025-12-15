package display

import (
	"fmt"
	"strings"

	"github.com/ariel-frischer/claude-clean/internal/parser"
)

func displayMessagePlain(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	switch msg.Type {
	case "system":
		displaySystemMessagePlain(msg, lineNum, cfg)
	case "assistant":
		displayAssistantMessagePlain(msg, lineNum, cfg)
	case "user":
		displayUserMessagePlain(msg, lineNum, cfg)
	case "result":
		displayResultMessagePlain(msg, lineNum, cfg)
	}
}

func displaySystemMessagePlain(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	fmt.Printf("SYSTEM")
	if msg.Subtype != "" {
		fmt.Printf(" [%s]", msg.Subtype)
	}
	fmt.Printf("%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

	if msg.CWD != "" {
		fmt.Printf("  Working Directory: %s\n", msg.CWD)
	}
	if msg.Model != "" {
		fmt.Printf("  Model: %s\n", msg.Model)
	}
	if msg.ClaudeCodeVersion != "" {
		fmt.Printf("  Claude Code: v%s\n", msg.ClaudeCodeVersion)
	}
	if len(msg.Tools) > 0 {
		fmt.Printf("  Tools: %d available\n", len(msg.Tools))
	}
	fmt.Println()
}

func displayAssistantMessagePlain(msg *parser.StreamMessage, lineNum int, cfg *Config) {
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
		fmt.Printf("ASSISTANT%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		for _, text := range textBlocks {
			fmt.Printf("  %s\n", text)
		}

		if cfg.Verbose && msg.Message.Usage != nil {
			fmt.Printf("  Tokens: in=%d out=%d", msg.Message.Usage.InputTokens, msg.Message.Usage.OutputTokens)
			if msg.Message.Usage.CacheReadInputTokens > 0 {
				fmt.Printf(" cache_read=%d", msg.Message.Usage.CacheReadInputTokens)
			}
			if msg.Message.Usage.CacheCreationInputTokens > 0 {
				fmt.Printf(" cache_create=%d", msg.Message.Usage.CacheCreationInputTokens)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Display tool uses
	for _, tool := range toolUses {
		displayToolUsePlain(&tool, lineNum, cfg)
	}
}

func displayToolUsePlain(tool *parser.ContentBlock, lineNum int, cfg *Config) {
	fmt.Printf("TOOL: %s%s\n", tool.Name, FormatLineNum(lineNum, cfg.ShowLineNum))

	if cfg.Verbose {
		fmt.Printf("  ID: %s\n", tool.ID)
	}

	if tool.Input != nil {
		fmt.Println("  Input:")
		for key, value := range tool.Input {
			fmt.Printf("    %s: ", key)

			switch v := value.(type) {
			case string:
				if len(v) > 300 {
					fmt.Printf("%s ... (%d chars omitted) ... %s\n", v[:200], len(v)-300, v[len(v)-100:])
				} else {
					fmt.Println(v)
				}
			case []interface{}:
				if tool.Name == "TodoWrite" && key == "todos" {
					fmt.Println()
					DisplayTodosPlain(v)
				} else {
					fmt.Printf("[%d items]\n", len(v))
				}
			case map[string]interface{}:
				fmt.Println("{...}")
			default:
				fmt.Printf("%v\n", v)
			}
		}
	}
	fmt.Println()
}

func displayUserMessagePlain(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		if block.Type == "tool_result" {
			displayToolResultPlain(&block, lineNum, cfg)
		}
	}
}

func displayToolResultPlain(block *parser.ContentBlock, lineNum int, cfg *Config) {
	if block.IsError {
		fmt.Printf("TOOL RESULT ERROR%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		if cfg.Verbose {
			fmt.Printf("  Tool ID: %s\n", block.ToolUseID)
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

		fmt.Printf("  %s\n", contentStr)
	} else {
		fmt.Printf("TOOL RESULT%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))

		if cfg.Verbose {
			fmt.Printf("  Tool ID: %s\n", block.ToolUseID)
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
			fmt.Println("  (no output)")
		} else {
			lines := strings.Split(contentStr, "\n")
			firstLines := parser.FirstLines
			lastLines := parser.LastLines
			totalLines := len(lines)

			if totalLines <= firstLines+lastLines {
				for _, line := range lines {
					fmt.Printf("  %s\n", line)
				}
			} else {
				for i := 0; i < firstLines; i++ {
					fmt.Printf("  %s\n", lines[i])
				}
				fmt.Printf("  ... (%d more lines) ...\n", totalLines-firstLines-lastLines)
				for i := totalLines - lastLines; i < totalLines; i++ {
					fmt.Printf("  %s\n", lines[i])
				}
			}
		}
	}
	fmt.Println()
}

func displayResultMessagePlain(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.IsError {
		fmt.Printf("RESULT: ERROR%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))
	} else {
		fmt.Printf("RESULT: SUCCESS%s\n", FormatLineNum(lineNum, cfg.ShowLineNum))
	}

	if msg.NumTurns > 0 {
		fmt.Printf("  Turns: %d\n", msg.NumTurns)
	}
	if msg.DurationMS > 0 {
		fmt.Printf("  Duration: %.2fs", float64(msg.DurationMS)/1000.0)
		if msg.DurationAPIMS > 0 {
			fmt.Printf(" (API: %.2fs)", float64(msg.DurationAPIMS)/1000.0)
		}
		fmt.Println()
	}
	if msg.TotalCostUSD > 0 {
		fmt.Printf("  Cost: $%.4f\n", msg.TotalCostUSD)
	}

	if msg.Usage != nil {
		fmt.Printf("  Tokens: in=%d out=%d", msg.Usage.InputTokens, msg.Usage.OutputTokens)
		if msg.Usage.CacheReadInputTokens > 0 {
			fmt.Printf(" cache_read=%d", msg.Usage.CacheReadInputTokens)
		}
		if msg.Usage.CacheCreationInputTokens > 0 {
			fmt.Printf(" cache_create=%d", msg.Usage.CacheCreationInputTokens)
		}
		fmt.Println()
	}

	if cfg.Verbose && msg.ModelUsage != nil && len(msg.ModelUsage) > 0 {
		fmt.Println()
		fmt.Println("  Model Usage:")
		for model, usageData := range msg.ModelUsage {
			fmt.Printf("    %s:\n", model)
			if usageMap, ok := usageData.(map[string]interface{}); ok {
				if inputTokens, ok := usageMap["inputTokens"].(float64); ok {
					fmt.Printf("      Input: %.0f tokens\n", inputTokens)
				}
				if outputTokens, ok := usageMap["outputTokens"].(float64); ok {
					fmt.Printf("      Output: %.0f tokens\n", outputTokens)
				}
				if cost, ok := usageMap["costUSD"].(float64); ok {
					fmt.Printf("      Cost: $%.4f\n", cost)
				}
			}
		}
	}

	if len(msg.PermissionDenials) > 0 {
		fmt.Println()
		fmt.Printf("  Permission Denials: %d\n", len(msg.PermissionDenials))
		if cfg.Verbose {
			for i, denial := range msg.PermissionDenials {
				fmt.Printf("    [%d] %v\n", i+1, denial)
			}
		}
	}

	if msg.Result != "" {
		fmt.Println()
		lines := strings.Split(msg.Result, "\n")
		for _, line := range lines {
			fmt.Printf("  %s\n", line)
		}
	}

	fmt.Println()
}
