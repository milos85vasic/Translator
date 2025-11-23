package ebook

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"io"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Set up uniDoc license (community edition)
	err := license.SetCommunityLicense()
	if err != nil {
		// Log error but continue - community edition will still work
		fmt.Printf("Warning: Could not set community license: %v\n", err)
	}
}

type PDFParser struct {
	config *PDFConfig
}

type PDFConfig struct {
	ExtractImages     bool  `yaml:"extract_images"`
	ImageFormat       string `yaml:"image_format"`
	OcrEnabled       bool   `yaml:"ocr_enabled"`
	OcrLanguage      string `yaml:"ocr_language"`
	PreserveLayout    bool   `yaml:"preserve_layout"`
	MinTextLength    int    `yaml:"min_text_length"`
	ExtractMetadata   bool   `yaml:"extract_metadata"`
	ExtractTables    bool   `yaml:"extract_tables"`
}

type PDFContent struct {
	Text       string      `json:"text"`
	Images     []PDFImage  `json:"images,omitempty"`
	Tables     []PDFTable  `json:"tables,omitempty"`
	Metadata   PDFMetadata `json:"metadata,omitempty"`
	Layout     PDFLayout   `json:"layout,omitempty"`
}

type PDFImage struct {
	Data     []byte `json:"data"`
	Format   string `json:"format"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Page     int    `json:"page"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"`
}

type PDFTable struct {
	Cells    [][]string `json:"cells"`
	Page     int       `json:"page"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"position"`
}

type PDFMetadata struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Subject     string    `json:"subject"`
	Creator     string    `json:"creator"`
	Producer    string    `json:"producer"`
	CreationDate string    `json:"creation_date"`
	ModDate     string    `json:"mod_date"`
	Pages       int       `json:"pages"`
	Language    string    `json:"language"`
	Keywords    []string  `json:"keywords"`
}

type PDFLayout struct {
	PageWidth  float64 `json:"page_width"`
	PageHeight float64 `json:"page_height"`
	Margin     struct {
		Top    float64 `json:"top"`
		Bottom float64 `json:"bottom"`
		Left   float64 `json:"left"`
		Right  float64 `json:"right"`
	} `json:"margin"`
	Columns   int `json:"columns"`
}

func NewPDFParser(config *PDFConfig) *PDFParser {
	if config == nil {
		config = &PDFConfig{
			ExtractImages:   true,
			ImageFormat:     "png",
			OcrEnabled:      false,
			OcrLanguage:     "eng",
			PreserveLayout:   true,
			MinTextLength:   5,
			ExtractMetadata:  true,
			ExtractTables:    true,
		}
	}
	return &PDFParser{config: config}
}

