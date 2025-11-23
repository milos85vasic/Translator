package ebook

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/schema/soo/wml"
)

func init() {
	// Set up uniDoc license (community edition)
	err := license.SetCommunityLicense()
	if err != nil {
		// Log error but continue - community edition will still work
		fmt.Printf("Warning: Could not set community license: %v\n", err)
	}
}

type DOCXParser struct {
	config *DOCXConfig
}

type DOCXConfig struct {
	ExtractImages     bool     `yaml:"extract_images"`
	ImageFormat       string   `yaml:"image_format"`
	ExtractTables     bool     `yaml:"extract_tables"`
	ExtractFootnotes  bool     `yaml:"extract_footnotes"`
	ExtractHeaders    bool     `yaml:"extract_headers"`
	ExtractFooters    bool     `yaml:"extract_footers"`
	ExtractComments   bool     `yaml:"extract_comments"`
	PreserveFormatting bool     `yaml:"preserve_formatting"`
	ExtractMetadata   bool     `yaml:"extract_metadata"`
	MinTextLength     int      `yaml:"min_text_length"`
	IgnoreStyles      []string `yaml:"ignore_styles"`
}

type DOCXContent struct {
	Text       string        `json:"text"`
	Images     []DOCXImage  `json:"images,omitempty"`
	Tables     []DOCXTable  `json:"tables,omitempty"`
	Metadata   DOCXMetadata `json:"metadata,omitempty"`
	Styles     DOCXStyles   `json:"styles,omitempty"`
	Comments   []DOCXComment `json:"comments,omitempty"`
	Footnotes  []DOCXFootnote `json:"footnotes,omitempty"`
	Headers    []DOCXHeader  `json:"headers,omitempty"`
	Footers    []DOCXFooter  `json:"footers,omitempty"`
}

type DOCXImage struct {
	Data     []byte `json:"data"`
	Format   string `json:"format"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	ID       string `json:"id"`
	AltText  string `json:"alt_text"`
}

type DOCXTable struct {
	Rows    [][]string `json:"rows"`
	Columns int        `json:"columns"`
	ID      string     `json:"id"`
	Style   string     `json:"style"`
}

type DOCXMetadata struct {
	Title        string    `json:"title"`
	Subject      string    `json:"subject"`
	Creator      string    `json:"creator"`
	Keywords     []string  `json:"keywords"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Status       string    `json:"status"`
	Language     string    `json:"language"`
	CreationDate time.Time `json:"creation_date"`
	ModDate      time.Time `json:"mod_date"`
	LastModified string    `json:"last_modified_by"`
	Revision     int       `json:"revision"`
	Pages        int       `json:"pages"`
	Words        int       `json:"words"`
	Characters   int       `json:"characters"`
}

type DOCXStyles struct {
	Paragraphs []DOCXStyle `json:"paragraphs"`
	Characters  []DOCXStyle `json:"characters"`
	Tables      []DOCXStyle `json:"tables"`
}

type DOCXStyle struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
	Bold   bool   `json:"bold"`
	Italic bool   `json:"italic"`
	Underline bool `json:"underline"`
	Size   int    `json:"size"`
	Color  string `json:"color"`
	Font   string `json:"font"`
}

type DOCXComment struct {
	ID       string `json:"id"`
	Author   string `json:"author"`
	Date     string `json:"date"`
	Text     string `json:"text"`
	Position struct {
		Paragraph int `json:"paragraph"`
		Character  int `json:"character"`
	} `json:"position"`
}

type DOCXFootnote struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type DOCXHeader struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Type  string `json:"type"` // primary, even, odd
}

type DOCXFooter struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Type  string `json:"type"` // primary, even, odd
}

func NewDOCXParser(config *DOCXConfig) *DOCXParser {
	if config == nil {
		config = &DOCXConfig{
			ExtractImages:      true,
			ImageFormat:        "png",
			ExtractTables:      true,
			ExtractFootnotes:   true,
			ExtractHeaders:      true,
			ExtractFooters:      true,
			ExtractComments:     true,
			PreserveFormatting:  true,
			ExtractMetadata:     true,
			MinTextLength:      5,
			IgnoreStyles:       []string{},
		}
	}
	return &DOCXParser{config: config}
}

