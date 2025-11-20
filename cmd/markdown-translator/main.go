package main

import (
	"context"
	"digital.vasic.translator/pkg/markdown"
	"digital.vasic.translator/pkg/translator"
	"digital.vasic.translator/pkg/translator/llm"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Command line flags
	inputFile := flag.String("input", "", "Input EPUB file")
	outputFile := flag.String("output", "", "Output EPUB file (optional, auto-generated if not provided)")
	targetLang := flag.String("lang", "sr", "Target language code (default: Serbian)")
	provider := flag.String("provider", "deepseek", "LLM provider (deepseek, openai, anthropic)")
	model := flag.String("model", "", "LLM model (optional, uses provider default)")
	keepMarkdown := flag.Bool("keep-md", true, "Keep intermediate markdown files")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Usage: markdown-translator -input <epub_file> [-output <output_file>] [-lang <language>] [-provider <provider>] [-keep-md]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate input file
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", *inputFile)
	}

	// Generate output filename if not provided
	if *outputFile == "" {
		base := strings.TrimSuffix(filepath.Base(*inputFile), filepath.Ext(*inputFile))
		*outputFile = fmt.Sprintf("%s_%s_md.epub", base, *targetLang)
	}

	// Generate intermediate markdown filenames
	sourceMD := strings.TrimSuffix(*outputFile, ".epub") + "_source.md"
	translatedMD := strings.TrimSuffix(*outputFile, ".epub") + "_translated.md"

	fmt.Printf("🚀 Markdown-Based Translation Pipeline\n\n")
	fmt.Printf("Input:  %s\n", *inputFile)
	fmt.Printf("Output: %s\n\n", *outputFile)

	// Step 1: EPUB → Markdown
	fmt.Println("📖 Step 1/4: Converting EPUB to Markdown...")
	converter := markdown.NewEPUBToMarkdownConverter(false, "")
	if err := converter.ConvertEPUBToMarkdown(*inputFile, sourceMD); err != nil {
		log.Fatalf("Failed to convert EPUB to Markdown: %v", err)
	}
	fmt.Printf("✓ Source markdown saved: %s\n\n", sourceMD)

	// Step 2: Create translator
	fmt.Println("🔧 Step 2/4: Initializing translator...")
	llmTranslator, err := createTranslator(*provider, *model, *targetLang)
	if err != nil {
		log.Fatalf("Failed to create translator: %v", err)
	}
	fmt.Printf("✓ Using provider: %s\n\n", *provider)

	// Step 3: Translate Markdown
	fmt.Println("🌍 Step 3/4: Translating markdown content...")
	ctx := context.Background()
	mdTranslator := markdown.NewMarkdownTranslator(func(text string) (string, error) {
		return llmTranslator.Translate(ctx, text, "")
	})

	if err := mdTranslator.TranslateMarkdownFile(sourceMD, translatedMD); err != nil {
		log.Fatalf("Failed to translate markdown: %v", err)
	}
	fmt.Printf("✓ Translated markdown saved: %s\n\n", translatedMD)

	// Step 4: Markdown → EPUB
	fmt.Println("📚 Step 4/4: Converting translated markdown to EPUB...")
	epubConverter := markdown.NewMarkdownToEPUBConverter()
	if err := epubConverter.ConvertMarkdownToEPUB(translatedMD, *outputFile); err != nil {
		log.Fatalf("Failed to convert markdown to EPUB: %v", err)
	}
	fmt.Printf("✓ Final EPUB created: %s\n\n", *outputFile)

	// Cleanup markdown files if requested
	if !*keepMarkdown {
		fmt.Println("🧹 Cleaning up intermediate files...")
		os.Remove(sourceMD)
		os.Remove(translatedMD)
		fmt.Println("✓ Cleanup complete\n")
	}

	fmt.Println("✅ Translation complete!")
	fmt.Printf("\nFiles generated:\n")
	if *keepMarkdown {
		fmt.Printf("  - Source MD:      %s\n", sourceMD)
		fmt.Printf("  - Translated MD:  %s\n", translatedMD)
	}
	fmt.Printf("  - Final EPUB:     %s\n", *outputFile)
}

// createTranslator creates an LLM translator based on provider
func createTranslator(provider, model, targetLang string) (translator.Translator, error) {
	// Get API keys from environment
	var apiKey string
	var defaultModel string

	switch provider {
	case "deepseek":
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
		defaultModel = "deepseek-chat"
	case "openai":
		apiKey = os.Getenv("OPENAI_API_KEY")
		defaultModel = "gpt-4"
	case "anthropic":
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
		defaultModel = "claude-3-sonnet-20240229"
	case "zhipu":
		apiKey = os.Getenv("ZHIPU_API_KEY")
		defaultModel = "glm-4"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key not set for provider %s (check environment variables)", provider)
	}

	// Use default model if not specified
	if model == "" {
		model = defaultModel
	}

	// Create translator config
	config := translator.TranslationConfig{
		SourceLang: "",
		TargetLang: targetLang,
		Provider:   provider,
		Model:      model,
		APIKey:     apiKey,
	}

	// Create LLM translator (it handles all providers internally)
	return llm.NewLLMTranslator(config)
}
