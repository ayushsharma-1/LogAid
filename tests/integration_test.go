package tests

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ayush-1/logaid/internal/config"
	"github.com/ayush-1/logaid/internal/engine"
	"github.com/ayush-1/logaid/internal/plugins"
)

// TestIntegrationRedisCliScenario tests the complete flow for rediscli typo
func TestIntegrationRedisCliScenario(t *testing.T) {
	// Set up test environment
	os.Setenv("AI_PROVIDER", "gemini")
	os.Setenv("GEMINI_API_KEY", "test-key")
	os.Setenv("TEST_MODE", "true")

	// Initialize config
	if err := config.Init(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	// Create engine
	eng := engine.New()
	if eng == nil {
		t.Fatal("Failed to create engine")
	}

	// Test command and output
	command := "sudo apt install rediscli"
	output := "E: Unable to locate package rediscli"

	// Process the error
	ctx := context.Background()
	suggestion, err := eng.ProcessError(ctx, command, output)

	if err != nil {
		t.Logf("Error processing command (expected in test mode): %v", err)
	}

	// Verify suggestion was generated
	if suggestion == "" {
		t.Error("Expected a suggestion but got empty string")
	}

	t.Logf("Command: %s", command)
	t.Logf("Output: %s", output)
	t.Logf("Suggestion: %s", suggestion)
}

// TestRealWorldScenarios tests multiple real-world error scenarios
func TestRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		command     string
		output      string
		expectFix   bool
		description string
	}{
		// APT scenarios
		{
			name:        "Redis CLI typo",
			command:     "sudo apt install rediscli",
			output:      "E: Unable to locate package rediscli",
			expectFix:   true,
			description: "Most common Redis CLI installation typo",
		},
		{
			name:        "Node.js installation",
			command:     "sudo apt install node",
			output:      "E: Unable to locate package node",
			expectFix:   true,
			description: "Node.js package name confusion",
		},
		{
			name:        "Python pip installation",
			command:     "sudo apt install pip",
			output:      "E: Unable to locate package pip",
			expectFix:   true,
			description: "Python pip package name confusion",
		},
		{
			name:        "Docker installation typo",
			command:     "sudo apt install docker",
			output:      "E: Unable to locate package docker",
			expectFix:   true,
			description: "Docker package name typo",
		},
		{
			name:        "MySQL installation",
			command:     "sudo apt install mysql",
			output:      "E: Unable to locate package mysql",
			expectFix:   true,
			description: "MySQL package name confusion",
		},
		{
			name:        "APT permission error",
			command:     "apt install htop",
			output:      "E: Could not open lock file /var/lib/dpkg/lock-frontend - open (13: Permission denied)",
			expectFix:   true,
			description: "Missing sudo for apt install",
		},
		{
			name:        "APT lock error",
			command:     "sudo apt install vim",
			output:      "E: Could not get lock /var/lib/dpkg/lock-frontend - open (11: Resource temporarily unavailable)",
			expectFix:   true,
			description: "DPKG lock conflict",
		},

		// Git scenarios
		{
			name:        "Git checkout typo",
			command:     "git checout main",
			output:      "git: 'checout' is not a git command. See 'git --help'.",
			expectFix:   true,
			description: "Most common git checkout typo",
		},
		{
			name:        "Git commit typo",
			command:     "git comit -m 'fix bug'",
			output:      "git: 'comit' is not a git command. See 'git --help'.",
			expectFix:   true,
			description: "Git commit command typo",
		},
		{
			name:        "Git status typo",
			command:     "git stat",
			output:      "git: 'stat' is not a git command. See 'git --help'.",
			expectFix:   true,
			description: "Git status command typo",
		},
		{
			name:        "Git push typo",
			command:     "git pus origin main",
			output:      "git: 'pus' is not a git command. See 'git --help'.",
			expectFix:   true,
			description: "Git push command typo",
		},

		// Docker scenarios
		{
			name:        "Docker Ubuntu typo",
			command:     "docker run ubntu",
			output:      "Unable to find image 'ubntu:latest' locally",
			expectFix:   true,
			description: "Ubuntu Docker image typo",
		},
		{
			name:        "Docker run typo",
			command:     "docker ru nginx",
			output:      "docker: 'ru' is not a docker command.",
			expectFix:   true,
			description: "Docker run command typo",
		},
		{
			name:        "Docker permission error",
			command:     "docker ps",
			output:      "docker: Got permission denied while trying to connect to the Docker daemon socket",
			expectFix:   true,
			description: "Docker permission denied",
		},
		{
			name:        "Docker nginx typo",
			command:     "docker run ngnix",
			output:      "Unable to find image 'ngnix:latest' locally",
			expectFix:   true,
			description: "Nginx Docker image typo",
		},

		// NPM scenarios
		{
			name:        "NPM install typo",
			command:     "npm instal express",
			output:      "Unknown command: \"instal\"",
			expectFix:   true,
			description: "NPM install command typo",
		},
		{
			name:        "NPM package typo - express",
			command:     "npm install expres",
			output:      "npm ERR! 404 Not Found - GET https://registry.npmjs.org/expres - Not found",
			expectFix:   true,
			description: "Express package name typo",
		},
		{
			name:        "NPM package typo - lodash",
			command:     "npm install lodas",
			output:      "npm ERR! 404 Not Found - GET https://registry.npmjs.org/lodas - Not found",
			expectFix:   true,
			description: "Lodash package name typo",
		},
		{
			name:        "NPM start typo",
			command:     "npm stat",
			output:      "Unknown command: \"stat\"",
			expectFix:   true,
			description: "NPM start command typo",
		},

		// Python pip scenarios
		{
			name:        "Pip requests typo",
			command:     "pip install request",
			output:      "ERROR: Could not find a version that satisfies the requirement request",
			expectFix:   true,
			description: "Python requests package typo",
		},
		{
			name:        "Pip beautiful soup typo",
			command:     "pip install beautifulsoup",
			output:      "ERROR: Could not find a version that satisfies the requirement beautifulsoup",
			expectFix:   true,
			description: "BeautifulSoup package name typo",
		},
		{
			name:        "Pip permission error",
			command:     "pip install -g numpy",
			output:      "ERROR: Permission denied: '/usr/local/lib/python3.9/site-packages/'",
			expectFix:   true,
			description: "Pip global install permission error",
		},

		// Complex scenarios that should be handled by AI
		{
			name:        "Complex dependency conflict",
			command:     "sudo apt install libssl-dev",
			output:      "The following packages have unmet dependencies:\n libssl-dev : Depends: libssl1.1 (= 1.1.1f-1ubuntu2.16) but 1.1.1f-1ubuntu2.17 is to be installed",
			expectFix:   true,
			description: "Complex dependency conflict requiring AI",
		},
		{
			name:        "Git merge conflict",
			command:     "git merge feature-branch",
			output:      "Auto-merging file.txt\nCONFLICT (content): Merge conflict in file.txt\nAutomatic merge failed; fix conflicts and then commit the result.",
			expectFix:   true,
			description: "Git merge conflict requiring AI guidance",
		},
		{
			name:        "Docker daemon not running",
			command:     "docker run ubuntu",
			output:      "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?",
			expectFix:   true,
			description: "Docker daemon service issue",
		},
	}

	// Set up test environment
	os.Setenv("AI_PROVIDER", "gemini")
	os.Setenv("GEMINI_API_KEY", "test-key")
	os.Setenv("TEST_MODE", "true")

	config.Init()
	eng := engine.New()

	if eng == nil {
		t.Fatal("Failed to create engine")
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			suggestion, err := eng.ProcessError(ctx, scenario.command, scenario.output)

			if scenario.expectFix {
				if suggestion == "" {
					t.Errorf("Expected a suggestion for %s but got empty string", scenario.description)
				} else {
					t.Logf("✅ %s", scenario.description)
					t.Logf("   Command: %s", scenario.command)
					t.Logf("   Error: %s", strings.ReplaceAll(scenario.output, "\n", "\\n"))
					t.Logf("   Suggestion: %s", suggestion)
				}
			}

			// Log any errors for debugging (expected in test mode)
			if err != nil {
				t.Logf("   Error: %v (expected in test mode)", err)
			}
		})
	}
}

