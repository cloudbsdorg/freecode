package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BuiltinMCP struct {
	httpClient *http.Client
}

func NewBuiltinMCP() *BuiltinMCP {
	return &BuiltinMCP{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type ExaSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type exaSearchResponse struct {
	Results []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Text    string `json:"text"`
		Snippet string `json:"snippet"`
	} `json:"results"`
}

// ExaSearch performs a web search using the Exa API.
// API key can be set via EXA_API_KEY environment variable.
func (m *BuiltinMCP) ExaSearch(query string, numResults int) ([]ExaSearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}
	if numResults <= 0 {
		numResults = 10
	}
	if numResults > 100 {
		numResults = 100
	}

	apiKey := os.Getenv("EXA_API_KEY")
	if apiKey == "" {
		return []ExaSearchResult{
			{Title: "Exa Search Not Configured", URL: "https://exa.ai", Snippet: "Set EXA_API_KEY environment variable to enable web search. Sign up at https://exa.ai"},
		}, nil
	}

	apiURL := "https://api.exa.ai/search"

	body, err := json.Marshal(map[string]interface{}{
		"query":      query,
		"numResults": numResults,
		"text":       true,
		"highlight":  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "freecode/1.0")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("exa API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var exaResp exaSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&exaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	results := make([]ExaSearchResult, 0, len(exaResp.Results))
	for _, r := range exaResp.Results {
		snippet := r.Snippet
		if snippet == "" {
			snippet = r.Text
		}
		if len(snippet) > 500 {
			snippet = snippet[:500] + "..."
		}
		results = append(results, ExaSearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: snippet,
		})
	}

	return results, nil
}

type context7SearchResponse struct {
	Results []struct {
		Content string `json:"content"`
	} `json:"results"`
}

// Context7Docs searches for library documentation using Context7 API.
// API key can be set via CONTEXT7_API_KEY environment variable.
func (m *BuiltinMCP) Context7Docs(query string) (string, error) {
	if query == "" {
		return "", fmt.Errorf("query cannot be empty")
	}

	apiKey := os.Getenv("CONTEXT7_API_KEY")
	if apiKey == "" {
		return "Context7 Docs Not Configured\n\nSet CONTEXT7_API_KEY environment variable to enable documentation search.\n\nTip: As a workaround, you can use GrepApp (grep_app) to search for code examples in public GitHub repositories.", nil
	}

	apiURL := "https://api.context7.io/v1/search"

	body, err := json.Marshal(map[string]interface{}{
		"query": query,
		"limit": 5,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "freecode/1.0")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("context7 API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var ctx7Resp context7SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&ctx7Resp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(ctx7Resp.Results) == 0 {
		return fmt.Sprintf("No documentation found for: %s\n\nTip: Try a more specific library name, e.g., 'langchain openai' instead of just 'langchain'", query), nil
	}

	var content string
	for i, r := range ctx7Resp.Results {
		if i > 0 {
			content += "\n\n---\n\n"
		}
		content += r.Content
	}

	return content, nil
}

type grepAppSearchResponse struct {
	Hits struct {
		Hits []struct {
			Repo struct {
				Raw string `json:"raw"`
			} `json:"repo"`
			Path struct {
				Raw string `json:"raw"`
			} `json:"path"`
			Content struct {
				Snippet string `json:"snippet"`
			} `json:"content"`
		} `json:"hits"`
	} `json:"hits"`
	Facets struct {
		Count int `json:"count"`
	} `json:"facets"`
}

// GrepApp searches for code across public GitHub repositories using grep.app API.
func (m *BuiltinMCP) GrepApp(query string) ([]string, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	apiURL := "https://grep.app/api/search"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Set("q", query)
	q.Set("limit", "10")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", "freecode/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("grep.app API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var grepResp grepAppSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&grepResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(grepResp.Hits.Hits) == 0 {
		return []string{fmt.Sprintf("No code results found for: %s", query)}, nil
	}

	results := make([]string, 0, len(grepResp.Hits.Hits))
	for _, hit := range grepResp.Hits.Hits {
		url := fmt.Sprintf("https://github.com/%s/blob/main/%s", hit.Repo.Raw, hit.Path.Raw)
		results = append(results, url)
	}

	return results, nil
}

// ExecuteBuiltinTool runs a built-in tool by name with the given arguments.
// This is used by the MCP server to route tool calls to built-in implementations.
func (m *BuiltinMCP) ExecuteBuiltinTool(ctx context.Context, name string, args map[string]interface{}) (interface{}, error) {
	switch name {
	case "exa_search":
		query, ok := args["query"].(string)
		if !ok {
			return nil, fmt.Errorf("query must be a string")
		}
		numResults := 10
		if n, ok := args["num_results"].(int); ok {
			numResults = n
		}
		return m.ExaSearch(query, numResults)

	case "context7_docs":
		query, ok := args["query"].(string)
		if !ok {
			return nil, fmt.Errorf("query must be a string")
		}
		return m.Context7Docs(query)

	case "grep_app":
		query, ok := args["query"].(string)
		if !ok {
			return nil, fmt.Errorf("query must be a string")
		}
		return m.GrepApp(query)

	default:
		return nil, fmt.Errorf("unknown built-in tool: %s", name)
	}
}
