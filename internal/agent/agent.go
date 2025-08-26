package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ayushsharma-1/LogAid/internal/ai"
	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/ayushsharma-1/LogAid/internal/ui"
	"github.com/creack/pty"
	"golang.org/x/term"
)

// Agent represents the LogAid terminal wrapper agent
type Agent struct {
	config   *config.Config
	aiClient ai.Client
	ui       *ui.Manager
	pty      *os.File
	shell    *exec.Cmd
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a new LogAid agent
func New(cfg *config.Config) (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create AI client
	aiClient, err := ai.NewClient(cfg)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}

	// Create UI manager
	uiManager := ui.NewManager(cfg)

	return &Agent{
		config:   cfg,
		aiClient: aiClient,
		ui:       uiManager,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start starts the LogAid agent
func (a *Agent) Start() error {
	defer a.cancel()

	// Show welcome message
	a.ui.ShowWelcome()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Start the shell in a pty
	if err := a.startShell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}
	defer a.cleanup()

	// Give shell a moment to initialize and capture any startup errors
	time.Sleep(100 * time.Millisecond)

	// Setup terminal in raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	// Start the main event loop
	errCh := make(chan error, 1)
	go a.eventLoop(errCh)

	// Wait for signal or error
	select {
	case <-sigCh:
		a.ui.ShowInfo("LogAid agent shutting down...")
		return nil
	case err := <-errCh:
		return err
	case <-a.ctx.Done():
		return nil
	}
}

// startShell starts the user's shell in a pseudo-terminal
func (a *Agent) startShell() error {
	// Create shell command with clean environment
	a.shell = exec.Command(a.config.Shell)

	// Set a minimal, clean environment to avoid shell config issues
	cleanEnv := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"USER=" + os.Getenv("USER"),
		"TERM=xterm-256color",
		"SHELL=" + a.config.Shell,
		"PWD=" + func() string {
			if pwd, err := os.Getwd(); err == nil {
				return pwd
			}
			return os.Getenv("HOME")
		}(),
	}

	// Add any LogAid-specific environment variables
	cleanEnv = append(cleanEnv, "LOGAID_ACTIVE=1")

	a.shell.Env = cleanEnv

	// Start shell with pty
	var err error
	a.pty, err = pty.Start(a.shell)
	if err != nil {
		return fmt.Errorf("failed to start shell with pty: %w", err)
	}

	return nil
}

// eventLoop handles the main I/O between user, pty, and AI
func (a *Agent) eventLoop(errCh chan<- error) {
	// Channel for PTY output
	ptyOutput := make(chan []byte, 1024)

	// Channel for user input
	userInput := make(chan []byte, 1024)

	// Start PTY output reader
	go a.readPTYOutput(ptyOutput, errCh)

	// Start user input reader
	go a.readUserInput(userInput, errCh)

	// Command tracking
	var currentCommand strings.Builder
	var commandStarted bool
	var outputBuffer strings.Builder
	var waitingForOutput bool

	for {
		select {
		case <-a.ctx.Done():
			return

		case data := <-userInput:
			// Forward user input to PTY
			if _, err := a.pty.Write(data); err != nil {
				errCh <- fmt.Errorf("failed to write to pty: %w", err)
				return
			}

			// Track command input
			for _, b := range data {
				switch b {
				case '\r', '\n': // Enter pressed
					if commandStarted && currentCommand.Len() > 0 {
						// Command submitted, start monitoring for output
						cmd := strings.TrimSpace(currentCommand.String())
						commandStarted = false
						waitingForOutput = true
						outputBuffer.Reset()

						// Start monitoring with a timeout
						go func(command string) {
							time.Sleep(500 * time.Millisecond)
							if waitingForOutput {
								waitingForOutput = false
								output := outputBuffer.String()
								if a.containsError(output) {
									errorMsg := a.extractErrorMessage(output)
									go a.analyzeAndSuggest(command, errorMsg, 1)
								}
							}
						}(cmd)

						currentCommand.Reset()
					}
				case 127, 8: // Backspace
					if currentCommand.Len() > 0 {
						cmd := currentCommand.String()
						currentCommand.Reset()
						currentCommand.WriteString(cmd[:len(cmd)-1])
					}
				default:
					if b >= 32 && b <= 126 { // Printable ASCII
						currentCommand.WriteByte(b)
						commandStarted = true
					}
				}
			}

		case data := <-ptyOutput:
			// Capture output for error analysis
			if waitingForOutput {
				outputBuffer.Write(data)
			}

			// Forward PTY output to user
			if _, err := os.Stdout.Write(data); err != nil {
				errCh <- fmt.Errorf("failed to write to stdout: %w", err)
				return
			}
		}
	}
}

