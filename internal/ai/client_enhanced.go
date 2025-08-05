package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ayushsharma-1/LogAid/internal/config"
	"github.com/ayushsharma-1/LogAid/internal/logger"
)

// AIClient represents the AI service client
type AIClient struct {
	Provider string
	APIKey   string
	Model    string
	BaseURL  string
	Timeout  time.Duration
}

// NewAIClient creates a new AI client based on configuration
func NewAIClient() *AIClient {
	var provider string

	// Use config if available, otherwise fall back to environment variables
	if config.AppConfig != nil {
		provider = config.AppConfig.AIProvider
		if provider == "" {
			provider = "gemini" // default
		}
	} else {
		provider = os.Getenv("AI_PROVIDER")
		if provider == "" {
			provider = "gemini" // default
		}
	}

	timeout := 15 * time.Second
	if config.AppConfig != nil && config.AppConfig.AIRequestTimeout > 0 {
		timeout = time.Duration(config.AppConfig.AIRequestTimeout) * time.Second
	}

	client := &AIClient{
		Provider: provider,
		Timeout:  timeout,
	}

	switch provider {
	case "gemini":
		if config.AppConfig != nil {
			client.APIKey = config.AppConfig.GeminiAPIKey
			client.Model = config.AppConfig.GeminiModel
		} else {
			client.APIKey = os.Getenv("GEMINI_API_KEY")
			client.Model = os.Getenv("GEMINI_MODEL")
		}
		if client.Model == "" {
			client.Model = "gemini-2.0-flash-exp"
		}
		client.BaseURL = "https://generativelanguage.googleapis.com/v1beta/models"
	case "openai":
		if config.AppConfig != nil {
			client.APIKey = config.AppConfig.OpenAIAPIKey
			client.Model = config.AppConfig.OpenAIModel
		} else {
			client.APIKey = os.Getenv("OPENAI_API_KEY")
			client.Model = os.Getenv("OPENAI_MODEL")
		}
		if client.Model == "" {
			client.Model = "gpt-4o"
		}
		client.BaseURL = "https://api.openai.com/v1/chat/completions"
	default:
		logger.Error(fmt.Sprintf("Unsupported AI provider: %s", provider))
		return nil
	}

	if client.APIKey == "" {
		logger.Error(fmt.Sprintf("API key not found for provider: %s", provider))
		return nil
	}

	return client
}

// GetSuggestion generates a command suggestion using AI
func GetSuggestion(ctx context.Context, prompt string) (string, error) {
	client := NewAIClient()
	if client == nil {
		return "", fmt.Errorf("failed to initialize AI client")
	}

	return client.GenerateSuggestion(ctx, prompt)
}

// GenerateSuggestion generates a suggestion using the configured AI provider
func (c *AIClient) GenerateSuggestion(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	switch c.Provider {
	case "gemini":
		return c.callGemini(ctx, prompt)
	case "openai":
		return c.callOpenAI(ctx, prompt)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", c.Provider)
	}
}

// GeminiRequest represents the request structure for Gemini API
type GeminiRequest struct {
	Contents         []GeminiContent        `json:"contents"`
	GenerationConfig GeminiGenerationConfig `json:"generationConfig"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
	TopP            float64 `json:"topP"`
	TopK            int     `json:"topK"`
}

// GeminiResponse represents the response structure from Gemini API
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content      GeminiContent `json:"content"`
	FinishReason string        `json:"finishReason"`
}

// callGemini makes a request to the Gemini API
func (c *AIClient) callGemini(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s/%s:generateContent?key=%s", c.BaseURL, c.Model, c.APIKey)

	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: GeminiGenerationConfig{
			Temperature:     0.1,
			MaxOutputTokens: 500,
			TopP:            0.8,
			TopK:            10,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	suggestion := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)

	// Clean up the response to extract just the command
	suggestion = c.extractCommand(suggestion)

	logger.Debug(fmt.Sprintf("AI suggestion: %s", suggestion))
	return suggestion, nil
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIChoice struct {
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// callOpenAI makes a request to the OpenAI API
func (c *AIClient) callOpenAI(ctx context.Context, prompt string) (string, error) {
	requestBody := OpenAIRequest{
		Model: c.Model,
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: "You are a Linux command-line expert. Provide only the corrected command, no explanations.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.1,
		MaxTokens:   500,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	suggestion := strings.TrimSpace(openaiResp.Choices[0].Message.Content)

	// Clean up the response to extract just the command
	suggestion = c.extractCommand(suggestion)

	logger.Debug(fmt.Sprintf("AI suggestion: %s", suggestion))
	return suggestion, nil
}

// extractCommand extracts the actual command from AI response
func (c *AIClient) extractCommand(response string) string {
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and explanations
		if line == "" || strings.HasPrefix(line, "Explanation:") ||
			strings.HasPrefix(line, "Note:") || strings.HasPrefix(line, "The") ||
			strings.HasPrefix(line, "This") || strings.HasPrefix(line, "Here") {
			continue
		}

		// Remove markdown code block markers
		if strings.HasPrefix(line, "```") {
			continue
		}

		// Look for actual command patterns
		if strings.Contains(line, "sudo") || strings.Contains(line, "apt") ||
			strings.Contains(line, "npm") || strings.Contains(line, "git") ||
			strings.Contains(line, "docker") || strings.Contains(line, "pip") {
			return line
		}
	}

	// If no command pattern found, return the first non-empty line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "```") {
			return line
		}
	}

	return response
}
