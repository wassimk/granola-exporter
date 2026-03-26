package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// cacheVersionRegex extracts the version number from cache-vN.json filenames.
var cacheVersionRegex = regexp.MustCompile(`cache-v(\d+)\.json$`)

// FindCacheFile finds the latest Granola cache file.
// Returns the path to the cache file with the highest version number.
func FindCacheFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	granolaDir := filepath.Join(homeDir, "Library", "Application Support", "Granola")

	// Find all cache-v*.json files
	pattern := filepath.Join(granolaDir, "cache-v*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for cache files: %w", err)
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no Granola cache files found in %s\nExpected to find cache-v*.json files", granolaDir)
	}

	// Sort by version number (extract number from cache-vN.json)
	sort.Slice(matches, func(i, j int) bool {
		vi := extractVersion(matches[i])
		vj := extractVersion(matches[j])
		return vi < vj
	})

	// Use the highest version
	latestCache := matches[len(matches)-1]

	// Verify it exists
	if _, err := os.Stat(latestCache); os.IsNotExist(err) {
		return "", fmt.Errorf("cache file not found: %s", latestCache)
	}

	return latestCache, nil
}

// extractVersion extracts the version number from a cache filename.
func extractVersion(path string) int {
	matches := cacheVersionRegex.FindStringSubmatch(path)
	if len(matches) < 2 {
		return 0
	}
	v, _ := strconv.Atoi(matches[1])
	return v
}

// outerCacheRaw is used to detect whether the "cache" field is a string or object.
type outerCacheRaw struct {
	Cache json.RawMessage `json:"cache"`
}

// innerCache represents the inner JSON structure containing the state.
type innerCache struct {
	State CacheState `json:"state"`
}

// LoadCache loads and parses the Granola cache from a file path.
func LoadCache(path string) (*CacheState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	return ParseCache(data)
}

// ParseCache parses the Granola cache from raw JSON bytes.
// Supports two formats:
//   - Legacy (cache-v5 and earlier): "cache" is a JSON string containing nested JSON
//   - Current (cache-v6+): "cache" is a direct JSON object
func ParseCache(data []byte) (*CacheState, error) {
	var outer outerCacheRaw
	if err := json.Unmarshal(data, &outer); err != nil {
		return nil, fmt.Errorf("failed to parse cache JSON: %w", err)
	}

	if len(outer.Cache) == 0 {
		return nil, fmt.Errorf("cache field is empty")
	}

	var innerData []byte

	// If cache is a JSON string, unwrap it; otherwise use it directly
	if outer.Cache[0] == '"' {
		var cacheStr string
		if err := json.Unmarshal(outer.Cache, &cacheStr); err != nil {
			return nil, fmt.Errorf("failed to parse cache string: %w", err)
		}
		if cacheStr == "" {
			return nil, fmt.Errorf("cache field is empty")
		}
		innerData = []byte(cacheStr)
	} else {
		innerData = outer.Cache
	}

	var inner innerCache
	if err := json.Unmarshal(innerData, &inner); err != nil {
		return nil, fmt.Errorf("failed to parse inner cache JSON: %w", err)
	}

	// Ensure maps are initialized
	if inner.State.Documents == nil {
		inner.State.Documents = make(map[string]Document)
	}
	if inner.State.SharedDocuments == nil {
		inner.State.SharedDocuments = make(map[string]Document)
	}
	if inner.State.Transcripts == nil {
		inner.State.Transcripts = make(map[string][]TranscriptEntry)
	}

	return &inner.State, nil
}

// GetCacheSize returns the size of a file in bytes.
func GetCacheSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
