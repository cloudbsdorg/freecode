package acp

import (
	"context"
	"testing"
)

func TestNewMemoryACP(t *testing.T) {
	acp := NewMemoryACP()
	if acp == nil {
		t.Fatal("NewMemoryACP() returned nil")
	}
}

func TestCreateSession(t *testing.T) {
	acp := NewMemoryACP()
	ctx := context.Background()

	session, err := acp.CreateSession(ctx)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	if session == nil {
		t.Fatal("CreateSession() returned nil session")
	}
	if session.ID != "session-1" {
		t.Errorf("session.ID = %q, want %q", session.ID, "session-1")
	}
	if !session.Agent.Ready {
		t.Error("session.Agent.Ready = false, want true")
	}
}

func TestGetSession(t *testing.T) {
	acp := NewMemoryACP()
	ctx := context.Background()

	created, err := acp.CreateSession(ctx)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	retrieved, err := acp.GetSession(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetSession() returned nil")
	}
	if retrieved.ID != created.ID {
		t.Errorf("retrieved.ID = %q, want %q", retrieved.ID, created.ID)
	}

	nonexistent, err := acp.GetSession(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("GetSession() for nonexistent returned error: %v", err)
	}
	if nonexistent != nil {
		t.Error("GetSession() for nonexistent should return nil")
	}
}

func TestPauseSession(t *testing.T) {
	acp := NewMemoryACP()
	ctx := context.Background()

	session, _ := acp.CreateSession(ctx)
	if session.Agent.Paused {
		t.Error("session should not be paused initially")
	}

	err := acp.PauseSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("PauseSession() error = %v", err)
	}

	retrieved, _ := acp.GetSession(ctx, session.ID)
	if !retrieved.Agent.Paused {
		t.Error("session should be paused after PauseSession")
	}
}

func TestResumeSession(t *testing.T) {
	acp := NewMemoryACP()
	ctx := context.Background()

	session, _ := acp.CreateSession(ctx)
	acp.PauseSession(ctx, session.ID)

	err := acp.ResumeSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("ResumeSession() error = %v", err)
	}

	retrieved, _ := acp.GetSession(ctx, session.ID)
	if retrieved.Agent.Paused {
		t.Error("session should not be paused after ResumeSession")
	}
}

func TestStopSession(t *testing.T) {
	acp := NewMemoryACP()
	ctx := context.Background()

	session, _ := acp.CreateSession(ctx)
	if session.Agent.Stopped {
		t.Error("session should not be stopped initially")
	}

	err := acp.StopSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("StopSession() error = %v", err)
	}

	retrieved, _ := acp.GetSession(ctx, session.ID)
	if !retrieved.Agent.Stopped {
		t.Error("session should be stopped after StopSession")
	}
}

func TestSessionStateTransitions(t *testing.T) {
	acp := NewMemoryACP()
	ctx := context.Background()

	session, _ := acp.CreateSession(ctx)

	acp.PauseSession(ctx, session.ID)
	retrieved, _ := acp.GetSession(ctx, session.ID)
	if !retrieved.Agent.Paused || retrieved.Agent.Stopped {
		t.Error("after pause: Paused=true, Stopped=false")
	}

	acp.ResumeSession(ctx, session.ID)
	retrieved, _ = acp.GetSession(ctx, session.ID)
	if retrieved.Agent.Paused || retrieved.Agent.Stopped {
		t.Error("after resume: Paused=false, Stopped=false")
	}

	acp.StopSession(ctx, session.ID)
	retrieved, _ = acp.GetSession(ctx, session.ID)
	if !retrieved.Agent.Stopped {
		t.Error("after stop: Stopped=true")
	}
}