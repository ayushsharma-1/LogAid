package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ayushsharma-1/LogAid/pkg/config"
)

// Interface handles user interface interactions
type Interface struct {
	config *config.Config
	reader *bufio.Reader
}

// NewInterface creates a new UI interface
func NewInterface(cfg *config.Config) *Interface {
	return &Interface{
		config: cfg,
		reader: bufio.NewReader(os.Stdin),
	}
}

// ShowLogo displays the LogAid ASCII logo
func (ui *Interface) ShowLogo() {
	logo := `
   _                _    _     _ 
  | |    ___   __ _| | _(_) __| |
  | |   / _ \ / _` + "`" + ` | |/ / |/ _` + "`" + ` |
  | |__| (_) | (_| | | <| | (_| |
  |_____\___/ \__,_|_|\_\_|\__,_|
             |___/               
      LogAid: Your CLI Guardian    
`
	if ui.config.EnableColors {
		fmt.Printf("\033[36m%s\033[0m\n", logo) // Cyan
	} else {
		fmt.Println(logo)
	}
}

// ShowPrompt displays the command prompt
func (ui *Interface) ShowPrompt() {
	if ui.config.EnableColors {
		fmt.Print("\033[32m❯\033[0m ") // Green arrow
	} else {
		fmt.Print("$ ")
	}
}

// PrintInfo prints an informational message
func (ui *Interface) PrintInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if ui.config.EnableColors {
		fmt.Printf("\033[34m[INFO]\033[0m %s\n", message) // Blue
	} else {
		fmt.Printf("[INFO] %s\n", message)
	}
}

// PrintError prints an error message
func (ui *Interface) PrintError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if ui.config.EnableColors {
		fmt.Printf("\033[31m[ERROR]\033[0m %s\n", message) // Red
	} else {
		fmt.Printf("[ERROR] %s\n", message)
	}
}

// PrintWarning prints a warning message
func (ui *Interface) PrintWarning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if ui.config.EnableColors {
		fmt.Printf("\033[33m[WARN]\033[0m %s\n", message) // Yellow
	} else {
		fmt.Printf("[WARN] %s\n", message)
	}
}

// PrintSuccess prints a success message
func (ui *Interface) PrintSuccess(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if ui.config.EnableColors {
		fmt.Printf("\033[32m[SUCCESS]\033[0m %s\n", message) // Green
	} else {
		fmt.Printf("[SUCCESS] %s\n", message)
	}
}

// ShowSuggestion displays a suggestion and prompts for user input
func (ui *Interface) ShowSuggestion(command, explanation string) bool {
	// Auto-confirm if enabled (for testing or automated scenarios)
	if ui.config.AutoConfirm {
		ui.PrintInfo("Auto-confirming suggestion: %s", command)
		return true
	}

	// Display suggestion
	if ui.config.EnableColors {
		fmt.Printf("\n\033[36m💡 LogAid Suggestion:\033[0m %s\n", command)
		if explanation != "" {
			fmt.Printf("\033[37m   Explanation: %s\033[0m\n", explanation)
		}
	} else {
		fmt.Printf("\nLogAid Suggestion: %s\n", command)
		if explanation != "" {
			fmt.Printf("   Explanation: %s\n", explanation)
		}
	}

	// Prompt for confirmation
	return ui.promptYesNo("Execute this command?")
}

