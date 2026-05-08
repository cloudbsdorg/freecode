package sync

import (
	"context"
	"testing"
	"time"
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
	if retrieved.Version != 1 {
		t.Errorf("expected version 1, got %d", retrieved.Version)
	}

	if err := s.Delete(ctx, "test"); err != nil {
		t.Errorf("Delete error: %v", err)
	}

	retrieved, _ = s.Get(ctx, "test")
	if retrieved != nil {
		t.Error("expected nil after delete")
	}
}

func TestSyncNoConflict(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	session := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "value"},
		Vector:  VersionVector{"client1": 1},
	}

	result, err := s.Sync(ctx, session, LastWriteWins)
	if err != nil {
		t.Errorf("Sync error: %v", err)
	}
	if result.Version != 1 {
		t.Errorf("expected version 1, got %d", result.Version)
	}
}

func TestSyncNewSession(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	session := &Session{
		ID:      "new-session",
		Data:    map[string]any{"key": "value"},
		Vector:  VersionVector{"client1": 1},
	}

	result, err := s.Sync(ctx, session, LastWriteWins)
	if err != nil {
		t.Errorf("Sync error: %v", err)
	}
	if result.ID != "new-session" {
		t.Errorf("expected ID new-session, got %s", result.ID)
	}
	if result.Version != 1 {
		t.Errorf("expected version 1, got %d", result.Version)
	}
}

func TestSyncConflictLastWriteWins(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	serverSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "server-value"},
		Updated: time.Now().Add(-1 * time.Hour),
		Vector:  VersionVector{"client1": 1, "server": 1},
	}
	s.Set(ctx, serverSession)

	clientSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "client-value"},
		Updated: time.Now(),
		Vector:  VersionVector{"client1": 1},
	}

	result, err := s.Sync(ctx, clientSession, LastWriteWins)
	if err != nil {
		t.Errorf("Sync error: %v", err)
	}
	if result.Data["key"] != "client-value" {
		t.Errorf("expected client-value (newer), got %v", result.Data["key"])
	}
}

func TestSyncConflictServerWins(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	serverSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "server-value"},
		Updated: time.Now(),
		Vector:  VersionVector{"client1": 1, "server": 1},
	}
	s.Set(ctx, serverSession)

	clientSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "client-value"},
		Updated: time.Now(),
		Vector:  VersionVector{"client1": 2},
	}

	_, err := s.Sync(ctx, clientSession, ServerWins)
	if err == nil {
		t.Error("expected conflict error")
	}
	conflictErr, ok := err.(*ConflictError)
	if !ok {
		t.Errorf("expected ConflictError, got %T", err)
	}
	if conflictErr.SessionID != "test" {
		t.Errorf("expected session ID test, got %s", conflictErr.SessionID)
	}
}

func TestSyncConflictMerge(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	serverSession := &Session{
		ID:      "test",
		Data:    map[string]any{"server-key": "server-value", "common": "unchanged"},
		Updated: time.Now(),
		Vector:  VersionVector{"client1": 1, "server": 1},
	}
	s.Set(ctx, serverSession)

	clientSession := &Session{
		ID:      "test",
		Data:    map[string]any{"client-key": "client-value", "common": "unchanged"},
		Updated: time.Now(),
		Vector:  VersionVector{"client1": 1},
	}

	result, err := s.Sync(ctx, clientSession, Merge)
	if err != nil {
		t.Errorf("Sync error: %v", err)
	}
	if result.Data["server-key"] != "server-value" {
		t.Errorf("expected server-key to be preserved, got %v", result.Data["server-key"])
	}
	if result.Data["client-key"] != "client-value" {
		t.Errorf("expected client-key to be merged, got %v", result.Data["client-key"])
	}
	if result.Data["common"] != "unchanged" {
		t.Errorf("expected common key to be unchanged, got %v", result.Data["common"])
	}
}

func TestSyncConflictClientWins(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	serverSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "server-value"},
		Updated: time.Now(),
		Vector:  VersionVector{"client1": 1, "server": 1},
	}
	s.Set(ctx, serverSession)

	clientSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "client-value"},
		Updated: time.Now().Add(-1 * time.Hour),
		Vector:  VersionVector{"client1": 1},
	}

	result, err := s.Sync(ctx, clientSession, ClientWins)
	if err != nil {
		t.Errorf("Sync error: %v", err)
	}
	if result.Data["key"] != "client-value" {
		t.Errorf("expected client-value (client wins), got %v", result.Data["key"])
	}
}

func TestVersionIncrement(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	session := &Session{ID: "test", Data: map[string]any{"key": "value1"}}
	s.Set(ctx, session)
	v1 := session.Version

	session.Data["key"] = "value2"
	s.Set(ctx, session)
	v2 := session.Version

	if v2 <= v1 {
		t.Errorf("expected version to increment, got v1=%d v2=%d", v1, v2)
	}
}

func TestVersionVectorConflict(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	serverSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "server-value"},
		Vector:  VersionVector{"client1": 2, "server": 1},
	}
	s.Set(ctx, serverSession)

	clientSession := &Session{
		ID:      "test",
		Data:    map[string]any{"key": "client-value"},
		Vector:  VersionVector{"client1": 1},
	}

	result, err := s.Sync(ctx, clientSession, LastWriteWins)
	if err != nil {
		t.Errorf("Sync error: %v", err)
	}
	if result.Vector["client1"] != 2 {
		t.Errorf("expected vector client1=2, got %d", result.Vector["client1"])
	}
}

func TestGetVersion(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	session := &Session{ID: "test", Data: map[string]any{"key": "value"}}
	s.Set(ctx, session)

	version, err := s.GetVersion(ctx, "test")
	if err != nil {
		t.Errorf("GetVersion error: %v", err)
	}
	if version != 1 {
		t.Errorf("expected version 1, got %d", version)
	}

	_, err = s.GetVersion(ctx, "nonexistent")
	if err != nil {
		t.Errorf("GetVersion for nonexistent should not error, got: %v", err)
	}
}

func TestGetVector(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	session := &Session{
		ID:     "test",
		Data:   map[string]any{"key": "value"},
		Vector: VersionVector{"client1": 1, "client2": 2},
	}
	s.Set(ctx, session)

	vector, err := s.GetVector(ctx, "test")
	if err != nil {
		t.Errorf("GetVector error: %v", err)
	}
	if vector["client1"] != 1 || vector["client2"] != 2 {
		t.Errorf("expected vector {client1:1, client2:2}, got %v", vector)
	}
}

func TestConflictError(t *testing.T) {
	err := &ConflictError{
		SessionID:     "test-session",
		ClientVersion: VersionVector{"c1": 1},
		ServerVersion: VersionVector{"c1": 2, "s1": 1},
		ClientData:    map[string]any{"key": "client"},
		ServerData:    map[string]any{"key": "server"},
	}

	expected := "conflict detected for session test-session"
	if err.Error() != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, err.Error())
	}
}