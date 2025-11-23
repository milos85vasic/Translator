package llm

import (
	"context"
	"testing"
	"time"

	"digital.vasic.translator/pkg/translator"
)

func TestGeminiClient(t *testing.T) {
	tests := []struct {
		name    string
		config  translator.TranslationConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: translator.TranslationConfig{
				Provider: "gemini",
				APIKey:   "test-key",
				Model:    "gemini-pro",
				BaseURL:  "https://generativelanguage.googleapis.com",
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: translator.TranslationConfig{
				Provider: "gemini",
				Model:    "gemini-pro",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewGeminiClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGeminiClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if client.GetProviderName() != "gemini" {
					t.Errorf("GetProviderName() = %v, want %v", client.GetProviderName(), "gemini")
				}
			}
		})
	}
}

func TestGeminiClient_Translate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if no real API key
	apiKey := getTestAPIKey("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("No GEMINI_API_KEY provided for integration test")
	}

	client, err := NewGeminiClient(translator.TranslationConfig{
		Provider: "gemini",
		APIKey:   apiKey,
		Model:    "gemini-pro",
		BaseURL:  "https://generativelanguage.googleapis.com",
	})
	if err != nil {
		t.Fatalf("Failed to create Gemini client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
			prompt:  "Translate this text to Russian",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 0 && result != "Hello world"
			},
		},
		{
			name:    "empty text",
			text:    "",
			prompt:  "Translate this text",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) == 0
			},
		},
		{
			name:    "long text",
			text:    "This is a very long text that should test the translation capabilities of the Gemini model when handling larger inputs.",
			prompt:  "Translate this text to French",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 10
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

func TestGeminiClient_RequestStructure(t *testing.T) {
	client := &GeminiClient{
		config: translator.TranslationConfig{
			Provider: "gemini",
			Model:    "gemini-pro",
		},
	}

	// Test request structure
	req := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: "Test content",
					},
				},
			},
		},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature: 0.7,
			MaxTokens:   1024,
		},
	}

	if len(req.Contents) == 0 {
		t.Error("Contents should not be empty")
	}

	if req.GenerationConfig.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", req.GenerationConfig.Temperature)
	}

	if req.GenerationConfig.MaxTokens != 1024 {
		t.Errorf("Expected MaxTokens 1024, got %d", req.GenerationConfig.MaxTokens)
	}
}

func TestGeminiClient_ErrorHandling(t *testing.T) {
	client := &GeminiClient{
		config: translator.TranslationConfig{
			Provider: "gemini",
			APIKey:   "invalid-key",
			Model:    "gemini-pro",
			BaseURL:  "https://invalid-url.com",
		},
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.Translate(ctx, "test", "test prompt")
	if err == nil {
		t.Error("Expected error for invalid configuration")
	}
}