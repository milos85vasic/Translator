package verification

import (
	"context"
	"digital.vasic.translator/pkg/translator/llm"
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// NoteType represents different types of literary notes
type NoteType string

const (
	NoteTypeCharacter  NoteType = "character"  // Character development, traits, arcs
	NoteTypeTone       NoteType = "tone"       // Tone, atmosphere, mood
	NoteTypeTheme      NoteType = "theme"      // Themes, motifs, symbols
	NoteTypeCulture    NoteType = "culture"    // Cultural references, idioms
	NoteTypeStyle      NoteType = "style"      // Literary style, techniques
	NoteTypeContext    NoteType = "context"    // Historical, social context
	NoteTypeVocabulary NoteType = "vocabulary" // Key terms, specialized vocabulary
	NoteTypeStructure  NoteType = "structure"  // Narrative structure, pacing
)

// ImportanceLevel represents the importance of a note
type ImportanceLevel string

const (
	ImportanceCritical ImportanceLevel = "critical" // Must preserve exactly
	ImportanceHigh     ImportanceLevel = "high"     // Very important to preserve
	ImportanceMedium   ImportanceLevel = "medium"   // Important context
	ImportanceLow      ImportanceLevel = "low"      // Minor observation
)

// LiteraryNote represents an observation about the text
type LiteraryNote struct {
	ID           string          `json:"id"`
	PassNumber   int             `json:"pass_number"`
	SectionID    string          `json:"section_id"`
	Location     string          `json:"location"`
	Provider     string          `json:"provider"`
	NoteType     NoteType        `json:"note_type"`
	Importance   ImportanceLevel `json:"importance"`
	Title        string          `json:"title"`
	Content      string          `json:"content"`
	Examples     []string        `json:"examples"`
	Implications string          `json:"implications"`
	CreatedAt    time.Time       `json:"created_at"`
}

// NoteCollection groups notes for efficient access
type NoteCollection struct {
	ByType     map[NoteType][]*LiteraryNote
	BySection  map[string][]*LiteraryNote
	ByProvider map[string][]*LiteraryNote
	ByPass     map[int][]*LiteraryNote
	All        []*LiteraryNote
}

// NewNoteCollection creates a new note collection
func NewNoteCollection() *NoteCollection {
	return &NoteCollection{
		ByType:     make(map[NoteType][]*LiteraryNote),
		BySection:  make(map[string][]*LiteraryNote),
		ByProvider: make(map[string][]*LiteraryNote),
		ByPass:     make(map[int][]*LiteraryNote),
		All:        make([]*LiteraryNote, 0),
	}
}

// Add adds a note to the collection
func (nc *NoteCollection) Add(note *LiteraryNote) {
	nc.All = append(nc.All, note)
	nc.ByType[note.NoteType] = append(nc.ByType[note.NoteType], note)
	nc.BySection[note.SectionID] = append(nc.BySection[note.SectionID], note)
	nc.ByProvider[note.Provider] = append(nc.ByProvider[note.Provider], note)
	nc.ByPass[note.PassNumber] = append(nc.ByPass[note.PassNumber], note)
}

// GetForSection retrieves all notes for a specific section
func (nc *NoteCollection) GetForSection(sectionID string) []*LiteraryNote {
	return nc.BySection[sectionID]
}

// GetByType retrieves all notes of a specific type
func (nc *NoteCollection) GetByType(noteType NoteType) []*LiteraryNote {
	return nc.ByType[noteType]
}

// GetCritical retrieves all critical notes
func (nc *NoteCollection) GetCritical() []*LiteraryNote {
	var critical []*LiteraryNote
	for _, note := range nc.All {
		if note.Importance == ImportanceCritical {
			critical = append(critical, note)
		}
	}
	return critical
}

// GetByPass retrieves all notes from a specific pass
func (nc *NoteCollection) GetByPass(passNumber int) []*LiteraryNote {
	return nc.ByPass[passNumber]
}

// Summary generates a text summary of the collection
func (nc *NoteCollection) Summary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Notes: %d\n\n", len(nc.All)))

	// By type
	sb.WriteString("By Type:\n")
	for noteType, notes := range nc.ByType {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", noteType, len(notes)))
	}
	sb.WriteString("\n")

	// By importance
	critical := nc.GetCritical()
	sb.WriteString(fmt.Sprintf("Critical Notes: %d\n", len(critical)))

	return sb.String()
}

// NoteTaker generates literary notes using LLMs
type NoteTaker struct {
	translator *llm.LLMTranslator
	provider   string
}

// NewNoteTaker creates a new note taker
func NewNoteTaker(translator *llm.LLMTranslator, provider string) *NoteTaker {
	return &NoteTaker{
		translator: translator,
		provider:   provider,
	}
}

