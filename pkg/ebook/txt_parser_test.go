package ebook

import (
	"os"
	"strings"
	"testing"

	"digital.vasic.translator/pkg/format"
)

func TestNewTXTParser(t *testing.T) {
	parser := NewTXTParser()
	if parser == nil {
		t.Fatal("NewTXTParser returned nil")
	}
}

func TestTXTParser_GetFormat(t *testing.T) {
	parser := NewTXTParser()
	if parser.GetFormat() != format.FormatTXT {
		t.Errorf("GetFormat() = %s, want %s", parser.GetFormat(), format.FormatTXT)
	}
}

func TestTXTParser_Parse_NonExistentFile(t *testing.T) {
	parser := NewTXTParser()
	_, err := parser.Parse("/non/existent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestTXTParser_Parse_EmptyFile(t *testing.T) {
	parser := NewTXTParser()

	tmpFile := createTempTXTFile(t, "test_empty.txt", "")
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Check metadata
	expectedTitle := strings.TrimSuffix(tmpFile, ".txt")
	if !strings.Contains(book.Metadata.Title, expectedTitle) {
		t.Errorf("Title = %s, expected to contain %s", book.Metadata.Title, expectedTitle)
	}

	// Check chapters
	if len(book.Chapters) != 1 {
		t.Errorf("Chapters count = %d, want 1", len(book.Chapters))
	}

	if book.Chapters[0].Title != "Content" {
		t.Errorf("Chapter title = %s, want Content", book.Chapters[0].Title)
	}

	// Check sections
	if len(book.Chapters[0].Sections) != 1 {
		t.Errorf("Sections count = %d, want 1", len(book.Chapters[0].Sections))
	}

	content := book.Chapters[0].Sections[0].Content
	if content != "" {
		t.Errorf("Content = %s, want empty string", content)
	}
}

func TestTXTParser_Parse_SimpleContent(t *testing.T) {
	parser := NewTXTParser()

	content := "This is a simple text file.\nIt has multiple lines.\nAnd some content."
	tmpFile := createTempTXTFile(t, "test_simple.txt", content)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Check content
	sectionContent := book.Chapters[0].Sections[0].Content
	expectedContent := content + "\n" // Parser adds newline at the end
	if sectionContent != expectedContent {
		t.Errorf("Content = %q, want %q", sectionContent, expectedContent)
	}
}

func TestTXTParser_Parse_MultilineContent(t *testing.T) {
	parser := NewTXTParser()

	content := `Chapter 1: The Beginning
This is the first chapter.
It has multiple paragraphs.

Chapter 2: The Middle
This is the second chapter.
More content here.

Chapter 3: The End
This is the final chapter.`

	tmpFile := createTempTXTFile(t, "test_multiline.txt", content)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Check that all content is preserved
	sectionContent := book.Chapters[0].Sections[0].Content
	expectedContent := content + "\n" // Parser adds newline at the end
	if sectionContent != expectedContent {
		t.Errorf("Content mismatch.\nGot: %q\nWant: %q", sectionContent, expectedContent)
	}

	// Check that chapter markers are preserved
	if !strings.Contains(sectionContent, "Chapter 1: The Beginning") {
		t.Error("Chapter 1 marker not found in content")
	}

	if !strings.Contains(sectionContent, "Chapter 2: The Middle") {
		t.Error("Chapter 2 marker not found in content")
	}

	if !strings.Contains(sectionContent, "Chapter 3: The End") {
		t.Error("Chapter 3 marker not found in content")
	}
}

func TestTXTParser_Parse_SpecialCharacters(t *testing.T) {
	parser := NewTXTParser()

	content := "Text with special characters: àáâãäåæçèéêë\nUnicode: 你好世界\nEmojis: 📚📖✨"
	tmpFile := createTempTXTFile(t, "test_special.txt", content)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	sectionContent := book.Chapters[0].Sections[0].Content
	expectedContent := content + "\n"
	if sectionContent != expectedContent {
		t.Errorf("Content = %q, want %q", sectionContent, expectedContent)
	}

	// Check special characters are preserved
	if !strings.Contains(sectionContent, "àáâãäåæçèéêë") {
		t.Error("Accented characters not preserved")
	}

	if !strings.Contains(sectionContent, "你好世界") {
		t.Error("Unicode characters not preserved")
	}

	if !strings.Contains(sectionContent, "📚📖✨") {
		t.Error("Emoji characters not preserved")
	}
}

func TestTXTParser_Parse_LongLines(t *testing.T) {
	parser := NewTXTParser()

	// Create a very long line
	longLine := strings.Repeat("This is a very long line. ", 1000)
	content := longLine + "\nShort line.\n" + longLine

	tmpFile := createTempTXTFile(t, "test_long.txt", content)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	sectionContent := book.Chapters[0].Sections[0].Content
	expectedContent := content + "\n"
	if sectionContent != expectedContent {
		t.Errorf("Long line content not preserved correctly")
	}

	// Check that the long line is intact
	if !strings.Contains(sectionContent, longLine) {
		t.Error("Long line content not found or truncated")
	}
}

func TestTXTParser_Parse_WhitespaceHandling(t *testing.T) {
	parser := NewTXTParser()

	content := "  Line with leading spaces\n\tLine with tab\nLine with trailing spaces   \n\n\nMultiple empty lines"
	tmpFile := createTempTXTFile(t, "test_whitespace.txt", content)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	sectionContent := book.Chapters[0].Sections[0].Content
	expectedContent := content + "\n"
	if sectionContent != expectedContent {
		t.Errorf("Whitespace not preserved correctly.\nGot: %q\nWant: %q", sectionContent, expectedContent)
	}

	// Check that whitespace is preserved
	if !strings.Contains(sectionContent, "  Line with leading spaces") {
		t.Error("Leading spaces not preserved")
	}

	if !strings.Contains(sectionContent, "\tLine with tab") {
		t.Error("Tab characters not preserved")
	}

	if !strings.Contains(sectionContent, "trailing spaces   ") {
		t.Error("Trailing spaces not preserved")
	}
}

func TestTXTParser_Parse_FileStructure(t *testing.T) {
	parser := NewTXTParser()

	content := "Test content for file structure verification."
	tmpFile := createTempTXTFile(t, "test_structure.txt", content)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check book structure
	if book.Format != format.FormatTXT {
		t.Errorf("Book format = %s, want %s", book.Format, format.FormatTXT)
	}

	// Check metadata structure
	if book.Metadata.Title == "" {
		t.Error("Book title is empty")
	}

	// Check chapters structure
	if len(book.Chapters) != 1 {
		t.Fatalf("Chapters count = %d, want 1", len(book.Chapters))
	}

	chapter := book.Chapters[0]
	if chapter.Title != "Content" {
		t.Errorf("Chapter title = %s, want Content", chapter.Title)
	}

	// Check sections structure
	if len(chapter.Sections) != 1 {
		t.Fatalf("Sections count = %d, want 1", len(chapter.Sections))
	}

	section := chapter.Sections[0]
	if section.Title != "" {
		t.Errorf("Section title = %s, want empty string", section.Title)
	}

	if len(section.Subsections) != 0 {
		t.Errorf("Subsections count = %d, want 0", len(section.Subsections))
	}
}

// Helper function to create temporary text files
func createTempTXTFile(t *testing.T, filename, content string) string {
	tmpFile, err := os.CreateTemp("", filename)
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}
