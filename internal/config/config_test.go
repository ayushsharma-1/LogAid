package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("LOGAID_AI_PROVIDER", "gemini")
	os.Setenv("GEMINI_API_KEY", "test-key")
	os.Setenv("LOGAID_ENABLE_COLORS", "true")

	defer func() {
		os.Unsetenv("LOGAID_AI_PROVIDER")
		os.Unsetenv("GEMINI_API_KEY")
		os.Unsetenv("LOGAID_ENABLE_COLORS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.AIProvider != "gemini" {
		t.Errorf("Expected AIProvider to be 'gemini', got '%s'", cfg.AIProvider)
	}

	if cfg.APIKey != "test-key" {
		t.Errorf("Expected APIKey to be 'test-key', got '%s'", cfg.APIKey)
	}

	if !cfg.EnableColors {
		t.Error("Expected EnableColors to be true")
	}
}

func TestLoadConfigMissingAPIKey(t *testing.T) {
	// Clear any existing API key environment variables
	os.Unsetenv("LOGAID_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")

	_, err := Load()
	if err == nil {
		t.Error("Expected error when API key is missing")
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnvWithDefault("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with environment variable not set
	result = getEnvWithDefault("NON_EXISTENT_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}

func TestGetEnvInt(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := getEnvInt("TEST_INT", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid integer (should return default)
	os.Setenv("TEST_INT_INVALID", "not_a_number")
	defer os.Unsetenv("TEST_INT_INVALID")

	result = getEnvInt("TEST_INT_INVALID", 10)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}
}

func TestGetEnvBool(t *testing.T) {
	// Test with valid boolean
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	result := getEnvBool("TEST_BOOL", false)
	if !result {
		t.Error("Expected true")
	}

	// Test with invalid boolean (should return default)
	os.Setenv("TEST_BOOL_INVALID", "not_a_bool")
	defer os.Unsetenv("TEST_BOOL_INVALID")

	result = getEnvBool("TEST_BOOL_INVALID", true)
	if !result {
		t.Error("Expected true (default value)")
	}
}
