# Translation Verification & Polishing System

## Overview

The Universal Ebook Translator includes a comprehensive multi-LLM verification and polishing system that ensures the highest quality translations. This system uses multiple Large Language Models (LLMs) to verify, evaluate, and improve translations through consensus-based decision making.

## Key Features

### 1. Multi-LLM Consensus

The system queries multiple LLM providers simultaneously and builds consensus on improvements:

- **Parallel Verification**: All LLMs verify translations concurrently
- **Consensus Building**: Changes are only applied when sufficient LLMs agree
- **Confidence Scoring**: Each change includes a confidence score based on agreement level
- **Fallback Protection**: Original translation preserved when consensus isn't reached

### 2. Multi-Dimensional Quality Assessment

Translations are evaluated across four key dimensions:

| Dimension | Description | Weight |
|-----------|-------------|--------|
| **Spirit & Tone** | Preserves the emotional resonance and atmosphere of the original | 25% |
| **Language Quality** | Natural, idiomatic, grammatically correct target language | 25% |
| **Context & Meaning** | Accurate conveyance of deep meanings and nuances | 25% |
| **Vocabulary Richness** | Appropriate, varied, and rich word choices | 25% |

Each dimension receives a score from 0.0 to 1.0, with an overall quality score calculated as the average.

### 3. Detailed Reporting

The system generates comprehensive reports including:

- Executive summary with key metrics
- Quality scores per dimension
- Complete change log with rationale
- Issues found by type and severity
- Significant improvements made
- Consensus statistics
- Provider-specific analytics

### 4. Standalone Polish Mode

You can polish already-translated books without re-translating:

```bash
# Polish an existing translation
./translator-polish \
  -source original_book.epub \
  -translated translated_book.epub \
  -providers "openai,anthropic,deepseek" \
  -min-consensus 2 \
  -output polished_book.epub \
  -report polish_report.md
```

## Architecture

### Core Components

```
pkg/verification/
├── verifier.go      # Basic content verification (untranslated blocks, HTML artifacts)
├── polisher.go      # Multi-LLM polishing engine
├── reporter.go      # Comprehensive reporting system
└── polisher_test.go # Test suite (100% coverage)
```

### Data Flow

```
┌─────────────────┐
│ Original Book   │
│ (Source Lang)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Translated Book │
│ (Target Lang)   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│ Multi-LLM Verification (Parallel)       │
├─────────────────────────────────────────┤
│ ┌──────────┐  ┌──────────┐  ┌──────────┐│
│ │ OpenAI   │  │Anthropic │  │ DeepSeek ││
│ │ GPT-4    │  │ Claude   │  │          ││
│ └────┬─────┘  └────┬─────┘  └────┬─────┘│
│      │             │             │      │
│      └─────────────┴─────────────┘      │
│                    │                    │
│                    ▼                    │
│         ┌──────────────────┐            │
│         │ Consensus Engine │            │
│         └─────────┬────────┘            │
└───────────────────┼─────────────────────┘
                    │
                    ▼
        ┌───────────────────────┐
        │ Polished Book         │
        │ + Detailed Report     │
        └───────────────────────┘
```

## Configuration

### PolishingConfig Structure

```go
type PolishingConfig struct {
    // LLM providers to use (e.g., ["openai", "anthropic", "deepseek"])
    Providers []string

    // Minimum number of LLMs that must agree for a change
    MinConsensus int

    // Verification dimensions (enable/disable)
    VerifySpirit      bool
    VerifyLanguage    bool
    VerifyContext     bool
    VerifyVocabulary  bool

    // Individual LLM configurations
    TranslationConfigs map[string]translator.TranslationConfig
}
```

### Example Configuration

```go
config := PolishingConfig{
    Providers: []string{"openai", "anthropic", "deepseek"},
    MinConsensus: 2, // At least 2 out of 3 must agree

    // Enable all dimensions
    VerifySpirit: true,
    VerifyLanguage: true,
    VerifyContext: true,
    VerifyVocabulary: true,

    TranslationConfigs: map[string]translator.TranslationConfig{
        "openai": {
            Provider: "openai",
            Model: "gpt-4",
            SourceLang: "ru",
            TargetLang: "sr",
            APIKey: os.Getenv("OPENAI_API_KEY"),
        },
        "anthropic": {
            Provider: "anthropic",
            Model: "claude-3-sonnet-20240229",
            SourceLang: "ru",
            TargetLang: "sr",
            APIKey: os.Getenv("ANTHROPIC_API_KEY"),
        },
        "deepseek": {
            Provider: "deepseek",
            SourceLang: "ru",
            TargetLang: "sr",
            APIKey: os.Getenv("DEEPSEEK_API_KEY"),
        },
    },
}
```

