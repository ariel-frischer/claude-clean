package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Color definitions
var (
	boldCyan    = color.New(color.FgCyan, color.Bold)
	boldGreen   = color.New(color.FgGreen, color.Bold)
	boldYellow  = color.New(color.FgYellow, color.Bold)
	boldRed     = color.New(color.FgRed, color.Bold)
	boldMagenta = color.New(color.FgMagenta, color.Bold)
	boldBlue    = color.New(color.FgBlue, color.Bold)
	cyan        = color.New(color.FgCyan)
	green       = color.New(color.FgGreen)
	yellow      = color.New(color.FgYellow)
	red         = color.New(color.FgRed)
	blue        = color.New(color.FgBlue)
	gray        = color.New(color.FgHiBlack)
	white       = color.New(color.FgWhite)
)

// OutputStyle represents different output formatting styles
type OutputStyle string

const (
	StyleDefault OutputStyle = "default"
	StyleCompact OutputStyle = "compact"
	StyleMinimal OutputStyle = "minimal"
	StylePlain   OutputStyle = "plain"
)

// Command-line flags
var (
	verbose = flag.Bool("v", false, "Show verbose output (tool IDs, token usage)")
	help    = flag.Bool("h", false, "Show help message")
	style   = flag.String("s", "default", "Output style: default, compact, minimal, plain")
)

