package plugins

import (
	"fmt"
	"strings"

	"github.com/ayush-1/logaid/internal/config"
	"github.com/ayush-1/logaid/internal/logger"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Match(cmd string, output string) bool     // When to trigger this plugin
	Suggest(cmd string, output string) string // Generate suggestion
	Name() string                             // Plugin identifier
}

// LoadAllPlugins loads all enabled plugins
func LoadAllPlugins() []Plugin {
	var plugins []Plugin

	if config.AppConfig == nil {
		return plugins
	}

	enabledPlugins := strings.Split(config.AppConfig.EnablePlugins, ",")
	enabledMap := make(map[string]bool)
	for _, plugin := range enabledPlugins {
		enabledMap[strings.TrimSpace(plugin)] = true
	}

	// Load built-in plugins
	if enabledMap["apt"] {
		plugins = append(plugins, &AptPlugin{})
		logger.Debug("Loaded apt plugin")
	}

	if enabledMap["npm"] {
		plugins = append(plugins, &NpmPlugin{})
		logger.Debug("Loaded npm plugin")
	}

	if enabledMap["git"] {
		plugins = append(plugins, &GitPlugin{})
		logger.Debug("Loaded git plugin")
	}

	if enabledMap["docker"] {
		plugins = append(plugins, &DockerPlugin{})
		logger.Debug("Loaded docker plugin")
	}

	if enabledMap["pip"] {
		plugins = append(plugins, &PipPlugin{})
		logger.Debug("Loaded pip plugin")
	}

	if enabledMap["systemctl"] {
		plugins = append(plugins, &SystemctlPlugin{})
		logger.Debug("Loaded systemctl plugin")
	}

	logger.Info(fmt.Sprintf("Loaded %d plugins", len(plugins)))
	return plugins
}

// Helper function to check if output contains any of the given strings
func containsAny(text string, patterns []string) bool {
	lowerText := strings.ToLower(text)
	for _, pattern := range patterns {
		if strings.Contains(lowerText, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}
