package cmd

import (
	"fmt"
	"log"

	"github.com/ayushsharma-1/LogAid/internal/agent"
	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the LogAid terminal wrapper",
	Long: `Start the LogAid agent which wraps your terminal and provides
AI-powered error detection and suggestions. This command will:

- Create a pseudo-terminal wrapper around your shell
- Monitor command execution for errors
- Provide real-time AI-generated suggestions for fixes
- Maintain your flow state with non-intrusive UX`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Override shell if specified
		if shell != "" {
			cfg.Shell = shell
		}

		// Create and start the agent
		agent, err := agent.New(cfg)
		if err != nil {
			log.Fatalf("Failed to create agent: %v", err)
		}

		if minimal {
			fmt.Println("🔧 Using minimal shell environment to avoid configuration conflicts")
		}
		fmt.Println("🚀 LogAid Flow-State Agent starting...")
		fmt.Println("   Press Ctrl+C to exit")
		fmt.Println()

		if err := agent.Start(); err != nil {
			log.Fatalf("Agent failed: %v", err)
		}
	},
}

var shell string
var minimal bool

func init() {
	rootCmd.AddCommand(runCmd)

	// Add flags for the run command
	runCmd.Flags().StringVarP(&shell, "shell", "s", "", "Shell to wrap (default: $SHELL or /bin/bash)")
	runCmd.Flags().BoolVarP(&minimal, "minimal", "m", false, "Use minimal shell environment (recommended for avoiding config conflicts)")
}
