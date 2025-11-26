package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type TranslationEvent struct {
	Type         string  `json:"type"`
	SessionID    string  `json:"session_id"`
	Step         string  `json:"step,omitempty"`
	Message      string  `json:"message,omitempty"`
	Progress     float64 `json:"progress,omitempty"`
	Error        string  `json:"error,omitempty"`
	CurrentItem  string  `json:"current_item,omitempty"`
	TotalItems   int     `json:"total_items,omitempty"`
	Timestamp    int64    `json:"timestamp"`
}

func main() {
	sessionID := "live-demo-" + fmt.Sprintf("%d", time.Now().Unix())
	inputFile := "test/fixtures/ebooks/russian_sample.txt"
	outputFile := "demo_translation_output.md"

	fmt.Printf("🚀 Starting Live Translation Demo with WebSocket Monitoring\n")
	fmt.Printf("📊 Session ID: %s\n", sessionID)
	fmt.Printf("📄 Input: %s\n", inputFile)
	fmt.Printf("📄 Output: %s\n", outputFile)
	fmt.Printf("🔗 WebSocket: ws://localhost:8090/ws?session_id=%s\n\n", sessionID)

	// Connect to WebSocket monitoring
	ws, err := connectWebSocket(sessionID)
	if err != nil {
		log.Printf("Warning: Could not connect to monitoring: %v", err)
	} else {
		defer ws.Close()
		fmt.Printf("✅ Connected to monitoring server\n")
		
		// Start listening for messages in background
		go listenForWebSocketEvents(ws)
	}

	// Emit translation started event
	emitEvent(ws, TranslationEvent{
		Type:      "translation_started",
		SessionID: sessionID,
		Message:   "Translation job started",
		Progress:  0,
		Timestamp: time.Now().Unix(),
	})

	// Read input file
	emitEvent(ws, TranslationEvent{
		Type:       "translation_progress",
		SessionID:  sessionID,
		Step:       "reading",
		Message:    "Reading input file...",
		Progress:   5,
		Timestamp:  time.Now().Unix(),
	})

	content, err := os.ReadFile(inputFile)
	if err != nil {
		emitErrorEvent(ws, sessionID, "reading", fmt.Sprintf("Failed to read input file: %v", err))
		log.Fatalf("Failed to read input file: %v", err)
	}

	text := string(content)
	lines := strings.Split(text, "\n")
	totalLines := len(lines)

	emitEvent(ws, TranslationEvent{
		Type:        "step_completed",
		SessionID:   sessionID,
		Step:        "reading",
		Message:     fmt.Sprintf("File read successfully (%d lines)", totalLines),
		Progress:    10,
		CurrentItem: "file_read",
		TotalItems:  totalLines,
		Timestamp:   time.Now().Unix(),
	})

	// Parse and prepare content
	emitEvent(ws, TranslationEvent{
		Type:       "translation_progress",
		SessionID:  sessionID,
		Step:       "parsing",
		Message:    "Parsing content structure...",
		Progress:   15,
		Timestamp:  time.Now().Unix(),
	})

	time.Sleep(1 * time.Second) // Simulate processing time

	emitEvent(ws, TranslationEvent{
		Type:        "step_completed",
		SessionID:   sessionID,
		Step:        "parsing",
		Message:     "Content parsed successfully",
		Progress:    20,
		CurrentItem: "content_parsed",
		TotalItems:  totalLines,
		Timestamp:   time.Now().Unix(),
	})

	// Convert to markdown
	emitEvent(ws, TranslationEvent{
		Type:       "translation_progress",
		SessionID:  sessionID,
		Step:       "conversion",
		Message:    "Converting to markdown format...",
		Progress:   25,
		Timestamp:  time.Now().Unix(),
	})

	markdownContent := convertToMarkdown(text)
	time.Sleep(1 * time.Second)

	emitEvent(ws, TranslationEvent{
		Type:        "step_completed",
		SessionID:   sessionID,
		Step:        "conversion",
		Message:     "Markdown conversion completed",
		Progress:    30,
		CurrentItem: "markdown_ready",
		TotalItems:  totalLines,
		Timestamp:   time.Now().Unix(),
	})

	// Simulate translation progress line by line
	emitEvent(ws, TranslationEvent{
		Type:       "translation_progress",
		SessionID:  sessionID,
		Step:       "translation",
		Message:    "Starting translation...",
		Progress:   35,
		Timestamp:  time.Now().Unix(),
	})

	translatedLines := make([]string, 0, len(lines))
	for i, line := range lines {
		// Simulate translation time
		time.Sleep(300 * time.Millisecond)

		// Simple demo translation (in reality, this would call LLM APIs)
		translatedLine := translateDemoText(line)
		translatedLines = append(translatedLines, translatedLine)

		// Calculate progress
		progress := 35.0 + (float64(i+1)/float64(len(lines)) * 50.0)

		emitEvent(ws, TranslationEvent{
			Type:        "translation_progress",
			SessionID:   sessionID,
			Step:        "translation",
			Message:     fmt.Sprintf("Translating line %d/%d", i+1, len(lines)),
			Progress:    progress,
			CurrentItem: fmt.Sprintf("line_%d", i+1),
			TotalItems:  len(lines),
			Timestamp:   time.Now().Unix(),
		})
	}

	emitEvent(ws, TranslationEvent{
		Type:        "step_completed",
		SessionID:   sessionID,
		Step:        "translation",
		Message:     "Translation completed",
		Progress:    85,
		CurrentItem: "translation_complete",
		TotalItems:  len(lines),
		Timestamp:   time.Now().Unix(),
	})

	// Generate output
	emitEvent(ws, TranslationEvent{
		Type:       "translation_progress",
		SessionID:  sessionID,
		Step:       "generation",
		Message:    "Generating output file...",
		Progress:   90,
		Timestamp:  time.Now().Unix(),
	})

	translatedMarkdown := strings.Join(translatedLines, "\n")
	err = os.WriteFile(outputFile, []byte(translatedMarkdown), 0644)
	if err != nil {
		emitErrorEvent(ws, sessionID, "generation", fmt.Sprintf("Failed to write output: %v", err))
		log.Fatalf("Failed to write output: %v", err)
	}

	time.Sleep(1 * time.Second)

	emitEvent(ws, TranslationEvent{
		Type:        "translation_completed",
		SessionID:   sessionID,
		Message:     fmt.Sprintf("Translation completed successfully! Output saved to %s", outputFile),
		Progress:    100,
		CurrentItem: "output_generated",
		TotalItems:  len(lines),
		Timestamp:   time.Now().Unix(),
	})

	fmt.Printf("\n🎉 Translation completed!\n")
	fmt.Printf("📁 Output file: %s\n", outputFile)
	fmt.Printf("📊 View progress at: http://localhost:8090/monitor\n")
	fmt.Printf("🔗 Monitor this session: ws://localhost:8090/ws?session_id=%s\n", sessionID)

	// Keep running to allow monitoring
	fmt.Println("\nPress Ctrl+C to exit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n👋 Demo completed!")
}

