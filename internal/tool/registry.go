package tool

import (
	"context"
	"fmt"
	"sync"
)

type Request struct {
	Name      string
	Arguments map[string]interface{}
	SessionID string
}

type Response struct {
	Result string
	Error  error
}

type Tool interface {
	Name() string
	Description() string
	Schema() ToolSchema
	Execute(ctx context.Context, req Request) (*Response, error)
}

type ToolSchema struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Parameters  map[string]Parameter `json:"parameters"`
}

type Parameter struct {
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Default     any        `json:"default,omitempty"`
	Required    bool       `json:"required"`
	Enum        []string   `json:"enum,omitempty"`
	Items       *Parameter `json:"items,omitempty"`
}

type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(t Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[t.Name()] = t
}

func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tools := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

func (r *Registry) Execute(ctx context.Context, req Request) (*Response, error) {
	t, ok := r.Get(req.Name)
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", req.Name)
	}
	return t.Execute(ctx, req)
}

func (r *Registry) Schema() []ToolSchema {
	tools := r.List()
	schemas := make([]ToolSchema, len(tools))
	for i, t := range tools {
		schemas[i] = t.Schema()
	}
	return schemas
}
