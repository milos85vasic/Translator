package security

import (
	"context"
	"strings"

	"digital.vasic.translator/pkg/events"
	"digital.vasic.translator/pkg/translator"
)

// MockTranslator implements the translator.Translator interface for security testing
type MockTranslator struct {
	stats translator.TranslationStats
	name  string
}

func NewMockTranslator() *MockTranslator {
	return &MockTranslator{
		stats: translator.TranslationStats{},
		name:  "mock",
	}
}

func (m *MockTranslator) Translate(ctx context.Context, text string, context string) (string, error) {
	// For security testing, return a simple translation that doesn't contain malicious content
	m.stats.Total++
	m.stats.Translated++

	// Sanitize input for security testing
	sanitized := text

	// Remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	// Basic XSS sanitization for testing
	sanitized = strings.ReplaceAll(sanitized, "<script>", "[script]")
	sanitized = strings.ReplaceAll(sanitized, "</script>", "[/script]")
	sanitized = strings.ReplaceAll(sanitized, "javascript:", "[javascript:]")
	sanitized = strings.ReplaceAll(sanitized, "onerror=", "[onerror=]")
	sanitized = strings.ReplaceAll(sanitized, "onload=", "[onload=]")
	sanitized = strings.ToLower(sanitized)
	// Check for remaining dangerous patterns
	if strings.Contains(sanitized, "<script") ||
		strings.Contains(sanitized, "javascript:") ||
		strings.Contains(sanitized, "onerror") ||
		strings.Contains(sanitized, "onload") {
		sanitized = "[sanitized]"
	}

	return "[translated: " + sanitized + "]", nil
}

func (m *MockTranslator) TranslateWithProgress(ctx context.Context, text string, context string, eventBus *events.EventBus, sessionID string) (string, error) {
	return m.Translate(ctx, text, context)
}

func (m *MockTranslator) GetStats() translator.TranslationStats {
	return m.stats
}

func (m *MockTranslator) GetName() string {
	return m.name
}
