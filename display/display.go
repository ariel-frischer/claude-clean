package display

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ariel-frischer/claude-clean/parser"
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
	Style          OutputStyle
	Verbose        bool
	ShowLineNum    bool
	ShowTimestamps bool
	StartTime      time.Time
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
	if !cfg.Verbose && isRoutineEvent(msg) {
		return
	}
	switch msg.Type {
	case "stream_event":
		// --include-partial-messages token deltas; the complete assistant message follows
		return
	case "rate_limit_event":
		displayRateLimit(msg, lineNum, cfg)
		return
	}
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

// isRoutineEvent reports bookkeeping events that carry no news unless something went wrong:
// hook start/success notices, API status pings, subagent progress ticks, and rate-limit
// updates that are still "allowed".
func isRoutineEvent(msg *parser.StreamMessage) bool {
	switch msg.Type {
	case "rate_limit_event":
		return msg.RateLimitInfo == nil || msg.RateLimitInfo.Status == "allowed"
	case "system":
		switch msg.Subtype {
		case "hook_started", "status", "task_progress", "task_updated":
			return true
		case "hook_response":
			return msg.Outcome == "success"
		}
	}
	return false
}

func displayRateLimit(msg *parser.StreamMessage, lineNum int, cfg *Config) {
	line := "RATE LIMIT"
	if info := msg.RateLimitInfo; info != nil {
		line += fmt.Sprintf(" [%s]", info.Status)
		if info.RateLimitType != "" {
			line += " " + info.RateLimitType
		}
		if info.ResetsAt > 0 {
			line += " resets " + time.Unix(info.ResetsAt, 0).Local().Format("Jan 2 15:04")
		}
	}
	suffix := FormatElapsed(cfg) + FormatLineNum(lineNum, cfg.ShowLineNum)
	if cfg.Style == StylePlain {
		fmt.Println(line + suffix)
		return
	}
	BoldYellow.Print(line)
	Gray.Println(suffix)
}

// systemSummary describes hook and subagent events for system message headers, e.g.
// " SessionStart:startup", " Stop: error (exit 1): boom", or " Run tests (general-purpose)".
// Empty for other system messages.
func systemSummary(msg *parser.StreamMessage) string {
	switch msg.Subtype {
	case "hook_started", "hook_response":
		if msg.HookName == "" {
			return ""
		}
		s := " " + msg.HookName
		if msg.Subtype == "hook_response" && msg.Outcome != "success" {
			s += fmt.Sprintf(": %s (exit %d)", msg.Outcome, msg.ExitCode)
			if stderr := firstLine(msg.Stderr); stderr != "" {
				s += ": " + stderr
			}
		}
		return s
	case "task_started":
		s := " " + msg.Description
		if msg.SubagentType != "" {
			s += " (" + msg.SubagentType + ")"
		}
		return s
	case "task_notification":
		s := " " + msg.Status
		if summary := firstLine(msg.Summary); summary != "" {
			s += ": " + summary
		}
		return s
	case "status":
		return " " + msg.Status
	}
	return ""
}

// firstLine returns the first non-empty line of s, cut to 100 runes
func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			if r := []rune(line); len(r) > 100 {
				return string(r[:100]) + "..."
			}
			return line
		}
	}
	return ""
}

// sortedKeys returns map keys in stable order so output is deterministic
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
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

// FormatElapsed returns the elapsed time since cfg.StartTime if ShowTimestamps is enabled
func FormatElapsed(cfg *Config) string {
	if !cfg.ShowTimestamps || cfg.StartTime.IsZero() {
		return ""
	}
	d := time.Since(cfg.StartTime)
	secs := d.Seconds()
	if secs < 60 {
		return fmt.Sprintf(" +%.1fs", secs)
	}
	mins := int(secs) / 60
	remSecs := int(secs) % 60
	return fmt.Sprintf(" +%dm%ds", mins, remSecs)
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
