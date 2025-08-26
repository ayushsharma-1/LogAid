package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ayushsharma-1/LogAid/internal/ai"
	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test LogAid AI integration",
	Long: `Test the LogAid AI integration by analyzing a sample error.
This command helps verify that your API keys and configuration are working correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧪 Testing LogAid AI Integration...")
		fmt.Println()

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ Failed to load configuration: %v\n", err)
			return
		}

		fmt.Printf("✅ Configuration loaded successfully\n")
		fmt.Printf("   Provider: %s\n", cfg.AIProvider)
		fmt.Printf("   Model: %s\n", cfg.AIModel)
		fmt.Println()

		// Create AI client
		aiClient, err := ai.NewClient(cfg)
		if err != nil {
			fmt.Printf("❌ Failed to create AI client: %v\n", err)
			return
		}

		fmt.Printf("✅ AI client created successfully\n")
		fmt.Println()

		// Test with a sample error
		fmt.Println("🔍 Analyzing sample error...")

		sampleCommand := "ls /nonexistent/directory"
		sampleStderr := "ls: cannot access '/nonexistent/directory': No such file or directory"
		sampleExitCode := 2

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		suggestion, err := aiClient.AnalyzeError(ctx, sampleCommand, sampleStderr, sampleExitCode)
		if err != nil {
			fmt.Printf("❌ AI analysis failed: %v\n", err)
			return
		}

		fmt.Printf("✅ AI analysis completed successfully\n")
		fmt.Println()

		// Display the suggestion
		fmt.Println("💡 Sample AI Suggestion:")
		fmt.Println("========================")
		fmt.Printf("Command: %s\n", sampleCommand)
		fmt.Printf("Error: %s\n", sampleStderr)
		fmt.Println()
		fmt.Printf("🤖 AI Suggestion:\n")
		fmt.Printf("   Explanation: %s\n", suggestion.Explanation)
		if suggestion.Command != "" {
			fmt.Printf("   Fix Command: %s\n", suggestion.Command)
		}
		fmt.Printf("   Confidence: %.2f\n", suggestion.Confidence)
		fmt.Println()
		fmt.Printf("✅ LogAid is ready to use! Run 'logaid run' to start the agent.\n")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
