package markdown

import (
	"digital.vasic.translator/pkg/ebook"
	"os"
	"strings"
	"testing"
)

// Test Markdown Translation with formatting preservation
func TestMarkdownTranslation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple paragraph",
			input:    "This is a simple paragraph.",
			expected: "TRANSLATED: This is a simple paragraph.",
		},
		{
			name:     "Bold text",
			input:    "This is **bold** text.",
			expected: "TRANSLATED: This is**TRANSLATED: bold**TRANSLATED: text.",
		},
		{
			name:     "Italic text",
			input:    "This is *italic* text.",
			expected: "TRANSLATED: This is*TRANSLATED: italic*TRANSLATED: text.",
		},
		{
			name:     "Header",
			input:    "# Main Title",
			expected: "# TRANSLATED: Main Title",
		},
		{
			name:     "List item",
			input:    "- First item",
			expected: "- TRANSLATED: First item",
		},
		{
			name:     "Code block (not translated)",
			input:    "```\ncode here\n```",
			expected: "```\ncode here\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := NewMarkdownTranslator(func(text string) (string, error) {
				return "TRANSLATED: " + text, nil
			})

			result, err := translator.TranslateMarkdown(tt.input)
			if err != nil {
				t.Fatalf("Translation failed: %v", err)
			}

			result = strings.TrimSpace(result)
			tt.expected = strings.TrimSpace(tt.expected)

			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// Test markdown inline formatting preservation
func TestInlineFormatting(t *testing.T) {
	translator := NewMarkdownTranslator(func(text string) (string, error) {
		return "TR[" + text + "]", nil
	})

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "Mixed formatting",
			input:    "Text with **bold** and *italic* and `code`.",
			contains: []string{"**TR[bold]**", "*TR[italic]*", "`code`"},
		},
		{
			name:     "Link",
			input:    "Check [this link](http://example.com).",
			contains: []string{"[TR[this link]](http://example.com)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := translator.TranslateMarkdown(tt.input)
			if err != nil {
				t.Fatalf("Translation failed: %v", err)
			}

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Result should contain '%s'\nGot: %s", expected, result)
				}
			}
		})
	}
}

// Test frontmatter preservation
func TestFrontmatterPreservation(t *testing.T) {
	input := `---
title: Test Book
authors: Author Name
language: en
---

# Chapter 1

This is content.`

	translator := NewMarkdownTranslator(func(text string) (string, error) {
		return "TRANSLATED: " + text, nil
	})

	result, err := translator.TranslateMarkdown(input)
	if err != nil {
		t.Fatalf("Translation failed: %v", err)
	}

	// Frontmatter should be preserved exactly
	if !strings.Contains(result, "title: Test Book") {
		t.Error("Frontmatter 'title' was not preserved")
	}

	if !strings.Contains(result, "authors: Author Name") {
		t.Error("Frontmatter 'authors' was not preserved")
	}

	// Content should be translated
	if !strings.Contains(result, "TRANSLATED:") {
		t.Error("Content was not translated")
	}
}

// Test Book to Markdown conversion
func TestBookToMarkdown(t *testing.T) {
	book := &ebook.Book{
		Metadata: ebook.Metadata{
			Title:    "Test Book",
			Authors:  []string{"Author 1", "Author 2"},
			Language: "en",
		},
		Chapters: []ebook.Chapter{
			{
				Title: "Chapter 1",
				Sections: []ebook.Section{
					{Content: "First paragraph.\n\nSecond paragraph."},
				},
			},
			{
				Title: "Chapter 2",
				Sections: []ebook.Section{
					{Content: "Another paragraph."},
				},
			},
		},
	}

	// Create temp file
	tmpFile := "/tmp/test_book.md"
	defer os.Remove(tmpFile)

	err := ConvertBookToMarkdown(book, tmpFile)
	if err != nil {
		t.Fatalf("Failed to convert book to markdown: %v", err)
	}

	// Read the result
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	result := string(content)

	// Check frontmatter
	if !strings.Contains(result, "title: Test Book") {
		t.Error("Missing title in frontmatter")
	}

	if !strings.Contains(result, "authors: Author 1, Author 2") {
		t.Error("Missing authors in frontmatter")
	}

	// Check chapters
	if !strings.Contains(result, "## Chapter 1") {
		t.Error("Missing Chapter 1 header")
	}

	if !strings.Contains(result, "## Chapter 2") {
		t.Error("Missing Chapter 2 header")
	}

	// Check content
	if !strings.Contains(result, "First paragraph") {
		t.Error("Missing chapter content")
	}
}

