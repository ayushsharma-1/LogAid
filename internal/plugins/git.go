package plugins

import (
	"context"
	"regexp"
	"strings"
)

// GitPlugin handles Git-related errors
type GitPlugin struct {
	BasePlugin
}

func init() {
	// Initialize Git plugin patterns
}

// NewGitPlugin creates a new Git plugin
func NewGitPlugin() *GitPlugin {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`fatal: not a git repository`),
		regexp.MustCompile(`fatal: remote .+ already exists`),
		regexp.MustCompile(`error: failed to push some refs`),
		regexp.MustCompile(`fatal: unable to access`),
		regexp.MustCompile(`fatal: repository .+ does not exist`),
		regexp.MustCompile(`error: pathspec .+ did not match any file`),
		regexp.MustCompile(`fatal: branch .+ already exists`),
		regexp.MustCompile(`error: Your local changes to the following files would be overwritten`),
	}

	return &GitPlugin{
		BasePlugin: BasePlugin{
			name:        "git",
			description: "Handles Git version control errors",
			patterns:    patterns,
		},
	}
}

// CanHandle checks if this plugin can handle the given error
func (p *GitPlugin) CanHandle(command, stderr string, exitCode int) bool {
	// Check if it's a git command
	if !p.MatchesCommand(command, []string{"git"}) {
		return false
	}

	// Check if stderr matches git error patterns
	return p.MatchesPattern(stderr)
}

// Analyze analyzes the Git error
func (p *GitPlugin) Analyze(ctx context.Context, command, stderr string, exitCode int) (*Analysis, error) {
	analysis := &Analysis{
		Plugin:   p.Name(),
		Severity: "medium",
		Category: "version_control",
		Context:  make(map[string]string),
	}

	// Parse the command
	args := ExtractArgs(command)
	analysis.Context["git_command"] = strings.Join(args, " ")

	// Analyze specific error types
	switch {
	case strings.Contains(stderr, "not a git repository"):
		analysis.ErrorType = "not_git_repo"
		analysis.Description = "Current directory is not a Git repository"
		analysis.Confidence = 0.95

	case strings.Contains(stderr, "remote") && strings.Contains(stderr, "already exists"):
		analysis.ErrorType = "remote_exists"
		analysis.Description = "Git remote already exists"
		analysis.Confidence = 0.90

	case strings.Contains(stderr, "failed to push"):
		analysis.ErrorType = "push_failed"
		analysis.Description = "Git push operation failed"
		analysis.Confidence = 0.85

	case strings.Contains(stderr, "unable to access"):
		analysis.ErrorType = "access_denied"
		analysis.Description = "Unable to access Git repository"
		analysis.Confidence = 0.90

	case strings.Contains(stderr, "does not exist"):
		analysis.ErrorType = "repo_not_found"
		analysis.Description = "Git repository does not exist"
		analysis.Confidence = 0.95

	case strings.Contains(stderr, "pathspec") && strings.Contains(stderr, "did not match"):
		analysis.ErrorType = "pathspec_no_match"
		analysis.Description = "File or path pattern not found"
		analysis.Confidence = 0.90

	case strings.Contains(stderr, "branch") && strings.Contains(stderr, "already exists"):
		analysis.ErrorType = "branch_exists"
		analysis.Description = "Git branch already exists"
		analysis.Confidence = 0.95

	case strings.Contains(stderr, "would be overwritten"):
		analysis.ErrorType = "conflicts"
		analysis.Description = "Local changes would be overwritten"
		analysis.Confidence = 0.90

	default:
		analysis.ErrorType = "generic_git_error"
		analysis.Description = "Generic Git error"
		analysis.Confidence = 0.70
	}

	return analysis, nil
}

// GetSuggestions provides suggestions for Git errors
func (p *GitPlugin) GetSuggestions(ctx context.Context, analysis *Analysis) ([]*Suggestion, error) {
	var suggestions []*Suggestion

	switch analysis.ErrorType {
	case "not_git_repo":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git init",
			Description: "Initialize a new Git repository",
			Type:        "fix",
			Confidence:  0.95,
		})

	case "remote_exists":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git remote -v",
			Description: "List existing remotes",
			Type:        "info",
			Confidence:  0.90,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "git remote set-url origin <new-url>",
			Description: "Update existing remote URL",
			Type:        "fix",
			Confidence:  0.85,
		})

	case "push_failed":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git pull origin main",
			Description: "Pull latest changes before pushing",
			Type:        "fix",
			Confidence:  0.80,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "git push --force-with-lease",
			Description: "Force push with safety check",
			Type:        "alternative",
			Confidence:  0.70,
		})

	case "access_denied":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git config --list | grep credential",
			Description: "Check Git credentials configuration",
			Type:        "info",
			Confidence:  0.85,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "ssh -T git@github.com",
			Description: "Test SSH connection to GitHub",
			Type:        "info",
			Confidence:  0.80,
		})

	case "repo_not_found":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git remote -v",
			Description: "Check remote repository URL",
			Type:        "info",
			Confidence:  0.90,
		})

	case "pathspec_no_match":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git status",
			Description: "Check current repository status",
			Type:        "info",
			Confidence:  0.90,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "ls -la",
			Description: "List files to verify paths",
			Type:        "info",
			Confidence:  0.80,
		})

	case "branch_exists":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git checkout <branch-name>",
			Description: "Switch to existing branch",
			Type:        "fix",
			Confidence:  0.85,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "git branch -d <branch-name>",
			Description: "Delete existing branch (if safe)",
			Type:        "alternative",
			Confidence:  0.70,
		})

	case "conflicts":
		suggestions = append(suggestions, &Suggestion{
			Command:     "git stash",
			Description: "Stash local changes temporarily",
			Type:        "fix",
			Confidence:  0.90,
		})
		suggestions = append(suggestions, &Suggestion{
			Command:     "git add -A && git commit -m 'Save local changes'",
			Description: "Commit local changes",
			Type:        "alternative",
			Confidence:  0.85,
		})
	}

	return suggestions, nil
}
