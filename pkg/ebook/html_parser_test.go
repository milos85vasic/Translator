package ebook

import (
	"os"
	"strings"
	"testing"

	"digital.vasic.translator/pkg/format"
	"golang.org/x/net/html"
)

func TestNewHTMLParser(t *testing.T) {
	parser := NewHTMLParser()
	if parser == nil {
		t.Fatal("NewHTMLParser returned nil")
	}
}

func TestHTMLParser_GetFormat(t *testing.T) {
	parser := NewHTMLParser()
	if parser.GetFormat() != format.FormatHTML {
		t.Errorf("GetFormat() = %s, want %s", parser.GetFormat(), format.FormatHTML)
	}
}

func TestHTMLParser_Parse_NonExistentFile(t *testing.T) {
	parser := NewHTMLParser()
	_, err := parser.Parse("/non/existent/file.html")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestHTMLParser_Parse_EmptyFile(t *testing.T) {
	parser := NewHTMLParser()

	tmpFile := createTempHTMLFile(t, "test_empty.html", "")
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Check metadata - should use filename as title when no title found
	expectedTitle := strings.TrimSuffix(tmpFile, ".html")
	if !strings.Contains(book.Metadata.Title, expectedTitle) {
		t.Errorf("Title = %s, expected to contain %s", book.Metadata.Title, expectedTitle)
	}

	// Check chapters
	if len(book.Chapters) != 1 {
		t.Errorf("Chapters count = %d, want 1", len(book.Chapters))
	}

	// Check sections
	if len(book.Chapters[0].Sections) != 1 {
		t.Errorf("Sections count = %d, want 1", len(book.Chapters[0].Sections))
	}

	content := book.Chapters[0].Sections[0].Content
	if content != "" {
		t.Errorf("Content = %q, want empty string", content)
	}
}

func TestHTMLParser_Parse_SimpleHTML(t *testing.T) {
	parser := NewHTMLParser()

	htmlContent := `<!DOCTYPE html>
<html>
<head>
	<title>Test Document</title>
</head>
<body>
	<p>This is a simple HTML document.</p>
	<p>It has multiple paragraphs.</p>
</body>
</html>`

	tmpFile := createTempHTMLFile(t, "test_simple.html", htmlContent)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Check title
	if book.Metadata.Title != "Test Document" {
		t.Errorf("Title = %s, want Test Document", book.Metadata.Title)
	}

	// Check content
	content := book.Chapters[0].Sections[0].Content

	if !strings.Contains(content, "This is a simple HTML document.") {
		t.Error("First paragraph not found in content")
	}

	if !strings.Contains(content, "It has multiple paragraphs.") {
		t.Error("Second paragraph not found in content")
	}
}

func TestHTMLParser_Parse_ComplexHTML(t *testing.T) {
	parser := NewHTMLParser()

	htmlContent := `<!DOCTYPE html>
<html>
<head>
	<title>Complex Document</title>
	<style>
		body { font-family: Arial; }
	</style>
	<script>
		console.log("test");
	</script>
</head>
<body>
	<header>
		<h1>Main Title</h1>
	</header>
	<main>
		<section>
			<h2>Section 1</h2>
			<p>First section content.</p>
			<ul>
				<li>List item 1</li>
				<li>List item 2</li>
			</ul>
		</section>
		<section>
			<h2>Section 2</h2>
			<p>Second section content.</p>
			<blockquote>
				This is a quote.
			</blockquote>
		</section>
	</main>
	<footer>
		<p>Footer content</p>
	</footer>
</body>
</html>`

	tmpFile := createTempHTMLFile(t, "test_complex.html", htmlContent)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Check title
	if book.Metadata.Title != "Complex Document" {
		t.Errorf("Title = %s, want Complex Document", book.Metadata.Title)
	}

	// Check content
	content := book.Chapters[0].Sections[0].Content

	// Should contain main content but not script/style
	if !strings.Contains(content, "Main Title") {
		t.Error("Main title not found in content")
	}

	if !strings.Contains(content, "Section 1") {
		t.Error("Section 1 not found in content")
	}

	if !strings.Contains(content, "First section content.") {
		t.Error("First section content not found")
	}

	if !strings.Contains(content, "List item 1") {
		t.Error("List item 1 not found")
	}

	if !strings.Contains(content, "This is a quote.") {
		t.Error("Quote not found")
	}

	// Should not contain script or style content
	if strings.Contains(content, "console.log") {
		t.Error("Script content found in extracted text")
	}

	if strings.Contains(content, "font-family") {
		t.Error("Style content found in extracted text")
	}
}

func TestHTMLParser_Parse_NoTitle(t *testing.T) {
	parser := NewHTMLParser()

	htmlContent := `<!DOCTYPE html>
<html>
<body>
	<p>Content without title tag.</p>
</body>
</html>`

	tmpFile := createTempHTMLFile(t, "test_no_title.html", htmlContent)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Should use filename as title when no title tag found
	expectedTitle := strings.TrimSuffix(tmpFile, ".html")
	if !strings.Contains(book.Metadata.Title, expectedTitle) {
		t.Errorf("Title = %s, expected to contain %s", book.Metadata.Title, expectedTitle)
	}

	// Content should still be extracted
	content := book.Chapters[0].Sections[0].Content
	if !strings.Contains(content, "Content without title tag.") {
		t.Error("Content not found when no title tag present")
	}
}

func TestHTMLParser_findTitle(t *testing.T) {
	parser := NewHTMLParser()

	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Title in head",
			html:     "<html><head><title>Test Title</title></head><body></body></html>",
			expected: "Test Title",
		},
		{
			name:     "No title tag",
			html:     "<html><head></head><body></body></html>",
			expected: "",
		},
		{
			name:     "Empty title tag",
			html:     "<html><head><title></title></head><body></body></html>",
			expected: "",
		},
		{
			name:     "Title with whitespace",
			html:     "<html><head><title>  Title with spaces  </title></head><body></body></html>",
			expected: "  Title with spaces  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse HTML to get node
			doc, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			title := parser.findTitle(doc)
			if title != tt.expected {
				t.Errorf("findTitle() = %q, want %q", title, tt.expected)
			}
		})
	}
}

