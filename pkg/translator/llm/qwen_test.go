package llm

import (
	"context"
	"testing"
	"time"

	"digital.vasic.translator/pkg/translator"
)

func TestQwenClient(t *testing.T) {
	tests := []struct {
		name    string
		config  translator.TranslationConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: translator.TranslationConfig{
				Provider: "qwen",
				APIKey:   "test-key",
				Model:    "qwen-plus",
				BaseURL:  "https://dashscope.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: translator.TranslationConfig{
				Provider: "qwen",
				Model:    "qwen-plus",
			},
			wantErr: true,
		},
		{
			name: "invalid provider",
			config: translator.TranslationConfig{
				Provider: "invalid",
				APIKey:   "test-key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewQwenClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQwenClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && client != nil {
				if client.GetProviderName() != "qwen" {
					t.Errorf("GetProviderName() = %v, want %v", client.GetProviderName(), "qwen")
				}
			}
		})
	}
}

func TestQwenClient_Translate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if no real API key
	apiKey := getTestAPIKey("QWEN_API_KEY")
	if apiKey == "" {
		t.Skip("No QWEN_API_KEY provided for integration test")
	}

	client, err := NewQwenClient(translator.TranslationConfig{
		Provider: "qwen",
		APIKey:   apiKey,
		Model:    "qwen-plus",
		BaseURL:  "https://dashscope.aliyuncs.com",
	})
	if err != nil {
		t.Fatalf("Failed to create Qwen client: %v", err)
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
				return len(result) > 0 && containsChinese(result)
			},
		},
		{
			name:    "english to french",
			text:    "How are you?",
			prompt:  "Translate this text to French",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 0 && containsFrench(result)
			},
		},
		{
			name:    "empty text",
			text:    "",
			prompt:  "Translate this text",
			wantErr: false,
			validate: func(result string) bool {
				return true // API might handle empty text gracefully
			},
		},
		{
			name:    "long text",
			text:    "This is a comprehensive test of the Qwen translation capabilities with a longer text segment that includes multiple sentences and complex structures.",
			prompt:  "Translate this text to Japanese",
			wantErr: false,
			validate: func(result string) bool {
				return len(result) > 20
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

func TestQwenClient_ContextHandling(t *testing.T) {
	client := &QwenClient{
		config: translator.TranslationConfig{
			Provider: "qwen",
			Model:    "qwen-plus",
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

func TestQwenClient_ConfigurationValidation(t *testing.T) {
	tests := []struct {
		name   string
		config translator.TranslationConfig
		errMsg string
	}{
		{
			name: "missing model",
			config: translator.TranslationConfig{
				Provider: "qwen",
				APIKey:   "test-key",
			},
			errMsg: "model",
		},
		{
			name: "invalid base url",
			config: translator.TranslationConfig{
				Provider: "qwen",
				APIKey:   "test-key",
				Model:    "qwen-plus",
				BaseURL:  "invalid-url",
			},
			errMsg: "base URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQwenClient(tt.config)
			if err == nil {
				t.Error("Expected error for invalid config")
				return
			}
			if !contains(err.Error(), tt.errMsg) {
				t.Errorf("Error message should contain %s, got: %s", tt.errMsg, err.Error())
			}
		})
	}
}

// Helper functions
func containsChinese(s string) bool {
	for _, r := range s {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

func containsFrench(s string) bool {
	frenchChars := []rune{'é', 'è', 'à', 'ç', 'ê', 'î', 'ô', 'ù'}
	for _, char := range frenchChars {
		for _, r := range s {
			if r == char {
				return true
			}
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
			 s[len(s)-len(substr):] == substr ||
			 indexOf(s, substr) >= 0)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}