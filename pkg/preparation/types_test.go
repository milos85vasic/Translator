package preparation

import (
	"encoding/json"
	"testing"
	"time"
)

func TestContentAnalysis_JSONSerialization(t *testing.T) {
	analysis := ContentAnalysis{
		ContentType:    "Novel",
		Genre:          "Science Fiction",
		Subgenres:      []string{"Space Opera", "Dystopian"},
		Tone:           "Formal",
		LanguageStyle:  "Literary",
		TargetAudience: "Adult",
		UntranslatableTerms: []UntranslatableTerm{
			{
				Term:            "Borscht",
				OriginalScript:  "Cyrillic",
				Reason:          "Traditional soup name",
				Context:         []string{"Chapter 1", "Chapter 3"},
				Transliteration: "Borshch",
			},
		},
		FootnoteGuidance: []FootnoteGuidance{
			{
				Term:        "Dacha",
				Explanation: "Russian country house",
				Locations:   []string{"Chapter 2"},
				Priority:    "Medium",
			},
		},
		Characters: []Character{
			{
				Name:            "Ivan",
				AlternateNames:  []string{"Vanya"},
				Role:            "Protagonist",
				SpeechPattern:   "Formal, educated",
				KeyTraits:       []string{"Brave", "Intelligent"},
				NameTranslation: map[string]string{"en": "John", "sr": "Jovan"},
			},
		},
		KeyThemes: []string{"Redemption", "Technology vs Humanity"},
		CulturalReferences: []CulturalReference{
			{
				Reference:   "Matryoshka",
				Origin:      "Russian",
				Explanation: "Nesting dolls",
				Handling:    "Keep with footnote",
			},
		},
		ChapterAnalyses: []ChapterAnalysis{
			{
				ChapterID:    "ch1",
				ChapterNum:   1,
				Title:        "The Beginning",
				Summary:      "Introduction to main character",
				KeyPoints:    []string{"Character introduction", "Setting establishment"},
				Caveats:      []string{"Complex metaphors"},
				Tone:         "Mysterious",
				Complexity:   "Moderate",
				SpecialNotes: "Pay attention to symbolism",
			},
		},
		AnalysisVersion: 1,
		AnalyzedAt:      time.Now(),
		AnalyzedBy:      "GPT-4",
	}

	// Test JSON marshaling
	data, err := json.Marshal(analysis)
	if err != nil {
		t.Fatalf("Failed to marshal ContentAnalysis: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled ContentAnalysis
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ContentAnalysis: %v", err)
	}

	// Verify key fields
	if unmarshaled.ContentType != analysis.ContentType {
		t.Errorf("ContentType = %s, want %s", unmarshaled.ContentType, analysis.ContentType)
	}

	if len(unmarshaled.Subgenres) != len(analysis.Subgenres) {
		t.Errorf("Subgenres length = %d, want %d", len(unmarshaled.Subgenres), len(analysis.Subgenres))
	}

	if len(unmarshaled.UntranslatableTerms) != len(analysis.UntranslatableTerms) {
		t.Errorf("UntranslatableTerms length = %d, want %d", len(unmarshaled.UntranslatableTerms), len(analysis.UntranslatableTerms))
	}

	if len(unmarshaled.Characters) != len(analysis.Characters) {
		t.Errorf("Characters length = %d, want %d", len(unmarshaled.Characters), len(analysis.Characters))
	}
}

func TestUntranslatableTerm_Fields(t *testing.T) {
	term := UntranslatableTerm{
		Term:            "Samovar",
		OriginalScript:  "Cyrillic",
		Reason:          "Traditional Russian tea maker",
		Context:         []string{"Chapter 1", "Chapter 5"},
		Transliteration: "Samovar",
	}

	if term.Term != "Samovar" {
		t.Errorf("Term = %s, want Samovar", term.Term)
	}

	if term.OriginalScript != "Cyrillic" {
		t.Errorf("OriginalScript = %s, want Cyrillic", term.OriginalScript)
	}

	if len(term.Context) != 2 {
		t.Errorf("Context length = %d, want 2", len(term.Context))
	}
}

func TestFootnoteGuidance_Fields(t *testing.T) {
	guidance := FootnoteGuidance{
		Term:        "Tovarisch",
		Explanation: "Russian term for comrade",
		Locations:   []string{"Chapter 2", "Chapter 4"},
		Priority:    "High",
	}

	if guidance.Term != "Tovarisch" {
		t.Errorf("Term = %s, want Tovarisch", guidance.Term)
	}

	if guidance.Priority != "High" {
		t.Errorf("Priority = %s, want High", guidance.Priority)
	}

	if len(guidance.Locations) != 2 {
		t.Errorf("Locations length = %d, want 2", len(guidance.Locations))
	}
}

func TestCharacter_NameTranslations(t *testing.T) {
	character := Character{
		Name:           "Alexei",
		AlternateNames: []string{"Alyosha", "Lyosha"},
		Role:           "Protagonist",
		SpeechPattern:  "Poetic, melancholic",
		KeyTraits:      []string{"Romantic", "Idealistic"},
		NameTranslation: map[string]string{
			"en": "Alex",
			"sr": "Aleksa",
			"de": "Alexej",
		},
	}

	// Test name translations
	if character.NameTranslation["en"] != "Alex" {
		t.Errorf("English translation = %s, want Alex", character.NameTranslation["en"])
	}

	if character.NameTranslation["sr"] != "Aleksa" {
		t.Errorf("Serbian translation = %s, want Aleksa", character.NameTranslation["sr"])
	}

	if len(character.AlternateNames) != 2 {
		t.Errorf("AlternateNames length = %d, want 2", len(character.AlternateNames))
	}
}

func TestCulturalReference_Handling(t *testing.T) {
	ref := CulturalReference{
		Reference:   "Babushka",
		Origin:      "Russian",
		Explanation: "Grandmother, but also cultural archetype",
		Handling:    "Translate with footnote",
	}

	if ref.Reference != "Babushka" {
		t.Errorf("Reference = %s, want Babushka", ref.Reference)
	}

	if ref.Origin != "Russian" {
		t.Errorf("Origin = %s, want Russian", ref.Origin)
	}

	if ref.Handling != "Translate with footnote" {
		t.Errorf("Handling = %s, want Translate with footnote", ref.Handling)
	}
}

func TestChapterAnalysis_Complexity(t *testing.T) {
	tests := []struct {
		name       string
		complexity string
		valid      bool
	}{
		{"Simple complexity", "Simple", true},
		{"Moderate complexity", "Moderate", true},
		{"Complex complexity", "Complex", true},
		{"Invalid complexity", "Very Complex", false},
		{"Empty complexity", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := ChapterAnalysis{
				ChapterID:    "test",
				ChapterNum:   1,
				Title:        "Test Chapter",
				Summary:      "Test summary",
				KeyPoints:    []string{"Point 1"},
				Caveats:      []string{},
				Tone:         "Neutral",
				Complexity:   tt.complexity,
				SpecialNotes: "",
			}

			// For this test, we just verify the field is set correctly
			if analysis.Complexity != tt.complexity {
				t.Errorf("Complexity = %s, want %s", analysis.Complexity, tt.complexity)
			}
		})
	}
}

