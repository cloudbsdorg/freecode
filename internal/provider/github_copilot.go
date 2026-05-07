package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GitHubCopilotProvider struct {
	Token   string
	BaseURL string
	Client  *http.Client
}

func NewGitHubCopilotProvider(token string) *GitHubCopilotProvider {
	return &GitHubCopilotProvider{
		Token:   token,
		BaseURL: "https://api.githubcopilot.com",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *GitHubCopilotProvider) Name() string {
	return "github_copilot"
}

func (p *GitHubCopilotProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.BaseURL + "/chat/completions"

	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.Token)
	httpReq.Header.Set("Editor-Version", "freecode/1.0")

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

func (p *GitHubCopilotProvider) ListModels(ctx context.Context) ([]Model, error) {
	return []Model{
		{ID: "gpt-4o", Name: "GPT-4o", Provider: "github_copilot"},
		{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Provider: "github_copilot"},
		{ID: "claude-3.5-sonnet", Name: "Claude 3.5 Sonnet", Provider: "github_copilot"},
		{ID: "claude-3.7-sonnet", Name: "Claude 3.7 Sonnet", Provider: "github_copilot"},
	}, nil
}
