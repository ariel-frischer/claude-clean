package display

import (
	"fmt"
	"strings"

	"github.com/ariel-frischer/claude-clean/internal/parser"
	"github.com/fatih/color"
)

// OutputStyle represents different output formatting styles
type OutputStyle string

const (
	StyleDefault OutputStyle = "default"
	StyleCompact OutputStyle = "compact"
	StyleMinimal OutputStyle = "minimal"
	StylePlain   OutputStyle = "plain"
)

// Config holds display configuration options
type Config struct {
	Style       OutputStyle
	Verbose     bool
	ShowLineNum bool
}

// Color definitions
var (
	BoldCyan    = color.New(color.FgCyan, color.Bold)
	BoldGreen   = color.New(color.FgGreen, color.Bold)
	BoldYellow  = color.New(color.FgYellow, color.Bold)
	BoldRed     = color.New(color.FgRed, color.Bold)
	BoldMagenta = color.New(color.FgMagenta, color.Bold)
	BoldBlue    = color.New(color.FgBlue, color.Bold)
	Cyan        = color.New(color.FgCyan)
	Green       = color.New(color.FgGreen)
	Yellow      = color.New(color.FgYellow)
	Red         = color.New(color.FgRed)
	Blue        = color.New(color.FgBlue)
	Gray        = color.New(color.FgHiBlack)
	White       = color.New(color.FgWhite)
)

// DisplayMessage routes to the appropriate formatter based on style
func DisplayMessage(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	switch cfg.Style {
	case StyleCompact:
		displayMessageCompact(msg, lineNum, cfg)
	case StyleMinimal:
		displayMessageMinimal(msg, lineNum, cfg)
	case StylePlain:
		displayMessagePlain(msg, lineNum, cfg)
	default: // StyleDefault
		displayMessageDefault(msg, lineNum, cfg)
	}
}

// FormatLineNum returns a formatted line number string if showLineNum is enabled
func FormatLineNum(lineNum int, showLineNum bool) string {
	if showLineNum {
		return fmt.Sprintf(" (line %d)", lineNum)
	}
	return ""
}

// FormatLineNumCompact returns a compact line number format if showLineNum is enabled
func FormatLineNumCompact(lineNum int, showLineNum bool) string {
	if showLineNum {
		return fmt.Sprintf(" L%d", lineNum)
	}
	return ""
}

// DisplayUsage shows token usage statistics
func DisplayUsage(usage *parser.Usage) {
	Gray.Print("│ ")
	Gray.Printf("Tokens: in=%d out=%d", usage.InputTokens, usage.OutputTokens)

	if usage.CacheReadInputTokens > 0 {
		Gray.Printf(" cache_read=%d", usage.CacheReadInputTokens)
	}
	if usage.CacheCreationInputTokens > 0 {
		Gray.Printf(" cache_create=%d", usage.CacheCreationInputTokens)
	}

	Gray.Println()
}

// DisplayUsageInline shows usage inline with a specific color
func DisplayUsageInline(usage *parser.Usage, c *color.Color) {
	c.Print("│ ")
	c.Printf("Tokens: in=%d out=%d", usage.InputTokens, usage.OutputTokens)

	if usage.CacheReadInputTokens > 0 {
		c.Printf(" cache_read=%d", usage.CacheReadInputTokens)
	}
	if usage.CacheCreationInputTokens > 0 {
		c.Printf(" cache_create=%d", usage.CacheCreationInputTokens)
	}

	c.Println()
}

// DisplayTodos displays todo items with status icons
func DisplayTodos(todos []interface{}) {
	for i, todo := range todos {
		todoMap, ok := todo.(map[string]interface{})
		if !ok {
			continue
		}

		content, _ := todoMap["content"].(string)
		status, _ := todoMap["status"].(string)

		// Format status with color
		var statusIcon string
		switch status {
		case "completed":
			statusIcon = Green.Sprint("✓")
		case "in_progress":
			statusIcon = Yellow.Sprint("→")
		case "pending":
			statusIcon = Gray.Sprint("○")
		default:
			statusIcon = Gray.Sprint("-")
		}

		Yellow.Printf("│     %s %s\n", statusIcon, content)
		if i < len(todos)-1 {
			// Continue without extra line between todos
		}
	}
}

// DisplayTodosMinimal displays todos in minimal style
func DisplayTodosMinimal(todos []interface{}) {
	for _, todo := range todos {
		todoMap, ok := todo.(map[string]interface{})
		if !ok {
			continue
		}

		content, _ := todoMap["content"].(string)
		status, _ := todoMap["status"].(string)

		var statusIcon string
		switch status {
		case "completed":
			statusIcon = Green.Sprint("✓")
		case "in_progress":
			statusIcon = Yellow.Sprint("→")
		case "pending":
			statusIcon = Gray.Sprint("○")
		default:
			statusIcon = Gray.Sprint("-")
		}

		Yellow.Printf("      %s %s\n", statusIcon, content)
	}
}

// DisplayTodosPlain displays todos in plain style
func DisplayTodosPlain(todos []interface{}) {
	for _, todo := range todos {
		todoMap, ok := todo.(map[string]interface{})
		if !ok {
			continue
		}

		content, _ := todoMap["content"].(string)
		status, _ := todoMap["status"].(string)

		var statusIcon string
		switch status {
		case "completed":
			statusIcon = "[✓]"
		case "in_progress":
			statusIcon = "[→]"
		case "pending":
			statusIcon = "[○]"
		default:
			statusIcon = "[-]"
		}

		fmt.Printf("      %s %s\n", statusIcon, content)
	}
}

// TruncateLongOutput shows first N + last N lines for long content
func TruncateLongOutput(contentStr string, prefix string, printFn func(string)) {
	lines := strings.Split(contentStr, "\n")
	firstLines := parser.FirstLines
	lastLines := parser.LastLines
	totalLines := len(lines)

	if totalLines <= firstLines+lastLines {
		// Show all lines if content is short enough
		for _, line := range lines {
			printFn(prefix + line + "\n")
		}
	} else {
		// Show first N lines
		for i := 0; i < firstLines; i++ {
			printFn(prefix + lines[i] + "\n")
		}

		// Show summary of middle content
		printFn(fmt.Sprintf("%s... (%d more lines) ...\n", prefix, totalLines-firstLines-lastLines))

		// Show last N lines
		for i := totalLines - lastLines; i < totalLines; i++ {
			printFn(prefix + lines[i] + "\n")
		}
	}
}