// TestPluginPerformance benchmarks plugin performance
func TestPluginPerformance(t *testing.T) {
	plugins := []plugins.Plugin{
		&plugins.AptPlugin{},
		&plugins.GitPlugin{},
		&plugins.DockerPlugin{},
		&plugins.NpmPlugin{},
		&plugins.PipPlugin{},
		&plugins.SystemctlPlugin{},
	}

	testCommands := []struct {
		command string
		output  string
	}{
		{"sudo apt install rediscli", "E: Unable to locate package rediscli"},
		{"git checout main", "git: 'checout' is not a git command"},
		{"docker run ubntu", "Unable to find image 'ubntu:latest' locally"},
		{"npm instal express", "Unknown command: \"instal\""},
		{"pip install request", "ERROR: Could not find a version that satisfies the requirement request"},
		{"systemctl start apache", "Unit apache.service not found"},
	}

	for _, plugin := range plugins {
		t.Run(plugin.Name()+"_Performance", func(t *testing.T) {
			start := time.Now()
			iterations := 1000

			for i := 0; i < iterations; i++ {
				for _, cmd := range testCommands {
					plugin.Match(cmd.command, cmd.output)
				}
			}

			duration := time.Since(start)
			avgTime := duration / time.Duration(iterations*len(testCommands))

			t.Logf("Plugin %s: %d iterations in %v (avg: %v per match)",
				plugin.Name(), iterations*len(testCommands), duration, avgTime)

			// Performance should be under 1ms per match
			if avgTime > time.Millisecond {
				t.Errorf("Plugin %s is too slow: %v per match", plugin.Name(), avgTime)
			}
		})
	}
}

