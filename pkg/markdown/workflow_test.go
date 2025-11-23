package markdown

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"digital.vasic.translator/pkg/ebook"
	"digital.vasic.translator/pkg/translator"
)

// MockLLMProvider for testing
type MockLLMProvider struct {
	translations map[string]string
	delay       time.Duration
}

func (m *MockLLMProvider) Translate(ctx context.Context, req translator.TranslationRequest) (*translator.TranslationResponse, error) {
	// Add delay if specified
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Return mapped translation
	if translation, exists := m.translations[req.Text]; exists {
		return &translator.TranslationResponse{
			TargetText:     translation,
			SourceLang:     req.SourceLang,
			TargetLang:     req.TargetLang,
			Provider:       "mock",
			Confidence:     1.0,
			ProcessingTime: time.Millisecond * 100,
		}, nil
	}

	// Default translation
	return &translator.TranslationResponse{
		TargetText:     "[TRANSLATED: " + req.Text + "]",
		SourceLang:     req.SourceLang,
		TargetLang:     req.TargetLang,
		Provider:       "mock",
		Confidence:     0.8,
		ProcessingTime: time.Millisecond * 100,
	}, nil
}

func TestNewMarkdownWorkflow(t *testing.T) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	if workflow.llmProvider != provider {
		t.Error("LLM provider not set correctly")
	}
	
	if workflow.translator == nil {
		t.Error("Translator should not be nil")
	}
}

func TestMarkdownWorkflow_SetProgressCallback(t *testing.T) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	var callbackCalled bool
	callback := func(stage string, progress float64, message string) {
		callbackCalled = true
		if stage != "test" {
			t.Errorf("Expected stage 'test', got '%s'", stage)
		}
		if progress != 0.5 {
			t.Errorf("Expected progress 0.5, got %f", progress)
		}
		if message != "test message" {
			t.Errorf("Expected message 'test message', got '%s'", message)
		}
	}
	
	workflow.SetProgressCallback(callback)
	workflow.reportProgress("test", 0.5, "test message")
	
	if !callbackCalled {
		t.Error("Progress callback was not called")
	}
}

func TestMarkdownWorkflow_ConvertEbookToMarkdown(t *testing.T) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	// Create test book
	book := &ebook.Book{
		Metadata: ebook.Metadata{
			Title:       "Test Book",
			Description: "A test book for translation",
			Authors: []ebook.Author{
				{FirstName: "Test", LastName: "Author"},
			},
		},
		Chapters: []ebook.Chapter{
			{
				Title:   "Chapter 1",
				Content: "This is the first chapter with some content.",
			},
			{
				Title:   "Chapter 2",
				Content: "This is the second chapter with more content.",
			},
		},
	}
	
	markdown, err := workflow.convertEbookToMarkdown(book)
	if err != nil {
		t.Fatalf("Failed to convert ebook to markdown: %v", err)
	}
	
	// Check title
	if !strings.Contains(markdown, "# Test Book") {
		t.Error("Markdown should contain book title")
	}
	
	// Check authors
	if !strings.Contains(markdown, "Test Author") {
		t.Error("Markdown should contain author name")
	}
	
	// Check chapters
	if !strings.Contains(markdown, "## Chapter 1: Chapter 1") {
		t.Error("Markdown should contain chapter 1 title")
	}
	
	if !strings.Contains(markdown, "## Chapter 2: Chapter 2") {
		t.Error("Markdown should contain chapter 2 title")
	}
	
	// Check content
	if !strings.Contains(markdown, "This is the first chapter") {
		t.Error("Markdown should contain chapter 1 content")
	}
	
	if !strings.Contains(markdown, "This is the second chapter") {
		t.Error("Markdown should contain chapter 2 content")
	}
}

