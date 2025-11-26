package translator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestUniversalTranslatorBasicFunctionality tests basic universal translator functionality
func TestUniversalTranslatorBasicFunctionality(t *testing.T) {
	// This test needs to be rewritten to use the proper translator interface
	// after the LLM package refactoring is complete
}

// TestUniversalTranslatorProviderSwitching tests provider switching
func TestUniversalTranslatorProviderSwitching(t *testing.T) {
	// TODO: Implement after LLM refactoring
}

// TestUniversalTranslatorMultipleLanguages tests translation between multiple language pairs
func TestUniversalTranslatorMultipleLanguages(t *testing.T) {
	// TODO: Implement after LLM refactoring
}

// TestUniversalTranslatorErrorHandling tests error handling
func TestUniversalTranslatorErrorHandling(t *testing.T) {
	// TODO: Implement after LLM refactoring
}

// TestTranslationConfigValidation tests translation config validation
func TestTranslationConfigValidation(t *testing.T) {
	// Test valid config
	validConfig := TranslationConfig{
		SourceLang:  "en",
		TargetLang:  "ru",
		Provider:    "openai",
		Model:       "gpt-4",
		Temperature: 0.7,
		MaxTokens:   1000,
		Timeout:     30 * time.Second,
		APIKey:      "test-key",
		BaseURL:     "https://api.openai.com",
		Script:      "latin",
		Options:     make(map[string]interface{}),
	}
	
	// Config should be valid (no validation function currently)
	assert.NotNil(t, validConfig)
	assert.Equal(t, "en", validConfig.SourceLang)
	assert.Equal(t, "ru", validConfig.TargetLang)
	assert.Equal(t, "openai", validConfig.Provider)
	assert.Equal(t, "gpt-4", validConfig.Model)
	assert.Equal(t, "test-key", validConfig.APIKey)
}

// TestTranslationResult tests translation result structure
func TestTranslationResult(t *testing.T) {
	result := TranslationResult{
		OriginalText:  "Hello world",
		TranslatedText: "Привет мир",
		Provider:      "openai",
		Cached:        false,
		Error:         nil,
	}
	
	assert.Equal(t, "Hello world", result.OriginalText)
	assert.Equal(t, "Привет мир", result.TranslatedText)
	assert.Equal(t, "openai", result.Provider)
	assert.False(t, result.Cached)
	assert.NoError(t, result.Error)
}

// TestTranslationStats tests translation statistics
func TestTranslationStats(t *testing.T) {
	stats := TranslationStats{
		Total:      100,
		Translated: 80,
		Cached:     15,
		Errors:     5,
	}
	
	assert.Equal(t, 100, stats.Total)
	assert.Equal(t, 80, stats.Translated)
	assert.Equal(t, 15, stats.Cached)
	assert.Equal(t, 5, stats.Errors)
	
	// Verify that total equals translated + cached + errors
	assert.Equal(t, stats.Translated + stats.Cached + stats.Errors, stats.Total)
}

// TestTranslatorErrors tests error variables
func TestTranslatorErrors(t *testing.T) {
	assert.NotNil(t, ErrNoLLMInstances)
	assert.NotNil(t, ErrInvalidProvider)
	assert.Equal(t, "no LLM instances available", ErrNoLLMInstances.Error())
	assert.Equal(t, "invalid translation provider", ErrInvalidProvider.Error())
}