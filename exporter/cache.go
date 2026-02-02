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

// outerCache represents the outer JSON structure of the cache file.
type outerCache struct {
	Cache string `json:"cache"`
}

// innerCache represents the inner JSON structure (parsed from the cache string).
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
// The cache has a nested structure: outer JSON contains a "cache" key with a JSON string value.
func ParseCache(data []byte) (*CacheState, error) {
	// Parse outer JSON
	var outer outerCache
	if err := json.Unmarshal(data, &outer); err != nil {
		return nil, fmt.Errorf("failed to parse outer cache JSON: %w", err)
	}

	if outer.Cache == "" {
		return nil, fmt.Errorf("cache field is empty")
	}

	// Parse inner JSON (the cache string)
	var inner innerCache
	if err := json.Unmarshal([]byte(outer.Cache), &inner); err != nil {
		return nil, fmt.Errorf("failed to parse inner cache JSON: %w", err)
	}

	// Ensure maps are initialized
	if inner.State.Documents == nil {
		inner.State.Documents = make(map[string]Document)
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
