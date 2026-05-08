package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	headers    map[string]string
}

type ClientOption func(*Client)

func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = key
	}
}

func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}

func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + path
	if path != "" && !strings.HasPrefix(path, "/") {
		url = c.baseURL + "/" + path
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, method, path string, body any) ([]byte, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	return c.Do(ctx, "GET", path, nil)
}

func (c *Client) Post(ctx context.Context, path string, body any) ([]byte, error) {
	return c.Do(ctx, "POST", path, body)
}

func (c *Client) Put(ctx context.Context, path string, body any) ([]byte, error) {
	return c.Do(ctx, "PUT", path, body)
}

func (c *Client) Delete(ctx context.Context, path string) ([]byte, error) {
	return c.Do(ctx, "DELETE", path, nil)
}

func (c *Client) GetJSON(ctx context.Context, path string, result any) error {
	data, err := c.Get(ctx, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func (c *Client) PostJSON(ctx context.Context, path string, body, result any) error {
	data, err := c.Post(ctx, path, body)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}

func (c *Client) PutJSON(ctx context.Context, path string, body, result any) error {
	data, err := c.Put(ctx, path, body)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}

func (c *Client) DeleteJSON(ctx context.Context, path string, result any) error {
	data, err := c.Delete(ctx, path)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}

type APIError struct {
	StatusCode int
	Message    string
	Details    map[string]any
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error: status %d, message: %s", e.StatusCode, e.Message)
}

func ParseAPIError(statusCode int, body []byte) *APIError {
	var errResp struct {
		Error   string         `json:"error"`
		Message string         `json:"message"`
		Details map[string]any `json:"details"`
	}
	if json.Unmarshal(body, &errResp) == nil {
		msg := errResp.Message
		if msg == "" {
			msg = errResp.Error
		}
		return &APIError{
			StatusCode: statusCode,
			Message:    msg,
			Details:    errResp.Details,
		}
	}
	return &APIError{
		StatusCode: statusCode,
		Message:    string(body),
	}
}

type PageResponse[T any] struct {
	Data       []T
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

func (c *Client) GetPage(ctx context.Context, path string, page, perPage int, result any) error {
	url := fmt.Sprintf("%s?page=%d&per_page=%d", path, page, perPage)
	data, err := c.Get(ctx, url)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}