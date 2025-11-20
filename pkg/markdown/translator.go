package markdown

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// MarkdownTranslator translates markdown content while preserving formatting
type MarkdownTranslator struct {
	translateFunc func(text string) (string, error)
}

// NewMarkdownTranslator creates a new markdown translator
func NewMarkdownTranslator(translateFunc func(string) (string, error)) *MarkdownTranslator {
	return &MarkdownTranslator{
		translateFunc: translateFunc,
	}
}

// TranslateMarkdownFile translates a markdown file
func (mt *MarkdownTranslator) TranslateMarkdownFile(inputPath, outputPath string) error {
	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Translate content
	translated, err := mt.TranslateMarkdown(string(content))
	if err != nil {
		return fmt.Errorf("translation failed: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, []byte(translated), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// TranslateMarkdown translates markdown content while preserving formatting
func (mt *MarkdownTranslator) TranslateMarkdown(content string) (string, error) {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))

	inFrontmatter := false
	inCodeBlock := false
	frontmatterCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Handle frontmatter (YAML between ---)
		if line == "---" {
			frontmatterCount++
			result.WriteString(line + "\n")
			if frontmatterCount == 1 {
				inFrontmatter = true
			} else if frontmatterCount == 2 {
				inFrontmatter = false
			}
			continue
		}

		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			result.WriteString(line + "\n")
			continue
		}

		// Don't translate frontmatter or code blocks
		if inFrontmatter || inCodeBlock {
			result.WriteString(line + "\n")
			continue
		}

		// Translate the line while preserving markdown syntax
		translated, err := mt.translateLine(line)
		if err != nil {
			return "", fmt.Errorf("failed to translate line: %w", err)
		}

		result.WriteString(translated + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading content: %w", err)
	}

	return result.String(), nil
}

// translateLine translates a single line while preserving markdown formatting
func (mt *MarkdownTranslator) translateLine(line string) (string, error) {
	// Empty lines
	if strings.TrimSpace(line) == "" {
		return line, nil
	}

	// Headers
	if strings.HasPrefix(line, "#") {
		return mt.translateHeader(line)
	}

	// Horizontal rules
	if matched, _ := regexp.MatchString(`^[-*_]{3,}$`, strings.TrimSpace(line)); matched {
		return line, nil
	}

	// Lists
	if matched, _ := regexp.MatchString(`^\s*[-*+]\s`, line); matched {
		return mt.translateList(line)
	}

	// Numbered lists
	if matched, _ := regexp.MatchString(`^\s*\d+\.\s`, line); matched {
		return mt.translateNumberedList(line)
	}

	// Blockquotes
	if strings.HasPrefix(strings.TrimSpace(line), ">") {
		return mt.translateBlockquote(line)
	}

	// Image references (don't translate alt text for now, keep original)
	if matched, _ := regexp.MatchString(`!\[.*?\]\(.*?\)`, line); matched {
		return line, nil
	}

	// Regular paragraph with inline formatting
	return mt.translateInlineFormatting(line)
}

// translateHeader translates a header line
func (mt *MarkdownTranslator) translateHeader(line string) (string, error) {
	// Extract header level and text
	match := regexp.MustCompile(`^(#{1,6})\s+(.+)$`).FindStringSubmatch(line)
	if len(match) != 3 {
		return line, nil
	}

	headerLevel := match[1]
	headerText := match[2]

	// Translate header text (preserve inline formatting)
	translated, err := mt.translateInlineFormatting(headerText)
	if err != nil {
		return "", err
	}

	return headerLevel + " " + translated, nil
}

// translateList translates a list item
func (mt *MarkdownTranslator) translateList(line string) (string, error) {
	// Extract indentation, bullet, and text
	match := regexp.MustCompile(`^(\s*)([-*+])\s+(.+)$`).FindStringSubmatch(line)
	if len(match) != 4 {
		return line, nil
	}

	indent := match[1]
	bullet := match[2]
	text := match[3]

	// Translate text
	translated, err := mt.translateInlineFormatting(text)
	if err != nil {
		return "", err
	}

	return indent + bullet + " " + translated, nil
}

// translateNumberedList translates a numbered list item
func (mt *MarkdownTranslator) translateNumberedList(line string) (string, error) {
	// Extract indentation, number, and text
	match := regexp.MustCompile(`^(\s*)(\d+)\.\s+(.+)$`).FindStringSubmatch(line)
	if len(match) != 4 {
		return line, nil
	}

	indent := match[1]
	number := match[2]
	text := match[3]

	// Translate text
	translated, err := mt.translateInlineFormatting(text)
	if err != nil {
		return "", err
	}

	return indent + number + ". " + translated, nil
}

