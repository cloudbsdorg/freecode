package tool

import (
	"context"
	"sync"

	"github.com/freecode/freecode/internal/config"
)

type Response struct {
	Result string
	Error  error
}

type Request struct {
	Name      string
	Arguments map[string]interface{}
	SessionID string
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

type ToolFactory func() Tool

var (
	mu      sync.RWMutex
	tools   = make(map[string]ToolFactory)
	enabled = make(map[string]bool)
)

func Register(name string, factory ToolFactory) {
	mu.Lock()
	defer mu.Unlock()
	tools[name] = factory
	enabled[name] = true
}

func IsEnabled(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return enabled[name]
}

func SetEnabled(name string, en bool) {
	mu.Lock()
	defer mu.Unlock()
	enabled[name] = en
}

func SetEnabledFromConfig(states map[string]bool) {
	mu.Lock()
	defer mu.Unlock()
	for name, state := range states {
		enabled[name] = state
	}
}

func ListTools() []string {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]string, 0, len(tools))
	for name := range tools {
		list = append(list, name)
	}
	return list
}

func GetFactory(name string) (ToolFactory, bool) {
	mu.RLock()
	defer mu.RUnlock()
	f, ok := tools[name]
	return f, ok
}

type Registry struct {
	mu       sync.RWMutex
	tools    map[string]Tool
	disabled map[string]bool
}

func NewRegistry() *Registry {
	return &Registry{
		tools:    make(map[string]Tool),
		disabled: make(map[string]bool),
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
	if r.disabled[name] {
		return nil, false
	}
	t, ok := r.tools[name]
	return t, ok
}

func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]Tool, 0, len(r.tools))
	for name, t := range r.tools {
		if !r.disabled[name] {
			list = append(list, t)
		}
	}
	return list
}

func (r *Registry) Execute(ctx context.Context, req Request) (*Response, error) {
	t, ok := r.Get(req.Name)
	if !ok {
		return nil, &ToolNotFoundError{Name: req.Name}
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

func (r *Registry) Disable(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.disabled[name] = true
}

func (r *Registry) Enable(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.disabled, name)
}

func (r *Registry) IsDisabled(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.disabled[name]
}

type ToolNotFoundError struct {
	Name string
}

func (e *ToolNotFoundError) Error() string {
	return "tool not found: " + e.Name
}

func AllTools() []Tool {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]Tool, 0, len(tools))
	for _, factory := range tools {
		list = append(list, factory())
	}
	return list
}

func CreateAll(cfg *config.Config) []*ToolWithConfig {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]*ToolWithConfig, 0, len(tools))
	for name, factory := range tools {
		t := factory()
		list = append(list, &ToolWithConfig{
			Name:    name,
			Tool:    t,
			Enabled: enabled[name],
		})
	}
	return list
}

type ToolWithConfig struct {
	Name        string
	Tool        Tool
	Enabled     bool
	Description string
}

type Engine interface {
	RegisterTool(t Tool)
}

func ListToolsWithStatus() []ToolInfo {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]ToolInfo, 0, len(tools))
	for name, factory := range tools {
		t := factory()
		list = append(list, ToolInfo{
			Name:        name,
			Description: t.Description(),
			Enabled:     enabled[name],
		})
	}
	return list
}

type ToolInfo struct {
	Name        string
	Description string
	Enabled     bool
}