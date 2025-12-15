package display

import (
	"fmt"
	"strings"

	"github.com/ariel-frischer/claude-clean/internal/parser"
)

func displayMessageCompact(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	switch msg.Type {
	case "system":
		displaySystemMessageCompact(msg, lineNum, cfg)
	case "assistant":
		displayAssistantMessageCompact(msg, lineNum, cfg)
	case "user":
		displayUserMessageCompact(msg, lineNum, cfg)
	case "result":
		displayResultMessageCompact(msg, lineNum, cfg)
	}
}

func displaySystemMessageCompact(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	BoldCyan.Print("SYS")
	if msg.Subtype != "" {
		Cyan.Printf("[%s]", msg.Subtype)
	}
	Gray.Printf("%s", FormatLineNumCompact(lineNum, cfg.ShowLineNum))
	if msg.Model != "" {
		Cyan.Printf(" %s", msg.Model)
	}
	if msg.CWD != "" {
		Cyan.Printf(" @%s", msg.CWD)
	}
	fmt.Println()
}

func displayAssistantMessageCompact(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		switch block.Type {
		case "text":
			if block.Text != "" {
				BoldGreen.Print("AST")
				Gray.Printf("%s ", FormatLineNumCompact(lineNum, cfg.ShowLineNum))
				// Truncate long text to single line
				text := strings.ReplaceAll(block.Text, "\n", " ")
				if len(text) > 100 {
					White.Printf("%s...\n", text[:100])
				} else {
					White.Println(text)
				}
			}
		case "tool_use":
			displayToolUseCompact(&block, lineNum, cfg)
		}
	}
}

func displayToolUseCompact(tool *parser.ContentBlock, lineNum int, cfg *Config) {
	BoldYellow.Printf("TOOL")
	Gray.Printf("%s ", FormatLineNumCompact(lineNum, cfg.ShowLineNum))
	Yellow.Printf("%s", tool.Name)

	// Show key inputs in compact form
	if tool.Input != nil {
		Yellow.Print(" {")
		first := true
		for key, value := range tool.Input {
			if !first {
				Yellow.Print(", ")
			}
			first = false

			switch v := value.(type) {
			case string:
				if len(v) > 50 {
					Yellow.Printf("%s: \"%.50s...\"", key, v)
				} else {
					Yellow.Printf("%s: \"%s\"", key, v)
				}
			case []interface{}:
				Yellow.Printf("%s: [%d items]", key, len(v))
			default:
				Yellow.Printf("%s: %v", key, v)
			}
		}
		Yellow.Print("}")
	}
	fmt.Println()
}

func displayUserMessageCompact(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		if block.Type == "tool_result" {
			displayToolResultCompact(&block, lineNum, cfg)
		}
	}
}

func displayToolResultCompact(block *parser.ContentBlock, lineNum int, cfg *Config) {
	if block.IsError {
		BoldRed.Print("ERR")
	} else {
		BoldMagenta.Print("RES")
	}
	Gray.Printf("%s ", FormatLineNumCompact(lineNum, cfg.ShowLineNum))

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

	// Compact output - single line summary
	contentStr = strings.ReplaceAll(contentStr, "\n", " ")
	if contentStr == "" {
		Gray.Println("(no output)")
	} else if len(contentStr) > 100 {
		White.Printf("%.100s...\n", contentStr)
	} else {
		White.Println(contentStr)
	}
}

func displayResultMessageCompact(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	if msg.IsError {
		BoldRed.Print("FAIL")
	} else {
		BoldBlue.Print("OK")
	}
	Gray.Printf("%s", FormatLineNumCompact(lineNum, cfg.ShowLineNum))

	if msg.NumTurns > 0 {
		Blue.Printf(" turns=%d", msg.NumTurns)
	}
	if msg.DurationMS > 0 {
		Blue.Printf(" %.2fs", float64(msg.DurationMS)/1000.0)
	}
	if msg.TotalCostUSD > 0 {
		Blue.Printf(" $%.4f", msg.TotalCostUSD)
	}
	if msg.Usage != nil {
		Blue.Printf(" in=%d out=%d", msg.Usage.InputTokens, msg.Usage.OutputTokens)
	}
	fmt.Println()

	// Show result text if present
	if msg.Result != "" {
		result := strings.ReplaceAll(msg.Result, "\n", " ")
		if len(result) > 200 {
			White.Printf("  %s...\n", result[:200])
		} else {
			White.Printf("  %s\n", result)
		}
	}
}
