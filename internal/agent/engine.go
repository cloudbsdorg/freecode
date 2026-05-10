package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	Parts    []MessagePart
}

type MessagePart struct {
	Type    string
	Content string
	Tool    string
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]interface{}
}

func NewEngine(cfg *config.Config) *Engine {
	eng := &Engine{
		config:   cfg,
		tools:    tool.NewRegistry(),
		hooks:    hook.NewRegistry(),
		sessions: session.NewManager(cfg),
		agents:   make(map[string]Agent),
	}

	tool.SetEnabledFromConfig(cfg.Tools.ToolStates)

	homeDir, _ := os.UserHomeDir()
	externalToolsDir := filepath.Join(homeDir, ".config", "freecode", "tools")
	tool.LoadExternalTools(externalToolsDir)
	tool.CompileExternalTools(externalToolsDir, "")

	eng.registerTools()

	return eng
}

func (e *Engine) registerTools() {
	for _, name := range tool.ListTools() {
		factory, ok := tool.GetFactory(name)
		if !ok {
			continue
		}

		if !tool.IsEnabled(name) {
			continue
		}

		if len(e.config.Tools.Allowed) > 0 {
			allowed := false
			for _, a := range e.config.Tools.Allowed {
				if a == name {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}

		for _, d := range e.config.Tools.Denied {
			if d == name {
				continue
			}
		}

		e.RegisterTool(factory())
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

func (e *Engine) EnableTool(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	factory, ok := tool.GetFactory(name)
	if !ok {
		return fmt.Errorf("tool not found: %s", name)
	}

	tool.SetEnabled(name, true)
	e.tools.Enable(name)

	if e.config.Tools.ToolStates == nil {
		e.config.Tools.ToolStates = make(map[string]bool)
	}
	e.config.Tools.ToolStates[name] = true

	e.RegisterTool(factory())

	return nil
}

func (e *Engine) DisableTool(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	tool.SetEnabled(name, false)
	e.tools.Disable(name)

	if e.config.Tools.ToolStates == nil {
		e.config.Tools.ToolStates = make(map[string]bool)
	}
	e.config.Tools.ToolStates[name] = false

	return nil
}

func (e *Engine) SaveToolStates() error {
	return e.config.Save()
}

func (e *Engine) SessionManager() *session.Manager {
	return e.sessions
}

func (e *Engine) HookRegistry() *hook.Registry {
	return e.hooks
}
