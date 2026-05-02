package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OllamaProvider struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func NewOllamaProvider(baseURL, apiKey string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaProvider{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaRequest struct {
	Model    string           `json:"model"`
	Messages []ollamaMessage  `json:"messages"`
	Stream   bool             `json:"stream"`
}

type ollamaResponse struct {
	Message struct {
		Content string `json:"content"`
		Role    string `json:"role"`
	} `json:"message"`
	Done           bool `json:"done"`
	TotalDuration int64 `json:"total_duration"`
	EvalCount     int   `json:"eval_count"`
	PromptEvalCount int `json:"prompt_eval_count"`
}

type ollamaTagsResponse struct {
	Models []struct {
		Name       string `json:"name"`
		Model      string `json:"model"`
		Size       int64  `json:"size"`
		Digest     string `json:"digest"`
		ModifiedAt string `json:"modified_at"`
	} `json:"models"`
}

func (p *OllamaProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.BaseURL + "/api/chat"

	messages := make([]ollamaMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		role := m.Role
		if role == "assistant" {
			role = "assistant"
		} else if role == "system" {
			role = "system"
		} else {
			role = "user"
		}
		messages = append(messages, ollamaMessage{
			Role:    role,
			Content: m.Content,
		})
	}

	payload := ollamaRequest{
		Model:    req.Model,
		Messages: messages,
		Stream:   false,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if p.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
	}

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Response{
		Content:    result.Message.Content,
		StopReason: "stop",
		Usage: Usage{
			InputTokens:  result.PromptEvalCount,
			OutputTokens: result.EvalCount,
		},
	}, nil
}

func (p *OllamaProvider) ListModels(ctx context.Context) ([]Model, error) {
	url := p.BaseURL + "/api/tags"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if p.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.APIKey)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]Model, 0, len(result.Models))
	for _, m := range result.Models {
		models = append(models, Model{
			ID:       m.Name,
			Name:     m.Name,
			Provider: "ollama",
			Created:  0,
			OwnedBy:  "ollama",
			Capabilities: ModelCapabilities{
				Temperature: true,
				Reasoning:   false,
				ToolCall:    false,
				Vision:      false,
				Audio:      false,
			},
			Cost: ModelCost{Input: 0, Output: 0},
			Limit: ModelLimit{
				Context: 8192,
				Input:   0,
				Output:  0,
			},
		})
	}
	return models, nil
}