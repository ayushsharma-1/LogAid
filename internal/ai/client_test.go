package ai

import (
	"testing"
)

func TestParseAISuggestion(t *testing.T) {
	// Test valid JSON response
	jsonResponse := `{
		"command": "mkdir -p /path/to/directory",
		"explanation": "The directory does not exist. Create it first.",
		"confidence": 0.95
	}`

	suggestion, err := parseAISuggestion(jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse valid JSON: %v", err)
	}

	if suggestion.Command != "mkdir -p /path/to/directory" {
		t.Errorf("Expected command 'mkdir -p /path/to/directory', got '%s'", suggestion.Command)
	}

	if suggestion.Confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", suggestion.Confidence)
	}
}

func TestParseAISuggestionWithMarkdown(t *testing.T) {
	// Test JSON response wrapped in markdown code blocks
	markdownResponse := "```json\n" + `{
		"command": "ls -la",
		"explanation": "List files with details",
		"confidence": 0.8
	}` + "\n```"

	suggestion, err := parseAISuggestion(markdownResponse)
	if err != nil {
		t.Fatalf("Failed to parse markdown JSON: %v", err)
	}

	if suggestion.Command != "ls -la" {
		t.Errorf("Expected command 'ls -la', got '%s'", suggestion.Command)
	}
}

func TestParseAISuggestionInvalidJSON(t *testing.T) {
	// Test invalid JSON (should create fallback suggestion)
	invalidJSON := "This is not JSON but still useful advice"

	suggestion, err := parseAISuggestion(invalidJSON)
	if err != nil {
		t.Fatalf("Should not error on invalid JSON: %v", err)
	}

	if suggestion.Explanation != invalidJSON {
		t.Errorf("Expected explanation '%s', got '%s'", invalidJSON, suggestion.Explanation)
	}

	if suggestion.Confidence != 0.5 {
		t.Errorf("Expected fallback confidence 0.5, got %f", suggestion.Confidence)
	}
}

func TestBuildErrorAnalysisPrompt(t *testing.T) {
	command := "ls /nonexistent"
	stderr := "ls: cannot access '/nonexistent': No such file or directory"
	exitCode := 2

	prompt := buildErrorAnalysisPrompt(command, stderr, exitCode)

	// Check that all components are included in the prompt
	if !contains(prompt, command) {
		t.Error("Prompt should contain the command")
	}

	if !contains(prompt, stderr) {
		t.Error("Prompt should contain the stderr output")
	}

	if !contains(prompt, "EXIT CODE: 2") {
		t.Error("Prompt should contain the exit code")
	}

	if !contains(prompt, "JSON format") {
		t.Error("Prompt should request JSON format response")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOfSubstring(s, substr) >= 0))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
