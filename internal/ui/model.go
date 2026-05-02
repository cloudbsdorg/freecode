package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/freecode/freecode/internal/session"
)

type Model struct {
	width             int
	height            int
	tabs              []TabModel
	active            int
	session           *session.Session
	yolo              bool
	quitting          bool
	commandPaletteOpen bool
	sidebarOpen       bool
	fleetPanelOpen    bool
}

type TabModel struct {
	ID      string
	Name    string
	Session *session.Session
}

func NewModel() *Model {
	return &Model{
		tabs:   make([]TabModel, 0),
		active: 0,
		yolo:   true,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case NewTabMsg:
		return m.addTab()

	case CloseTabMsg:
		return m.closeTab(msg.Tab)

	case SwitchTabMsg:
		return m.switchTab(msg.Tab)

	case ToggleYolo:
		m.yolo = !m.yolo
		return m, nil

	case quitMsg:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	s := m.tabBar()
	s += "\n"
	s += m.sessionContent()
	s += m.statusBar()

	return s
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "ctrl+t":
		return m.addTab()
	case "ctrl+w":
		return m.closeTab(m.active)
	case "tab":
		return m.switchTab((m.active + 1) % len(m.tabs))
	case "shift+tab":
		return m.switchTab((m.active - 1 + len(m.tabs)) % len(m.tabs))
	case "ctrl+y":
		m.yolo = !m.yolo
		return m, nil
	}
	return m, nil
}

func (m *Model) addTab() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) closeTab(idx int) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) switchTab(idx int) (tea.Model, tea.Cmd) {
	if idx >= 0 && idx < len(m.tabs) {
		m.active = idx
	}
	return m, nil
}

func (m *Model) tabBar() string {
	return "[Tab 1] [Tab 2] [+]"
}

func (m *Model) sessionContent() string {
	return "Session content here..."
}

func (m *Model) statusBar() string {
	yoloStatus := "YOLO: ON"
	if !m.yolo {
		yoloStatus = "YOLO: OFF"
	}
	return "\n" + yoloStatus + " | Tab: " + string(rune(m.active+'1'))
}

func (m *Model) OpenCommandPalette() tea.Cmd {
	m.commandPaletteOpen = !m.commandPaletteOpen
	return nil
}

func (m *Model) ToggleSidebar() tea.Cmd {
	m.sidebarOpen = !m.sidebarOpen
	return nil
}

func (m *Model) ToggleFleetPanel() tea.Cmd {
	m.fleetPanelOpen = !m.fleetPanelOpen
	return nil
}

type NewTabMsg struct{}

type CloseTabMsg struct {
	Tab int
}

type SwitchTabMsg struct {
	Tab int
}

type ToggleYolo struct{}

type quitMsg struct{}
