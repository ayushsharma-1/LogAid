package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ayush-1/logaid/internal/config"
	"github.com/ayush-1/logaid/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "logaid",
	Short: "AI-powered Linux CLI assistant",
	Long: `LogAid is a CLI-first AI assistant that intercepts shell commands and error logs 
in real time, identifies mistakes (typos, wrong package names, syntax errors, etc.), 
and suggests or auto-applies corrections with user confirmation.`,
	Run: func(cmd *cobra.Command, args []string) {
		showLogo()
		startInteractiveShell()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

func showLogo() {
	logoFile := "assets/logo.txt"
	if _, err := os.Stat(logoFile); err == nil {
		content, err := ioutil.ReadFile(logoFile)
		if err == nil {
			if config.AppConfig != nil && config.AppConfig.EnableColors {
				logger.InfoColor.Println(string(content))
			} else {
				fmt.Println(string(content))
			}
			return
		}
	}

	// Fallback ASCII logo
	logo := `  _                _    _     _ 
 | |    ___   __ _| | _(_) __| |
 | |   / _ \ / _` + "`" + ` | |/ / |/ _` + "`" + ` |
 | |__| (_) | (_| |   <| | (_| |
 |_____\___/ \__, |_|\_\_|\__,_|
             |___/              
       LogAid CLI Companion      
`
	if config.AppConfig != nil && config.AppConfig.EnableColors {
		logger.InfoColor.Println(logo)
	} else {
		fmt.Println(logo)
	}
}

func startInteractiveShell() {
	logger.Info("Starting LogAid interactive shell...")
	logger.Info("Type 'exit' to quit")

	// TODO: Implement interactive shell with PTY
	fmt.Println("Interactive shell not yet implemented. Use 'logaid exec <command>' for now.")
}