// translateBlockquote translates a blockquote
func (mt *MarkdownTranslator) translateBlockquote(line string) (string, error) {
	// Extract quote marker and text
	match := regexp.MustCompile(`^(\s*>+)\s*(.*)$`).FindStringSubmatch(line)
	if len(match) != 3 {
		return line, nil
	}

	quoteMarker := match[1]
	text := match[2]

	if text == "" {
		return line, nil
	}

	// Translate text
	translated, err := mt.translateInlineFormatting(text)
	if err != nil {
		return "", err
	}

	return quoteMarker + " " + translated, nil
}

// translateInlineFormatting translates text while preserving inline markdown formatting
func (mt *MarkdownTranslator) translateInlineFormatting(text string) (string, error) {
	// Pattern to match markdown inline formatting
	// Matches: **bold**, *italic*, `code`, [link](url), etc.
	pattern := regexp.MustCompile(`(\*\*.*?\*\*|\*.*?\*|__.*?__|_.*?_|` + "`" + `.*?` + "`" + `|\[.*?\]\(.*?\))`)

	// Find all formatted segments
	segments := pattern.FindAllStringIndex(text, -1)

	if len(segments) == 0 {
		// No formatting, translate entire text
		return mt.translateText(text)
	}

	// Build result preserving formatting
	var result strings.Builder
	lastEnd := 0

	for _, seg := range segments {
		start, end := seg[0], seg[1]

		// Translate text before formatted segment
		if start > lastEnd {
			plainText := text[lastEnd:start]
			if strings.TrimSpace(plainText) != "" {
				translated, err := mt.translateText(plainText)
				if err != nil {
					return "", err
				}
				result.WriteString(translated)
			}
		}

		// Handle formatted segment
		formatted := text[start:end]
		translatedFormatted, err := mt.translateFormattedSegment(formatted)
		if err != nil {
			return "", err
		}
		result.WriteString(translatedFormatted)

		lastEnd = end
	}

	// Translate remaining text
	if lastEnd < len(text) {
		plainText := text[lastEnd:]
		if strings.TrimSpace(plainText) != "" {
			translated, err := mt.translateText(plainText)
			if err != nil {
				return "", err
			}
			result.WriteString(translated)
		}
	}

	return result.String(), nil
}

// translateFormattedSegment translates a formatted markdown segment
func (mt *MarkdownTranslator) translateFormattedSegment(segment string) (string, error) {
	// Bold: **text** or __text__
	if strings.HasPrefix(segment, "**") && strings.HasSuffix(segment, "**") {
		inner := segment[2 : len(segment)-2]
		translated, err := mt.translateText(inner)
		if err != nil {
			return "", err
		}
		return "**" + translated + "**", nil
	}

	if strings.HasPrefix(segment, "__") && strings.HasSuffix(segment, "__") {
		inner := segment[2 : len(segment)-2]
		translated, err := mt.translateText(inner)
		if err != nil {
			return "", err
		}
		return "__" + translated + "__", nil
	}

	// Italic: *text* or _text_
	if strings.HasPrefix(segment, "*") && strings.HasSuffix(segment, "*") && !strings.HasPrefix(segment, "**") {
		inner := segment[1 : len(segment)-1]
		translated, err := mt.translateText(inner)
		if err != nil {
			return "", err
		}
		return "*" + translated + "*", nil
	}

	if strings.HasPrefix(segment, "_") && strings.HasSuffix(segment, "_") && !strings.HasPrefix(segment, "__") {
		inner := segment[1 : len(segment)-1]
		translated, err := mt.translateText(inner)
		if err != nil {
			return "", err
		}
		return "_" + translated + "_", nil
	}

	// Code: `text`
	if strings.HasPrefix(segment, "`") && strings.HasSuffix(segment, "`") {
		// Don't translate code
		return segment, nil
	}

	// Links: [text](url)
	linkPattern := regexp.MustCompile(`^\[(.*?)\]\((.*?)\)$`)
	if match := linkPattern.FindStringSubmatch(segment); len(match) == 3 {
		linkText := match[1]
		linkURL := match[2]

		// Translate link text only
		translated, err := mt.translateText(linkText)
		if err != nil {
			return "", err
		}
		return "[" + translated + "](" + linkURL + ")", nil
	}

	// Unknown format, translate as-is
	return mt.translateText(segment)
}

// translateText translates plain text using the provided translation function
func (mt *MarkdownTranslator) translateText(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return text, nil
	}

	translated, err := mt.translateFunc(text)
	if err != nil {
		return "", fmt.Errorf("translation error: %w", err)
	}

	return translated, nil
}
