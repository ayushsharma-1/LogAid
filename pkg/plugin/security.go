package plugin

import (
	"regexp"
	"strings"
)

// SecuritySanitizer handles sensitive data detection and sanitization
type SecuritySanitizer struct {
	// Patterns for detecting sensitive information
	sensitivePatterns map[string]*regexp.Regexp
	
	// Safe command prefixes that are generally okay to send to AI
	safeCommandPrefixes []string
	
	// Dangerous commands that should never be sent to AI
	dangerousCommands []string
}

// NewSecuritySanitizer creates a new security sanitizer
func NewSecuritySanitizer() *SecuritySanitizer {
	sanitizer := &SecuritySanitizer{
		sensitivePatterns: make(map[string]*regexp.Regexp),
		safeCommandPrefixes: []string{
			"ls", "cd", "pwd", "echo", "cat", "less", "more", "head", "tail",
			"grep", "find", "which", "whereis", "man", "help", "history",
			"ps", "top", "htop", "free", "df", "du", "uname", "date", "uptime",
			"git status", "git log", "git diff", "git branch",
			"npm list", "npm info", "pip list", "pip show",
		},
		dangerousCommands: []string{
			"ssh", "scp", "rsync", "curl", "wget", "mysql", "psql", "mongo",
			"docker login", "kubectl", "aws", "gcloud", "az",
		},
	}

	// Define regex patterns for sensitive data
	sanitizer.sensitivePatterns = map[string]*regexp.Regexp{
		// API Keys and Tokens
		"api_key":         regexp.MustCompile(`(?i)(api[_-]?key|token|secret)[=:\s]+['""]?([a-zA-Z0-9_\-/+]{16,})['""]?`),
		"bearer_token":    regexp.MustCompile(`(?i)bearer[:\s]+([a-zA-Z0-9_\-/+\.]{16,})`),
		"authorization":   regexp.MustCompile(`(?i)authorization[:\s]+([a-zA-Z0-9_\-/+\.]{16,})`),
		
		// AWS Credentials
		"aws_access_key":  regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"aws_secret_key":  regexp.MustCompile(`(?i)(aws[_-]?secret[_-]?access[_-]?key)[=:\s]+['""]?([a-zA-Z0-9/+=]{40})['""]?`),
		
		// SSH Keys
		"ssh_private_key": regexp.MustCompile(`-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`),
		"ssh_public_key":  regexp.MustCompile(`ssh-(rsa|dss|ed25519|ecdsa)\s+[A-Za-z0-9+/]+[=]{0,2}`),
		
		// Database Connection Strings
		"db_connection":   regexp.MustCompile(`(?i)(mysql|postgresql|mongodb)://[^@]+:[^@]+@`),
		"db_password":     regexp.MustCompile(`(?i)(password|pwd)[=:\s]+['""]?([^'""\s]{4,})['""]?`),
		
		// Common Passwords
		"password":        regexp.MustCompile(`(?i)(password|passwd|pwd)[=:\s]+['""]?([^'""\s]{3,})['""]?`),
		
		// Credit Card Numbers
		"credit_card":     regexp.MustCompile(`\b(?:\d{4}[-\s]?){3}\d{4}\b`),
		
		// IP Addresses (private ranges might be sensitive)
		"private_ip":      regexp.MustCompile(`\b(?:10\.|172\.(?:1[6-9]|2[0-9]|3[0-1])\.|192\.168\.)\d+\.\d+\b`),
		
		// Email addresses (might contain sensitive info)
		"email":           regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		
		// File paths that might be sensitive
		"home_path":       regexp.MustCompile(`/home/[^/\s]+/\.(ssh|aws|config|env)`),
		"config_file":     regexp.MustCompile(`\.(env|config|ini|yaml|yml|json)$`),
		
		// URLs with credentials
		"url_with_creds":  regexp.MustCompile(`https?://[^@/]+:[^@/]+@`),
		
		// Environment variables
		"env_var":         regexp.MustCompile(`(?i)(export\s+)?[A-Z_]+=(.*)`),
	}

	return sanitizer
}

// SanitizationResult contains the result of sanitization
type SanitizationResult struct {
	IsSafe              bool     // Whether the data is safe to send to AI
	SanitizedCommand    string   // Command with sensitive parts redacted
	SanitizedOutput     string   // Output with sensitive parts redacted
	DetectedPatterns    []string // List of detected sensitive patterns
	RiskLevel           RiskLevel // Overall risk level
	UserConsentRequired bool     // Whether user consent is needed
}

// RiskLevel represents the security risk level
type RiskLevel int

const (
	RiskNone RiskLevel = iota
	RiskLow
	RiskMedium  
	RiskHigh
	RiskCritical
)

