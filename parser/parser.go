package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// MaxBufferCapacity is the maximum buffer size for handling large JSON lines (10MB)
const MaxBufferCapacity = 10 * 1024 * 1024

// FirstLines is the number of lines to show at the start of truncated output
const FirstLines = 20

// LastLines is the number of lines to show at the end of truncated output
const LastLines = 20

var (
	systemReminderRe = regexp.MustCompile(`(?s)<system-reminder>.*?</system-reminder>`)
	extraBlankRe     = regexp.MustCompile(`\n\n\n+`)
)

// StripSystemReminders removes <system-reminder>...</system-reminder> tags from content
func StripSystemReminders(content string) string {
	result := systemReminderRe.ReplaceAllString(content, "")
	result = extraBlankRe.ReplaceAllString(result, "\n\n")
	return strings.TrimSpace(result)
}

// ToolResultText flattens a tool_result content field to display text.
// Content is either a string or an array of blocks (text, image, tool_reference, ...).
func ToolResultText(content interface{}) string {
	switch v := content.(type) {
	case nil:
		return ""
	case string:
		return v
	case []interface{}:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			block, ok := item.(map[string]interface{})
			if !ok {
				parts = append(parts, fmt.Sprint(item))
				continue
			}
			switch block["type"] {
			case "text":
				if text, _ := block["text"].(string); text != "" {
					parts = append(parts, text)
				}
			case "image":
				parts = append(parts, "[image]")
			case "tool_reference":
				parts = append(parts, fmt.Sprintf("[loaded tool: %v]", block["tool_name"]))
			default:
				raw, _ := json.Marshal(block)
				parts = append(parts, string(raw))
			}
		}
		return strings.Join(parts, "\n")
	default:
		raw, _ := json.Marshal(v)
		return string(raw)
	}
}
