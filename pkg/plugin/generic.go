package plugin

import (
	"context"
	"strings"
)

// GenericPlugin handles general shell command errors
type GenericPlugin struct {
	*BasePlugin
}

// NewGenericPlugin creates a new generic plugin
func NewGenericPlugin() *GenericPlugin {
	return &GenericPlugin{
		BasePlugin: NewBasePlugin(
			"generic",
			"Handles general shell command errors and common mistakes",
			10, // Low priority - fallback for other plugins
		),
	}
}

// Match determines if this plugin should handle the command
func (p *GenericPlugin) Match(cmd string, output string, exitCode int) bool {
	// Generic plugin matches any failed command as a fallback
	return exitCode != 0
}

// Suggest generates a suggestion for fixing the generic error
func (p *GenericPlugin) Suggest(ctx context.Context, cmd string, output string, exitCode int) (string, string, error) {
	outputLower := strings.ToLower(output)
	cmdLower := strings.ToLower(strings.TrimSpace(cmd))
	
	// Common shell errors
	
	// Command not found
	if strings.Contains(outputLower, "command not found") ||
	   strings.Contains(outputLower, "not found") {
		commandName := extractCommandName(cmd)
		if commandName != "" {
			return p.suggestCommandInstallation(commandName)
		}
	}
	
	// Permission denied
	if strings.Contains(outputLower, "permission denied") ||
	   strings.Contains(outputLower, "operation not permitted") {
		if !strings.HasPrefix(cmdLower, "sudo ") {
			return "sudo " + cmd, "Adding sudo for administrative privileges", nil
		}
	}
	
	// File or directory not found
	if strings.Contains(outputLower, "no such file or directory") {
		return "ls -la", "Check if the file or directory exists", nil
	}
	
	// Directory not empty
	if strings.Contains(outputLower, "directory not empty") {
		return "rm -rf " + extractDirectoryName(cmd), "Force remove directory and contents", nil
	}
	
	// Disk space issues
	if strings.Contains(outputLower, "no space left on device") {
		return "df -h", "Check disk space usage", nil
	}
	
	// Network issues
	if strings.Contains(outputLower, "network is unreachable") ||
	   strings.Contains(outputLower, "connection refused") {
		return "ping google.com", "Test network connectivity", nil
	}
	
	// Process/port issues
	if strings.Contains(outputLower, "address already in use") {
		return "netstat -tlnp", "Check which process is using the port", nil
	}
	
	// File already exists
	if strings.Contains(outputLower, "file exists") {
		return "ls -la " + extractFilename(cmd), "Check existing file", nil
	}
	
	// Syntax errors
	if strings.Contains(outputLower, "syntax error") {
		return "", "Check command syntax and try again", nil
	}
	
	// Check for common command typos
	if correctedCmd, explanation := p.checkForCommonTypos(cmd); correctedCmd != "" {
		return correctedCmd, explanation, nil
	}
	
	return "", "No specific suggestion available", nil
}

// suggestCommandInstallation suggests how to install a missing command
func (p *GenericPlugin) suggestCommandInstallation(commandName string) (string, string, error) {
	// Common command installation suggestions
	installations := map[string]string{
		"curl":      "sudo apt install curl",
		"wget":      "sudo apt install wget",
		"git":       "sudo apt install git",
		"vim":       "sudo apt install vim",
		"nano":      "sudo apt install nano",
		"htop":      "sudo apt install htop",
		"tree":      "sudo apt install tree",
		"zip":       "sudo apt install zip",
		"unzip":     "sudo apt install unzip",
		"ssh":       "sudo apt install openssh-client",
		"python":    "sudo apt install python3",
		"python3":   "sudo apt install python3",
		"pip":       "sudo apt install python3-pip",
		"pip3":      "sudo apt install python3-pip",
		"node":      "sudo apt install nodejs",
		"npm":       "sudo apt install npm",
		"docker":    "sudo apt install docker.io",
		"kubectl":   "sudo snap install kubectl --classic",
		"helm":      "sudo snap install helm --classic",
		"code":      "sudo snap install code --classic",
		"firefox":   "sudo apt install firefox",
		"chrome":    "wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | sudo apt-key add - && sudo apt install google-chrome-stable",
		"make":      "sudo apt install build-essential",
		"gcc":       "sudo apt install build-essential",
		"g++":       "sudo apt install build-essential",
		"java":      "sudo apt install openjdk-11-jdk",
		"javac":     "sudo apt install openjdk-11-jdk",
		"redis-cli": "sudo apt install redis-tools",
		"mysql":     "sudo apt install mysql-client",
		"psql":      "sudo apt install postgresql-client",
		"go":        "sudo apt install golang-go",
		"rust":      "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh",
	}

	if installation, exists := installations[strings.ToLower(commandName)]; exists {
		return installation, "Install " + commandName, nil
	}

	// Fallback: suggest using apt search
	return "apt search " + commandName, "Search for package containing " + commandName, nil
}

// checkForCommonTypos checks for common command typos
func (p *GenericPlugin) checkForCommonTypos(cmd string) (string, string) {
	corrections := map[string]string{
		"sl":     "ls",
		"dc":     "cd",
		"pwdd":   "pwd",
		"mkdri":  "mkdir",
		"mkidr":  "mkdir",
		"rmdri":  "rmdir",
		"mr":     "rm",
		"pc":     "cp",
		"vm":     "mv",
		"act":    "cat",
		"lses":   "less",
		"moer":   "more",
		"haed":   "head",
		"tial":   "tail",
		"gerp":   "grep",
		"grpe":   "grep",
		"fdin":   "find",
		"fidn":   "find",
		"whcih":  "which",
		"pot":    "top",
		"sp":     "ps",
		"kil":    "kill",
		"atr":    "tar",
		"shh":    "ssh",
		"psc":    "scp",
		"clar":   "clear",
		"cls":    "clear",
		"eixt":   "exit",
		"suod":   "sudo",
		"sduo":   "sudo",
	}

	words := strings.Fields(cmd)
	if len(words) == 0 {
		return "", ""
	}

	// Check the first word (command)
	firstWord := strings.ToLower(words[0])
	if correction, exists := corrections[firstWord]; exists && firstWord != correction {
		words[0] = correction
		newCmd := strings.Join(words, " ")
		return newCmd, "Fixed typo: '" + firstWord + "' → '" + correction + "'"
	}

	return "", ""
}

// extractCommandName extracts the command name from the command string
func extractCommandName(cmd string) string {
	words := strings.Fields(cmd)
	if len(words) > 0 {
		return words[0]
	}
	return ""
}

// extractDirectoryName extracts directory name from rm/rmdir commands
func extractDirectoryName(cmd string) string {
	words := strings.Fields(cmd)
	for i, word := range words {
		if (word == "rm" || word == "rmdir") && i+1 < len(words) {
			return words[i+1]
		}
	}
	return ""
}

// extractFilename extracts filename from various commands
func extractFilename(cmd string) string {
	words := strings.Fields(cmd)
	if len(words) > 1 {
		// Return the last argument which is likely the filename
		return words[len(words)-1]
	}
	return ""
}
