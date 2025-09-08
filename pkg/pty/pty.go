package pty

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ayushsharma-1/LogAid/pkg/ai"
	"github.com/ayushsharma-1/LogAid/pkg/config"
	"github.com/ayushsharma-1/LogAid/pkg/logger"
	"github.com/ayushsharma-1/LogAid/pkg/plugin"
)

// Wrapper wraps a shell with LogAid functionality
type Wrapper struct {
	config        *config.Config
	logger        *logger.Logger
	pluginManager *plugin.Manager
	aiClient      *ai.Client
	shell         string
	cmd           *exec.Cmd
	sigChan       chan os.Signal
}

// NewWrapper creates a new PTY wrapper
func NewWrapper(cfg *config.Config, log *logger.Logger, pluginManager *plugin.Manager, aiClient *ai.Client) *Wrapper {
	return &Wrapper{
		config:        cfg,
		logger:        log,
		pluginManager: pluginManager,
		aiClient:      aiClient,
		shell:         cfg.Shell,
	}
}

// Run starts the wrapped shell session it is here because it graceful startup and shutdown
func (w *Wrapper) Run(ctx context.Context) error {
	w.logger.Info("Starting LogAid shell wrapper")

	// Set up signal handling for control+c and termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start shell in background
	shellCtx, shellCancel := context.WithCancel(ctx)
	defer shellCancel()

	go func() {
		select {
		case <-sigChan:
			w.logger.Info("Received termination signal")
			shellCancel()
		case <-ctx.Done():
			shellCancel()
		}
	}()
		// delegates to run main logic
	return w.runShellLoop(shellCtx)
}

// Start starts the wrapper with context
func (w *Wrapper) Start(ctx context.Context) error {
	return w.Run(ctx)
}

// SetupSignalHandling sets up signal handling and returns the signal channel creates and return a channel for OS signal notifications
func (w *Wrapper) SetupSignalHandling() chan os.Signal {
	if w.sigChan == nil {
		w.sigChan = make(chan os.Signal, 1)
		signal.Notify(w.sigChan, syscall.SIGINT, syscall.SIGTERM)
	}
	return w.sigChan
}

// Shutdown gracefully shuts down the wrapper
func (w *Wrapper) Shutdown() error {
	if w.cmd != nil && w.cmd.Process != nil {
		w.logger.Info("Shutting down shell process")
		return w.cmd.Process.Signal(syscall.SIGTERM)
	}
	return nil
}

// runShellLoop runs the main shell interaction loop
func (w *Wrapper) runShellLoop(ctx context.Context) error {
	// For this implementation, we'll use a simpler approach:
	// Monitor commands by intercepting bash history and checking exit codes
	
	fmt.Println("LogAid is monitoring your commands...")
	fmt.Println("Type 'exit' to quit LogAid")
	fmt.Println()

	// Start the actual shell as a subprocess
	cmd := exec.CommandContext(ctx, w.shell)
	w.cmd = cmd  // Store reference for shutdown
	cmd.Env = append(os.Environ(), 
		"PS1=[LogAid] "+os.Getenv("PS1"),
		"PROMPT_COMMAND=echo \"__LOGAID_CMD__:$?:$(history 1)\" >&2; "+os.Getenv("PROMPT_COMMAND"),
	)
	// This is where we creates channles to communicate with the shell process and logAid
	// Set up pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the shell
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	// Handle I/O
	go w.handleStdout(stdout)
	go w.handleStderr(stderr)
	go w.handleStdin(stdin)

	// Wait for shell to exit
	err = cmd.Wait()
	if err != nil && ctx.Err() == nil {
		w.logger.Error("Shell exited with error: %v", err)
	}

	return nil
}

// handleStdout forwards stdout to the terminal
func (w *Wrapper) handleStdout(stdout io.Reader) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

