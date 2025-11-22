# Preparation Phase Documentation

## Overview

The Preparation Phase is a multi-pass, multi-LLM content analysis system that runs **before** translation to extract critical context and metadata. This significantly improves translation quality by providing translators (LLMs) with deep understanding of:

- Content type and genre
- Characters and their speech patterns
- Untranslatable terms to preserve
- Cultural references needing footnotes
- Chapter-specific context and caveats

## Architecture

### Components

1. **PreparationCoordinator** (`pkg/preparation/coordinator.go`)
   - Orchestrates multi-pass analysis
   - Rotates between LLM providers
   - Consolidates results from all passes

2. **PreparationPromptBuilder** (`pkg/preparation/prompts.go`)
   - Generates prompts for initial analysis
   - Creates refinement prompts for subsequent passes
   - Builds chapter-specific analysis prompts
   - Constructs consolidation prompts

3. **PreparationAwareTranslator** (`pkg/preparation/translator.go`)
   - Wraps standard translator with preparation capabilities
   - Runs preparation phase before translation
   - Injects chapter-specific context into translation

4. **Utility Functions** (`pkg/preparation/utils.go`)
   - Save/load preparation results
   - Format human-readable summaries
   - Generate translation context strings

### Data Structures

```go
type ContentAnalysis struct {
    // Classification
    ContentType    string
    Genre          string
    Subgenres      []string

    // Style
    Tone           string
    LanguageStyle  string
    TargetAudience string

    // Translation guidance
    UntranslatableTerms []UntranslatableTerm
    FootnoteGuidance    []FootnoteGuidance
    Characters          []Character
    CulturalReferences  []CulturalReference
    KeyThemes           []string

    // Chapter analysis
    ChapterAnalyses []ChapterAnalysis
}
```

## Multi-Pass Analysis Flow

### Pass 1: Initial Analysis
- Full content analysis from scratch
- Identifies all key elements
- JSON-structured output

### Pass 2-N: Refinement Passes
- Reviews previous pass results
- Validates and corrects errors
- Adds missed elements
- Improves explanations
- Each pass uses different LLM provider for diversity

### Final: Consolidation
- Merges all passes into definitive analysis
- Resolves conflicts between passes
- Prioritizes most important items
- Creates final guide for translation

## Configuration

```go
type PreparationConfig struct {
    PassCount      int      // Number of analysis passes (default: 2)
    Providers      []string // LLM providers to rotate through

    // Analysis focus
    AnalyzeContentType  bool // Determine type/genre
    AnalyzeCharacters   bool // Extract character info
    AnalyzeTerminology  bool // Find untranslatable terms
    AnalyzeCulture      bool // Identify cultural references
    AnalyzeChapters     bool // Analyze each chapter

    DetailLevel    string // "basic", "standard", "comprehensive"
    SourceLanguage string
    TargetLanguage string
}
```

## Usage

### Command-Line Tool

```bash
# Build the tool
go build -o build/preparation-translator ./cmd/preparation-translator

# Run with preparation + translation
export DEEPSEEK_API_KEY="your-key"
export ZHIPU_API_KEY="your-key"

./build/preparation-translator \
  -input Books/MyBook.epub \
  -output Books/MyBook_translated.epub \
  -analysis Books/MyBook_analysis.json \
  -source "Russian" \
  -target "Serbian" \
  -passes 3
```

### Programmatic Usage

```go
package main

import (
    "context"
    "digital.vasic.translator/pkg/ebook"
    "digital.vasic.translator/pkg/language"
    "digital.vasic.translator/pkg/preparation"
    "digital.vasic.translator/pkg/translator/llm"
)

func main() {
    // Configure preparation
    prepConfig := &preparation.PreparationConfig{
        PassCount:          2,
        Providers:          []string{"deepseek", "zhipu"},
        AnalyzeContentType: true,
        AnalyzeCharacters:  true,
        AnalyzeTerminology: true,
        AnalyzeCulture:     true,
        AnalyzeChapters:    true,
        DetailLevel:        "comprehensive",
        SourceLanguage:     "Russian",
        TargetLanguage:     "Serbian",
    }

    // Create base translator
    baseTranslator, _ := llm.NewLLMTranslator(translator.TranslationConfig{
        SourceLang: "ru",
        TargetLang: "sr",
        Provider:   "deepseek",
    })

    // Create preparation-aware translator
    prepTranslator := preparation.NewPreparationAwareTranslator(
        baseTranslator,
        nil,
        language.Language{Code: "ru", Name: "Russian"},
        language.Language{Code: "sr", Name: "Serbian"},
        prepConfig,
    )

    // Parse book
    parser := ebook.NewUniversalParser()
    book, _ := parser.Parse("book.epub")

    // Run preparation + translation
    ctx := context.Background()
    err := prepTranslator.TranslateBook(ctx, book, nil, "session-id")

    // Save preparation analysis
    prepTranslator.SavePreparationAnalysis("analysis.json")

    // Save translated book
    writer := ebook.NewEPUBWriter()
    writer.Write(book, "translated.epub")
}
```