// String returns string representation of risk level
func (r RiskLevel) String() string {
	switch r {
	case RiskNone:
		return "None"
	case RiskLow:
		return "Low"
	case RiskMedium:
		return "Medium"
	case RiskHigh:
		return "High"
	case RiskCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// SanitizeData performs comprehensive security sanitization
func (s *SecuritySanitizer) SanitizeData(command, output string, exitCode int) *SanitizationResult {
	result := &SanitizationResult{
		IsSafe:           true,
		SanitizedCommand: command,
		SanitizedOutput:  output,
		DetectedPatterns: []string{},
		RiskLevel:        RiskNone,
	}

	// Check if command is in dangerous list
	if s.isDangerousCommand(command) {
		result.IsSafe = false
		result.RiskLevel = RiskCritical
		result.UserConsentRequired = true
		result.DetectedPatterns = append(result.DetectedPatterns, "dangerous_command")
		return result
	}

	// Check for sensitive patterns in command
	commandPatterns := s.detectSensitivePatterns(command)
	if len(commandPatterns) > 0 {
		result.DetectedPatterns = append(result.DetectedPatterns, commandPatterns...)
		result.SanitizedCommand = s.redactSensitiveData(command)
		result.RiskLevel = s.calculateRiskLevel(commandPatterns)
	}

	// Check for sensitive patterns in output
	outputPatterns := s.detectSensitivePatterns(output)
	if len(outputPatterns) > 0 {
		result.DetectedPatterns = append(result.DetectedPatterns, outputPatterns...)
		result.SanitizedOutput = s.redactSensitiveData(output)
		currentRisk := s.calculateRiskLevel(outputPatterns)
		if currentRisk > result.RiskLevel {
			result.RiskLevel = currentRisk
		}
	}

	// Determine if data is safe and if consent is required
	if result.RiskLevel >= RiskHigh {
		result.IsSafe = false
		result.UserConsentRequired = true
	} else if result.RiskLevel >= RiskMedium {
		result.UserConsentRequired = true
	}

	return result
}

// isDangerousCommand checks if a command is potentially dangerous
func (s *SecuritySanitizer) isDangerousCommand(command string) bool {
	cmdLower := strings.ToLower(strings.TrimSpace(command))
	
	for _, dangerous := range s.dangerousCommands {
		if strings.HasPrefix(cmdLower, dangerous) {
			return true
		}
	}
	
	// Check for commands with credentials in them
	if strings.Contains(cmdLower, "password") || 
	   strings.Contains(cmdLower, "token") ||
	   strings.Contains(cmdLower, "key") ||
	   strings.Contains(cmdLower, "secret") {
		return true
	}
	
	return false
}

// detectSensitivePatterns finds sensitive patterns in text
func (s *SecuritySanitizer) detectSensitivePatterns(text string) []string {
	var detected []string
	
	for patternName, pattern := range s.sensitivePatterns {
		if pattern.MatchString(text) {
			detected = append(detected, patternName)
		}
	}
	
	return detected
}

// redactSensitiveData replaces sensitive data with placeholders
func (s *SecuritySanitizer) redactSensitiveData(text string) string {
	sanitized := text
	
	for patternName, pattern := range s.sensitivePatterns {
		switch patternName {
		case "api_key", "bearer_token", "authorization":
			sanitized = pattern.ReplaceAllString(sanitized, "$1 [REDACTED_API_KEY]")
		case "aws_access_key":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_AWS_ACCESS_KEY]")
		case "aws_secret_key":
			sanitized = pattern.ReplaceAllString(sanitized, "$1 [REDACTED_AWS_SECRET_KEY]")
		case "ssh_private_key":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_SSH_PRIVATE_KEY]")
		case "ssh_public_key":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_SSH_PUBLIC_KEY]")
		case "db_connection":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_DB_CONNECTION]")
		case "password", "db_password":
			sanitized = pattern.ReplaceAllString(sanitized, "$1 [REDACTED_PASSWORD]")
		case "credit_card":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_CREDIT_CARD]")
		case "email":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_EMAIL]")
		case "private_ip":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_PRIVATE_IP]")
		case "url_with_creds":
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_URL_WITH_CREDENTIALS]")
		case "env_var":
			sanitized = pattern.ReplaceAllString(sanitized, "$1[REDACTED_ENV_VAR]")
		default:
			sanitized = pattern.ReplaceAllString(sanitized, "[REDACTED_SENSITIVE_DATA]")
		}
	}
	
	return sanitized
}

// calculateRiskLevel determines the risk level based on detected patterns
func (s *SecuritySanitizer) calculateRiskLevel(patterns []string) RiskLevel {
	maxRisk := RiskNone
	
	riskMap := map[string]RiskLevel{
		"api_key":         RiskHigh,
		"bearer_token":    RiskHigh,
		"authorization":   RiskHigh,
		"aws_access_key":  RiskCritical,
		"aws_secret_key":  RiskCritical,
		"ssh_private_key": RiskCritical,
		"ssh_public_key":  RiskMedium,
		"db_connection":   RiskHigh,
		"password":        RiskHigh,
		"db_password":     RiskHigh,
		"credit_card":     RiskCritical,
		"email":           RiskLow,
		"private_ip":      RiskMedium,
		"url_with_creds":  RiskHigh,
		"env_var":         RiskMedium,
		"home_path":       RiskMedium,
		"config_file":     RiskMedium,
	}
	
	for _, pattern := range patterns {
		if risk, exists := riskMap[pattern]; exists && risk > maxRisk {
			maxRisk = risk
		}
	}
	
	return maxRisk
}

// IsSafeForAI determines if data is safe to send to AI without user consent
func (s *SecuritySanitizer) IsSafeForAI(command, output string) bool {
	result := s.SanitizeData(command, output, 1)
	return result.IsSafe && !result.UserConsentRequired
}

// GetSanitizedDataForAI returns sanitized data that can be sent to AI
func (s *SecuritySanitizer) GetSanitizedDataForAI(command, output string) (string, string) {
	result := s.SanitizeData(command, output, 1)
	return result.SanitizedCommand, result.SanitizedOutput
}