// Test Markdown to EPUB metadata parsing
func TestMarkdownToEPUBMetadata(t *testing.T) {
	markdown := `---
title: My Test Book
authors: John Doe, Jane Smith
language: en
---

# My Test Book

**By John Doe, Jane Smith**

---

## Chapter 1

Content here.`

	converter := NewMarkdownToEPUBConverter()

	// Create temp files
	tmpMD := "/tmp/test_source.md"
	tmpEPUB := "/tmp/test_output.epub"
	defer os.Remove(tmpMD)
	defer os.Remove(tmpEPUB)

	// Write markdown
	if err := os.WriteFile(tmpMD, []byte(markdown), 0644); err != nil {
		t.Fatalf("Failed to write temp markdown: %v", err)
	}

	// Convert to EPUB
	if err := converter.ConvertMarkdownToEPUB(tmpMD, tmpEPUB); err != nil {
		t.Fatalf("Failed to convert markdown to EPUB: %v", err)
	}

	// Verify EPUB was created
	if _, err := os.Stat(tmpEPUB); os.IsNotExist(err) {
		t.Fatal("EPUB file was not created")
	}

	// Check file size is reasonable (should be more than mimetype)
	info, _ := os.Stat(tmpEPUB)
	if info.Size() < 100 {
		t.Errorf("EPUB file is too small: %d bytes", info.Size())
	}
}

// Test empty and edge cases
func TestEdgeCases(t *testing.T) {
	translator := NewMarkdownTranslator(func(text string) (string, error) {
		return "TRANSLATED", nil
	})

	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Only whitespace", "   \n  \n  "},
		{"Only frontmatter", "---\ntitle: Test\n---"},
		{"Only code block", "```\ncode\n```"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := translator.TranslateMarkdown(tt.input)
			if err != nil {
				t.Errorf("Should handle edge case without error: %v", err)
			}
		})
	}
}

// Test complex markdown structure
func TestComplexMarkdown(t *testing.T) {
	input := `---
title: Complex Document
language: en
---

# Main Title

## Chapter 1: Introduction

This is a paragraph with **bold text** and *italic text*.

### Section 1.1

Here's a list:

- Item 1 with **bold**
- Item 2 with *italic*
- Item 3 with [link](http://example.com)

### Section 1.2

> This is a blockquote.

` + "```" + `
This is code that should not be translated.
int main() { return 0; }
` + "```" + `

## Chapter 2: Conclusion

Final paragraph.`

	translator := NewMarkdownTranslator(func(text string) (string, error) {
		return "TR:" + text, nil
	})

	result, err := translator.TranslateMarkdown(input)
	if err != nil {
		t.Fatalf("Translation failed: %v", err)
	}

	// Verify structure is preserved
	checks := []struct {
		description string
		shouldHave  string
	}{
		{"Frontmatter", "title: Complex Document"},
		{"Main header", "# TR:Main Title"},
		{"Chapter header", "## TR:Chapter 1: Introduction"},
		{"Subsection", "### TR:Section 1.1"},
		{"Bold formatting", "**TR:bold**"},
		{"Italic formatting", "*TR:italic*"},
		{"List structure", "- TR:Item 1 with**TR:bold**"},
		{"Code block untranslated", "int main() { return 0; }"},
		{"Link structure", "[TR:link](http://example.com)"},
	}

	for _, check := range checks {
		if !strings.Contains(result, check.shouldHave) {
			t.Errorf("%s check failed.\nExpected to find: %s\nIn result:\n%s",
				check.description, check.shouldHave, result)
		}
	}
}

// Benchmark markdown translation
func BenchmarkMarkdownTranslation(b *testing.B) {
	input := `# Chapter Title

This is a paragraph with **bold** and *italic* text.

## Section

- Item 1
- Item 2
- Item 3`

	translator := NewMarkdownTranslator(func(text string) (string, error) {
		return "TRANSLATED: " + text, nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := translator.TranslateMarkdown(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark inline formatting preservation
func BenchmarkInlineFormatting(b *testing.B) {
	translator := NewMarkdownTranslator(func(text string) (string, error) {
		return "TR:" + text, nil
	})

	input := "Text with **bold**, *italic*, and `code` formatting."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		translator.translateInlineFormatting(input)
	}
}
