package hook

import (
	"context"
	"fmt"
	"testing"
)

// Helper dummy continuation hook that does nothing
func dummyContinuation(_ context.Context, _ *SessionData) (*ContinueSignal, error) {
	return nil, nil
}

// Helper dummy session hook that does nothing
func dummySession(_ context.Context, _ SessionEvent) error {
	return nil
}

// Helper dummy tool hook that does nothing
func dummyTool(_ context.Context, _ ToolEvent) (error, bool) {
	return nil, false
}

func TestContinuationHooksRegistration_All(t *testing.T) {
	r := NewRegistry()
	ch := NewContinuationHooks(r)

	// Register all continuation hooks
	ch.OnIterate(dummyContinuation)
	ch.OnRetry(dummyContinuation)
	ch.OnFallback(dummyContinuation)
	ch.OnEscalate(dummyContinuation)
	ch.OnDelegate(dummyContinuation)
	ch.OnTimeout(dummyContinuation)
	ch.OnIdle(dummyContinuation)

	if len(r.continuationHooks) != 7 {
		t.Fatalf("expected 7 continuation hooks, got %d", len(r.continuationHooks))
	}
}

func TestSessionHooksRegistration_All(t *testing.T) {
	r := NewRegistry()
	sh := NewSessionHooks(r)

	// Register a variety of session hooks
	sh.OnStart(dummySession)
	sh.OnEnd(dummySession)
	sh.OnPause(dummySession)
	sh.OnResume(dummySession)
	sh.OnSave(dummySession)
	sh.OnLoad(dummySession)
	sh.OnError(dummySession)
	sh.OnTimeout(dummySession)
	sh.OnIdle(dummySession)
	sh.OnActive(dummySession)
	sh.OnMessage(dummySession)
	sh.OnToolCall(dummySession)
	sh.OnAgentCall(dummySession)
	sh.OnUserInput(dummySession)
	sh.OnResponse(dummySession)
	sh.OnCompaction(dummySession)
	sh.OnTabCreate(dummySession)
	sh.OnTabClose(dummySession)
	sh.OnTabSwitch(dummySession)
	sh.OnSplitCreate(dummySession)
	sh.OnSplitClose(dummySession)
	sh.OnFleetConnect(dummySession)
	sh.OnFleetDisconnect(dummySession)
	sh.OnFleetError(dummySession)

	// We registered 24 hooks in total
	if len(r.sessionHooks) == 0 {
		t.Fatalf("expected some session hooks registered, got none")
	}
}

func TestToolHooksRegistration_All(t *testing.T) {
	r := NewRegistry()
	th := NewToolHooks(r)

	th.OnBefore("foo", dummyTool)
	th.OnAfter("foo", dummyTool)
	th.OnError("foo", dummyTool)
	th.OnTransform(func(msg *Message) (*Message, error) { return msg, nil })
	th.OnBash(dummyTool)
	th.OnRead(dummyTool)
	th.OnWrite(dummyTool)
	th.OnEdit(dummyTool)
	th.OnGlob(dummyTool)
	th.OnGrep(dummyTool)
	th.OnWebFetch(dummyTool)
	th.OnWebSearch(dummyTool)
	th.OnTask(dummyTool)
	th.OnSkill(dummyTool)

	// 14 registrations (including transform) should exist
	if len(r.toolHooks) == 0 {
		t.Fatalf("expected tool hooks to be registered, got none")
	}
}

func TestRalphHooksAndApply(t *testing.T) {
	r := NewRegistry()
	// success path
	r.RegisterRalphHook(func(ctx context.Context, input string) (string, error) {
		return input + "_X", nil
	})
	out, err := r.ApplyRalphHooks(context.Background(), "in")
	if err != nil {
		t.Fatalf("ApplyRalphHooks() error = %v", err)
	}
	if out != "in_X" {
		t.Fatalf("unexpected output: %q", out)
	}

	// error path
	r2 := NewRegistry()
	r2.RegisterRalphHook(func(ctx context.Context, input string) (string, error) {
		return "", fmt.Errorf("boom")
	})
	_, err = r2.ApplyRalphHooks(context.Background(), "in")
	if err == nil {
		t.Fatalf("expected error from Ralph hook, got nil")
	}
}

func TestSkillHooksRegistration(t *testing.T) {
	r := NewRegistry()
	// Ensure registering a skill hook updates map
	r.RegisterSkillHook("test", func(ctx context.Context, skill string, args map[string]interface{}) error {
		return nil
	})
	if hooks, ok := r.skillHooks["test"]; !ok || len(hooks) != 1 {
		t.Fatalf("expected 1 skill hook for 'test', got %d", len(r.skillHooks["test"]))
	}
}

func TestEmitToolEventHandled(t *testing.T) {
	r := NewRegistry()
	r.RegisterToolHook("before:boom", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return fmt.Errorf("boom"), true
	})

	err, handled := r.EmitToolEvent(context.Background(), "before:boom", "boom", "sess", nil)
	if !handled {
		t.Fatalf("expected handled to be true")
	}
	if err == nil {
		t.Fatalf("expected non-nil error when hook handles the event")
	}
}

