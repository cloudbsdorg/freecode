package bus

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestPublishSubscribe(t *testing.T) {
	Close()
	var received int32

	def := Define("test.event", struct{ Message string }{})

	ctx := context.Background()

	unsubscribe := Subscribe(def, func(payload Payload) {
		atomic.AddInt32(&received, 1)
	})

	Publish(ctx, def, struct{ Message string }{Message: "hello"})

	if received != 1 {
		t.Errorf("expected 1, got %d", received)
	}

	unsubscribe()

	Publish(ctx, def, struct{ Message string }{Message: "world"})

	if received != 1 {
		t.Errorf("expected 1 after unsubscribe, got %d", received)
	}
}

func TestSubscribeAll(t *testing.T) {
	Close()
	var received int32

	unsubscribe := SubscribeAll(func(event GlobalEvent) {
		atomic.AddInt32(&received, 1)
	})

	Publish(context.Background(), Define("test.1", nil), nil)
	Publish(context.Background(), Define("test.2", nil), nil)

	if received != 2 {
		t.Errorf("expected 2, got %d", received)
	}

	unsubscribe()

	Publish(context.Background(), Define("test.3", nil), nil)

	if received != 2 {
		t.Errorf("expected 2 after unsubscribe, got %d", received)
	}
}

func TestGlobalBus(t *testing.T) {
	bus := GetGlobalBus()
	var received int32

	cleanup := bus.On(func(event GlobalEvent) {
		atomic.AddInt32(&received, 1)
	})

	bus.Emit(GlobalEvent{Directory: "/test", Project: "proj", Workspace: "ws"})

	if received != 1 {
		t.Errorf("expected 1, got %d", received)
	}

	cleanup()

	bus.Emit(GlobalEvent{Directory: "/test2"})

	if received != 1 {
		t.Errorf("expected 1 after cleanup, got %d", received)
	}
}
