package sync

import (
	"context"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	session := &Session{ID: "test", Data: map[string]any{"key": "value"}}
	if err := s.Set(ctx, session); err != nil {
		t.Errorf("Set error: %v", err)
	}

	retrieved, err := s.Get(ctx, "test")
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if retrieved == nil {
		t.Error("expected session")
	}

	if err := s.Delete(ctx, "test"); err != nil {
		t.Errorf("Delete error: %v", err)
	}

	retrieved, _ = s.Get(ctx, "test")
	if retrieved != nil {
		t.Error("expected nil after delete")
	}
}
