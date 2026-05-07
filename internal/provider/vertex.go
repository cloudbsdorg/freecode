package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type VertexProvider struct {
	ProjectID string
	Location  string
	Token     string
	BaseURL   string
	Client    *http.Client
}

func NewVertexProvider(projectID, location, token string) *VertexProvider {
	baseURL := "https://" + location + "-aiplatform.googleapis.com/v1"
	if projectID != "" {
		baseURL = "https://" + location + "-aiplatform.googleapis.com/v1/projects/" + projectID + "/locations/" + location
	}
	return &VertexProvider{
		ProjectID: projectID,
		Location:  location,
		Token:     token,
		BaseURL:   baseURL,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *VertexProvider) Name() string {
	return "vertex"
}

func (p *VertexProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.BaseURL + "/publishers/google/models/" + req.Model + ":predict"

	payload := map[string]interface{}{
		"instances": []map[string]interface{}{
			{
				"messages": req.Messages,
			},
		},
		"parameters": map[string]interface{}{
			"temperature": req.Temperature,
			"maxTokens":   req.MaxTokens,
		},
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
		Predictions []struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		} `json:"predictions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Predictions) == 0 {
		return nil, fmt.Errorf("no predictions in response")
	}

	content := ""
	if len(result.Predictions[0].Messages) > 0 {
		content = result.Predictions[0].Messages[0].Content
	}

	return &Response{
		Content:    content,
		StopReason: "stop",
		Usage: Usage{
			InputTokens:  0,
			OutputTokens: 0,
		},
	}, nil
}

func (p *VertexProvider) ListModels(ctx context.Context) ([]Model, error) {
	url := p.BaseURL + "/models"
	return fetchModelsFromURL(ctx, url, p.Token)
}
