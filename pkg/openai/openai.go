package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	chatAPIURL     = "https://api.openai.com/v1/chat/completions"
	defaultTimeout = 30 * time.Second
)

// Config holds the OpenAI configuration
type Config struct {
	APIKey    string
	Model     string
	ShellType string
	Enabled   bool
}

// Client represents an OpenAI API client
type Client struct {
	config     Config
	httpClient *http.Client
}

// Response represents the OpenAI API response for chat completions
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewClient creates a new OpenAI client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// GenerateCompletion generates a completion using OpenAI Chat API
func (c *Client) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	if !c.config.Enabled {
		return "", fmt.Errorf("OpenAI is not enabled")
	}

	// Check if this is a workflow generation request (longer, more complex)
	isWorkflowRequest := strings.Contains(strings.ToLower(prompt), "workflow") ||
		strings.Contains(strings.ToLower(prompt), "yaml") ||
		len(prompt) > 200

	var messages []map[string]string
	var maxTokens int
	var temperature float64

	if isWorkflowRequest {
		// For workflow generation: use system prompt and more tokens
		messages = []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert in creating runfromyaml workflow configurations. Generate complete, production-ready YAML workflows. Only return valid YAML, no explanations.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		}
		maxTokens = 2000
		temperature = 0.3
	} else {
		// For shell commands: use original format
		messages = []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf("%s. show a %s example. Please do not write explanations. Please just a suggestion as %s code.", prompt, c.config.ShellType, c.config.ShellType),
			},
		}
		maxTokens = 100
		temperature = 0
	}

	reqBody := map[string]interface{}{
		"model":       c.config.Model,
		"messages":    messages,
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"top_p":       1.0,
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", chatAPIURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return strings.TrimSpace(strings.ReplaceAll(response.Choices[0].Message.Content, "`", "")), nil
}

// Legacy support for backward compatibility
var (
	IsAiEnabled bool
	Key         string
	Model       string
	ShellType   string
)

// OpenAI is a legacy function that uses the new client internally
func OpenAI(apiKey string, model string, prompt string, cmdtype string) (map[string][]interface{}, error) {
	client := NewClient(Config{
		APIKey:    apiKey,
		Model:     model,
		ShellType: cmdtype,
		Enabled:   true,
	})

	response, err := client.GenerateCompletion(context.Background(), prompt)
	if err != nil {
		return nil, err
	}

	// Convert to legacy format
	return map[string][]interface{}{
		"choices": {map[string]interface{}{
			"text": response,
		}},
	}, nil
}

// PrintAiResponse is kept for backward compatibility
func PrintAiResponse(response map[string][]interface{}) string {
	if len(response["choices"]) == 0 {
		return ""
	}

	if choice, ok := response["choices"][0].(map[string]interface{}); ok {
		if text, ok := choice["text"].(string); ok {
			return strings.TrimSpace(strings.ReplaceAll(text, "`", ""))
		}
	}
	return ""
}