func (p *PDFParser) Parse(ctx context.Context, data []byte) (*Ebook, error) {
	pdfReader, err := model.NewPdfReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	ebook := &Ebook{
		Format: "pdf",
		Metadata: Metadata{
			Format: "pdf",
		},
	}

	// Extract metadata
	if p.config.ExtractMetadata {
		if err := p.extractMetadata(pdfReader, ebook); err != nil {
			return nil, fmt.Errorf("failed to extract metadata: %w", err)
		}
	}

	// Get page count
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	if numPages == 0 {
		return nil, fmt.Errorf("PDF has no pages")
	}

	// Extract content from all pages
	var allText strings.Builder
	var images []PDFImage
	var tables []PDFTable

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		// Extract text
		pageText, err := p.extractPageText(page, pageNum)
		if err != nil {
			return nil, fmt.Errorf("failed to extract text from page %d: %w", pageNum, err)
		}

		if len(pageText) >= p.config.MinTextLength {
			if allText.Len() > 0 {
				allText.WriteString("\n\n") // Page separator
			}
			allText.WriteString(pageText)
		}

		// Extract images if enabled
		if p.config.ExtractImages {
			pageImages, err := p.extractPageImages(page, pageNum)
			if err != nil {
				// Log error but continue
				fmt.Printf("Warning: Failed to extract images from page %d: %v\n", pageNum, err)
			} else {
				images = append(images, pageImages...)
			}
		}

		// Extract tables if enabled
		if p.config.ExtractTables {
			pageTables, err := p.extractPageTables(page, pageNum)
			if err != nil {
				// Log error but continue
				fmt.Printf("Warning: Failed to extract tables from page %d: %v\n", pageNum, err)
			} else {
				tables = append(tables, pageTables...)
			}
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	// Set main content
	if allText.Len() == 0 && !p.config.OcrEnabled {
		return nil, fmt.Errorf("no text found in PDF and OCR is disabled")
	}

	// Fallback to OCR if enabled and no text found
	if allText.Len() == 0 && p.config.OcrEnabled {
		ocrText, err := p.performOCR(ctx, pdfReader)
		if err != nil {
			return nil, fmt.Errorf("OCR failed: %w", err)
		}
		allText.WriteString(ocrText)
	}

	ebook.Content = allText.String()

	// Add extracted resources
	if len(images) > 0 {
		ebook.Resources = make(map[string]interface{})
		ebook.Resources["images"] = images
	}

	if len(tables) > 0 {
		if ebook.Resources == nil {
			ebook.Resources = make(map[string]interface{})
		}
		ebook.Resources["tables"] = tables
	}

	return ebook, nil
}

func (p *PDFParser) Validate(data []byte) error {
	// Check for PDF header
	if len(data) < 5 {
		return fmt.Errorf("file too small to be a valid PDF")
	}

	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return fmt.Errorf("invalid PDF header")
	}

	// Try to create PDF reader to validate structure
	_, err := model.NewPdfReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("invalid PDF structure: %w", err)
	}

	return nil
}

func (p *PDFParser) SupportedFormats() []string {
	return []string{"pdf", "application/pdf"}
}

func (p *PDFParser) GetMetadata(data []byte) (*Metadata, error) {
	pdfReader, err := model.NewPdfReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	metadata := &Metadata{Format: "pdf"}
	err = p.extractMetadata(pdfReader, &Ebook{Metadata: *metadata})
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}

	return metadata, nil
}

func (p *PDFParser) extractMetadata(pdfReader *model.PdfReader, ebook *Ebook) error {
	pdfInfo, err := pdfReader.GetPdfInfo()
	if err != nil {
		return fmt.Errorf("failed to get PDF info: %w", err)
	}

	// Extract basic metadata
	if title, err := pdfInfo.GetTitle(); err == nil && title != "" {
		ebook.Metadata.Title = title
	}

	if author, err := pdfInfo.GetAuthor(); err == nil && author != "" {
		ebook.Metadata.Author = []string{author}
	}

	if subject, err := pdfInfo.GetSubject(); err == nil && subject != "" {
		ebook.Metadata.Subject = subject
	}

	if creator, err := pdfInfo.GetCreator(); err == nil && creator != "" {
		ebook.Metadata.Publisher = creator
	}

	if producer, err := pdfInfo.GetProducer(); err == nil && producer != "" {
		if ebook.Metadata.Publisher == "" {
			ebook.Metadata.Publisher = producer
		}
	}

	// Extract dates
	if creationDate, err := pdfInfo.GetCreationDate(); err == nil {
		ebook.Metadata.Date = creationDate.String()
	}

	if modDate, err := pdfInfo.GetModDate(); err == nil {
		ebook.Metadata.Modified = modDate.String()
	}

	// Extract page count
	if numPages, err := pdfReader.GetNumPages(); err == nil {
		ebook.Metadata.Pages = &numPages
	}

	// Extract keywords
	if keywords, err := pdfInfo.GetKeywords(); err == nil && keywords != "" {
		ebook.Metadata.Keywords = strings.Split(keywords, ",")
		for i, kw := range ebook.Metadata.Keywords {
			ebook.Metadata.Keywords[i] = strings.TrimSpace(kw)
		}
	}

	// Extract language (if available)
	if lang, err := pdfInfo.GetLanguage(); err == nil && lang != "" {
		ebook.Metadata.Language = lang
	}

	return nil
}

