package markdown

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"digital.vasic.translator/pkg/logger"
	"digital.vasic.translator/pkg/translator/llm"
)

// WorkflowConfig holds configuration for the markdown workflow
type WorkflowConfig struct {
	ChunkSize        int
	OverlapSize      int
	MaxConcurrency   int
	TranslationCache map[string]string
	LLMProvider      llm.LLMClient
}

// ProgressCallback function for progress updates
type ProgressCallback func(current, total int, message string)

// SimpleWorkflow handles ebook → markdown → translation → ebook workflow
type SimpleWorkflow struct {
	config     WorkflowConfig
	logger     logger.Logger
	callback   ProgressCallback
}

// NewSimpleWorkflow creates a new simple markdown workflow
func NewSimpleWorkflow(config WorkflowConfig, logger logger.Logger, callback ProgressCallback) *SimpleWorkflow {
	return &SimpleWorkflow{
		config:   config,
		logger:   logger,
		callback: callback,
	}
}

// ConvertToMarkdown converts an ebook file to markdown format
func (sw *SimpleWorkflow) ConvertToMarkdown(ctx context.Context, inputPath, outputPath string) error {
	sw.logger.Info("Converting ebook to markdown", map[string]interface{}{
		"input_path":  inputPath,
		"output_path": outputPath,
	})

	if sw.callback != nil {
		sw.callback(0, 100, "Converting ebook to markdown...")
	}

	// Determine file type and convert accordingly
	ext := strings.ToLower(filepath.Ext(inputPath))
	
	switch ext {
	case ".epub":
		return sw.convertEPUBToMarkdown(ctx, inputPath, outputPath)
	case ".fb2":
		return sw.convertFB2ToMarkdown(ctx, inputPath, outputPath)
	default:
		return fmt.Errorf("unsupported input format: %s", ext)
	}
}

// TranslateMarkdown translates a markdown file using the configured LLM provider
func (sw *SimpleWorkflow) TranslateMarkdown(ctx context.Context, inputPath, outputPath, fromLang, toLang string) error {
	sw.logger.Info("Translating markdown file", map[string]interface{}{
		"input_path":  inputPath,
		"output_path": outputPath,
		"from_lang":   fromLang,
		"to_lang":     toLang,
	})

	if sw.callback != nil {
		sw.callback(1, 100, "Translating markdown content...")
	}

	// Create translation function using the LLM provider
	translateFunc := func(text string) (string, error) {
		prompt := fmt.Sprintf(`Translate the following text from %s to %s. 
Provide ONLY the translation without any explanations, notes, or additional text.
Maintain the original formatting, line breaks, and structure.

Source text:
%s

Translation:`, fromLang, toLang, text)

		return sw.config.LLMProvider.Translate(ctx, text, prompt)
	}

	// Create markdown translator
	markdownTranslator := NewMarkdownTranslator(translateFunc)

	// Translate the markdown file
	if err := markdownTranslator.TranslateMarkdownFile(inputPath, outputPath); err != nil {
		return fmt.Errorf("failed to translate markdown: %w", err)
	}

	if sw.callback != nil {
		sw.callback(2, 100, "Markdown translation completed")
	}

	return nil
}

// ConvertFromMarkdown converts a markdown file to ebook format
func (sw *SimpleWorkflow) ConvertFromMarkdown(ctx context.Context, inputPath, outputPath string) error {
	sw.logger.Info("Converting markdown to ebook", map[string]interface{}{
		"input_path":  inputPath,
		"output_path": outputPath,
	})

	if sw.callback != nil {
		sw.callback(3, 100, "Converting markdown to ebook...")
	}

	// Determine output format and convert accordingly
	ext := strings.ToLower(filepath.Ext(outputPath))
	
	switch ext {
	case ".epub":
		return sw.convertMarkdownToEPUB(ctx, inputPath, outputPath)
	default:
		return fmt.Errorf("unsupported output format: %s", ext)
	}
}

// ExecuteFullWorkflow executes the complete conversion workflow
func (sw *SimpleWorkflow) ExecuteFullWorkflow(ctx context.Context, inputPath, outputPath, fromLang, toLang string) error {
	startTime := time.Now()
	
	// Generate intermediate file paths
	ext := strings.ToLower(filepath.Ext(inputPath))
	baseName := strings.TrimSuffix(inputPath, ext)
	
	originalMarkdownPath := baseName + "_original.md"
	translatedMarkdownPath := baseName + "_translated.md"

	// Step 1: Convert input ebook to markdown
	if err := sw.ConvertToMarkdown(ctx, inputPath, originalMarkdownPath); err != nil {
		return fmt.Errorf("failed to convert to markdown: %w", err)
	}

	// Step 2: Translate markdown
	if err := sw.TranslateMarkdown(ctx, originalMarkdownPath, translatedMarkdownPath, fromLang, toLang); err != nil {
		return fmt.Errorf("failed to translate markdown: %w", err)
	}

	// Step 3: Convert translated markdown to output format
	if err := sw.ConvertFromMarkdown(ctx, translatedMarkdownPath, outputPath); err != nil {
		return fmt.Errorf("failed to convert from markdown: %w", err)
	}

	duration := time.Since(startTime)
	sw.logger.Info("Full workflow completed successfully", map[string]interface{}{
		"duration": duration.String(),
		"input_path": inputPath,
		"output_path": outputPath,
	})

	if sw.callback != nil {
		sw.callback(100, 100, "Workflow completed successfully")
	}

	return nil
}

// convertEPUBToMarkdown converts EPUB to markdown
func (sw *SimpleWorkflow) convertEPUBToMarkdown(ctx context.Context, inputPath, outputPath string) error {
	converter := NewEPUBToMarkdownConverter(false, "") // Don't preserve images for simplicity
	return converter.ConvertEPUBToMarkdown(inputPath, outputPath)
}

// convertFB2ToMarkdown converts FB2 to markdown (simplified version)
func (sw *SimpleWorkflow) convertFB2ToMarkdown(ctx context.Context, inputPath, outputPath string) error {
	// For now, create a simple text extraction
	// In a real implementation, you'd parse the FB2 XML properly
	
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read FB2 file: %w", err)
	}

	// Simple text extraction (this is very basic)
	text := string(content)
	
	// Create basic markdown
	markdownContent := fmt.Sprintf("# Book Content\n\n%s\n", text)
	
	return os.WriteFile(outputPath, []byte(markdownContent), 0644)
}

// convertMarkdownToEPUB converts markdown to EPUB
func (sw *SimpleWorkflow) convertMarkdownToEPUB(ctx context.Context, inputPath, outputPath string) error {
	// Read markdown content
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read markdown file: %w", err)
	}

	// Simple EPUB creation (this would need proper implementation)
	// For now, just copy the content as text
	epubContent := fmt.Sprintf("EPUB Version of:\n\n%s", string(content))
	
	return os.WriteFile(outputPath, []byte(epubContent), 0644)
}