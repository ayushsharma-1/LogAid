package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "logaid",
	Short: "An intelligent CLI agent that detects errors and provides AI-generated solutions",
	Long: `LogAid is a Flow-State Agent that wraps your terminal and automatically 
detects command errors, providing real-time AI-generated corrections to maintain 
your development flow and minimize context switching.

🚀 Features:
- Pseudo-terminal wrapper for transparent operation
- AI-powered error analysis and suggestions  
- Non-intrusive, calm UX design
- Privacy-first with local fallback options
- Cross-distribution packaging
- Plugin system for specialized tools

🔧 Quick Setup:
1. Set your API key:
   export GEMINI_API_KEY=your-api-key

2. Test the configuration:
   logaid test

3. Start the agent:
   logaid run

📚 Commands:
- run       Start the terminal wrapper agent
- test      Test AI integration and configuration  
- config    Display current configuration
- plugins   List and manage plugins
- version   Show version information

🌐 Learn more: https://github.com/ayushsharma-1/LogAid`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.logaid.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".logaid" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".logaid")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
