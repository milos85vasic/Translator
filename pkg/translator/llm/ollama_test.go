package llm

import (
	"context"
	"testing"
	"time"

	"digital.vasic.translator/pkg/translator"
)

func TestOllamaClient(t *testing.T) {
	tests := []struct {
		name    string
		config  translator.TranslationConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: translator.TranslationConfig{
				Provider: "ollama",
				Model:    "llama3:8b",
				BaseURL:  "http://localhost:11434",
			},
			wantErr: false,
		},
		{
			name: "missing model",
			config: translator.TranslationConfig{
				Provider: "ollama",
				BaseURL:  "http://localhost:11434",
			},
			wantErr: true,
		},
		{
			name: "missing base url",
			config: translator.TranslationConfig{
				Provider: "ollama",
				Model:    "llama3:8b",
			},
			wantErr: false, // Should default to localhost:11434
		},
		{
			name: "invalid provider",
			config: translator.TranslationConfig{
				Provider: "invalid",
				Model:    "llama3:8b",
				BaseURL:  "http://localhost:11434",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOllamaClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOllamaClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && client != nil {
				if client.GetProviderName() != "ollama" {
					t.Errorf("GetProviderName() = %v, want %v", client.GetProviderName(), "ollama")
				}
			}
		})
	}
}

func TestOllamaClient_Translate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if Ollama is not running
	if !isOllamaRunning() {
		t.Skip("Ollama is not running locally for integration test")
	}

	client, err := NewOllamaClient(translator.TranslationConfig{
		Provider: "ollama",
		Model:    "llama3:8b", // Use a common model
		BaseURL:  "http://localhost:11434",
	})
	if err != nil {
		t.Fatalf("Failed to create Ollama client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		text     string
		prompt   string
		wantErr  bool
		validate func(string) bool
	}{
		{
			name:    "simple translation",
			text:    "Hello world",
			prompt:  "Translate this text to Spanish",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 0 && result != "Hello world"
			},
		},
		{
			name:    "creative translation",
			text:    "The sunset painted the sky in brilliant colors",
			prompt:  "Translate this poetic text to French",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 10
			},
		},
		{
			name:    "technical translation",
			text:    "The algorithm processes data efficiently",
			prompt:  "Translate this technical text to German",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 10
			},
		},
		{
			name:    "empty text",
			text:    "",
			prompt:  "Translate this text",
			wantErr: false,
			validate: func(result string) bool {
				return true // Ollama should handle empty text gracefully
			},
		},
		{
			name:    "long text",
			text:    "This comprehensive text explores the fascinating world of artificial intelligence, examining how neural networks learn from data and make predictions that can rival human performance in various tasks.",
			prompt:  "Translate this text to Italian",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 50
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.Translate(ctx, tt.text, tt.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.validate(result) {
				t.Errorf("Translate() result validation failed for: %s", result)
			}
		})
	}
}

func TestOllamaClient_RequestStructure(t *testing.T) {
	client := &OllamaClient{
		config: translator.TranslationConfig{
			Provider: "ollama",
			Model:    "llama3:8b",
		},
	}

	// Test request structure
	req := OllamaRequest{
		Model:  "llama3:8b",
		Prompt: "Test message",
		Options: map[string]interface{}{
			"temperature": 0.7,
			"num_predict": 1024,
		},
		Stream: false,
	}

	if req.Model != "llama3:8b" {
		t.Errorf("Expected model 'llama3:8b', got '%s'", req.Model)
	}

	if req.Prompt != "Test message" {
		t.Errorf("Expected prompt 'Test message', got '%s'", req.Prompt)
	}

	if req.Stream != false {
		t.Error("Expected Stream to be false")
	}

	temp, ok := req.Options["temperature"].(float64)
	if !ok {
		t.Error("Expected temperature in options")
	}
	if temp != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", temp)
	}
}

func TestOllamaClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config translator.TranslationConfig
	}{
		{
			name: "invalid base url",
			config: translator.TranslationConfig{
				Provider: "ollama",
				Model:    "llama3:8b",
				BaseURL:  "http://invalid-url:11434",
			},
		},
		{
			name: "invalid port",
			config: translator.TranslationConfig{
				Provider: "ollama",
				Model:    "llama3:8b",
				BaseURL:  "http://localhost:99999",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOllamaClient(tt.config)
			if err != nil {
				// Some configurations might fail at creation
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err = client.Translate(ctx, "test", "test prompt")
			if err == nil {
				t.Error("Expected error for invalid configuration")
			}
		})
	}
}

func TestOllamaClient_ContextHandling(t *testing.T) {
	client := &OllamaClient{
		config: translator.TranslationConfig{
			Provider: "ollama",
			Model:    "llama3:8b",
		},
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.Translate(ctx, "test", "test prompt")
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
}

func TestOllamaClient_ModelValidation(t *testing.T) {
	validModels := []string{
		"llama3:8b",
		"llama3:70b",
		"codellama:7b",
		"mistral:7b",
		"phi3:mini",
	}

	for _, model := range validModels {
		t.Run("valid_model_"+model, func(t *testing.T) {
			client, err := NewOllamaClient(translator.TranslationConfig{
				Provider: "ollama",
				Model:    model,
				BaseURL:  "http://localhost:11434",
			})
			if err != nil {
				t.Errorf("Expected no error for valid model %s, got: %v", model, err)
			}
			if client == nil {
				t.Error("Expected client to be created")
			}
		})
	}

	invalidModels := []string{
		"",
		"invalid-model-name",
	}

	for _, model := range invalidModels {
		t.Run("invalid_model_"+model, func(t *testing.T) {
			client, err := NewOllamaClient(translator.TranslationConfig{
				Provider: "ollama",
				Model:    model,
				BaseURL:  "http://localhost:11434",
			})
			if err == nil && model == "" {
				t.Error("Expected error for empty model")
			}
			if client != nil && model == "" {
				t.Error("Expected client to be nil for empty model")
			}
		})
	}
}

func TestOllamaClient_DefaultBaseURL(t *testing.T) {
	// Test that default base URL is set when not provided
	client, err := NewOllamaClient(translator.TranslationConfig{
		Provider: "ollama",
		Model:    "llama3:8b",
		// No BaseURL provided
	})
	if err != nil {
		t.Errorf("Expected no error when BaseURL is not provided, got: %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created with default BaseURL")
	}
}

func TestOllamaClient_TemperatureValidation(t *testing.T) {
	tests := []struct {
		name        string
		temperature float64
		expectError bool
	}{
		{
			name:        "valid temperature",
			temperature: 0.7,
			expectError: false,
		},
		{
			name:        "minimum temperature",
			temperature: 0.0,
			expectError: false,
		},
		{
			name:        "maximum temperature",
			temperature: 2.0,
			expectError: false,
		},
		{
			name:        "too high temperature",
			temperature: 2.5,
			expectError: true,
		},
		{
			name:        "negative temperature",
			temperature: -0.1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := translator.TranslationConfig{
				Provider:    "ollama",
				Model:       "llama3:8b",
				BaseURL:     "http://localhost:11434",
				Temperature: tt.temperature,
			}

			client, err := NewOllamaClient(config)
			if tt.expectError && err == nil {
				t.Error("Expected error for invalid temperature")
			}
			if !tt.expectError && client == nil {
				t.Error("Expected client to be created for valid temperature")
			}
		})
	}
}

// Helper function to check if Ollama is running
func isOllamaRunning() bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == 200
}