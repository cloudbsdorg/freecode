package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BedrockProvider struct {
	Region      string
	Profile     string
	Endpoint    string
	AccessKey   string
	SecretKey   string
	Client      *http.Client
}

func NewBedrockProvider(region, profile, endpoint, accessKey, secretKey string) *BedrockProvider {
	return &BedrockProvider{
		Region:    region,
		Profile:   profile,
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *BedrockProvider) Name() string {
	return "bedrock"
}

func (p *BedrockProvider) GetEndpoint() string {
	if p.Endpoint != "" {
		return p.Endpoint
	}
	return "https://bedrock." + p.Region + ".amazonaws.com"
}

func (p *BedrockProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	url := p.GetEndpoint() + "/model/" + req.Model + "/invoke"

	payload := map[string]interface{}{
		"inputText": req.Messages[len(req.Messages)-1].Content,
		"text":      req.Model,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if p.AccessKey != "" && p.SecretKey != "" {
		httpReq.Header.Set("X-Access-Key", p.AccessKey)
		httpReq.Header.Set("X-Secret-Key", p.SecretKey)
	}

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Body struct {
			Text string `json:"text"`
		} `json:"body"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Response{
		Content:    result.Body.Text,
		StopReason: "stop",
		Usage: Usage{
			InputTokens:  0,
			OutputTokens: 0,
		},
	}, nil
}

func (p *BedrockProvider) ListModels(ctx context.Context) ([]Model, error) {
	return nil, nil
}
