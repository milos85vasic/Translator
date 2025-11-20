package language

import (
	"context"
	"fmt"
	"strings"
)

// SimpleLLMDetector implements LLM-based language detection
type SimpleLLMDetector struct {
	apiKey   string
	provider string
}

// NewSimpleLLMDetector creates a new LLM detector
func NewSimpleLLMDetector(provider, apiKey string) *SimpleLLMDetector {
	return &SimpleLLMDetector{
		apiKey:   apiKey,
		provider: provider,
	}
}

// DetectLanguage detects language using LLM
func (d *SimpleLLMDetector) DetectLanguage(ctx context.Context, text string) (string, error) {
	// Sample text (first 500 characters)
	sample := text
	if len(text) > 500 {
		sample = text[:500]
	}

	// Create prompt for language detection
	prompt := fmt.Sprintf(`Identify the language of the following text.
Respond with ONLY the ISO 639-1 language code (e.g., "en" for English, "ru" for Russian, "de" for German).
Do not include any explanation, just the 2-letter code.

Text:
%s

Language code:`, sample)

	// For now, return empty to use fallback
	// In a full implementation, this would call the LLM API
	_ = prompt

	return "", fmt.Errorf("LLM detection not implemented")
}

// FormatLanguageCode normalizes language codes
func FormatLanguageCode(code string) string {
	code = strings.TrimSpace(strings.ToLower(code))

	// Handle common variations
	if len(code) > 2 {
		code = code[:2]
	}

	return code
}
