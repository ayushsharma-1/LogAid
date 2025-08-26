package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ayushsharma-1/LogAid/internal/ai"
	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/ayushsharma-1/LogAid/internal/plugins"
	"github.com/spf13/cobra"
)

// testNpmCmd represents the test-npm command
var testNpmCmd = &cobra.Command{
	Use:   "test-npm",
	Short: "Test LogAid with NPM error scenarios",
	Long:  `Test LogAid's ability to analyze and provide suggestions for NPM-related errors.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧪 Testing LogAid with NPM Error Scenarios...")
		fmt.Println()

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ Failed to load configuration: %v\n", err)
			return
		}

		// Create AI client
		aiClient, err := ai.NewClient(cfg)
		if err != nil {
			fmt.Printf("❌ Failed to create AI client: %v\n", err)
			return
		}

		// Create plugin manager
		pluginManager := plugins.NewManager(cfg.EnabledPlugins)

		// Test scenarios
		testScenarios := []struct {
			name     string
			command  string
			stderr   string
			exitCode int
		}{
			{
				name:     "Unknown NPM Command",
				command:  "npm d",
				stderr:   "Unknown command: \"d\"\n\nTo see a list of supported npm commands, run:\n  npm help",
				exitCode: 1,
			},
			{
				name:     "NPM Install Error",
				command:  "npm install nonexistent-package-xyz123",
				stderr:   "npm ERR! code E404\nnpm ERR! 404 Not Found - GET https://registry.npmjs.org/nonexistent-package-xyz123 - Not found",
				exitCode: 1,
			},
			{
				name:     "NPM Permission Error",
				command:  "npm install -g some-package",
				stderr:   "npm ERR! code EACCES\nnpm ERR! syscall mkdir\nnpm ERR! path /usr/local/lib/node_modules/some-package\nnpm ERR! errno -13\nnpm ERR! Error: EACCES: permission denied, mkdir '/usr/local/lib/node_modules/some-package'",
				exitCode: 1,
			},
		}

		for i, scenario := range testScenarios {
			fmt.Printf("📦 Test %d: %s\n", i+1, scenario.name)
			fmt.Printf("   Command: %s\n", scenario.command)
			fmt.Printf("   Error: %s\n", scenario.stderr)
			fmt.Println()

			// Test with AI client
			fmt.Println("🤖 AI Analysis:")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			suggestion, err := aiClient.AnalyzeError(ctx, scenario.command, scenario.stderr, scenario.exitCode)
			cancel()

			if err != nil {
				fmt.Printf("   ❌ AI analysis failed: %v\n", err)
			} else {
				fmt.Printf("   💡 Explanation: %s\n", suggestion.Explanation)
				if suggestion.Command != "" {
					fmt.Printf("   🔧 Suggested fix: %s\n", suggestion.Command)
				}
				fmt.Printf("   📊 Confidence: %.2f\n", suggestion.Confidence)
			}

			// Test with plugin system
			fmt.Println("🔌 Plugin Analysis:")
			pluginSuggestions, err := pluginManager.GetSuggestions(context.Background(), scenario.command, scenario.stderr, scenario.exitCode)
			if err != nil {
				fmt.Printf("   ❌ Plugin analysis failed: %v\n", err)
			} else if len(pluginSuggestions) > 0 {
				for _, pluginSuggestion := range pluginSuggestions {
					fmt.Printf("   🔧 %s: %s\n", pluginSuggestion.Type, pluginSuggestion.Command)
					fmt.Printf("      📝 %s\n", pluginSuggestion.Description)
				}
			} else {
				fmt.Println("   ℹ️  No specific plugin suggestions available")
			}

			fmt.Println()
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()
		}

		fmt.Println("✅ NPM error testing completed!")
		fmt.Println()
		fmt.Println("🚀 To see LogAid in action with real errors:")
		fmt.Println("   1. Run: logaid run")
		fmt.Println("   2. Try some commands that fail")
		fmt.Println("   3. Watch LogAid provide instant suggestions!")
	},
}

func init() {
	rootCmd.AddCommand(testNpmCmd)
}