func TestMarkdownWorkflow_ConvertTextToMarkdown(t *testing.T) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Simple text",
			expected: "Simple text",
		},
		{
			input:    "Text with\ttabs",
			expected: "Text with  tabs",
		},
		{
			input:    "Multiple\n\n\nlines",
			expected: "Multiple\n\nlines",
		},
	}
	
	for _, test := range tests {
		result := workflow.convertTextToMarkdown(test.input)
		if result != test.expected {
			t.Errorf("convertTextToMarkdown(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestMarkdownWorkflow_SplitMarkdownIntoChunks(t *testing.T) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	// Test with small content
	small := "This is a small content."
	chunks := workflow.splitMarkdownIntoChunks(small)
	if len(chunks) != 1 {
		t.Errorf("Expected 1 chunk for small content, got %d", len(chunks))
	}
	
	// Test with large content (simulate with repeated text)
	large := strings.Repeat("This is a longer content that should be split into multiple chunks. ", 100)
	chunks = workflow.splitMarkdownIntoChunks(large)
	if len(chunks) < 2 {
		t.Errorf("Expected at least 2 chunks for large content, got %d", len(chunks))
	}
	
	// Check that chunks are reasonable size
	for i, chunk := range chunks {
		if len(chunk) > 5000 { // Allow some flexibility over 3000 target
			t.Errorf("Chunk %d is too large: %d characters", i, len(chunk))
		}
	}
}

func TestMarkdownWorkflow_MergeSmallChunks(t *testing.T) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	// Create many small chunks
	smallChunks := make([]string, 60)
	for i := 0; i < 60; i++ {
		smallChunks[i] = "Small chunk " + string(rune(i))
	}
	
	merged := workflow.mergeSmallChunks(smallChunks)
	
	// Should have fewer chunks after merging
	if len(merged) >= len(smallChunks) {
		t.Error("Expected fewer chunks after merging")
	}
	
	// Check that merged chunks are reasonable size
	for i, chunk := range merged {
		if len(chunk) > 3000 {
			t.Errorf("Merged chunk %d is too large: %d characters", i, len(chunk))
		}
	}
}

func TestMarkdownWorkflow_ParseMarkdownToEbook(t *testing.T) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	markdown := `# Test Book

**Authors:** Test Author

---

## Chapter 1: First Chapter

This is the content of chapter 1.

## Chapter 2: Second Chapter

This is the content of chapter 2.
`
	
	book, err := workflow.parseMarkdownToEbook(markdown, "en")
	if err != nil {
		t.Fatalf("Failed to parse markdown to ebook: %v", err)
	}
	
	// Check metadata
	if book.Metadata.Title != "Test Book" {
		t.Errorf("Expected title 'Test Book', got '%s'", book.Metadata.Title)
	}
	
	if book.Metadata.Language != "en" {
		t.Errorf("Expected language 'en', got '%s'", book.Metadata.Language)
	}
	
	if len(book.Metadata.Authors) != 1 {
		t.Errorf("Expected 1 author, got %d", len(book.Metadata.Authors))
	}
	
	if book.Metadata.Authors[0].FirstName != "Test" || book.Metadata.Authors[0].LastName != "Author" {
		t.Errorf("Expected author 'Test Author', got '%s %s'", 
			book.Metadata.Authors[0].FirstName, book.Metadata.Authors[0].LastName)
	}
	
	// Check chapters
	if len(book.Chapters) != 2 {
		t.Errorf("Expected 2 chapters, got %d", len(book.Chapters))
	}
	
	if book.Chapters[0].Title != "First Chapter" {
		t.Errorf("Expected chapter 1 title 'First Chapter', got '%s'", book.Chapters[0].Title)
	}
	
	if book.Chapters[1].Title != "Second Chapter" {
		t.Errorf("Expected chapter 2 title 'Second Chapter', got '%s'", book.Chapters[1].Title)
	}
}

