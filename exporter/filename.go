package exporter

import (
	"strings"
)

// unsafeChars contains characters that are unsafe for filenames on most filesystems.
const unsafeChars = `<>:"/\|?*`

// SafeFilename generates a safe filename from a title and date string.
// Format: YYYY-MM-DD_Title.md
func SafeFilename(title, dateStr string) string {
	// Handle nil/empty/"None" titles
	if title == "" || title == "None" || strings.TrimSpace(title) == "" {
		title = "Untitled"
	}

	// Remove unsafe characters
	safeTitle := removeUnsafeChars(title)

	// Trim whitespace from ends
	safeTitle = strings.TrimSpace(safeTitle)

	// If title becomes empty after removing unsafe chars, use "Untitled"
	if safeTitle == "" {
		safeTitle = "Untitled"
	}

	// Truncate to 100 characters
	if len(safeTitle) > 100 {
		safeTitle = safeTitle[:100]
	}

	return dateStr + "_" + safeTitle + ".md"
}

// removeUnsafeChars removes characters that are unsafe for filenames.
func removeUnsafeChars(s string) string {
	var result strings.Builder
	for _, r := range s {
		if !strings.ContainsRune(unsafeChars, r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}
