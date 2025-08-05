package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ayushsharma-1/LogAid/internal/engine"
	"github.com/ayushsharma-1/LogAid/internal/logger"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [command]",
	Short: "Execute a command with LogAid monitoring",
	Long: `Execute a command with LogAid monitoring. LogAid will intercept the command output
and provide AI-powered suggestions if errors are detected.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand(args)
	},
}

func executeCommand(args []string) {
	// Join arguments back into a single command string for parsing
	cmdStr := strings.Join(args, " ")
	logger.Info(fmt.Sprintf("Executing command: %s", cmdStr))

	// Split the command string into parts for proper execution
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		logger.Error("No command provided")
		os.Exit(1)
	}

	// Create command
	var cmd *exec.Cmd
	if len(parts) > 1 {
		cmd = exec.Command(parts[0], parts[1:]...)
	} else {
		cmd = exec.Command(parts[0])
	}

	// Set up environment
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	// Execute with monitoring
	if err := engine.ExecuteWithMonitoring(cmd); err != nil {
		logger.Error(fmt.Sprintf("Command execution failed: %v", err))
		os.Exit(1)
	}
}
