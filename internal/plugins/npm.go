package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/ayushsharma-1/LogAid/internal/ai"
)

// NpmPlugin handles NPM command errors with AI-powered suggestions
type NpmPlugin struct{}

func (p *NpmPlugin) Name() string {
	return "npm"
}

// Match checks if this plugin should handle the command/output
func (p *NpmPlugin) Match(cmd string, output string) bool {
	// Check if command uses npm
	if !strings.Contains(strings.ToLower(cmd), "npm") {
		return false
	}

	// Check for common npm errors
	npmErrors := []string{
		"unknown command:",
		"npm err! 404",
		"not found",
		"eacces: permission denied",
		"network request",
		"enotfound",
		"timeout",
		"npm err! missing script:",
		"cannot resolve dependency:",
		"peer dep warning",
		"deprecated warning",
		"audit found",
		"vulnerabilities found",
		"npm err! code enoent",
		"npm err! errno -4058",
		"npm err! path",
		"operation not permitted",
	}

	return containsAny(output, npmErrors)
}

// Suggest generates an AI-powered suggestion for the error
func (p *NpmPlugin) Suggest(cmd string, output string) string {
	// First try manual corrections for speed
	if quickFix := p.getQuickFix(cmd, output); quickFix != "" {
		return quickFix
	}

	// Use AI for complex suggestions
	return p.getAISuggestion(cmd, output)
}

// getQuickFix provides immediate fixes for common issues
func (p *NpmPlugin) getQuickFix(cmd string, output string) string {
	outputLower := strings.ToLower(output)

	// Handle permission errors
	if strings.Contains(outputLower, "eacces") || strings.Contains(outputLower, "permission denied") {
		if strings.Contains(cmd, "-g") && !strings.Contains(cmd, "sudo") {
			return "sudo " + cmd
		}
	}

	// Handle command typos
	if strings.Contains(outputLower, "unknown command:") {
		return p.correctNpmCommand(cmd)
	}

	// Handle package not found (404 errors)
	if strings.Contains(outputLower, "404") && strings.Contains(outputLower, "not found") {
		return p.correctPackageName(cmd, output)
	}

	// Handle missing script
	if strings.Contains(outputLower, "missing script:") {
		return p.suggestScriptCommand(cmd, output)
	}

	return ""
}

// correctNpmCommand fixes common NPM command typos
func (p *NpmPlugin) correctNpmCommand(cmd string) string {
	corrections := map[string]string{
		"instal":    "install",
		"instll":    "install",
		"insall":    "install",
		"isntall":   "install",
		"instlal":   "install",
		"intall":    "install",
		"i":         "install",
		"stat":      "start",
		"strt":      "start",
		"str":       "start",
		"tes":       "test",
		"tst":       "test",
		"ru":        "run",
		"rn":        "run",
		"updat":     "update",
		"updte":     "update",
		"upgrd":     "upgrade",
		"uninsta":   "uninstall",
		"uninstal":  "uninstall",
		"remove":    "uninstall",
		"rm":        "uninstall",
		"lst":       "list",
		"ls":        "list",
		"info":      "info",
		"inf":       "info",
		"view":      "view",
		"vw":        "view",
		"search":    "search",
		"find":      "search",
		"audit":     "audit",
		"audt":      "audit",
		"outdated":  "outdated",
		"outdate":   "outdated",
		"init":      "init",
		"int":       "init",
		"publish":   "publish",
		"pub":       "publish",
		"unpublish": "unpublish",
		"version":   "version",
		"ver":       "version",
		"link":      "link",
		"lnk":       "link",
		"unlink":    "unlink",
		"config":    "config",
		"conf":      "config",
		"cache":     "cache",
		"chche":     "cache",
	}

	parts := strings.Fields(cmd)
	if len(parts) >= 2 {
		command := parts[1]
		if correction, exists := corrections[command]; exists {
			parts[1] = correction
			return strings.Join(parts, " ")
		}
	}

	return cmd
}

