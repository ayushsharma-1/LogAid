package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ayushsharma-1/LogAid/pkg/cli"
	"github.com/ayushsharma-1/LogAid/pkg/config"
	"github.com/ayushsharma-1/LogAid/pkg/logger"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := logger.New(cfg.LogLevel, cfg.LogPath)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Show ASCII logo
	showLogo()

	// Create and run CLI application
	app := cli.NewApp(cfg, logger)
	if err := app.Run(os.Args); err != nil {
		logger.Error("Application error: %v", err)
		os.Exit(1)
	}
}

func showLogo() {
	fmt.Println(`
   _                _    _     _ 
  | |    ___   __ _| | _(_) __| |
  | |   / _ \ / _` + "`" + ` | |/ / |/ _` + "`" + ` |
  | |__| (_) | (_| | | <| | (_| |
  |_____\___/ \__,_|_|\_\_|\__,_|
             |___/               
      LogAid: Your CLI Guardian    
`)
}
