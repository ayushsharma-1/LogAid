package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	// AI Configuration
	AIProvider       string  `mapstructure:"AI_PROVIDER"`
	GeminiAPIKey     string  `mapstructure:"GEMINI_API_KEY"`
	GeminiModel      string  `mapstructure:"GEMINI_MODEL"`
	OpenAIAPIKey     string  `mapstructure:"OPENAI_API_KEY"`
	OpenAIModel      string  `mapstructure:"OPENAI_MODEL"`
	AIRequestTimeout int     `mapstructure:"AI_REQUEST_TIMEOUT"`
	MaxAIRetries     int     `mapstructure:"MAX_AI_RETRIES"`
	AITemperature    float64 `mapstructure:"AI_TEMPERATURE"`
	AIMaxTokens      int     `mapstructure:"AI_MAX_TOKENS"`

	// Logging Configuration
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	LogFile         string `mapstructure:"LOG_FILE"`
	EnableDebugLogs bool   `mapstructure:"ENABLE_DEBUG_LOGS"`
	LogRotation     bool   `mapstructure:"LOG_ROTATION"`
	MaxLogSize      string `mapstructure:"MAX_LOG_SIZE"`
	MaxLogFiles     int    `mapstructure:"MAX_LOG_FILES"`

	// Plugin Configuration
	PluginsDir             string `mapstructure:"PLUGINS_DIR"`
	EnablePlugins          string `mapstructure:"ENABLE_PLUGINS"`
	PluginTimeout          int    `mapstructure:"PLUGIN_TIMEOUT"`
	APTSearchSuggestions   bool   `mapstructure:"APT_SEARCH_SUGGESTIONS"`
	APTEnableBackports     bool   `mapstructure:"APT_ENABLE_BACKPORTS"`
	GitAutoCorrect         bool   `mapstructure:"GIT_AUTO_CORRECT"`
	GitSuggestAliases      bool   `mapstructure:"GIT_SUGGEST_ALIASES"`
	DockerHubSearch        bool   `mapstructure:"DOCKER_HUB_SEARCH"`
	DockerSuggestTags      bool   `mapstructure:"DOCKER_SUGGEST_TAGS"`
	NPMSuggestAlternatives bool   `mapstructure:"NPM_SUGGEST_ALTERNATIVES"`
	PipSuggestVersions     bool   `mapstructure:"PIP_SUGGEST_VERSIONS"`

	// UI Configuration
	EnableColors        bool   `mapstructure:"ENABLE_COLORS"`
	EnableASCIILogo     bool   `mapstructure:"ENABLE_ASCII_LOGO"`
	AutoConfirm         bool   `mapstructure:"AUTO_CONFIRM"`
	SuggestionTimeout   int    `mapstructure:"SUGGESTION_TIMEOUT"`
	MaxSuggestions      int    `mapstructure:"MAX_SUGGESTIONS"`
	ShowConfidenceScore bool   `mapstructure:"SHOW_CONFIDENCE_SCORE"`
	EnableSoundAlerts   bool   `mapstructure:"ENABLE_SOUND_ALERTS"`
	ColorError          string `mapstructure:"COLOR_ERROR"`
	ColorSuggestion     string `mapstructure:"COLOR_SUGGESTION"`
	ColorSuccess        string `mapstructure:"COLOR_SUCCESS"`
	ColorWarning        string `mapstructure:"COLOR_WARNING"`

	// History & Caching
	HistoryFile         string `mapstructure:"HISTORY_FILE"`
	MaxHistoryEntries   int    `mapstructure:"MAX_HISTORY_ENTRIES"`
	EnableHistorySearch bool   `mapstructure:"ENABLE_HISTORY_SEARCH"`
	CacheSuggestions    bool   `mapstructure:"CACHE_SUGGESTIONS"`
	CacheDuration       int    `mapstructure:"CACHE_DURATION"`
	CacheDir            string `mapstructure:"CACHE_DIR"`

	// Security & Safety
	DangerousCommandsCheck  bool   `mapstructure:"DANGEROUS_COMMANDS_CHECK"`
	RequireSudoConfirmation bool   `mapstructure:"REQUIRE_SUDO_CONFIRMATION"`
	SandboxMode             bool   `mapstructure:"SANDBOX_MODE"`
	WhitelistCommands       bool   `mapstructure:"WHITELIST_COMMANDS"`
	BlacklistCommands       string `mapstructure:"BLACKLIST_COMMANDS"`

	// Performance Settings
	PTYBufferSize     int    `mapstructure:"PTY_BUFFER_SIZE"`
	ConcurrentPlugins bool   `mapstructure:"CONCURRENT_PLUGINS"`
	EnableAsyncAI     bool   `mapstructure:"ENABLE_ASYNC_AI"`
	MemoryLimit       string `mapstructure:"MEMORY_LIMIT"`

	// Development & Testing
	DebugMode              bool   `mapstructure:"DEBUG_MODE"`
	TestMode               bool   `mapstructure:"TEST_MODE"`
	MockAIResponses        bool   `mapstructure:"MOCK_AI_RESPONSES"`
	EnableTelemetry        bool   `mapstructure:"ENABLE_TELEMETRY"`
	TelemetryEndpoint      string `mapstructure:"TELEMETRY_ENDPOINT"`
	TestDataDir            string `mapstructure:"TEST_DATA_DIR"`
	IntegrationTestTimeout int    `mapstructure:"INTEGRATION_TEST_TIMEOUT"`
	E2ETestContainers      bool   `mapstructure:"E2E_TEST_CONTAINERS"`
}

var AppConfig *Config

// Init initializes the configuration
func Init() error {
	// Set default values
	setDefaults()

	// Create config directory if it doesn't exist
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load .env file if it exists
	envFile := filepath.Join(configDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Configure viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand home directory in paths
	if err := expandPaths(); err != nil {
		return fmt.Errorf("failed to expand paths: %w", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("AI_PROVIDER", "gemini")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE", "~/.logaid/logs/logaid.log")
	viper.SetDefault("PLUGINS_DIR", "~/.logaid/plugins")
	viper.SetDefault("ENABLE_PLUGINS", "apt,npm,git,docker,pip,systemctl")
	viper.SetDefault("ENABLE_COLORS", true)
	viper.SetDefault("AUTO_CONFIRM", false)
	viper.SetDefault("SUGGESTION_TIMEOUT", 30)
	viper.SetDefault("HISTORY_FILE", "~/.logaid/logs/history.json")
	viper.SetDefault("MAX_HISTORY_ENTRIES", 1000)
	viper.SetDefault("PTY_BUFFER_SIZE", 4096)
	viper.SetDefault("AI_REQUEST_TIMEOUT", 10)
	viper.SetDefault("ENABLE_TELEMETRY", false)
}

func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".logaid"
	}
	return filepath.Join(homeDir, ".logaid")
}

func expandPaths() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Expand LogFile path
	if filepath.HasPrefix(AppConfig.LogFile, "~/") {
		AppConfig.LogFile = filepath.Join(homeDir, AppConfig.LogFile[2:])
	}

	// Expand PluginsDir path
	if filepath.HasPrefix(AppConfig.PluginsDir, "~/") {
		AppConfig.PluginsDir = filepath.Join(homeDir, AppConfig.PluginsDir[2:])
	}

	// Expand HistoryFile path
	if filepath.HasPrefix(AppConfig.HistoryFile, "~/") {
		AppConfig.HistoryFile = filepath.Join(homeDir, AppConfig.HistoryFile[2:])
	}

	return nil
}
