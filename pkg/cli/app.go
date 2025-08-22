package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ayushsharma-1/LogAid/pkg/ai"
	"github.com/ayushsharma-1/LogAid/pkg/config"
	"github.com/ayushsharma-1/LogAid/pkg/logger"
	"github.com/ayushsharma-1/LogAid/pkg/plugin"
	"github.com/ayushsharma-1/LogAid/pkg/pty"
)

// App represents the main CLI application
type App struct {
	config *config.Config
	logger *logger.Logger
}

// NewApp creates a new CLI application
func NewApp(cfg *config.Config, log *logger.Logger) *App {
	return &App{
		config: cfg,
		logger: log,
	}
}

// Run runs the CLI application
func (a *App) Run(args []string) error {
	if len(args) < 2 {
		return a.runInteractiveMode()
	}

	command := args[1]
	switch command {
	case "start":
		return a.runInteractiveMode()
	case "config":
		return a.showConfig()
	case "logs":
		return a.showLogs()
	case "help":
		return a.showHelp()
	case "version":
		return a.showVersion()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// runInteractiveMode starts the interactive monitoring mode
func (a *App) runInteractiveMode() error {
	// Initialize AI client
	aiClient, err := ai.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize AI client: %w", err)
	}

	// Initialize plugin manager
	pluginManager := plugin.NewManager()
	pluginManager.LoadBuiltinPlugins()

	// Create PTY wrapper
	wrapper := pty.NewWrapper(a.config, a.logger, pluginManager, aiClient)

	// Setup signal handling
	sigChan := wrapper.SetupSignalHandling()

	// Start monitoring in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- wrapper.Start(ctx)
	}()

	// Wait for completion or signal
	select {
	case err := <-errChan:
		if err != nil {
			a.logger.Error("PTY wrapper error: %v", err)
			return err
		}
	case sig := <-sigChan:
		a.logger.Info("Received signal: %v", sig)
		cancel()
		
		// Wait for graceful shutdown
		if err := wrapper.Shutdown(); err != nil {
			a.logger.Error("Shutdown error: %v", err)
			return err
		}
	}

	return nil
}

// showConfig displays the current configuration
func (a *App) showConfig() error {
	fmt.Println("LogAid Configuration:")
	fmt.Printf("  AI Provider: %s\n", a.config.AIProvider)
	fmt.Printf("  Log Level: %s\n", a.config.LogLevel)
	fmt.Printf("  Log Path: %s\n", a.config.LogPath)
	fmt.Printf("  Shell: %s\n", a.config.Shell)
	fmt.Printf("  Plugin Directory: %s\n", a.config.PluginDir)
	fmt.Printf("  Enabled Plugins: %v\n", a.config.EnabledPlugins)
	fmt.Printf("  Colors Enabled: %t\n", a.config.EnableColors)
	fmt.Printf("  Debug Mode: %t\n", a.config.DebugMode)
	return nil
}

// showLogs displays recent logs
func (a *App) showLogs() error {
	logFile := a.config.HistoryFile
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		fmt.Printf("No log file found at: %s\n", logFile)
		return nil
	}

	// Read and display recent log entries
	data, err := os.ReadFile(logFile)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	fmt.Printf("Recent LogAid Activity (from %s):\n", logFile)
	fmt.Println(string(data))
	return nil
}

// showHelp displays help information
func (a *App) showHelp() error {
	help := `LogAid - Your CLI Guardian

USAGE:
    logaid [COMMAND]

COMMANDS:
    start       Start interactive monitoring mode (default)
    config      Show current configuration
    logs        Show recent activity logs
    help        Show this help message
    version     Show version information

EXAMPLES:
    logaid              # Start monitoring
    logaid start        # Start monitoring (explicit)
    logaid config       # Show configuration
    logaid logs         # Show logs

CONFIGURATION:
    LogAid can be configured via environment variables or a .env file.
    See .env.example for all available options.

    Key variables:
    - GEMINI_API_KEY: Google Gemini API key
    - OPENAI_API_KEY: OpenAI API key
    - AI_PROVIDER: Preferred AI provider (gemini/openai)
    - LOG_LEVEL: Logging level (debug/info/warn/error)
    - ENABLE_COLORS: Enable colored output (true/false)

MORE INFO:
    Visit: https://github.com/ayushsharma-1/LogAid
`
	fmt.Println(help)
	return nil
}

// showVersion displays version information
func (a *App) showVersion() error {
	version := "LogAid v1.0.0"
	buildInfo := "Built with Go"
	
	fmt.Printf("%s\n", version)
	fmt.Printf("%s\n", buildInfo)
	fmt.Printf("Repository: https://github.com/ayushsharma-1/LogAid\n")
	
	// Show active providers
	if aiClient, err := ai.NewClient(); err == nil {
		providers := aiClient.GetActiveProviders()
		if len(providers) > 0 {
			fmt.Printf("Active AI Providers: %v\n", providers)
		} else {
			fmt.Println("No AI providers configured")
		}
	}
	
	return nil
}

// ValidateEnvironment validates the environment and configuration
func (a *App) ValidateEnvironment() error {
	// Check if we're in a terminal
	if !isTerminal() {
		return fmt.Errorf("LogAid must be run in a terminal")
	}

	// Check log directory permissions
	logDir := filepath.Dir(a.config.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("cannot create log directory: %w", err)
	}

	// Test log file writing
	testFile := filepath.Join(logDir, ".logaid_test")
	if file, err := os.Create(testFile); err != nil {
		return fmt.Errorf("cannot write to log directory: %w", err)
	} else {
		file.Close()
		os.Remove(testFile)
	}

	return nil
}

// isTerminal checks if the current process is running in a terminal
func isTerminal() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
