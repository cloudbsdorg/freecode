package tool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WebFetchTool struct{}

func NewWebFetchTool() *WebFetchTool {
	return &WebFetchTool{}
}

func (t *WebFetchTool) Name() string {
	return "webfetch"
}

func (t *WebFetchTool) Description() string {
	return "Fetch content from a URL"
}

func (t *WebFetchTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "webfetch",
		Description: "Fetch content from a URL",
		Parameters: map[string]Parameter{
			"url": {
				Type:        "string",
				Description: "URL to fetch",
				Required:    true,
			},
			"timeout": {
				Type:        "integer",
				Description: "Request timeout in seconds",
				Default:     30,
			},
		},
	}
}

func (t *WebFetchTool) Execute(ctx context.Context, req Request) (*Response, error) {
	urlStr, ok := req.Arguments["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url must be a string")
	}

	timeoutSecs := 30
	if to, ok := req.Arguments["timeout"].(int); ok {
		timeoutSecs = to
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: time.Duration(timeoutSecs) * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &Response{
		Result: string(body),
	}, nil
}
