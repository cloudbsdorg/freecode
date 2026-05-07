package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Client struct {
	mu      sync.RWMutex
	servers map[string]*ServerConn
}

type ServerConn struct {
	ID      string
	Name    string
	URL     string
	Auth    string
	Headers map[string]string
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type CallRequest struct {
	Method string `json:"method"`
	Params *struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments,omitempty"`
	} `json:"params,omitempty"`
}

type CallResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	IsError bool `json:"isError,omitempty"`
}

func NewClient() *Client {
	return &Client{
		servers: make(map[string]*ServerConn),
	}
}

func (c *Client) AddServer(id, name, url, auth string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.servers[id] = &ServerConn{
		ID:      id,
		Name:    name,
		URL:     url,
		Auth:    auth,
		Headers: make(map[string]string),
	}
}

func (c *Client) RemoveServer(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.servers, id)
}

func (c *Client) ListTools() ([]Tool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var allTools []Tool
	for _, srv := range c.servers {
		tools, err := c.listServerTools(srv)
		if err != nil {
			continue
		}
		allTools = append(allTools, tools...)
	}
	return allTools, nil
}

func (c *Client) listServerTools(srv *ServerConn) ([]Tool, error) {
	return []Tool{}, nil
}

func (c *Client) CallTool(ctx context.Context, serverID, toolName string, args map[string]interface{}) (*CallResponse, error) {
	c.mu.RLock()
	srv, ok := c.servers[serverID]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("server not found: %s", serverID)
	}

	reqBody := CallRequest{
		Method: "tools/call",
		Params: &struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}{
			Name:      toolName,
			Arguments: args,
		},
	}

	reqBodyJSON, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", srv.URL, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if srv.Auth != "" {
		req.Header.Set("Authorization", "Bearer "+srv.Auth)
	}

	return &CallResponse{}, nil
}