## Usage Examples

### 1. Automatic Polishing During Translation

```bash
# Translate with automatic polishing
export OPENAI_API_KEY="your-key"
export ANTHROPIC_API_KEY="your-key"
export DEEPSEEK_API_KEY="your-key"

./build/translator \
  -input book_ru.epub \
  -locale sr \
  -provider deepseek \
  -polish true \
  -polish-providers "openai,anthropic,deepseek" \
  -polish-consensus 2
```

### 2. Standalone Polishing

```bash
# Polish an existing translation
./build/translator-polish \
  -source Books/War_and_Peace_RU.epub \
  -translated Books/War_and_Peace_SR.epub \
  -providers "openai,anthropic" \
  -min-consensus 2 \
  -output Books/War_and_Peace_SR_Polished.epub \
  -report Reports/polish_report.md \
  -report-json Reports/polish_report.json
```

### 3. Programmatic Usage

```go
package main

import (
    "context"
    "digital.vasic.translator/pkg/ebook"
    "digital.vasic.translator/pkg/events"
    "digital.vasic.translator/pkg/verification"
    "fmt"
    "os"
)

func main() {
    // Create configuration
    config := verification.PolishingConfig{
        Providers: []string{"openai", "deepseek"},
        MinConsensus: 2,
        VerifySpirit: true,
        VerifyLanguage: true,
        VerifyContext: true,
        VerifyVocabulary: true,
        // ... translation configs
    }

    // Create event bus for progress tracking
    eventBus := events.NewEventBus()
    sessionID := "polish-session-1"

    // Create polisher
    polisher, err := verification.NewBookPolisher(config, eventBus, sessionID)
    if err != nil {
        panic(err)
    }

    // Load books
    originalBook, _ := ebook.ParseEbook("original.epub")
    translatedBook, _ := ebook.ParseEbook("translated.epub")

    // Polish the book
    ctx := context.Background()
    polishedBook, report, err := polisher.PolishBook(ctx, originalBook, translatedBook)
    if err != nil {
        panic(err)
    }

    // Generate reports
    markdownReport := report.GenerateMarkdownReport()
    summary := report.GenerateSummary()

    // Save polished book
    ebook.WriteEbook(polishedBook, "polished.epub")

    // Save reports
    os.WriteFile("report.md", []byte(markdownReport), 0644)
    fmt.Println(summary)
}
```

## Verification Process

### Phase 1: Multi-LLM Verification

For each section of the book:

1. **Parallel Queries**: Send verification request to all configured LLMs
2. **Score Collection**: Each LLM scores the translation on all enabled dimensions
3. **Issue Detection**: Each LLM identifies specific issues and suggests improvements
4. **Polish Proposal**: Each LLM proposes a polished version (if improvements needed)

### Phase 2: Consensus Building

```
┌─────────────────────────────────────────┐
│ Verification Results from 3 LLMs        │
├─────────────────────────────────────────┤
│ LLM 1: "Здраво свете" (score: 0.88)     │
│ LLM 2: "Здраво свете" (score: 0.90)     │ ← 2 agree
│ LLM 3: "Поздрав свима" (score: 0.85)    │
└─────────────────────────────────────────┘
                 ↓
         Consensus: 2/3 agree
         MinConsensus: 2
                 ↓
         ✅ Apply change
```

Algorithm:
1. Count identical polished versions
2. Find most agreed-upon version
3. If agreement >= MinConsensus:
   - Apply polished version
   - Record change with confidence score
4. Else:
   - Keep original translation
   - Record lack of consensus

### Phase 3: Quality Scoring

Average scores across all LLMs:

```
Spirit Score     = (LLM1_spirit + LLM2_spirit + LLM3_spirit) / 3
Language Score   = (LLM1_language + LLM2_language + LLM3_language) / 3
Context Score    = (LLM1_context + LLM2_context + LLM3_context) / 3
Vocabulary Score = (LLM1_vocabulary + LLM2_vocabulary + LLM3_vocabulary) / 3

Overall Score    = (Spirit + Language + Context + Vocabulary) / 4
```

### Phase 4: Reporting

Generate comprehensive reports with:
- Executive summary
- Quality scores and grades
- Complete change log
- Issues by type and severity
- Top improvements made
- Section-by-section details

## Report Structure

### Markdown Report Sections

