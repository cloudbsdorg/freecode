package ui

import (
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
