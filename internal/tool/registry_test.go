package tool

import (
	"context"
	"testing"
)

func TestRegistryNew(t *testing.T) {
	r := NewRegistry()

	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	if r.tools == nil {
		t.Error("Registry.tools is nil")
	}

	if len(r.tools) != 0 {
		t.Errorf("len(Registry.tools) = %d, want 0", len(r.tools))
	}
}

type mockTool struct {
	name        string
	description string
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Description() string {
	return m.description
}

func (m *mockTool) Schema() ToolSchema {
	return ToolSchema{Name: m.name, Description: m.description}
}

func (m *mockTool) Execute(ctx context.Context, req Request) (*Response, error) {
	return &Response{Result: "result"}, nil
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()

	tool := &mockTool{name: "test", description: "test tool"}
	r.Register(tool)

	if len(r.tools) != 1 {
		t.Errorf("len(Registry.tools) = %d, want 1", len(r.tools))
	}

	if _, ok := r.tools["test"]; !ok {
		t.Error("tool 'test' not found in registry")
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()

	tool := &mockTool{name: "gettest", description: "get test tool"}
	r.Register(tool)

	got, ok := r.Get("gettest")
	if !ok {
		t.Error("Get() returned false for registered tool")
	}

	if got.Name() != "gettest" {
		t.Errorf("got.Name() = %q, want %q", got.Name(), "gettest")
	}
}

func TestRegistryGetNotFound(t *testing.T) {
	r := NewRegistry()

	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("Get() returned true for nonexistent tool")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewRegistry()

	r.Register(&mockTool{name: "tool1", description: "tool 1"})
	r.Register(&mockTool{name: "tool2", description: "tool 2"})
	r.Register(&mockTool{name: "tool3", description: "tool 3"})

	tools := r.List()

	if len(tools) != 3 {
		t.Errorf("len(List()) = %d, want 3", len(tools))
	}

	names := make(map[string]bool)
	for _, t := range tools {
		names[t.Name()] = true
	}

	for _, name := range []string{"tool1", "tool2", "tool3"} {
		if !names[name] {
			t.Errorf("List() missing tool %q", name)
		}
	}
}

func TestRegistryExecute(t *testing.T) {
	r := NewRegistry()

	r.Register(&mockTool{name: "exec", description: "exec tool"})

	resp, err := r.Execute(context.Background(), Request{Name: "exec"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if resp.Result != "result" {
		t.Errorf("Result = %q, want %q", resp.Result, "result")
	}
}

func TestRegistryExecuteNotFound(t *testing.T) {
	r := NewRegistry()

	_, err := r.Execute(context.Background(), Request{Name: "nonexistent"})
	if err == nil {
		t.Error("Execute() should error for nonexistent tool")
	}
}
