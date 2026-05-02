package hook

import (
	"context"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()

	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	if r.sessionHooks == nil {
		t.Error("Registry.sessionHooks is nil")
	}

	if r.toolHooks == nil {
		t.Error("Registry.toolHooks is nil")
	}

	if r.transformHooks == nil {
		t.Error("Registry.transformHooks is nil")
	}

	if r.continuationHooks == nil {
		t.Error("Registry.continuationHooks is nil")
	}
}

func TestRegistryRegisterSessionHook(t *testing.T) {
	r := NewRegistry()

	r.RegisterSessionHook("test", func(ctx context.Context, evt SessionEvent) error {
		return nil
	})

	hooks := r.sessionHooks["test"]
	if len(hooks) != 1 {
		t.Errorf("len(sessionHooks[test]) = %d, want 1", len(hooks))
	}
}

func TestRegistryEmitSessionEvent(t *testing.T) {
	r := NewRegistry()
	called := false

	r.RegisterSessionHook("start", func(ctx context.Context, evt SessionEvent) error {
		called = true
		if evt.Type != "start" {
			t.Errorf("Event.Type = %q, want %q", evt.Type, "start")
		}
		if evt.SessionID != "sess123" {
			t.Errorf("Event.SessionID = %q, want %q", evt.SessionID, "sess123")
		}
		return nil
	})

	err := r.EmitSessionEvent(context.Background(), "start", "sess123", nil)
	if err != nil {
		t.Errorf("EmitSessionEvent() error = %v", err)
	}

	if !called {
		t.Error("Hook was not called")
	}
}

func TestRegistryEmitSessionEventNoHook(t *testing.T) {
	r := NewRegistry()

	err := r.EmitSessionEvent(context.Background(), "nonexistent", "sess", nil)
	if err != nil {
		t.Errorf("EmitSessionEvent() error = %v, want nil for unregistered event", err)
	}
}

func TestRegistryRegisterToolHook(t *testing.T) {
	r := NewRegistry()

	r.RegisterToolHook("before:bash", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return nil, false
	})

	hooks := r.toolHooks["before:bash"]
	if len(hooks) != 1 {
		t.Errorf("len(toolHooks[before:bash]) = %d, want 1", len(hooks))
	}
}

func TestRegistryEmitToolEventNoHook(t *testing.T) {
	r := NewRegistry()

	err, handled := r.EmitToolEvent(context.Background(), "before", "bash", "sess", nil)
	if err != nil {
		t.Errorf("EmitToolEvent() error = %v, want nil", err)
	}
	if handled {
		t.Error("EmitToolEvent() returned handled = true for unregistered hook")
	}
}

func TestRegistryRegisterTransformHook(t *testing.T) {
	r := NewRegistry()

	r.RegisterTransformHook(func(msg *Message) (*Message, error) {
		return msg, nil
	})

	if len(r.transformHooks) != 1 {
		t.Errorf("len(transformHooks) = %d, want 1", len(r.transformHooks))
	}
}

func TestRegistryApplyTransformHooks(t *testing.T) {
	r := NewRegistry()

	msg := &Message{Role: "user", Content: "original"}

	r.RegisterTransformHook(func(msg *Message) (*Message, error) {
		msg.Content = "transformed"
		return msg, nil
	})

	result, err := r.ApplyTransformHooks(msg)
	if err != nil {
		t.Errorf("ApplyTransformHooks() error = %v", err)
	}

	if result.Content != "transformed" {
		t.Errorf("Content = %q, want %q", result.Content, "transformed")
	}
}

func TestRegistryRegisterContinuationHook(t *testing.T) {
	r := NewRegistry()

	r.RegisterContinuationHook(func(ctx context.Context, session *SessionData) (*ContinueSignal, error) {
		return &ContinueSignal{Continue: true}, nil
	})

	if len(r.continuationHooks) != 1 {
		t.Errorf("len(continuationHooks) = %d, want 1", len(r.continuationHooks))
	}
}

func TestRegistryCheckContinuation(t *testing.T) {
	r := NewRegistry()

	r.RegisterContinuationHook(func(ctx context.Context, session *SessionData) (*ContinueSignal, error) {
		return &ContinueSignal{Continue: true, Reason: "continue"}, nil
	})

	signal, err := r.CheckContinuation(context.Background(), &SessionData{})
	if err != nil {
		t.Errorf("CheckContinuation() error = %v", err)
	}

	if !signal.Continue {
		t.Error("signal.Continue = false, want true")
	}
}

func TestRegistryCheckContinuationStop(t *testing.T) {
	r := NewRegistry()

	r.RegisterContinuationHook(func(ctx context.Context, session *SessionData) (*ContinueSignal, error) {
		return &ContinueSignal{Continue: false, Reason: "stop"}, nil
	})

	signal, err := r.CheckContinuation(context.Background(), &SessionData{})
	if err != nil {
		t.Errorf("CheckContinuation() error = %v", err)
	}

	if signal.Continue {
		t.Error("signal.Continue = true, want false")
	}
}
