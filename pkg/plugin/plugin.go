package plugin

import (
	"context"
)

// Plugin defines the interface for LogAid plugins
type Plugin interface {
	// Name returns the plugin name
	Name() string
	
	// Description returns a brief description of what the plugin handles
	Description() string
	
	// Match determines if this plugin should handle the given command and output
	Match(cmd string, output string, exitCode int) bool
	
	// Suggest generates a suggestion for fixing the error
	// Returns suggested command, explanation, and any error
	Suggest(ctx context.Context, cmd string, output string, exitCode int) (string, string, error)
	
	// Priority returns the plugin priority (higher numbers = higher priority)
	Priority() int
	
	// Enabled returns whether the plugin is enabled
	Enabled() bool
	
	// SetEnabled enables or disables the plugin
	SetEnabled(enabled bool)
}

// Manager manages plugin loading and execution
type Manager struct {
	plugins []Plugin
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins: make([]Plugin, 0),
	}
}

// Register registers a plugin with the manager
func (m *Manager) Register(plugin Plugin) {
	m.plugins = append(m.plugins, plugin)
}

// GetPlugins returns all registered plugins
func (m *Manager) GetPlugins() []Plugin {
	return m.plugins
}

// FindMatchingPlugin finds the first plugin that matches the given command and output
func (m *Manager) FindMatchingPlugin(cmd string, output string, exitCode int) Plugin {
	var bestMatch Plugin
	bestPriority := -1

	for _, plugin := range m.plugins {
		if !plugin.Enabled() {
			continue
		}
		
		if plugin.Match(cmd, output, exitCode) {
			if plugin.Priority() > bestPriority {
				bestMatch = plugin
				bestPriority = plugin.Priority()
			}
		}
	}

	return bestMatch
}

// LoadBuiltinPlugins loads all built-in plugins
func (m *Manager) LoadBuiltinPlugins() {
	plugins := []Plugin{
		NewAptPlugin(),
		NewGitPlugin(),
		NewNpmPlugin(),
		NewDockerPlugin(),
		NewKubernetesPlugin(),
		NewGenericPlugin(),
	}

	for _, plugin := range plugins {
		m.Register(plugin)
	}
}

// BasePlugin provides common functionality for plugins
type BasePlugin struct {
	name        string
	description string
	priority    int
	enabled     bool
}

// NewBasePlugin creates a new base plugin
func NewBasePlugin(name, description string, priority int) *BasePlugin {
	return &BasePlugin{
		name:        name,
		description: description,
		priority:    priority,
		enabled:     true,
	}
}

// Name returns the plugin name
func (p *BasePlugin) Name() string {
	return p.name
}

// Description returns the plugin description
func (p *BasePlugin) Description() string {
	return p.description
}

// Priority returns the plugin priority
func (p *BasePlugin) Priority() int {
	return p.priority
}

// Enabled returns whether the plugin is enabled
func (p *BasePlugin) Enabled() bool {
	return p.enabled
}

// SetEnabled enables or disables the plugin
func (p *BasePlugin) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// SuggestionRequest represents a request for an AI suggestion
type SuggestionRequest struct {
	Command    string
	Output     string
	ExitCode   int
	Plugin     string
	Context    string
}

// SuggestionResponse represents a response from an AI suggestion
type SuggestionResponse struct {
	SuggestedCommand string `json:"suggested_command"`
	Explanation      string `json:"explanation"`
	Confidence       float64 `json:"confidence,omitempty"`
}
