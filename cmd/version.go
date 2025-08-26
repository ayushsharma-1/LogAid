package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display LogAid version information",
	Long:  `Display the current version of LogAid and build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 LogAid Flow-State Agent")
		fmt.Println("Version: 2.0.1")
		fmt.Println("Build: Development")
		fmt.Println("Author: Ayush Sharma")
		fmt.Println("Repository: https://github.com/ayushsharma-1/LogAid")
		fmt.Println()
		fmt.Println("An intelligent CLI agent that detects errors and provides")
		fmt.Println("AI-generated solutions to maintain your development flow.")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
