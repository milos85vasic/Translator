package main

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug_epub.go <epub_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	
	r, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()

	fmt.Println("Files in EPUB:")
	generalIndicators := []string{
		"mimetype",
		"OEBPS",
		"META-INF",
	}

	azw3Specific := []string{
		"kindle:embed",
		"amzn-eastock",
		"kindle-fonts",
		"kindle:enclosure",
		"kindle:meta",
	}

	hasGeneral := 0
	hasSpecific := false

	for _, f := range r.File {
		fmt.Printf("- %s\n", f.Name)
		
		// Check for general indicators
		for _, indicator := range generalIndicators {
			if strings.Contains(f.Name, indicator) {
				hasGeneral++
				fmt.Printf("  General indicator found: %s\n", indicator)
				break
			}
		}
		
		// Check for AZW3-specific indicators
		for _, specific := range azw3Specific {
			if strings.Contains(f.Name, specific) {
				hasSpecific = true
				fmt.Printf("  AZW3-specific indicator found: %s\n", specific)
				break
			}
		}
	}

	fmt.Printf("\nAnalysis:\n")
	fmt.Printf("- General indicators: %d\n", hasGeneral)
	fmt.Printf("- AZW3-specific: %v\n", hasSpecific)
	fmt.Printf("- Would be detected as AZW3: %v\n", hasGeneral >= 2 && hasSpecific)
}