func (p *DOCXParser) Parse(ctx context.Context, data []byte) (*Ebook, error) {
	doc, err := document.Read(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX document: %w", err)
	}

	ebook := &Ebook{
		Format: "docx",
		Metadata: Metadata{
			Format: "docx",
		},
	}

	// Extract metadata
	if p.config.ExtractMetadata {
		if err := p.extractMetadata(doc, ebook); err != nil {
			return nil, fmt.Errorf("failed to extract metadata: %w", err)
		}
	}

	// Extract main content
	content := &DOCXContent{}
	var allText strings.Builder

	// Extract paragraphs
	paragraphs := doc.Paragraphs()
	for i, para := range paragraphs {
		// Check for ignored styles
		styleName := para.Style()
		shouldIgnore := false
		for _, ignoredStyle := range p.config.IgnoreStyles {
			if styleName == ignoredStyle {
				shouldIgnore = true
				break
			}
		}

		if shouldIgnore {
			continue
		}

		paraText := p.extractParagraphText(para)
		if len(paraText) >= p.config.MinTextLength {
			if allText.Len() > 0 {
				allText.WriteString("\n\n") // Paragraph separator
			}
			allText.WriteString(paraText)
		}

		// Extract tables from paragraphs
		if p.config.ExtractTables {
			tables := p.extractTablesFromParagraph(para)
			content.Tables = append(content.Tables, tables...)
		}

		// Check for context cancellation
		if i%10 == 0 { // Check every 10 paragraphs
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
	}

	// Extract images if enabled
	if p.config.ExtractImages {
		images := p.extractImages(doc)
		content.Images = images
	}

	// Extract headers if enabled
	if p.config.ExtractHeaders {
		headers := p.extractHeaders(doc)
		content.Headers = headers
	}

	// Extract footers if enabled
	if p.config.ExtractFooters {
		footers := p.extractFooters(doc)
		content.Footers = footers
	}

	// Extract footnotes if enabled
	if p.config.ExtractFootnotes {
		footnotes := p.extractFootnotes(doc)
		content.Footnotes = footnotes
	}

	// Extract comments if enabled
	if p.config.ExtractComments {
		comments := p.extractComments(doc)
		content.Comments = comments
	}

	// Extract styles if preserving formatting
	if p.config.PreserveFormatting {
		styles := p.extractStyles(doc)
		content.Styles = styles
	}

	// Set main content
	ebook.Content = allText.String()

	// Add extracted resources to ebook
	if len(content.Images) > 0 || len(content.Tables) > 0 || 
	   len(content.Comments) > 0 || len(content.Footnotes) > 0 {
		ebook.Resources = make(map[string]interface{})
		
		if len(content.Images) > 0 {
			ebook.Resources["images"] = content.Images
		}
		if len(content.Tables) > 0 {
			ebook.Resources["tables"] = content.Tables
		}
		if len(content.Comments) > 0 {
			ebook.Resources["comments"] = content.Comments
		}
		if len(content.Footnotes) > 0 {
			ebook.Resources["footnotes"] = content.Footnotes
		}
	}

	return ebook, nil
}

func (p *DOCXParser) Validate(data []byte) error {
	// Check for DOCX file signature (should be a ZIP file)
	if len(data) < 4 {
		return fmt.Errorf("file too small to be a valid DOCX")
	}

	// DOCX files are ZIP archives, check for ZIP signature
	if !bytes.Equal(data[:4], []byte{0x50, 0x4B, 0x03, 0x04}) &&
	   !bytes.Equal(data[:4], []byte{0x50, 0x4B, 0x05, 0x06}) &&
	   !bytes.Equal(data[:4], []byte{0x50, 0x4B, 0x07, 0x08}) {
		return fmt.Errorf("invalid DOCX signature (not a ZIP archive)")
	}

	// Try to parse document structure
	_, err := document.Read(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("invalid DOCX structure: %w", err)
	}

	return nil
}

func (p *DOCXParser) SupportedFormats() []string {
	return []string{"docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"}
}

func (p *DOCXParser) GetMetadata(data []byte) (*Metadata, error) {
	doc, err := document.Read(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX document: %w", err)
	}

	metadata := &Metadata{Format: "docx"}
	err = p.extractMetadata(doc, &Ebook{Metadata: *metadata})
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}

	return metadata, nil
}

