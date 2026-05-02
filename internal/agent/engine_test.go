package agent

import (
	"context"
	"testing"

	"github.com/freecode/freecode/internal/config"
)

func TestNewEngine(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	if e == nil {
		t.Fatal("NewEngine() returned nil")
	}

	if e.config == nil {
		t.Error("Engine.config is nil")
	}

	if e.tools == nil {
		t.Error("Engine.tools is nil")
	}

	if e.hooks == nil {
		t.Error("Engine.hooks is nil")
	}

	if e.sessions == nil {
		t.Error("Engine.sessions is nil")
	}

	if e.agents == nil {
		t.Error("Engine.agents is nil")
	}
}

type mockAgent struct {
	nameVal string
}

func (m *mockAgent) Name() string {
	return m.nameVal
}

func (m *mockAgent) Run(ctx context.Context, req Request) (*Response, error) {
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: "mock response",
		},
	}, nil
}

func TestEngineRegisterAgent(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	agent := &mockAgent{nameVal: "test-agent"}
	e.RegisterAgent(agent)

	if len(e.agents) != 1 {
		t.Errorf("len(Engine.agents) = %d, want 1", len(e.agents))
	}
}

func TestEngineGetAgent(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	agent := &mockAgent{nameVal: "get-agent"}
	e.RegisterAgent(agent)

	got, ok := e.GetAgent("get-agent")
	if !ok {
		t.Error("GetAgent() returned false for registered agent")
	}

	if got.Name() != "get-agent" {
		t.Errorf("got.Name() = %q, want %q", got.Name(), "get-agent")
	}
}

func TestEngineGetAgentNotFound(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	_, ok := e.GetAgent("nonexistent")
	if ok {
		t.Error("GetAgent() returned true for nonexistent agent")
	}
}

func TestEngineRun(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	agent := &mockAgent{nameVal: "run-agent"}
	e.RegisterAgent(agent)

	req := Request{
		SessionID: "sess123",
		AgentName: "run-agent",
		Message: Message{
			Role:    "user",
			Content: "hello",
		},
	}

	resp, err := e.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if resp.SessionID != "sess123" {
		t.Errorf("Response.SessionID = %q, want %q", resp.SessionID, "sess123")
	}

	if resp.Message.Content != "mock response" {
		t.Errorf("Response.Message.Content = %q, want %q", resp.Message.Content, "mock response")
	}
}

func TestEngineRunAgentNotFound(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	req := Request{
		SessionID: "sess123",
		AgentName: "nonexistent",
	}

	_, err := e.Run(context.Background(), req)
	if err == nil {
		t.Error("Run() should error for nonexistent agent")
	}
}

func TestEngineSessionManager(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	sm := e.SessionManager()
	if sm == nil {
		t.Error("SessionManager() returned nil")
	}
}

func TestEngineHookRegistry(t *testing.T) {
	cfg := config.DefaultConfig()
	e := NewEngine(cfg)

	hr := e.HookRegistry()
	if hr == nil {
		t.Error("HookRegistry() returned nil")
	}
}
