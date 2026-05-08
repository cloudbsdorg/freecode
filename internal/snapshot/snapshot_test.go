package snapshot

import (
	"context"
	"testing"
	"time"
)

func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	if store == nil {
		t.Fatal("NewMemoryStore() returned nil")
	}
}

func TestMemoryStoreCreate(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	data := map[string]any{"key": "value"}
	snap, err := store.Create(ctx, data)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if snap == nil {
		t.Fatal("Create() returned nil")
	}
	if snap.ID == "" {
		t.Error("Snapshot ID is empty")
	}
	if snap.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
	if snap.Data["key"] != "value" {
		t.Errorf("Data[key] = %v, want %v", snap.Data["key"], "value")
	}
}

func TestMemoryStoreGet(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	data := map[string]any{"key": "value"}
	created, _ := store.Create(ctx, data)

	retrieved, err := store.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if retrieved == nil {
		t.Fatal("Get() returned nil")
	}
	if retrieved.ID != created.ID {
		t.Errorf("ID = %q, want %q", retrieved.ID, created.ID)
	}

	nonexistent, err := store.Get(ctx, "nonexistent-id")
	if err != nil {
		t.Fatalf("Get() for nonexistent returned error: %v", err)
	}
	if nonexistent != nil {
		t.Error("Get() for nonexistent should return nil")
	}
}

func TestMemoryStoreList(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	snaps, err := store.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(snaps) != 0 {
		t.Errorf("List() on empty store returned %d snaps, want 0", len(snaps))
	}

	store.Create(ctx, map[string]any{"key1": "value1"})
	store.Create(ctx, map[string]any{"key2": "value2"})

	snaps, err = store.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(snaps) != 2 {
		t.Errorf("List() returned %d snaps, want 2", len(snaps))
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	snap, _ := store.Create(ctx, map[string]any{"key": "value"})

	err := store.Delete(ctx, snap.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	retrieved, _ := store.Get(ctx, snap.ID)
	if retrieved != nil {
		t.Error("Get() after Delete() should return nil")
	}

	err = store.Delete(ctx, "nonexistent-id")
	if err != nil {
		t.Errorf("Delete() for nonexistent returned error: %v", err)
	}
}

func TestSnapshot(t *testing.T) {
	now := time.Now()
	snap := &Snapshot{
		ID:        "test-id",
		CreatedAt: now,
		Data:      map[string]any{"key": "value"},
	}

	if snap.ID != "test-id" {
		t.Errorf("ID = %q, want %q", snap.ID, "test-id")
	}
	if !snap.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", snap.CreatedAt, now)
	}
	if snap.Data["key"] != "value" {
		t.Errorf("Data[key] = %v, want %v", snap.Data["key"], "value")
	}
}