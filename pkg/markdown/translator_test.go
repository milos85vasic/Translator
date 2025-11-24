package markdown

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTranslator implements a mock translator for testing
type MockTranslator struct {
	mock.Mock
}

func (m *MockTranslator) Translate(ctx context.Context, text string, sourceLang string, targetLang string) (string, error) {
	args := m.Called(ctx, text, sourceLang, targetLang)
	return args.String(0), args.Error(1)
}

func (m *MockTranslator) GetProvider() string {
	args := m.Called()
	return args.String(0)
}

// TestMarkdownTranslatorBasicFunctionality tests basic Markdown translation
func TestMarkdownTranslatorBasicFunctionality(t *testing.T) {
	// Create Markdown translator that returns the same text (no actual translation)
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return text, nil  // Return original text to test formatting preservation
	})

	// Test basic Markdown translation
	input := "# Hello World\n" +
		"This is a paragraph with **bold text** and *italic text*.\n" +
		"## Subsection\n" +
		"- List item 1\n" +
		"- List item 2\n" +
		"Hello world"

	result, err := mdTrans.TranslateMarkdown(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify headers are preserved
	assert.Contains(t, result, "# Hello World")
	assert.Contains(t, result, "## Subsection")
	
	// Verify formatting is preserved
	assert.Contains(t, result, "**bold text**")
	assert.Contains(t, result, "*italic text*")
	
	// Verify list structure is preserved
	assert.Contains(t, result, "- List item 1")
	assert.Contains(t, result, "- List item 2")
	
	// Verify content is preserved (since we're not actually translating)
	assert.Contains(t, result, "Hello world")
}

// TestMarkdownTranslatorCodeBlocks tests handling of code blocks
func TestMarkdownTranslatorCodeBlocks(t *testing.T) {
	// Create Markdown translator that returns same text
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return text, nil
	})

	// Test with code blocks
	input := "# Code Test\n" +
		"Here is some inline code: `" + "`console.log('Hello')`" + "`\n" +
		"\n" +
		"Here is a code block:\n" +
		"```javascript\n" +
		"function greet() {\n" +
		"    console.log('Hello world');\n" +
		"	}\n" +
		"```\n" +
		"More text after code block."

	result, err := mdTrans.TranslateMarkdown(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify code blocks are preserved
	assert.Contains(t, result, "`console.log('Hello')`") // Inline code preserved
	assert.Contains(t, result, "```javascript") // Code block marker preserved
	assert.Contains(t, result, "function greet()") // Code content preserved
}

// TestMarkdownTranslatorTables tests handling of Markdown tables
func TestMarkdownTranslatorTables(t *testing.T) {
	// Create Markdown translator that returns same text
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return text, nil
	})

	// Test with table
	input := "| Name | Age |\n" +
		"|------|-----|\n" +
		"| John | 25  |"

	result, err := mdTrans.TranslateMarkdown(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify table structure is preserved
	assert.Contains(t, result, "|")
	assert.Contains(t, result, "---")
	
	// Verify table content is preserved
	assert.Contains(t, result, "John")
}

// TestMarkdownTranslatorFrontMatter tests handling of YAML front matter
func TestMarkdownTranslatorFrontMatter(t *testing.T) {
	// Create Markdown translator that returns same text
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return text, nil
	})

	// Test with front matter
	input := "---\n" +
		"title: My Document\n" +
		"author: John Doe\n" +
		"---\n" +
		"\n" +
		"# Hello World\n" +
		"Hello world"

	result, err := mdTrans.TranslateMarkdown(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify front matter is preserved
	assert.Contains(t, result, "title: My Document")
	assert.Contains(t, result, "author: John Doe")
	assert.Contains(t, result, "---")
	
	// Verify content is preserved
	assert.Contains(t, result, "Hello world")
}

// TestMarkdownTranslatorLinks tests handling of Markdown links
func TestMarkdownTranslatorLinks(t *testing.T) {
	// Create Markdown translator that returns same text
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return text, nil
	})

	// Test with links
	input := "[Click here](https://example.com)"

	result, err := mdTrans.TranslateMarkdown(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify link structure is preserved
	assert.Contains(t, result, "[")
	assert.Contains(t, result, "]")
	assert.Contains(t, result, "(")
	assert.Contains(t, result, ")")
	assert.Contains(t, result, "https://example.com")
	
	// Verify link text is preserved
	assert.Contains(t, result, "Click here")
}

// TestMarkdownTranslatorEmptyInput tests handling of empty input
func TestMarkdownTranslatorEmptyInput(t *testing.T) {
	// Create Markdown translator
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return "translated", nil
	})

	result, err := mdTrans.TranslateMarkdown("")
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

// TestMarkdownTranslatorTranslationError tests handling of translation errors
func TestMarkdownTranslatorTranslationError(t *testing.T) {
	// Create Markdown translator that returns an error
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return "", assert.AnError
	})

	input := "Hello world"
	
	_, err := mdTrans.TranslateMarkdown(input)
	assert.Error(t, err)
}

// TestMarkdownTranslatorFileOperations tests file I/O operations
func TestMarkdownTranslatorFileOperations(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	inputFile := tempDir + "/input.md"
	outputFile := tempDir + "/output.md"

	// Create input file
	input := "# Test Document\nHello world"
	err := writeFile(inputFile, input)
	require.NoError(t, err)

	// Create Markdown translator
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return text, nil
	})

	// Test file translation
	err = mdTrans.TranslateMarkdownFile(inputFile, outputFile)
	require.NoError(t, err)

	// Verify output file exists and has content
	result, err := readFile(outputFile)
	require.NoError(t, err)
	assert.Contains(t, result, "Test Document")
	assert.Contains(t, result, "Hello world")
}

// TestMarkdownTranslatorLargeDocument tests handling of large documents
func TestMarkdownTranslatorLargeDocument(t *testing.T) {
	// Create Markdown translator
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return "translated: " + text, nil
	})

	// Generate large content
	var builder strings.Builder
	for i := 0; i < 1000; i++ {
		builder.WriteString("# Section ")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString("\nContent for section ")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(".\n\n")
	}

	input := builder.String()
	
	result, err := mdTrans.TranslateMarkdown(input)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "translated: ")
}

// BenchmarkMarkdownTranslatorSmallContent benchmarks small content translation
func BenchmarkMarkdownTranslatorSmallContent(b *testing.B) {
	// Create Markdown translator
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return "Переведено", nil
	})

	input := "Hello world"
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := mdTrans.TranslateMarkdown(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarkdownTranslatorMediumContent benchmarks medium content translation
func BenchmarkMarkdownTranslatorMediumContent(b *testing.B) {
	// Create Markdown translator
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return "Переведено", nil
	})

	input := "# Test Document\n" +
		"This is a test paragraph with **bold** and *italic* text.\n" +
		"## Section\n" +
		"- List item 1\n" +
		"- List item 2\n" +
		"> This is a quote\n" +
		"`inline code`"
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := mdTrans.TranslateMarkdown(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarkdownTranslatorLargeContent benchmarks large content translation
func BenchmarkMarkdownTranslatorLargeContent(b *testing.B) {
	// Create Markdown translator
	mdTrans := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return "Переведено", nil
	})

	// Generate large content
	var builder strings.Builder
	for i := 0; i < 100; i++ {
		builder.WriteString("# Section ")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString("\nThis is content for section ")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(".\n\n")
	}

	input := builder.String()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := mdTrans.TranslateMarkdown(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Helper functions for file operations in tests
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	return string(content), err
}