func TestHTMLParser_extractText(t *testing.T) {
	parser := NewHTMLParser()

	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple text",
			html:     "<p>Simple paragraph</p>",
			expected: "Simple paragraph",
		},
		{
			name:     "Nested elements",
			html:     "<div><p>Nested <strong>text</strong> here</p></div>",
			expected: "Nested text here",
		},
		{
			name:     "Multiple paragraphs",
			html:     "<p>First paragraph</p><p>Second paragraph</p>",
			expected: "First paragraph\n\nSecond paragraph",
		},
		{
			name:     "Script and style should be ignored",
			html:     "<p>Visible text</p><script>var x = 1;</script><style>body { color: red; }</style><p>More text</p>",
			expected: "Visible text\n\nMore text",
		},
		{
			name:     "Headings",
			html:     "<h1>Main Title</h1><p>Content</p><h2>Subtitle</h2><p>More content</p>",
			expected: "Main Title\n\nContent\n\nSubtitle\n\nMore content",
		},
		{
			name:     "List items",
			html:     "<ul><li>Item 1</li><li>Item 2</li></ul>",
			expected: "Item 1\n\nItem 2",
		},
		{
			name:     "Blockquote",
			html:     "<p>Regular text</p><blockquote>Quoted text</blockquote><p>More text</p>",
			expected: "Regular text\n\nQuoted text\n\nMore text",
		},
		{
			name:     "Preformatted text",
			html:     "<pre>Preformatted\n  text\n    here</pre>",
			expected: "Preformatted\n  text\n    here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse HTML to get node
			doc, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			text := parser.extractText(doc)
			// Trim trailing whitespace for comparison
			text = strings.TrimSpace(text)
			expected := strings.TrimSpace(tt.expected)

			if text != expected {
				t.Errorf("extractText() = %q, want %q", text, expected)
			}
		})
	}
}

func TestIsBlockElement(t *testing.T) {
	tests := []struct {
		tag      string
		expected bool
	}{
		{"p", true},
		{"div", true},
		{"h1", true},
		{"h2", true},
		{"h3", true},
		{"h4", true},
		{"h5", true},
		{"h6", true},
		{"li", true},
		{"section", true},
		{"article", true},
		{"header", true},
		{"footer", true},
		{"blockquote", true},
		{"pre", true},
		{"span", false},
		{"strong", false},
		{"em", false},
		{"a", false},
		{"img", false},
		{"br", false},
		{"hr", false},
		{"script", false},
		{"style", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			result := isBlockElement(tt.tag)
			if result != tt.expected {
				t.Errorf("isBlockElement(%s) = %v, want %v", tt.tag, result, tt.expected)
			}
		})
	}
}

func TestHTMLParser_Parse_InvalidHTML(t *testing.T) {
	parser := NewHTMLParser()

	// Create a file with invalid HTML
	invalidHTML := "<html><body><p>Unclosed paragraph<div>Unclosed div"
	tmpFile := createTempHTMLFile(t, "test_invalid.html", invalidHTML)
	defer os.Remove(tmpFile)

	book, err := parser.Parse(tmpFile)
	// HTML parser is forgiving, so it should still parse
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if book == nil {
		t.Fatal("Book is nil")
	}

	// Should still extract some content
	content := book.Chapters[0].Sections[0].Content
	if !strings.Contains(content, "Unclosed paragraph") {
		t.Error("Content not extracted from invalid HTML")
	}
}

// Helper function to create temporary HTML files
func createTempHTMLFile(t *testing.T, filename, content string) string {
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
