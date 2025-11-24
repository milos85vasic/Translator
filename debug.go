package main

import (
	"fmt"
	"log"
	"digital.vasic.translator/pkg/format"
	"digital.vasic.translator/pkg/markdown"
	"os"
)

func main() {
	// Test creating an EPUB using the MarkdownToEPUBConverter
	converter := markdown.NewMarkdownToEPUBConverter()
	
	// Create a simple markdown file
	mdContent := `# Test Book

This is a test book with some content.

## Chapter 1

This is the first chapter.

## Chapter 2

This is the second chapter.
`
	
	// Write markdown to temp file
	tmpDir := "/tmp/test_epub"
	os.MkdirAll(tmpDir, 0755)
	mdPath := tmpDir + "/test.md"
	os.WriteFile(mdPath, []byte(mdContent), 0644)
	
	// Convert to EPUB
	epubPath := tmpDir + "/test.epub"
	err := converter.ConvertMarkdownToEPUB(mdPath, epubPath)
	if err != nil {
		log.Fatalf("Failed to convert to EPUB: %v", err)
	}
	
	// Check what format is detected
	detector := format.NewDetector()
	detectedFormat, err := detector.DetectFile(epubPath)
	if err != nil {
		log.Fatalf("Failed to detect format: %v", err)
	}
	
	fmt.Printf("Detected format: %s\n", detectedFormat)
	
	// Clean up
	os.RemoveAll(tmpDir)
}