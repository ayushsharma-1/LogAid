package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ayushsharma-1/LogAid/pkg/plugin"
)

// Config holds the application configuration
type Config struct {
	// AI Configuration
	AIProvider      string
	GeminiAPIKey    string
	GeminiModel     string
	OpenAIAPIKey    string
	OpenAIModel     string
	AIRequestTimeout int
	MaxAIRetries    int
	AITemperature   float64
	AIMaxTokens     int
	
	// Logging Configuration
	LogLevel        string
	LogPath         string
	EnableDebugLogs bool
	LogRotation     bool
	MaxLogSize      string
	MaxLogFiles     int
	
	// Plugin Configuration
	PluginDir       string
	EnabledPlugins  []string
	PluginTimeout   int
	
	// Terminal Configuration
	Shell           string
	PromptTimeout   int // seconds
	
	// UI Configuration
	EnableColors    bool
	EnableASCIILogo bool
	AutoConfirm     bool
	SuggestionTimeout int
	MaxSuggestions  int
	ShowConfidenceScore bool
	EnableSoundAlerts bool
	
	// History & Caching
	HistoryFile     string
	MaxHistoryEntries int
	EnableHistorySearch bool
	CacheSuggestions bool
	CacheDuration   int
	CacheDir        string
	
	// Security & Safety
	DangerousCommandsCheck bool
	RequireSudoConfirmation bool
	SandboxMode     bool
	WhitelistCommands bool
	BlacklistCommands []string
	
	// Security Sanitization
	EnableSecuritySanitization bool
	AutoSanitizeForAI         bool
	RequireConsentForAI       bool
	RequireConsentForSensitive bool
	MaxRiskLevelForAI         string
	BlockHighRiskCommands     bool
	LogSecurityEvents         bool
	
	// Performance Settings
	PTYBufferSize   int
	ConcurrentPlugins bool
	EnableAsyncAI   bool
	MemoryLimit     string
	
	// Development & Testing
	DebugMode       bool
	TestMode        bool
	MockAIResponses bool
	EnableTelemetry bool
	TelemetryEndpoint string
}

// Default configuration values
const (
	DefaultAIProvider = "gemini"
	DefaultAIModel    = "gemini-1.5-flash"
	DefaultLogLevel   = "info"
	DefaultShell      = "/bin/bash"
	DefaultTimeout    = 30
)

