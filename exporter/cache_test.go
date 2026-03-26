package exporter

import (
	"testing"
)

func TestParseCache(t *testing.T) {
	t.Run("parses valid nested JSON structure", func(t *testing.T) {
		data := []byte(`{
			"cache": "{\"state\":{\"documents\":{\"doc1\":{\"id\":\"doc1\",\"title\":\"Test\",\"created_at\":\"2026-01-21T10:00:00Z\",\"notes_markdown\":\"# Notes\",\"notes_plain\":\"Notes\"}},\"transcripts\":{\"doc1\":[{\"id\":\"t1\",\"document_id\":\"doc1\",\"text\":\"Hello\",\"source\":\"microphone\"}]}}}"
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(state.Documents) != 1 {
			t.Errorf("Expected 1 document, got %d", len(state.Documents))
		}
		if state.Documents["doc1"].Title != "Test" {
			t.Errorf("Expected title 'Test', got %s", state.Documents["doc1"].Title)
		}
		if len(state.Transcripts) != 1 {
			t.Errorf("Expected 1 transcript, got %d", len(state.Transcripts))
		}
		if len(state.Transcripts["doc1"]) != 1 {
			t.Errorf("Expected 1 transcript entry, got %d", len(state.Transcripts["doc1"]))
		}
	})

	t.Run("handles missing documents key", func(t *testing.T) {
		data := []byte(`{
			"cache": "{\"state\":{\"transcripts\":{}}}"
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if state.Documents == nil {
			t.Error("Documents map should be initialized")
		}
		if len(state.Documents) != 0 {
			t.Error("Documents map should be empty")
		}
	})

	t.Run("handles missing transcripts key", func(t *testing.T) {
		data := []byte(`{
			"cache": "{\"state\":{\"documents\":{}}}"
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if state.Transcripts == nil {
			t.Error("Transcripts map should be initialized")
		}
		if len(state.Transcripts) != 0 {
			t.Error("Transcripts map should be empty")
		}
	})

	t.Run("handles empty documents", func(t *testing.T) {
		data := []byte(`{
			"cache": "{\"state\":{\"documents\":{},\"transcripts\":{}}}"
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(state.Documents) != 0 {
			t.Errorf("Expected 0 documents, got %d", len(state.Documents))
		}
	})

	t.Run("handles empty transcripts", func(t *testing.T) {
		data := []byte(`{
			"cache": "{\"state\":{\"documents\":{},\"transcripts\":{}}}"
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(state.Transcripts) != 0 {
			t.Errorf("Expected 0 transcripts, got %d", len(state.Transcripts))
		}
	})

	t.Run("parses cache as direct object (v6+ format)", func(t *testing.T) {
		data := []byte(`{
			"cache": {"state":{"documents":{"doc1":{"id":"doc1","title":"Test","created_at":"2026-01-21T10:00:00Z","notes_markdown":"# Notes","notes_plain":"Notes"}},"transcripts":{"doc1":[{"id":"t1","document_id":"doc1","text":"Hello","source":"microphone"}]}}}
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(state.Documents) != 1 {
			t.Errorf("Expected 1 document, got %d", len(state.Documents))
		}
		if state.Documents["doc1"].Title != "Test" {
			t.Errorf("Expected title 'Test', got %s", state.Documents["doc1"].Title)
		}
		if len(state.Transcripts) != 1 {
			t.Errorf("Expected 1 transcript, got %d", len(state.Transcripts))
		}
	})

	t.Run("parses sharedDocuments", func(t *testing.T) {
		data := []byte(`{
			"cache": {"state":{"documents":{"doc1":{"id":"doc1","title":"My Meeting","created_at":"2026-01-21T10:00:00Z","notes_markdown":"# Notes"}},"sharedDocuments":{"doc2":{"id":"doc2","title":"Shared Meeting","created_at":"2026-01-21T11:00:00Z","notes_markdown":"# Shared Notes"}},"transcripts":{}}}
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(state.Documents) != 1 {
			t.Errorf("Expected 1 document, got %d", len(state.Documents))
		}
		if len(state.SharedDocuments) != 1 {
			t.Errorf("Expected 1 shared document, got %d", len(state.SharedDocuments))
		}
		if state.SharedDocuments["doc2"].Title != "Shared Meeting" {
			t.Errorf("Expected shared doc title 'Shared Meeting', got %s", state.SharedDocuments["doc2"].Title)
		}
	})

	t.Run("handles missing sharedDocuments key", func(t *testing.T) {
		data := []byte(`{
			"cache": {"state":{"documents":{},"transcripts":{}}}
		}`)

		state, err := ParseCache(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if state.SharedDocuments == nil {
			t.Error("SharedDocuments map should be initialized")
		}
		if len(state.SharedDocuments) != 0 {
			t.Error("SharedDocuments map should be empty")
		}
	})

	t.Run("returns error for invalid outer JSON", func(t *testing.T) {
		data := []byte(`not valid json`)

		_, err := ParseCache(data)
		if err == nil {
			t.Error("Expected error for invalid outer JSON")
		}
	})

	t.Run("returns error for invalid inner JSON", func(t *testing.T) {
		data := []byte(`{
			"cache": "not valid json"
		}`)

		_, err := ParseCache(data)
		if err == nil {
			t.Error("Expected error for invalid inner JSON")
		}
	})

	t.Run("returns error for empty cache field", func(t *testing.T) {
		data := []byte(`{
			"cache": ""
		}`)

		_, err := ParseCache(data)
		if err == nil {
			t.Error("Expected error for empty cache field")
		}
	})
}

func TestAllDocuments(t *testing.T) {
	t.Run("merges owned and shared documents", func(t *testing.T) {
		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Owned"},
			},
			SharedDocuments: map[string]Document{
				"doc2": {ID: "doc2", Title: "Shared"},
			},
		}

		all := state.AllDocuments()
		if len(all) != 2 {
			t.Errorf("Expected 2 documents, got %d", len(all))
		}
		if all["doc1"].Title != "Owned" {
			t.Error("Expected owned document in result")
		}
		if all["doc2"].Title != "Shared" {
			t.Error("Expected shared document in result")
		}
	})

	t.Run("owned documents take precedence over shared", func(t *testing.T) {
		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Owned Version"},
			},
			SharedDocuments: map[string]Document{
				"doc1": {ID: "doc1", Title: "Shared Version"},
			},
		}

		all := state.AllDocuments()
		if len(all) != 1 {
			t.Errorf("Expected 1 document (deduped), got %d", len(all))
		}
		if all["doc1"].Title != "Owned Version" {
			t.Errorf("Expected owned version to take precedence, got %q", all["doc1"].Title)
		}
	})

	t.Run("handles nil shared documents", func(t *testing.T) {
		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Owned"},
			},
		}

		all := state.AllDocuments()
		if len(all) != 1 {
			t.Errorf("Expected 1 document, got %d", len(all))
		}
	})
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		{"/path/to/cache-v1.json", 1},
		{"/path/to/cache-v3.json", 3},
		{"/path/to/cache-v10.json", 10},
		{"/path/to/cache-v123.json", 123},
		{"/path/to/other.json", 0},
		{"cache-v5.json", 5},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := extractVersion(tt.path)
			if result != tt.expected {
				t.Errorf("extractVersion(%q) = %d, want %d", tt.path, result, tt.expected)
			}
		})
	}
}
