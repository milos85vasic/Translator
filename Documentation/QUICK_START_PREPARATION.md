# Quick Start: Preparation Phase

## What is the Preparation Phase?

The preparation phase analyzes your book **before** translation to extract:
- Content type, genre, and tone
- Character speech patterns
- Untranslatable terms to preserve
- Cultural references needing footnotes
- Chapter-specific translation guidance

This context dramatically improves translation quality.

## Installation

```bash
# Build the preparation-aware translator
cd /path/to/Translate
go build -o build/preparation-translator ./cmd/preparation-translator
```

## Basic Usage

### 1. Set Your API Keys

```bash
export DEEPSEEK_API_KEY="your-deepseek-key"
export ZHIPU_API_KEY="your-zhipu-key"
```

### 2. Run Translation with Preparation

```bash
./build/preparation-translator \
  -input Books/MyBook.epub \
  -output Books/MyBook_translated.epub \
  -analysis Books/MyBook_analysis.json
```

That's it! The tool will:
1. Run 2 passes of content analysis (default)
2. Analyze each chapter in parallel
3. Translate using the preparation context
4. Save the analysis as JSON for reference

## Advanced Options

### More Analysis Passes

More passes = better analysis quality:

```bash
./build/preparation-translator \
  -input Books/MyBook.epub \
  -output Books/MyBook_translated.epub \
  -analysis Books/MyBook_analysis.json \
  -passes 3
```

### Specify Languages

```bash
./build/preparation-translator \
  -input Books/MyBook.epub \
  -output Books/MyBook_translated.epub \
  -source "Russian" \
  -target "Serbian"
```

### All Options

```
-input string
    Input ebook path (default "/tmp/markdown_e2e_source.md")
-output string
    Output EPUB path (default "/tmp/prepared_translated.epub")
-analysis string
    Preparation analysis output path (default "/tmp/preparation_analysis.json")
-source string
    Source language (default "Russian")
-target string
    Target language (default "Serbian")
-passes int
    Number of preparation passes (default 2)
-providers string
    Comma-separated list of LLM providers (default "deepseek,zhipu")
```

## Understanding the Analysis Output

The analysis JSON file contains:

### Content Classification
```json
{
  "content_type": "Novel",
  "genre": "Detective Fiction",
  "subgenres": ["Psychological Thriller"],
  "tone": "Dark, suspenseful",
  "target_audience": "Adults"
}
```

### Untranslatable Terms
```json
{
  "untranslatable_terms": [
    {
      "term": "матрёшка",
      "reason": "Culture-specific object",
      "transliteration": "matryoshka"
    }
  ]
}
```

### Characters
```json
{
  "characters": [
    {
      "name": "Иван Петрович",
      "role": "Protagonist",
      "speech_pattern": "Formal, educated"
    }
  ]
}
```

### Cultural References
```json
{
  "cultural_references": [
    {
      "reference": "Dostoevsky's Crime and Punishment",
      "explanation": "Classic novel referenced as metaphor",
      "handling": "Keep reference, add footnote"
    }
  ]
}
```

### Chapter Analyses
```json
{
  "chapter_analyses": [
    {
      "chapter_num": 1,
      "title": "Пролог",
      "summary": "Introduction of main character",
      "caveats": ["Archaic legal terminology"],
      "complexity": "moderate"
    }
  ]
}
```

## How It Improves Translation

### Without Preparation
- Generic translation
- Character voices inconsistent
- Cultural terms mistranslated
- No context awareness

### With Preparation
- Context-aware translation
- Character voices preserved
- Cultural terms handled appropriately
- Chapter-specific guidance applied

## Example Workflow

```bash
# 1. Build the tool
go build -o build/preparation-translator ./cmd/preparation-translator

# 2. Set API keys
export DEEPSEEK_API_KEY="sk-xxx"
export ZHIPU_API_KEY="xxx.xxx"

# 3. Run preparation + translation
./build/preparation-translator \
  -input Books/Crime_and_Punishment_ru.epub \
  -output Books/Crime_and_Punishment_sr.epub \
  -analysis Books/Crime_and_Punishment_analysis.json \
  -passes 3

# 4. Review the analysis
cat Books/Crime_and_Punishment_analysis.json | jq '.final_analysis | keys'

# 5. Check specific aspects
cat Books/Crime_and_Punishment_analysis.json | \
  jq '.final_analysis.untranslatable_terms[] | .term'

# 6. Review character speech patterns
cat Books/Crime_and_Punishment_analysis.json | \
  jq '.final_analysis.characters[] | {name, speech_pattern}'
```