// readPTYOutput reads output from the PTY
func (a *Agent) readPTYOutput(output chan<- []byte, errCh chan<- error) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			n, err := a.pty.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return // Shell exited
				}
				errCh <- fmt.Errorf("failed to read from pty: %w", err)
				return
			}
			if n > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])
				output <- data
			}
		}
	}
}

// readUserInput reads input from the user
func (a *Agent) readUserInput(input chan<- []byte, errCh chan<- error) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			n, err := os.Stdin.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				}
				errCh <- fmt.Errorf("failed to read from stdin: %w", err)
				return
			}
			if n > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])
				input <- data
			}
		}
	}
}

// analyzeAndSuggest analyzes command errors and provides AI suggestions
func (a *Agent) analyzeAndSuggest(command, stderr string, exitCode int) {
	// Only provide suggestions for failed commands or when we have error output
	if exitCode == 0 && stderr == "" {
		return
	}

	// Clear the line and move cursor to start
	fmt.Printf("\r\033[K")

	// Show compact loading indicator
	fmt.Printf("🔍 Analyzing...")

	// Check for quick fixes first (faster response)
	quickFix := a.getQuickFix(command, stderr)
	if quickFix != "" {
		// Clear the loading message
		fmt.Printf("\r\033[K")
		fmt.Printf("💡 \033[1;36mLogAid:\033[0m %s\n", quickFix)
		return
	}

	// Get AI suggestion for complex cases
	ctx, cancel := context.WithTimeout(a.ctx, time.Duration(a.config.PromptTimeout)*time.Second)
	defer cancel()

	suggestion, err := a.aiClient.AnalyzeError(ctx, command, stderr, exitCode)

	// Clear the loading message
	fmt.Printf("\r\033[K")

	if err != nil {
		fmt.Printf("💡 \033[1;36mLogAid:\033[0m Check command syntax and try again\n")
		return
	}

	// Show AI suggestion
	fmt.Printf("💡 \033[1;36mLogAid:\033[0m %s\n", suggestion.Explanation)
}

// getQuickFix provides instant suggestions for common errors
func (a *Agent) getQuickFix(command, stderr string) string {
	cmdLower := strings.ToLower(command)
	errLower := strings.ToLower(stderr)

	// ls command fixes
	if strings.HasPrefix(cmdLower, "ls ") {
		if strings.Contains(errLower, "no such file or directory") {
			if strings.Contains(command, "ltr") {
				return "Try \033[1;32mls -ltr\033[0m (list files by time)"
			}
			if strings.Contains(command, "la") && !strings.Contains(command, "-la") {
				return "Try \033[1;32mls -la\033[0m (list all files)"
			}
		}
	}

	// npm command fixes
	if strings.HasPrefix(cmdLower, "npm ") {
		if strings.Contains(errLower, "unknown command") {
			if strings.Contains(command, " d") {
				return "Try \033[1;32mnpm install\033[0m or \033[1;32mnpm run dev\033[0m"
			}
		}
	}

	// git command fixes
	if strings.HasPrefix(cmdLower, "git ") {
		if strings.Contains(errLower, "not a git command") {
			if strings.Contains(command, "stauts") {
				return "Try \033[1;32mgit status\033[0m"
			}
			if strings.Contains(command, "committ") {
				return "Try \033[1;32mgit commit\033[0m"
			}
		}
	}

	// Command not found
	if strings.Contains(errLower, "command not found") {
		parts := strings.Fields(command)
		if len(parts) > 0 {
			cmd := parts[0]
			return fmt.Sprintf("Command '\033[1;31m%s\033[0m' not found. Check if it's installed", cmd)
		}
	}

	return ""
}

// containsError checks if the output contains error indicators
func (a *Agent) containsError(output string) bool {
	errorPatterns := []string{
		"command not found",
		"No such file or directory",
		"Permission denied",
		"not a git command",
		"Unknown command",
		"invalid option",
		"Error:",
		"error:",
		"FAILED",
		"failed",
		"Cannot",
		"cannot",
		"bash:",
		"sh:",
		"zsh:",
	}

	outputLower := strings.ToLower(output)
	for _, pattern := range errorPatterns {
		if strings.Contains(outputLower, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// extractErrorMessage extracts the relevant error message from output
func (a *Agent) extractErrorMessage(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && a.containsError(line) {
			return line
		}
	}
	return output
}

// cleanup cleans up resources
func (a *Agent) cleanup() {
	if a.pty != nil {
		_ = a.pty.Close()
	}
	if a.shell != nil && a.shell.Process != nil {
		_ = a.shell.Process.Kill()
	}
}
