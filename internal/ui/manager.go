package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ayushsharma-1/LogAid/internal/ai"
	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/fatih/color"
)

// Manager handles the user interface and display logic
type Manager struct {
	config       *config.Config
	colorEnabled bool
	loading      bool
	loadingDone  chan bool
}

// NewManager creates a new UI manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:       cfg,
		colorEnabled: cfg.EnableColors,
		loadingDone:  make(chan bool, 1),
	}
}

// ShowLoading displays a loading indicator
func (m *Manager) ShowLoading(message string) {
	if m.loading {
		return
	}
	m.loading = true

	// Start loading animation in a goroutine
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0

		for {
			select {
			case <-m.loadingDone:
				// Clear the loading line
				fmt.Print("\r" + strings.Repeat(" ", len(message)+3) + "\r")
				return
			default:
				if m.colorEnabled {
					fmt.Printf("\r%s %s", color.CyanString(spinner[i]), color.WhiteString(message))
				} else {
					fmt.Printf("\r%s %s", spinner[i], message)
				}
				i = (i + 1) % len(spinner)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// HideLoading hides the loading indicator
func (m *Manager) HideLoading() {
	if !m.loading {
		return
	}
	m.loading = false
	m.loadingDone <- true
}

// ShowSuggestion displays an AI-generated suggestion with calm UX
func (m *Manager) ShowSuggestion(suggestion *ai.Suggestion) {
	if suggestion == nil {
		return
	}

	// Clear any previous output and add some spacing
	fmt.Println()

	// Show the suggestion with calm, non-intrusive styling
	if m.colorEnabled {
		// Header
		fmt.Printf("%s %s\n",
			color.New(color.FgBlue, color.Bold).Sprint("💡"),
			color.New(color.FgBlue, color.Bold).Sprint("LogAid Suggestion"))

		// Explanation
		if suggestion.Explanation != "" {
			fmt.Printf("%s %s\n",
				color.New(color.FgWhite).Sprint("   "),
				color.New(color.FgWhite).Sprint(suggestion.Explanation))
		}

		// Command (if available)
		if suggestion.Command != "" {
			fmt.Printf("%s %s\n",
				color.New(color.FgGreen, color.Bold).Sprint("   $"),
				color.New(color.FgGreen).Sprint(suggestion.Command))
		}

		// Instructions
		if suggestion.Command != "" {
			fmt.Printf("%s %s\n",
				color.New(color.FgYellow).Sprint("   "),
				color.New(color.FgYellow).Sprint("Press ↑ and Enter to run this command, or continue typing"))
		}

		// Confidence indicator (subtle)
		if suggestion.Confidence >= 0.8 {
			fmt.Printf("%s %s\n",
				color.New(color.FgGreen).Sprint("   "),
				color.New(color.FgGreen).Sprint("High confidence"))
		} else if suggestion.Confidence >= 0.6 {
			fmt.Printf("%s %s\n",
				color.New(color.FgYellow).Sprint("   "),
				color.New(color.FgYellow).Sprint("Medium confidence"))
		} else {
			fmt.Printf("%s %s\n",
				color.New(color.FgRed).Sprint("   "),
				color.New(color.FgRed).Sprint("Low confidence - please verify"))
		}
	} else {
		// Plain text version
		fmt.Printf("💡 LogAid Suggestion\n")
		if suggestion.Explanation != "" {
			fmt.Printf("   %s\n", suggestion.Explanation)
		}
		if suggestion.Command != "" {
			fmt.Printf("   $ %s\n", suggestion.Command)
			fmt.Printf("   Press ↑ and Enter to run this command, or continue typing\n")
		}

		// Confidence indicator
		if suggestion.Confidence >= 0.8 {
			fmt.Printf("   High confidence\n")
		} else if suggestion.Confidence >= 0.6 {
			fmt.Printf("   Medium confidence\n")
		} else {
			fmt.Printf("   Low confidence - please verify\n")
		}
	}

	fmt.Println() // Add spacing after suggestion
}

// ShowError displays an error message
func (m *Manager) ShowError(message string) {
	if m.colorEnabled {
		fmt.Printf("%s %s\n",
			color.New(color.FgRed, color.Bold).Sprint("❌"),
			color.New(color.FgRed).Sprint(message))
	} else {
		fmt.Printf("❌ %s\n", message)
	}
}

// ShowInfo displays an informational message
func (m *Manager) ShowInfo(message string) {
	if m.colorEnabled {
		fmt.Printf("%s %s\n",
			color.New(color.FgBlue).Sprint("ℹ️"),
			color.New(color.FgWhite).Sprint(message))
	} else {
		fmt.Printf("ℹ️  %s\n", message)
	}
}

// ShowSuccess displays a success message
func (m *Manager) ShowSuccess(message string) {
	if m.colorEnabled {
		fmt.Printf("%s %s\n",
			color.New(color.FgGreen, color.Bold).Sprint("✅"),
			color.New(color.FgGreen).Sprint(message))
	} else {
		fmt.Printf("✅ %s\n", message)
	}
}

// ShowWelcome displays the welcome message
func (m *Manager) ShowWelcome() {
	if m.colorEnabled {
		fmt.Printf("%s %s\n",
			color.New(color.FgMagenta, color.Bold).Sprint("🚀"),
			color.New(color.FgMagenta, color.Bold).Sprint("LogAid Flow-State Agent"))
		fmt.Printf("%s %s\n",
			color.New(color.FgWhite).Sprint("   "),
			color.New(color.FgWhite).Sprint("AI-powered CLI error detection and suggestions"))
		fmt.Printf("%s %s\n",
			color.New(color.FgHiBlack).Sprint("   "),
			color.New(color.FgHiBlack).Sprint("Press Ctrl+C to exit"))
	} else {
		fmt.Printf("🚀 LogAid Flow-State Agent\n")
		fmt.Printf("   AI-powered CLI error detection and suggestions\n")
		fmt.Printf("   Press Ctrl+C to exit\n")
	}
	fmt.Println()
}

// PromptUser prompts the user for input with a message
func (m *Manager) PromptUser(message string) (string, error) {
	if m.colorEnabled {
		fmt.Printf("%s %s: ",
			color.New(color.FgYellow).Sprint("?"),
			color.New(color.FgWhite).Sprint(message))
	} else {
		fmt.Printf("? %s: ", message)
	}

	var input string
	_, err := fmt.Scanln(&input)
	return input, err
}

// ClearScreen clears the terminal screen
func (m *Manager) ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// MoveCursorUp moves the cursor up by n lines
func (m *Manager) MoveCursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}

// MoveCursorDown moves the cursor down by n lines
func (m *Manager) MoveCursorDown(n int) {
	fmt.Printf("\033[%dB", n)
}

// ClearLine clears the current line
func (m *Manager) ClearLine() {
	fmt.Print("\033[2K\r")
}
