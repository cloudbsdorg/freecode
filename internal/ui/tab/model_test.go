package tab

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
)

func TestNewTabState(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.tabs == nil {
		t.Error("tabs is nil")
	}
	if m.active != 0 {
		t.Errorf("active = %d, want 0", m.active)
	}
}

func TestTabStateAddTab(t *testing.T) {
	m := New()
	tab := m.AddTab("Test Tab")

	if tab.Name != "Test Tab" {
		t.Errorf("Name = %q, want %q", tab.Name, "Test Tab")
	}
	if tab.ID == "" {
		t.Error("ID is empty")
	}
}

func TestTabStateCloseTab(t *testing.T) {
	m := New()
	tab := m.AddTab("Test Tab")

	result := m.CloseTab(tab.ID)
	if !result {
		t.Error("CloseTab() returned false")
	}

	if len(m.List()) != 0 {
		t.Errorf("len(List()) = %d, want 0", len(m.List()))
	}
}

func TestTabStateCloseTabNotFound(t *testing.T) {
	m := New()
	m.AddTab("Test Tab")

	result := m.CloseTab("nonexistent-id")
	if result {
		t.Error("CloseTab() should return false for nonexistent tab")
	}
}

func TestTabStateSetActive(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")
	m.AddTab("Tab 2")

	m.SetActive(1)
	if m.active != 1 {
		t.Errorf("active = %d, want 1", m.active)
	}

	m.SetActive(-1)
	if m.active != 1 {
		t.Errorf("active = %d, want 1 (should not change for negative)", m.active)
	}
}

func TestTabStateGetActive(t *testing.T) {
	m := New()
	tab := m.AddTab("Active Tab")

	active := m.GetActive()
	if active.ID != tab.ID {
		t.Errorf("GetActive().ID = %q, want %q", active.ID, tab.ID)
	}
}

func TestTabStateGetActiveEmpty(t *testing.T) {
	m := New()
	active := m.GetActive()
	if active.ID != "" {
		t.Errorf("GetActive().ID = %q, want empty", active.ID)
	}
}

func TestTabStateList(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")
	m.AddTab("Tab 2")

	tabs := m.List()
	if len(tabs) != 2 {
		t.Errorf("len(List()) = %d, want 2", len(tabs))
	}
}

func TestTabStateAddSession(t *testing.T) {
	m := New()
	tab := m.AddTab("Test Tab")

	result := m.AddSession(tab.ID, "session-1")
	if !result {
		t.Error("AddSession() returned false")
	}

	active := m.GetActive()
	if len(active.Sessions) != 1 {
		t.Errorf("len(Sessions) = %d, want 1", len(active.Sessions))
	}
}

func TestTabStateAddSessionNotFound(t *testing.T) {
	m := New()
	result := m.AddSession("nonexistent", "session-1")
	if result {
		t.Error("AddSession() should return false for nonexistent tab")
	}
}

func TestTabStateRemoveSession(t *testing.T) {
	m := New()
	tab := m.AddTab("Test Tab")
	m.AddSession(tab.ID, "session-1")

	result := m.RemoveSession(tab.ID, "session-1")
	if !result {
		t.Error("RemoveSession() returned false")
	}
}

func TestTabStateSetSplit(t *testing.T) {
	m := New()
	tab := m.AddTab("Test Tab")

	result := m.SetSplit(tab.ID, true, 0.5)
	if !result {
		t.Error("SetSplit() returned false")
	}

	active := m.GetActive()
	if !active.SplitVertical {
		t.Error("SplitVertical should be true")
	}
}

func TestTabStateSetSplitNotFound(t *testing.T) {
	m := New()
	result := m.SetSplit("nonexistent", true, 0.5)
	if result {
		t.Error("SetSplit() should return false for nonexistent tab")
	}
}

func TestGenerateID(t *testing.T) {
	id := generateID()
	if id == "" {
		t.Error("generateID() returned empty")
	}
	if len(id) < 10 {
		t.Errorf("generateID() length = %d, want > 10", len(id))
	}
}

func TestRandomString(t *testing.T) {
	s := randomString(8)
	if len(s) != 8 {
		t.Errorf("len(randomString(8)) = %d, want 8", len(s))
	}
}

func TestNewKeyHandler(t *testing.T) {
	h := NewKeyHandler()
	if h == nil {
		t.Fatal("NewKeyHandler() returned nil")
	}
}

func TestKeyHandlerHandleCtrlT(t *testing.T) {
	h := NewKeyHandler()
	m := New()

	_, cmd := h.Handle(tea.KeyMsg{}, m)
	if cmd != nil {
		t.Error("Handle() should return nil cmd for unknown key")
	}
}

func TestTabStateView(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestTabStateViewEmpty(t *testing.T) {
	m := New()

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string for empty tab state")
	}
}

func TestTabStateTabBar(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")
	m.AddTab("Tab 2")

	bar := m.tabBar()
	if bar == "" {
		t.Error("tabBar() returned empty string")
	}
}

func TestTabStateTabBarEmpty(t *testing.T) {
	m := New()

	bar := m.tabBar()
	if bar == "" {
		t.Error("tabBar() returned empty string")
	}
}

func TestTabStateContent(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")
	m.AddSession(m.List()[0].ID, "session-1")

	content := m.content()
	if content == "" {
		t.Error("content() returned empty string")
	}
}

func TestTabStateContentNoSessions(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")

	content := m.content()
	if content != "(no session)" {
		t.Errorf("content() = %q, want %q", content, "(no session)")
	}
}

func TestTabStateStatusLine(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")
	m.AddSession(m.List()[0].ID, "session-1")

	status := m.statusLine()
	if status == "" {
		t.Error("statusLine() returned empty string")
	}
}

func TestTabStateSetSize(t *testing.T) {
	m := New()
	m.SetSize(100, 50)

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}

func TestTabStateStatusLineHSplit(t *testing.T) {
	m := New()
	m.AddTab("Tab 1")
	m.SetSplit(m.List()[0].ID, false, 0.5)

	status := m.statusLine()
	if status == "" {
		t.Error("statusLine() returned empty string")
	}
}