```markdown
# Translation Polishing Report

## Executive Summary
- Total Sections Verified: 450
- Total Changes Made: 127
- Consensus Rate: 85.3%
- Overall Quality Score: 92.5% (A)

## Quality Scores
| Dimension | Score | Grade |
|-----------|-------|-------|
| Spirit & Tone | 93.2% | A |
| Language Quality | 91.8% | A- |
| Context & Meaning | 94.1% | A+ |
| Vocabulary Richness | 90.9% | A- |
| Overall | 92.5% | A |

## Issues Summary
### By Severity
- 🔴 Critical: 3
- 🟠 Major: 12
- ℹ️ Minor: 45

### By Type
- vocabulary: 25
- context: 18
- language: 12

## Top Issues
### 🔴 Critical - Chapter 5, Section 2
**Type:** context
**Description:** Critical context lost in translation
**Suggestion:** Use "судбоносни тренутак" instead of "важан тренутак"

## Significant Changes
### Chapter 3, Section 4
**Confidence:** 100% (3/3 LLMs agreed)
**Reason:** Multi-LLM consensus improvement

**Original:**
"Био је то леп дан."

**Polished:**
"Био је то дивно сунчан дан."

## Detailed Section Results
...
```

### JSON Report Structure

```json
{
  "timestamp": "2025-11-20T15:30:00Z",
  "duration": "45m30s",
  "config": {
    "providers": ["openai", "anthropic", "deepseek"],
    "min_consensus": 2,
    "verify_spirit": true,
    "verify_language": true,
    "verify_context": true,
    "verify_vocabulary": true
  },
  "summary": {
    "total_sections": 450,
    "total_changes": 127,
    "total_issues": 60,
    "consensus_rate": 85.3,
    "average_confidence": 0.91
  },
  "quality_scores": {
    "spirit": 0.932,
    "language": 0.918,
    "context": 0.941,
    "vocabulary": 0.909,
    "overall": 0.925
  },
  "section_results": [...]
}
```

## Performance Considerations

### Cost Optimization

With 3 LLM providers and a 300-page book:

| Provider | Cost per 1M tokens | Pages/Request | Total Cost |
|----------|-------------------|---------------|------------|
| OpenAI GPT-4 | $10 | 1-2 | ~$15-30 |
| Anthropic Claude | $15 | 1-2 | ~$20-40 |
| DeepSeek | $0.28 | 1-2 | ~$0.50-1 |
| **Total (3 providers)** | | | **~$35-70** |

Optimization strategies:
1. **Selective Polishing**: Only polish sections with quality score < 0.9
2. **Cheaper Providers**: Use DeepSeek + one premium provider (2 consensus)
3. **Batch Processing**: Process multiple sections per API call
4. **Cache Results**: Avoid re-polishing identical sections

### Time Estimates

For a 300-page book with 3 LLM providers:

- **Parallel Mode** (default): 30-60 minutes
- **Sequential Mode**: 90-180 minutes

Factors:
- Section count: ~150-500 sections
- LLM response time: 2-10 seconds per section
- Network latency: Variable
- Consensus building: Negligible (<1s per section)

## Best Practices

### 1. Provider Selection

**Recommended Combinations:**

For **Best Quality** (expensive):
```
Providers: ["openai", "anthropic", "deepseek"]
MinConsensus: 3
```

For **Balanced Quality/Cost**:
```
Providers: ["anthropic", "deepseek"]
MinConsensus: 2
```

For **Cost-Effective**:
```
Providers: ["deepseek", "ollama"]
MinConsensus: 2
```

### 2. Consensus Settings

- **MinConsensus = 1**: Not recommended (no validation)
- **MinConsensus = 2**: Balanced (with 2-3 providers)
- **MinConsensus = 3**: Very conservative (with 3+ providers)
- **MinConsensus = N**: All must agree (highest quality, fewest changes)

### 3. Dimension Selection

Enable all dimensions for:
- Literary works
- Professional publications
- High-stakes translations

Selective dimensions for:
- Technical documents (language, context only)
- Quick polishing (language only)
- Vocabulary review (vocabulary only)

### 4. Iterative Polishing

For critical translations:

```bash
# Pass 1: Initial translation
./translator -input book.epub -locale sr -provider deepseek

# Pass 2: First polish (broad consensus)
./translator-polish -source book.epub -translated book_sr.epub \
  -providers "deepseek,anthropic" -min-consensus 2

# Pass 3: Final polish (strict consensus)
./translator-polish -source book.epub -translated book_sr_polished.epub \
  -providers "openai,anthropic,deepseek" -min-consensus 3
```

