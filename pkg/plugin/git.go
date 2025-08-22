package plugin

import (
	"context"
	"regexp"
	"strings"
)

// GitPlugin handles git command errors
type GitPlugin struct {
	*BasePlugin
}

// NewGitPlugin creates a new git plugin
func NewGitPlugin() *GitPlugin {
	return &GitPlugin{
		BasePlugin: NewBasePlugin(
			"git",
			"Handles git command errors and common mistakes",
			90,
		),
	}
}

// Match determines if this plugin should handle the command
func (p *GitPlugin) Match(cmd string, output string, exitCode int) bool {
	if exitCode == 0 {
		return false
	}

	cmdLower := strings.ToLower(strings.TrimSpace(cmd))
	return strings.HasPrefix(cmdLower, "git ")
}

// Suggest generates a suggestion for fixing the git error
func (p *GitPlugin) Suggest(ctx context.Context, cmd string, output string, exitCode int) (string, string, error) {
	outputLower := strings.ToLower(output)
	
	// Check for common git errors
	
	// Command typos
	if correctedCmd, explanation := p.checkForTypos(cmd); correctedCmd != "" {
		return correctedCmd, explanation, nil
	}
	
	// No upstream branch
	if strings.Contains(outputLower, "no upstream branch") || 
	   strings.Contains(outputLower, "set-upstream") {
		branchName := extractCurrentBranch(output)
		if branchName == "" {
			branchName = "main" // fallback
		}
		return "git push --set-upstream origin " + branchName, 
			   "Set upstream branch for first push", nil
	}
	
	// Authentication failed
	if strings.Contains(outputLower, "authentication failed") ||
	   strings.Contains(outputLower, "permission denied") {
		return "git config --list | grep user", 
			   "Check git user configuration and authentication", nil
	}
	
	// Nothing to commit
	if strings.Contains(outputLower, "nothing to commit") {
		return "git status", "Check repository status for uncommitted changes", nil
	}
	
	// Not a git repository
	if strings.Contains(outputLower, "not a git repository") {
		return "git init", "Initialize a new git repository", nil
	}
	
	// Branch doesn't exist
	if strings.Contains(outputLower, "pathspec") && strings.Contains(outputLower, "did not match") {
		return "git branch -a", "List all available branches", nil
	}
	
	// Merge conflicts
	if strings.Contains(outputLower, "merge conflict") ||
	   strings.Contains(outputLower, "unmerged paths") {
		return "git status", "Check files with merge conflicts", nil
	}
	
	// Dirty working directory
	if strings.Contains(outputLower, "your local changes would be overwritten") {
		return "git stash", "Stash local changes before switching branches", nil
	}
	
	// Remote doesn't exist
	if strings.Contains(outputLower, "remote") && strings.Contains(outputLower, "does not exist") {
		return "git remote -v", "List configured remotes", nil
	}
	
	// Detached HEAD
	if strings.Contains(outputLower, "detached head") {
		return "git checkout main", "Return to main branch", nil
	}
	
	return "", "No specific git suggestion available", nil
}

// checkForTypos checks for common typos in git commands
func (p *GitPlugin) checkForTypos(cmd string) (string, string) {
	corrections := map[string]string{
		"checout":     "checkout",
		"checkot":     "checkout",
		"chekout":     "checkout",
		"chekcout":    "checkout",
		"committ":     "commit",
		"comit":       "commit",
		"comitt":      "commit",
		"commiy":      "commit",
		"statu":       "status",
		"stats":       "status",
		"stauts":      "status",
		"statsu":      "status",
		"pul":         "pull",
		"pll":         "pull",
		"puul":        "pull",
		"push":        "push",
		"psh":         "push",
		"pusj":        "push",
		"pusn":        "push",
		"branch":      "branch",
		"brach":       "branch",
		"branc":       "branch",
		"branchh":     "branch",
		"merge":       "merge",
		"merg":        "merge",
		"mrege":       "merge",
		"mereg":       "merge",
		"add":         "add",
		"ad":          "add",
		"addd":        "add",
		"remote":      "remote",
		"remot":       "remote",
		"rmote":       "remote",
		"remte":       "remote",
		"clone":       "clone",
		"clon":        "clone",
		"cloen":       "clone",
		"clonee":      "clone",
		"fetch":       "fetch",
		"fech":        "fetch",
		"fetsh":       "fetch",
		"log":         "log",
		"lgo":         "log",
		"logg":        "log",
		"diff":        "diff",
		"dif":         "diff",
		"difff":       "diff",
		"reset":       "reset",
		"rset":        "reset",
		"rese":        "reset",
		"resett":      "reset",
		"rebase":      "rebase",
		"rebas":       "rebase",
		"reabse":      "rebase",
		"stash":       "stash",
		"stas":        "stash",
		"stsh":        "stash",
		"stasj":       "stash",
		"tag":         "tag",
		"tg":          "tag",
		"tagg":        "tag",
	}

	words := strings.Fields(cmd)
	if len(words) < 2 {
		return "", ""
	}

	// Check the git subcommand (second word)
	if words[0] == "git" && len(words) > 1 {
		subcommand := words[1]
		if correction, exists := corrections[strings.ToLower(subcommand)]; exists {
			words[1] = correction
			newCmd := strings.Join(words, " ")
			return newCmd, "Fixed typo: '" + subcommand + "' → '" + correction + "'"
		}
	}

	return "", ""
}

// extractCurrentBranch extracts the current branch name from git output
func extractCurrentBranch(output string) string {
	// Try to extract branch name from various git messages
	patterns := []string{
		`current branch (.+?)[\s\n]`,
		`branch '(.+?)'`,
		`on branch (.+?)[\s\n]`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(output)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}
