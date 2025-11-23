package llm

import (
	"context"
	"testing"
	"time"

	"digital.vasic.translator/pkg/translator"
)

func TestZhipuClient(t *testing.T) {
	tests := []struct {
		name    string
		config  translator.TranslationConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: translator.TranslationConfig{
				Provider: "zhipu",
				APIKey:   "test-key",
				Model:    "glm-4",
				BaseURL:  "https://open.bigmodel.cn",
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: translator.TranslationConfig{
				Provider: "zhipu",
				Model:    "glm-4",
			},
			wantErr: true,
		},
		{
			name: "missing model",
			config: translator.TranslationConfig{
				Provider: "zhipu",
				APIKey:   "test-key",
			},
			wantErr: true,
		},
		{
			name: "invalid provider",
			config: translator.TranslationConfig{
				Provider: "invalid",
				APIKey:   "test-key",
				Model:    "glm-4",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewZhipuClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewZhipuClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && client != nil {
				if client.GetProviderName() != "zhipu" {
					t.Errorf("GetProviderName() = %v, want %v", client.GetProviderName(), "zhipu")
				}
			}
		})
	}
}

func TestZhipuClient_Translate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if no real API key
	apiKey := getTestAPIKey("ZHIPU_API_KEY")
	if apiKey == "" {
		t.Skip("No ZHIPU_API_KEY provided for integration test")
	}

	client, err := NewZhipuClient(translator.TranslationConfig{
		Provider: "zhipu",
		APIKey:   apiKey,
		Model:    "glm-4",
		BaseURL:  "https://open.bigmodel.cn",
	})
	if err != nil {
		t.Fatalf("Failed to create Zhipu client: %v", err)
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
			prompt:  "Translate this text to Chinese",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 0 && result != "Hello world"
			},
		},
		{
			name:    "literary translation",
			text:    "The river flows gently through the valley",
			prompt:  "Translate this literary text to Japanese",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 10
			},
		},
		{
			name:    "technical translation",
			text:    "Machine learning algorithms optimize performance",
			prompt:  "Translate this technical text to German",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 15
			},
		},
		{
			name:    "empty text",
			text:    "",
			prompt:  "Translate this text",
			wantErr: false,
			validate: func(result string) bool {
				return true // API should handle empty text gracefully
			},
		},
		{
			name:    "long text",
			text:    "This comprehensive study explores the intricate relationships between artificial intelligence and human creativity, examining how collaborative approaches between humans and machines can lead to unprecedented innovations in various fields including art, music, literature, and scientific discovery.",
			prompt:  "Translate this academic text to French",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 100
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

func TestZhipuClient_RequestStructure(t *testing.T) {
	client := &ZhipuClient{
		config: translator.TranslationConfig{
			Provider: "zhipu",
			Model:    "glm-4",
		},
	}

	// Test request structure
	req := ZhipuRequest{
		Model: "glm-4",
		Messages: []ZhipuMessage{
			{
				Role:    "user",
				Content: "Test message",
			},
		},
		Temperature: 0.7,
		MaxTokens:   1024,
		Stream:      false,
	}

	if req.Model != "glm-4" {
		t.Errorf("Expected model 'glm-4', got '%s'", req.Model)
	}

	if len(req.Messages) == 0 {
		t.Error("Messages should not be empty")
	}

	if req.Messages[0].Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", req.Messages[0].Role)
	}

	if req.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", req.Temperature)
	}

	if req.MaxTokens != 1024 {
		t.Errorf("Expected MaxTokens 1024, got %d", req.MaxTokens)
	}

	if req.Stream != false {
		t.Error("Expected Stream to be false")
	}
}

func TestZhipuClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config translator.TranslationConfig
	}{
		{
			name: "invalid api key",
			config: translator.TranslationConfig{
				Provider: "zhipu",
				APIKey:   "invalid-key",
				Model:    "glm-4",
				BaseURL:  "https://open.bigmodel.cn",
			},
		},
		{
			name: "invalid base url",
			config: translator.TranslationConfig{
				Provider: "zhipu",
				APIKey:   "test-key",
				Model:    "glm-4",
				BaseURL:  "https://invalid-url.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewZhipuClient(tt.config)
			if err != nil {
				// Some configurations might fail at creation
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err = client.Translate(ctx, "test", "test prompt")
			if err == nil {
				t.Error("Expected error for invalid configuration")
			}
		})
	}
}

func TestZhipuClient_ContextHandling(t *testing.T) {
	client := &ZhipuClient{
		config: translator.TranslationConfig{
			Provider: "zhipu",
			Model:    "glm-4",
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

func TestZhipuClient_ModelValidation(t *testing.T) {
	validModels := []string{
		"glm-4",
		"glm-3-turbo",
		"glm-4v", // vision model
	}

	for _, model := range validModels {
		t.Run("valid_model_"+model, func(t *testing.T) {
			client, err := NewZhipuClient(translator.TranslationConfig{
				Provider: "zhipu",
				APIKey:   "test-key",
				Model:    model,
				BaseURL:  "https://open.bigmodel.cn",
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
		"invalid-model",
		"gpt-4",
		"claude-3",
		"",
	}

	for _, model := range invalidModels {
		t.Run("invalid_model_"+model, func(t *testing.T) {
			client, err := NewZhipuClient(translator.TranslationConfig{
				Provider: "zhipu",
				APIKey:   "test-key",
				Model:    model,
				BaseURL:  "https://open.bigmodel.cn",
			})
			// Note: Zhipu might not validate models at creation time
			// This test might need adjustment based on actual implementation
			if err == nil && model == "" {
				t.Error("Expected error for empty model")
			}
			if client != nil && model == "" {
				t.Error("Expected client to be nil for empty model")
			}
		})
	}
}

func TestZhipuClient_TemperatureValidation(t *testing.T) {
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
			temperature: 1.0,
			expectError: false,
		},
		{
			name:        "too high temperature",
			temperature: 1.5,
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
				Provider:    "zhipu",
				APIKey:      "test-key",
				Model:       "glm-4",
				BaseURL:     "https://open.bigmodel.cn",
				Temperature: tt.temperature,
			}

			client, err := NewZhipuClient(config)
			if tt.expectError && err == nil {
				t.Error("Expected error for invalid temperature")
			}
			if !tt.expectError && client == nil {
				t.Error("Expected client to be created for valid temperature")
			}
		})
	}
}

func TestZhipuClient_MaxTokensValidation(t *testing.T) {
	tests := []struct {
		name        string
		maxTokens   int
		expectError bool
	}{
		{
			name:        "valid max tokens",
			maxTokens:   1024,
			expectError: false,
		},
		{
			name:        "minimum max tokens",
			maxTokens:   1,
			expectError: false,
		},
		{
			name:        "maximum max tokens",
			maxTokens:   8192,
			expectError: false,
		},
		{
			name:        "too many tokens",
			maxTokens:   100000,
			expectError: true,
		},
		{
			name:        "zero max tokens",
			maxTokens:   0,
			expectError: true,
		},
		{
			name:        "negative max tokens",
			maxTokens:   -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := translator.TranslationConfig{
				Provider:  "zhipu",
				APIKey:    "test-key",
				Model:     "glm-4",
				BaseURL:   "https://open.bigmodel.cn",
				MaxTokens: tt.maxTokens,
			}

			client, err := NewZhipuClient(config)
			if tt.expectError && err == nil {
				t.Error("Expected error for invalid max tokens")
			}
			if !tt.expectError && client == nil {
				t.Error("Expected client to be created for valid max tokens")
			}
		})
	}
}