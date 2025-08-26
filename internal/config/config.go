package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config represents the LogAid configuration
type Config struct {
	// AI Configuration
	AIProvider  string  `json:"ai_provider"`
	APIKey      string  `json:"api_key"`
	AIModel     string  `json:"ai_model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`

	// Logging Configuration
	LogLevel string `json:"log_level"`
	LogPath  string `json:"log_path"`

	// Plugin Configuration
	PluginDir      string   `json:"plugin_dir"`
	EnabledPlugins []string `json:"enabled_plugins"`

	// Terminal Configuration
	Shell         string `json:"shell"`
	PromptTimeout int    `json:"prompt_timeout"`

	// Feature Flags
	EnableLocalFallback bool `json:"enable_local_fallback"`
	EnableLogging       bool `json:"enable_logging"`
	EnableColors        bool `json:"enable_colors"`

	// Gemini API Key (for backward compatibility)
	GeminiAPIKey string `json:"gemini_api_key"`
	OpenAIAPIKey string `json:"openai_api_key"`
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file (ignore errors if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		// Default values
		AIProvider:          getEnvWithDefault("LOGAID_AI_PROVIDER", "gemini"),
		AIModel:             getEnvWithDefault("LOGAID_AI_MODEL", "gemini-1.5-flash"),
		MaxTokens:           getEnvInt("LOGAID_MAX_TOKENS", 1000),
		Temperature:         getEnvFloat("LOGAID_TEMPERATURE", 0.3),
		LogLevel:            getEnvWithDefault("LOGAID_LOG_LEVEL", "info"),
		LogPath:             expandPath(getEnvWithDefault("LOGAID_LOG_PATH", "~/.logaid/logs/history.json")),
		PluginDir:           expandPath(getEnvWithDefault("LOGAID_PLUGIN_DIR", "~/.logaid/plugins")),
		Shell:               getEnvWithDefault("LOGAID_SHELL", getDefaultShell()),
		PromptTimeout:       getEnvInt("LOGAID_PROMPT_TIMEOUT", 30),
		EnableLocalFallback: getEnvBool("LOGAID_ENABLE_LOCAL_FALLBACK", false),
		EnableLogging:       getEnvBool("LOGAID_ENABLE_LOGGING", true),
		EnableColors:        getEnvBool("LOGAID_ENABLE_COLORS", true),
	}

	// Load enabled plugins
	pluginsStr := getEnvWithDefault("LOGAID_ENABLED_PLUGINS", "apt,git,npm,docker,kubernetes,generic")
	config.EnabledPlugins = strings.Split(pluginsStr, ",")
	for i := range config.EnabledPlugins {
		config.EnabledPlugins[i] = strings.TrimSpace(config.EnabledPlugins[i])
	}

	// Load API keys
	config.APIKey = getEnvWithDefault("LOGAID_API_KEY", "")
	config.GeminiAPIKey = getEnvWithDefault("GEMINI_API_KEY", "")
	config.OpenAIAPIKey = getEnvWithDefault("OPENAI_API_KEY", "")

	// Use provider-specific API key if main API key is not set
	if config.APIKey == "" {
		switch config.AIProvider {
		case "gemini":
			config.APIKey = config.GeminiAPIKey
		case "openai":
			config.APIKey = config.OpenAIAPIKey
		}
	}

	// Validate required configuration
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required. Set LOGAID_API_KEY, GEMINI_API_KEY, or OPENAI_API_KEY")
	}

	return config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDefaultShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}
	return "/bin/bash"
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
