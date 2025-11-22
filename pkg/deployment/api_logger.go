package deployment

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// APICommunicationLogger logs all REST API communications between distributed nodes
type APICommunicationLogger struct {
	logFile *os.File
	logger  *log.Logger
	mu      sync.Mutex
}

// NewAPICommunicationLogger creates a new API communication logger
func NewAPICommunicationLogger(logPath string) (*APICommunicationLogger, error) {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New(file, "", 0) // No prefix, we'll format ourselves

	return &APICommunicationLogger{
		logFile: file,
		logger:  logger,
	}, nil
}

// LogCommunication logs an API communication event
func (acl *APICommunicationLogger) LogCommunication(logEntry *APICommunicationLog) error {
	acl.mu.Lock()
	defer acl.mu.Unlock()

	// Format as JSON for structured logging
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	acl.logger.Println(string(jsonData))
	return nil
}

// LogRequest logs an outgoing API request
func (acl *APICommunicationLogger) LogRequest(sourceHost string, sourcePort int, targetHost string, targetPort int, method, url string, requestSize int64) *APICommunicationLog {
	entry := &APICommunicationLog{
		Timestamp:   time.Now(),
		SourceHost:  sourceHost,
		SourcePort:  sourcePort,
		TargetHost:  targetHost,
		TargetPort:  targetPort,
		Method:      method,
		URL:         url,
		RequestSize: requestSize,
	}

	// Log asynchronously to avoid blocking
	go func() {
		if err := acl.LogCommunication(entry); err != nil {
			log.Printf("Failed to log API request: %v", err)
		}
	}()

	return entry
}

// LogResponse logs the response for a previously logged request
func (acl *APICommunicationLogger) LogResponse(entry *APICommunicationLog, statusCode int, responseSize int64, duration time.Duration, err error) {
	entry.StatusCode = statusCode
	entry.ResponseSize = responseSize
	entry.Duration = duration

	if err != nil {
		entry.Error = err.Error()
	}

	// Update the existing log entry
	go func() {
		if logErr := acl.LogCommunication(entry); err != nil {
			log.Printf("Failed to log API response: %v", logErr)
		}
	}()
}

// GetLogs retrieves recent log entries
func (acl *APICommunicationLogger) GetLogs(limit int) ([]*APICommunicationLog, error) {
	acl.mu.Lock()
	defer acl.mu.Unlock()

	// This is a simplified implementation
	// In a real system, you'd want to read from the log file and parse recent entries
	return []*APICommunicationLog{}, nil
}

// Close closes the logger
func (acl *APICommunicationLogger) Close() error {
	acl.mu.Lock()
	defer acl.mu.Unlock()

	if acl.logFile != nil {
		return acl.logFile.Close()
	}
	return nil
}

// GetStats returns communication statistics
func (acl *APICommunicationLogger) GetStats() map[string]interface{} {
	// This would parse the log file to generate statistics
	// For now, return empty stats
	return map[string]interface{}{
		"total_requests":  0,
		"total_responses": 0,
		"error_count":     0,
		"avg_duration":    "0s",
	}
}