// handleStderr monitors stderr for command info and errors
func (w *Wrapper) handleStderr(stderr io.Reader) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check for our command marker
		if strings.HasPrefix(line, "__LOGAID_CMD__:") {
			w.processCommandInfo(line)
			continue
		}
		
		// Forward other stderr output
		fmt.Fprintln(os.Stderr, line)
	}
}

// handleStdin forwards user input to the shell
func (w *Wrapper) handleStdin(stdin io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(stdin, line)
	}
}
// It helps whetere user need help or not
// processCommandInfo processes command information from the shell
func (w *Wrapper) processCommandInfo(line string) {
	// Parse: __LOGAID_CMD__:exitcode:command
	parts := strings.SplitN(line, ":", 3)
	if len(parts) != 3 {
		return
	}

	exitCodeStr := parts[1]
	historyLine := parts[2]

	// Parse exit code
	var exitCode int
	if _, err := fmt.Sscanf(exitCodeStr, "%d", &exitCode); err != nil {
		return
	}

	// Extract command from history (format: "  123  command")
	command := w.extractCommandFromHistory(historyLine)
	if command == "" {
		return
	}

	// Log the command
	w.logger.LogCommand(command, "", "", exitCode)

	// If command failed, try to suggest a fix
	if exitCode != 0 {
		w.handleFailedCommand(command, "", exitCode)
	}
}

// extractCommandFromHistory extracts the command from bash history line
func (w *Wrapper) extractCommandFromHistory(historyLine string) string {
	// History format: "  123  command args"
	parts := strings.Fields(historyLine)
	if len(parts) < 2 {
		return ""
	}

	// Skip the history number (first field)
	return strings.Join(parts[1:], " ")
}

// handleFailedCommand handles a failed command by suggesting fixes
func (w *Wrapper) handleFailedCommand(command, stderr string, exitCode int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find matching plugin
	matchingPlugin := w.pluginManager.FindMatchingPlugin(command, stderr, exitCode)
	
	var suggestion string
	var explanation string
	var pluginName string
	var processingTime time.Duration
	var sanitizationResult *plugin.SanitizationResult

	start := time.Now()
	
	if matchingPlugin != nil {
		pluginName = matchingPlugin.Name()
		w.logger.Debug("Using plugin: %s for command: %s", pluginName, command)
		
		// Try secure plugin-specific suggestion first
		if secureResult, cmd, exp, err := matchingPlugin.SuggestSecure(ctx, command, stderr, exitCode); err == nil {
			sanitizationResult = secureResult
			
			// Check if we need user consent for security reasons
			if sanitizationResult.UserConsentRequired {
				if !w.requestSecurityConsent(command, sanitizationResult) {
					w.logger.Info("User declined security consent for command: %s", command)
					return
				}
			}
			
			if cmd != "" {
				suggestion = cmd
				explanation = exp
			}
		}
	}

	// If no plugin suggestion or plugin failed, try AI with sanitization
	if suggestion == "" && w.aiClient != nil {
		w.logger.Debug("Requesting AI suggestion for command: %s", command)
		
		// If we don't have sanitization result yet, create one
		if sanitizationResult == nil {
			sanitizer := plugin.NewSecuritySanitizer()
			sanitizationResult = sanitizer.SanitizeData(command, stderr, exitCode)
		}
		
		// Check if it's safe to send to AI
		if sanitizationResult.RiskLevel >= plugin.RiskHigh {
			w.logger.Warn("Command contains sensitive data - skipping AI suggestion: %s", command)
			w.presentSecurityWarning(command, sanitizationResult)
			return
		}
		
		// Request user consent if needed
		if sanitizationResult.UserConsentRequired {
			if !w.requestAIConsent(command, sanitizationResult) {
				w.logger.Info("User declined AI sharing consent for command: %s", command)
				return
			}
		}
		
		// Use sanitized data for AI request
		sanitizedCmd := sanitizationResult.SanitizedCommand
		sanitizedOutput := sanitizationResult.SanitizedOutput
		
		req := &plugin.SuggestionRequest{
			Command:  sanitizedCmd,
			Output:   sanitizedOutput,
			ExitCode: exitCode,
			Plugin:   pluginName,
			Context:  "Linux terminal error (sanitized)",
		}

		if resp, err := w.aiClient.Suggest(ctx, req); err == nil && resp.SuggestedCommand != "" {
			suggestion = resp.SuggestedCommand
			explanation = resp.Explanation
		} else if err != nil {
			w.logger.Warn("AI suggestion failed: %v", err)
		}
	}

	processingTime = time.Since(start)

	// If we have a suggestion, present it to the user
	if suggestion != "" {
		w.presentSuggestion(command, stderr, suggestion, explanation, pluginName, processingTime)
	}
}

