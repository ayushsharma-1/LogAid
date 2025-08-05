package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/ayushsharma-1/LogAid/internal/logger"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage LogAid configuration",
	Long:  `Manage LogAid configuration settings`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize LogAid configuration",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
}

func showConfig() {
	if config.AppConfig == nil {
		logger.Error("Configuration not initialized")
		return
	}

	fmt.Println("LogAid Configuration:")
	fmt.Printf("AI Provider: %s\n", config.AppConfig.AIProvider)
	fmt.Printf("Log Level: %s\n", config.AppConfig.LogLevel)
	fmt.Printf("Log File: %s\n", config.AppConfig.LogFile)
	fmt.Printf("Plugins Directory: %s\n", config.AppConfig.PluginsDir)
	fmt.Printf("Enabled Plugins: %s\n", config.AppConfig.EnablePlugins)
	fmt.Printf("Enable Colors: %t\n", config.AppConfig.EnableColors)
	fmt.Printf("Auto Confirm: %t\n", config.AppConfig.AutoConfirm)
	fmt.Printf("History File: %s\n", config.AppConfig.HistoryFile)
}

func initConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Failed to get home directory")
		return
	}

	configDir := filepath.Join(homeDir, ".logaid")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create config directory: %v", err))
		return
	}

	// Create logs directory
	logsDir := filepath.Join(configDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create logs directory: %v", err))
		return
	}

	// Create plugins directory
	pluginsDir := filepath.Join(configDir, "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create plugins directory: %v", err))
		return
	}

	// Copy .env.example to .env if it doesn't exist
	envFile := filepath.Join(configDir, ".env")
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		exampleFile := ".env.example"
		if content, err := os.ReadFile(exampleFile); err == nil {
			if err := os.WriteFile(envFile, content, 0644); err != nil {
				logger.Error(fmt.Sprintf("Failed to create .env file: %v", err))
				return
			}
		}
	}

	logger.Success("LogAid configuration initialized successfully!")
	logger.Info(fmt.Sprintf("Configuration directory: %s", configDir))
	logger.Info(fmt.Sprintf("Edit %s to configure your API keys", envFile))
}
