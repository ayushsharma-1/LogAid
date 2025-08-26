package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ayushsharma-1/LogAid/internal/config"
)

// Client represents an AI client interface
type Client interface {
	AnalyzeError(ctx context.Context, command, stderr string, exitCode int) (*Suggestion, error)
}

// Suggestion represents an AI-generated suggestion
type Suggestion struct {
	Command     string  `json:"command"`
	Explanation string  `json:"explanation"`
	Confidence  float64 `json:"confidence"`
}

// geminiClient implements the Client interface for Google's Gemini API
type geminiClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// openaiClient implements the Client interface for OpenAI's API
type openaiClient struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewClient creates a new AI client based on the configuration
func NewClient(cfg *config.Config) (Client, error) {
	switch strings.ToLower(cfg.AIProvider) {
	case "gemini":
		return newGeminiClient(cfg)
	case "openai":
		return newOpenAIClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.AIProvider)
	}
}

// newGeminiClient creates a new Gemini client
func newGeminiClient(cfg *config.Config) (*geminiClient, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	return &geminiClient{
		apiKey:  cfg.APIKey,
		model:   cfg.AIModel,
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// newOpenAIClient creates a new OpenAI client
func newOpenAIClient(cfg *config.Config) (*openaiClient, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	model := cfg.AIModel
	if model == "" || strings.Contains(model, "gemini") {
		model = "gpt-3.5-turbo" // Default OpenAI model
	}

	return &openaiClient{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: "https://api.openai.com/v1",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// AnalyzeError analyzes a command error using Gemini API
func (c *geminiClient) AnalyzeError(ctx context.Context, command, stderr string, exitCode int) (*Suggestion, error) {
	prompt := buildErrorAnalysisPrompt(command, stderr, exitCode)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.3,
			"maxOutputTokens": 1000,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini API")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	return parseAISuggestion(responseText)
}

// AnalyzeError analyzes a command error using OpenAI API
func (c *openaiClient) AnalyzeError(ctx context.Context, command, stderr string, exitCode int) (*Suggestion, error) {
	prompt := buildErrorAnalysisPrompt(command, stderr, exitCode)

	reqBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "You are LogAid, an expert CLI assistant that helps developers fix command-line errors quickly and efficiently.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature":     0.3,
		"max_tokens":      1000,
		"response_format": map[string]string{"type": "text"},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI API")
	}

	responseText := openaiResp.Choices[0].Message.Content
	return parseAISuggestion(responseText)
}

// buildErrorAnalysisPrompt creates a comprehensive prompt for error analysis
func buildErrorAnalysisPrompt(command, stderr string, exitCode int) string {
	prompt := fmt.Sprintf(`You are LogAid, an expert CLI assistant. A user ran a command that failed. Please analyze the error and provide a concise, actionable suggestion.

COMMAND: %s
EXIT CODE: %d
ERROR OUTPUT: %s

Please respond in the following JSON format:
{
  "command": "the exact command to fix the issue",
  "explanation": "a brief explanation of what went wrong and how the fix works",
  "confidence": 0.95
}

Focus on:
1. Providing the exact command to run (if applicable)
2. Being concise but informative
3. Considering common issues like permissions, missing packages, typos, etc.
4. If no specific fix command is needed, set "command" to an empty string
5. Set confidence between 0.0 and 1.0 based on how certain you are about the fix

Respond only with valid JSON.`, command, exitCode, stderr)

	return prompt
}

// parseAISuggestion parses the AI response into a Suggestion struct
func parseAISuggestion(responseText string) (*Suggestion, error) {
	// Clean up the response text (remove markdown code blocks if present)
	responseText = strings.TrimSpace(responseText)
	if strings.HasPrefix(responseText, "```json") {
		responseText = strings.TrimPrefix(responseText, "```json")
	}
	if strings.HasPrefix(responseText, "```") {
		responseText = strings.TrimPrefix(responseText, "```")
	}
	if strings.HasSuffix(responseText, "```") {
		responseText = strings.TrimSuffix(responseText, "```")
	}
	responseText = strings.TrimSpace(responseText)

	var suggestion Suggestion
	if err := json.Unmarshal([]byte(responseText), &suggestion); err != nil {
		// If JSON parsing fails, create a fallback suggestion
		return &Suggestion{
			Command:     "",
			Explanation: responseText,
			Confidence:  0.5,
		}, nil
	}

	return &suggestion, nil
}