func (p *DOCXParser) extractMetadata(doc *document.Document, ebook *Ebook) error {
	// Extract core properties
	props := doc.CoreProperties()

	if props.Title != nil && *props.Title != "" {
		ebook.Metadata.Title = *props.Title
	}

	if props.Subject != nil && *props.Subject != "" {
		ebook.Metadata.Subject = *props.Subject
	}

	if props.Creator != nil && *props.Creator != "" {
		ebook.Metadata.Author = []string{*props.Creator}
	}

	if props.Description != nil && *props.Description != "" {
		ebook.Metadata.Description = *props.Description
	}

	if props.Keywords != nil && *props.Keywords != "" {
		keywords := strings.Split(*props.Keywords, ",")
		for i, kw := range keywords {
			keywords[i] = strings.TrimSpace(kw)
		}
		ebook.Metadata.Keywords = keywords
	}

	if props.Language != nil && *props.Language != "" {
		ebook.Metadata.Language = *props.Language
	}

	if props.Created != nil {
		ebook.Metadata.Date = props.Created.Format(time.RFC3339)
	}

	if props.Modified != nil {
		ebook.Metadata.Modified = props.Modified.Format(time.RFC3339)
	}

	// Extract extended properties
	extProps := doc.ExtendedProperties()

	if extProps.Company != nil && *extProps.Company != "" {
		if ebook.Metadata.Publisher == "" {
			ebook.Metadata.Publisher = *extProps.Company
		}
	}

	if extProps.Pages != nil {
		ebook.Metadata.Pages = extProps.Pages
	}

	if extProps.Words != nil {
		ebook.Metadata.Words = extProps.Words
	}

	if extProps.Characters != nil {
		ebook.Metadata.Characters = extProps.Characters
	}

	if extProps.Application != nil && *extProps.Application != "" {
		// Could be useful for format info
		_ = *extProps.Application
	}

	return nil
}

func (p *DOCXParser) extractParagraphText(para *document.Paragraph) string {
	var text strings.Builder

	// Extract runs from paragraph
	runs := para.Runs()
	for _, run := range runs {
		switch r := run.(type) {
		case *document.Run:
			text.WriteString(r.Text())
		case *document.Break:
			// Handle line breaks
			if r.Type() == wml.ST_BrTypeTextWrapping {
				text.WriteString(" ")
			} else {
				text.WriteString("\n")
			}
		}
	}

	return p.cleanText(text.String())
}

func (p *DOCXParser) extractTablesFromParagraph(para *document.Paragraph) []DOCXTable {
	var tables []DOCXTable

	// Extract tables in paragraph
	tablesInPara := para.Tables()
	for i, tbl := range tablesInPara {
		docxTable := DOCXTable{
			ID:    fmt.Sprintf("table_%d", i),
			Style:  tbl.Properties().Style(),
			Rows:  [][]string{},
		}

		// Extract table data
		for _, row := range tbl.Rows() {
			var rowData []string
			for _, cell := range row.Cells() {
				cellText := p.extractCellText(cell)
				rowData = append(rowData, cellText)
			}
			if len(rowData) > 0 {
				docxTable.Rows = append(docxTable.Rows, rowData)
			}
		}

		if len(docxTable.Rows) > 0 {
			docxTable.Columns = len(docxTable.Rows[0])
			tables = append(tables, docxTable)
		}
	}

	return tables
}

func (p *DOCXParser) extractCellText(cell *document.Cell) string {
	var text strings.Builder
	
	for _, para := range cell.Paragraphs() {
		paraText := p.extractParagraphText(para)
		text.WriteString(paraText)
		if para != cell.Paragraphs()[len(cell.Paragraphs())-1] {
			text.WriteString(" ")
		}
	}

	return strings.TrimSpace(text.String())
}

func (p *DOCXParser) extractImages(doc *document.Document) []DOCXImage {
	var images []DOCXImage

	// Extract images from document
	docImages := doc.Images()
	for i, img := range docImages {
		// Get image data
		imgData, err := img.Data()
		if err != nil {
			continue // Skip images that can't be extracted
		}

		// Determine format
		format := p.getImageFormat(imgData)

		images = append(images, DOCXImage{
			Data:    imgData,
			Format:  format,
			Width:   0, // Would need to decode to get actual dimensions
			Height:  0,
			ID:      fmt.Sprintf("img_%d", i),
			AltText: "", // Would need to extract from document XML
		})
	}

	return images
}

func (p *DOCXParser) extractHeaders(doc *document.Document) []DOCXHeader {
	var headers []DOCXHeader

	// Extract headers from different sections
	headerSections := doc.Headers()
	for i, hdr := range headerSections {
		docxHeader := DOCXHeader{
			ID:   fmt.Sprintf("header_%d", i),
			Text:  p.extractHeaderText(hdr),
			Type:  "primary", // Would need to determine actual type
		}
		headers = append(headers, docxHeader)
	}

	return headers
}

func (p *DOCXParser) extractFooters(doc *document.Document) []DOCXFooter {
	var footers []DOCXFooter

	// Extract footers from different sections
	footerSections := doc.Footers()
	for i, ftr := range footerSections {
		docxFooter := DOCXFooter{
			ID:   fmt.Sprintf("footer_%d", i),
			Text:  p.extractFooterText(ftr),
			Type:  "primary", // Would need to determine actual type
		}
		footers = append(footers, docxFooter)
	}

	return footers
}

