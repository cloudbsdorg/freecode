package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OllamaCloudProvider struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewOllamaCloudProvider(apiKey string) *OllamaCloudProvider {
	return &OllamaCloudProvider{
		APIKey:  apiKey,
		BaseURL: "https://cloud.ollama.ai",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *OllamaCloudProvider) Name() string {
	return "ollamacloud"
}

func (p *OllamaCloudProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.BaseURL + "/api/chat"

	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Done bool `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Response{
		Content:    result.Message.Content,
		StopReason: "stop",
		Usage: Usage{
			InputTokens:  0,
			OutputTokens: 0,
		},
	}, nil
}

func (p *OllamaCloudProvider) ListModels(ctx context.Context) ([]Model, error) {
	return nil, nil
}