// Global style setting
var currentStyle OutputStyle

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage: claude-clean-output [OPTIONS] [FILE]")
		fmt.Println("\nParse and beautify Claude Code's streaming JSON output")
		fmt.Println("\nOptions:")
		fmt.Println("  -v            Show verbose output (tool IDs, token usage)")
		fmt.Println("  -s STYLE      Output style: default, compact, minimal, plain (default: default)")
		fmt.Println("  -h            Show this help message")
		fmt.Println("\nStyles:")
		fmt.Println("  default       Full boxed format with colors (current format)")
		fmt.Println("  compact       Minimal single-line format")
		fmt.Println("  minimal       Simple indented format without boxes")
		fmt.Println("  plain         No colors, no boxes, just text")
		fmt.Println("\nExamples:")
		fmt.Println("  claude-clean-output log.jsonl")
		fmt.Println("  cat log.jsonl | claude-clean-output")
		fmt.Println("  claude-clean-output -v log.jsonl          # Show detailed info")
		fmt.Println("  claude-clean-output -s compact log.jsonl  # Use compact style")
		fmt.Println("  claude-clean-output -s minimal log.jsonl  # Use minimal style")
		os.Exit(0)
	}

	// Validate and set the style
	switch OutputStyle(*style) {
	case StyleDefault, StyleCompact, StyleMinimal, StylePlain:
		currentStyle = OutputStyle(*style)
	default:
		fmt.Fprintf(os.Stderr, "Invalid style: %s. Valid styles are: default, compact, minimal, plain\n", *style)
		os.Exit(1)
	}

	var reader io.Reader

	// Check if we have a file argument or should read from stdin
	args := flag.Args()
	if len(args) > 0 {
		file, err := os.Open(args[0])
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
	// Increase buffer size to handle very large JSON lines (e.g., large file contents)
	// Default is 64KB, we set it to 10MB max
	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	lineNum := 0
	var lastAssistantMsg *StreamMessage
	var lastAssistantLine int

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

		// Handle duplicate detection between assistant and result messages
		if msg.Type == "result" && lastAssistantMsg != nil && !*verbose {
			// Check if the result message contains the same content as the last assistant message
			if msg.Result != "" && lastAssistantMsg.Message != nil && len(lastAssistantMsg.Message.Content) > 0 {
				// Get the last text from assistant message
				for _, block := range lastAssistantMsg.Message.Content {
					if block.Type == "text" && block.Text == msg.Result {
						// Skip the duplicate assistant message
						lastAssistantMsg = nil
						break
					}
				}
			}
			// Display the buffered assistant message if it wasn't a duplicate
			if lastAssistantMsg != nil {
				displayMessage(lastAssistantMsg, lastAssistantLine)
				lastAssistantMsg = nil
			}
		} else if lastAssistantMsg != nil {
			// Display the buffered assistant message
			displayMessage(lastAssistantMsg, lastAssistantLine)
			lastAssistantMsg = nil
		}

		// Buffer assistant messages to check for duplicates with result
		if msg.Type == "assistant" && !*verbose {
			lastAssistantMsg = &msg
			lastAssistantLine = lineNum
		} else {
			// Display other message types immediately
			displayMessage(&msg, lineNum)
		}
	}

	// Display any remaining buffered assistant message
	if lastAssistantMsg != nil {
		displayMessage(lastAssistantMsg, lastAssistantLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func displayMessage(msg *StreamMessage, lineNum int) {
	// Route to appropriate formatter based on style
	switch currentStyle {
	case StyleCompact:
		displayMessageCompact(msg, lineNum)
	case StyleMinimal:
		displayMessageMinimal(msg, lineNum)
	case StylePlain:
		displayMessagePlain(msg, lineNum)
	default: // StyleDefault
		displayMessageDefault(msg, lineNum)
	}
}

// stripSystemReminders removes <system-reminder>...</system-reminder> tags from content
func stripSystemReminders(content string) string {
	// Use regex to remove system-reminder tags and their content
	re := regexp.MustCompile(`(?s)<system-reminder>.*?</system-reminder>`)
	result := re.ReplaceAllString(content, "")

	// Clean up any resulting multiple blank lines
	result = regexp.MustCompile(`\n\n\n+`).ReplaceAllString(result, "\n\n")

	// Trim leading/trailing whitespace
	return strings.TrimSpace(result)
}

func displayMessageDefault(msg *StreamMessage, lineNum int) {
	switch msg.Type {
	case "system":
		displaySystemMessage(msg, lineNum)
	case "assistant":
		displayAssistantMessage(msg, lineNum)
	case "user":
		displayUserMessage(msg, lineNum)
	case "result":
		displayResultMessage(msg, lineNum)
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

		if *verbose && msg.Message.Usage != nil {
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

	if *verbose {
		yellow.Printf("│ ID: %s\n", tool.ID)
	}

	if tool.Input != nil {
		yellow.Println("│ Input:")
		for key, value := range tool.Input {
			// Pretty print the value
			yellow.Printf("│   %s: ", key)

			switch v := value.(type) {
			case string:
				// Show more context for strings - first 200 + last 100 chars
				if len(v) > 300 {
					white.Printf("%s ... (%d chars omitted) ... %s\n",
						v[:200], len(v)-300, v[len(v)-100:])
				} else {
					white.Println(v)
				}
			case []interface{}:
				// Special handling for todos array in TodoWrite tool
				if tool.Name == "TodoWrite" && key == "todos" {
					white.Println()
					displayTodos(v)
				} else {
					white.Printf("[%d items]\n", len(v))
				}
			case map[string]interface{}:
				white.Println("{...}")
			default:
				white.Printf("%v\n", v)
			}
		}
	}

	yellow.Println("└─")
}

func displayTodos(todos []interface{}) {
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
			statusIcon = green.Sprint("✓")
		case "in_progress":
			statusIcon = yellow.Sprint("→")
		case "pending":
			statusIcon = gray.Sprint("○")
		default:
			statusIcon = gray.Sprint("-")
		}

		yellow.Printf("│     %s %s\n", statusIcon, content)
		if i < len(todos)-1 {
			// Continue without extra line between todos
		}
	}
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

		if *verbose {
			red.Printf("│ Tool ID: %s\n", block.ToolUseID)
		}

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		// Strip system reminders in non-verbose mode
		if !*verbose {
			contentStr = stripSystemReminders(contentStr)
		}

		red.Print("│ ")
		white.Println(contentStr)
		red.Println("└─")
	} else {
		boldMagenta.Print("┌─ ")
		boldMagenta.Print("TOOL RESULT")
		gray.Printf(" (line %d)\n", lineNum)

		if *verbose {
			gray.Printf("│ Tool ID: %s\n", block.ToolUseID)
		}

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		// Strip system reminders in non-verbose mode
		if !*verbose {
			contentStr = stripSystemReminders(contentStr)
		}

		if contentStr == "" {
			gray.Println("│ (no output)")
		} else {
			// Show first 20 + last 20 lines for long output
			lines := strings.Split(contentStr, "\n")
			firstLines := 20
			lastLines := 20
			totalLines := len(lines)

			if totalLines <= firstLines+lastLines {
				// Show all lines if content is short enough
				for _, line := range lines {
					gray.Print("│ ")
					white.Println(line)
				}
			} else {
				// Show first 20 lines
				for i := 0; i < firstLines; i++ {
					gray.Print("│ ")
					white.Println(lines[i])
				}

				// Show summary of middle content
				gray.Printf("│ ... (%d more lines) ...\n", totalLines-firstLines-lastLines)

				// Show last 20 lines
				for i := totalLines - lastLines; i < totalLines; i++ {
					gray.Print("│ ")
					white.Println(lines[i])
				}
			}
		}

		gray.Println("└─")
	}
}

func displayResultMessage(msg *StreamMessage, lineNum int) {
	if msg.IsError {
		boldRed.Print("┌─ ")
		boldRed.Print("RESULT: ERROR")
	} else {
		boldBlue.Print("┌─ ")
		boldBlue.Print("RESULT: SUCCESS")
	}
	gray.Printf(" (line %d)\n", lineNum)

	// Show summary stats
	if msg.NumTurns > 0 {
		blue.Printf("│ Turns: %d\n", msg.NumTurns)
	}
	if msg.DurationMS > 0 {
		blue.Printf("│ Duration: %.2fs", float64(msg.DurationMS)/1000.0)
		if msg.DurationAPIMS > 0 {
			blue.Printf(" (API: %.2fs)", float64(msg.DurationAPIMS)/1000.0)
		}
		blue.Println()
	}
	if msg.TotalCostUSD > 0 {
		blue.Printf("│ Cost: $%.4f\n", msg.TotalCostUSD)
	}

	// Show detailed token usage
	if msg.Usage != nil {
		blue.Println("│")
		blue.Print("│ ")
		blue.Printf("Tokens: in=%d out=%d", msg.Usage.InputTokens, msg.Usage.OutputTokens)
		if msg.Usage.CacheReadInputTokens > 0 {
			blue.Printf(" cache_read=%d", msg.Usage.CacheReadInputTokens)
		}
		if msg.Usage.CacheCreationInputTokens > 0 {
			blue.Printf(" cache_create=%d", msg.Usage.CacheCreationInputTokens)
		}
		blue.Println()
	}

	// Show per-model usage in verbose mode
	if *verbose && msg.ModelUsage != nil && len(msg.ModelUsage) > 0 {
		blue.Println("│")
		blue.Println("│ Model Usage:")
		for model, usageData := range msg.ModelUsage {
			blue.Printf("│   %s:\n", model)
			if usageMap, ok := usageData.(map[string]interface{}); ok {
				if inputTokens, ok := usageMap["inputTokens"].(float64); ok {
					blue.Printf("│     Input: %.0f tokens\n", inputTokens)
				}
				if outputTokens, ok := usageMap["outputTokens"].(float64); ok {
					blue.Printf("│     Output: %.0f tokens\n", outputTokens)
				}
				if cost, ok := usageMap["costUSD"].(float64); ok {
					blue.Printf("│     Cost: $%.4f\n", cost)
				}
			}
		}
	}

	// Show permission denials if present
	if len(msg.PermissionDenials) > 0 {
		blue.Println("│")
		red.Printf("│ Permission Denials: %d\n", len(msg.PermissionDenials))
		if *verbose {
			for i, denial := range msg.PermissionDenials {
				red.Printf("│   [%d] %v\n", i+1, denial)
			}
		}
	}

	// Show result content if present
	if msg.Result != "" {
		blue.Println("│")
		// Split result into lines and display
		lines := strings.Split(msg.Result, "\n")
		for _, line := range lines {
			blue.Print("│ ")
			white.Println(line)
		}
	}

	blue.Println("└─")
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

func displayUsageInline(usage *Usage, c *color.Color) {
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

// ============================================================================
// COMPACT STYLE FORMATTERS
// ============================================================================

func displayMessageCompact(msg *StreamMessage, lineNum int) {
	switch msg.Type {
	case "system":
		displaySystemMessageCompact(msg, lineNum)
	case "assistant":
		displayAssistantMessageCompact(msg, lineNum)
	case "user":
		displayUserMessageCompact(msg, lineNum)
	case "result":
		displayResultMessageCompact(msg, lineNum)
	}
}

func displaySystemMessageCompact(msg *StreamMessage, lineNum int) {
	boldCyan.Print("SYS")
	if msg.Subtype != "" {
		cyan.Printf("[%s]", msg.Subtype)
	}
	gray.Printf(" L%d", lineNum)
	if msg.Model != "" {
		cyan.Printf(" %s", msg.Model)
	}
	if msg.CWD != "" {
		cyan.Printf(" @%s", msg.CWD)
	}
	fmt.Println()
}

func displayAssistantMessageCompact(msg *StreamMessage, lineNum int) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		switch block.Type {
		case "text":
			if block.Text != "" {
				boldGreen.Print("AST ")
				gray.Printf("L%d ", lineNum)
				// Truncate long text to single line
				text := strings.ReplaceAll(block.Text, "\n", " ")
				if len(text) > 100 {
					white.Printf("%s...\n", text[:100])
				} else {
					white.Println(text)
				}
			}
		case "tool_use":
			displayToolUseCompact(&block, lineNum)
		}
	}
}

func displayToolUseCompact(tool *ContentBlock, lineNum int) {
	boldYellow.Printf("TOOL ")
	gray.Printf("L%d ", lineNum)
	yellow.Printf("%s", tool.Name)

	// Show key inputs in compact form
	if tool.Input != nil {
		yellow.Print(" {")
		first := true
		for key, value := range tool.Input {
			if !first {
				yellow.Print(", ")
			}
			first = false

			switch v := value.(type) {
			case string:
				if len(v) > 50 {
					yellow.Printf("%s: \"%.50s...\"", key, v)
				} else {
					yellow.Printf("%s: \"%s\"", key, v)
				}
			case []interface{}:
				yellow.Printf("%s: [%d items]", key, len(v))
			default:
				yellow.Printf("%s: %v", key, v)
			}
		}
		yellow.Print("}")
	}
	fmt.Println()
}

func displayUserMessageCompact(msg *StreamMessage, lineNum int) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		if block.Type == "tool_result" {
			displayToolResultCompact(&block, lineNum)
		}
	}
}