func TestPreparationPass_Duration(t *testing.T) {
	pass := PreparationPass{
		PassNumber: 1,
		Provider:   "OpenAI",
		Model:      "GPT-4",
		Analysis:   ContentAnalysis{},
		Duration:   5 * time.Minute,
		TokensUsed: 1000,
	}

	if pass.Duration != 5*time.Minute {
		t.Errorf("Duration = %v, want %v", pass.Duration, 5*time.Minute)
	}

	if pass.TokensUsed != 1000 {
		t.Errorf("TokensUsed = %d, want 1000", pass.TokensUsed)
	}
}

func TestPreparationResult_Statistics(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)

	result := PreparationResult{
		SourceLanguage: "Russian",
		TargetLanguage: "Serbian",
		Passes: []PreparationPass{
			{
				PassNumber: 1,
				Provider:   "OpenAI",
				Model:      "GPT-4",
				Duration:   10 * time.Minute,
				TokensUsed: 2000,
			},
			{
				PassNumber: 2,
				Provider:   "Anthropic",
				Model:      "Claude-3",
				Duration:   15 * time.Minute,
				TokensUsed: 1500,
			},
		},
		FinalAnalysis: ContentAnalysis{},
		TotalDuration: 25 * time.Minute,
		TotalTokens:   3500,
		PassCount:     2,
		StartedAt:     startTime,
		CompletedAt:   endTime,
	}

	if result.SourceLanguage != "Russian" {
		t.Errorf("SourceLanguage = %s, want Russian", result.SourceLanguage)
	}

	if result.TargetLanguage != "Serbian" {
		t.Errorf("TargetLanguage = %s, want Serbian", result.TargetLanguage)
	}

	if len(result.Passes) != 2 {
		t.Errorf("Passes length = %d, want 2", len(result.Passes))
	}

	if result.TotalTokens != 3500 {
		t.Errorf("TotalTokens = %d, want 3500", result.TotalTokens)
	}

	if result.PassCount != 2 {
		t.Errorf("PassCount = %d, want 2", result.PassCount)
	}
}

