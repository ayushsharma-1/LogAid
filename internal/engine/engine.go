package engine

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ayush-1/logaid/internal/ai"
	"github.com/ayush-1/logaid/internal/config"
	"github.com/ayush-1/logaid/internal/logger"
	"github.com/ayush-1/logaid/internal/plugins"
)

// Engine represents the core LogAid engine
type Engine struct {
	plugins []plugins.Plugin
}

// New creates a new Engine instance
func New() *Engine {
	return &Engine{
		plugins: plugins.LoadAllPlugins(),
	}
}

// ProcessError processes a command error and returns a suggestion
func (e *Engine) ProcessError(ctx context.Context, command, output string) (string, error) {
	// Try plugins first
	for _, plugin := range e.plugins {
		if plugin.Match(command, output) {
			suggestion := plugin.Suggest(command, output)
			if suggestion != "" {
				return suggestion, nil
			}
		}
	}

	// If no plugin matched, use AI directly
	suggestion, err := ai.GetSuggestion(ctx, fmt.Sprintf("Command: %s\nError: %s\nProvide a corrected command:", command, output))
	if err != nil {
		return "", fmt.Errorf("failed to get AI suggestion: %w", err)
	}

	return suggestion, nil
}

// detectError checks if the output contains error indicators
func (e *Engine) detectError(output string) bool {
	errorIndicators := []string{
		"error:",
		"Error:",
		"ERROR:",
		"failed",
		"Failed",
		"FAILED",
		"not found",
		"Not found",
		"command not found",
		"is not a git command",
		"is not a docker command",
		"permission denied",
		"Permission denied",
		"E: Unable to locate package",
		"npm ERR!",
		"fatal:",
		"Fatal:",
	}

	lowerOutput := strings.ToLower(output)
	for _, indicator := range errorIndicators {
		if strings.Contains(lowerOutput, strings.ToLower(indicator)) {
			return true
		}
	}

	return false
}

func (e *Engine) handleError(command, output string) bool {
	logger.Warn("Error detected in command output")

	// Try plugins first
	for _, plugin := range e.plugins {
		if plugin.Match(command, output) {
			suggestion := plugin.Suggest(command, output)
			if suggestion != "" {
				return e.presentSuggestion(command, output, suggestion, plugin.Name())
			}
		}
	}

	// If no plugin matched, use AI
	ctx := context.Background()
	suggestion, err := ai.GetSuggestion(ctx, fmt.Sprintf("Command: %s\nError: %s\nProvide a corrected command:", command, output))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get AI suggestion: %v", err))
		return false
	}

	if suggestion != "" {
		return e.presentSuggestion(command, output, suggestion, "AI")
	}

	return false
}

func (e *Engine) presentSuggestion(command, output, suggestion, source string) bool {
	logger.Warn(fmt.Sprintf("Suggestion from %s:", source))
	logger.Info(fmt.Sprintf("💡 %s", suggestion))

	// Check if auto-confirm is enabled
	if config.AppConfig != nil && config.AppConfig.AutoConfirm {
		logger.Info("Auto-confirm enabled, executing suggestion...")
		return e.executeSuggestion(suggestion)
	}

	// Prompt user for confirmation
	logger.Info("Execute this suggestion? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read user input: %v", err))
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input == "y" || input == "yes" {
		logger.Info("Executing suggestion...")
		return e.executeSuggestion(suggestion)
	} else {
		logger.Info("Suggestion ignored.")
		return false
	}
}

func (e *Engine) executeSuggestion(suggestion string) bool {
	// Parse the suggestion into command and args
	parts := strings.Fields(suggestion)
	if len(parts) == 0 {
		logger.Error("Invalid suggestion: empty command")
		return false
	}

	var cmd *exec.Cmd
	if len(parts) > 1 {
		cmd = exec.Command(parts[0], parts[1:]...)
	} else {
		cmd = exec.Command(parts[0])
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info(fmt.Sprintf("Running: %s", suggestion))
	err := cmd.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("Suggestion execution failed: %v", err))
		return false
	} else {
		logger.Info("Suggestion executed successfully!")
		return true
	}
}

// ExecuteWithMonitoring executes a command with LogAid monitoring
func ExecuteWithMonitoring(cmd *exec.Cmd) error {
	engine := New()

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	// Execute the command
	err := cmd.Run()

	// Combine command for logging
	command := strings.Join(cmd.Args, " ")

	if err != nil {
		// Command failed, analyze the error
		output := stderr.String()
		if output == "" {
			output = stdout.String()
		}

		logger.Error(fmt.Sprintf("Command failed: %s", command))

		if engine.detectError(output) {
			// If we successfully handle the error (user accepts and suggestion works), return success
			if engine.handleError(command, output) {
				return nil // Suggestion executed successfully, don't return original error
			}
		}

		return err // Return original error if no suggestion or suggestion failed
	}

	// Check stdout for potential issues even if command succeeded
	output := stdout.String()
	if engine.detectError(output) {
		logger.Warn("Potential issues detected in command output")
		engine.handleError(command, output)
	}

	return nil
}