// correctPackageName fixes common package name typos
func (p *NpmPlugin) correctPackageName(cmd string, output string) string {
	packageCorrections := map[string]string{
		// Popular packages with common typos
		"expres":       "express",
		"exprees":      "express",
		"expresss":     "express",
		"lodas":        "lodash",
		"lodsh":        "lodash",
		"lodassh":      "lodash",
		"reac":         "react",
		"react":        "react",
		"reactt":       "react",
		"axio":         "axios",
		"axois":        "axios",
		"axioss":       "axios",
		"momen":        "moment",
		"momnet":       "moment",
		"momentt":      "moment",
		"nod-fetch":    "node-fetch",
		"node-fech":    "node-fetch",
		"nodefetch":    "node-fetch",
		"cheerio":      "cheerio",
		"cherio":       "cheerio",
		"cheeio":       "cheerio",
		"socket.i":     "socket.io",
		"socketio":     "socket.io",
		"socket-io":    "socket.io",
		"uuid":         "uuid",
		"uui":          "uuid",
		"uuuid":        "uuid",
		"bcryp":        "bcrypt",
		"bcrypt":       "bcrypt",
		"bcryptjs":     "bcryptjs",
		"jsonwebtoken": "jsonwebtoken",
		"jwt":          "jsonwebtoken",
		"mongoose":     "mongoose",
		"mongose":      "mongoose",
		"mungoose":     "mongoose",
		"sequelize":    "sequelize",
		"sequlize":     "sequelize",
		"sequeize":     "sequelize",
		"cors":         "cors",
		"cor":          "cors",
		"corss":        "cors",
		"helmet":       "helmet",
		"helmt":        "helmet",
		"helnet":       "helmet",
		"morgan":       "morgan",
		"morga":        "morgan",
		"morganr":      "morgan",
		"nodemon":      "nodemon",
		"nodmon":       "nodemon",
		"nodemn":       "nodemon",
		"pm2":          "pm2",
		"pm":           "pm2",
		"dotenv":       "dotenv",
		"dotev":        "dotenv",
		"dontenv":      "dotenv",
		"chalk":        "chalk",
		"chlk":         "chalk",
		"chalck":       "chalk",
		"commander":    "commander",
		"comander":     "commander",
		"comandr":      "commander",
		"inquirer":     "inquirer",
		"inquierer":    "inquirer",
		"inquirr":      "inquirer",
		"fs-extra":     "fs-extra",
		"fs-ext":       "fs-extra",
		"fsextra":      "fs-extra",
		"glob":         "glob",
		"globb":        "glob",
		"globo":        "glob",
		"rimraf":       "rimraf",
		"rimaf":        "rimraf",
		"rmraf":        "rimraf",
	}

	// Try to extract package name and correct it
	parts := strings.Fields(cmd)
	for i, part := range parts {
		if part == "install" || part == "i" {
			if i+1 < len(parts) {
				packageName := parts[i+1]
				// Remove flags and get clean package name
				cleanPackage := strings.Split(packageName, "@")[0]
				if correction, exists := packageCorrections[cleanPackage]; exists {
					parts[i+1] = strings.Replace(packageName, cleanPackage, correction, 1)
					return strings.Join(parts, " ")
				}
			}
		}
	}

	return cmd
}

// suggestScriptCommand suggests npm run scripts
func (p *NpmPlugin) suggestScriptCommand(cmd string, output string) string {
	// Common script names
	commonScripts := []string{
		"start", "dev", "build", "test", "lint", "serve", "watch", "deploy",
		"production", "development", "server", "client", "clean", "install",
		"update", "postinstall", "preinstall", "prebuild", "postbuild",
		"prestart", "poststart", "pretest", "posttest",
	}

	parts := strings.Fields(cmd)
	if len(parts) >= 2 && (parts[1] == "run" && len(parts) >= 3) {
		scriptName := parts[2]
		// Suggest similar script names
		for _, common := range commonScripts {
			if strings.Contains(common, scriptName) || strings.Contains(scriptName, common) {
				parts[2] = common
				return strings.Join(parts, " ")
			}
		}
	}

	return "npm run-script --help # List available scripts"
}

// getAISuggestion uses AI to generate intelligent suggestions
func (p *NpmPlugin) getAISuggestion(cmd string, output string) string {
	prompt := p.buildAIPrompt(cmd, output)

	ctx := context.Background()
	suggestion, err := ai.GetSuggestion(ctx, prompt)
	if err != nil {
		// Fallback to generic suggestion
		return "npm --help # Check the correct NPM command syntax"
	}

	return suggestion
}

// buildAIPrompt creates a detailed prompt for the AI
func (p *NpmPlugin) buildAIPrompt(cmd string, output string) string {
	return fmt.Sprintf(`
You are an expert Node.js and NPM package manager specialist.

CONTEXT:
- User executed command: %s
- Command output/error: %s
- System: Node.js environment with NPM package manager
- Goal: Provide the EXACT corrected command to fix the issue

TASK:
Analyze the command and error, then provide a single, executable command that will resolve the issue.

RULES:
1. Return ONLY the corrected command, no explanations
2. Use proper NPM syntax and package names
3. Include sudo if needed for global installations
4. Handle common issues: typos, missing packages, permission errors, network issues
5. If package doesn't exist, suggest the closest alternative
6. For script issues, suggest the complete fix
7. Always prioritize safety and standard practices

COMMON NPM PATTERNS TO CONSIDER:
- Command typos (instal → install, stat → start, ru → run)
- Package name typos (expres → express, lodas → lodash, reac → react)
- Missing sudo for global installations
- Network connectivity issues
- Registry configuration problems
- Script name typos
- Version conflicts
- Permission issues with node_modules

EXAMPLES:
- Input: "npm instal express" + "Unknown command: instal"
- Output: "npm install express"

- Input: "npm install expres" + "404 Not Found - GET https://registry.npmjs.org/expres"
- Output: "npm install express"

- Input: "npm install -g typescript" + "EACCES: permission denied"
- Output: "sudo npm install -g typescript"

Provide the corrected command:`, cmd, output)
}
