package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of LogAid",
	Long:  `Print the version number of LogAid`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("LogAid v1.0.0")
		fmt.Println("AI-powered Linux CLI assistant")
		fmt.Println("Built with ❤️  in Go")
	},
}
