package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/spf13/cobra"
)

// diagCmd represents the diag command
var diagCmd = &cobra.Command{
	Use:   "diag",
	Short: "Diagnose LogAid and shell configuration issues",
	Long: `Diagnose potential issues with LogAid setup and shell configuration.
This command helps identify problems that might prevent LogAid from starting correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 LogAid Diagnostic Report")
		fmt.Println("===========================")
		fmt.Println()

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ Configuration Error: %v\n", err)
			return
		}

		fmt.Println("✅ Configuration loaded successfully")
		fmt.Printf("   Shell: %s\n", cfg.Shell)
		fmt.Printf("   AI Provider: %s\n", cfg.AIProvider)
		fmt.Println()

		// Check shell availability
		fmt.Println("🐚 Shell Diagnostics:")
		checkShell(cfg.Shell)
		fmt.Println()

		// Check shell configuration files
		fmt.Println("📄 Shell Configuration Files:")
		checkShellConfig(cfg.Shell)
		fmt.Println()

		// Check AI connectivity
		fmt.Println("🤖 AI Connectivity:")
		if cfg.APIKey != "" {
			fmt.Printf("✅ API key configured (%s...%s)\n",
				cfg.APIKey[:4], cfg.APIKey[len(cfg.APIKey)-4:])
		} else {
			fmt.Println("❌ No API key configured")
		}
		fmt.Println()

		// Environment checks
		fmt.Println("🌍 Environment:")
		checkEnvironment()
		fmt.Println()

		// Recommendations
		fmt.Println("💡 Recommendations:")
		fmt.Println("   1. If you see shell errors, try running:")
		fmt.Println("      logaid run --minimal")
		fmt.Println()
		fmt.Println("   2. To fix .bashrc issues, backup and edit:")
		fmt.Println("      cp ~/.bashrc ~/.bashrc.backup")
		fmt.Println("      nano ~/.bashrc")
		fmt.Println()
		fmt.Println("   3. For VS Code integration issues, restart VS Code")
		fmt.Println("      or run LogAid outside of VS Code terminal")
		fmt.Println()
		fmt.Println("   4. Test with a clean shell:")
		fmt.Println("      bash --norc --noprofile")
		fmt.Println("      logaid run")
	},
}

func checkShell(shellPath string) {
	// Check if shell exists
	if _, err := os.Stat(shellPath); err != nil {
		fmt.Printf("❌ Shell not found: %s\n", shellPath)
		return
	}
	fmt.Printf("✅ Shell found: %s\n", shellPath)

	// Check if shell is executable
	if err := exec.Command(shellPath, "--version").Run(); err != nil {
		fmt.Printf("⚠️  Warning: Shell may not be working correctly: %v\n", err)
	} else {
		fmt.Printf("✅ Shell is executable\n")
	}

	// Get shell version
	if output, err := exec.Command(shellPath, "--version").Output(); err == nil {
		version := strings.Split(strings.TrimSpace(string(output)), "\n")[0]
		fmt.Printf("   Version: %s\n", version)
	}
}

func checkShellConfig(shellPath string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("❌ Cannot determine home directory: %v\n", err)
		return
	}

	var configFiles []string
	shellName := filepath.Base(shellPath)

	switch shellName {
	case "bash":
		configFiles = []string{".bashrc", ".bash_profile", ".profile"}
	case "zsh":
		configFiles = []string{".zshrc", ".zprofile", ".profile"}
	case "fish":
		configFiles = []string{".config/fish/config.fish"}
	default:
		configFiles = []string{".profile"}
	}

	for _, configFile := range configFiles {
		fullPath := filepath.Join(home, configFile)
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("📄 Found: %s\n", configFile)

			// Check for common problematic patterns
			checkConfigFileIssues(fullPath)
		} else {
			fmt.Printf("   Missing: %s (this is often okay)\n", configFile)
		}
	}
}

func checkConfigFileIssues(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("   ⚠️  Cannot read file: %v\n", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	issues := 0

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Check for common issues
		if strings.Contains(trimmed, "local ") && !strings.Contains(trimmed, "function") {
			fmt.Printf("   ⚠️  Line %d: 'local' used outside function\n", lineNum)
			issues++
		}

		if strings.Contains(trimmed, "__vsc_prompt_cmd") {
			fmt.Printf("   ⚠️  Line %d: VS Code terminal integration detected\n", lineNum)
			issues++
		}
	}

	if issues == 0 {
		fmt.Printf("   ✅ No obvious issues detected\n")
	} else {
		fmt.Printf("   ❌ Found %d potential issues\n", issues)
	}
}

func checkEnvironment() {
	// Check important environment variables
	envVars := []string{"PATH", "HOME", "SHELL", "TERM", "USER"}

	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value != "" {
			fmt.Printf("✅ %s: %s\n", envVar, value)
		} else {
			fmt.Printf("❌ %s: not set\n", envVar)
		}
	}

	// Check if we're in VS Code
	if os.Getenv("VSCODE_INJECTION") != "" || os.Getenv("TERM_PROGRAM") == "vscode" {
		fmt.Printf("📱 Running in VS Code terminal\n")
	}
}

func init() {
	rootCmd.AddCommand(diagCmd)
}
