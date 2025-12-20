package parser

import (
	"regexp"
	"strings"
)

// MaxBufferCapacity is the maximum buffer size for handling large JSON lines (10MB)
const MaxBufferCapacity = 10 * 1024 * 1024

// FirstLines is the number of lines to show at the start of truncated output
const FirstLines = 20

// LastLines is the number of lines to show at the end of truncated output
const LastLines = 20

// StripSystemReminders removes <system-reminder>...</system-reminder> tags from content
func StripSystemReminders(content string) string {
	// Use regex to remove system-reminder tags and their content
	re := regexp.MustCompile(`(?s)<system-reminder>.*?</system-reminder>`)
	result := re.ReplaceAllString(content, "")

	// Clean up any resulting multiple blank lines
	result = regexp.MustCompile(`\n\n\n+`).ReplaceAllString(result, "\n\n")

	// Trim leading/trailing whitespace
	return strings.TrimSpace(result)
}