## Analysis Output

### JSON Structure

The preparation analysis is saved as JSON with the following structure:

```json
{
  "source_language": "Russian",
  "target_language": "Serbian",
  "passes": [
    {
      "pass_number": 1,
      "provider": "deepseek",
      "analysis": { ... },
      "duration": "45.2s",
      "tokens_used": 12500
    }
  ],
  "final_analysis": {
    "content_type": "Novel",
    "genre": "Detective Fiction",
    "subgenres": ["Psychological Thriller", "Crime"],
    "tone": "Dark, suspenseful, introspective",
    "language_style": "Literary, complex sentence structures",
    "target_audience": "Adults interested in crime fiction",

    "untranslatable_terms": [
      {
        "term": "матрёшка",
        "original_script": "Cyrillic",
        "reason": "Culture-specific object, no direct equivalent",
        "context": ["Chapter 3, paragraph 5"],
        "transliteration": "matryoshka"
      }
    ],

    "footnote_guidance": [
      {
        "term": "Белые ночи",
        "explanation": "White Nights - natural phenomenon in St. Petersburg",
        "locations": ["Chapter 7"],
        "priority": "high"
      }
    ],

    "characters": [
      {
        "name": "Иван Петрович",
        "alternate_names": ["Ваня"],
        "role": "Protagonist",
        "speech_pattern": "Formal, educated, uses literary references",
        "key_traits": ["Observant", "Melancholic"],
        "name_translation": {
          "sr": "Иван Петрович"
        }
      }
    ],

    "key_themes": [
      "Guilt and redemption",
      "Truth vs. appearance",
      "Social justice"
    ],

    "cultural_references": [
      {
        "reference": "Dostoevsky's Crime and Punishment",
        "origin": "Russian literature",
        "explanation": "Classic novel referenced as metaphor",
        "handling": "Keep reference, add footnote"
      }
    ],

    "chapter_analyses": [
      {
        "chapter_id": "chapter_1",
        "chapter_num": 1,
        "title": "Пролог",
        "summary": "Introduction of main character and setting",
        "key_points": [
          "Establishes dark atmosphere",
          "Introduces protagonist's background"
        ],
        "caveats": [
          "Contains archaic legal terminology",
          "Multiple literary allusions"
        ],
        "tone": "Ominous, foreboding",
        "complexity": "moderate"
      }
    ]
  },

  "total_duration": "120.5s",
  "total_tokens": 25000,
  "pass_count": 2
}
```

### Human-Readable Summary

```
=== PREPARATION ANALYSIS SUMMARY ===

Languages: Russian → Serbian
Duration: 120.50 seconds
Passes: 2
Total Tokens: 25000

--- CONTENT CLASSIFICATION ---
Type: Novel
Genre: Detective Fiction
Subgenres: [Psychological Thriller Crime]
Tone: Dark, suspenseful, introspective
Target Audience: Adults interested in crime fiction

--- KEY FINDINGS ---
Untranslatable Terms: 15
Footnotes Needed: 8
Characters: 12
Cultural References: 23
Key Themes: 5
Chapters Analyzed: 38

--- KEY THEMES ---
• Guilt and redemption
• Truth vs. appearance
• Social justice
• Memory and trauma
• Moral ambiguity

--- UNTRANSLATABLE TERMS (showing first 10) ---
• матрёшка: Culture-specific object, no direct equivalent
• самовар: Traditional Russian tea urn, cultural significance
• ...

--- CHARACTERS ---
• Иван Петрович (Protagonist)
  Speech: Formal, educated, uses literary references
• Анна Сергеевна (Supporting Character)
  Speech: Informal, emotional, colloquial
...

=== END SUMMARY ===
```

## Translation Context Generation

For each chapter, the preparation system generates a context string passed to the translation LLM:

```
## TRANSLATION CONTEXT

**Content Type**: Novel
**Genre**: Detective Fiction
**Tone**: Dark, suspenseful

**Terms to Keep in Original**:
- матрёшка: Culture-specific object, no direct equivalent
- самовар: Traditional Russian tea urn

**Character Speech Patterns**:
- Иван Петрович: Formal, educated, uses literary references
- Анна Сергеевна: Informal, emotional, colloquial

**Chapter 5 Context**:
Summary: Pivotal confrontation between protagonist and antagonist
Translation Caveats:
- Contains complex wordplay on "правда" (truth/justice)
- References to Pushkin require careful handling
- Legal terminology must be precise
```

