package progress

import (
	"strconv"
	"sync"
	"time"
)

// TranslationProgress tracks detailed translation progress
type TranslationProgress struct {
	mu sync.RWMutex

	// Book information
	BookTitle       string    `json:"book_title"`
	TotalChapters   int       `json:"total_chapters"`
	CurrentChapter  int       `json:"current_chapter"`
	ChapterTitle    string    `json:"chapter_title"`
	CurrentSection  int       `json:"current_section"`
	TotalSections   int       `json:"total_sections"`

	// Progress metrics
	PercentComplete float64   `json:"percent_complete"`
	ItemsTotal      int       `json:"items_total"`
	ItemsCompleted  int       `json:"items_completed"`
	ItemsFailed     int       `json:"items_failed"`

	// Time tracking
	StartTime       time.Time `json:"start_time"`
	ElapsedTime     string    `json:"elapsed_time"`
	EstimatedETA    string    `json:"estimated_eta"`

	// Translation details
	SourceLanguage  string    `json:"source_language"`
	TargetLanguage  string    `json:"target_language"`
	Provider        string    `json:"provider"`
	Model           string    `json:"model"`

	// Status
	Status          string    `json:"status"` // "initializing", "translating", "completed", "error"
	CurrentTask     string    `json:"current_task"`
	SessionID       string    `json:"session_id"`
}

// Tracker manages translation progress
type Tracker struct {
	mu       sync.RWMutex
	progress *TranslationProgress
}

// NewTracker creates a new progress tracker
func NewTracker(sessionID, bookTitle string, totalChapters int, sourceLanguage, targetLanguage, provider, model string) *Tracker {
	return &Tracker{
		progress: &TranslationProgress{
			SessionID:      sessionID,
			BookTitle:      bookTitle,
			TotalChapters:  totalChapters,
			SourceLanguage: sourceLanguage,
			TargetLanguage: targetLanguage,
			Provider:       provider,
			Model:          model,
			StartTime:      time.Now(),
			Status:         "initializing",
		},
	}
}

// UpdateChapter updates the current chapter being translated
func (t *Tracker) UpdateChapter(chapterNum int, chapterTitle string, totalSections int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.CurrentChapter = chapterNum
	t.progress.ChapterTitle = chapterTitle
	t.progress.TotalSections = totalSections
	t.progress.CurrentSection = 0
	t.progress.Status = "translating"
	t.progress.CurrentTask = "Translating chapter " + chapterTitle

	t.updateProgress()
}

// UpdateSection updates the current section being translated
func (t *Tracker) UpdateSection(sectionNum int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.CurrentSection = sectionNum
	t.updateProgress()
}

// IncrementCompleted increments the completed items counter
func (t *Tracker) IncrementCompleted() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.ItemsCompleted++
	t.updateProgress()
}

// IncrementFailed increments the failed items counter
func (t *Tracker) IncrementFailed() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.ItemsFailed++
	t.updateProgress()
}

// SetTotal sets the total number of items to translate
func (t *Tracker) SetTotal(total int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.ItemsTotal = total
}

// SetStatus updates the status
func (t *Tracker) SetStatus(status, task string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.Status = status
	t.progress.CurrentTask = task
}

// Complete marks the translation as completed
func (t *Tracker) Complete() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.Status = "completed"
	t.progress.CurrentTask = "Translation completed"
	t.progress.PercentComplete = 100.0
	t.updateProgress()
}

// Error marks the translation as errored
func (t *Tracker) Error(errorMsg string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.progress.Status = "error"
	t.progress.CurrentTask = "Error: " + errorMsg
}

// GetProgress returns a copy of the current progress
func (t *Tracker) GetProgress() TranslationProgress {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Update time fields before returning
	elapsed := time.Since(t.progress.StartTime)
	t.progress.ElapsedTime = formatDuration(elapsed)

	if t.progress.ItemsCompleted > 0 && t.progress.ItemsTotal > 0 {
		avgTimePerItem := elapsed / time.Duration(t.progress.ItemsCompleted)
		remainingItems := t.progress.ItemsTotal - t.progress.ItemsCompleted
		estimatedRemaining := avgTimePerItem * time.Duration(remainingItems)
		t.progress.EstimatedETA = formatDuration(estimatedRemaining)
	}

	return *t.progress
}

// updateProgress calculates percentage and updates progress (must be called with lock held)
func (t *Tracker) updateProgress() {
	if t.progress.TotalChapters > 0 {
		t.progress.PercentComplete = float64(t.progress.CurrentChapter-1) / float64(t.progress.TotalChapters) * 100.0

		// Add section progress within current chapter
		if t.progress.TotalSections > 0 {
			sectionPercent := float64(t.progress.CurrentSection) / float64(t.progress.TotalSections) / float64(t.progress.TotalChapters) * 100.0
			t.progress.PercentComplete += sectionPercent
		}

		// Cap at 100%
		if t.progress.PercentComplete > 100.0 {
			t.progress.PercentComplete = 100.0
		}
	}

	// Update elapsed time
	elapsed := time.Since(t.progress.StartTime)
	t.progress.ElapsedTime = formatDuration(elapsed)

	// Calculate ETA
	if t.progress.PercentComplete > 0 && t.progress.PercentComplete < 100 {
		totalEstimated := elapsed / time.Duration(t.progress.PercentComplete) * 100
		remaining := totalEstimated - elapsed
		t.progress.EstimatedETA = formatDuration(remaining)
	} else if t.progress.PercentComplete >= 100 {
		t.progress.EstimatedETA = "Completed"
	} else {
		t.progress.EstimatedETA = "Calculating..."
	}
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return formatTime(hours, "hour") + " " + formatTime(minutes, "minute")
	} else if minutes > 0 {
		return formatTime(minutes, "minute") + " " + formatTime(seconds, "second")
	} else {
		return formatTime(seconds, "second")
	}
}

// formatTime formats a time value with proper singular/plural
func formatTime(value int, unit string) string {
	if value == 0 {
		return ""
	}
	if value == 1 {
		return "1 " + unit
	}
	return formatInt(value) + " " + unit + "s"
}

// formatInt formats an integer
func formatInt(n int) string {
	return strconv.Itoa(n)
}
