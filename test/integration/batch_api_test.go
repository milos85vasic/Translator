//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"digital.vasic.translator/internal/cache"
	"digital.vasic.translator/internal/config"
	"digital.vasic.translator/pkg/api"
	"digital.vasic.translator/pkg/events"
	"digital.vasic.translator/pkg/security"
	"digital.vasic.translator/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func setupTestAPI() (*gin.Engine, *api.Handler, *events.EventBus) {
	gin.SetMode(gin.TestMode)

	cfg := config.DefaultConfig()
	cfg.Translation.DefaultProvider = "dictionary"
	cfg.Translation.DefaultModel = "default"

	eventBus := events.NewEventBus()
	cacheStore := cache.NewCache(time.Hour, true)
	authService := security.NewAuthService("test-secret", 24*time.Hour)
	wsHub := websocket.NewHub(eventBus)

	handler := api.NewHandler(cfg, eventBus, cacheStore, authService, wsHub, nil)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterBatchRoutes(v1)

	return router, handler, eventBus
}

func TestStringTranslationAPI(t *testing.T) {
	router, _, _ := setupTestAPI()

	t.Run("TranslateString", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"text":            "Hello, world!",
			"target_language": "sr",
			"provider":        "dictionary",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/string", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response api.TranslateStringResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.TranslatedText == "" {
			t.Error("Expected translated text, got empty string")
		}

		if response.TargetLanguage != "sr" {
			t.Errorf("Expected target language 'sr', got '%s'", response.TargetLanguage)
		}

		if response.SessionID == "" {
			t.Error("Expected session ID, got empty string")
		}

		if response.Duration <= 0 {
			t.Error("Expected positive duration")
		}
	})

	t.Run("TranslateStringMissingText", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"target_language": "sr",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/string", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("TranslateStringInvalidLanguage", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"text":            "Hello",
			"target_language": "invalid",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/string", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("TranslateStringWithSourceLanguage", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"text":            "Hello",
			"source_language": "en",
			"target_language": "sr",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/string", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response api.TranslateStringResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.SourceLanguage != "en" {
			t.Errorf("Expected source language 'en', got '%s'", response.SourceLanguage)
		}
	})
}

func TestDirectoryTranslationAPI(t *testing.T) {
	router, _, _ := setupTestAPI()

	// Create temporary test directory
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	os.MkdirAll(inputDir, 0755)

	// Create test files
	os.WriteFile(filepath.Join(inputDir, "test1.txt"), []byte("Test content 1"), 0644)
	os.WriteFile(filepath.Join(inputDir, "test2.txt"), []byte("Test content 2"), 0644)

	t.Run("TranslateDirectory", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"input_path":      inputDir,
			"output_path":     outputDir,
			"target_language": "sr",
			"provider":        "dictionary",
			"recursive":       false,
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/directory", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var response api.TranslateDirectoryResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.TotalFiles != 2 {
			t.Errorf("Expected 2 files, got %d", response.TotalFiles)
		}

		if response.SessionID == "" {
			t.Error("Expected session ID")
		}

		if response.Duration <= 0 {
			t.Error("Expected positive duration")
		}

		if len(response.Results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(response.Results))
		}
	})

	t.Run("TranslateDirectoryRecursive", func(t *testing.T) {
		// Create subdirectory
		subDir := filepath.Join(inputDir, "subdir")
		os.MkdirAll(subDir, 0755)
		os.WriteFile(filepath.Join(subDir, "test3.txt"), []byte("Test content 3"), 0644)

		reqBody := map[string]interface{}{
			"input_path":      inputDir,
			"output_path":     outputDir,
			"target_language": "sr",
			"provider":        "dictionary",
			"recursive":       true,
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/directory", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response api.TranslateDirectoryResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.TotalFiles < 3 {
			t.Errorf("Expected at least 3 files with recursive, got %d", response.TotalFiles)
		}
	})

	t.Run("TranslateDirectoryParallel", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"input_path":      inputDir,
			"output_path":     outputDir,
			"target_language": "sr",
			"provider":        "dictionary",
			"parallel":        true,
			"max_concurrency": 2,
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/directory", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("TranslateDirectoryMissingPath", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"target_language": "sr",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/directory", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("TranslateDirectoryInvalidPath", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"input_path":      "/nonexistent/path",
			"target_language": "sr",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/translate/directory", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})
}

func TestAPIEventEmission(t *testing.T) {
	_, handler, eventBus := setupTestAPI()

	eventReceived := make(chan bool, 1)

	// Subscribe to events
	eventBus.Subscribe(events.EventTranslationStarted, func(event events.Event) {
		eventReceived <- true
	})

	// Create test request
	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterBatchRoutes(v1)

	reqBody := map[string]interface{}{
		"text":            "Test",
		"target_language": "sr",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/translate/string", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Wait for event
	select {
	case <-eventReceived:
		t.Log("Event received successfully")
	case <-time.After(1 * time.Second):
		t.Error("Expected translation_started event not received")
	}
}
