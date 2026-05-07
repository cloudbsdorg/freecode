package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/google/uuid"
)

type Route string

const (
	RouteHome    Route = "home"
	RouteSession Route = "session"
)

type Model struct {
	width           int
	height         int
	route           Route
	tabBar          *TabBarComponent
	statusBar       *StatusBar
	messageList     *MessageList
	inputArea      *InputArea
	commandPalette *CommandPalette
	sidebar        *Sidebar
	quitting       bool
	focus          focusArea
	yolo           bool
	activeTabIdx   int
	tabs           []*TabState
	banner         string
}

func getBanner() string {
	f := figure.NewFigure("FREECODE", "cosmike", true)
	return f.String()
}

type focusArea int

const (
	focusInput focusArea = iota
	focusPalette
	focusSidebar
)

type TabState struct {
	ID        string
	Name      string
	SessionID string
}

func NewModel() *Model {
	m := &Model{
		width:           80,
		height:         24,
		route:           RouteHome,
		tabBar:          NewTabBar(),
		statusBar:       NewStatusBar(),
		messageList:     NewMessageList(),
		inputArea:       NewInputArea(),
		commandPalette:  NewCommandPalette(),
		sidebar:         NewSidebar(),
		quitting:        false,
		focus:           focusInput,
		yolo:           false,
		activeTabIdx:   0,
		tabs:           make([]*TabState, 0),
		banner:         getBanner(),
	}

	m.tabBar.AddTab("main", "main")
	m.tabs = append(m.tabs, &TabState{ID: uuid.New().String(), Name: "main", SessionID: ""})
	m.statusBar.SetTabCount(1)

	m.registerCommands()

	return m
}

