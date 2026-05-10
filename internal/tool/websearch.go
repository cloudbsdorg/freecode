package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WebSearchTool struct{}

func init() {
	Register("websearch", func() Tool { return &WebSearchTool{} })
}

func NewWebSearchTool() *WebSearchTool {
	return &WebSearchTool{}
}

func (t *WebSearchTool) Name() string {
	return "websearch"
}

func (t *WebSearchTool) Description() string {
	return "Search the web using Exa"
}

func (t *WebSearchTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "websearch",
		Description: "Search the web using Exa",
		Parameters: map[string]Parameter{
			"query": {
				Type:        "string",
				Description: "Search query",
				Required:    true,
			},
			"num_results": {
				Type:        "integer",
				Description: "Number of results to return",
				Default:     10,
			},
		},
	}
}

func (t *WebSearchTool) Execute(ctx context.Context, req Request) (*Response, error) {
	query, ok := req.Arguments["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query must be a string")
	}

	numResults := 10
	if n, ok := req.Arguments["num_results"].(int); ok {
		numResults = n
	}

	apiKey := ""
	apiURL := "https://api.exa.ai/search"

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	}

	body, _ := json.Marshal(map[string]interface{}{
		"query":      query,
		"numResults": numResults,
	})
	httpReq.Body = io.NopCloser(bytes.NewReader(body))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &Response{
			Result: fmt.Sprintf("Web search: query=%s num_results=%d (API not configured)", query, numResults),
		}, nil
	}

	return &Response{
		Result: fmt.Sprintf("Web search: query=%s num_results=%d", query, numResults),
	}, nil
}
