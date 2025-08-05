package plugins

import (
	"strings"
)

// GitPlugin handles Git command errors
type GitPlugin struct{}

func (p *GitPlugin) Name() string {
	return "git"
}

func (p *GitPlugin) Match(cmd string, output string) bool {
	// Check if this is a git command
	if !strings.HasPrefix(cmd, "git ") {
		return false
	}

	// Check for common git errors
	errorPatterns := []string{
		"git: command not found",
		"is not a git command",
		"unknown option",
		"invalid command",
		"not a git repository",
		"pathspec",
		"did not match any file",
		"fatal:",
		"error:",
	}

	return containsAny(output, errorPatterns)
}

func (p *GitPlugin) Suggest(cmd string, output string) string {
	// Common git command typos
	commandCorrections := map[string]string{
		"checout":  "checkout",
		"checkuot": "checkout",
		"chekout":  "checkout",
		"committ":  "commit",
		"comit":    "commit",
		"stauts":   "status",
		"stats":    "status",
		"stat":     "status",
		"brach":    "branch",
		"branc":    "branch",
		"branh":    "branch",
		"pul":      "pull",
		"pus":      "push",
		"pussh":    "push",
		"fetch":    "fetch",
		"merg":     "merge",
		"merge":    "merge",
		"rebase":   "rebase",
		"rebas":    "rebase",
		"clon":     "clone",
		"cloen":    "clone",
		"ad":       "add",
		"remot":    "remote",
		"reste":    "reset",
		"resett":   "reset",
		"dif":      "diff",
		"lo":       "log",
		"sho":      "show",
		"tag":      "tag",
		"stash":    "stash",
		"stas":     "stash",
	}

	// Parse the git command
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		return ""
	}

	gitCommand := parts[1]

	// Check for direct command corrections
	if correction, exists := commandCorrections[gitCommand]; exists {
		return strings.Replace(cmd, "git "+gitCommand, "git "+correction, 1)
	}

	// Handle specific error cases
	if strings.Contains(output, "not a git repository") {
		return "git init"
	}

	if strings.Contains(output, "pathspec") && strings.Contains(output, "did not match") {
		// Suggest git branch to show available branches
		return "git branch -a"
	}

	if strings.Contains(cmd, "checkout") && strings.Contains(output, "pathspec") {
		// Extract branch name and suggest creating it
		for i, part := range parts {
			if part == "checkout" && i+1 < len(parts) {
				branchName := parts[i+1]
				return "git checkout -b " + branchName
			}
		}
	}

	return ""
}
