package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/freecode/freecode/internal/config"
	"github.com/freecode/freecode/internal/hook"
	"github.com/freecode/freecode/internal/session"
	"github.com/freecode/freecode/internal/tool"
)

type Engine struct {
	config   *config.Config
	tools    *tool.Registry
	hooks    *hook.Registry
	sessions *session.Manager
	agents   map[string]Agent
	mu       sync.RWMutex
}

type Agent interface {
	Name() string
	Run(ctx context.Context, req Request) (*Response, error)
}

type Request struct {
	SessionID string
	AgentName string
	Model     string
	Message   Message
	Tools     []string
	Stream    bool
}

type Response struct {
	SessionID    string
	Message      Message
	ToolCalls    []ToolCall
	Error        error
	AgentName    string
	SystemPrompt string
}

type Message struct {
	Role     string
	Content  string
	Thinking string
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]interface{}
}

func NewEngine(cfg *config.Config) *Engine {
	return &Engine{
		config:   cfg,
		tools:    tool.NewRegistry(),
		hooks:    hook.NewRegistry(),
		sessions: session.NewManager(cfg),
		agents:   make(map[string]Agent),
	}
}

func (e *Engine) RegisterAgent(agent Agent) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.agents[agent.Name()] = agent
}

func (e *Engine) GetAgent(name string) (Agent, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	agent, ok := e.agents[name]
	return agent, ok
}

func (e *Engine) Run(ctx context.Context, req Request) (*Response, error) {
	e.mu.RLock()
	agent, ok := e.agents[req.AgentName]
	e.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("agent not found: %s", req.AgentName)
	}

	return agent.Run(ctx, req)
}

func (e *Engine) RegisterTool(t tool.Tool) {
	e.tools.Register(t)
}

func (e *Engine) GetTool(name string) (tool.Tool, bool) {
	return e.tools.Get(name)
}

func (e *Engine) ListTools() []tool.Tool {
	return e.tools.List()
}

func (e *Engine) SessionManager() *session.Manager {
	return e.sessions
}

func (e *Engine) HookRegistry() *hook.Registry {
	return e.hooks
}