This context helps the translation LLM:
1. Maintain character voice consistency
2. Preserve untranslatable terms
3. Handle cultural references appropriately
4. Match the appropriate tone and register
5. Be aware of chapter-specific challenges

## Benefits

### Translation Quality Improvements

1. **Character Consistency**
   - Each character's speech pattern is maintained
   - Formal/informal register preserved
   - Unique speech quirks retained

2. **Cultural Preservation**
   - Untranslatable terms kept in original
   - Cultural references identified for footnotes
   - Context provided for better understanding

3. **Genre Awareness**
   - Translation matches genre conventions
   - Tone appropriate for content type
   - Target audience considerations

4. **Context-Aware Translation**
   - Chapter-specific guidance
   - Awareness of literary allusions
   - Better handling of complex terminology

### Cost Efficiency

- Analysis done once, used for all chapters
- Prevents translation errors requiring expensive re-work
- Parallel chapter analysis with concurrency control
- Provider rotation reduces API rate limiting

## Advanced Features

### Parallel Chapter Analysis

The system analyzes chapters in parallel with a semaphore limiting concurrency:

```go
semaphore := make(chan struct{}, 3) // Max 3 concurrent analyses
for i, chapter := range book.Chapters {
    wg.Add(1)
    go func(chapterNum int, ch ebook.Chapter) {
        defer wg.Done()
        semaphore <- struct{}{}        // Acquire
        defer func() { <-semaphore }() // Release

        // Analyze chapter...
    }(i, chapter)
}
wg.Wait()
```

### Provider Rotation

Different LLM providers are used across passes for diversity:

```go
providerIndex := (passNum - 1) % len(providers)
provider := providers[providerIndex]
```

This ensures:
- Different perspectives on the content
- Reduced bias from single provider
- Better consolidation through diversity

### Analysis Consolidation

The final pass merges multiple analyses using an LLM:

```
Create the DEFINITIVE, HIGHEST-QUALITY analysis by:
1. Merging: Combine insights from all passes
2. Validating: Include only accurate information
3. Deduplicating: Remove redundant entries
4. Prioritizing: Keep the most important items
5. Clarifying: Use the clearest explanations
```

## Troubleshooting

### Common Issues

**Issue**: Preparation phase fails with "no valid LLM providers"
**Solution**: Ensure API keys are set as environment variables:
```bash
export DEEPSEEK_API_KEY="your-key"
export ZHIPU_API_KEY="your-key"
```

**Issue**: JSON parsing errors in analysis output
**Solution**: LLM may include extra text. The system extracts JSON automatically, but verify prompt configuration.

**Issue**: Chapter analysis takes too long
**Solution**: Reduce concurrency limit or disable chapter analysis:
```go
config.AnalyzeChapters = false
```

**Issue**: Consolidation produces inconsistent results
**Solution**: Increase pass count for more data points:
```go
config.PassCount = 3  // or more
```

## Future Enhancements

Planned improvements:

1. **Caching**: Store analyses for reuse across translation attempts
2. **Incremental Analysis**: Analyze only changed chapters
3. **Quality Scoring**: Rate analysis quality and re-run if needed
4. **Custom Prompts**: Allow user-defined analysis prompts
5. **Analysis Visualization**: Web UI for exploring analysis results
6. **Multi-Language Support**: Extend beyond Russian-Serbian
7. **Integration with Translation Memory**: Leverage TMX files
8. **Automated Footnote Generation**: Create footnotes from guidance

## Performance Metrics

Typical performance (2 passes, 38 chapters):

- **Preparation Time**: 2-5 minutes
- **Token Usage**: 20,000-50,000 tokens
- **Chapter Analysis**: 30-60 seconds (parallel)
- **Memory Usage**: < 100MB
- **API Calls**: ~80-150 (depends on passes and chapters)

## Best Practices

1. **Use at least 2 passes** for comprehensive analysis
2. **Enable chapter analysis** for long-form content (novels)
3. **Rotate providers** for diverse perspectives
4. **Save analysis JSON** for future reference and debugging
5. **Review untranslatable terms** before translation starts
6. **Validate cultural references** with native speakers
7. **Test with sample chapters** before full book translation

## See Also

- [Translation Pipeline Documentation](TRANSLATION.md)
- [LLM Provider Configuration](LLM_PROVIDERS.md)
- [Quality Verification Guide](VERIFICATION.md)
- [API Reference](API.md)
