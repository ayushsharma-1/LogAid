package tests

import (
	"strings"
	"testing"

	"github.com/ayushsharma-1/LogAid/internal/plugins"
)

// TestAptPlugin tests the APT plugin with comprehensive test cases
func TestAptPlugin(t *testing.T) {
	plugin := &plugins.AptPlugin{}

	testCases := []struct {
		name        string
		command     string
		output      string
		shouldMatch bool
		expectedFix string
		description string
	}{
		// Package not found errors
		{
			name:        "redis-cli typo",
			command:     "sudo apt install rediscli",
			output:      "E: Unable to locate package rediscli",
			shouldMatch: true,
			expectedFix: "sudo apt install redis-tools",
			description: "Most common Redis CLI typo",
		},
		{
			name:        "redis-client typo",
			command:     "apt install redis-client",
			output:      "E: Unable to locate package redis-client",
			shouldMatch: true,
			expectedFix: "apt install redis-tools",
			description: "Alternative Redis client typo",
		},
		{
			name:        "mysql typo",
			command:     "sudo apt install mysql",
			output:      "E: Unable to locate package mysql",
			shouldMatch: true,
			expectedFix: "sudo apt install mysql-server",
			description: "MySQL package name correction",
		},
		{
			name:        "nodejs typo",
			command:     "sudo apt install node",
			output:      "E: Unable to locate package node",
			shouldMatch: true,
			expectedFix: "sudo apt install nodejs npm",
			description: "Node.js package name correction",
		},
		{
			name:        "python typo",
			command:     "sudo apt install python",
			output:      "E: Unable to locate package python",
			shouldMatch: true,
			expectedFix: "sudo apt install python3",
			description: "Python 3 package correction",
		},
		{
			name:        "pip typo",
			command:     "sudo apt install pip",
			output:      "E: Unable to locate package pip",
			shouldMatch: true,
			expectedFix: "sudo apt install python3-pip",
			description: "Python pip package correction",
		},
		{
			name:        "docker typo",
			command:     "sudo apt install docker",
			output:      "E: Unable to locate package docker",
			shouldMatch: true,
			expectedFix: "sudo apt install docker.io",
			description: "Docker package name correction",
		},
		{
			name:        "java typo",
			command:     "sudo apt install java",
			output:      "E: Unable to locate package java",
			shouldMatch: true,
			expectedFix: "sudo apt install openjdk-11-jdk",
			description: "Java JDK package correction",
		},
		{
			name:        "gcc typo",
			command:     "sudo apt install gcc",
			output:      "E: Unable to locate package gcc",
			shouldMatch: true,
			expectedFix: "sudo apt install build-essential",
			description: "GCC build tools correction",
		},
		{
			name:        "make typo",
			command:     "sudo apt install make",
			output:      "E: Unable to locate package make",
			shouldMatch: true,
			expectedFix: "sudo apt install build-essential",
			description: "Make build tools correction",
		},
		{
			name:        "vim typo",
			command:     "sudo apt install vim",
			output:      "E: Unable to locate package vim",
			shouldMatch: true,
			expectedFix: "sudo apt install vim-gtk3",
			description: "Vim editor correction",
		},
		{
			name:        "chrome typo",
			command:     "sudo apt install chrome",
			output:      "E: Unable to locate package chrome",
			shouldMatch: true,
			expectedFix: "sudo apt install google-chrome-stable",
			description: "Chrome browser correction",
		},
		{
			name:        "vscode typo",
			command:     "sudo apt install vscode",
			output:      "E: Unable to locate package vscode",
			shouldMatch: true,
			expectedFix: "sudo apt install code",
			description: "VSCode package correction",
		},
		{
			name:        "postgres typo",
			command:     "sudo apt install postgres",
			output:      "E: Unable to locate package postgres",
			shouldMatch: true,
			expectedFix: "sudo apt install postgresql postgresql-contrib",
			description: "PostgreSQL package correction",
		},
		{
			name:        "sqlite typo",
			command:     "sudo apt install sqlite",
			output:      "E: Unable to locate package sqlite",
			shouldMatch: true,
			expectedFix: "sudo apt install sqlite3",
			description: "SQLite package correction",
		},

		// Permission errors
		{
			name:        "permission denied",
			command:     "apt install redis-tools",
			output:      "E: Could not open lock file /var/lib/dpkg/lock-frontend - open (13: Permission denied)",
			shouldMatch: true,
			expectedFix: "sudo apt install redis-tools",
			description: "Missing sudo for install",
		},
		{
			name:        "permission denied update",
			command:     "apt update",
			output:      "Reading package lists... Done\nE: Could not open lock file /var/lib/dpkg/lock-frontend - open (13: Permission denied)",
			shouldMatch: true,
			expectedFix: "sudo apt update",
			description: "Missing sudo for update",
		},

		// Lock errors
		{
			name:        "dpkg lock error",
			command:     "sudo apt install htop",
			output:      "E: Could not get lock /var/lib/dpkg/lock-frontend - open (11: Resource temporarily unavailable)",
			shouldMatch: true,
			expectedFix: "sudo killall apt apt-get dpkg && sudo dpkg --configure -a && sudo apt update && sudo apt install htop",
			description: "DPKG lock conflict resolution",
		},
		{
			name:        "apt lock error",
			command:     "sudo apt update",
			output:      "E: Could not get lock /var/lib/apt/lists/lock - open (11: Resource temporarily unavailable)",
			shouldMatch: true,
			expectedFix: "sudo killall apt apt-get dpkg && sudo dpkg --configure -a && sudo apt update && sudo apt update",
			description: "APT lock conflict resolution",
		},

		// Repository errors
		{
			name:        "404 repository error",
			command:     "sudo apt install some-package",
			output:      "Err:1 http://archive.ubuntu.com/ubuntu focal/main amd64 Packages\n  404  Not Found [IP: 91.189.91.38 80]",
			shouldMatch: true,
			expectedFix: "sudo apt update && sudo apt install some-package",
			description: "Repository 404 error - needs update",
		},
		{
			name:        "failed to fetch",
			command:     "sudo apt install curl",
			output:      "E: Failed to fetch http://archive.ubuntu.com/ubuntu/dists/focal/Release",
			shouldMatch: true,
			expectedFix: "sudo apt update && sudo apt install curl",
			description: "Failed to fetch - needs update",
		},

		// Complex cases that should use AI
		{
			name:        "complex dependency error",
			command:     "sudo apt install libssl-dev",
			output:      "The following packages have unmet dependencies:\n libssl-dev : Depends: libssl1.1 (= 1.1.1f-1ubuntu2.16) but 1.1.1f-1ubuntu2.17 is to be installed",
			shouldMatch: true,
			expectedFix: "", // Should fallback to AI
			description: "Complex dependency conflict",
		},
		{
			name:        "broken packages",
			command:     "sudo apt install python3-dev",
			output:      "You have held broken packages.\nThe following packages have unmet dependencies:",
			shouldMatch: true,
			expectedFix: "", // Should fallback to AI
			description: "Broken packages issue",
		},

		// Non-matching cases
		{
			name:        "successful install",
			command:     "sudo apt install htop",
			output:      "Reading package lists... Done\nBuilding dependency tree... Done\nThe following NEW packages will be installed:\n  htop",
			shouldMatch: false,
			expectedFix: "",
			description: "Successful installation - should not match",
		},
		{
			name:        "non-apt command",
			command:     "ls -la",
			output:      "total 48\ndrwxr-xr-x 12 user user 4096 Jan 26 10:00 .",
			shouldMatch: false,
			expectedFix: "",
			description: "Non-APT command - should not match",
		},
		{
			name:        "npm command with error",
			command:     "npm install express",
			output:      "npm ERR! code ENOTFOUND",
			shouldMatch: false,
			expectedFix: "",
			description: "NPM command - should not match APT plugin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Match function
			matches := plugin.Match(tc.command, tc.output)
			if matches != tc.shouldMatch {
				t.Errorf("Match() = %v, want %v for case: %s", matches, tc.shouldMatch, tc.description)
			}

			// Test Suggest function (only if it should match)
			if tc.shouldMatch && tc.expectedFix != "" {
				suggestion := plugin.Suggest(tc.command, tc.output)
				if suggestion != tc.expectedFix {
					t.Errorf("Suggest() = %q, want %q for case: %s", suggestion, tc.expectedFix, tc.description)
				}
			}
		})
	}
}