// promptYesNo prompts the user for a yes/no response
func (ui *Interface) promptYesNo(question string) bool {
	if ui.config.EnableColors {
		fmt.Printf("\033[33m%s [y/N]:\033[0m ", question) // Yellow
	} else {
		fmt.Printf("%s [y/N]: ", question)
	}

	// Handle timeout
	if ui.config.SuggestionTimeout > 0 {
		return ui.promptWithTimeout()
	}

	// Read user input
	input, err := ui.reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

// promptWithTimeout prompts with a timeout (simplified implementation)
func (ui *Interface) promptWithTimeout() bool {
	// Create a channel for user input
	inputChan := make(chan string, 1)
	
	go func() {
		input, err := ui.reader.ReadString('\n')
		if err == nil {
			inputChan <- strings.ToLower(strings.TrimSpace(input))
		}
	}()

	// Wait for input or timeout
	select {
	case input := <-inputChan:
		return input == "y" || input == "yes"
	case <-time.After(time.Duration(ui.config.SuggestionTimeout) * time.Second):
		fmt.Println("\nTimeout reached. Defaulting to 'no'.")
		return false
	}
}

// ShowProgress displays a progress indicator
func (ui *Interface) ShowProgress(message string) {
	if ui.config.EnableColors {
		fmt.Printf("\033[34m⏳ %s...\033[0m", message) // Blue
	} else {
		fmt.Printf("⏳ %s...", message)
	}
}

// ClearProgress clears the progress indicator
func (ui *Interface) ClearProgress() {
	fmt.Print("\r\033[K") // Clear line
}

// ShowSpinner shows a simple spinner animation
func (ui *Interface) ShowSpinner(ctx <-chan struct{}, message string) {
	if !ui.config.EnableColors {
		fmt.Printf("%s...", message)
		return
	}

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	for {
		select {
		case <-ctx:
			fmt.Print("\r\033[K") // Clear line
			return
		default:
			fmt.Printf("\r\033[34m%s %s\033[0m", spinner[i%len(spinner)], message)
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}

// ShowConfidence displays confidence score if enabled
func (ui *Interface) ShowConfidence(confidence float64) {
	if !ui.config.ShowConfidenceScore {
		return
	}

	var color string
	var indicator string

	if confidence >= 0.8 {
		color = "\033[32m" // Green
		indicator = "🟢"
	} else if confidence >= 0.6 {
		color = "\033[33m" // Yellow
		indicator = "🟡"
	} else {
		color = "\033[31m" // Red
		indicator = "🔴"
	}

	if ui.config.EnableColors {
		fmt.Printf("%s   Confidence: %s%.1f%%\033[0m\n", indicator, color, confidence*100)
	} else {
		fmt.Printf("   Confidence: %.1f%%\n", confidence*100)
	}
}

// ShowMultipleSuggestions displays multiple suggestions and lets user choose
func (ui *Interface) ShowMultipleSuggestions(suggestions []SuggestionOption) int {
	if len(suggestions) == 0 {
		return -1
	}

	if len(suggestions) == 1 {
		if ui.ShowSuggestion(suggestions[0].Command, suggestions[0].Explanation) {
			return 0
		}
		return -1
	}

	// Display multiple options
	fmt.Println("\nMultiple suggestions available:")
	for i, suggestion := range suggestions {
		if ui.config.EnableColors {
			fmt.Printf("\033[36m%d.\033[0m %s\n", i+1, suggestion.Command)
			if suggestion.Explanation != "" {
				fmt.Printf("   \033[37m%s\033[0m\n", suggestion.Explanation)
			}
		} else {
			fmt.Printf("%d. %s\n", i+1, suggestion.Command)
			if suggestion.Explanation != "" {
				fmt.Printf("   %s\n", suggestion.Explanation)
			}
		}
		ui.ShowConfidence(suggestion.Confidence)
	}

	// Prompt for selection
	if ui.config.EnableColors {
		fmt.Printf("\033[33mSelect option (1-%d) or 0 to skip:\033[0m ", len(suggestions))
	} else {
		fmt.Printf("Select option (1-%d) or 0 to skip: ", len(suggestions))
	}

	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 0 || choice > len(suggestions) {
		return -1
	}

	if choice == 0 {
		return -1
	}

	return choice - 1
}

// SuggestionOption represents a suggestion option
type SuggestionOption struct {
	Command     string
	Explanation string
	Confidence  float64
}

// ShowBanner displays a banner message
func (ui *Interface) ShowBanner(message string) {
	border := strings.Repeat("=", len(message)+4)
	
	if ui.config.EnableColors {
		fmt.Printf("\033[36m%s\033[0m\n", border)
		fmt.Printf("\033[36m  %s  \033[0m\n", message)
		fmt.Printf("\033[36m%s\033[0m\n", border)
	} else {
		fmt.Println(border)
		fmt.Printf("  %s  \n", message)
		fmt.Println(border)
	}
}
