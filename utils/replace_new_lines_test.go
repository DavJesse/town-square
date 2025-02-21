package utils

import (
	"html/template"
	"testing"
)

func TestReplaceNewlines(t *testing.T) {
	table := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{
			name:     "No newlines",
			input:    "Hello, World!",
			expected: template.HTML("Hello, World!"),
		},
		{
			name:     "Single newline",
			input:    "Hello,\nWorld!",
			expected: template.HTML("Hello,<br>World!"),
		},
		{
			name:     "Multiple newlines",
			input:    "Hello,\nWorld!\nThis is a test.",
			expected: template.HTML("Hello,<br>World!<br>This is a test."),
		},
		{
			name:     "Newlines at the start",
			input:    "\nHello, World!",
			expected: template.HTML("<br>Hello, World!"),
		},
		{
			name:     "Newlines at the end",
			input:    "Hello, World!\n",
			expected: template.HTML("Hello, World!<br>"),
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			result := ReplaceNewlines(entry.input)
			if result != entry.expected {
				t.Errorf("For input '%s', expected '%s' but got '%s'", entry.input, entry.expected, result)
			}
		})
	}
}
