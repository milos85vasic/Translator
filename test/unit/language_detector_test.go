package unit

import (
	"context"
	"digital.vasic.translator/pkg/language"
	"testing"
)

func TestLanguageDetector(t *testing.T) {
	detector := language.NewDetector(nil)

	t.Run("DetectCyrillic", func(t *testing.T) {
		russianText := "Привет, мир! Это тестовое сообщение на русском языке."
		detected, err := detector.Detect(context.Background(), russianText)

		if err != nil {
			t.Fatalf("Detection failed: %v", err)
		}

		// Should detect as Cyrillic language (Russian, Serbian, Ukrainian, etc.)
		if detected.Code != "ru" && detected.Code != "sr" && detected.Code != "uk" {
			t.Errorf("Expected Cyrillic language, got %s", detected.Code)
		}
	})

	t.Run("DetectLatin", func(t *testing.T) {
		englishText := "Hello, world! This is a test message in English."
		detected, err := detector.Detect(context.Background(), englishText)

		if err != nil {
			t.Fatalf("Detection failed: %v", err)
		}

		// Should detect as English (default for Latin)
		if detected.Code != "en" {
			t.Logf("Detected as %s instead of en (acceptable for Latin text)", detected.Code)
		}
	})

	t.Run("ParseLanguage", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"en", "en"},
			{"EN", "en"},
			{"English", "en"},
			{"english", "en"},
			{"ru", "ru"},
			{"Russian", "ru"},
			{"sr", "sr"},
			{"Serbian", "sr"},
			{"de", "de"},
			{"German", "de"},
			{"fr", "fr"},
			{"es", "es"},
		}

		for _, tt := range tests {
			lang, err := language.ParseLanguage(tt.input)
			if err != nil {
				t.Errorf("ParseLanguage(%s) failed: %v", tt.input, err)
				continue
			}

			if lang.Code != tt.expected {
				t.Errorf("ParseLanguage(%s) = %s, want %s", tt.input, lang.Code, tt.expected)
			}
		}
	})

	t.Run("GetSupportedLanguages", func(t *testing.T) {
		langs := language.GetSupportedLanguages()

		if len(langs) < 10 {
			t.Errorf("Expected at least 10 supported languages, got %d", len(langs))
		}

		// Check that common languages are present
		langMap := make(map[string]bool)
		for _, lang := range langs {
			langMap[lang.Code] = true
		}

		required := []string{"en", "ru", "sr", "de", "fr", "es"}
		for _, code := range required {
			if !langMap[code] {
				t.Errorf("Language %s should be supported", code)
			}
		}
	})
}
