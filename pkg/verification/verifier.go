package verification

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"digital.vasic.translator/pkg/ebook"
	"digital.vasic.translator/pkg/events"
	"digital.vasic.translator/pkg/language"
)

// VerificationResult represents the result of content verification
type VerificationResult struct {
	IsValid            bool
	UntranslatedBlocks []UntranslatedBlock
	HTMLArtifacts      []HTMLArtifact
	QualityScore       float64
	Warnings           []string
	Errors             []string
}

// UntranslatedBlock represents a piece of content that wasn't translated
type UntranslatedBlock struct {
	Location    string // e.g., "Chapter 5, Section 2, Paragraph 3"
	OriginalText string
	Language     string
	Length       int
}

// HTMLArtifact represents HTML/XML found in translated content
type HTMLArtifact struct {
	Location string
	Content  string
	Type     string // "tag", "entity", "attribute"
}

// Verifier validates translation quality
type Verifier struct {
	sourceLanguage language.Language
	targetLanguage language.Language
	eventBus       *events.EventBus
	sessionID      string
}

// NewVerifier creates a new content verifier
func NewVerifier(
	sourceLanguage, targetLanguage language.Language,
	eventBus *events.EventBus,
	sessionID string,
) *Verifier {
	return &Verifier{
		sourceLanguage: sourceLanguage,
		targetLanguage: targetLanguage,
		eventBus:       eventBus,
		sessionID:      sessionID,
	}
}

// VerifyBook performs comprehensive verification of translated book
func (v *Verifier) VerifyBook(ctx context.Context, book *ebook.Book) (*VerificationResult, error) {
	result := &VerificationResult{
		IsValid:            true,
		UntranslatedBlocks: make([]UntranslatedBlock, 0),
		HTMLArtifacts:      make([]HTMLArtifact, 0),
		Warnings:           make([]string, 0),
		Errors:             make([]string, 0),
	}

	// Emit verification start event
	v.emitEvent(events.Event{
		Type:      "verification_started",
		SessionID: v.sessionID,
		Message:   "Starting translation verification",
	})

	// Verify metadata
	if err := v.verifyMetadata(&book.Metadata, result); err != nil {
		return result, err
	}

	// Verify chapters
	totalChapters := len(book.Chapters)
	for i := range book.Chapters {
		location := fmt.Sprintf("Chapter %d/%d", i+1, totalChapters)

		v.emitEvent(events.Event{
			Type:      "verification_progress",
			SessionID: v.sessionID,
			Message:   fmt.Sprintf("Verifying %s", location),
			Data: map[string]interface{}{
				"chapter":        i + 1,
				"total_chapters": totalChapters,
				"progress":       float64(i+1) / float64(totalChapters) * 100,
			},
		})

		if err := v.verifyChapter(&book.Chapters[i], i+1, result); err != nil {
			return result, err
		}
	}

	// Calculate quality score
	result.QualityScore = v.calculateQualityScore(result, book)

	// Determine if valid
	result.IsValid = len(result.Errors) == 0 && result.QualityScore >= 0.95

	// Emit completion event
	completionEvent := events.NewEvent(
		"verification_completed",
		fmt.Sprintf("Verification completed - Score: %.2f%%", result.QualityScore*100),
		map[string]interface{}{
			"quality_score":       result.QualityScore,
			"is_valid":            result.IsValid,
			"untranslated_blocks": len(result.UntranslatedBlocks),
			"html_artifacts":      len(result.HTMLArtifacts),
			"warnings":            len(result.Warnings),
			"errors":              len(result.Errors),
		},
	)
	completionEvent.SessionID = v.sessionID
	v.emitEvent(completionEvent)

	// Emit warnings for untranslated content
	if len(result.UntranslatedBlocks) > 0 {
		v.emitWarning(fmt.Sprintf("Found %d untranslated blocks", len(result.UntranslatedBlocks)))
		for i, block := range result.UntranslatedBlocks {
			if i < 10 { // Limit to first 10 warnings
				v.emitWarning(fmt.Sprintf("Untranslated: %s - %s", block.Location, truncate(block.OriginalText, 100)))
			}
		}
	}

	// Emit warnings for HTML artifacts
	if len(result.HTMLArtifacts) > 0 {
		v.emitWarning(fmt.Sprintf("Found %d HTML artifacts in translation", len(result.HTMLArtifacts)))
		for i, artifact := range result.HTMLArtifacts {
			if i < 10 { // Limit to first 10 warnings
				v.emitWarning(fmt.Sprintf("HTML in %s: %s", artifact.Location, artifact.Content))
			}
		}
	}

	return result, nil
}