func TestMarkdownWorkflow_TranslateMarkdown(t *testing.T) {
	translations := map[string]string{
		"Hello": "Привет",
		"World": "Мир",
	}
	
	provider := &MockLLMProvider{
		translations: translations,
		delay:       time.Millisecond * 10, // Small delay for testing
	}
	
	workflow := NewMarkdownWorkflow(provider)
	
	markdown := "# Hello\n\nThis is a test World."
	
	ctx := context.Background()
	translated, err := workflow.translateMarkdown(ctx, markdown, "en", "ru")
	if err != nil {
		t.Fatalf("Failed to translate markdown: %v", err)
	}
	
	// Check translations were applied
	if !strings.Contains(translated, "Привет") {
		t.Error("Expected translation of 'Hello' to be present")
	}
	
	if !strings.Contains(translated, "Мир") {
		t.Error("Expected translation of 'World' to be present")
	}
}

func TestMarkdownWorkflow_TranslateEbook(t *testing.T) {
	translations := map[string]string{
		"Test Book":       "Тестовая Книга",
		"Test Author":     "Тестовый Автор",
		"Chapter content": "Содержание главы",
	}
	
	provider := &MockLLMProvider{
		translations: translations,
		delay:       time.Millisecond * 10,
	}
	
	workflow := NewMarkdownWorkflow(provider)
	
	// Create test input file
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "test.fb2")
	outputPath := filepath.Join(tempDir, "test_sr.epub")
	
	// Create a simple test file
	testContent := `<FictionBook>
	<description>
	<title-info>
		<genre>test</genre>
		<author>
			<first-name>Test</first-name>
			<last-name>Author</last-name>
		</author>
		<book-title>Test Book</book-title>
	</title-info>
</description>
<body>
	<section>
	<title>Test Chapter</title>
	<p>Chapter content</p>
</section>
</body>
</FictionBook>`
	
	err := os.WriteFile(inputPath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}
	
	ctx := context.Background()
	err = workflow.TranslateEbook(ctx, inputPath, outputPath, "en", "sr")
	if err != nil {
		t.Fatalf("Failed to translate ebook: %v", err)
	}
	
	// Check that output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	// Check that markdown files were created
	originalMarkdownPath := filepath.Join(tempDir, "test_original.md")
	translatedMarkdownPath := filepath.Join(tempDir, "test_translated.md")
	
	if _, err := os.Stat(originalMarkdownPath); os.IsNotExist(err) {
		t.Error("Original markdown file was not created")
	}
	
	if _, err := os.Stat(translatedMarkdownPath); os.IsNotExist(err) {
		t.Error("Translated markdown file was not created")
	}
}

func TestWorkflowSummary(t *testing.T) {
	summary := &WorkflowSummary{
		InputPath:  "/path/to/input.fb2",
		OutputPath: "/path/to/output.epub",
		Stages:     []string{"Parse", "Translate", "Convert"},
		StartTime:   time.Now(),
	}
	
	// Test completion
	summary.Complete(true, nil)
	
	if !summary.Success {
		t.Error("Expected success to be true")
	}
	
	if summary.Duration == 0 {
		t.Error("Expected duration to be set")
	}
	
	if summary.Error != "" {
		t.Error("Expected error to be empty on success")
	}
	
	// Test completion with error
	summary.StartTime = time.Now()
	summary.Complete(false, fmt.Errorf("test error"))
	
	if summary.Success {
		t.Error("Expected success to be false on error")
	}
	
	if summary.Error != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", summary.Error)
	}
}

// Benchmark tests
func BenchmarkMarkdownWorkflow_ConvertEbookToMarkdown(b *testing.B) {
	provider := &MockLLMProvider{}
	workflow := NewMarkdownWorkflow(provider)
	
	// Create a test book
	book := &ebook.Book{
		Metadata: ebook.Metadata{
			Title: "Benchmark Test Book",
			Authors: []ebook.Author{
				{FirstName: "Test", LastName: "Author"},
			},
		},
		Chapters: make([]ebook.Chapter, 10), // 10 chapters
	}
	
	for i := range book.Chapters {
		book.Chapters[i] = ebook.Chapter{
			Title:   fmt.Sprintf("Chapter %d", i+1),
			Content: strings.Repeat("This is test content for benchmarking. ", 50),
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := workflow.convertEbookToMarkdown(book)
		if err != nil {
			b.Fatalf("Failed to convert ebook to markdown: %v", err)
		}
	}
}