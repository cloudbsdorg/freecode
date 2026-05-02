package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m == nil {
		t.Fatal("NewModel() returned nil")
	}

	if m.tabs == nil {
		t.Error("Model.tabs is nil")
	}

	if m.active != 0 {
		t.Errorf("Model.active = %d, want 0", m.active)
	}

	if !m.yolo {
		t.Error("Model.yolo should be true by default")
	}
}

func TestModelInit(t *testing.T) {
	m := NewModel()

	cmd := m.Init()
	if cmd != nil {
		t.Error("Init() should return nil cmd")
	}
}

func TestModelView(t *testing.T) {
	m := NewModel()

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestModelViewWhenQuitting(t *testing.T) {
	m := NewModel()
	m.quitting = true

	view := m.View()
	if view != "Goodbye!\n" {
		t.Errorf("View() = %q, want %q", view, "Goodbye!\n")
	}
}

func TestModelStatusBar(t *testing.T) {
	m := NewModel()
	m.yolo = true
	m.active = 0

	status := m.statusBar()
	if status == "" {
		t.Error("statusBar() returned empty string")
	}
}

func TestModelTabBar(t *testing.T) {
	m := NewModel()

	bar := m.tabBar()
	if bar == "" {
		t.Error("tabBar() returned empty string")
	}
}

func TestModelSessionContent(t *testing.T) {
	m := NewModel()

	content := m.sessionContent()
	if content == "" {
		t.Error("sessionContent() returned empty string")
	}
}

func TestModelOpenCommandPalette(t *testing.T) {
	m := NewModel()
	m.commandPaletteOpen = false

	cmd := m.OpenCommandPalette()
	if cmd != nil {
		t.Error("OpenCommandPalette() should return nil cmd")
	}

	if !m.commandPaletteOpen {
		t.Error("commandPaletteOpen should be true after OpenCommandPalette()")
	}
}

func TestModelToggleSidebar(t *testing.T) {
	m := NewModel()
	m.sidebarOpen = false

	cmd := m.ToggleSidebar()
	if cmd != nil {
		t.Error("ToggleSidebar() should return nil cmd")
	}

	if !m.sidebarOpen {
		t.Error("sidebarOpen should be true after ToggleSidebar()")
	}
}

func TestModelToggleFleetPanel(t *testing.T) {
	m := NewModel()
	m.fleetPanelOpen = false

	cmd := m.ToggleFleetPanel()
	if cmd != nil {
		t.Error("ToggleFleetPanel() should return nil cmd")
	}

	if !m.fleetPanelOpen {
		t.Error("fleetPanelOpen should be true after ToggleFleetPanel()")
	}
}

func TestModelUpdateWindowSize(t *testing.T) {
	m := NewModel()

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update with WindowSizeMsg should return nil cmd")
	}

	if m.width != 100 {
		t.Errorf("m.width = %d, want 100", m.width)
	}

	if m.height != 50 {
		t.Errorf("m.height = %d, want 50", m.height)
	}
}

func TestModelUpdateToggleYolo(t *testing.T) {
	m := NewModel()
	m.yolo = false

	msg := ToggleYolo{}
	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update with ToggleYolo should return nil cmd")
	}

	if !m.yolo {
		t.Error("m.yolo should be true after ToggleYolo")
	}
}

func TestModelUpdateQuit(t *testing.T) {
	m := NewModel()
	m.quitting = false

	msg := quitMsg{}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("Update with quitMsg should return tea.Quit cmd")
	}

	if !m.quitting {
		t.Error("m.quitting should be true after quitMsg")
	}
}

func TestModelUpdateNewTab(t *testing.T) {
	m := NewModel()

	msg := NewTabMsg{}
	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update with NewTabMsg should return nil cmd")
	}
}

func TestModelUpdateCloseTab(t *testing.T) {
	m := NewModel()

	msg := CloseTabMsg{Tab: 0}
	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update with CloseTabMsg should return nil cmd")
	}
}

func TestModelUpdateSwitchTab(t *testing.T) {
	m := NewModel()

	msg := SwitchTabMsg{Tab: 1}
	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update with SwitchTabMsg should return nil cmd")
	}
}

func TestModelUpdateUnknownMsg(t *testing.T) {
	m := NewModel()

	msg := "unknown"
	_, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Update with unknown msg should return nil cmd")
	}
}

