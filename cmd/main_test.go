package main

import (
	"os"
	"testing"
)

// TestMain sets up test environment
func TestMain(m *testing.M) {
	// Setup: ensure we have required environment for tests
	if os.Getenv("GEMINI_API_KEY") == "" && os.Getenv("OPENAI_API_KEY") == "" {
		os.Setenv("GEMINI_API_KEY", "test-key-for-testing")
	}
	
	// Run tests
	code := m.Run()
	
	// Cleanup if needed
	os.Exit(code)
}

func TestShowLogo(t *testing.T) {
	// Test that showLogo function exists and doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("showLogo() panicked: %v", r)
		}
	}()
	
	// This will test that the function can be called without error
	showLogo()
}
