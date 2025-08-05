package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/ayush-1/logaid/internal/ai"
)

// AptPlugin handles APT package manager errors with AI-powered suggestions
type AptPlugin struct{}

func (p *AptPlugin) Name() string {
	return "apt"
}

// Match checks if this plugin should handle the command/output
func (p *AptPlugin) Match(cmd string, output string) bool {
	// Check if command uses apt/apt-get
	if !strings.Contains(strings.ToLower(cmd), "apt") {
		return false
	}

	// Check for common apt errors
	aptErrors := []string{
		"unable to locate package",
		"package not found",
		"e: could not get lock",
		"e: package",
		"has no installation candidate",
		"depends:",
		"unmet dependencies",
		"permission denied",
		"command not found",
		"broken packages",
		"held broken packages",
		"404 not found",
		"signature verification failed",
		"repository does not have a release file",
	}

	return containsAny(output, aptErrors)
}

// Suggest generates an AI-powered suggestion for the error
func (p *AptPlugin) Suggest(cmd string, output string) string {
	// First try manual corrections for speed
	if quickFix := p.getQuickFix(cmd, output); quickFix != "" {
		return quickFix
	}

	// Use AI for complex suggestions
	return p.getAISuggestion(cmd, output)
}

// getQuickFix provides immediate fixes for common issues
func (p *AptPlugin) getQuickFix(cmd string, output string) string {
	outputLower := strings.ToLower(output)

	// Handle lock errors
	if strings.Contains(outputLower, "could not get lock") {
		return "sudo killall apt apt-get dpkg && sudo dpkg --configure -a && sudo apt update && " + cmd
	}

	// Handle permission errors
	if strings.Contains(outputLower, "permission denied") && !strings.Contains(cmd, "sudo") {
		return "sudo " + cmd
	}

	// Handle update needed
	if strings.Contains(outputLower, "404") || strings.Contains(outputLower, "failed to fetch") {
		return "sudo apt update && " + cmd
	}

	// Common package name corrections
	if strings.Contains(outputLower, "unable to locate package") {
		parts := strings.Fields(cmd)
		for i, part := range parts {
			if (part == "install" || part == "show" || part == "search") && i+1 < len(parts) {
				packageName := parts[i+1]
				if correction := p.getPackageCorrection(packageName); correction != "" {
					return strings.Replace(cmd, packageName, correction, 1)
				}
			}
		}
	}

	return ""
}

// getPackageCorrection provides manual corrections for common package name typos
func (p *AptPlugin) getPackageCorrection(packageName string) string {
	corrections := map[string]string{
		"rediscli":     "redis-tools",
		"redis-cli":    "redis-tools",
		"redisclient":  "redis-tools",
		"mysql":        "mysql-server",
		"mysqlserver":  "mysql-server",
		"nodejs":       "nodejs npm",
		"node":         "nodejs npm",
		"python":       "python3",
		"pip":          "python3-pip",
		"python-pip":   "python3-pip",
		"docker":       "docker.io",
		"dockerio":     "docker.io",
		"java":         "openjdk-11-jdk",
		"jdk":          "openjdk-11-jdk",
		"gcc":          "build-essential",
		"make":         "build-essential",
		"cmake":        "cmake build-essential",
		"vim":          "vim-gtk3",
		"emacs":        "emacs-gtk",
		"firefox":      "firefox-esr",
		"chrome":       "google-chrome-stable",
		"chromium":     "chromium-browser",
		"vscode":       "code",
		"visualstudio": "code",
		"git":          "git-all",
		"curl":         "curl wget",
		"ssh":          "openssh-client openssh-server",
		"sshd":         "openssh-server",
		"nginx":        "nginx-full",
		"apache":       "apache2",
		"httpd":        "apache2",
		"postgres":     "postgresql postgresql-contrib",
		"postgresql":   "postgresql postgresql-contrib",
		"sqlite":       "sqlite3",
		"htop":         "htop",
		"top":          "htop",
		"nano":         "nano",
		"unzip":        "unzip zip",
		"tar":          "tar gzip",
		"screen":       "screen tmux",
		"tmux":         "tmux",
		"tree":         "tree",
		"less":         "less",
		"jq":           "jq",
		"net-tools":    "net-tools",
		"ifconfig":     "net-tools",
		"netstat":      "net-tools",
		"telnet":       "telnet",
		"ftp":          "ftp",
		"rsync":        "rsync",
		"scp":          "openssh-client",
	}

	return corrections[strings.ToLower(packageName)]
}

// getAISuggestion uses AI to generate intelligent suggestions
func (p *AptPlugin) getAISuggestion(cmd string, output string) string {
	prompt := p.buildAIPrompt(cmd, output)

	ctx := context.Background()
	suggestion, err := ai.GetSuggestion(ctx, prompt)
	if err != nil {
		// Fallback to generic suggestion
		return "sudo apt update && apt search <package-name> && " + cmd
	}

	return suggestion
}

// buildAIPrompt creates a detailed prompt for the AI
func (p *AptPlugin) buildAIPrompt(cmd string, output string) string {
	return fmt.Sprintf(`
You are an expert Linux system administrator specializing in APT package management on Debian/Ubuntu systems.

CONTEXT:
- User executed command: %s
- Command output/error: %s
- System: Debian/Ubuntu with APT package manager
- Goal: Provide the EXACT corrected command to fix the issue

TASK:
Analyze the command and error, then provide a single, executable command that will resolve the issue.

RULES:
1. Return ONLY the corrected command, no explanations
2. Use proper APT syntax and package names
3. Include sudo if needed for permissions
4. Handle common issues: typos, missing packages, lock files, repository updates
5. If package doesn't exist, suggest the closest alternative
6. For dependency issues, suggest the complete fix
7. Always prioritize safety and standard practices

COMMON APT PATTERNS TO CONSIDER:
- Package name typos (redis-cli → redis-tools)
- Missing sudo for install operations
- Need to update package lists first
- Lock file conflicts requiring cleanup
- Missing repositories or keys
- Dependency conflicts
- Alternative package names

EXAMPLES:
- Input: "apt install rediscli" + "Unable to locate package rediscli"
- Output: "sudo apt install redis-tools"

- Input: "apt update" + "Permission denied"  
- Output: "sudo apt update"

Provide the corrected command:`, cmd, output)
}
