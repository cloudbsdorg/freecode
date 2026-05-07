package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/freecode/freecode/internal/session"
)

type Model struct {
	width            int
	height          int
	tabBar          *TabBarComponent
	statusBar       *StatusBar
	messageList     *MessageList
	inputArea       *InputArea
	commandPalette  *CommandPalette
	sidebar         *Sidebar
	sessionManager  *session.Manager
	quitting        bool
	focus           focusArea
	yolo            bool
	activeTabIdx    int
	tabs            []*TabState
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
		tabBar:         NewTabBar(),
		statusBar:      NewStatusBar(),
		messageList:    NewMessageList(),
		inputArea:      NewInputArea(),
		commandPalette: NewCommandPalette(),
		sidebar:        NewSidebar(),
		quitting:       false,
		focus:          focusInput,
		yolo:           false,
		activeTabIdx:   0,
		tabs:           make([]*TabState, 0),
	}

	m.tabBar.AddTab("main", "main")
	m.tabs = append(m.tabs, &TabState{ID: uuid.New().String(), Name: "main", SessionID: ""})

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
		Name:        "Split Vertical",
		Description: "Split the current view vertically",
		Keybind:     "Ctrl+Shift+V",
		Category:    "Tab",
		Handler: func() {
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Split Horizontal",
		Description: "Split the current view horizontally",
		Keybind:     "Ctrl+Shift+H",
		Category:    "Tab",
		Handler: func() {
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
		m.tabBar.SetWidth(m.width)
		m.statusBar.SetWidth(m.width)
		m.messageList.SetWidth(m.width - (m.sidebar.Width() + 4))
		m.messageList.SetHeight(m.height - 6)
		m.inputArea.SetWidth(m.width - 4)
		m.commandPalette.SetWidth(m.width / 2)
		m.commandPalette.SetHeight(m.height / 2)
		m.sidebar.SetWidth(42)
		m.sidebar.SetHeight(m.height - 3)

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

	case "ctrl+y":
		m.yolo = !m.yolo
		m.statusBar.SetYOLO(m.yolo)

	case "tab":
		m.tabBar.NextTab()
		m.activeTabIdx = m.tabBar.activeIdx

	case "shift+tab":
		m.tabBar.PrevTab()
		m.activeTabIdx = m.tabBar.activeIdx

	case "ctrl+shift+v":
	case "ctrl+shift+h":

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

	case "enter":
		if m.focus == focusSidebar {
			selected := m.sidebar.SelectedItem()
			if selected != nil {
				m.switchTabByID(selected.ID)
			}
		}

	case "g":
		m.messageList.ScrollToBottom()

	case "G":
		m.messageList.scrollToTop()
	}

	if m.focus == focusInput {
		switch msg.Type {
		case tea.KeyRunes, tea.KeyLeft, tea.KeyRight, tea.KeyHome, tea.KeyEnd,
			tea.KeyBackspace, tea.KeyDelete, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			m.inputArea.HandleKey(msg)

			if msg.Type == tea.KeyEnter && msg.Runes != nil {
				value := m.inputArea.Submit()
				if value != "" {
					m.addUserMessage(value)
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

	tabBar := m.tabBar.Render()

	var content string
	if m.sidebar.IsOpen() {
		content = m.sidebar.Render() + "\n" + m.messageList.Render()
	} else {
		content = m.messageList.Render()
	}

	content += "\n" + m.inputArea.Render()

	statusBar := m.statusBar.Render()

	paletteView := ""
	if m.commandPalette.IsOpen() {
		paletteView = m.commandPalette.Render()
	}

	result := tabBar + "\n" + content + "\n" + statusBar

	if paletteView != "" {
		result += "\n" + paletteView
	}

	return result
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
