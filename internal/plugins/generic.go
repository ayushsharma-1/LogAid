package plugins

import (
	"context"
	"regexp"
	"strings"
)

// GenericPlugin handles common shell errors
type GenericPlugin struct {
	BasePlugin
}

// NewGenericPlugin creates a new generic plugin
func NewGenericPlugin() *GenericPlugin {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`command not found`),
		regexp.MustCompile(`No such file or directory`),
		regexp.MustCompile(`Permission denied`),
		regexp.MustCompile(`cannot access`),
		regexp.MustCompile(`syntax error`),
		regexp.MustCompile(`Segmentation fault`),
	}

	return &GenericPlugin{
		BasePlugin: BasePlugin{
			name:        "generic",
			description: "Handles common shell and system errors",
			patterns:    patterns,
		},
	}
}

// CanHandle checks if this plugin can handle the given error
func (p *GenericPlugin) CanHandle(command, stderr string, exitCode int) bool {
	// Generic plugin can handle most errors as a fallback
	return p.MatchesPattern(stderr) || exitCode != 0
}

// Analyze analyzes generic errors
func (p *GenericPlugin) Analyze(ctx context.Context, command, stderr string, exitCode int) (*Analysis, error) {
	analysis := &Analysis{
		Plugin:   p.Name(),
		Severity: "medium",
		Category: "system",
		Context:  make(map[string]string),
	}

	cmd := ExtractCommand(command)
	analysis.Context["command"] = cmd
	analysis.Context["exit_code"] = string(rune(exitCode))

	// Analyze specific error types
	switch {
	case strings.Contains(stderr, "command not found"):
		analysis.ErrorType = "command_not_found"
		analysis.Description = "Command not found in PATH"
		analysis.Confidence = 0.95
		analysis.Context["missing_command"] = cmd

	case strings.Contains(stderr, "No such file or directory"):
		analysis.ErrorType = "file_not_found"
		analysis.Description = "File or directory does not exist"
		analysis.Confidence = 0.90

	case strings.Contains(stderr, "Permission denied"):
		analysis.ErrorType = "permission_denied"
		analysis.Description = "Insufficient permissions"
		analysis.Confidence = 0.95

	case strings.Contains(stderr, "cannot access"):
		analysis.ErrorType = "access_error"
		analysis.Description = "Cannot access file or resource"
		analysis.Confidence = 0.85

	case strings.Contains(stderr, "syntax error"):
		analysis.ErrorType = "syntax_error"
		analysis.Description = "Syntax error in command or script"
		analysis.Confidence = 0.80

	case strings.Contains(stderr, "Segmentation fault"):
		analysis.ErrorType = "segfault"
		analysis.Description = "Program crashed with segmentation fault"
		analysis.Confidence = 0.95
		analysis.Severity = "high"

	default:
		analysis.ErrorType = "generic_error"
		analysis.Description = "Generic command execution error"
		analysis.Confidence = 0.60
	}

	return analysis, nil
}

// GetSuggestions provides suggestions for generic errors
func (p *GenericPlugin) GetSuggestions(ctx context.Context, analysis *Analysis) ([]*Suggestion, error) {
	var suggestions []*Suggestion

	switch analysis.ErrorType {
	case "command_not_found":
		cmd := analysis.Context["missing_command"]
		suggestions = append(suggestions, &Suggestion{
			Command:     "which " + cmd,
			Description: "Check if command exists in PATH",
			Type:        "info",
			Confidence:  0.90,
		})

		// Common package suggestions
		if strings.Contains(cmd, "curl") {
			suggestions = append(suggestions, &Suggestion{
				Command:     "sudo apt-get install curl",
				Description: "Install curl (Ubuntu/Debian)",
				Type:        "fix",
				Confidence:  0.80,
			})
		} else if strings.Contains(cmd, "git") {
			suggestions = append(suggestions, &Suggestion{
				Command:     "sudo apt-get install git",
				Description: "Install Git (Ubuntu/Debian)",
				Type:        "fix",
				Confidence:  0.80,
			})
		} else if strings.Contains(cmd, "node") || strings.Contains(cmd, "npm") {
			suggestions = append(suggestions, &Suggestion{
				Command:     "curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash - && sudo apt-get install -y nodejs",
				Description: "Install Node.js and npm",
				Type:        "fix",
				Confidence:  0.75,
			})
		}

	case "file_not_found":
		suggestions = append(suggestions, &Suggestion{
			Command:     "ls -la",
			Description: "List files in current directory",
			Type:        "info",
			Confidence:  0.85,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "pwd",
			Description: "Show current working directory",
			Type:        "info",
			Confidence:  0.90,
		})

		// Check if this might be a missing flag for ls command
		if strings.Contains(analysis.Context["command"], "ls") {
			suggestions = append(suggestions, &Suggestion{
				Command:     "ls -ltr",
				Description: "List files in long format, sorted by time (newest last)",
				Type:        "fix",
				Confidence:  0.95,
			})
			suggestions = append(suggestions, &Suggestion{
				Command:     "ls -la",
				Description: "List all files including hidden ones in long format",
				Type:        "alternative",
				Confidence:  0.90,
			})
		}

	case "permission_denied":
		suggestions = append(suggestions, &Suggestion{
			Command:     "ls -la",
			Description: "Check file permissions",
			Type:        "info",
			Confidence:  0.90,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "sudo " + analysis.Context["command"],
			Description: "Run with elevated privileges",
			Type:        "fix",
			Confidence:  0.75,
		})

	case "access_error":
		suggestions = append(suggestions, &Suggestion{
			Command:     "ls -la",
			Description: "Check file existence and permissions",
			Type:        "info",
			Confidence:  0.85,
		})

	case "syntax_error":
		suggestions = append(suggestions, &Suggestion{
			Command:     "man " + analysis.Context["command"],
			Description: "Check command manual for correct syntax",
			Type:        "info",
			Confidence:  0.80,
		})

	case "segfault":
		suggestions = append(suggestions, &Suggestion{
			Command:     "ulimit -c unlimited && " + analysis.Context["command"],
			Description: "Enable core dumps for debugging",
			Type:        "info",
			Confidence:  0.70,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "gdb " + analysis.Context["command"],
			Description: "Debug with GDB (if available)",
			Type:        "info",
			Confidence:  0.60,
		})
	}

	return suggestions, nil
}