func (p *PDFParser) extractPageText(page *model.PdfPage, pageNum int) (string, error) {
	ex, err := extractor.New(page)
	if err != nil {
		return "", fmt.Errorf("failed to create extractor: %w", err)
	}

	text, err := ex.ExtractText()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	// Clean up extracted text
	text = p.cleanText(text)

	if p.config.PreserveLayout {
		return p.preserveLayout(text, pageNum), nil
	}

	return text, nil
}

func (p *PDFParser) extractPageImages(page *model.PdfPage, pageNum int) ([]PDFImage, error) {
	var images []PDFImage

	// Get page images using unidoc
	pageImages, err := extractor.GetImages(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get page images: %w", err)
	}

	for i, pageImage := range pageImages {
		// Convert to standard image format
		img, format, err := pageImage.Image.GetImage()
		if err != nil {
			continue // Skip images that can't be converted
		}

		// Get image bounds
		bounds := img.Bounds()

		// Encode image
		var buf bytes.Buffer
		switch p.config.ImageFormat {
		case "png":
			err = p.encodePNG(&buf, img)
		case "jpeg":
			err = p.encodeJPEG(&buf, img)
		default:
			err = p.encodePNG(&buf, img)
		}

		if err != nil {
			continue // Skip images that can't be encoded
		}

		images = append(images, PDFImage{
			Data:   buf.Bytes(),
			Format: format,
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
			Page:   pageNum,
			Position: struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			}{
				X: pageImage.X,
				Y: pageImage.Y,
			},
		})
	}

	return images, nil
}

func (p *PDFParser) extractPageTables(page *model.PdfPage, pageNum int) ([]PDFTable, error) {
	// This is a simplified table extraction
	// In production, you would use more sophisticated table detection
	var tables []PDFTable

	// Try to extract tabular content from the page
	ex, err := extractor.New(page)
	if err != nil {
		return nil, fmt.Errorf("failed to create extractor: %w", err)
	}

	// This is a placeholder - real table extraction is complex
	// You would need to analyze text positions and identify tabular structures
	_ = ex

	return tables, nil
}

func (p *PDFParser) performOCR(ctx context.Context, pdfReader *model.PdfReader) (string, error) {
	// OCR implementation would go here
	// This requires OCR engine integration (Tesseract, etc.)
	// For now, return empty string as OCR is not implemented
	return "", fmt.Errorf("OCR not implemented in this build")
}

func (p *PDFParser) cleanText(text string) string {
	// Remove excessive whitespace
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\n ", "\n")
	text = strings.ReplaceAll(text, "  ", " ")

	// Remove common PDF artifacts
	text = strings.ReplaceAll(text, "", "") // Form feed
	text = strings.ReplaceAll(text, "", "\"")
	text = strings.ReplaceAll(text, "", "...")
	text = strings.ReplaceAll(text, "", "'")
	text = strings.ReplaceAll(text, "", "'")
	text = strings.ReplaceAll(text, "", "\"")
	text = strings.ReplaceAll(text, "", "\"")

	return strings.TrimSpace(text)
}

func (p *PDFParser) preserveLayout(text string, pageNum int) string {
	// Basic layout preservation
	// In production, you would preserve more layout information
	lines := strings.Split(text, "\n")
	
	var preservedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			preservedLines = append(preservedLines, line)
		}
	}

	return strings.Join(preservedLines, "\n")
}

func (p *PDFParser) encodePNG(w io.Writer, img image.Image) error {
	// This would use png.Encode in real implementation
	// Placeholder for encoding
	return fmt.Errorf("PNG encoding not implemented")
}

func (p *PDFParser) encodeJPEG(w io.Writer, img image.Image) error {
	// This would use jpeg.Encode in real implementation
	// Placeholder for encoding
	return fmt.Errorf("JPEG encoding not implemented")
}