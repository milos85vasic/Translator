package ebook

import (
	"os"
	"strings"
	"testing"

	"digital.vasic.translator/pkg/format"
)

func TestNewUniversalParser(t *testing.T) {
	parser := NewUniversalParser()

	if parser == nil {
		t.Fatal("NewUniversalParser returned nil")
	}

	if parser.detector == nil {
		t.Error("Parser detector is nil")
	}

	if parser.parsers == nil {
		t.Error("Parser parsers map is nil")
	}

	// Check that all expected parsers are registered
	expectedFormats := []format.Format{
		format.FormatFB2,
		format.FormatEPUB,
		format.FormatTXT,
		format.FormatHTML,
	}

	for _, fmt := range expectedFormats {
		if _, exists := parser.parsers[fmt]; !exists {
			t.Errorf("Parser for format %s not registered", fmt)
		}
	}
}

func TestUniversalParser_GetSupportedFormats(t *testing.T) {
	parser := NewUniversalParser()
	formats := parser.GetSupportedFormats()

	if len(formats) == 0 {
		t.Error("GetSupportedFormats returned empty list")
	}

	// Check that our expected formats are included
	expectedFormats := []format.Format{
		format.FormatFB2,
		format.FormatEPUB,
		format.FormatTXT,
		format.FormatHTML,
	}

	for _, expected := range expectedFormats {
		found := false
		for _, fmt := range formats {
			if fmt == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected format %s not found in supported formats", expected)
		}
	}
}

func TestUniversalParser_Parse_UnsupportedFormat(t *testing.T) {
	parser := NewUniversalParser()

	// Create a temporary file with unknown extension and binary content
	tmpFile, err := os.CreateTemp("", "test.unknown")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write some binary data that won't match any known format
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}
	_, err = tmpFile.Write(binaryData)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close() // Close before parsing

	_, err = parser.Parse(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for unknown format")
	}

	if !strings.Contains(err.Error(), "unknown or unsupported format") {
		t.Errorf("Expected 'unknown or unsupported format' error, got: %v", err)
	}
}

func TestConvertBook(t *testing.T) {
	originalBook := &Book{
		Metadata: Metadata{
			Title:    "Test Book",
			Authors:  []string{"Test Author"},
			Language: "en",
		},
		Chapters: []Chapter{
			{
				Title: "Chapter 1",
				Sections: []Section{
					{
						Title:   "Section 1",
						Content: "This is test content",
					},
				},
			},
		},
		Format:   format.FormatTXT,
		Language: "en",
	}

	convertedBook, err := ConvertBook(originalBook, format.FormatEPUB)
	if err != nil {
		t.Fatalf("ConvertBook failed: %v", err)
	}

	// Check that the book structure is preserved
	if convertedBook.Metadata.Title != originalBook.Metadata.Title {
		t.Errorf("Title not preserved: got %s, want %s", convertedBook.Metadata.Title, originalBook.Metadata.Title)
	}

	if len(convertedBook.Chapters) != len(originalBook.Chapters) {
		t.Errorf("Chapter count not preserved: got %d, want %d", len(convertedBook.Chapters), len(originalBook.Chapters))
	}

	// Check that format was changed
	if convertedBook.Format != format.FormatEPUB {
		t.Errorf("Format not changed: got %s, want %s", convertedBook.Format, format.FormatEPUB)
	}

	// Check that it's a different object (not the same reference)
	if convertedBook == originalBook {
		t.Error("ConvertBook returned the same object reference")
	}
}

func TestBook_ExtractText(t *testing.T) {
	book := &Book{
		Metadata: Metadata{
			Title: "Test Book",
		},
		Chapters: []Chapter{
			{
				Title: "Chapter 1",
				Sections: []Section{
					{
						Title:   "Section 1.1",
						Content: "This is the first section content.",
						Subsections: []Section{
							{
								Title:   "Subsection 1.1.1",
								Content: "This is a subsection.",
							},
						},
					},
					{
						Title:   "Section 1.2",
						Content: "This is the second section content.",
					},
				},
			},
			{
				Title: "Chapter 2",
				Sections: []Section{
					{
						Title:   "",
						Content: "This section has no title.",
					},
				},
			},
		},
	}

	text := book.ExtractText()

	// Check that title is included
	if !strings.Contains(text, "Test Book") {
		t.Error("Book title not found in extracted text")
	}

	// Check that chapter titles are included
	if !strings.Contains(text, "Chapter 1") {
		t.Error("Chapter 1 title not found in extracted text")
	}

	if !strings.Contains(text, "Chapter 2") {
		t.Error("Chapter 2 title not found in extracted text")
	}

	// Check that section titles are included
	if !strings.Contains(text, "Section 1.1") {
		t.Error("Section 1.1 title not found in extracted text")
	}

	if !strings.Contains(text, "Subsection 1.1.1") {
		t.Error("Subsection title not found in extracted text")
	}

	// Check that content is included
	if !strings.Contains(text, "This is the first section content.") {
		t.Error("Section content not found in extracted text")
	}

	if !strings.Contains(text, "This is a subsection.") {
		t.Error("Subsection content not found in extracted text")
	}
}

