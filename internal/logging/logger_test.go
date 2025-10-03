package logging

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"achievement-management/internal/config"
)

func TestNewLogger_JSONFormat(t *testing.T) {
	config := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if logger == nil {
		t.Fatal("Expected logger to be created")
	}
}

func TestNewLogger_TextFormat(t *testing.T) {
	config := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		},
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if logger == nil {
		t.Fatal("Expected logger to be created")
	}
}

func TestLogger_LogLevels(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "debug",
			Format: "json",
			Output: "stdout",
		},
	}
	
	logger := NewLoggerWithOutput(config, &buf)
	
	// Test different log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	if len(lines) != 4 {
		t.Errorf("Expected 4 log lines, got %d", len(lines))
	}
	
	// Check if each line is valid JSON
	for i, line := range lines {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Errorf("Line %d is not valid JSON: %v", i+1, err)
		}
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
	
	logger := NewLoggerWithOutput(config, &buf)
	
	logger.WithFields(map[string]interface{}{
		"user_id": "123",
		"action":  "create",
	}).Info("User action")
	
	output := buf.String()
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}
	
	if logEntry["user_id"] != "123" {
		t.Errorf("Expected user_id '123', got %v", logEntry["user_id"])
	}
	
	if logEntry["action"] != "create" {
		t.Errorf("Expected action 'create', got %v", logEntry["action"])
	}
}

func TestAccessLogger_LogRequest(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
	
	logger := NewLoggerWithOutput(config, &buf)
	accessLogger := &AccessLogger{logger: logger}
	
	duration := 100 * time.Millisecond
	accessLogger.LogRequest("GET", "/api/achievements", "127.0.0.1", 200, duration)
	
	output := buf.String()
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}
	
	if logEntry["method"] != "GET" {
		t.Errorf("Expected method 'GET', got %v", logEntry["method"])
	}
	
	if logEntry["path"] != "/api/achievements" {
		t.Errorf("Expected path '/api/achievements', got %v", logEntry["path"])
	}
	
	if logEntry["status_code"] != float64(200) {
		t.Errorf("Expected status_code 200, got %v", logEntry["status_code"])
	}
	
	if logEntry["type"] != "access" {
		t.Errorf("Expected type 'access', got %v", logEntry["type"])
	}
}

func TestErrorLogger_LogError(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "error",
			Format: "json",
			Output: "stdout",
		},
	}
	
	logger := NewLoggerWithOutput(config, &buf)
	errorLogger := &ErrorLogger{logger: logger}
	
	testErr := &testError{message: "test error"}
	errorLogger.LogError("test_operation", "test_component", testErr, map[string]interface{}{
		"extra_field": "extra_value",
	})
	
	output := buf.String()
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}
	
	if logEntry["operation"] != "test_operation" {
		t.Errorf("Expected operation 'test_operation', got %v", logEntry["operation"])
	}
	
	if logEntry["component"] != "test_component" {
		t.Errorf("Expected component 'test_component', got %v", logEntry["component"])
	}
	
	if logEntry["error"] != "test error" {
		t.Errorf("Expected error 'test error', got %v", logEntry["error"])
	}
	
	if logEntry["type"] != "error" {
		t.Errorf("Expected type 'error', got %v", logEntry["type"])
	}
	
	if logEntry["extra_field"] != "extra_value" {
		t.Errorf("Expected extra_field 'extra_value', got %v", logEntry["extra_field"])
	}
}

func TestErrorLogger_LogDatabaseError(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "error",
			Format: "json",
			Output: "stdout",
		},
	}
	
	logger := NewLoggerWithOutput(config, &buf)
	errorLogger := &ErrorLogger{logger: logger}
	
	testErr := &testError{message: "database connection failed"}
	errorLogger.LogDatabaseError("query", "achievements", testErr)
	
	output := buf.String()
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}
	
	if logEntry["operation"] != "query" {
		t.Errorf("Expected operation 'query', got %v", logEntry["operation"])
	}
	
	if logEntry["component"] != "database" {
		t.Errorf("Expected component 'database', got %v", logEntry["component"])
	}
	
	if logEntry["table"] != "achievements" {
		t.Errorf("Expected table 'achievements', got %v", logEntry["table"])
	}
}

// testError テスト用のエラー型
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}