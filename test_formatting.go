package main

import (
	"fmt"
)

// Simulate the new LogAid formatting
func simulateLogAidOutput() {
	fmt.Println("=== Before (Bad Formatting) ===")
	fmt.Printf("\n🔍 LogAid is analyzing the error...\n")
	fmt.Printf("\n💡 LogAid Suggestion:\n")
	fmt.Printf("   Command: ls ltr\n")
	fmt.Printf("   Error: ls: cannot access 'ltr': No such file or directory\n")
	fmt.Printf("\n🤖 AI Analysis:\n")
	fmt.Printf("   The error message indicates that a file or directory named 'ltr' does not exist...\n")
	fmt.Printf("   Confidence: 0.95\n\n")

	fmt.Println("=== After (Clean Formatting) ===")
	// Clear line and show loading
	fmt.Printf("🔍 Analyzing...")

	// Simulate quick processing
	fmt.Printf("\r\033[K") // Clear line

	// Show clean suggestion
	fmt.Printf("💡 \033[1;36mLogAid:\033[0m Try \033[1;32mls -ltr\033[0m (list files by time)\n")

	fmt.Println("\n=== Other Examples ===")
	fmt.Printf("💡 \033[1;36mLogAid:\033[0m Try \033[1;32mnpm install\033[0m or \033[1;32mnpm run dev\033[0m\n")
	fmt.Printf("💡 \033[1;36mLogAid:\033[0m Try \033[1;32mgit status\033[0m\n")
	fmt.Printf("💡 \033[1;36mLogAid:\033[0m Command '\033[1;31mdocker\033[0m' not found. Check if it's installed\n")
}

func main() {
	simulateLogAidOutput()
}