func TestBook_GetChapterCount(t *testing.T) {
	tests := []struct {
		name     string
		chapters []Chapter
		expected int
	}{
		{
			name:     "Empty book",
			chapters: []Chapter{},
			expected: 0,
		},
		{
			name: "Single chapter",
			chapters: []Chapter{
				{Title: "Chapter 1"},
			},
			expected: 1,
		},
		{
			name: "Multiple chapters",
			chapters: []Chapter{
				{Title: "Chapter 1"},
				{Title: "Chapter 2"},
				{Title: "Chapter 3"},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := &Book{Chapters: tt.chapters}
			count := book.GetChapterCount()
			if count != tt.expected {
				t.Errorf("GetChapterCount() = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestBook_GetWordCount(t *testing.T) {
	tests := []struct {
		name     string
		book     *Book
		expected int
	}{
		{
			name: "Empty book",
			book: &Book{
				Metadata: Metadata{Title: ""},
				Chapters: []Chapter{},
			},
			expected: 0,
		},
		{
			name: "Simple book",
			book: &Book{
				Metadata: Metadata{Title: "Test"},
				Chapters: []Chapter{
					{
						Title: "Chapter",
						Sections: []Section{
							{
								Content: "This is a test sentence.",
							},
						},
					},
				},
			},
			expected: 7, // Test + Chapter + This + is + a + test + sentence
		},
		{
			name: "Book with multiple spaces",
			book: &Book{
				Metadata: Metadata{Title: "Test  Book"},
				Chapters: []Chapter{
					{
						Title: "Chapter   One",
						Sections: []Section{
							{
								Content: "This    is    a    test.",
							},
						},
					},
				},
			},
			expected: 8, // Test + Book + Chapter + One + This + is + a + test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := tt.book.GetWordCount()
			if count != tt.expected {
				t.Errorf("GetWordCount() = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestExtractSectionText(t *testing.T) {
	section := &Section{
		Title:   "Main Section",
		Content: "This is the main content.",
		Subsections: []Section{
			{
				Title:   "Subsection 1",
				Content: "This is subsection content.",
				Subsections: []Section{
					{
						Title:   "",
						Content: "Nested subsection without title.",
					},
				},
			},
			{
				Title:   "Subsection 2",
				Content: "More subsection content.",
			},
		},
	}

	text := extractSectionText(section)

	// Check main section title and content
	if !strings.Contains(text, "Main Section") {
		t.Error("Main section title not found")
	}

	if !strings.Contains(text, "This is the main content.") {
		t.Error("Main section content not found")
	}

	// Check subsections
	if !strings.Contains(text, "Subsection 1") {
		t.Error("Subsection 1 title not found")
	}

	if !strings.Contains(text, "This is subsection content.") {
		t.Error("Subsection 1 content not found")
	}

	if !strings.Contains(text, "Nested subsection without title.") {
		t.Error("Nested subsection content not found")
	}

	if !strings.Contains(text, "Subsection 2") {
		t.Error("Subsection 2 title not found")
	}
}

func TestBook_Structure(t *testing.T) {
	book := &Book{
		Metadata: Metadata{
			Title:       "Test Book",
			Authors:     []string{"Author 1", "Author 2"},
			Description: "A test book",
			Publisher:   "Test Publisher",
			Language:    "en",
			ISBN:        "1234567890",
			Date:        "2023-01-01",
		},
		Chapters: []Chapter{
			{
				Title: "Chapter 1",
				Sections: []Section{
					{
						Title:   "Introduction",
						Content: "This is an introduction",
					},
				},
			},
		},
		Format:   format.FormatEPUB,
		Language: "en",
	}

	// Test metadata access
	if book.Metadata.Title != "Test Book" {
		t.Errorf("Title = %s, want Test Book", book.Metadata.Title)
	}

	if len(book.Metadata.Authors) != 2 {
		t.Errorf("Authors count = %d, want 2", len(book.Metadata.Authors))
	}

	if book.Metadata.Authors[0] != "Author 1" {
		t.Errorf("First author = %s, want Author 1", book.Metadata.Authors[0])
	}

	// Test chapter access
	if len(book.Chapters) != 1 {
		t.Errorf("Chapters count = %d, want 1", len(book.Chapters))
	}

	if book.Chapters[0].Title != "Chapter 1" {
		t.Errorf("Chapter title = %s, want Chapter 1", book.Chapters[0].Title)
	}

	// Test section access
	if len(book.Chapters[0].Sections) != 1 {
		t.Errorf("Sections count = %d, want 1", len(book.Chapters[0].Sections))
	}

	if book.Chapters[0].Sections[0].Title != "Introduction" {
		t.Errorf("Section title = %s, want Introduction", book.Chapters[0].Sections[0].Title)
	}
}

// Test with temporary files
func createTempFile(t *testing.T, content, extension string) string {
	tmpFile, err := os.CreateTemp("", "test*"+extension)
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func TestUniversalParser_Parse_NonExistentFile(t *testing.T) {
	parser := NewUniversalParser()

	_, err := parser.Parse("/non/existent/file.fb2")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestUniversalParser_Parse_EmptyFile(t *testing.T) {
	parser := NewUniversalParser()

	tmpFile := createTempFile(t, "", ".txt")
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// TXT parser sets title to filename, so we expect the filename
	expectedTitle := strings.TrimSuffix(tmpFile, ".txt")
	if !strings.Contains(book.Metadata.Title, expectedTitle) {
		t.Errorf("Expected title to contain filename %s, got %s", expectedTitle, book.Metadata.Title)
	}

	// TXT parser always creates one chapter called "Content"
	if len(book.Chapters) != 1 {
		t.Errorf("Expected 1 chapter for TXT file, got %d", len(book.Chapters))
	}

	if len(book.Chapters) > 0 && book.Chapters[0].Title != "Content" {
		t.Errorf("Expected chapter title 'Content', got %s", book.Chapters[0].Title)
	}
}