// BenchmarkAptPlugin benchmarks the APT plugin performance
func BenchmarkAptPlugin(b *testing.B) {
	plugin := &plugins.AptPlugin{}
	command := "sudo apt install rediscli"
	output := "E: Unable to locate package rediscli"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plugin.Match(command, output)
		plugin.Suggest(command, output)
	}
}

// TestAptPluginEdgeCases tests edge cases and error conditions
func TestAptPluginEdgeCases(t *testing.T) {
	plugin := &plugins.AptPlugin{}

	edgeCases := []struct {
		name        string
		command     string
		output      string
		description string
	}{
		{
			name:        "empty command",
			command:     "",
			output:      "E: Unable to locate package test",
			description: "Empty command should not match",
		},
		{
			name:        "empty output",
			command:     "sudo apt install test",
			output:      "",
			description: "Empty output should not match",
		},
		{
			name:        "malformed command",
			command:     "apt",
			output:      "E: Unable to locate package",
			description: "Incomplete command",
		},
		{
			name:        "very long package name",
			command:     "sudo apt install " + strings.Repeat("a", 1000),
			output:      "E: Unable to locate package " + strings.Repeat("a", 1000),
			description: "Very long package name",
		},
		{
			name:        "special characters in command",
			command:     "sudo apt install 'test-package'",
			output:      "E: Unable to locate package 'test-package'",
			description: "Special characters in package name",
		},
		{
			name:        "multiple packages",
			command:     "sudo apt install rediscli mysql nodejs",
			output:      "E: Unable to locate package rediscli",
			description: "Multiple packages with one error",
		},
		{
			name:        "case sensitivity",
			command:     "sudo APT INSTALL REDISCLI",
			output:      "E: Unable to locate package REDISCLI",
			description: "Uppercase command should still match",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Plugin panicked on edge case %s: %v", tc.description, r)
				}
			}()

			matches := plugin.Match(tc.command, tc.output)
			if matches {
				plugin.Suggest(tc.command, tc.output) // Should not panic
			}
		})
	}
}