func connectWebSocket(sessionID string) (*websocket.Conn, error) {
	u := fmt.Sprintf("ws://localhost:8090/ws?session_id=%s&client_id=demo-translator", sessionID)
	
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	
	return conn, nil
}

func listenForWebSocketEvents(ws *websocket.Conn) {
	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}
		
		log.Printf("📩 Received WebSocket message: %+v", msg)
	}
}

func emitEvent(ws *websocket.Conn, event TranslationEvent) {
	if ws == nil {
		return
	}
	
	event.Timestamp = time.Now().Unix()
	
	if err := ws.WriteJSON(event); err != nil {
		log.Printf("Failed to emit event: %v", err)
	}
	
	fmt.Printf("📤 Event: %s (%.1f%%) - %s\n", event.Type, event.Progress, event.Message)
}

func emitErrorEvent(ws *websocket.Conn, sessionID, step, message string) {
	event := TranslationEvent{
		Type:      "translation_error",
		SessionID: sessionID,
		Step:      step,
		Error:     message,
		Timestamp: time.Now().Unix(),
	}
	emitEvent(ws, event)
}

func convertToMarkdown(text string) string {
	lines := strings.Split(text, "\n")
	var markdownLines []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			markdownLines = append(markdownLines, trimmed)
		} else {
			markdownLines = append(markdownLines, "")
		}
	}
	
	return strings.Join(markdownLines, "\n")
}

// Simple demo translation (Russian to Serbian Cyrillic)
func translateDemoText(text string) string {
	// This is a simplified demo - in reality this would use LLM APIs
	replacements := map[string]string{
		"Это":         "Ово",
		"образец":     "пример",
		"русского":    "руског",
		"текста":      "текста",
		"для":         "за",
		"проверки":    "проверу",
		"функции":     "функције",
		"перевода":    "превода",
		"Он":          "Он",
		"содержит":    "садржи",
		"несколько":   "неколико",
		"предложений":  "реченица",
		"и":           "и",
		"подходит":    "одговара",
		"для":         "за",
		"тестирования": "тестирања",
		"перевода":    "превода",
		"Текст":       "Текст",
		"включает":    "укључује",
		"различные":   "различите",
		"знаки":      "знакове",
		"препинания":  "знаке",
		"структуры":   "структуре",
		"предложений": "реченица",
		".":           ".",
	}
	
	result := text
	for russian, serbian := range replacements {
		result = strings.ReplaceAll(result, russian, serbian)
	}
	
	return result
}