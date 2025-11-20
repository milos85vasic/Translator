package coordination

import (
	"context"
	"digital.vasic.translator/pkg/events"
	"digital.vasic.translator/pkg/translator"
)

// MultiLLMTranslatorWrapper wraps MultiLLMCoordinator to implement the Translator interface
type MultiLLMTranslatorWrapper struct {
	Coordinator *MultiLLMCoordinator // Exported so CLI can access instance count
	config      translator.TranslationConfig
}

// NewMultiLLMTranslatorWrapper creates a new wrapper
func NewMultiLLMTranslatorWrapper(config translator.TranslationConfig, eventBus *events.EventBus, sessionID string) (*MultiLLMTranslatorWrapper, error) {
	coordinator := NewMultiLLMCoordinator(CoordinatorConfig{
		MaxRetries: 3,
		RetryDelay: 0, // No delay between retries with different instances
		EventBus:   eventBus,
		SessionID:  sessionID,
	})

	if coordinator.GetInstanceCount() == 0 {
		// Fall back to single translator if no instances available
		return nil, translator.ErrNoLLMInstances
	}

	return &MultiLLMTranslatorWrapper{
		Coordinator: coordinator,
		config:      config,
	}, nil
}

// Translate implements translator.Translator
func (w *MultiLLMTranslatorWrapper) Translate(ctx context.Context, text string, context string) (string, error) {
	return w.Coordinator.TranslateWithRetry(ctx, text, context)
}

// TranslateWithProgress implements translator.Translator
func (w *MultiLLMTranslatorWrapper) TranslateWithProgress(
	ctx context.Context,
	text string,
	contextHint string,
	eventBus *events.EventBus,
	sessionID string,
) (string, error) {
	return w.Coordinator.TranslateWithRetry(ctx, text, contextHint)
}

// GetName implements translator.Translator
func (w *MultiLLMTranslatorWrapper) GetName() string {
	return "multi-llm-coordinator"
}

// GetStats implements translator.Translator
func (w *MultiLLMTranslatorWrapper) GetStats() translator.TranslationStats {
	// Multi-LLM coordinator doesn't track individual stats the same way
	// Return zero stats for now - proper stats tracking can be added later
	return translator.TranslationStats{
		Total:      0,
		Translated: 0,
		Cached:     0,
		Errors:     0,
	}
}
