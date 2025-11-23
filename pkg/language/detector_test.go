package language

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockLLMDetector for testing
type MockLLMDetector struct {
	detectFunc func(ctx context.Context, text string) (string, error)
}

func (m *MockLLMDetector) DetectLanguage(ctx context.Context, text string) (string, error) {
	if m.detectFunc != nil {
		return m.detectFunc(ctx, text)
	}
	return "en", nil
}

func TestNewDetector(t *testing.T) {
	mockLLM := &MockLLMDetector{}
	detector := NewDetector(mockLLM)

	assert.NotNil(t, detector)
	assert.Equal(t, mockLLM, detector.llmDetector)
}

func TestNewDetector_NilLLM(t *testing.T) {
	detector := NewDetector(nil)

	assert.NotNil(t, detector)
	assert.Nil(t, detector.llmDetector)
}

func TestDetector_Detect_WithLLM(t *testing.T) {
	tests := []struct {
		name        string
		llmResponse string
		llmError    error
		expected    Language
		expectError bool
	}{
		{
			name:        "LLM returns valid language",
			llmResponse: "es",
			llmError:    nil,
			expected:    Spanish,
			expectError: false,
		},
		{
			name:        "LLM returns invalid language",
			llmResponse: "invalid",
			llmError:    nil,
			expected:    English, // Should fallback to heuristic
			expectError: false,
		},
		{
			name:        "LLM returns error",
			llmResponse: "",
			llmError:    assert.AnError,
			expected:    English, // Should fallback to heuristic
			expectError: false,
		},
		{
			name:        "LLM returns empty response",
			llmResponse: "",
			llmError:    nil,
			expected:    English, // Should fallback to heuristic
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLLM := &MockLLMDetector{
				detectFunc: func(ctx context.Context, text string) (string, error) {
					return tt.llmResponse, tt.llmError
				},
			}
			detector := NewDetector(mockLLM)

			result, err := detector.Detect(context.Background(), "test text")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetector_Detect_WithoutLLM(t *testing.T) {
	detector := NewDetector(nil)

	result, err := detector.Detect(context.Background(), "test text")

	assert.NoError(t, err)
	// Should use heuristic detection
	assert.NotEmpty(t, result.Code)
}

func TestDetector_Detect_EmptyText(t *testing.T) {
	detector := NewDetector(nil)

	result, err := detector.Detect(context.Background(), "")

	assert.NoError(t, err)
	assert.Equal(t, English, result) // Default to English
}

func TestDetector_Detect_Heuristic_Cyrillic(t *testing.T) {
	detector := NewDetector(nil)

	tests := []struct {
		name     string
		text     string
		expected Language
	}{
		{
			name:     "Russian text",
			text:     "Привет мир",
			expected: Russian,
		},
		{
			name:     "Serbian Cyrillic text",
			text:     "Здраво свет",
			expected: Serbian,
		},
		{
			name:     "Bulgarian text",
			text:     "Здравей свят",
			expected: Bulgarian,
		},
		{
			name:     "Ukrainian text",
			text:     "Привіт світ",
			expected: Ukrainian,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(context.Background(), tt.text)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetector_Detect_Heuristic_Latin(t *testing.T) {
	detector := NewDetector(nil)

	tests := []struct {
		name     string
		text     string
		expected Language
	}{
		{
			name:     "English text",
			text:     "Hello world",
			expected: English,
		},
		{
			name:     "Spanish text",
			text:     "Hola mundo",
			expected: Spanish,
		},
		{
			name:     "French text",
			text:     "Bonjour le monde",
			expected: French,
		},
		{
			name:     "German text",
			text:     "Hallo Welt",
			expected: German,
		},
		{
			name:     "Italian text",
			text:     "Ciao mondo",
			expected: Italian,
		},
		{
			name:     "Portuguese text",
			text:     "Olá mundo",
			expected: Portuguese,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(context.Background(), tt.text)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetector_Detect_Heuristic_SpecificPatterns(t *testing.T) {
	detector := NewDetector(nil)

	tests := []struct {
		name     string
		text     string
		expected Language
	}{
		{
			name:     "Polish text with specific characters",
			text:     "Witaj świecie",
			expected: Polish,
		},
		{
			name:     "Czech text with háček",
			text:     "Dobrý den",
			expected: Czech,
		},
		{
			name:     "Slovak text",
			text:     "Dobrý deň",
			expected: Slovak,
		},
		{
			name:     "Croatian text",
			text:     "Pozdrav svijetu",
			expected: Croatian,
		},
		{
			name:     "Chinese text",
			text:     "你好世界",
			expected: Chinese,
		},
		{
			name:     "Japanese text",
			text:     "こんにちは世界",
			expected: Japanese,
		},
		{
			name:     "Korean text",
			text:     "안녕하세요 세계",
			expected: Korean,
		},
		{
			name:     "Arabic text",
			text:     "مرحبا بالعالم",
			expected: Arabic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(context.Background(), tt.text)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetector_Detect_Heuristic_MixedText(t *testing.T) {
	detector := NewDetector(nil)

	tests := []struct {
		name     string
		text     string
		expected Language
	}{
		{
			name:     "Mostly Cyrillic with some Latin",
			text:     "Привет world, как дела?",
			expected: Russian,
		},
		{
			name:     "Mostly Latin with some Cyrillic",
			text:     "Hello мир, how are you?",
			expected: English,
		},
		{
			name:     "Equal mix - should default to Latin",
			text:     "Привет Hello",
			expected: English,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(context.Background(), tt.text)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetector_Detect_Heuristic_ShortText(t *testing.T) {
	detector := NewDetector(nil)

	tests := []struct {
		name     string
		text     string
		expected Language
	}{
		{
			name:     "Single Cyrillic character",
			text:     "п",
			expected: Russian,
		},
		{
			name:     "Single Latin character",
			text:     "a",
			expected: English,
		},
		{
			name:     "Two characters - Cyrillic",
			text:     "пр",
			expected: Russian,
		},
		{
			name:     "Two characters - Latin",
			text:     "ab",
			expected: English,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(context.Background(), tt.text)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetector_Detect_Heuristic_LongText(t *testing.T) {
	detector := NewDetector(nil)

	// Test with text longer than 1000 characters
	longText := "Это очень длинный текст на русском языке, который используется для проверки того, " +
		"как детектор работает с большими объемами текста. Текст должен быть достаточно длинным, " +
		"чтобы превысить лимит в 1000 символов и проверить, что детектор правильно обрезает " +
		"текст перед анализом. Здесь мы добавляем еще больше символов, чтобы убедиться, " +
		"что текст действительно длинный и содержит достаточно кириллических символов " +
		"для точного определения языка."

	result, err := detector.Detect(context.Background(), longText)

	assert.NoError(t, err)
	assert.Equal(t, Russian, result)
}

func TestDetector_Detect_Heuristic_SpecialCharacters(t *testing.T) {
	detector := NewDetector(nil)

	tests := []struct {
		name     string
		text     string
		expected Language
	}{
		{
			name:     "Text with numbers and punctuation",
			text:     "Hello, world! 123 testing...",
			expected: English,
		},
		{
			name:     "Text with emojis",
			text:     "Hello world! 😊 🌍",
			expected: English,
		},
		{
			name:     "Text with mixed scripts",
			text:     "Привет! Hello! 123",
			expected: Russian, // Cyrillic should dominate
		},
		{
			name:     "Only punctuation and numbers",
			text:     "123!@#$%^&*()",
			expected: English, // Default to Latin
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.Detect(context.Background(), tt.text)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLanguageMap_Coverage(t *testing.T) {
	// Test that all common languages are in the map
	expectedLanguages := []Language{
		English, Russian, Serbian, German, French, Spanish, Italian,
		Portuguese, Chinese, Japanese, Korean, Arabic, Polish,
		Ukrainian, Czech, Slovak, Croatian, Bulgarian,
	}

	for _, lang := range expectedLanguages {
		// Test lowercase code
		assert.Contains(t, languageMap, lang.Code, "Language code %s should be in map", lang.Code)

		// Test lowercase full name
		assert.Contains(t, languageMap, lang.Name, "Language name %s should be in map", lang.Name)
	}
}

func TestLanguageMap_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected Language
	}{
		{"en", English},
		{"EN", English},
		{"En", English},
		{"english", English},
		{"ENGLISH", English},
		{"English", English},
		{"ru", Russian},
		{"RU", Russian},
		{"russian", Russian},
		{"RUSSIAN", Russian},
		{"Russian", Russian},
		{"sr", Serbian},
		{"SR", Serbian},
		{"serbian", Serbian},
		{"SERBIAN", Serbian},
		{"Serbian", Serbian},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, ok := languageMap[strings.ToLower(tt.input)]
			assert.True(t, ok, "Input %s should be found in language map", tt.input)
			assert.Equal(t, tt.expected, result, "Input %s should map to correct language", tt.input)
		})
	}
}

func TestDetector_ContextCancellation(t *testing.T) {
	detector := NewDetector(nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := detector.Detect(ctx, "test text")

	// Should still work since heuristic detection doesn't check context
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Code)
}

func TestDetector_ConcurrentAccess(t *testing.T) {
	detector := NewDetector(nil)

	// Test concurrent detection
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := detector.Detect(context.Background(), "test text")
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
		case <-time.After(time.Second * 5):
			t.Fatal("Timeout waiting for concurrent detections")
		}
	}
}

func BenchmarkDetector_Detect_Heuristic(b *testing.B) {
	detector := NewDetector(nil)
	text := "This is a test text for benchmarking language detection performance."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.Detect(context.Background(), text)
	}
}

func BenchmarkDetector_Detect_WithLLM(b *testing.B) {
	mockLLM := &MockLLMDetector{
		detectFunc: func(ctx context.Context, text string) (string, error) {
			return "en", nil
		},
	}
	detector := NewDetector(mockLLM)
	text := "This is a test text for benchmarking language detection performance."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.Detect(context.Background(), text)
	}
}
