package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type VercelProvider struct {
	Token   string
	BaseURL string
	Client  *http.Client
}

func NewVercelProvider(token string) *VercelProvider {
	return &VercelProvider{
		Token:   token,
		BaseURL: "https://api.vercel.com",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *VercelProvider) Name() string {
	return "vercel"
}

func (p *VercelProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.BaseURL + "/ai/v1/chat/completions"

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
	httpReq.Header.Set("Authorization", "Bearer "+p.Token)

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &Response{
		Content:    result.Choices[0].Message.Content,
		StopReason: result.Choices[0].FinishReason,
		Usage: Usage{
			InputTokens:  result.Usage.PromptTokens,
			OutputTokens: result.Usage.CompletionTokens,
		},
	}, nil
}

func (p *VercelProvider) ListModels(ctx context.Context) ([]Model, error) {
	return fetchModelsFromURL(ctx, p.BaseURL+"/ai/v1/models", p.Token)
}
