package ebook

import (
	"digital.vasic.translator/pkg/format"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// HTMLParser implements Parser for HTML format
type HTMLParser struct{}

// NewHTMLParser creates a new HTML parser
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// Parse parses an HTML file into universal Book structure
func (p *HTMLParser) Parse(filename string) (*Book, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	book := &Book{
		Metadata: Metadata{
			Title: filename,
		},
		Chapters: make([]Chapter, 0),
		Format:   format.FormatHTML,
	}

	// Extract title
	title := p.findTitle(doc)
	if title != "" {
		book.Metadata.Title = title
	}

	// Extract content
	content := p.extractText(doc)

	// Create single chapter
	chapter := Chapter{
		Title: book.Metadata.Title,
		Sections: []Section{
			{
				Content: content,
			},
		},
	}

	book.Chapters = append(book.Chapters, chapter)

	return book, nil
}

// findTitle finds the title in HTML
func (p *HTMLParser) findTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil {
			return n.FirstChild.Data
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := p.findTitle(c); title != "" {
			return title
		}
	}

	return ""
}

// extractText extracts text content from HTML
func (p *HTMLParser) extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			return text + " "
		}
	}

	var content strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// Skip script and style tags
		if c.Type == html.ElementNode && (c.Data == "script" || c.Data == "style") {
			continue
		}

		text := p.extractText(c)
		content.WriteString(text)

		// Add newlines after block elements
		if c.Type == html.ElementNode && isBlockElement(c.Data) {
			content.WriteString("\n\n")
		}
	}

	return content.String()
}

// isBlockElement checks if HTML element is a block element
func isBlockElement(tag string) bool {
	blockElements := []string{
		"p", "div", "h1", "h2", "h3", "h4", "h5", "h6",
		"li", "section", "article", "header", "footer",
		"blockquote", "pre",
	}

	for _, elem := range blockElements {
		if tag == elem {
			return true
		}
	}
	return false
}

// GetFormat returns the format
func (p *HTMLParser) GetFormat() format.Format {
	return format.FormatHTML
}
