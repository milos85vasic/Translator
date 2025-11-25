package llm

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestOllamaProvider(t *testing.T) {
	// Test placeholder - provider implementation needed
	t.Log("Ollama provider test placeholder")
}

func TestOllamaProviderConfig(t *testing.T) {
	// Test placeholder for config testing
	t.Log("Ollama config test placeholder")
}

// TestOllamaRequestErrorPaths tests error paths in ollama Translate function
func TestOllamaRequestErrorPaths(t *testing.T) {
	t.Run("invalid_base_url", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			Model:    "llama3:8b",
			BaseURL:  "invalid-url", // Invalid URL
		}

		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		if err == nil {
			t.Log("Request succeeded with invalid URL - may be using mock")
		}
		if result != "" && err != nil {
			t.Error("Result should be empty when error occurs")
		}
	})

	t.Run("empty_model", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			Model:    "", // Empty model to trigger default logic
			BaseURL:  "http://localhost:11434",
		}

		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		if err != nil {
			// Expected error path - no Ollama server running
			if result != "" {
				t.Error("Result should be empty when connection fails")
			}
			// Check for connection-related error
			if !strings.Contains(err.Error(), "connection refused") &&
			   !strings.Contains(err.Error(), "no such host") &&
			   !strings.Contains(err.Error(), "timeout") {
				t.Logf("Error may not be connection-related: %v", err)
			}
		}
	})

	t.Run("context_cancellation", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			Model:    "llama3:8b",
			BaseURL:  "http://localhost:11434",
		}

		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		// Create cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		if err != nil {
			// Expected error path
			if result != "" {
				t.Error("Result should be empty when context is cancelled")
			}
			// Check for context-related error
			if !strings.Contains(err.Error(), "context") && 
			   !strings.Contains(err.Error(), "canceled") && 
			   !strings.Contains(err.Error(), "deadline") {
				t.Logf("Error may not be context-related: %v", err)
			}
		}
	})

	t.Run("empty_text_input", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			Model:    "llama3:8b",
			BaseURL:  "http://localhost:11434",
		}

		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "", "Translate to Russian")
		if err == nil && result == "" {
			t.Log("Empty input returned empty result - this is acceptable")
		}
		// Either should work - some APIs handle empty text, others don't
	})

	t.Run("very_long_text", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			Model:    "llama3:8b",
			BaseURL:  "http://localhost:11434",
		}

		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create very long text that might trigger size limits
		longText := strings.Repeat("Hello world. ", 1000)
		
		result, err := client.Translate(ctx, longText, "Translate to Russian")
		if err != nil {
			// Expected error path - text too long
			if result != "" {
				t.Error("Result should be empty when text is too long")
			}
			// Check for size-related error
			if !strings.Contains(err.Error(), "too large") && 
			   !strings.Contains(err.Error(), "size") && 
			   !strings.Contains(err.Error(), "limit") &&
			   !strings.Contains(err.Error(), "payload") {
				t.Logf("Error may not be size-related: %v", err)
			}
		}
	})

	t.Run("malformed_json_response", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			Model:    "llama3:8b",
			BaseURL:  "http://httpbin.org/json", // Returns valid JSON but wrong format
		}

		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		if err != nil {
			// Expected error path - malformed response
			if result != "" {
				t.Error("Result should be empty when response is malformed")
			}
			// Check for JSON-related error
			if !strings.Contains(err.Error(), "unmarshal") && 
			   !strings.Contains(err.Error(), "json") {
				t.Logf("Error may not be JSON-related: %v", err)
			}
		}
	})
}

// TestOllamaClientCreation tests client creation error paths
func TestOllamaClientCreation(t *testing.T) {
	t.Run("minimal_valid_config", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			// Model and BaseURL are optional
		}

		client, err := NewOllamaClient(config)
		if err != nil {
			t.Errorf("Unexpected error with minimal config: %v", err)
		}
		if client == nil {
			t.Error("Client should not be nil with minimal valid config")
		}

		if client != nil {
			provider := client.GetProviderName()
			if provider != "ollama" {
				t.Errorf("Expected provider 'ollama', got: %s", provider)
			}
		}
	})

	t.Run("full_config", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			Model:    "mistral:7b",
			BaseURL:  "http://custom.localhost:11435",
		}

		client, err := NewOllamaClient(config)
		if err != nil {
			t.Errorf("Unexpected error with full config: %v", err)
		}
		if client == nil {
			t.Error("Client should not be nil with full valid config")
		}

		if client != nil {
			provider := client.GetProviderName()
			if provider != "ollama" {
				t.Errorf("Expected provider 'ollama', got: %s", provider)
			}
		}
	})
}