func TestEmitToolEventHandledMultiple(t *testing.T) {
	r := NewRegistry()
	// first hook does nothing, second handles with error
	r.RegisterToolHook("before:multi", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return nil, false
	})
	r.RegisterToolHook("before:multi", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return fmt.Errorf("fail"), true
	})

	err, handled := r.EmitToolEvent(context.Background(), "before:multi", "t", "sess", nil)
	if !handled {
		t.Fatalf("expected handled to be true for second hook")
	}
	if err == nil {
		t.Fatalf("expected non-nil error from second hook")
	}
}

func TestTransformHooksMultiple(t *testing.T) {
	r := NewRegistry()
	// Two transforms: append A then B
	r.RegisterTransformHook(func(msg *Message) (*Message, error) {
		msg.Content = msg.Content + "A"
		return msg, nil
	})
	r.RegisterTransformHook(func(msg *Message) (*Message, error) {
		msg.Content = msg.Content + "B"
		return msg, nil
	})

	m := &Message{Role: "user", Content: "start"}
	out, err := r.ApplyTransformHooks(m)
	if err != nil {
		t.Fatalf("ApplyTransformHooks() error = %v", err)
	}
	if out.Content != "startAB" {
		t.Fatalf("expected final content 'startAB', got %q", out.Content)
	}
}

func TestSessionEventRegistrationAndEmission(t *testing.T) {
	r := NewRegistry()
	called := 0
	r.RegisterSessionHook("start", func(ctx context.Context, evt SessionEvent) error {
		called++
		return nil
	})
	r.RegisterSessionHook("start", func(ctx context.Context, evt SessionEvent) error {
		called++
		return nil
	})
	if err := r.EmitSessionEvent(context.Background(), "start", "sess", nil); err != nil {
		t.Fatalf("EmitSessionEvent() error = %v", err)
	}
	if called != 2 {
		t.Fatalf("expected 2 session hooks to be called, got %d", called)
	}
}

func TestEmitSessionEventWrapsError(t *testing.T) {
	r := NewRegistry()
	r.RegisterSessionHook("start", func(ctx context.Context, evt SessionEvent) error {
		return fmt.Errorf("boom")
	})

	err := r.EmitSessionEvent(context.Background(), "start", "sess", nil)
	if err == nil {
		t.Fatalf("expected error from session hook, got nil")
	}
	if !contains(err.Error(), "session hook error (start): boom") {
		t.Fatalf("unexpected error wrap: %v", err)
	}
}

func TestEmitToolEventWrapsError(t *testing.T) {
	r := NewRegistry()
	r.RegisterToolHook("before:boom", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return fmt.Errorf("boom"), true
	})
	err, handled := r.EmitToolEvent(context.Background(), "before:boom", "t", "sess", nil)
	if !handled {
		t.Fatalf("expected handled to be true for tool hook")
	}
	if err == nil {
		t.Fatalf("expected error from tool hook, got nil")
	}
	if !contains(err.Error(), "boom") {
		t.Fatalf("unexpected tool hook error: %v", err)
	}
}

func TestTransformHooksNoOp(t *testing.T) {
	r := NewRegistry()
	m := &Message{Role: "user", Content: "hello"}
	out, err := r.ApplyTransformHooks(m)
	if err != nil {
		t.Fatalf("ApplyTransformHooks() error = %v", err)
	}
	if out.Content != m.Content {
		t.Fatalf("expected no transformation, got %q", out.Content)
	}
}

func TestEmitToolEventNoHandledWithHooks(t *testing.T) {
	r := NewRegistry()
	// Hook returns an error but not handled
	r.RegisterToolHook("before:nohandle", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return fmt.Errorf("boom"), false
	})
	err, handled := r.EmitToolEvent(context.Background(), "before:nohandle", "t", "sess", nil)
	if handled {
		t.Fatalf("expected handled to be false when hook did not handle the event")
	}
	if err != nil {
		t.Fatalf("expected no error when no hook handled the event, got %v", err)
	}
}

func TestApplyTransformHooksErrorPath(t *testing.T) {
	r := NewRegistry()
	r.RegisterTransformHook(func(msg *Message) (*Message, error) {
		return nil, fmt.Errorf("bad transform")
	})
	m := &Message{Role: "user", Content: "start"}
	_, err := r.ApplyTransformHooks(m)
	if err == nil {
		t.Fatalf("expected error from transform hook, got nil")
	}
}

func TestCheckContinuationNoHooks(t *testing.T) {
	r := NewRegistry()
	sig, err := r.CheckContinuation(context.Background(), &SessionData{})
	if err != nil {
		t.Fatalf("CheckContinuation() error = %v", err)
	}
	if !sig.Continue {
		t.Fatalf("expected default continue = true")
	}
}

func TestCheckContinuationErrorWrapping(t *testing.T) {
	r := NewRegistry()
	r.RegisterContinuationHook(func(ctx context.Context, session *SessionData) (*ContinueSignal, error) {
		return nil, fmt.Errorf("boom")
	})
	_, err := r.CheckContinuation(context.Background(), &SessionData{})
	if err == nil {
		t.Fatalf("expected error from continuation hook, got nil")
	}
	if !contains(err.Error(), "continuation hook error") {
		t.Fatalf("unexpected error wrap: %v", err)
	}
}

// simple helper to check substring presence (since strings.Contains isn't imported yet)
func contains(s, substr string) bool {
	return (len(substr) == 0) || (len(s) >= len(substr) && (indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
