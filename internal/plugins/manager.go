package plugins

import (
	"context"
	"regexp"
	"strings"
)

// Plugin represents a LogAid plugin interface
type Plugin interface {
	// Name returns the plugin name
	Name() string

	// Description returns a brief description of the plugin
	Description() string

	// CanHandle checks if the plugin can handle the given command/error
	CanHandle(command, stderr string, exitCode int) bool

	// Analyze analyzes the error and provides specific insights
	Analyze(ctx context.Context, command, stderr string, exitCode int) (*Analysis, error)

	// GetSuggestions provides specific suggestions for the error
	GetSuggestions(ctx context.Context, analysis *Analysis) ([]*Suggestion, error)
}

// Analysis represents the analysis result from a plugin
type Analysis struct {
	Plugin      string            `json:"plugin"`
	ErrorType   string            `json:"error_type"`
	Severity    string            `json:"severity"`
	Category    string            `json:"category"`
	Context     map[string]string `json:"context"`
	Confidence  float64           `json:"confidence"`
	Description string            `json:"description"`
}

// Suggestion represents a specific suggestion from a plugin
type Suggestion struct {
	Command     string            `json:"command"`
	Description string            `json:"description"`
	Type        string            `json:"type"` // "fix", "alternative", "info"
	Confidence  float64           `json:"confidence"`
	Context     map[string]string `json:"context"`
}

// Manager manages all plugins
type Manager struct {
	plugins map[string]Plugin
	enabled map[string]bool
}

// NewManager creates a new plugin manager
func NewManager(enabledPlugins []string) *Manager {
	manager := &Manager{
		plugins: make(map[string]Plugin),
		enabled: make(map[string]bool),
	}

	// Register built-in plugins
	manager.registerBuiltinPlugins()

	// Enable specified plugins
	for _, name := range enabledPlugins {
		manager.enabled[strings.TrimSpace(name)] = true
	}

	return manager
}

// registerBuiltinPlugins registers all built-in plugins
func (m *Manager) registerBuiltinPlugins() {
	plugins := []Plugin{
		NewGitPlugin(),
		NewDockerPlugin(),
		NewNPMPlugin(),
		NewAptPlugin(),
		NewKubernetesPlugin(),
		NewGenericPlugin(),
	}

	for _, plugin := range plugins {
		m.plugins[plugin.Name()] = plugin
	}
}

// RegisterPlugin registers a custom plugin
func (m *Manager) RegisterPlugin(plugin Plugin) {
	m.plugins[plugin.Name()] = plugin
}

// AnalyzeError analyzes an error using all applicable plugins
func (m *Manager) AnalyzeError(ctx context.Context, command, stderr string, exitCode int) ([]*Analysis, error) {
	var analyses []*Analysis

	for name, plugin := range m.plugins {
		// Skip if plugin is not enabled
		if !m.enabled[name] {
			continue
		}

		// Check if plugin can handle this error
		if !plugin.CanHandle(command, stderr, exitCode) {
			continue
		}

		// Analyze the error
		analysis, err := plugin.Analyze(ctx, command, stderr, exitCode)
		if err != nil {
			// Log error but continue with other plugins
			continue
		}

		if analysis != nil {
			analyses = append(analyses, analysis)
		}
	}

	return analyses, nil
}

// GetSuggestions gets suggestions from all applicable plugins
func (m *Manager) GetSuggestions(ctx context.Context, command, stderr string, exitCode int) ([]*Suggestion, error) {
	var allSuggestions []*Suggestion

	for name, plugin := range m.plugins {
		// Skip if plugin is not enabled
		if !m.enabled[name] {
			continue
		}

		// Check if plugin can handle this error
		if !plugin.CanHandle(command, stderr, exitCode) {
			continue
		}

		// Analyze first
		analysis, err := plugin.Analyze(ctx, command, stderr, exitCode)
		if err != nil || analysis == nil {
			continue
		}

		// Get suggestions
		suggestions, err := plugin.GetSuggestions(ctx, analysis)
		if err != nil {
			continue
		}

		allSuggestions = append(allSuggestions, suggestions...)
	}

	return allSuggestions, nil
}

// GetPluginInfo returns information about all registered plugins
func (m *Manager) GetPluginInfo() map[string]PluginInfo {
	info := make(map[string]PluginInfo)

	for name, plugin := range m.plugins {
		info[name] = PluginInfo{
			Name:        plugin.Name(),
			Description: plugin.Description(),
			Enabled:     m.enabled[name],
		}
	}

	return info
}

// PluginInfo contains information about a plugin
type PluginInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// BasePlugin provides common functionality for plugins
type BasePlugin struct {
	name        string
	description string
	patterns    []*regexp.Regexp
}

// Name returns the plugin name
func (b *BasePlugin) Name() string {
	return b.name
}

// Description returns the plugin description
func (b *BasePlugin) Description() string {
	return b.description
}

// MatchesPattern checks if stderr matches any of the plugin's patterns
func (b *BasePlugin) MatchesPattern(stderr string) bool {
	for _, pattern := range b.patterns {
		if pattern.MatchString(stderr) {
			return true
		}
	}
	return false
}

// MatchesCommand checks if the command matches expected patterns
func (b *BasePlugin) MatchesCommand(command string, expectedCommands []string) bool {
	cmdParts := strings.Fields(command)
	if len(cmdParts) == 0 {
		return false
	}

	mainCmd := cmdParts[0]
	for _, expected := range expectedCommands {
		if strings.Contains(mainCmd, expected) {
			return true
		}
	}
	return false
}

// ExtractCommand extracts the main command from a command string
func ExtractCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

// ExtractArgs extracts arguments from a command string
func ExtractArgs(command string) []string {
	parts := strings.Fields(command)
	if len(parts) <= 1 {
		return []string{}
	}
	return parts[1:]
}