func displayToolResultCompact(block *ContentBlock, lineNum int) {
	if block.IsError {
		boldRed.Print("ERR ")
	} else {
		boldMagenta.Print("RES ")
	}
	gray.Printf("L%d ", lineNum)

	contentStr := ""
	switch v := block.Content.(type) {
	case string:
		contentStr = v
	default:
		contentStr = fmt.Sprintf("%v", v)
	}

	// Strip system reminders in non-verbose mode
	if !*verbose {
		contentStr = stripSystemReminders(contentStr)
	}

	// Compact output - single line summary
	contentStr = strings.ReplaceAll(contentStr, "\n", " ")
	if contentStr == "" {
		gray.Println("(no output)")
	} else if len(contentStr) > 100 {
		white.Printf("%.100s...\n", contentStr)
	} else {
		white.Println(contentStr)
	}
}

func displayResultMessageCompact(msg *StreamMessage, lineNum int) {
	if msg.IsError {
		boldRed.Print("FAIL ")
	} else {
		boldBlue.Print("OK ")
	}
	gray.Printf("L%d", lineNum)

	if msg.NumTurns > 0 {
		blue.Printf(" turns=%d", msg.NumTurns)
	}
	if msg.DurationMS > 0 {
		blue.Printf(" %.2fs", float64(msg.DurationMS)/1000.0)
	}
	if msg.TotalCostUSD > 0 {
		blue.Printf(" $%.4f", msg.TotalCostUSD)
	}
	if msg.Usage != nil {
		blue.Printf(" in=%d out=%d", msg.Usage.InputTokens, msg.Usage.OutputTokens)
	}
	fmt.Println()
}

