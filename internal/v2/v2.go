package v2

import (
	"context"
	"net/http"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *Client) Do(ctx context.Context, method, path string, body any) ([]byte, error) {
	return []byte{}, nil
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
