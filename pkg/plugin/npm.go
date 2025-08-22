package plugin

import (
	"context"
	"strings"
)

// NpmPlugin handles npm/Node.js errors
type NpmPlugin struct {
	*BasePlugin
}

// NewNpmPlugin creates a new npm plugin
func NewNpmPlugin() *NpmPlugin {
	return &NpmPlugin{
		BasePlugin: NewBasePlugin(
			"npm",
			"Handles npm and Node.js package manager errors",
			85,
		),
	}
}

// Match determines if this plugin should handle the command
func (p *NpmPlugin) Match(cmd string, output string, exitCode int) bool {
	if exitCode == 0 {
		return false
	}

	cmdLower := strings.ToLower(strings.TrimSpace(cmd))
	npmCommands := []string{"npm ", "npx ", "yarn ", "pnpm "}
	
	for _, npmCmd := range npmCommands {
		if strings.HasPrefix(cmdLower, npmCmd) {
			return true
		}
	}

	return false
}

// Suggest generates a suggestion for fixing the npm error
func (p *NpmPlugin) Suggest(ctx context.Context, cmd string, output string, exitCode int) (string, string, error) {
	outputLower := strings.ToLower(output)
	
	// Check for common npm errors
	
	// Command typos
	if correctedCmd, explanation := p.checkForTypos(cmd); correctedCmd != "" {
		return correctedCmd, explanation, nil
	}
	
	// Package not found
	if strings.Contains(outputLower, "404 not found") ||
	   strings.Contains(outputLower, "package not found") {
		return "npm search " + extractNpmPackageName(cmd), "Search for the correct package name", nil
	}
	
	// Permission errors
	if strings.Contains(outputLower, "eacces") || 
	   strings.Contains(outputLower, "permission denied") {
		if strings.Contains(cmd, "npm install -g") {
			return strings.Replace(cmd, "npm install -g", "sudo npm install -g", 1), 
				   "Adding sudo for global npm install", nil
		}
		return "npm config set prefix ~/.npm-global", 
			   "Configure npm to use local prefix to avoid permission issues", nil
	}
	
	// Network/registry issues
	if strings.Contains(outputLower, "network error") ||
	   strings.Contains(outputLower, "etimedout") ||
	   strings.Contains(outputLower, "enotfound") {
		return "npm config set registry https://registry.npmjs.org/", 
			   "Reset npm registry to default", nil
	}
	
	// Outdated npm
	if strings.Contains(outputLower, "npm warn") && strings.Contains(outputLower, "version") {
		return "npm install -g npm@latest", "Update npm to latest version", nil
	}
	
	// Missing package.json
	if strings.Contains(outputLower, "no such file") && strings.Contains(outputLower, "package.json") {
		return "npm init -y", "Initialize a new package.json file", nil
	}
	
	// Dependency conflicts
	if strings.Contains(outputLower, "peer dep") || strings.Contains(outputLower, "conflicting") {
		return "npm install --legacy-peer-deps", "Install with legacy peer dependency resolution", nil
	}
	
	// Cache issues
	if strings.Contains(outputLower, "cache") && strings.Contains(outputLower, "corrupt") {
		return "npm cache clean --force", "Clean npm cache", nil
	}
	
	// Lock file issues
	if strings.Contains(outputLower, "lockfileversion") ||
	   strings.Contains(outputLower, "package-lock.json") {
		return "rm package-lock.json && npm install", "Remove lock file and reinstall", nil
	}
	
	// Node version issues
	if strings.Contains(outputLower, "unsupported engine") ||
	   strings.Contains(outputLower, "node version") {
		return "node --version", "Check current Node.js version", nil
	}
	
	return "", "No specific npm suggestion available", nil
}

// checkForTypos checks for common typos in npm commands
func (p *NpmPlugin) checkForTypos(cmd string) (string, string) {
	corrections := map[string]string{
		"instal":      "install",
		"instll":      "install",
		"isntall":     "install",
		"installll":   "install",
		"intsall":     "install",
		"uninstal":    "uninstall",
		"unintsall":   "uninstall",
		"unisntall":   "uninstall",
		"udpate":      "update",
		"upate":       "update",
		"updtae":      "update",
		"update":      "update",
		"run":         "run",
		"rn":          "run",
		"runn":        "run",
		"start":       "start",
		"stat":        "start",
		"satrt":       "start",
		"startt":      "start",
		"test":        "test",
		"tst":         "test",
		"testt":       "test",
		"build":       "build",
		"buil":        "build",
		"buidl":       "build",
		"biuld":       "build",
		"audit":       "audit",
		"auidt":       "audit",
		"auidit":      "audit",
		"init":        "init",
		"ini":         "init",
		"initt":       "init",
		"publish":     "publish",
		"publsh":      "publish",
		"pubish":      "publish",
		"list":        "list",
		"lst":         "list",
		"lis":         "list",
		"info":        "info",
		"inf":         "info",
		"view":        "view",
		"vew":         "view",
		"search":      "search",
		"serach":      "search",
		"seach":       "search",
		"outdated":    "outdated",
		"outdate":     "outdated",
		"outdatd":     "outdated",
	}

	words := strings.Fields(cmd)
	changed := false
	
	for i, word := range words {
		cleanWord := strings.ToLower(word)
		if correction, exists := corrections[cleanWord]; exists {
			words[i] = correction
			changed = true
		}
	}

	if changed {
		newCmd := strings.Join(words, " ")
		return newCmd, "Fixed typos in npm command"
	}

	return "", ""
}

// extractNpmPackageName extracts package name from npm command
func extractNpmPackageName(cmd string) string {
	words := strings.Fields(cmd)
	
	for i, word := range words {
		if (word == "install" || word == "i" || word == "add" || word == "remove" || word == "uninstall") && i+1 < len(words) {
			return words[i+1]
		}
	}
	
	return ""
}
