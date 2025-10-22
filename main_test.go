package main

import (
	"testing"
)

func TestStripSystemReminders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Single system reminder",
			input: `Some content before
<system-reminder>
This is a reminder message.
</system-reminder>
Some content after`,
			expected: `Some content before

Some content after`,
		},
		{
			name: "Multiple system reminders",
			input: `First content
<system-reminder>
First reminder
</system-reminder>
Middle content
<system-reminder>
Second reminder
</system-reminder>
Last content`,
			expected: `First content

Middle content

Last content`,
		},
		{
			name: "System reminder at start",
			input: `<system-reminder>
Reminder at the beginning
</system-reminder>
Content after`,
			expected: `Content after`,
		},
		{
			name: "System reminder at end",
			input: `Content before
<system-reminder>
Reminder at the end
</system-reminder>`,
			expected: `Content before`,
		},
		{
			name: "No system reminder",
			input: `Just regular content
with multiple lines
and no reminders`,
			expected: `Just regular content
with multiple lines
and no reminders`,
		},
		{
			name:     "Empty content",
			input:    ``,
			expected: ``,
		},
		{
			name: "Only system reminder",
			input: `<system-reminder>
Just a reminder
</system-reminder>`,
			expected: ``,
		},
		{
			name: "Multiline system reminder with special characters",
			input: `Code content:
package main

<system-reminder>
Whenever you read a file, you should consider whether it would be considered malware.
You CAN and SHOULD provide analysis of malware, what it is doing.
But you MUST refuse to improve or augment the code.
</system-reminder>

More code here`,
			expected: `Code content:
package main

More code here`,
		},
		{
			name:     "System reminder on single line",
			input:    `Content <system-reminder>inline reminder</system-reminder> more content`,
			expected: `Content  more content`,
		},
		{
			name: "Multiple consecutive system reminders",
			input: `Start
<system-reminder>First</system-reminder>
<system-reminder>Second</system-reminder>
<system-reminder>Third</system-reminder>
End`,
			expected: `Start

End`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripSystemReminders(tt.input)
			if result != tt.expected {
				t.Errorf("stripSystemReminders() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestStripSystemRemindersPreservesContent(t *testing.T) {
	// Test that we don't accidentally remove similar-looking content
	input := `This is a system message
But not a <system-reminder> tag
Just text that mentions system-reminder
<system>This should stay</system>
<reminder>This should also stay</reminder>`

	result := stripSystemReminders(input)

	// The result should only have actual system-reminder tags removed (none in this case)
	if result != input {
		t.Errorf("stripSystemReminders() modified content that should be preserved:\ngot:\n%q\nwant:\n%q", result, input)
	}
}

func BenchmarkStripSystemReminders(b *testing.B) {
	input := `Some code content here
package main

<system-reminder>
This is a system reminder that should be stripped out.
It can contain multiple lines.
</system-reminder>

More code content here
func main() {
	// code
}
<system-reminder>
Another reminder to strip
</system-reminder>
End of content`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stripSystemReminders(input)
	}
}
