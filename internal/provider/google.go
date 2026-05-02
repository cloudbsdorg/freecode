package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GoogleProvider struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewGoogleProvider(apiKey string) *GoogleProvider {
	return &GoogleProvider{
		APIKey:  apiKey,
		BaseURL: "https://generativelanguage.googleapis.com",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *GoogleProvider) Name() string {
	return "google"
}

func (p *GoogleProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.BaseURL + "/v1beta/models/" + req.Model + ":generateContent"

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": req.Messages[len(req.Messages)-1].Content},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":  req.Temperature,
			"maxOutputTokens": req.MaxTokens,
		},
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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Candidates) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := ""
	if len(result.Candidates[0].Content.Parts) > 0 {
		content = result.Candidates[0].Content.Parts[0].Text
	}

	return &Response{
		Content:    content,
		StopReason: "stop",
		Usage: Usage{
			InputTokens:  result.UsageMetadata.PromptTokenCount,
			OutputTokens: result.UsageMetadata.CandidatesTokenCount,
		},
	}, nil
}

func (p *GoogleProvider) ListModels(ctx context.Context) ([]Model, error) {
	url := p.BaseURL + "/v1beta/models"
	return fetchModelsFromURL(ctx, url, p.APIKey)
}
