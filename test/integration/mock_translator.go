//go:build integration
// +build integration

package integration

import (
	"context"
	"digital.vasic.translator/pkg/events"
	"digital.vasic.translator/pkg/translator"
)

// MockTranslator for testing
type MockTranslator struct{}

func (m *MockTranslator) Translate(ctx context.Context, text, context string) (string, error) {
	return "translated: " + text, nil
}

func (m *MockTranslator) TranslateWithProgress(ctx context.Context, text, context string, eventBus *events.EventBus, sessionID string) (string, error) {
	return m.Translate(ctx, text, context)
}

func (m *MockTranslator) GetStats() translator.TranslationStats {
	return translator.TranslationStats{
		Total:      0,
		Translated: 0,
		Cached:     0,
		Errors:     0,
	}
}

func (m *MockTranslator) GetName() string {
	return "mock"
}
