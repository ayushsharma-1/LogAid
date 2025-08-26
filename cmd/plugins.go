package cmd

import (
	"fmt"

	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/ayushsharma-1/LogAid/internal/plugins"
	"github.com/spf13/cobra"
)

// pluginsCmd represents the plugins command
var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "List and manage LogAid plugins",
	Long: `List all available LogAid plugins and their status.
Plugins extend LogAid's error analysis capabilities for specific tools and frameworks.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			return
		}

		// Create plugin manager
		manager := plugins.NewManager(cfg.EnabledPlugins)
		pluginInfo := manager.GetPluginInfo()

		fmt.Println("🔌 LogAid Plugins")
		fmt.Println("================")
		fmt.Println()

		// Display enabled plugins
		fmt.Println("✅ Enabled Plugins:")
		hasEnabled := false
		for name, info := range pluginInfo {
			if info.Enabled {
				fmt.Printf("   • %s - %s\n", name, info.Description)
				hasEnabled = true
			}
		}
		if !hasEnabled {
			fmt.Println("   None")
		}
		fmt.Println()

		// Display disabled plugins
		fmt.Println("❌ Disabled Plugins:")
		hasDisabled := false
		for name, info := range pluginInfo {
			if !info.Enabled {
				fmt.Printf("   • %s - %s\n", name, info.Description)
				hasDisabled = true
			}
		}
		if !hasDisabled {
			fmt.Println("   None")
		}
		fmt.Println()

		fmt.Printf("📁 Plugin Directory: %s\n", cfg.PluginDir)
		fmt.Println()
		fmt.Println("To enable/disable plugins, edit the LOGAID_ENABLED_PLUGINS")
		fmt.Println("environment variable or your .env file.")
	},
}

func init() {
	rootCmd.AddCommand(pluginsCmd)
}
