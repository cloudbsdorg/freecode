package agent

import (
	"context"
	"fmt"
	"github.com/freecode/freecode/internal/config"
	"testing"
)

func TestMessageHandler_HandleAndFormat(t *testing.T) {
	eng := NewEngine(config.DefaultConfig())
	// Register a simple agent to handle the message
	eng.RegisterAgent(NewSisyphusAgent(eng))
	mh := NewMessageHandler()
	ctx := context.Background()
	msg := Message{Content: "do something"}
	resp, err := mh.Handle(ctx, msg, eng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.SessionID != "" {
		// session id is not set by the handler, only by agent; still ensure response exists
		if resp == nil {
			t.Fatalf("nil response")
		}
	}
	if resp.Message.Content == "" {
		t.Fatalf("empty content in response: %+v", resp)
	}

	// Test FormatResponse without error
	got := mh.FormatResponse(resp)
	if got != resp.Message.Content {
		t.Fatalf("expected content %q, got %q", resp.Message.Content, got)
	}

	// Test error formatting
	respErr := &Response{Error: fmt.Errorf("boom"), SessionID: resp.SessionID, Message: resp.Message}
	if sug := mh.FormatResponse(respErr); sug != "Error: boom" {
		t.Fatalf("expected error format, got %q", sug)
	}
}
