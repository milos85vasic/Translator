package language

import (
	"context"
	"fmt"
	"strings"
	"unicode"
)

// Language represents a language with its codes
type Language struct {
	Code   string // ISO 639-1 code (e.g., "en", "ru", "sr")
	Name   string // English name (e.g., "English", "Russian")
	Native string // Native name (e.g., "English", "Русский")
}

// Common languages
var (
	English    = Language{Code: "en", Name: "English", Native: "English"}
	Russian    = Language{Code: "ru", Name: "Russian", Native: "Русский"}
	Serbian    = Language{Code: "sr", Name: "Serbian", Native: "Српски"}
	German     = Language{Code: "de", Name: "German", Native: "Deutsch"}
	French     = Language{Code: "fr", Name: "French", Native: "Français"}
	Spanish    = Language{Code: "es", Name: "Spanish", Native: "Español"}
	Italian    = Language{Code: "it", Name: "Italian", Native: "Italiano"}
	Portuguese = Language{Code: "pt", Name: "Portuguese", Native: "Português"}
	Chinese    = Language{Code: "zh", Name: "Chinese", Native: "中文"}
	Japanese   = Language{Code: "ja", Name: "Japanese", Native: "日本語"}
	Korean     = Language{Code: "ko", Name: "Korean", Native: "한국어"}
	Arabic     = Language{Code: "ar", Name: "Arabic", Native: "العربية"}
	Polish     = Language{Code: "pl", Name: "Polish", Native: "Polski"}
	Ukrainian  = Language{Code: "uk", Name: "Ukrainian", Native: "Українська"}
	Czech      = Language{Code: "cs", Name: "Czech", Native: "Čeština"}
	Slovak     = Language{Code: "sk", Name: "Slovak", Native: "Slovenčina"}
	Croatian   = Language{Code: "hr", Name: "Croatian", Native: "Hrvatski"}
	Bulgarian  = Language{Code: "bg", Name: "Bulgarian", Native: "Български"}
)

// languageMap maps codes and names to Language structs
var languageMap = map[string]Language{
	// Codes
	"en": English, "eng": English,
	"ru": Russian, "rus": Russian,
	"sr": Serbian, "srp": Serbian,
	"de": German, "deu": German, "ger": German,
	"fr": French, "fra": French, "fre": French,
	"es": Spanish, "spa": Spanish,
	"it": Italian, "ita": Italian,
	"pt": Portuguese, "por": Portuguese,
	"zh": Chinese, "zho": Chinese, "chi": Chinese,
	"ja": Japanese, "jpn": Japanese,
	"ko": Korean, "kor": Korean,
	"ar": Arabic, "ara": Arabic,
	"pl": Polish, "pol": Polish,
	"uk": Ukrainian, "ukr": Ukrainian,
	"cs": Czech, "ces": Czech, "cze": Czech,
	"sk": Slovak, "slk": Slovak, "slo": Slovak,
	"hr": Croatian, "hrv": Croatian,
	"bg": Bulgarian, "bul": Bulgarian,

	// Names (lowercase)
	"english":    English,
	"russian":    Russian,
	"serbian":    Serbian,
	"german":     German,
	"french":     French,
	"spanish":    Spanish,
	"italian":    Italian,
	"portuguese": Portuguese,
	"chinese":    Chinese,
	"japanese":   Japanese,
	"korean":     Korean,
	"arabic":     Arabic,
	"polish":     Polish,
	"ukrainian":  Ukrainian,
	"czech":      Czech,
	"slovak":     Slovak,
	"croatian":   Croatian,
	"bulgarian":  Bulgarian,
}

// Detector handles language detection
type Detector struct {
	llmDetector LLMDetector
}

// LLMDetector interface for LLM-based language detection
type LLMDetector interface {
	DetectLanguage(ctx context.Context, text string) (string, error)
}

// NewDetector creates a new language detector
func NewDetector(llmDetector LLMDetector) *Detector {
	return &Detector{
		llmDetector: llmDetector,
	}
}

