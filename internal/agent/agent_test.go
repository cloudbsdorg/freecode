package agent

import (
    "context"
    "testing"
    "github.com/freecode/freecode/internal/config"
)

// Helper to create a test engine with defaults
func newTestEngine() *Engine {
    cfg := config.DefaultConfig()
    return NewEngine(cfg)
}

// Use the shared mockAgent type defined in engine_test.go

func TestEngine_RunAndRegisterAgent(t *testing.T) {
    eng := newTestEngine()
    // Register a mock agent and verify GetAgent works and Run dispatches
    ma := &mockAgent{nameVal: "mock-01"}
    eng.RegisterAgent(ma)
    a, ok := eng.GetAgent("mock-01")
    if !ok || a == nil {
        t.Fatalf("expected to find registered agent 'mock-01'")
    }

    ctx := context.Background()
    req := Request{SessionID: "sess-1", AgentName: "mock-01", Message: Message{Content: "hello"}}
    resp, err := eng.Run(ctx, req)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp == nil || resp.SessionID != "sess-1" {
        t.Fatalf("unexpected response: %+v", resp)
    }
}

func TestEngine_RunUnknownAgentErrors(t *testing.T) {
    eng := newTestEngine()
    ctx := context.Background()
    req := Request{SessionID: "s", AgentName: "not-found", Message: Message{Content: "x"}}
    if _, err := eng.Run(ctx, req); err == nil {
        t.Fatalf("expected error for unknown agent, got nil")
    }
}

func TestAllSisyphusLikeAgents_RunContent(t *testing.T) {
    eng := newTestEngine()
    // Register all agents to exercise their Run paths
    agents := []Agent{
        NewSisyphusAgent(eng),
        NewHephaestusAgent(eng),
        NewOracleAgent(eng),
        NewLibrarianAgent(eng),
        NewExploreAgent(eng),
        NewPrometheusAgent(eng),
        NewMetisAgent(eng),
        NewMomusAgent(eng),
        NewAtlasAgent(eng),
        NewMultimodalLookerAgent(eng),
        NewSisyphusJuniorAgent(eng),
    }
    ctx := context.Background()
    req := Request{SessionID: "sess-xx", Message: Message{Content: "do task"}}
    for _, ac := range agents {
        resp, err := ac.Run(ctx, req)
        if err != nil {
            t.Fatalf("agent %s returned error: %v", ac.Name(), err)
        }
        if resp == nil {
            t.Fatalf("agent %s returned nil response", ac.Name())
        }
        if resp.SessionID != req.SessionID {
            t.Fatalf("agent %s returned wrong session: %s", ac.Name(), resp.SessionID)
        }
        if resp.Message.Content == "" {
            t.Fatalf("agent %s returned empty content", ac.Name())
        }
    }
}

// Ensure exported constructors return non-nil names for a quick sanity check
func TestAgentNamesAreNonEmpty(t *testing.T) {
    eng := newTestEngine()
    tests := []Agent{
        NewSisyphusAgent(eng), NewHephaestusAgent(eng), NewOracleAgent(eng),
        NewLibrarianAgent(eng), NewExploreAgent(eng), NewPrometheusAgent(eng),
        NewMetisAgent(eng), NewMomusAgent(eng), NewAtlasAgent(eng),
        NewMultimodalLookerAgent(eng), NewSisyphusJuniorAgent(eng),
    }
    for _, a := range tests {
        if a.Name() == "" {
            t.Fatalf("agent has empty name: %#v", a)
        }
    }
}