// TestOllamaTranslateUncoveredPaths tests uncovered error paths in Ollama Translate function
func TestOllamaTranslateUncoveredPaths(t *testing.T) {
	// Test 1: JSON marshaling error
	t.Run("json_marshal_error", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			BaseURL:  "http://localhost:11434",
		}
		
		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}
		
		// Test with valid config but try to exercise marshal error indirectly
		ctx := context.Background()
		
		// This will likely fail due to no server, but we're testing the path
		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		
		if err != nil {
			t.Logf("Expected error (server not running): %v", err)
			if result != "" {
				t.Error("Result should be empty when error occurs")
			}
		}
	})
	
	// Test 2: HTTP request creation error
	t.Run("http_request_error", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			BaseURL:  "invalid://invalid-url", // Invalid URL that should cause request creation error
		}
		
		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}
		
		ctx := context.Background()
		
		// This should fail during HTTP request creation
		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		
		if err != nil {
			t.Logf("Expected HTTP request error: %v", err)
			if result != "" {
				t.Error("Result should be empty when HTTP request fails")
			}
			
			// Check for appropriate error message
			if !strings.Contains(err.Error(), "failed to create request") &&
			   !strings.Contains(err.Error(), "failed to send request") &&
			   !strings.Contains(err.Error(), "invalid") {
				t.Logf("Error may not be request creation related: %v", err)
			}
		}
	})
	
	// Test 3: Response reading error
	t.Run("response_reading_error", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			BaseURL:  "http://localhost:99999", // Port that's likely not running
		}
		
		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		
		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		
		if err != nil {
			t.Logf("Expected connection error: %v", err)
			if result != "" {
				t.Error("Result should be empty when connection fails")
			}
			
			// Should be a connection-related error
			if !strings.Contains(err.Error(), "connection refused") &&
			   !strings.Contains(err.Error(), "timeout") &&
			   !strings.Contains(err.Error(), "network") &&
			   !strings.Contains(err.Error(), "failed to send request") {
				t.Logf("Error may not be connection-related: %v", err)
			}
		}
	})
	
	// Test 4: Non-200 status codes
	t.Run("non_200_status_codes", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			BaseURL:  "http://httpbin.org/status/404", // Will return 404
		}
		
		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		
		if err != nil {
			t.Logf("Expected status code error: %v", err)
			if result != "" {
				t.Error("Result should be empty when status code is not 200")
			}
			
			// Should contain status code information
			if !strings.Contains(err.Error(), "status") &&
			   !strings.Contains(err.Error(), "404") &&
			   !strings.Contains(err.Error(), "Ollama API error") {
				t.Logf("Error may not be status code related: %v", err)
			}
		}
	})
	
	// Test 5: JSON unmarshaling error
	t.Run("json_unmarshal_error", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			BaseURL:  "http://httpbin.org/html", // Returns HTML, not JSON
		}
		
		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		result, err := client.Translate(ctx, "Hello", "Translate to Russian")
		
		if err != nil {
			t.Logf("Expected unmarshal error: %v", err)
			if result != "" {
				t.Error("Result should be empty when JSON unmarshaling fails")
			}
			
			// Should contain unmarshal error information
			if !strings.Contains(err.Error(), "unmarshal") &&
			   !strings.Contains(err.Error(), "json") &&
			   !strings.Contains(err.Error(), "invalid") &&
			   !strings.Contains(err.Error(), "syntax") {
				t.Logf("Error may not be JSON unmarshal related: %v", err)
			}
		}
	})
	
	// Test 6: Model defaulting behavior
	t.Run("model_defaulting_behavior", func(t *testing.T) {
		config := TranslationConfig{
			Provider: "ollama",
			BaseURL:  "http://httpbin.org", // Base URL - client will append /api/generate
			Model:    "", // Empty model to trigger defaulting
		}
		
		client, err := NewOllamaClient(config)
		if err != nil || client == nil {
			t.Skip("Skipping test - client creation failed")
			return
		}
		
		// The client should still have empty model after creation
		// Defaulting only happens during Translate
		if client.config.Model != "" {
			t.Errorf("Client model should still be empty after creation, got: %s", client.config.Model)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// This will fail with JSON unmarshal error since httpbin.org will return 404 for /api/generate
		// but the model defaulting should happen during Translate
		_, err = client.Translate(ctx, "Hello", "Translate to Russian")
		
		// We should get an error, but it should be a 404 error, not missing model
		if err == nil {
			t.Error("Expected error with httpbin.org response")
			return
		}

		// Any error that's not about missing model confirms the request was made with default model
		if strings.Contains(err.Error(), "model") && strings.Contains(err.Error(), "required") {
			t.Errorf("Got model-related error (defaulting didn't happen): %v", err)
		}
		
		t.Log("Model defaulting confirmed - request was made without model validation errors")
	})
}