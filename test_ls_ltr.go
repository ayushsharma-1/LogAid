package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ayushsharma-1/LogAid/internal/ai"
	"github.com/ayushsharma-1/LogAid/internal/config"
)

func testLsLtr() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create AI client
	aiClient, err := ai.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create AI client:", err)
	}

	// Simulate the "ls ltr" error
	command := "ls ltr"
	errorOutput := "ls: cannot access 'ltr': No such file or directory"

	fmt.Printf("🧪 Testing LogAid with command: %s\n", command)
	fmt.Printf("Error output: %s\n\n", errorOutput)

	// Get AI analysis
	suggestion, err := aiClient.AnalyzeError(context.Background(), command, errorOutput, 3)
	if err != nil {
		log.Fatal("Failed to analyze error:", err)
	}

	fmt.Printf("🤖 AI Analysis:\n")
	fmt.Printf("   Explanation: %s\n", suggestion.Explanation)
	fmt.Printf("   Confidence: %.2f\n\n", suggestion.Confidence)

	// Show manual suggestion for ls command
	if strings.Contains(command, "ls") && strings.Contains(errorOutput, "No such file or directory") {
		fmt.Printf("🔧 Smart Analysis:\n")
		fmt.Printf("   This looks like you meant 'ls -ltr' (list files in long format, sorted by time)\n")
		fmt.Printf("   💡 Try: ls -ltr\n")
		fmt.Printf("   💡 Or:  ls -la (list all files in long format)\n")
		fmt.Printf("   💡 Or:  ls -lt (list files in long format, sorted by time)\n")
	}
}

func main() {
	testLsLtr()
}
