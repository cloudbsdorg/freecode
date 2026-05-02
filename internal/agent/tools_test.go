package agent

import (
    "context"
    "testing"
    "github.com/freecode/freecode/internal/config"
)

// Use the shared fakeTool type defined in tool_registry_test.go

func TestToolCaller_CallAndCallBatch(t *testing.T) {
    eng := NewEngine(config.DefaultConfig())
    ft := &fakeTool{name: "t1"}
    eng.RegisterTool(ft)
    tc := NewToolCaller()
    ctx := context.Background()

    // Test Call
    resp, err := tc.Call(ctx, "t1", map[string]interface{}{"a": 1}, eng)
    if err != nil {
        t.Fatalf("Call() error: %v", err)
    }
    if resp == nil || resp.Result != "ok" {
        t.Fatalf("unexpected Call() response: %+v", resp)
    }

    // Test CallBatch with two calls
    calls := []ToolCall{
        {Name: "t1", Arguments: map[string]interface{}{"a": 1}},
        {Name: "t1", Arguments: map[string]interface{}{"b": 2}},
    }
    res, err := tc.CallBatch(ctx, calls, eng)
    if err != nil {
        t.Fatalf("CallBatch() error: %v", err)
    }
    if len(res) != 2 {
        t.Fatalf("expected 2 results, got %d", len(res))
    }
    for i, r := range res {
        if r == nil || r.Result != "ok" {
            t.Fatalf("unexpected batch result %d: %+v", i, r)
        }
    }
}

func TestToolCaller_CallNotFound(t *testing.T) {
    eng := NewEngine(config.DefaultConfig())
    tc := NewToolCaller()
    ctx := context.Background()
    _, err := tc.Call(ctx, "nonexistent", nil, eng)
    if err == nil {
        t.Fatalf("expected error for missing tool, got nil")
    }
}

func TestToolCaller_CallBatchWithMissing(t *testing.T) {
    eng := NewEngine(config.DefaultConfig())
    ft := &fakeTool{name: "t1"}
    eng.RegisterTool(ft)
    tc := NewToolCaller()
    ctx := context.Background()
    calls := []ToolCall{
        {Name: "t1", Arguments: map[string]interface{}{"a": 1}},
        {Name: "missing", Arguments: map[string]interface{}{"b": 2}},
    }
    res, err := tc.CallBatch(ctx, calls, eng)
    if err != nil {
        t.Fatalf("CallBatch() error: %v", err)
    }
    if len(res) != 2 {
        t.Fatalf("expected 2 results, got %d", len(res))
    }
    if res[0] == nil || res[0].Result != "ok" {
        t.Fatalf("unexpected first batch result: %#v", res[0])
    }
    if res[1] == nil || res[1].Error == nil {
        t.Fatalf("expected error in second batch result for missing tool, got %#v", res[1])
    }
}
