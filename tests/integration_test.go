package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestBuild(t *testing.T) {
	// Test that the application builds successfully
	cmd := exec.Command("go", "build", "-o", "logaid_test", "../cmd")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Build failed: %v\nOutput: %s", err, output)
	}
	
	// Clean up
	defer os.Remove("logaid_test")
	
	// Verify the binary was created
	if _, err := os.Stat("logaid_test"); os.IsNotExist(err) {
		t.Fatal("Binary was not created")
	}
}

func TestVersion(t *testing.T) {
	// First build the application
	buildCmd := exec.Command("go", "build", "-o", "logaid_test", "../cmd")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}
	defer os.Remove("logaid_test")
	
	// Test version command with minimal environment
	cmd := exec.Command("./logaid_test", "version")
	cmd.Env = append(os.Environ(), "GEMINI_API_KEY=test-key")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Version command failed: %v\nOutput: %s", err, output)
	}
	
	outputStr := string(output)
	if !contains(outputStr, "LogAid v1.0.0") {
		t.Error("Version output doesn't contain expected version string")
	}
}

func TestHelp(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "logaid_test", "../cmd")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}
	defer os.Remove("logaid_test")
	
	cmd := exec.Command("./logaid_test", "help")
	cmd.Env = append(os.Environ(), "GEMINI_API_KEY=test-key")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Help command failed: %v\nOutput: %s", err, output)
	}
	
	outputStr := string(output)
	expectedStrings := []string{"LogAid - Your CLI Guardian", "USAGE:", "COMMANDS:", "start", "config", "help", "version"}
	
	for _, expected := range expectedStrings {
		if !contains(outputStr, expected) {
			t.Errorf("Help output doesn't contain expected string: %s", expected)
		}
	}
}

func TestConfig(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "logaid_test", "../cmd")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}
	defer os.Remove("logaid_test")
	
	cmd := exec.Command("./logaid_test", "config")
	cmd.Env = append(os.Environ(), 
		"GEMINI_API_KEY=test-key",
		"AI_PROVIDER=gemini",
		"LOG_LEVEL=debug",
	)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Config command failed: %v\nOutput: %s", err, output)
	}
	
	outputStr := string(output)
	expectedStrings := []string{"LogAid Configuration:", "AI Provider:", "Log Level:", "Shell:"}
	
	for _, expected := range expectedStrings {
		if !contains(outputStr, expected) {
			t.Errorf("Config output doesn't contain expected string: %s", expected)
		}
	}
}

func TestInvalidCommand(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "logaid_test", "../cmd")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}
	defer os.Remove("logaid_test")
	
	cmd := exec.Command("./logaid_test", "invalid-command")
	cmd.Env = append(os.Environ(), "GEMINI_API_KEY=test-key")
	
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for invalid command, but command succeeded")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    contains(s[1:], substr) || 
		    (len(s) > 0 && s[:len(substr)] == substr))
}