// GenerateNotes generates literary notes for a text section
func (nt *NoteTaker) GenerateNotes(
	ctx context.Context,
	passNumber int,
	sectionID string,
	location string,
	originalText string,
	translatedText string,
	previousNotes []*LiteraryNote,
) ([]*LiteraryNote, error) {
	// Create note-taking prompt
	prompt := nt.createNotePrompt(originalText, translatedText, previousNotes)

	// Get LLM analysis
	response, err := nt.translator.Translate(ctx, prompt, location)
	if err != nil {
		return nil, fmt.Errorf("note generation failed: %w", err)
	}

	// Parse notes from response
	notes := nt.parseNotes(response, passNumber, sectionID, location)

	return notes, nil
}

// createNotePrompt creates the prompt for note generation
func (nt *NoteTaker) createNotePrompt(
	originalText string,
	translatedText string,
	previousNotes []*LiteraryNote,
) string {
	var sb strings.Builder

	sb.WriteString(`You are a literary analyst reviewing a translation. Generate detailed notes about important aspects that must be preserved or improved.

**Original Text:**
`)
	sb.WriteString(originalText)
	sb.WriteString("\n\n**Current Translation:**\n")
	sb.WriteString(translatedText)
	sb.WriteString("\n\n")

	// Include previous notes if available
	if len(previousNotes) > 0 {
		sb.WriteString("**Previous Analysis (from earlier pass):**\n")
		for _, note := range previousNotes {
			sb.WriteString(fmt.Sprintf("- [%s] %s: %s\n", note.NoteType, note.Title, note.Content))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(`**Generate notes for the following dimensions:**

1. **CHARACTER**: Character traits, development, voice, relationships
2. **TONE**: Atmosphere, mood, emotional tone, narrative voice
3. **THEME**: Themes, motifs, symbols, deeper meanings
4. **CULTURE**: Cultural references, idioms, historical context
5. **STYLE**: Literary techniques, sentence structure, rhythm
6. **VOCABULARY**: Key terms, specialized vocabulary, word choice significance

**Response Format:**

NOTE: [TYPE]
IMPORTANCE: [critical/high/medium/low]
TITLE: [Brief title]
CONTENT: [Detailed observation]
EXAMPLES: [Specific examples from text, one per line]
IMPLICATIONS: [Why this matters for translation]
---

Provide 3-10 notes covering different aspects. Focus on elements that are critical for translation quality.
`)

	return sb.String()
}

// parseNotes parses notes from LLM response
func (nt *NoteTaker) parseNotes(
	response string,
	passNumber int,
	sectionID string,
	location string,
) []*LiteraryNote {
	notes := make([]*LiteraryNote, 0)

	// Split by note separator
	noteSections := strings.Split(response, "---")

	for _, noteSection := range noteSections {
		noteSection = strings.TrimSpace(noteSection)
		if noteSection == "" {
			continue
		}

		note := nt.parseNote(noteSection, passNumber, sectionID, location)
		if note != nil {
			notes = append(notes, note)
		}
	}

	return notes
}

// parseNote parses a single note from text
func (nt *NoteTaker) parseNote(
	text string,
	passNumber int,
	sectionID string,
	location string,
) *LiteraryNote {
	note := &LiteraryNote{
		ID:         fmt.Sprintf("%s_%s_%d_%d", nt.provider, sectionID, passNumber, time.Now().UnixNano()),
		PassNumber: passNumber,
		SectionID:  sectionID,
		Location:   location,
		Provider:   nt.provider,
		Examples:   make([]string, 0),
		CreatedAt:  time.Now(),
	}

	lines := strings.Split(text, "\n")
	var currentField string
	var contentBuilder strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for field markers
		if strings.HasPrefix(line, "NOTE:") {
			noteType := strings.TrimSpace(strings.TrimPrefix(line, "NOTE:"))
			noteType = strings.Trim(noteType, "[]")
			note.NoteType = NoteType(strings.ToLower(noteType))
			currentField = "note_type"
		} else if strings.HasPrefix(line, "IMPORTANCE:") {
			importance := strings.TrimSpace(strings.TrimPrefix(line, "IMPORTANCE:"))
			importance = strings.Trim(importance, "[]")
			note.Importance = ImportanceLevel(strings.ToLower(importance))
			currentField = "importance"
		} else if strings.HasPrefix(line, "TITLE:") {
			note.Title = strings.TrimSpace(strings.TrimPrefix(line, "TITLE:"))
			note.Title = strings.Trim(note.Title, "[]")
			currentField = "title"
		} else if strings.HasPrefix(line, "CONTENT:") {
			contentBuilder.Reset()
			content := strings.TrimSpace(strings.TrimPrefix(line, "CONTENT:"))
			if content != "" {
				contentBuilder.WriteString(content)
			}
			currentField = "content"
		} else if strings.HasPrefix(line, "EXAMPLES:") {
			currentField = "examples"
		} else if strings.HasPrefix(line, "IMPLICATIONS:") {
			note.Content = contentBuilder.String()
			contentBuilder.Reset()
			implications := strings.TrimSpace(strings.TrimPrefix(line, "IMPLICATIONS:"))
			if implications != "" {
				contentBuilder.WriteString(implications)
			}
			currentField = "implications"
		} else {
			// Continuation of current field
			switch currentField {
			case "content":
				if contentBuilder.Len() > 0 {
					contentBuilder.WriteString(" ")
				}
				contentBuilder.WriteString(line)
			case "examples":
				if line != "" && !strings.HasPrefix(line, "IMPLICATIONS:") {
					note.Examples = append(note.Examples, line)
				}
			case "implications":
				if contentBuilder.Len() > 0 {
					contentBuilder.WriteString(" ")
				}
				contentBuilder.WriteString(line)
			}
		}
	}

	// Finalize implications
	if currentField == "implications" {
		note.Implications = contentBuilder.String()
	}

	// Validate note has minimum required fields
	if note.NoteType == "" || note.Title == "" || note.Content == "" {
		return nil
	}

	// Default importance
	if note.Importance == "" {
		note.Importance = ImportanceMedium
	}

	return note
}

// FormatNotesForContext formats notes for inclusion in polishing prompt
func FormatNotesForContext(notes []*LiteraryNote) string {
	if len(notes) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("**Previous Literary Analysis:**\n\n")

	// Group by type
	byType := make(map[NoteType][]*LiteraryNote)
	for _, note := range notes {
		byType[note.NoteType] = append(byType[note.NoteType], note)
	}

	// Format each type
	for noteType, typeNotes := range byType {
		sb.WriteString(fmt.Sprintf("### %s\n", cases.Title(language.English, cases.Compact).String(string(noteType))))

		for _, note := range typeNotes {
			importance := ""
			if note.Importance == ImportanceCritical {
				importance = " ⚠️"
			} else if note.Importance == ImportanceHigh {
				importance = " ⭐"
			}

			sb.WriteString(fmt.Sprintf("- **%s**%s: %s\n", note.Title, importance, note.Content))

			if len(note.Examples) > 0 {
				sb.WriteString("  Examples: ")
				sb.WriteString(strings.Join(note.Examples, "; "))
				sb.WriteString("\n")
			}

			if note.Implications != "" {
				sb.WriteString(fmt.Sprintf("  → %s\n", note.Implications))
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// FilterNotesByImportance filters notes by minimum importance level
func FilterNotesByImportance(notes []*LiteraryNote, minImportance ImportanceLevel) []*LiteraryNote {
	importanceOrder := map[ImportanceLevel]int{
		ImportanceLow:      1,
		ImportanceMedium:   2,
		ImportanceHigh:     3,
		ImportanceCritical: 4,
	}

	minLevel := importanceOrder[minImportance]
	filtered := make([]*LiteraryNote, 0)

	for _, note := range notes {
		if importanceOrder[note.Importance] >= minLevel {
			filtered = append(filtered, note)
		}
	}

	return filtered
}

// MergeNotes merges notes from multiple passes, deduplicating similar ones
func MergeNotes(allNotes []*LiteraryNote) []*LiteraryNote {
	// Simple deduplication based on title similarity
	// More sophisticated merging can be added later
	seen := make(map[string]*LiteraryNote)
	merged := make([]*LiteraryNote, 0)

	for _, note := range allNotes {
		key := fmt.Sprintf("%s:%s:%s", note.SectionID, note.NoteType, strings.ToLower(note.Title))

		if existing, found := seen[key]; found {
			// Merge: keep higher importance, append examples
			if importanceLevel(note.Importance) > importanceLevel(existing.Importance) {
				existing.Importance = note.Importance
			}

			// Append unique examples
			for _, example := range note.Examples {
				if !contains(existing.Examples, example) {
					existing.Examples = append(existing.Examples, example)
				}
			}

			// Append implications
			if note.Implications != "" && !strings.Contains(existing.Implications, note.Implications) {
				if existing.Implications != "" {
					existing.Implications += " "
				}
				existing.Implications += note.Implications
			}
		} else {
			seen[key] = note
			merged = append(merged, note)
		}
	}

	return merged
}

// Helper functions

func importanceLevel(importance ImportanceLevel) int {
	levels := map[ImportanceLevel]int{
		ImportanceLow:      1,
		ImportanceMedium:   2,
		ImportanceHigh:     3,
		ImportanceCritical: 4,
	}
	return levels[importance]
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