func TestModelHandleKeyCtrlC(t *testing.T) {
	m := NewModel()

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.handleKey(msg)

	if cmd == nil {
		t.Error("handleKey with ctrl+c should return tea.Quit cmd")
	}
}

func TestModelHandleKeyQ(t *testing.T) {
	m := NewModel()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.handleKey(msg)

	if cmd == nil {
		t.Error("handleKey with q should return tea.Quit cmd")
	}
}

func TestModelHandleKeyCtrlT(t *testing.T) {
	m := NewModel()

	msg := tea.KeyMsg{Type: tea.KeyCtrlT}
	_, cmd := m.handleKey(msg)

	if cmd != nil {
		t.Error("handleKey with ctrl+t should return nil cmd")
	}
}

func TestModelHandleKeyCtrlW(t *testing.T) {
	m := NewModel()

	msg := tea.KeyMsg{Type: tea.KeyCtrlW}
	_, cmd := m.handleKey(msg)

	if cmd != nil {
		t.Error("handleKey with ctrl+w should return nil cmd")
	}
}

func TestModelHandleKeyTab(t *testing.T) {
	m := NewModel()
	m.tabs = []TabModel{{ID: "1"}, {ID: "2"}}

	msg := tea.KeyMsg{Type: tea.KeyTab}
	_, cmd := m.handleKey(msg)

	if cmd != nil {
		t.Error("handleKey with tab should return nil cmd")
	}
}

func TestModelHandleKeyShiftTab(t *testing.T) {
	m := NewModel()
	m.tabs = []TabModel{{ID: "1"}, {ID: "2"}}

	msg := tea.KeyMsg{Type: tea.KeyShiftTab}
	_, cmd := m.handleKey(msg)

	if cmd != nil {
		t.Error("handleKey with shift+tab should return nil cmd")
	}
}

func TestModelHandleKeyCtrlY(t *testing.T) {
	m := NewModel()
	m.yolo = false

	msg := tea.KeyMsg{Type: tea.KeyCtrlY}
	_, cmd := m.handleKey(msg)

	if cmd != nil {
		t.Error("handleKey with ctrl+y should return nil cmd")
	}

	if !m.yolo {
		t.Error("m.yolo should be toggled to true")
	}
}

func TestModelHandleKeyUnknown(t *testing.T) {
	m := NewModel()
	oldYolo := m.yolo

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	_, cmd := m.handleKey(msg)

	if cmd != nil {
		t.Error("handleKey with unknown key should return nil cmd")
	}

	if m.yolo != oldYolo {
		t.Error("m.yolo should not change for unknown key")
	}
}

func TestModelSwitchTab(t *testing.T) {
	m := NewModel()
	m.tabs = []TabModel{{ID: "1"}, {ID: "2"}}
	m.active = 0

	_, _ = m.switchTab(1)

	if m.active != 1 {
		t.Errorf("m.active = %d, want 1", m.active)
	}
}

func TestModelSwitchTabInvalidIndex(t *testing.T) {
	m := NewModel()
	m.tabs = []TabModel{{ID: "1"}, {ID: "2"}}
	m.active = 0

	_, _ = m.switchTab(5)

	if m.active != 0 {
		t.Errorf("m.active = %d, want 0 (should not change for invalid index)", m.active)
	}
}

func TestModelSwitchTabNegativeIndex(t *testing.T) {
	m := NewModel()
	m.tabs = []TabModel{{ID: "1"}, {ID: "2"}}
	m.active = 0

	_, _ = m.switchTab(-1)

	if m.active != 0 {
		t.Errorf("m.active = %d, want 0 (should not change for negative index)", m.active)
	}
}

func TestModelAddTab(t *testing.T) {
	m := NewModel()

	_, _ = m.addTab()

	if len(m.tabs) != 0 {
		t.Errorf("len(m.tabs) = %d, want 0 (addTab is a stub)", len(m.tabs))
	}
}

func TestModelCloseTab(t *testing.T) {
	m := NewModel()

	_, _ = m.closeTab(0)

	if len(m.tabs) != 0 {
		t.Errorf("len(m.tabs) = %d, want 0 (closeTab is a stub)", len(m.tabs))
	}
}