// verifyMetadata checks if metadata is properly translated
func (v *Verifier) verifyMetadata(metadata *ebook.Metadata, result *VerificationResult) error {
	if metadata.Title != "" {
		if v.isSourceLanguage(metadata.Title) {
			result.UntranslatedBlocks = append(result.UntranslatedBlocks, UntranslatedBlock{
				Location:    "Book Title",
				OriginalText: metadata.Title,
				Language:     v.sourceLanguage.Code,
				Length:       len(metadata.Title),
			})
			result.Errors = append(result.Errors, "Book title not translated")
		}
	}

	if metadata.Description != "" {
		if v.isSourceLanguage(metadata.Description) {
			result.UntranslatedBlocks = append(result.UntranslatedBlocks, UntranslatedBlock{
				Location:    "Book Description",
				OriginalText: truncate(metadata.Description, 200),
				Language:     v.sourceLanguage.Code,
				Length:       len(metadata.Description),
			})
			result.Warnings = append(result.Warnings, "Book description not translated")
		}
	}

	return nil
}

// verifyChapter checks if chapter is properly translated
func (v *Verifier) verifyChapter(chapter *ebook.Chapter, chapterNum int, result *VerificationResult) error {
	location := fmt.Sprintf("Chapter %d", chapterNum)

	// Verify chapter title
	if chapter.Title != "" {
		if v.isSourceLanguage(chapter.Title) {
			result.UntranslatedBlocks = append(result.UntranslatedBlocks, UntranslatedBlock{
				Location:    location + " - Title",
				OriginalText: chapter.Title,
				Language:     v.sourceLanguage.Code,
				Length:       len(chapter.Title),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("%s title not translated", location))
		}
	}

	// Verify sections
	for i := range chapter.Sections {
		sectionLoc := fmt.Sprintf("%s, Section %d", location, i+1)
		if err := v.verifySection(&chapter.Sections[i], sectionLoc, result); err != nil {
			return err
		}
	}

	return nil
}

// verifySection checks if section is properly translated
func (v *Verifier) verifySection(section *ebook.Section, location string, result *VerificationResult) error {
	// Verify section title
	if section.Title != "" {
		if v.isSourceLanguage(section.Title) {
			result.UntranslatedBlocks = append(result.UntranslatedBlocks, UntranslatedBlock{
				Location:    location + " - Title",
				OriginalText: section.Title,
				Language:     v.sourceLanguage.Code,
				Length:       len(section.Title),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("%s title not translated", location))
		}
	}

	// Verify section content
	if section.Content != "" {
		// Check if content is translated
		if v.isSourceLanguage(section.Content) {
			result.UntranslatedBlocks = append(result.UntranslatedBlocks, UntranslatedBlock{
				Location:    location + " - Content",
				OriginalText: truncate(section.Content, 500),
				Language:     v.sourceLanguage.Code,
				Length:       len(section.Content),
			})
			result.Errors = append(result.Errors, fmt.Sprintf("%s content not translated", location))
		}

		// Check for HTML artifacts
		htmlArtifacts := v.detectHTMLArtifacts(section.Content)
		for _, artifact := range htmlArtifacts {
			artifact.Location = location
			result.HTMLArtifacts = append(result.HTMLArtifacts, artifact)
			result.Warnings = append(result.Warnings, fmt.Sprintf("HTML artifact in %s: %s", location, artifact.Content))
		}

		// Verify paragraphs
		paragraphs := v.splitIntoParagraphs(section.Content)
		for pi, para := range paragraphs {
			if v.isSourceLanguage(para) {
				paraLoc := fmt.Sprintf("%s, Paragraph %d", location, pi+1)
				result.UntranslatedBlocks = append(result.UntranslatedBlocks, UntranslatedBlock{
					Location:    paraLoc,
					OriginalText: truncate(para, 200),
					Language:     v.sourceLanguage.Code,
					Length:       len(para),
				})
			}
		}
	}

	// Verify subsections recursively
	for i := range section.Subsections {
		subLoc := fmt.Sprintf("%s, Subsection %d", location, i+1)
		if err := v.verifySection(&section.Subsections[i], subLoc, result); err != nil {
			return err
		}
	}

	return nil
}

// isSourceLanguage detects if text is in source language (not translated)
func (v *Verifier) isSourceLanguage(text string) bool {
	if text == "" {
		return false
	}

	// Clean text
	cleanText := strings.TrimSpace(text)

	// For Cyrillic-to-Cyrillic (e.g., Russian to Serbian), check specific characters
	// This check doesn't need minimum length as finding even one Russian-specific char is conclusive
	if v.sourceLanguage.Code == "ru" && v.targetLanguage.Code == "sr" {
		// Russian-specific letters that don't exist in Serbian
		russianOnlyChars := []rune{'ы', 'э', 'Ы', 'Э'}
		for _, char := range cleanText {
			for _, rusChar := range russianOnlyChars {
				if char == rusChar {
					return true // Definitely Russian
				}
			}
		}
	}

	// Check script - if source is Cyrillic and target is Latin (or vice versa)
	hasCyrillic := false
	hasLatin := false
	charCount := 0

	for _, r := range cleanText {
		if unicode.IsLetter(r) {
			charCount++
			if unicode.Is(unicode.Cyrillic, r) {
				hasCyrillic = true
			} else if unicode.Is(unicode.Latin, r) {
				hasLatin = true
			}
		}
	}

	if charCount < 10 {
		return false // Too few letters
	}

	// If we expect Cyrillic but got Latin, or vice versa
	targetCyrillic := v.targetLanguage.Code == "sr" || v.targetLanguage.Code == "ru" ||
	                   v.targetLanguage.Code == "bg" || v.targetLanguage.Code == "uk"
	sourceCyrillic := v.sourceLanguage.Code == "ru" || v.sourceLanguage.Code == "sr" ||
	                   v.sourceLanguage.Code == "bg" || v.sourceLanguage.Code == "uk"

	if sourceCyrillic && !targetCyrillic {
		// Source is Cyrillic, target is not - if we have Cyrillic, not translated
		return hasCyrillic
	}

	if !sourceCyrillic && targetCyrillic {
		// Source is Latin, target is Cyrillic - if we have Latin, not translated
		return hasLatin
	}

	// Default: assume if mostly Cyrillic and source is Cyrillic, might be untranslated
	// This is a heuristic and may need refinement
	return false
}

// detectHTMLArtifacts finds HTML/XML tags in content
func (v *Verifier) detectHTMLArtifacts(content string) []HTMLArtifact {
	artifacts := make([]HTMLArtifact, 0)

	// Regex patterns for HTML detection
	tagPattern := regexp.MustCompile(`<[^>]+>`)
	entityPattern := regexp.MustCompile(`&[a-zA-Z]+;|&#[0-9]+;`)

	// Find HTML tags
	tags := tagPattern.FindAllString(content, -1)
	for _, tag := range tags {
		// Skip common allowed tags if any
		if !strings.Contains(tag, "<!") && !strings.Contains(tag, "<?") {
			artifacts = append(artifacts, HTMLArtifact{
				Content: tag,
				Type:    "tag",
			})
		}
	}

	// Find HTML entities
	entities := entityPattern.FindAllString(content, -1)
	for _, entity := range entities {
		artifacts = append(artifacts, HTMLArtifact{
			Content: entity,
			Type:    "entity",
		})
	}

	return artifacts
}

// splitIntoParagraphs splits content into paragraphs
func (v *Verifier) splitIntoParagraphs(content string) []string {
	// Split by double newlines or paragraph breaks
	paragraphs := regexp.MustCompile(`\n\n+`).Split(content, -1)
	result := make([]string, 0, len(paragraphs))

	for _, para := range paragraphs {
		cleaned := strings.TrimSpace(para)
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}

	return result
}

// calculateQualityScore computes overall translation quality
func (v *Verifier) calculateQualityScore(result *VerificationResult, book *ebook.Book) float64 {
	// Count total translatable items
	totalItems := 0
	totalChars := 0

	// Count book elements
	if book.Metadata.Title != "" {
		totalItems++
		totalChars += len(book.Metadata.Title)
	}
	if book.Metadata.Description != "" {
		totalItems++
		totalChars += len(book.Metadata.Description)
	}

	for _, chapter := range book.Chapters {
		if chapter.Title != "" {
			totalItems++
			totalChars += len(chapter.Title)
		}
		totalItems += v.countSectionItems(&chapter.Sections, &totalChars)
	}

	if totalItems == 0 || totalChars == 0 {
		return 0.0
	}

	// Calculate untranslated character count
	untranslatedChars := 0
	for _, block := range result.UntranslatedBlocks {
		untranslatedChars += block.Length
	}

	// Calculate character-based quality score
	charScore := 1.0 - (float64(untranslatedChars) / float64(totalChars))

	// Penalize for HTML artifacts
	htmlPenalty := float64(len(result.HTMLArtifacts)) * 0.01
	if htmlPenalty > 0.1 {
		htmlPenalty = 0.1 // Cap at 10% penalty
	}

	// Penalize for errors more than warnings
	errorPenalty := float64(len(result.Errors)) * 0.05
	if errorPenalty > 0.3 {
		errorPenalty = 0.3 // Cap at 30% penalty
	}

	finalScore := charScore - htmlPenalty - errorPenalty
	if finalScore < 0 {
		finalScore = 0
	}

	return finalScore
}

// countSectionItems recursively counts sections for quality calculation
func (v *Verifier) countSectionItems(sections *[]ebook.Section, totalChars *int) int {
	count := 0
	for i := range *sections {
		section := &(*sections)[i]
		if section.Title != "" {
			count++
			*totalChars += len(section.Title)
		}
		if section.Content != "" {
			count++
			*totalChars += len(section.Content)
		}
		count += v.countSectionItems(&section.Subsections, totalChars)
	}
	return count
}

// emitEvent emits a verification event
func (v *Verifier) emitEvent(event events.Event) {
	if v.eventBus != nil {
		v.eventBus.Publish(event)
	}
}

// emitWarning emits a warning event
func (v *Verifier) emitWarning(message string) {
	if v.eventBus != nil {
		warningEvent := events.NewEvent("verification_warning", message, nil)
		warningEvent.SessionID = v.sessionID
		v.eventBus.Publish(warningEvent)
	}
}

// truncate truncates text to specified length
func truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
