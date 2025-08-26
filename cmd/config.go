package cmd

import (
	"fmt"

	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display current LogAid configuration",
	Long: `Display the current LogAid configuration including:
- AI provider and model settings
- API key status (masked for security)
- Feature flags and logging settings
- Plugin configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			return
		}

		fmt.Println("📋 LogAid Configuration")
		fmt.Println("========================")
		fmt.Printf("AI Provider:         %s\n", cfg.AIProvider)
		fmt.Printf("AI Model:            %s\n", cfg.AIModel)
		fmt.Printf("Max Tokens:          %d\n", cfg.MaxTokens)
		fmt.Printf("Temperature:         %.2f\n", cfg.Temperature)
		fmt.Println()

		// Mask API key for security
		apiKeyMasked := "Not set"
		if cfg.APIKey != "" {
			if len(cfg.APIKey) > 8 {
				apiKeyMasked = cfg.APIKey[:4] + "..." + cfg.APIKey[len(cfg.APIKey)-4:]
			} else {
				apiKeyMasked = "***"
			}
		}
		fmt.Printf("API Key:             %s\n", apiKeyMasked)
		fmt.Println()

		fmt.Printf("Log Level:           %s\n", cfg.LogLevel)
		fmt.Printf("Log Path:            %s\n", cfg.LogPath)
		fmt.Printf("Shell:               %s\n", cfg.Shell)
		fmt.Printf("Prompt Timeout:      %ds\n", cfg.PromptTimeout)
		fmt.Println()

		fmt.Printf("Enable Local Fallback: %t\n", cfg.EnableLocalFallback)
		fmt.Printf("Enable Logging:        %t\n", cfg.EnableLogging)
		fmt.Printf("Enable Colors:         %t\n", cfg.EnableColors)
		fmt.Println()

		fmt.Printf("Plugin Directory:    %s\n", cfg.PluginDir)
		fmt.Printf("Enabled Plugins:     %v\n", cfg.EnabledPlugins)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
