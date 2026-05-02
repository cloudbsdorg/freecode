package agent

import (
    "context"
    "testing"
    "github.com/freecode/freecode/internal/config"
    "github.com/freecode/freecode/internal/tool"
)

// fakeTool implements the tool.Tool interface for testing
type fakeTool struct{ name string }

func (f *fakeTool) Name() string { return f.name }
func (f *fakeTool) Description() string { return "fake tool" }
func (f *fakeTool) Schema() tool.ToolSchema {
    return tool.ToolSchema{Name: f.name, Description: "fake tool"}
}
func (f *fakeTool) Execute(ctx context.Context, req tool.Request) (*tool.Response, error) {
    return &tool.Response{Result: "ok"}, nil
}

func TestEngineToolRegistry(t *testing.T) {
    eng := NewEngine(config.DefaultConfig())
    ft := &fakeTool{name: "fake-tool"}
    eng.RegisterTool(ft)

    // Get and execute the tool
    ttt, ok := eng.GetTool("fake-tool")
    if !ok {
        t.Fatalf("expected to find registered tool")
    }
    resp, err := ttt.Execute(context.Background(), tool.Request{Name: "fake-tool"})
    if err != nil {
        t.Fatalf("tool execute error: %v", err)
    }
    if resp == nil || resp.Result != "ok" {
        t.Fatalf("unexpected tool response: %+v", resp)
    }

    // ListTools should include the tool
    tools := eng.ListTools()
    found := false
    for _, tt := range tools {
        if tt.Name() == "fake-tool" {
            found = true
            break
        }
    }
    if !found {
        t.Fatalf("registered tool not found in ListTools() output")
    }
}