func (m *Model) registerCommands() {
	m.commandPalette.Register(PaletteCommand{
		Name:        "New Tab",
		Description: "Create a new tab",
		Keybind:     "Ctrl+T",
		Category:    "Tab",
		Handler: func() {
			m.addTab()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Close Tab",
		Description: "Close the current tab",
		Keybind:     "Ctrl+W",
		Category:    "Tab",
		Handler: func() {
			m.closeCurrentTab()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Toggle Sidebar",
		Description: "Show or hide the session sidebar",
		Keybind:     "Ctrl+B",
		Category:    "View",
		Handler: func() {
			m.sidebar.Toggle()
			m.updateLayout()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Toggle YOLO Mode",
		Description: "Skip all confirmations",
		Keybind:     "Ctrl+Y",
		Category:    "General",
		Handler: func() {
			m.yolo = !m.yolo
			m.statusBar.SetYOLO(m.yolo)
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Command Palette",
		Description: "Open command palette",
		Keybind:     "Ctrl+P",
		Category:    "General",
		Handler: func() {
			m.commandPalette.Toggle()
			if m.commandPalette.IsOpen() {
				m.focus = focusPalette
			} else {
				m.focus = focusInput
			}
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Go to Home",
		Description: "Return to home screen",
		Keybind:     "Ctrl+H",
		Category:    "Navigation",
		Handler: func() {
			m.route = RouteHome
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Start Session",
		Description: "Start a new session",
		Keybind:     "Enter",
		Category:    "Navigation",
		Handler: func() {
			m.route = RouteSession
		},
	})
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()

	case tea.KeyMsg:
		return m.handleKey(msg)

	case ToggleYolo:
		m.yolo = !m.yolo
		m.statusBar.SetYOLO(m.yolo)

	case NewTabMsg:
		m.addTab()

	case CloseTabMsg:
		m.closeTab(msg.Tab)

	case SwitchTabMsg:
		m.switchTab(msg.Tab)

	case quitMsg:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) updateLayout() {
	sidebarWidth := 0
	if m.sidebar.IsOpen() {
		sidebarWidth = 42
		if sidebarWidth > m.width/4 {
			sidebarWidth = m.width / 4
		}
	}

	contentWidth := m.width - sidebarWidth - 4
	if contentWidth < 40 {
		contentWidth = m.width - 4
		sidebarWidth = 0
	}

	m.tabBar.SetWidth(m.width)
	m.statusBar.SetWidth(m.width)
	m.messageList.SetWidth(contentWidth)
	m.messageList.SetHeight(m.height - 6)
	m.inputArea.SetWidth(contentWidth)
	m.commandPalette.SetWidth(m.width / 2)
	m.commandPalette.SetHeight(m.height / 2)
	m.sidebar.SetWidth(sidebarWidth)
	m.sidebar.SetHeight(m.height - 3)
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.commandPalette.IsOpen() {
		handled := m.commandPalette.HandleKey(msg)
		if handled {
			return m, nil
		}
		if msg.Type == tea.KeyEsc {
			m.commandPalette.Close()
			m.focus = focusInput
			return m, nil
		}
	}

	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit

	case "ctrl+t":
		m.addTab()

	case "ctrl+w":
		m.closeCurrentTab()

	case "ctrl+p":
		m.commandPalette.Toggle()
		if m.commandPalette.IsOpen() {
			m.focus = focusPalette
		} else {
			m.focus = focusInput
		}

	case "ctrl+b":
		m.sidebar.Toggle()
		m.updateLayout()

	case "ctrl+y":
		m.yolo = !m.yolo
		m.statusBar.SetYOLO(m.yolo)

	case "ctrl+h":
		m.route = RouteHome

	case "tab":
		m.tabBar.NextTab()
		m.activeTabIdx = m.tabBar.activeIdx

	case "shift+tab":
		m.tabBar.PrevTab()
		m.activeTabIdx = m.tabBar.activeIdx

	case "j", "down":
		if m.focus == focusSidebar {
			m.sidebar.SelectNext()
		} else if m.focus == focusPalette {
			m.commandPalette.HandleKey(msg)
		} else {
			m.messageList.ScrollDown()
		}

	case "k", "up":
		if m.focus == focusSidebar {
			m.sidebar.SelectPrev()
		} else if m.focus == focusPalette {
			m.commandPalette.HandleKey(msg)
		} else {
			m.messageList.ScrollUp()
		}

	case "g":
		m.messageList.ScrollToBottom()

	case "G":
		m.messageList.scrollToTop()

	case "enter":
		if m.focus == focusSidebar {
			selected := m.sidebar.SelectedItem()
			if selected != nil {
				m.route = RouteSession
				m.switchTabByID(selected.ID)
			}
		} else if m.route == RouteHome {
			value := m.inputArea.Submit()
			if value != "" {
				m.route = RouteSession
				m.addUserMessage(value)
			}
		}
	}

	if m.focus == focusInput {
		switch msg.Type {
		case tea.KeyRunes, tea.KeyLeft, tea.KeyRight, tea.KeyHome, tea.KeyEnd,
			tea.KeyBackspace, tea.KeyDelete, tea.KeyUp, tea.KeyDown:
			m.inputArea.HandleKey(msg)

			if msg.Type == tea.KeyEnter && len(msg.Runes) > 0 {
				value := m.inputArea.Submit()
				if value != "" {
					m.addUserMessage(value)
					m.route = RouteSession
				}
			}
		}
	}

	return m, nil
}

func (m *Model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4EC9B0")).
			Render("Goodbye!\n")
	}

	switch m.route {
	case RouteHome:
		return m.renderHome()
	case RouteSession:
		return m.renderSession()
	default:
		return m.renderHome()
	}
}

func (m *Model) renderHome() string {
	bannerLines := strings.Split(m.banner, "\n")
	bannerHeight := len(bannerLines)
	bannerWidth := 0
	for _, line := range bannerLines {
		if len(line) > bannerWidth {
			bannerWidth = len(line)
		}
	}

	startY := (m.height - bannerHeight - 5) / 2
	if startY < 0 {
		startY = 0
	}

	var s strings.Builder

	for i := 0; i < startY; i++ {
		s.WriteString("\n")
	}

	for _, line := range bannerLines {
		padding := (m.width - len(line)) / 2
		if padding < 0 {
			padding = 0
		}
		s.WriteString(strings.Repeat(" ", padding))
		s.WriteString(line)
		s.WriteString("\n")
	}

	s.WriteString("\n")

	inputWidth := m.width / 2
	if inputWidth < 40 {
		inputWidth = 40
	}
	inputPadding := (m.width - inputWidth) / 2
	inputContent := m.inputArea.Render()
	for _, line := range strings.Split(inputContent, "\n") {
		s.WriteString(strings.Repeat(" ", inputPadding))
		s.WriteString(line)
		s.WriteString("\n")
	}

	s.WriteString("\n")

	hintY := m.height - 3
	currentY := startY + bannerHeight + 5
	for currentY < hintY {
		s.WriteString("\n")
		currentY++
	}

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#606060")).
		Align(lipgloss.Center)

	hintText := "Ctrl+P: Command Palette  |  Ctrl+B: Toggle Sidebar  |  Ctrl+H: Home  |  Ctrl+Q: Quit"
	hintPadding := (m.width - len(hintText)) / 2
	if hintPadding > 0 {
		hintText = strings.Repeat(" ", hintPadding) + hintText
	}
	s.WriteString(hintStyle.Render(hintText))
	s.WriteString("\n")

	s.WriteString(m.statusBar.Render())

	paletteView := ""
	if m.commandPalette.IsOpen() {
		paletteView = "\n" + m.commandPalette.Render()
	}

	return s.String() + paletteView
}

func (m *Model) renderSession() string {
	tabBar := m.tabBar.Render()

	var content string

	if m.sidebar.IsOpen() {
		sidebarWidth := m.sidebar.Width()
		if sidebarWidth > 0 {
			sidebarRender := m.sidebar.Render()
			content = sidebarRender + " " + m.renderSessionContent()
		} else {
			content = m.renderSessionContent()
		}
	} else {
		content = m.renderSessionContent()
	}

	input := m.inputArea.Render()
	status := m.statusBar.Render()

	paletteView := ""
	if m.commandPalette.IsOpen() {
		paletteView = "\n" + m.commandPalette.Render()
	}

	return tabBar + "\n" + content + "\n" + input + "\n" + status + paletteView
}

func (m *Model) renderSessionContent() string {
	sidebarWidth := 0
	if m.sidebar.IsOpen() {
		sidebarWidth = m.sidebar.Width()
	}
	contentWidth := m.width - sidebarWidth - 4

	msgs := m.messageList.Render()
	if len(msgs) > contentWidth {
		lines := strings.Split(msgs, "\n")
		for i, line := range lines {
			if len(line) > contentWidth {
				lines[i] = line[:contentWidth-3] + "..."
			}
		}
		msgs = strings.Join(lines, "\n")
	}

	return msgs
}

func (m *Model) addTab() {
	id := uuid.New().String()
	name := fmt.Sprintf("tab-%d", len(m.tabs)+1)
	m.tabBar.AddTab(id, name)
	m.tabs = append(m.tabs, &TabState{ID: id, Name: name, SessionID: ""})
	m.activeTabIdx = len(m.tabs) - 1
	m.statusBar.SetTabCount(len(m.tabs))
}

func (m *Model) closeCurrentTab() {
	if len(m.tabs) <= 1 {
		return
	}
	m.tabBar.CloseTab(m.activeTabIdx)
	m.tabs = append(m.tabs[:m.activeTabIdx], m.tabs[m.activeTabIdx+1:]...)
	if m.activeTabIdx >= len(m.tabs) {
		m.activeTabIdx = len(m.tabs) - 1
	}
	m.statusBar.SetTabCount(len(m.tabs))
}

func (m *Model) closeTab(idx int) {
	if len(m.tabs) <= 1 || idx < 0 || idx >= len(m.tabs) {
		return
	}
	m.tabBar.CloseTab(idx)
	m.tabs = append(m.tabs[:idx], m.tabs[idx+1:]...)
	if m.activeTabIdx >= len(m.tabs) {
		m.activeTabIdx = len(m.tabs) - 1
	}
	m.statusBar.SetTabCount(len(m.tabs))
}

func (m *Model) switchTab(idx int) {
	if idx >= 0 && idx < len(m.tabs) {
		m.activeTabIdx = idx
		m.tabBar.SetActive(idx)
	}
}

func (m *Model) switchTabByID(id string) {
	for i, t := range m.tabs {
		if t.ID == id {
			m.switchTab(i)
			break
		}
	}
}

func (m *Model) addUserMessage(content string) {
	msg := Message{
		ID:        uuid.New().String(),
		Role:      "user",
		Content:   content,
		Timestamp: time.Now(),
	}
	m.messageList.AddMessage(msg)
	m.messageList.ScrollToBottom()
}

func (m *Model) addAssistantMessage(content string) {
	msg := Message{
		ID:        uuid.New().String(),
		Role:      "assistant",
		Content:   content,
		Timestamp: time.Now(),
	}
	m.messageList.AddMessage(msg)
	m.messageList.ScrollToBottom()
}

func (m *Model) SetMessages(msgs []Message) {
	m.messageList.SetMessages(msgs)
}

func (m *Model) SetModel(model string) {
	m.statusBar.SetModel(model)
}

func (m *Model) SetAgent(agent string) {
	m.statusBar.SetAgent(agent)
}

func (m *Model) SetProvider(provider string) {
	m.statusBar.SetProvider(provider)
}

func (m *Model) SetSessions(sessions []SidebarItem) {
	m.sidebar.SetItems(sessions)
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
