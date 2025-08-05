package main

import (
	"fmt"
	"os"

	"github.com/ayush-1/logaid/cmd"
	"github.com/ayush-1/logaid/internal/config"
	"github.com/ayush-1/logaid/internal/logger"
)

func main() {
	// Initialize configuration
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Execute root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