func (p *DOCXParser) extractFootnotes(doc *document.Document) []DOCXFootnote {
	var footnotes []DOCXFootnote

	// Extract footnotes
	docFootnotes := doc.Footnotes()
	for i, fn := range docFootnotes {
		docxFootnote := DOCXFootnote{
			ID:   fmt.Sprintf("footnote_%d", i),
			Text:  p.extractFootnoteText(fn),
		}
		footnotes = append(footnotes, docxFootnote)
	}

	return footnotes
}

func (p *DOCXParser) extractComments(doc *document.Document) []DOCXComment {
	var comments []DOCXComment

	// Extract comments
	docComments := doc.Comments()
	for i, cmnt := range docComments {
		docxComment := DOCXComment{
			ID:     fmt.Sprintf("comment_%d", i),
			Author:  "", // Would need to extract from comment data
			Date:    "", // Would need to extract from comment data
			Text:    cmnt.Text(),
		}
		comments = append(comments, docxComment)
	}

	return comments
}

func (p *DOCXParser) extractStyles(doc *document.Document) DOCXStyles {
	styles := DOCXStyles{
		Paragraphs: []DOCXStyle{},
		Characters:  []DOCXStyle{},
		Tables:      []DOCXStyle{},
	}

	// Extract styles from document
	docStyles := doc.Styles()
	for _, style := range docStyles {
		docxStyle := DOCXStyle{
			Name:   style.Name(),
			Type:   p.getStyleType(style),
			Format: "", // Would need to extract from style properties
		}

		// This is a simplified style extraction
		// In production, you would extract all style properties
		switch style.Name() {
		case "Heading1":
			docxStyle.Size = 16
			docxStyle.Bold = true
		case "Heading2":
			docxStyle.Size = 14
			docxStyle.Bold = true
		case "Normal":
			docxStyle.Size = 11
		}

		styles.Paragraphs = append(styles.Paragraphs, docxStyle)
	}

	return styles
}

func (p *DOCXParser) extractHeaderText(hdr document.Header) string {
	var text strings.Builder
	
	for _, para := range hdr.Paragraphs() {
		paraText := p.extractParagraphText(para)
		text.WriteString(paraText)
		if para != hdr.Paragraphs()[len(hdr.Paragraphs())-1] {
			text.WriteString("\n")
		}
	}

	return strings.TrimSpace(text.String())
}

func (p *DOCXParser) extractFooterText(ftr document.Footer) string {
	var text strings.Builder
	
	for _, para := range ftr.Paragraphs() {
		paraText := p.extractParagraphText(para)
		text.WriteString(paraText)
		if para != ftr.Paragraphs()[len(ftr.Paragraphs())-1] {
			text.WriteString("\n")
		}
	}

	return strings.TrimSpace(text.String())
}

func (p *DOCXParser) extractFootnoteText(fn document.Footnote) string {
	var text strings.Builder
	
	for _, para := range fn.Paragraphs() {
		paraText := p.extractParagraphText(para)
		text.WriteString(paraText)
		if para != fn.Paragraphs()[len(fn.Paragraphs())-1] {
			text.WriteString("\n")
		}
	}

	return strings.TrimSpace(text.String())
}

func (p *DOCXParser) getStyleType(style document.Style) string {
	// Determine style type based on style properties
	// This is a simplified implementation
	switch style.Type() {
	case wml.ST_StyleTypeParagraph:
		return "paragraph"
	case wml.ST_StyleTypeCharacter:
		return "character"
	case wml.ST_StyleTypeTable:
		return "table"
	default:
		return "unknown"
	}
}

func (p *DOCXParser) getImageFormat(data []byte) string {
	// Simple image format detection based on file signatures
	if len(data) < 8 {
		return "unknown"
	}

	switch {
	case bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}):
		return "png"
	case bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}):
		return "jpeg"
	case bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}):
		return "gif"
	case bytes.HasPrefix(data, []byte{0x42, 0x4D}):
		return "bmp"
	default:
		return "unknown"
	}
}

func (p *DOCXParser) cleanText(text string) string {
	// Remove excessive whitespace
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\n ", "\n")
	text = strings.ReplaceAll(text, "  ", " ")

	// Remove common DOCX artifacts
	text = strings.ReplaceAll(text, "", "\"")
	text = strings.ReplaceAll(text, "", "...")
	text = strings.ReplaceAll(text, "", "'")
	text = strings.ReplaceAll(text, "", "'")
	text = strings.ReplaceAll(text, "", "\"")
	text = strings.ReplaceAll(text, "", "\"")

	return strings.TrimSpace(text)
}