// Load loads configuration from environment variables and defaults
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	cfg := &Config{
		// AI Configuration
		AIProvider:       getEnvOrDefault("AI_PROVIDER", DefaultAIProvider),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		GeminiModel:      getEnvOrDefault("GEMINI_MODEL", "gemini-2.0-flash-exp"),
		OpenAIAPIKey:     os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:      getEnvOrDefault("OPENAI_MODEL", "gpt-4o"),
		AIRequestTimeout: getEnvIntOrDefault("AI_REQUEST_TIMEOUT", 15),
		MaxAIRetries:     getEnvIntOrDefault("MAX_AI_RETRIES", 3),
		AITemperature:    getEnvFloatOrDefault("AI_TEMPERATURE", 0.1),
		AIMaxTokens:      getEnvIntOrDefault("AI_MAX_TOKENS", 500),
		
		// Logging Configuration
		LogLevel:         getEnvOrDefault("LOG_LEVEL", DefaultLogLevel),
		LogPath:          getEnvOrDefault("LOG_FILE", filepath.Join(homeDir, ".logaid", "logs", "logaid.log")),
		EnableDebugLogs:  getEnvBoolOrDefault("ENABLE_DEBUG_LOGS", false),
		LogRotation:      getEnvBoolOrDefault("LOG_ROTATION", true),
		MaxLogSize:       getEnvOrDefault("MAX_LOG_SIZE", "10MB"),
		MaxLogFiles:      getEnvIntOrDefault("MAX_LOG_FILES", 5),
		
		// Plugin Configuration
		PluginDir:        getEnvOrDefault("PLUGINS_DIR", filepath.Join(homeDir, ".logaid", "plugins")),
		EnabledPlugins:   getEnvSliceOrDefault("ENABLE_PLUGINS", []string{"apt", "npm", "git", "docker", "pip", "systemctl", "yarn", "cargo", "make", "ssh"}),
		PluginTimeout:    getEnvIntOrDefault("PLUGIN_TIMEOUT", 5),
		
		// Terminal Configuration
		Shell:            getEnvOrDefault("SHELL", DefaultShell),
		PromptTimeout:    getEnvIntOrDefault("SUGGESTION_TIMEOUT", DefaultTimeout),
		
		// UI Configuration
		EnableColors:        getEnvBoolOrDefault("ENABLE_COLORS", true),
		EnableASCIILogo:     getEnvBoolOrDefault("ENABLE_ASCII_LOGO", true),
		AutoConfirm:         getEnvBoolOrDefault("AUTO_CONFIRM", false),
		SuggestionTimeout:   getEnvIntOrDefault("SUGGESTION_TIMEOUT", 30),
		MaxSuggestions:      getEnvIntOrDefault("MAX_SUGGESTIONS", 5),
		ShowConfidenceScore: getEnvBoolOrDefault("SHOW_CONFIDENCE_SCORE", true),
		EnableSoundAlerts:   getEnvBoolOrDefault("ENABLE_SOUND_ALERTS", false),
		
		// History & Caching
		HistoryFile:         getEnvOrDefault("HISTORY_FILE", filepath.Join(homeDir, ".logaid", "history.json")),
		MaxHistoryEntries:   getEnvIntOrDefault("MAX_HISTORY_ENTRIES", 10000),
		EnableHistorySearch: getEnvBoolOrDefault("ENABLE_HISTORY_SEARCH", true),
		CacheSuggestions:    getEnvBoolOrDefault("CACHE_SUGGESTIONS", true),
		CacheDuration:       getEnvIntOrDefault("CACHE_DURATION", 3600),
		CacheDir:            getEnvOrDefault("CACHE_DIR", filepath.Join(homeDir, ".logaid", "cache")),
		
		// Security & Safety
		DangerousCommandsCheck:  getEnvBoolOrDefault("DANGEROUS_COMMANDS_CHECK", true),
		RequireSudoConfirmation: getEnvBoolOrDefault("REQUIRE_SUDO_CONFIRMATION", true),
		SandboxMode:             getEnvBoolOrDefault("SANDBOX_MODE", false),
		WhitelistCommands:       getEnvBoolOrDefault("WHITELIST_COMMANDS", false),
		BlacklistCommands:       getEnvSliceOrDefault("BLACKLIST_COMMANDS", []string{"rm -rf /", "dd if="}),
		
		// Security Sanitization
		EnableSecuritySanitization: getEnvBoolOrDefault("ENABLE_SECURITY_SANITIZATION", true),
		AutoSanitizeForAI:         getEnvBoolOrDefault("AUTO_SANITIZE_FOR_AI", true),
		RequireConsentForAI:       getEnvBoolOrDefault("REQUIRE_CONSENT_FOR_AI", true),
		RequireConsentForSensitive: getEnvBoolOrDefault("REQUIRE_CONSENT_FOR_SENSITIVE", true),
		MaxRiskLevelForAI:         getEnvOrDefault("MAX_RISK_LEVEL_FOR_AI", "Medium"),
		BlockHighRiskCommands:     getEnvBoolOrDefault("BLOCK_HIGH_RISK_COMMANDS", true),
		LogSecurityEvents:         getEnvBoolOrDefault("LOG_SECURITY_EVENTS", true),
		
		// Performance Settings
		PTYBufferSize:     getEnvIntOrDefault("PTY_BUFFER_SIZE", 8192),
		ConcurrentPlugins: getEnvBoolOrDefault("CONCURRENT_PLUGINS", true),
		EnableAsyncAI:     getEnvBoolOrDefault("ENABLE_ASYNC_AI", true),
		MemoryLimit:       getEnvOrDefault("MEMORY_LIMIT", "256MB"),
		
		// Development & Testing
		DebugMode:         getEnvBoolOrDefault("DEBUG_MODE", false),
		TestMode:          getEnvBoolOrDefault("TEST_MODE", false),
		MockAIResponses:   getEnvBoolOrDefault("MOCK_AI_RESPONSES", false),
		EnableTelemetry:   getEnvBoolOrDefault("ENABLE_TELEMETRY", false),
		TelemetryEndpoint: getEnvOrDefault("TELEMETRY_ENDPOINT", "https://api.logaid.ayushsharma.site/telemetry"),
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Ensure directories exist
	dirs := []string{
		filepath.Dir(cfg.LogPath),
		cfg.PluginDir,
		filepath.Dir(cfg.HistoryFile),
		cfg.CacheDir,
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return cfg, nil
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.AIProvider == "" {
		return fmt.Errorf("AI provider cannot be empty")
	}

	// Check that at least one AI provider has an API key
	if c.GeminiAPIKey == "" && c.OpenAIAPIKey == "" {
		return fmt.Errorf("at least one AI provider API key must be configured")
	}

	if c.Shell == "" {
		return fmt.Errorf("shell cannot be empty")
	}

	if c.PromptTimeout <= 0 {
		return fmt.Errorf("prompt timeout must be positive")
	}

	if c.AIRequestTimeout <= 0 {
		return fmt.Errorf("AI request timeout must be positive")
	}

	return nil
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