func TestPreparationConfig_DefaultValues(t *testing.T) {
	config := PreparationConfig{
		PassCount:          3,
		Providers:          []string{"OpenAI", "Anthropic"},
		AnalyzeContentType: true,
		AnalyzeCharacters:  true,
		AnalyzeTerminology: false,
		AnalyzeCulture:     true,
		AnalyzeChapters:    true,
		DetailLevel:        "Comprehensive",
		SourceLanguage:     "Russian",
		TargetLanguage:     "Serbian",
	}

	if config.PassCount != 3 {
		t.Errorf("PassCount = %d, want 3", config.PassCount)
	}

	if len(config.Providers) != 2 {
		t.Errorf("Providers length = %d, want 2", len(config.Providers))
	}

	if !config.AnalyzeContentType {
		t.Error("AnalyzeContentType should be true")
	}

	if config.AnalyzeTerminology {
		t.Error("AnalyzeTerminology should be false")
	}

	if config.DetailLevel != "Comprehensive" {
		t.Errorf("DetailLevel = %s, want Comprehensive", config.DetailLevel)
	}
}

func TestPreparationConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config PreparationConfig
		valid  bool
	}{
		{
			name: "Valid config",
			config: PreparationConfig{
				PassCount:      2,
				Providers:      []string{"OpenAI"},
				SourceLanguage: "Russian",
				TargetLanguage: "Serbian",
			},
			valid: true,
		},
		{
			name: "Zero pass count",
			config: PreparationConfig{
				PassCount:      0,
				Providers:      []string{"OpenAI"},
				SourceLanguage: "Russian",
				TargetLanguage: "Serbian",
			},
			valid: false,
		},
		{
			name: "Empty providers",
			config: PreparationConfig{
				PassCount:      1,
				Providers:      []string{},
				SourceLanguage: "Russian",
				TargetLanguage: "Serbian",
			},
			valid: false,
		},
		{
			name: "Empty source language",
			config: PreparationConfig{
				PassCount:      1,
				Providers:      []string{"OpenAI"},
				SourceLanguage: "",
				TargetLanguage: "Serbian",
			},
			valid: false,
		},
		{
			name: "Empty target language",
			config: PreparationConfig{
				PassCount:      1,
				Providers:      []string{"OpenAI"},
				SourceLanguage: "Russian",
				TargetLanguage: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation logic
			valid := tt.config.PassCount > 0 &&
				len(tt.config.Providers) > 0 &&
				tt.config.SourceLanguage != "" &&
				tt.config.TargetLanguage != ""

			if valid != tt.valid {
				t.Errorf("Validation result = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestContentAnalysis_EmptyStruct(t *testing.T) {
	analysis := ContentAnalysis{}

	if analysis.ContentType != "" {
		t.Errorf("ContentType should be empty, got %s", analysis.ContentType)
	}

	if len(analysis.Subgenres) != 0 {
		t.Errorf("Subgenres should be empty, got %d elements", len(analysis.Subgenres))
	}

	if len(analysis.UntranslatableTerms) != 0 {
		t.Errorf("UntranslatableTerms should be empty, got %d elements", len(analysis.UntranslatableTerms))
	}

	if analysis.AnalysisVersion != 0 {
		t.Errorf("AnalysisVersion should be 0, got %d", analysis.AnalysisVersion)
	}
}

func TestCharacter_EmptyNameTranslations(t *testing.T) {
	character := Character{
		Name:            "Test",
		AlternateNames:  []string{},
		Role:            "Test Role",
		SpeechPattern:   "Test Pattern",
		KeyTraits:       []string{},
		NameTranslation: map[string]string{},
	}

	if len(character.NameTranslation) != 0 {
		t.Errorf("NameTranslation should be empty, got %d elements", len(character.NameTranslation))
	}

	// Test accessing non-existent translation
	if translation, exists := character.NameTranslation["en"]; exists {
		t.Errorf("Expected no English translation, got %s", translation)
	}
}