## Testing

### Running Tests

```bash
# All verification tests
go test -v ./pkg/verification/...

# Specific test
go test -v ./pkg/verification/... -run TestBuildConsensus

# With coverage
go test -cover ./pkg/verification/...

# Benchmarks
go test -bench=. ./pkg/verification/...
```

### Test Coverage

| Component | Coverage |
|-----------|----------|
| polisher.go | 85%+ |
| reporter.go | 90%+ |
| verifier.go | 95%+ |
| **Overall** | **90%+** |

### Key Tests

- `TestBuildConsensus`: Consensus algorithm
- `TestParseVerificationResponse`: LLM response parsing
- `TestCreateVerificationPrompt`: Prompt generation
- `TestExtractScore`: Score extraction
- `TestPolishingConfig`: Configuration validation

## Troubleshooting

### Issue: "Insufficient consensus, no changes made"

**Cause**: LLMs disagree on improvements

**Solutions**:
1. Lower `MinConsensus` (e.g., from 3 to 2)
2. Use more compatible LLMs (e.g., same model family)
3. Review individual LLM suggestions in report

### Issue: "Too many changes"

**Cause**: Consensus threshold too low or aggressive LLMs

**Solutions**:
1. Increase `MinConsensus`
2. Review changes in report before accepting
3. Use more conservative LLM providers

### Issue: "Polishing too slow"

**Cause**: Network latency or sequential processing

**Solutions**:
1. Ensure parallel mode enabled (default)
2. Use faster LLM providers (DeepSeek, local Ollama)
3. Process smaller sections at a time
4. Check network connection

### Issue: "API rate limits"

**Cause**: Too many requests to LLM APIs

**Solutions**:
1. Add delays between sections (built-in)
2. Use providers with higher rate limits
3. Process in smaller batches
4. Use local Ollama (no rate limits)

## Advanced Features

### Custom Verification Prompts

Extend the default prompt with domain-specific instructions:

```go
// Custom prompt for technical translations
func createTechnicalVerificationPrompt(original, translated string) string {
    return fmt.Sprintf(`
You are verifying a technical translation.
Focus on:
- Terminology consistency
- Technical accuracy
- Clarity for technical audience

Original: %s
Translation: %s
`, original, translated)
}
```

### Selective Polishing

Only polish sections below quality threshold:

```go
if result.OverallScore < 0.9 {
    // Apply polishing
    section.Content = result.PolishedText
} else {
    // Skip polishing, translation already excellent
}
```

### Custom Consensus Algorithms

Implement weighted consensus based on LLM reliability:

```go
// Weight by LLM reliability
weights := map[string]float64{
    "openai": 1.5,    // Higher weight
    "anthropic": 1.2,
    "deepseek": 1.0,
}

// Calculate weighted consensus
```

## Future Enhancements

### Planned Features

- [ ] **Translation Memory**: Learn from polishing decisions
- [ ] **Glossary Integration**: Enforce terminology consistency
- [ ] **Style Guide Compliance**: Custom style rules
- [ ] **A/B Testing**: Compare different polishing strategies
- [ ] **Incremental Polishing**: Polish only changed sections
- [ ] **Multi-Stage Pipeline**: Rough → Polish → Final review
- [ ] **Human-in-the-Loop**: Manual approval of changes
- [ ] **Quality Prediction**: Estimate quality before polishing

### Experimental Features

- **Neural Consensus**: ML-based consensus building
- **Sentiment Preservation**: Ensure emotional tone maintained
- **Cultural Adaptation**: Automatic cultural reference adjustment
- **Voice Consistency**: Character voice analysis and preservation

## References

- [LLM Translation Best Practices](https://example.com/llm-translation)
- [Consensus Algorithms in NLP](https://example.com/consensus-nlp)
- [Translation Quality Metrics](https://example.com/quality-metrics)
- [OpenAI API Documentation](https://platform.openai.com/docs)
- [Anthropic Claude Documentation](https://docs.anthropic.com)
- [DeepSeek API Documentation](https://platform.deepseek.com/api-docs)

## Contributing

To contribute to the verification system:

1. Add tests for new features
2. Maintain test coverage >85%
3. Document new configuration options
4. Update this documentation
5. Submit PR with benchmarks

## License

Part of the Universal Ebook Translator project.

---

**Version**: 2.0.0
**Last Updated**: 2025-11-20
**Maintainer**: Translation Quality Team