// ============================================================================
// MINIMAL STYLE FORMATTERS
// ============================================================================

func displayMessageMinimal(msg *StreamMessage, lineNum int) {
	switch msg.Type {
	case "system":
		displaySystemMessageMinimal(msg, lineNum)
	case "assistant":
		displayAssistantMessageMinimal(msg, lineNum)
	case "user":
		displayUserMessageMinimal(msg, lineNum)
	case "result":
		displayResultMessageMinimal(msg, lineNum)
	}
}

func displaySystemMessageMinimal(msg *StreamMessage, lineNum int) {
	boldCyan.Printf("SYSTEM")
	if msg.Subtype != "" {
		cyan.Printf(" [%s]", msg.Subtype)
	}
	gray.Printf(" (line %d)\n", lineNum)

	if msg.CWD != "" {
		cyan.Printf("  Working Directory: %s\n", msg.CWD)
	}
	if msg.Model != "" {
		cyan.Printf("  Model: %s\n", msg.Model)
	}
	if msg.ClaudeCodeVersion != "" {
		cyan.Printf("  Claude Code: v%s\n", msg.ClaudeCodeVersion)
	}
	if len(msg.Tools) > 0 {
		cyan.Printf("  Tools: %d available\n", len(msg.Tools))
	}
	fmt.Println()
}