## Programmatic Usage

If you want to integrate preparation into your own code:

```go
package main

import (
    "context"
    "digital.vasic.translator/pkg/ebook"
    "digital.vasic.translator/pkg/language"
    "digital.vasic.translator/pkg/preparation"
    "digital.vasic.translator/pkg/translator"
    "digital.vasic.translator/pkg/translator/llm"
)

func main() {
    // 1. Configure preparation
    prepConfig := &preparation.PreparationConfig{
        PassCount:          2,
        Providers:          []string{"deepseek", "zhipu"},
        AnalyzeContentType: true,
        AnalyzeCharacters:  true,
        AnalyzeTerminology: true,
        AnalyzeCulture:     true,
        AnalyzeChapters:    true,
        SourceLanguage:     "Russian",
        TargetLanguage:     "Serbian",
    }

    // 2. Create base translator
    config := translator.TranslationConfig{
        SourceLang: "ru",
        TargetLang: "sr",
        Provider:   "deepseek",
    }
    baseTranslator, _ := llm.NewLLMTranslator(config)

    // 3. Create preparation-aware translator
    prepTranslator := preparation.NewPreparationAwareTranslator(
        baseTranslator,
        nil,
        language.Language{Code: "ru", Name: "Russian"},
        language.Language{Code: "sr", Name: "Serbian"},
        prepConfig,
    )

    // 4. Parse book
    parser := ebook.NewUniversalParser()
    book, _ := parser.Parse("input.epub")

    // 5. Run preparation + translation
    ctx := context.Background()
    prepTranslator.TranslateBook(ctx, book, nil, "session-id")

    // 6. Save analysis
    prepTranslator.SavePreparationAnalysis("analysis.json")

    // 7. Save translated book
    writer := ebook.NewEPUBWriter()
    writer.Write(book, "output.epub")
}
```

## Performance Tips

### For Large Books (100+ chapters)
- Use 2 passes (default) initially
- Disable chapter analysis if time-constrained:
  ```go
  prepConfig.AnalyzeChapters = false
  ```

### For Maximum Quality
- Use 3-4 passes
- Enable all analysis types
- Review analysis JSON before translation

### For Cost Optimization
- Use fewer passes (1-2)
- Use cost-effective providers (DeepSeek)
- Cache analysis for re-translation attempts

## Troubleshooting

### "No valid LLM providers available"
Make sure API keys are set:
```bash
echo $DEEPSEEK_API_KEY
echo $ZHIPU_API_KEY
```

### Slow Chapter Analysis
Reduce concurrency or disable:
```go
prepConfig.AnalyzeChapters = false
```

### JSON Parsing Errors
The tool automatically extracts JSON from LLM responses. If issues persist, check the raw analysis in the passes array.

### Translation Not Using Context
Verify the preparation phase completed successfully by checking:
```bash
cat analysis.json | jq '.final_analysis.content_type'
```

## Next Steps

- Read the [full documentation](PREPARATION_PHASE.md)
- Explore the [analysis JSON structure](PREPARATION_PHASE.md#analysis-output)
- Learn about [multi-pass analysis](PREPARATION_PHASE.md#multi-pass-analysis-flow)
- See [advanced features](PREPARATION_PHASE.md#advanced-features)

## Quick Reference Commands

```bash
# Basic run
./build/preparation-translator -input book.epub -output translated.epub

# With 3 passes
./build/preparation-translator -input book.epub -output translated.epub -passes 3

# Custom analysis location
./build/preparation-translator -input book.epub -analysis my_analysis.json

# View content type
cat analysis.json | jq '.final_analysis.content_type'

# List untranslatable terms
cat analysis.json | jq '.final_analysis.untranslatable_terms[].term'

# View characters
cat analysis.json | jq '.final_analysis.characters[] | {name, role}'

# Check preparation stats
cat analysis.json | jq '{passes: .pass_count, duration: .total_duration, tokens: .total_tokens}'
```
