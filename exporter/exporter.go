package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportResult holds statistics about an export operation.
type ExportResult struct {
	Written int
	Skipped int
	Empty   int
	Errors  []ExportError
}

// ExportError represents an error that occurred during export.
type ExportError struct {
	DocumentID string
	Title      string
	Error      string
}

// Exporter handles exporting Granola documents to markdown files.
type Exporter struct {
	OutputDir string
}

// NewExporter creates a new Exporter with the given output directory.
func NewExporter(outputDir string) *Exporter {
	return &Exporter{OutputDir: outputDir}
}

// DefaultOutputDir returns the default output directory path.
func DefaultOutputDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".local", "share", "granola-transcripts")
}

// Export exports all exportable documents from the cache state.
func (e *Exporter) Export(state *CacheState, verbose bool) (*ExportResult, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(e.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	result := &ExportResult{}

	// Collect exportable documents
	var exportable []Document
	for _, doc := range state.Documents {
		if doc.HasExportableContent(state.Transcripts) {
			exportable = append(exportable, doc)
		}
	}

	if verbose {
		fmt.Printf("Found %d documents with content to export\n\n", len(exportable))
		fmt.Println("Exporting Granola documents:")
		fmt.Println(strings.Repeat("=", 70))
	}

	// Export each document
	for _, doc := range exportable {
		err := e.exportDocument(&doc, state.Transcripts, result, verbose)
		if err != nil {
			result.Errors = append(result.Errors, ExportError{
				DocumentID: doc.ID,
				Title:      doc.Title,
				Error:      err.Error(),
			})
			if verbose {
				fmt.Printf("✗ Error with %s (%s): %s\n", doc.ID, doc.Title, err.Error())
			}
		}
	}

	if verbose {
		fmt.Printf("\n%s\n", strings.Repeat("=", 70))
	}

	return result, nil
}

func (e *Exporter) exportDocument(doc *Document, transcripts map[string][]TranscriptEntry, result *ExportResult, verbose bool) error {
	title := doc.Title
	if title == "" {
		title = "Untitled"
	}

	// Get date string for filename
	dateStr := FormatDateForFilename(doc.CreatedAt)

	// Get transcript if available
	transcript := transcripts[doc.ID]

	// Check if both notes and transcript are empty
	notes := doc.GetNotes()
	if (notes == "" || strings.TrimSpace(notes) == "") && len(transcript) == 0 {
		result.Empty++
		return nil
	}

	// Generate filename
	filename := SafeFilename(title, dateStr)
	outputPath := filepath.Join(e.OutputDir, filename)

	// If file exists and cache has no transcript, try to preserve transcript from file
	if _, err := os.Stat(outputPath); err == nil && len(transcript) == 0 {
		existingContent, err := os.ReadFile(outputPath)
		if err == nil && strings.Contains(string(existingContent), "## Transcript") {
			transcript = ExtractTranscriptFromMarkdown(string(existingContent))
		}
	}

	// Format content with latest notes and best available transcript
	content := FormatDocumentMarkdown(doc, transcript)

	// Check if file exists and content is identical
	if existingContent, err := os.ReadFile(outputPath); err == nil {
		if string(existingContent) == content {
			result.Skipped++
			return nil
		}
	}

	// Write the file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if verbose {
		// Count words
		wordCount := len(strings.Fields(content))

		// Get file size
		info, _ := os.Stat(outputPath)
		fileSize := info.Size()

		// Describe what was included
		var contentParts []string
		if notes != "" && strings.TrimSpace(notes) != "" {
			contentParts = append(contentParts, "notes")
		}
		if len(transcript) > 0 {
			contentParts = append(contentParts, fmt.Sprintf("transcript (%d entries)", len(transcript)))
		}

		fmt.Printf("✓ %s\n", filename)
		fmt.Printf("  [%s] %s words, %s bytes\n", strings.Join(contentParts, " + "), NumberWithCommas(wordCount), NumberWithCommas(int(fileSize)))
	}

	result.Written++
	return nil
}

// PrintSummary prints a summary of the export result.
func (r *ExportResult) PrintSummary(outputDir string) {
	fmt.Println("\nSummary:")
	fmt.Printf("  Written: %d documents\n", r.Written)
	fmt.Printf("  Skipped (unchanged): %d documents\n", r.Skipped)
	fmt.Printf("  Empty: %d documents\n", r.Empty)
	fmt.Printf("  Errors: %d\n", len(r.Errors))
	fmt.Printf("\nAll documents saved to: %s\n", outputDir)

	if len(r.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, e := range r.Errors {
			fmt.Printf("  %s: %s\n", e.DocumentID, e.Error)
		}
	}
}
