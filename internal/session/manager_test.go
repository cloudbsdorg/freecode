package session

import (
	"testing"
	"time"

	"github.com/freecode/freecode/internal/config"
)

func TestNewManager(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	if m == nil {
		t.Fatal("NewManager() returned nil")
	}

	if m.sessions == nil {
		t.Error("Manager.sessions is nil")
	}

	if m.tabs == nil {
		t.Error("Manager.tabs is nil")
	}

	if m.config == nil {
		t.Error("Manager.config is nil")
	}
}

func TestManagerCreateSession(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	sess, err := m.CreateSession("test-session", "gpt-4", "sisyphus")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if sess.Title != "test-session" {
		t.Errorf("Session.Title = %q, want %q", sess.Title, "test-session")
	}

	if sess.Model != "gpt-4" {
		t.Errorf("Session.Model = %q, want %q", sess.Model, "gpt-4")
	}

	if sess.Agent != "sisyphus" {
		t.Errorf("Session.Agent = %q, want %q", sess.Agent, "sisyphus")
	}

	if sess.ID == "" {
		t.Error("Session.ID is empty")
	}

	if len(sess.Messages) != 0 {
		t.Errorf("len(Session.Messages) = %d, want 0", len(sess.Messages))
	}
}

func TestManagerGetSession(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	created, _ := m.CreateSession("get-test", "", "")

	got, ok := m.GetSession(created.ID)
	if !ok {
		t.Error("GetSession() returned false for existing session")
	}

	if got.ID != created.ID {
		t.Errorf("GetSession().ID = %q, want %q", got.ID, created.ID)
	}
}

func TestManagerGetSessionNotFound(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	_, ok := m.GetSession("nonexistent-id")
	if ok {
		t.Error("GetSession() returned true for nonexistent session")
	}
}

func TestManagerListSessions(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	m.CreateSession("sess1", "", "")
	m.CreateSession("sess2", "", "")
	m.CreateSession("sess3", "", "")

	sessions := m.ListSessions()

	if len(sessions) != 3 {
		t.Errorf("len(ListSessions()) = %d, want 3", len(sessions))
	}
}

func TestManagerDeleteSession(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	sess, _ := m.CreateSession("delete-test", "", "")

	err := m.DeleteSession(sess.ID)
	if err != nil {
		t.Fatalf("DeleteSession() error = %v", err)
	}

	_, ok := m.GetSession(sess.ID)
	if ok {
		t.Error("GetSession() returned true after DeleteSession()")
	}
}

func TestManagerDeleteSessionNotFound(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	err := m.DeleteSession("nonexistent")
	if err == nil {
		t.Error("DeleteSession() should error for nonexistent session")
	}
}

func TestManagerAddMessage(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	sess, _ := m.CreateSession("msg-test", "", "")

	msg, err := m.AddMessage(sess.ID, "user", "Hello")
	if err != nil {
		t.Fatalf("AddMessage() error = %v", err)
	}

	if msg.Role != "user" {
		t.Errorf("Message.Role = %q, want %q", msg.Role, "user")
	}

	if msg.Content != "Hello" {
		t.Errorf("Message.Content = %q, want %q", msg.Content, "Hello")
	}

	if msg.ID == "" {
		t.Error("Message.ID is empty")
	}

	got, _ := m.GetSession(sess.ID)
	if len(got.Messages) != 1 {
		t.Errorf("len(Session.Messages) = %d, want 1", len(got.Messages))
	}
}

func TestManagerAddMessageToNonexistentSession(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	_, err := m.AddMessage("nonexistent", "user", "Hello")
	if err == nil {
		t.Error("AddMessage() should error for nonexistent session")
	}
}

func TestManagerCreateTab(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	tab, err := m.CreateTab("test-tab")
	if err != nil {
		t.Fatalf("CreateTab() error = %v", err)
	}

	if tab.Name != "test-tab" {
		t.Errorf("Tab.Name = %q, want %q", tab.Name, "test-tab")
	}

	if tab.ID == "" {
		t.Error("Tab.ID is empty")
	}
}

func TestManagerGetTab(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	created, _ := m.CreateTab("get-tab")

	got, ok := m.GetTab(created.ID)
	if !ok {
		t.Error("GetTab() returned false for existing tab")
	}

	if got.ID != created.ID {
		t.Errorf("GetTab().ID = %q, want %q", got.ID, created.ID)
	}
}

func TestManagerListTabs(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	m.CreateTab("tab1")
	m.CreateTab("tab2")

	tabs := m.ListTabs()

	if len(tabs) != 2 {
		t.Errorf("len(ListTabs()) = %d, want 2", len(tabs))
	}
}

func TestManagerCloseTab(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	tab, _ := m.CreateTab("close-tab")

	err := m.CloseTab(tab.ID)
	if err != nil {
		t.Fatalf("CloseTab() error = %v", err)
	}

	_, ok := m.GetTab(tab.ID)
	if ok {
		t.Error("GetTab() returned true after CloseTab()")
	}
}

func TestSessionTimestamps(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewManager(cfg)

	before := time.Now()
	sess, _ := m.CreateSession("time-test", "", "")
	after := time.Now()

	if sess.CreatedAt.Before(before) || sess.CreatedAt.After(after) {
		t.Errorf("Session.CreatedAt outside expected range")
	}

	if sess.UpdatedAt.Before(before) || sess.UpdatedAt.After(after) {
		t.Errorf("Session.UpdatedAt outside expected range")
	}
}
