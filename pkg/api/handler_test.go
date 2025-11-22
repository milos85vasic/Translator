package api

import (
	"testing"
)

func TestGenerateOutputFilename(t *testing.T) {
	tests := []struct {
		input    string
		provider string
		expected string
	}{
		{"book.fb2", "dictionary", "book_sr_dictionary.fb2"},
		{"test.b2", "openai", "test_sr_openai.b2"},
		{"novel.fb2", "anthropic", "novel_sr_anthropic.fb2"},
	}

	for _, tt := range tests {
		t.Run(tt.input+"_"+tt.provider, func(t *testing.T) {
			result := generateOutputFilename(tt.input, tt.provider)
			if result != tt.expected {
				t.Errorf("generateOutputFilename(%s, %s) = %s, want %s", tt.input, tt.provider, result, tt.expected)
			}
		})
	}
}