func displayAssistantMessageMinimal(msg *StreamMessage, lineNum int) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	// Group consecutive text blocks
	var textBlocks []string
	var toolUses []ContentBlock

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
		boldGreen.Printf("ASSISTANT")
		gray.Printf(" (line %d)\n", lineNum)

		for _, text := range textBlocks {
			white.Printf("  %s\n", text)
		}

		if *verbose && msg.Message.Usage != nil {
			gray.Printf("  Tokens: in=%d out=%d", msg.Message.Usage.InputTokens, msg.Message.Usage.OutputTokens)
			if msg.Message.Usage.CacheReadInputTokens > 0 {
				gray.Printf(" cache_read=%d", msg.Message.Usage.CacheReadInputTokens)
			}
			if msg.Message.Usage.CacheCreationInputTokens > 0 {
				gray.Printf(" cache_create=%d", msg.Message.Usage.CacheCreationInputTokens)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Display tool uses
	for _, tool := range toolUses {
		displayToolUseMinimal(&tool, lineNum)
	}
}

func displayToolUseMinimal(tool *ContentBlock, lineNum int) {
	boldYellow.Printf("TOOL: %s", tool.Name)
	gray.Printf(" (line %d)\n", lineNum)

	if *verbose {
		yellow.Printf("  ID: %s\n", tool.ID)
	}

	if tool.Input != nil {
		yellow.Println("  Input:")
		for key, value := range tool.Input {
			yellow.Printf("    %s: ", key)

			switch v := value.(type) {
			case string:
				if len(v) > 300 {
					white.Printf("%s ... (%d chars omitted) ... %s\n", v[:200], len(v)-300, v[len(v)-100:])
				} else {
					white.Println(v)
				}
			case []interface{}:
				if tool.Name == "TodoWrite" && key == "todos" {
					white.Println()
					displayTodosMinimal(v)
				} else {
					white.Printf("[%d items]\n", len(v))
				}
			case map[string]interface{}:
				white.Println("{...}")
			default:
				white.Printf("%v\n", v)
			}
		}
	}
	fmt.Println()
}

func displayTodosMinimal(todos []interface{}) {
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
			statusIcon = green.Sprint("✓")
		case "in_progress":
			statusIcon = yellow.Sprint("→")
		case "pending":
			statusIcon = gray.Sprint("○")
		default:
			statusIcon = gray.Sprint("-")
		}

		yellow.Printf("      %s %s\n", statusIcon, content)
	}
}

