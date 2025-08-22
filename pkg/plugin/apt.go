package plugin

import (
	"context"
	"regexp"
	"strings"
)

// AptPlugin handles apt package manager errors
type AptPlugin struct {
	*BasePlugin
}

// NewAptPlugin creates a new apt plugin
func NewAptPlugin() *AptPlugin {
	return &AptPlugin{
		BasePlugin: NewBasePlugin(
			"apt",
			"Handles apt package manager errors and typos",
			80,
		),
	}
}

// Match determines if this plugin should handle the command
func (p *AptPlugin) Match(cmd string, output string, exitCode int) bool {
	if exitCode == 0 {
		return false
	}

	// Check if it's an apt command
	aptCommands := []string{"apt", "apt-get", "apt-cache", "aptitude"}
	cmdLower := strings.ToLower(cmd)
	
	for _, aptCmd := range aptCommands {
		if strings.HasPrefix(cmdLower, aptCmd) {
			return true
		}
	}

	return false
}

// Suggest generates a suggestion for fixing the apt error
func (p *AptPlugin) Suggest(ctx context.Context, cmd string, output string, exitCode int) (string, string, error) {
	// Check for common apt errors and provide specific suggestions
	outputLower := strings.ToLower(output)
	
	// Package not found error
	if strings.Contains(outputLower, "unable to locate package") {
		packageName := extractPackageName(cmd, output)
		if packageName != "" {
			return p.suggestPackageAlternatives(cmd, packageName)
		}
	}
	
	// Permission denied
	if strings.Contains(outputLower, "permission denied") || 
	   strings.Contains(outputLower, "operation not permitted") {
		if !strings.HasPrefix(strings.TrimSpace(cmd), "sudo") {
			return "sudo " + cmd, "Adding sudo for administrative privileges", nil
		}
	}
	
	// Lock file error
	if strings.Contains(outputLower, "could not get lock") ||
	   strings.Contains(outputLower, "dpkg was interrupted") {
		return "sudo dpkg --configure -a", "Fixing interrupted dpkg process", nil
	}
	
	// Update required
	if strings.Contains(outputLower, "no installation candidate") {
		return "sudo apt update && " + cmd, "Updating package lists first", nil
	}
	
	// Generic typo check
	return p.checkForTypos(cmd)
}

// extractPackageName extracts the package name from the command or error output
func extractPackageName(cmd, output string) string {
	// Try to extract from error message first
	re := regexp.MustCompile(`unable to locate package (.+?)(?:\s|$)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	
	// Extract from command
	parts := strings.Fields(cmd)
	for i, part := range parts {
		if (part == "install" || part == "remove" || part == "purge" || part == "show") && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	
	return ""
}

// suggestPackageAlternatives suggests alternative package names
func (p *AptPlugin) suggestPackageAlternatives(cmd, packageName string) (string, string, error) {
	// Common package name corrections
	corrections := map[string]string{
		"rediscli":     "redis-tools",
		"redis-cli":    "redis-tools",
		"python":       "python3",
		"pip":          "python3-pip",
		"nodejs":       "node",
		"jdk":          "openjdk-11-jdk",
		"java":         "openjdk-11-jdk",
		"gcc":          "build-essential",
		"make":         "build-essential",
		"curl":         "curl",
		"wget":         "wget",
		"git":          "git",
		"vim":          "vim",
		"nano":         "nano",
		"htop":         "htop",
		"tree":         "tree",
		"zip":          "zip",
		"unzip":        "unzip",
		"ssh":          "openssh-client",
		"sshd":         "openssh-server",
		"docker":       "docker.io",
		"docker-ce":    "docker.io",
		"kubectl":      "kubectl",
		"terraform":    "terraform",
	}

	if corrected, exists := corrections[strings.ToLower(packageName)]; exists {
		newCmd := strings.Replace(cmd, packageName, corrected, 1)
		return newCmd, "Corrected package name from '" + packageName + "' to '" + corrected + "'", nil
	}

	// Suggest using apt search to find the correct package
	searchCmd := "apt search " + packageName
	return searchCmd, "Search for the correct package name", nil
}

// checkForTypos checks for common typos in apt commands
func (p *AptPlugin) checkForTypos(cmd string) (string, string, error) {
	corrections := map[string]string{
		"instal":      "install",
		"instll":      "install",
		"isntall":     "install",
		"update":      "update",
		"upate":       "update",
		"updtae":      "update",
		"remove":      "remove",
		"remov":       "remove",
		"rmove":       "remove",
		"purge":       "purge",
		"purg":        "purge",
		"search":      "search",
		"serach":      "search",
		"seach":       "search",
		"show":        "show",
		"sho":         "show",
		"upgrade":     "upgrade",
		"upgrad":      "upgrade",
		"upgade":      "upgrade",
		"autoremove":  "autoremove",
		"autormeove":  "autoremove",
		"autoclean":   "autoclean",
		"autoclen":    "autoclean",
	}

	words := strings.Fields(cmd)
	changed := false
	
	for i, word := range words {
		if correction, exists := corrections[strings.ToLower(word)]; exists {
			words[i] = correction
			changed = true
		}
	}

	if changed {
		newCmd := strings.Join(words, " ")
		return newCmd, "Fixed typos in command", nil
	}

	return "", "No specific suggestion available", nil
}