// TestEdgeCasesAndErrorHandling tests edge cases and error conditions
func TestEdgeCasesAndErrorHandling(t *testing.T) {
	edgeCases := []struct {
		name        string
		command     string
		output      string
		description string
	}{
		{
			name:        "Empty command",
			command:     "",
			output:      "bash: command not found",
			description: "Empty command string",
		},
		{
			name:        "Empty output",
			command:     "ls",
			output:      "",
			description: "Empty output string",
		},
		{
			name:        "Very long command",
			command:     "sudo apt install " + strings.Repeat("very-long-package-name-", 100),
			output:      "E: Unable to locate package " + strings.Repeat("very-long-package-name-", 100),
			description: "Extremely long command and output",
		},
		{
			name:        "Binary data in output",
			command:     "cat /bin/ls",
			output:      "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F",
			description: "Binary data in command output",
		},
		{
			name:        "Unicode characters",
			command:     "echo 'こんにちは 🌍'",
			output:      "bash: こんにちは: command not found",
			description: "Unicode and emoji characters",
		},
		{
			name:        "Special characters",
			command:     "ls -la ~/.ssh/id_rsa | grep '^-'",
			output:      "bash: syntax error near unexpected token '|'",
			description: "Special shell characters",
		},
		{
			name:        "Multiple error types",
			command:     "sudo apt install rediscli && git checout main && docker run ubntu",
			output:      "E: Unable to locate package rediscli",
			description: "Command with multiple potential errors",
		},
	}

	os.Setenv("TEST_MODE", "true")
	config.Init()
	eng := engine.New()

	if eng == nil {
		t.Fatal("Failed to create engine")
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Engine panicked on edge case %s: %v", tc.description, r)
				}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			suggestion, err := eng.ProcessError(ctx, tc.command, tc.output)

			// Log results for analysis
			t.Logf("Edge case: %s", tc.description)
			t.Logf("Command: %q", tc.command)
			t.Logf("Output: %q", tc.output)
			t.Logf("Suggestion: %q", suggestion)
			if err != nil {
				t.Logf("Error: %v", err)
			}
		})
	}
}

// TestConcurrentProcessing tests thread safety and concurrent processing
func TestConcurrentProcessing(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	config.Init()
	eng := engine.New()

	if eng == nil {
		t.Fatal("Failed to create engine")
	}

	// Test concurrent processing of different commands
	commands := []struct {
		command string
		output  string
	}{
		{"sudo apt install rediscli", "E: Unable to locate package rediscli"},
		{"git checout main", "git: 'checout' is not a git command"},
		{"docker run ubntu", "Unable to find image 'ubntu:latest' locally"},
		{"npm instal express", "Unknown command: \"instal\""},
		{"pip install request", "ERROR: Could not find a version that satisfies the requirement request"},
	}

	const numWorkers = 10
	const iterationsPerWorker = 10

	results := make(chan error, numWorkers*iterationsPerWorker)

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			for j := 0; j < iterationsPerWorker; j++ {
				cmd := commands[j%len(commands)]
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

				_, err := eng.ProcessError(ctx, cmd.command, cmd.output)
				results <- err
				cancel()
			}
		}(i)
	}

	// Collect results
	var errors []error
	for i := 0; i < numWorkers*iterationsPerWorker; i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	t.Logf("Concurrent processing completed. Errors: %d/%d", len(errors), numWorkers*iterationsPerWorker)

	// In test mode, we expect some errors due to mock AI responses
	// The important thing is that no panics occurred and the system remained stable
}