// presentSuggestion presents a suggestion to the user
func (w *Wrapper) presentSuggestion(originalCmd, error, suggestion, explanation, pluginName string, processingTime time.Duration) {
	// Use colors if enabled
	var (
		errorColor      = ""
		suggestionColor = ""
		promptColor     = ""
		resetColor      = ""
	)

	if w.config.EnableColors {
		errorColor = "\033[31m"      // Red
		suggestionColor = "\033[36m" // Cyan  
		promptColor = "\033[33m"     // Yellow
		resetColor = "\033[0m"       // Reset
	}

	fmt.Printf("\n%s╭─ LogAid Suggestion%s\n", suggestionColor, resetColor)
	fmt.Printf("%s│%s %s\n", suggestionColor, resetColor, suggestion)
	
	if explanation != "" {
		fmt.Printf("%s│%s Explanation: %s\n", suggestionColor, resetColor, explanation)
	}
	
	if pluginName != "" {
		fmt.Printf("%s│%s Plugin: %s\n", suggestionColor, resetColor, pluginName)
	}
	
	fmt.Printf("%s╰─%s\n", suggestionColor, resetColor)
	
	// Prompt user for confirmation
	fmt.Printf("%sExecute suggestion? [y/N]:%s ", promptColor, resetColor)
	
	// Read user input with timeout
	userChoice := w.readUserChoiceWithTimeout(w.config.PromptTimeout)
	
	var userApproved bool
	var outcome string
	
	if strings.ToLower(strings.TrimSpace(userChoice)) == "y" {
		userApproved = true
		fmt.Printf("Executing: %s%s%s\n", suggestionColor, suggestion, resetColor)
		
		// Execute the suggestion
		if err := w.executeCommand(suggestion); err != nil {
			outcome = "Failed: " + err.Error()
			fmt.Printf("%sExecution failed: %v%s\n", errorColor, err, resetColor)
		} else {
			outcome = "Success"
		}
	} else {
		userApproved = false
		outcome = "Skipped"
		fmt.Println("Suggestion skipped.")
	}

	// Log the suggestion and outcome
	if w.config.LogLevel != "" {
		aiProvider := ""
		if w.aiClient != nil {
			aiProvider = "gemini" // TODO: get from AI client
		}
		
		w.logger.LogSuggestion(
			originalCmd,
			error,
			suggestion,
			explanation,
			pluginName,
			aiProvider,
			userApproved,
			outcome,
			processingTime,
		)
	}
}

// readUserChoiceWithTimeout reads user input with a timeout
func (w *Wrapper) readUserChoiceWithTimeout(timeoutSeconds int) string {
	type result struct {
		input string
		err   error
	}

	ch := make(chan result, 1)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			ch <- result{input: scanner.Text(), err: scanner.Err()}
		}
	}()

	select {
	case res := <-ch:
		return res.input
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		fmt.Println("\nTimeout reached. Suggestion skipped.")
		return "n"
	}
}

// executeCommand executes a shell command
func (w *Wrapper) executeCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
}

