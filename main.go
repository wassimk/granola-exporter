package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wassimk/granola-exporter/exporter"
)

// version is set at build time via ldflags
var version = "dev"

func main() {
	// Define flags
	var (
		showVersion bool
		showHelp    bool
		outputDir   string
	)

	flag.BoolVar(&showVersion, "version", false, "Show version number")
	flag.BoolVar(&showVersion, "V", false, "Show version number (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (shorthand)")
	flag.StringVar(&outputDir, "output-dir", "", "Custom output directory")
	flag.StringVar(&outputDir, "o", "", "Custom output directory (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "granola-exporter - Export Granola meeting notes and transcripts to markdown\n\n")
		fmt.Fprintf(os.Stderr, "Usage: granola-exporter [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fmt.Fprintf(os.Stderr, "  -h, --help         Show this help message\n")
		fmt.Fprintf(os.Stderr, "  -V, --version      Show version number\n")
		fmt.Fprintf(os.Stderr, "  -o, --output-dir   Custom output directory (default: ~/.local/share/granola-transcripts)\n")
	}

	flag.Parse()

	// Handle --help
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Handle --version
	if showVersion {
		fmt.Printf("granola-exporter %s\n", version)
		os.Exit(0)
	}

	// Set default output directory
	if outputDir == "" {
		outputDir = exporter.DefaultOutputDir()
	}

	// Run the export
	if err := run(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(outputDir string) error {
	// Find cache file
	cachePath, err := exporter.FindCacheFile()
	if err != nil {
		return err
	}

	fmt.Printf("Loading cache from: %s\n", cachePath)

	// Get cache size
	cacheSize, err := exporter.GetCacheSize(cachePath)
	if err != nil {
		return fmt.Errorf("failed to get cache size: %w", err)
	}
	cacheSizeMB := float64(cacheSize) / 1024.0 / 1024.0
	fmt.Printf("Cache size: %.1f MB\n\n", cacheSizeMB)

	// Load and parse cache
	fmt.Println("Parsing cache...")
	state, err := exporter.LoadCache(cachePath)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d documents\n", len(state.Documents))
	fmt.Printf("Found %d transcripts\n\n", len(state.Transcripts))

	// Export documents
	exp := exporter.NewExporter(outputDir)
	result, err := exp.Export(state, true)
	if err != nil {
		return err
	}

	// Print summary
	result.PrintSummary(outputDir)

	return nil
}
