package exporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExporter(t *testing.T) {
	t.Run("filters documents with transcripts", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Doc with transcript", CreatedAt: "2026-01-21T10:00:00Z"},
				"doc2": {ID: "doc2", Title: "Doc without anything", CreatedAt: "2026-01-21T10:00:00Z"},
			},
			Transcripts: map[string][]TranscriptEntry{
				"doc1": {{Text: "Hello", Source: "microphone"}},
			},
		}

		result, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Written != 1 {
			t.Errorf("Expected 1 written, got %d", result.Written)
		}
	})

	t.Run("filters documents with notes > 10 chars", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Doc with notes", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "This is a long enough note to export"},
				"doc2": {ID: "doc2", Title: "Doc with short notes", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "Short"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		result, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Written != 1 {
			t.Errorf("Expected 1 written, got %d", result.Written)
		}
	})

	t.Run("skips documents with neither notes nor transcript", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Empty doc", CreatedAt: "2026-01-21T10:00:00Z"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		result, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Written != 0 {
			t.Errorf("Expected 0 written, got %d", result.Written)
		}
	})

	t.Run("preserves transcript from existing file when cache lacks it", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		// Create existing file with transcript
		existingContent := `# Existing Meeting
Date: 2026-01-21 10:00
Meeting ID: doc1

---

## Transcript

**Me:** Preserved transcript entry.

`
		existingPath := filepath.Join(tmpDir, "2026-01-21_Existing Meeting.md")
		if err := os.WriteFile(existingPath, []byte(existingContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Export with notes but no transcript in cache
		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Existing Meeting", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "New notes added"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		_, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Read the file and verify transcript was preserved
		content, err := os.ReadFile(existingPath)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(content), "## Transcript") {
			t.Error("Expected transcript section to be preserved")
		}
		if !strings.Contains(string(content), "Preserved transcript entry") {
			t.Error("Expected transcript content to be preserved")
		}
		if !strings.Contains(string(content), "New notes added") {
			t.Error("Expected new notes to be present")
		}
	})

	t.Run("skips writing unchanged files", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Test", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "Some notes here"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		// First export
		result1, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result1.Written != 1 {
			t.Errorf("First export: expected 1 written, got %d", result1.Written)
		}

		// Second export with same content
		result2, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result2.Skipped != 1 {
			t.Errorf("Second export: expected 1 skipped, got %d", result2.Skipped)
		}
		if result2.Written != 0 {
			t.Errorf("Second export: expected 0 written, got %d", result2.Written)
		}
	})

	t.Run("creates output directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "nested", "output", "dir")
		exp := NewExporter(outputDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Test", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "Some notes here"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		_, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify directory was created
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			t.Error("Output directory was not created")
		}
	})
}

func TestDefaultOutputDir(t *testing.T) {
	dir := DefaultOutputDir()
	if !strings.Contains(dir, ".local") || !strings.Contains(dir, "granola-transcripts") {
		t.Errorf("Unexpected default output dir: %s", dir)
	}
}
