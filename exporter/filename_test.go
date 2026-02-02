package exporter

import (
	"testing"
)

func TestSafeFilename(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		dateStr  string
		expected string
	}{
		{
			name:     "basic filename",
			title:    "Team Meeting",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Team Meeting.md",
		},
		{
			name:     "removes unsafe characters",
			title:    `Acme Corp <> Globex :: Weekly Sync`,
			dateStr:  "2025-10-23",
			expected: "2025-10-23_Acme Corp  Globex  Weekly Sync.md",
		},
		{
			name:     "empty title becomes Untitled",
			title:    "",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Untitled.md",
		},
		{
			name:     "None title becomes Untitled",
			title:    "None",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Untitled.md",
		},
		{
			name:     "truncates long titles",
			title:    "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.md",
		},
		{
			name:     "whitespace-only title becomes Untitled",
			title:    "   ",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Untitled.md",
		},
		{
			name:     "trims leading and trailing spaces",
			title:    "  Meeting Title  ",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Meeting Title.md",
		},
		{
			name:     "handles Person <> Person pattern",
			title:    "Alice <> Bob",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Alice  Bob.md",
		},
		{
			name:     "handles colons in titles",
			title:    "Acme Corp <> Globex: Weekly Sync",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Acme Corp  Globex Weekly Sync.md",
		},
		{
			name:     "preserves Unicode emojis",
			title:    "🏠 Personal Commitment",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_🏠 Personal Commitment.md",
		},
		{
			name:     "preserves accented characters",
			title:    "Dennis Schröpfer: Next Steps",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Dennis Schröpfer Next Steps.md",
		},
		{
			name:     "handles slash in titles",
			title:    "EM/PM Sync",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_EMPM Sync.md",
		},
		{
			name:     "title becomes empty after removing unsafe chars",
			title:    "<>:?*",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Untitled.md",
		},
		{
			name:     "handles question marks",
			title:    "Is this the right fit?",
			dateStr:  "2025-01-24",
			expected: "2025-01-24_Is this the right fit.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeFilename(tt.title, tt.dateStr)
			if result != tt.expected {
				t.Errorf("SafeFilename(%q, %q) = %q, want %q", tt.title, tt.dateStr, result, tt.expected)
			}
		})
	}
}

func TestSafeFilenameTruncation(t *testing.T) {
	longTitle := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	result := SafeFilename(longTitle, "2025-01-24")

	// Expected length: 10 (date) + 1 (_) + 100 (truncated title) + 3 (.md) = 114
	if len(result) != 114 {
		t.Errorf("Expected length 114, got %d for result %q", len(result), result)
	}
}
