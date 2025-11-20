package unit

import (
	"digital.vasic.translator/pkg/ebook"
	"testing"
)

func TestEbookStructure(t *testing.T) {
	t.Run("CreateBook", func(t *testing.T) {
		book := &ebook.Book{
			Metadata: ebook.Metadata{
				Title:   "Test Book",
				Authors: []string{"Test Author"},
				Language: "en",
			},
			Chapters: []ebook.Chapter{
				{
					Title: "Chapter 1",
					Sections: []ebook.Section{
						{
							Content: "This is test content.",
						},
					},
				},
			},
		}

		if book.Metadata.Title != "Test Book" {
			t.Errorf("Expected title 'Test Book', got '%s'", book.Metadata.Title)
		}

		if book.GetChapterCount() != 1 {
			t.Errorf("Expected 1 chapter, got %d", book.GetChapterCount())
		}
	})

	t.Run("ExtractText", func(t *testing.T) {
		book := &ebook.Book{
			Metadata: ebook.Metadata{
				Title: "Test Book",
			},
			Chapters: []ebook.Chapter{
				{
					Title: "Chapter 1",
					Sections: []ebook.Section{
						{
							Content: "First section content.",
						},
						{
							Content: "Second section content.",
						},
					},
				},
			},
		}

		text := book.ExtractText()

		if text == "" {
			t.Error("ExtractText returned empty string")
		}

		if !containsString(text, "Test Book") {
			t.Error("Extracted text should contain title")
		}

		if !containsString(text, "Chapter 1") {
			t.Error("Extracted text should contain chapter title")
		}

		if !containsString(text, "First section content") {
			t.Error("Extracted text should contain section content")
		}
	})

	t.Run("GetWordCount", func(t *testing.T) {
		book := &ebook.Book{
			Metadata: ebook.Metadata{
				Title: "Test",
			},
			Chapters: []ebook.Chapter{
				{
					Title: "Chapter",
					Sections: []ebook.Section{
						{
							Content: "One two three four five.",
						},
					},
				},
			},
		}

		wordCount := book.GetWordCount()

		if wordCount < 5 {
			t.Errorf("Expected word count >= 5, got %d", wordCount)
		}
	})
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