// Detect detects the language of the given text
func (d *Detector) Detect(ctx context.Context, text string) (Language, error) {
	// Try LLM detection first if available
	if d.llmDetector != nil {
		code, err := d.llmDetector.DetectLanguage(ctx, text)
		if err == nil && code != "" {
			if lang, ok := languageMap[strings.ToLower(code)]; ok {
				return lang, nil
			}
		}
	}

	// Fallback to heuristic detection
	return d.detectHeuristic(text), nil
}

// detectHeuristic uses character-based heuristics to detect language
func (d *Detector) detectHeuristic(text string) Language {
	if text == "" {
		return English // default
	}

	// Sample first 1000 characters
	sample := text
	if len(text) > 1000 {
		sample = text[:1000]
	}

	// Count character types
	var (
		cyrillic int
		latin    int
		cjk      int
		arabic   int
	)

	for _, r := range sample {
		switch {
		case isCyrillic(r):
			cyrillic++
		case isLatin(r):
			latin++
		case isCJK(r):
			cjk++
		case isArabic(r):
			arabic++
		}
	}

	// Determine language by character frequency
	total := cyrillic + latin + cjk + arabic
	if total == 0 {
		return English // default
	}

	// CJK languages
	if float64(cjk)/float64(total) > 0.3 {
		// Could be Chinese, Japanese, or Korean
		// For now, default to Chinese
		return Chinese
	}

	// Arabic
	if float64(arabic)/float64(total) > 0.3 {
		return Arabic
	}

	// Cyrillic scripts
	if float64(cyrillic)/float64(total) > 0.3 {
		// Try to distinguish between Russian, Serbian, Ukrainian, etc.
		return d.detectCyrillicLanguage(sample)
	}

	// Latin scripts - default to English
	// Could be improved with n-gram analysis
	return English
}

// detectCyrillicLanguage distinguishes between Cyrillic languages
func (d *Detector) detectCyrillicLanguage(text string) Language {
	// Count language-specific characters
	var (
		russianChars  int
		serbianChars  int
		ukrainianChars int
	)

	for _, r := range text {
		switch r {
		case 'Ё', 'ё', 'Ы', 'ы', 'Э', 'э':
			russianChars++
		case 'Ђ', 'ђ', 'Ћ', 'ћ', 'Љ', 'љ', 'Њ', 'њ', 'Џ', 'џ':
			serbianChars++
		case 'Є', 'є', 'І', 'і', 'Ї', 'ї', 'Ґ', 'ґ':
			ukrainianChars++
		}
	}

	// Return language with most specific characters
	if serbianChars > russianChars && serbianChars > ukrainianChars {
		return Serbian
	}
	if ukrainianChars > russianChars && ukrainianChars > serbianChars {
		return Ukrainian
	}

	// Default to Russian for Cyrillic
	return Russian
}

// ParseLanguage parses a language string (code or name)
func ParseLanguage(s string) (Language, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if lang, ok := languageMap[s]; ok {
		return lang, nil
	}
	return Language{}, fmt.Errorf("unknown language: %s", s)
}

// GetSupportedLanguages returns list of supported languages
func GetSupportedLanguages() []Language {
	return []Language{
		English, Russian, Serbian, German, French, Spanish,
		Italian, Portuguese, Chinese, Japanese, Korean, Arabic,
		Polish, Ukrainian, Czech, Slovak, Croatian, Bulgarian,
	}
}

// isCyrillic checks if a rune is Cyrillic
func isCyrillic(r rune) bool {
	return unicode.Is(unicode.Cyrillic, r)
}

// isLatin checks if a rune is Latin
func isLatin(r rune) bool {
	return unicode.Is(unicode.Latin, r)
}

// isCJK checks if a rune is CJK (Chinese, Japanese, Korean)
func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r)
}

// isArabic checks if a rune is Arabic
func isArabic(r rune) bool {
	return unicode.Is(unicode.Arabic, r)
}
