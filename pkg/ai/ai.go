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

	"github.com/ayushsharma-1/LogAid/pkg/plugin"
)

// Provider defines the interface for AI providers
type Provider interface {
	Name() string
	Suggest(ctx context.Context, req *plugin.SuggestionRequest) (*plugin.SuggestionResponse, error)
}

// Client manages AI providers with fallback support
type Client struct {
	providers []Provider
	timeout   time.Duration
}

// NewClient creates a new AI client with available providers
func NewClient() (*Client, error) {
	client := &Client{
		providers: make([]Provider, 0),
		timeout:   30 * time.Second,
	}

	// Initialize providers based on available API keys
	if geminiKey := os.Getenv("GEMINI_API_KEY"); geminiKey != "" {
		provider := NewGeminiProvider(geminiKey, "gemini-pro", client.timeout)
		client.providers = append(client.providers, provider)
	}

	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
		provider := NewOpenAIProvider(openaiKey, "gpt-3.5-turbo", client.timeout)
		client.providers = append(client.providers, provider)
	}

	if len(client.providers) == 0 {
		return nil, fmt.Errorf("no AI providers configured - please set GEMINI_API_KEY or OPENAI_API_KEY")
	}

	return client, nil
}

// Suggest gets a suggestion from the first available AI provider with fallback
func (c *Client) Suggest(ctx context.Context, req *plugin.SuggestionRequest) (*plugin.SuggestionResponse, error) {
	var lastErr error

	for i, provider := range c.providers {
		resp, err := provider.Suggest(ctx, req)
		if err == nil && resp.SuggestedCommand != "" {
			return resp, nil
		}

		lastErr = err
		if i < len(c.providers)-1 {
			// Log fallback attempt but continue to next provider
			continue
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all AI providers failed: %w", lastErr)
	}

	return nil, fmt.Errorf("no AI providers returned valid suggestions")
}

// GetActiveProviders returns the list of currently active providers
func (c *Client) GetActiveProviders() []string {
	names := make([]string, len(c.providers))
	for i, provider := range c.providers {
		names[i] = provider.Name()
	}
	return names
}

// GeminiProvider implements the Provider interface for Google Gemini
type GeminiProvider struct {
	apiKey   string
	model    string
	client   *http.Client
	endpoint string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey, model string, timeout time.Duration) *GeminiProvider {
	return &GeminiProvider{
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{Timeout: timeout},
		endpoint: "https://generativelanguage.googleapis.com/v1beta/models/" + model + ":generateContent",
	}
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "gemini"
}

// Suggest generates a suggestion using Google Gemini
func (p *GeminiProvider) Suggest(ctx context.Context, req *plugin.SuggestionRequest) (*plugin.SuggestionResponse, error) {
	prompt := fmt.Sprintf(`You are a Linux command-line expert. A user ran this command and got an error:

Command: %s
Output: %s
Exit Code: %d

Please provide a corrected command and explanation. Respond ONLY with valid JSON in this format:
{
  "suggested_command": "corrected command here",
  "explanation": "brief explanation of what was wrong and how the correction fixes it"
}

Do not include any text outside the JSON response.`, req.Command, req.Output, req.ExitCode)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":   0.3,
			"maxOutputTokens": 1024,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := p.endpoint + "?key=" + p.apiKey
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
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
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text
	return parseAIResponse(text)
}

// OpenAIProvider implements the Provider interface for OpenAI GPT
type OpenAIProvider struct {
	apiKey   string
	model    string
	client   *http.Client
	endpoint string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey, model string, timeout time.Duration) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:   apiKey,
		model:    model,
		client:   &http.Client{Timeout: timeout},
		endpoint: "https://api.openai.com/v1/chat/completions",
	}
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Suggest generates a suggestion using OpenAI GPT
func (p *OpenAIProvider) Suggest(ctx context.Context, req *plugin.SuggestionRequest) (*plugin.SuggestionResponse, error) {
	prompt := fmt.Sprintf(`You are a Linux command-line expert. A user ran this command and got an error:

Command: %s
Output: %s
Exit Code: %d

Please provide a corrected command and explanation. Respond ONLY with valid JSON in this format:
{
  "suggested_command": "corrected command here",
  "explanation": "brief explanation of what was wrong and how the correction fixes it"
}

Do not include any text outside the JSON response.`, req.Command, req.Output, req.ExitCode)

	reqBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  1024,
		"temperature": 0.3,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	text := openaiResp.Choices[0].Message.Content
	return parseAIResponse(text)
}

// parseAIResponse parses the AI response text and extracts JSON
func parseAIResponse(text string) (*plugin.SuggestionResponse, error) {
	// Extract JSON from the response
	jsonStr := extractJSON(text)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response: %s", text)
	}

	var response plugin.SuggestionResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if response.SuggestedCommand == "" {
		return nil, fmt.Errorf("empty suggested command in response")
	}

	return &response, nil
}

// extractJSON extracts JSON content from text that might contain extra content
func extractJSON(text string) string {
	text = strings.TrimSpace(text)

	// Look for JSON block markers
	if strings.Contains(text, "```json") {
		start := strings.Index(text, "```json") + 7
		end := strings.Index(text[start:], "```")
		if end != -1 {
			return strings.TrimSpace(text[start : start+end])
		}
	}

	// Look for JSON object brackets
	start := strings.Index(text, "{")
	if start == -1 {
		return ""
	}

	// Find the matching closing brace
	braceCount := 0
	for i := start; i < len(text); i++ {
		if text[i] == '{' {
			braceCount++
		} else if text[i] == '}' {
			braceCount--
			if braceCount == 0 {
				return text[start : i+1]
			}
		}
	}

	return ""
}
