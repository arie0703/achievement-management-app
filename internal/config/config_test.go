package config

import (
	"os"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear environment variables
	os.Clearenv()
	
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check default values
	if config.Environment != "development" {
		t.Errorf("Expected environment 'development', got '%s'", config.Environment)
	}
	
	if config.AWS.Region != "us-east-1" {
		t.Errorf("Expected AWS region 'us-east-1', got '%s'", config.AWS.Region)
	}
	
	if config.Tables.Achievements != "achievements" {
		t.Errorf("Expected achievements table 'achievements', got '%s'", config.Tables.Achievements)
	}
	
	if config.Server.Port != "8080" {
		t.Errorf("Expected server port '8080', got '%s'", config.Server.Port)
	}
	
	if config.Logging.Level != "info" {
		t.Errorf("Expected log level 'info', got '%s'", config.Logging.Level)
	}
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("AWS_REGION", "ap-northeast-1")
	os.Setenv("ACHIEVEMENTS_TABLE", "prod-achievements")
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("LOG_LEVEL", "error")
	
	defer func() {
		os.Clearenv()
	}()
	
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check environment variable overrides
	if config.Environment != "production" {
		t.Errorf("Expected environment 'production', got '%s'", config.Environment)
	}
	
	if config.AWS.Region != "ap-northeast-1" {
		t.Errorf("Expected AWS region 'ap-northeast-1', got '%s'", config.AWS.Region)
	}
	
	if config.Tables.Achievements != "prod-achievements" {
		t.Errorf("Expected achievements table 'prod-achievements', got '%s'", config.Tables.Achievements)
	}
	
	if config.Server.Port != "9000" {
		t.Errorf("Expected server port '9000', got '%s'", config.Server.Port)
	}
	
	if config.Logging.Level != "error" {
		t.Errorf("Expected log level 'error', got '%s'", config.Logging.Level)
	}
}

func TestValidateConfig_InvalidEnvironment(t *testing.T) {
	config := getDefaultConfig()
	config.Environment = "invalid"
	
	err := validateConfig(config)
	if err == nil {
		t.Error("Expected validation error for invalid environment")
	}
}

func TestValidateConfig_EmptyRegion(t *testing.T) {
	config := getDefaultConfig()
	config.AWS.Region = ""
	
	err := validateConfig(config)
	if err == nil {
		t.Error("Expected validation error for empty AWS region")
	}
}

func TestValidateConfig_InvalidLogLevel(t *testing.T) {
	config := getDefaultConfig()
	config.Logging.Level = "invalid"
	
	err := validateConfig(config)
	if err == nil {
		t.Error("Expected validation error for invalid log level")
	}
}

func TestCreateConfigFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)
	
	err := CreateConfigFile("test")
	if err != nil {
		t.Fatalf("Expected no error creating config file, got %v", err)
	}
	
	// Check if file was created
	if _, err := os.Stat("config/test.json"); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	
	result := getEnvAsInt("TEST_INT", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
	
	result = getEnvAsInt("NON_EXISTENT", 10)
	if result != 10 {
		t.Errorf("Expected default value 10, got %d", result)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}
	
	if !contains(slice, "b") {
		t.Error("Expected 'b' to be found in slice")
	}
	
	if contains(slice, "d") {
		t.Error("Expected 'd' not to be found in slice")
	}
}