package deployment

import (
	"os"
	"testing"
	"time"
)

func TestAPICommunicationLogger_LogCommunication(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := os.CreateTemp("", "api_log_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create logger
	logger, err := NewAPICommunicationLogger(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create API logger: %v", err)
	}
	defer logger.Close()

	// Create a test log entry
	entry := &APICommunicationLog{
		Timestamp:    time.Now(),
		SourceHost:   "192.168.1.100",
		SourcePort:   8443,
		TargetHost:   "192.168.1.101",
		TargetPort:   8444,
		Method:       "POST",
		URL:          "/api/v1/translate",
		StatusCode:   200,
		RequestSize:  1024,
		ResponseSize: 2048,
		Duration:     150 * time.Millisecond,
		UserAgent:    "translator/1.0",
	}

	// Log the communication
	err = logger.LogCommunication(entry)
	if err != nil {
		t.Fatalf("Failed to log communication: %v", err)
	}

	// Read the log file to verify
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Check that the log contains expected data
	contentStr := string(content)
	if len(contentStr) == 0 {
		t.Error("Log file is empty")
	}

	t.Logf("Log content: %s", contentStr)
}

func TestAPICommunicationLogger_LogRequest(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := os.CreateTemp("", "api_log_request_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create logger
	logger, err := NewAPICommunicationLogger(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create API logger: %v", err)
	}
	defer logger.Close()

	// Log a request
	entry := logger.LogRequest("192.168.1.100", 8443, "192.168.1.101", 8444, "POST", "/api/v1/translate", 1024)

	// Verify the entry was created
	if entry.SourceHost != "192.168.1.100" {
		t.Errorf("Expected source host '192.168.1.100', got '%s'", entry.SourceHost)
	}
	if entry.Method != "POST" {
		t.Errorf("Expected method 'POST', got '%s'", entry.Method)
	}
	if entry.RequestSize != 1024 {
		t.Errorf("Expected request size 1024, got %d", entry.RequestSize)
	}
}

func TestAPICommunicationLogger_LogResponse(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := os.CreateTemp("", "api_log_response_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create logger
	logger, err := NewAPICommunicationLogger(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create API logger: %v", err)
	}
	defer logger.Close()

	// Create a request entry first
	entry := logger.LogRequest("192.168.1.100", 8443, "192.168.1.101", 8444, "POST", "/api/v1/translate", 1024)

	// Log the response
	logger.LogResponse(entry, 200, 2048, 150*time.Millisecond, nil)

	// Verify the response data was added
	if entry.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", entry.StatusCode)
	}
	if entry.ResponseSize != 2048 {
		t.Errorf("Expected response size 2048, got %d", entry.ResponseSize)
	}
	if entry.Duration != 150*time.Millisecond {
		t.Errorf("Expected duration 150ms, got %v", entry.Duration)
	}
}

func TestAPICommunicationLogger_GetStats(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := os.CreateTemp("", "api_log_stats_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create logger
	logger, err := NewAPICommunicationLogger(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create API logger: %v", err)
	}
	defer logger.Close()

	// Get stats (should return empty stats for now)
	stats := logger.GetStats()

	// Verify stats structure
	if stats == nil {
		t.Error("Expected stats map, got nil")
	}

	// Check for expected keys
	expectedKeys := []string{"total_requests", "total_responses", "error_count", "avg_duration"}
	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected stats key '%s' not found", key)
		}
	}

	t.Logf("Stats: %v", stats)
}

func TestAPICommunicationLogger_ErrorHandling(t *testing.T) {
	// Test with invalid log file path
	logger, err := NewAPICommunicationLogger("/invalid/path/api_communication.log")
	if err == nil {
		t.Error("Expected error when creating logger with invalid path")
		logger.Close()
	}

	// Test logging after close
	tmpFile, err := os.CreateTemp("", "api_log_close_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	logger, err = NewAPICommunicationLogger(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create API logger: %v", err)
	}

	logger.Close()

	// Try to log after close (should not panic)
	entry := &APICommunicationLog{
		Timestamp:  time.Now(),
		SourceHost: "test",
		Method:     "GET",
		URL:        "/test",
	}

	err = logger.LogCommunication(entry)
	// We don't expect this to succeed, but it shouldn't panic
	t.Logf("Logging after close: %v", err)
}
