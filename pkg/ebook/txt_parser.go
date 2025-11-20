package ebook

import (
	"bufio"
	"digital.vasic.translator/pkg/format"
	"os"
	"strings"
)

// TXTParser implements Parser for plain text format
type TXTParser struct{}

// NewTXTParser creates a new TXT parser
func NewTXTParser() *TXTParser {
	return &TXTParser{}
}

// Parse parses a plain text file into universal Book structure
func (p *TXTParser) Parse(filename string) (*Book, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	book := &Book{
		Metadata: Metadata{
			Title: filename,
		},
		Chapters: make([]Chapter, 0),
		Format:   format.FormatTXT,
	}

	// Read content
	scanner := bufio.NewScanner(file)
	var content strings.Builder

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Create single chapter with all content
	chapter := Chapter{
		Title: "Content",
		Sections: []Section{
			{
				Content: content.String(),
			},
		},
	}

	book.Chapters = append(book.Chapters, chapter)

	return book, nil
}

// GetFormat returns the format
func (p *TXTParser) GetFormat() format.Format {
	return format.FormatTXT
}