// requestSecurityConsent asks user for consent when security issues are detected
func (w *Wrapper) requestSecurityConsent(command string, sanitizationResult *plugin.SanitizationResult) bool {
	// Use colors if enabled
	var (
		warningColor = ""
		promptColor  = ""
		resetColor   = ""
	)

	if w.config.EnableColors {
		warningColor = "\033[33m" // Yellow
		promptColor = "\033[36m"  // Cyan
		resetColor = "\033[0m"    // Reset
	}

	fmt.Printf("\n%s⚠️  SECURITY WARNING%s\n", warningColor, resetColor)
	fmt.Printf("%s│%s Command contains potentially sensitive information\n", warningColor, resetColor)
	fmt.Printf("%s│%s Risk Level: %s\n", warningColor, resetColor, sanitizationResult.RiskLevel.String())
	
	if len(sanitizationResult.DetectedPatterns) > 0 {
		fmt.Printf("%s│%s Detected patterns: %v\n", warningColor, resetColor, sanitizationResult.DetectedPatterns)
	}
	
	fmt.Printf("%s│%s Command: %s\n", warningColor, resetColor, command)
	fmt.Printf("%s╰─%s\n", warningColor, resetColor)
	
	fmt.Printf("%sProceed with plugin suggestion? [y/N]:%s ", promptColor, resetColor)
	
	userChoice := w.readUserChoiceWithTimeout(w.config.PromptTimeout)
	return strings.ToLower(strings.TrimSpace(userChoice)) == "y"
}

// requestAIConsent asks user for consent before sending data to AI
func (w *Wrapper) requestAIConsent(command string, sanitizationResult *plugin.SanitizationResult) bool {
	// Use colors if enabled
	var (
		infoColor   = ""
		promptColor = ""
		resetColor  = ""
	)

	if w.config.EnableColors {
		infoColor = "\033[34m"   // Blue
		promptColor = "\033[36m" // Cyan
		resetColor = "\033[0m"   // Reset
	}

	fmt.Printf("\n%s🤖 AI CONSENT REQUEST%s\n", infoColor, resetColor)
	fmt.Printf("%s│%s LogAid wants to send sanitized command data to AI for suggestions\n", infoColor, resetColor)
	fmt.Printf("%s│%s Risk Level: %s\n", infoColor, resetColor, sanitizationResult.RiskLevel.String())
	
	if len(sanitizationResult.DetectedPatterns) > 0 {
		fmt.Printf("%s│%s Sensitive data detected and will be redacted: %v\n", infoColor, resetColor, sanitizationResult.DetectedPatterns)
	}
	
	fmt.Printf("%s│%s Original: %s\n", infoColor, resetColor, command)
	fmt.Printf("%s│%s Sanitized: %s\n", infoColor, resetColor, sanitizationResult.SanitizedCommand)
	fmt.Printf("%s╰─%s\n", infoColor, resetColor)
	
	fmt.Printf("%sSend sanitized data to AI? [y/N]:%s ", promptColor, resetColor)
	
	userChoice := w.readUserChoiceWithTimeout(w.config.PromptTimeout)
	return strings.ToLower(strings.TrimSpace(userChoice)) == "y"
}

// presentSecurityWarning shows a security warning for high-risk commands
func (w *Wrapper) presentSecurityWarning(command string, sanitizationResult *plugin.SanitizationResult) {
	// Use colors if enabled
	var (
		errorColor = ""
		resetColor = ""
	)

	if w.config.EnableColors {
		errorColor = "\033[31m" // Red
		resetColor = "\033[0m"  // Reset
	}

	fmt.Printf("\n%s🛑 SECURITY ALERT%s\n", errorColor, resetColor)
	fmt.Printf("%s│%s Command contains sensitive information and cannot be sent to AI\n", errorColor, resetColor)
	fmt.Printf("%s│%s Risk Level: %s\n", errorColor, resetColor, sanitizationResult.RiskLevel.String())
	fmt.Printf("%s│%s Detected patterns: %v\n", errorColor, resetColor, sanitizationResult.DetectedPatterns)
	fmt.Printf("%s│%s Please review and fix the command manually\n", errorColor, resetColor)
	fmt.Printf("%s╰─%s\n", errorColor, resetColor)
}