func displayUserMessageMinimal(msg *StreamMessage, lineNum int) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		if block.Type == "tool_result" {
			displayToolResultMinimal(&block, lineNum)
		}
	}
}

func displayToolResultMinimal(block *ContentBlock, lineNum int) {
	if block.IsError {
		boldRed.Printf("TOOL RESULT ERROR")
		gray.Printf(" (line %d)\n", lineNum)

		if *verbose {
			red.Printf("  Tool ID: %s\n", block.ToolUseID)
		}

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		// Strip system reminders in non-verbose mode
		if !*verbose {
			contentStr = stripSystemReminders(contentStr)
		}

		white.Printf("  %s\n", contentStr)
	} else {
		boldMagenta.Printf("TOOL RESULT")
		gray.Printf(" (line %d)\n", lineNum)

		if *verbose {
			gray.Printf("  Tool ID: %s\n", block.ToolUseID)
		}

		contentStr := ""
		switch v := block.Content.(type) {
		case string:
			contentStr = v
		default:
			contentStr = fmt.Sprintf("%v", v)
		}

		// Strip system reminders in non-verbose mode
		if !*verbose {
			contentStr = stripSystemReminders(contentStr)
		}

		if contentStr == "" {
			gray.Println("  (no output)")
		} else {
			lines := strings.Split(contentStr, "\n")
			firstLines := 20
			lastLines := 20
			totalLines := len(lines)

			if totalLines <= firstLines+lastLines {
				for _, line := range lines {
					white.Printf("  %s\n", line)
				}
			} else {
				for i := 0; i < firstLines; i++ {
					white.Printf("  %s\n", lines[i])
				}
				gray.Printf("  ... (%d more lines) ...\n", totalLines-firstLines-lastLines)
				for i := totalLines - lastLines; i < totalLines; i++ {
					white.Printf("  %s\n", lines[i])
				}
			}
		}
	}
	fmt.Println()
}

func displayResultMessageMinimal(msg *StreamMessage, lineNum int) {
	if msg.IsError {
		boldRed.Printf("RESULT: ERROR")
	} else {
		boldBlue.Printf("RESULT: SUCCESS")
	}
	gray.Printf(" (line %d)\n", lineNum)

	if msg.NumTurns > 0 {
		blue.Printf("  Turns: %d\n", msg.NumTurns)
	}
	if msg.DurationMS > 0 {
		blue.Printf("  Duration: %.2fs", float64(msg.DurationMS)/1000.0)
		if msg.DurationAPIMS > 0 {
			blue.Printf(" (API: %.2fs)", float64(msg.DurationAPIMS)/1000.0)
		}
		blue.Println()
	}
	if msg.TotalCostUSD > 0 {
		blue.Printf("  Cost: $%.4f\n", msg.TotalCostUSD)
	}

	if msg.Usage != nil {
		blue.Printf("  Tokens: in=%d out=%d", msg.Usage.InputTokens, msg.Usage.OutputTokens)
		if msg.Usage.CacheReadInputTokens > 0 {
			blue.Printf(" cache_read=%d", msg.Usage.CacheReadInputTokens)
		}
		if msg.Usage.CacheCreationInputTokens > 0 {
			blue.Printf(" cache_create=%d", msg.Usage.CacheCreationInputTokens)
		}
		blue.Println()
	}

	if *verbose && msg.ModelUsage != nil && len(msg.ModelUsage) > 0 {
		blue.Println()
		blue.Println("  Model Usage:")
		for model, usageData := range msg.ModelUsage {
			blue.Printf("    %s:\n", model)
			if usageMap, ok := usageData.(map[string]interface{}); ok {
				if inputTokens, ok := usageMap["inputTokens"].(float64); ok {
					blue.Printf("      Input: %.0f tokens\n", inputTokens)
				}
				if outputTokens, ok := usageMap["outputTokens"].(float64); ok {
					blue.Printf("      Output: %.0f tokens\n", outputTokens)
				}
				if cost, ok := usageMap["costUSD"].(float64); ok {
					blue.Printf("      Cost: $%.4f\n", cost)
				}
			}
		}
	}

	if len(msg.PermissionDenials) > 0 {
		fmt.Println()
		red.Printf("  Permission Denials: %d\n", len(msg.PermissionDenials))
		if *verbose {
			for i, denial := range msg.PermissionDenials {
				red.Printf("    [%d] %v\n", i+1, denial)
			}
		}
	}

	if msg.Result != "" {
		fmt.Println()
		lines := strings.Split(msg.Result, "\n")
		for _, line := range lines {
			white.Printf("  %s\n", line)
		}
	}

	fmt.Println()
}

