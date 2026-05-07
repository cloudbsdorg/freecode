package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GitLabProvider struct {
	Token   string
	BaseURL string
	Client  *http.Client
}

func NewGitLabProvider(token string) *GitLabProvider {
	return &GitLabProvider{
		Token:   token,
		BaseURL: "https://gitlab.com",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *GitLabProvider) Name() string {
	return "gitlab"
}

func (p *GitLabProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.BaseURL + "/api/v4/ai/duo_chat"

	messages := make([]map[string]string, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = map[string]string{
			"role":    m.Role,
			"content": m.Content,
		}
	}

	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": messages,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("PRIVATE-TOKEN", p.Token)

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

func (p *GitLabProvider) ListModels(ctx context.Context) ([]Model, error) {
	return []Model{
		{ID: "duo-chat-haiku-4-5", Name: "GitLab Duo Chat Haiku", Provider: "gitlab"},
		{ID: "duo-chat-sonnet-4-5", Name: "GitLab Duo Chat Sonnet", Provider: "gitlab"},
		{ID: "duo-chat-opus-4-5", Name: "GitLab Duo Chat Opus", Provider: "gitlab"},
	}, nil
}
