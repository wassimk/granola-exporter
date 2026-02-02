package exporter

import (
	"regexp"
	"strings"
)

// transcriptEntryRegex matches transcript entries in the format: **Speaker:** text
// Matches text until double newline, single newline at end, or end of string
var transcriptEntryRegex = regexp.MustCompile(`\*\*(\w+):\*\* (.+?)(?:\n\n|\n$|$)`)

// ExtractTranscriptFromMarkdown extracts transcript entries from an existing markdown file.
// Returns nil if no transcript section exists or if the section is empty.
func ExtractTranscriptFromMarkdown(content string) []TranscriptEntry {
	// Check if transcript section exists
	if !strings.Contains(content, "## Transcript") {
		return nil
	}

	// Split on the transcript header
	parts := strings.SplitN(content, "## Transcript", 2)
	if len(parts) < 2 {
		return nil
	}

	transcriptSection := parts[1]

	// Find all transcript entries
	matches := transcriptEntryRegex.FindAllStringSubmatch(transcriptSection, -1)
	if len(matches) == 0 {
		return nil
	}

	var entries []TranscriptEntry
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		speaker := match[1]
		text := strings.TrimSpace(match[2])

		if text == "" {
			continue
		}

		source := SpeakerToSource(speaker)
		entries = append(entries, TranscriptEntry{
			Text:   text,
			Source: source,
		})
	}

	if len(entries) == 0 {
		return nil
	}

	return entries
}

// SpeakerToSource maps a speaker label back to a source.
func SpeakerToSource(speaker string) string {
	switch speaker {
	case "Me":
		return "microphone"
	case "Them":
		return "system"
	default:
		return strings.ToLower(speaker)
	}
}
