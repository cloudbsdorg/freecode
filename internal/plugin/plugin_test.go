package plugin

import (
	"context"
	"testing"
)

type testPlugin struct {
	name string
}

func (p *testPlugin) Name() string                    { return p.name }
func (p *testPlugin) Init(ctx context.Context) error  { return nil }
func (p *testPlugin) Close() error                    { return nil }

func TestRegistry(t *testing.T) {
	r := NewMemoryRegistry()

	err := r.Register(&testPlugin{name: "test"})
	if err != nil {
		t.Errorf("Register error: %v", err)
	}

	plugin, err := r.Get("test")
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if plugin.Name() != "test" {
		t.Errorf("expected test, got %s", plugin.Name())
	}

	list := r.List()
	if len(list) != 1 || list[0] != "test" {
		t.Errorf("expected [test], got %v", list)
	}

	err = r.Disable("test")
	if err != nil {
		t.Errorf("Disable error: %v", err)
	}

	_, err = r.Get("test")
	if err != ErrPluginDisabled {
		t.Errorf("expected ErrPluginDisabled, got %v", err)
	}

	err = r.Enable("test")
	if err != nil {
		t.Errorf("Enable error: %v", err)
	}

	err = r.Unregister("test")
	if err != nil {
		t.Errorf("Unregister error: %v", err)
	}

	_, err = r.Get("test")
	if err != ErrPluginNotFound {
		t.Errorf("expected ErrPluginNotFound, got %v", err)
	}
}
