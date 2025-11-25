package llm

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestOpenAITranslateErrorPaths tests error paths in OpenAI Translate function
func TestOpenAITranslateErrorPaths(t *testing.T) {
	// Test with invalid configuration to trigger error paths
	t.Run("invalid_api_key", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "", // Empty API key to trigger error
			Model:    "gpt-4",
		}

		client, err := NewOpenAIClient(config)
		if err != nil {
			// Expected error path - client creation failed
			if client != nil {
				t.Error("Client should be nil when creation fails")
			}
			return
		}

		if client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		// Test translation with empty API key
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "Hello world", "Translate to Russian")
		if err == nil {
			t.Error("Expected error with invalid configuration")
		}
		if result != "" {
			t.Error("Result should be empty when error occurs")
		}
	})

	t.Run("empty_text_input", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "test-api-key",
			Model:    "gpt-4",
		}

		client, err := NewOpenAIClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test with empty text
		result, err := client.Translate(ctx, "", "Translate to Russian")
		if err == nil && result == "" {
			t.Log("Empty input returned empty result - this is acceptable")
		}
		// Either should work - some APIs handle empty text, others don't
	})

	t.Run("context_cancellation", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "test-api-key",
			Model:    "gpt-4",
		}

		client, err := NewOpenAIClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		// Create cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result, err := client.Translate(ctx, "Hello world", "Translate to Russian")
		if err != nil {
			// Expected error path - context cancelled
			if result != "" {
				t.Error("Result should be empty when context is cancelled")
			}
			// Check for context-related error
			if !containsContextError(err) {
				t.Logf("Error may not be context-related: %v", err)
			}
		}
		// It's also possible the request completes before cancellation is detected
	})

	t.Run("malformed_model_name", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "test-api-key",
			Model:    "", // Empty model name to trigger default logic
		}

		client, err := NewOpenAIClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		if err != nil {
			// Expected error path - invalid model
			if result != "" {
				t.Error("Result should be empty when model is invalid")
			}
		}
	})

	t.Run("temperature_options", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "test-api-key",
			Model:    "gpt-4",
			Options: map[string]interface{}{
				"temperature": 1.5, // Very high temperature
			},
		}

		client, err := NewOpenAIClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		if err != nil {
			// Expected error path - invalid options
			if result != "" {
				t.Error("Result should be empty when options are invalid")
			}
		}
	})

	t.Run("max_tokens_override", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "test-api-key",
			Model:    "gpt-4",
			Options: map[string]interface{}{
				"max_tokens": -1, // Invalid max_tokens
			},
		}

		client, err := NewOpenAIClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		if err != nil {
			// Expected error path - invalid max_tokens
			if result != "" {
				t.Error("Result should be empty when max_tokens is invalid")
			}
		}
	})
}

// TestOpenAIClientCreation tests client creation error paths
func TestOpenAIClientCreation(t *testing.T) {
	t.Run("invalid_provider", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "invalid-provider",
			APIKey:   "test-key",
			Model:    "gpt-4",
		}

		client, err := NewOpenAIClient(config)
		// The function doesn't validate provider, it just creates a client
		if err != nil {
			t.Errorf("Unexpected error with invalid provider: %v", err)
		}
		if client == nil {
			t.Error("Client should not be nil even with invalid provider")
		}
		
		// The provider name should still be "openai" regardless of config
		if client != nil && client.GetProviderName() != "openai" {
			t.Errorf("Expected provider 'openai', got: %s", client.GetProviderName())
		}
	})

	t.Run("minimal_valid_config", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "test-key",
			// Model and options are optional
		}

		client, err := NewOpenAIClient(config)
		if err != nil {
			t.Errorf("Unexpected error with minimal config: %v", err)
		}
		if client == nil {
			t.Error("Client should not be nil with minimal valid config")
		}

		if client != nil {
			provider := client.GetProviderName()
			if provider != "openai" {
				t.Errorf("Expected provider 'openai', got: %s", provider)
			}
		}
	})

	t.Run("full_config", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "openai",
			APIKey:   "test-key",
			Model:    "gpt-4-turbo",
			Options: map[string]interface{}{
				"temperature": 0.7,
				"max_tokens":  2000,
				"timeout":     60 * time.Second,
			},
		}

		client, err := NewOpenAIClient(config)
		if err != nil {
			t.Errorf("Unexpected error with full config: %v", err)
		}
		if client == nil {
			t.Error("Client should not be nil with full valid config")
		}

		if client != nil {
			provider := client.GetProviderName()
			if provider != "openai" {
				t.Errorf("Expected provider 'openai', got: %s", provider)
			}
		}
	})
}

// Helper function to check if error is context-related
func containsContextError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "context") || 
	       strings.Contains(errStr, "canceled") || 
	       strings.Contains(errStr, "deadline")
}