// ============================================================================
// PLAIN STYLE FORMATTERS (No colors, no boxes)
// ============================================================================

func displayMessagePlain(msg *StreamMessage, lineNum int) {
	switch msg.Type {
	case "system":
		displaySystemMessagePlain(msg, lineNum)
	case "assistant":
		displayAssistantMessagePlain(msg, lineNum)
	case "user":
		displayUserMessagePlain(msg, lineNum)
	case "result":
		displayResultMessagePlain(msg, lineNum)
	}
}

func displaySystemMessagePlain(msg *StreamMessage, lineNum int) {
	fmt.Printf("SYSTEM")
	if msg.Subtype != "" {
		fmt.Printf(" [%s]", msg.Subtype)
	}
	fmt.Printf(" (line %d)\n", lineNum)

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

func displayAssistantMessagePlain(msg *StreamMessage, lineNum int) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	// Group consecutive text blocks
	var textBlocks []string
	var toolUses []ContentBlock

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
		fmt.Printf("ASSISTANT (line %d)\n", lineNum)

		for _, text := range textBlocks {
			fmt.Printf("  %s\n", text)
		}

		if *verbose && msg.Message.Usage != nil {
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
		displayToolUsePlain(&tool, lineNum)
	}
}

func displayToolUsePlain(tool *ContentBlock, lineNum int) {
	fmt.Printf("TOOL: %s (line %d)\n", tool.Name, lineNum)

	if *verbose {
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
					displayTodosPlain(v)
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

func displayTodosPlain(todos []interface{}) {
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

func displayUserMessagePlain(msg *StreamMessage, lineNum int) {
	if msg.Message == nil || len(msg.Message.Content) == 0 {
		return
	}

	for _, block := range msg.Message.Content {
		if block.Type == "tool_result" {
			displayToolResultPlain(&block, lineNum)
		}
	}
}

func displayToolResultPlain(block *ContentBlock, lineNum int) {
	if block.IsError {
		fmt.Printf("TOOL RESULT ERROR (line %d)\n", lineNum)

		if *verbose {
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
		if !*verbose {
			contentStr = stripSystemReminders(contentStr)
		}

		fmt.Printf("  %s\n", contentStr)
	} else {
		fmt.Printf("TOOL RESULT (line %d)\n", lineNum)

		if *verbose {
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
		if !*verbose {
			contentStr = stripSystemReminders(contentStr)
		}

		if contentStr == "" {
			fmt.Println("  (no output)")
		} else {
			lines := strings.Split(contentStr, "\n")
			firstLines := 20
			lastLines := 20
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

func displayResultMessagePlain(msg *StreamMessage, lineNum int) {
	if msg.IsError {
		fmt.Printf("RESULT: ERROR (line %d)\n", lineNum)
	} else {
		fmt.Printf("RESULT: SUCCESS (line %d)\n", lineNum)
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

	if *verbose && msg.ModelUsage != nil && len(msg.ModelUsage) > 0 {
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
		if *verbose {
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
