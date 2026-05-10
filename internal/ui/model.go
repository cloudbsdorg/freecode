package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/freecode/freecode/internal/agent"
	"github.com/freecode/freecode/internal/args"
	"github.com/freecode/freecode/internal/config"
	"github.com/freecode/freecode/internal/session"
	"github.com/freecode/freecode/internal/tool"
	"github.com/google/uuid"
)

type Route string

const (
	RouteHome    Route = "home"
	RouteSession Route = "session"
	RouteSetup   Route = "setup"
)

type Model struct {
	width             int
	height            int
	route             Route
	tabBar            *TabBarComponent
	statusBar         *StatusBar
	messageList       *MessageList
	inputArea         *InputArea
	commandPalette    *CommandPalette
	sidebar           *Sidebar
	toastManager      *ToastManager
	soundManager      *SoundManager
	animationManager  *AnimationManager
	helpDialog        *HelpDialog
	permissionDialog  *PermissionDialog
	questionDialog    *QuestionDialog
	selectDialog      *SelectDialog
	statusDialog      *StatusDialog
	exportDialog      *ExportDialog
	mcpDialog         *MCPDialog
	toolDialog        *ToolDialog
	consolePanel      *ConsolePanel
	autocompleteDialog *AutocompleteDialog
	fleetPanel        *FleetPanel
	setupDialog       *SetupDialog
	fleetTicking      bool
	quitting          bool
	focus             focusArea
	yolo              bool
	activeTabIdx      int
	tabs              []*TabState
	banner            string
	cliArgs           args.Args
	promptSubmitted   bool
	sessions          []*session.Session
	engine            *agent.Engine
	currentSession    *session.Session
	theme             Theme
	themeName         string
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

func NewModel(args args.Args) *Model {
	cfg := config.DefaultConfig()
	engine := agent.NewEngine(cfg)
	agent.RegisterBuiltinAgents(engine)

	m := &Model{
		width:             80,
		height:            24,
		route:             RouteHome,
		tabBar:            NewTabBar(),
		statusBar:         NewStatusBar(),
		messageList:       NewMessageList(),
		inputArea:         NewInputArea(),
		commandPalette:    NewCommandPalette(),
		sidebar:           NewSidebar(),
		toastManager:      NewToastManager(),
		soundManager:      NewSoundManager(),
		animationManager:  NewAnimationManager(),
		helpDialog:        NewHelpDialog(),
		permissionDialog:   NewPermissionDialog(),
		questionDialog:     NewQuestionDialog(),
		selectDialog:       NewSelectDialog(),
		statusDialog:       NewStatusDialog(),
		exportDialog:       NewExportDialog(),
		mcpDialog:          NewMCPDialog(),
		toolDialog:         NewToolDialog(),
		consolePanel:       NewConsolePanel(),
		autocompleteDialog: NewAutocompleteDialog(),
		fleetPanel:         NewFleetPanel(),
		setupDialog:        NewSetupDialog(),
		quitting:          false,
		focus:             focusInput,
		yolo:              false,
		activeTabIdx:      0,
		tabs:             make([]*TabState, 0),
		banner:           GetBanner(),
		cliArgs:          args,
		promptSubmitted:   false,
		engine:           engine,
		theme:            DarkTheme(),
		themeName:        "dark",
	}

	m.tabBar.AddTab("main", "main")
	m.tabs = append(m.tabs, &TabState{ID: uuid.New().String(), Name: "main", SessionID: ""})
	m.statusBar.SetTabCount(1)

	if args.Agent != "" {
		m.statusBar.SetAgent(args.Agent)
	}
	if args.Model != "" {
		m.statusBar.SetModel(args.Model)
	}
	if args.Continue || args.SessionID != "" {
		m.route = RouteSession
	}

	paths := config.PathsGet()
	configPath := paths.ConfigFile("config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		m.route = RouteSetup
	}

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
		Name:        "Toggle Animation",
		Description: fmt.Sprintf("Cycle animation level (current: %s)", m.animationManager.Level().String()),
		Keybind:     "Ctrl+Shift+A",
		Category:    "View",
		Handler: func() {
			m.animationManager.Toggle()
			m.toastManager.ShowInfo("Animation: " + m.animationManager.Level().String())
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

	m.commandPalette.Register(PaletteCommand{
		Name:        "Switch Theme",
		Description: fmt.Sprintf("Change UI theme (current: %s)", m.themeName),
		Keybind:     "",
		Category:    "View",
		Handler: func() {
			themes := ListThemes()
			if len(themes) == 0 {
				return
			}
			currentIdx := 0
			for i, t := range themes {
				if t == m.themeName {
					currentIdx = i
					break
				}
			}
			nextIdx := (currentIdx + 1) % len(themes)
			m.themeName = themes[nextIdx]
			m.theme = GetTheme(m.themeName)
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Switch Model",
		Description: fmt.Sprintf("Change AI model (current: %s)", m.cliArgs.Model),
		Keybind:     "",
		Category:    "Agent",
		Handler: func() {
			models := m.getAvailableModels()
			if len(models) == 0 {
				return
			}
			currentIdx := 0
			for i, model := range models {
				if model == m.cliArgs.Model {
					currentIdx = i
					break
				}
			}
			nextIdx := (currentIdx + 1) % len(models)
			m.cliArgs.Model = models[nextIdx]
			m.statusBar.SetModel(m.cliArgs.Model)
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Rename Session",
		Description: "Rename the current session",
		Keybind:     "",
		Category:    "Session",
		Handler: func() {
			if m.currentSession == nil {
				return
			}
			m.currentSession.Title = "Renamed Session"
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Fork Session",
		Description: "Create a copy of the current session",
		Keybind:     "",
		Category:    "Session",
		Handler: func() {
			if m.currentSession == nil {
				return
			}
			agentName := m.cliArgs.Agent
			if agentName == "" {
				agentName = "sisyphus"
			}
			model := m.cliArgs.Model
			if model == "" {
				model = "claude-opus-4-7"
			}
			newSess, _ := m.engine.SessionManager().CreateSession(m.currentSession.Title+" (fork)", model, agentName)
			for _, msg := range m.currentSession.Messages {
				m.engine.SessionManager().AddMessage(newSess.ID, msg.Role, msg.Content)
			}
			m.currentSession = newSess
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Undo",
		Description: "Remove last assistant message",
		Keybind:     "",
		Category:    "Session",
		Handler: func() {
			if m.currentSession == nil || len(m.currentSession.Messages) == 0 {
				return
			}
			m.messageList.RemoveLastAssistant()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Copy Transcript",
		Description: "Copy session transcript to clipboard",
		Keybind:     "",
		Category:    "Session",
		Handler: func() {
			if m.currentSession == nil {
				return
			}
			var transcript strings.Builder
			transcript.WriteString("# Session Transcript\n\n")
			transcript.WriteString(fmt.Sprintf("Title: %s\n", m.currentSession.Title))
			transcript.WriteString(fmt.Sprintf("Model: %s\n", m.currentSession.Model))
			transcript.WriteString(fmt.Sprintf("Agent: %s\n\n", m.currentSession.Agent))
			for _, msg := range m.currentSession.Messages {
				transcript.WriteString(fmt.Sprintf("## %s\n\n%s\n\n", strings.ToUpper(msg.Role), msg.Content))
			}
			_ = transcript.String()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Help",
		Description: "Show keyboard shortcuts and help",
		Keybind:     "",
		Category:    "System",
		Handler: func() {
			m.helpDialog.Open()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Show Status",
		Description: "Show current model/agent/provider status",
		Keybind:     "",
		Category:    "System",
		Handler: func() {
			m.statusDialog.SetInfo(m.cliArgs.Model, m.cliArgs.Agent, "anthropic")
			m.statusDialog.Open()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Export Session",
		Description: "Export session transcript",
		Keybind:     "",
		Category:    "Session",
		Handler: func() {
			title := "Session"
			if m.currentSession != nil {
				title = m.currentSession.Title
			}
			var msgs []string
			if m.currentSession != nil {
				for _, msg := range m.currentSession.Messages {
					msgs = append(msgs, msg.Content)
				}
			}
			m.exportDialog.SetSession(title, msgs)
			m.exportDialog.Open()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "MCP Servers",
		Description: "Manage MCP server connections",
		Keybind:     "",
		Category:    "System",
		Handler: func() {
			m.mcpDialog.Open()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Tools Configuration",
		Description: "Enable or disable built-in tools",
		Keybind:     "",
		Category:    "System",
		Handler: func() {
			m.toolDialog.SetToggleHandler(func(name string, enabled bool) {
				if enabled {
					m.engine.EnableTool(name)
				} else {
					m.engine.DisableTool(name)
				}
				m.engine.SaveToolStates()
			})
			m.toolDialog.SetTools(tool.ListToolsWithStatus())
			m.toolDialog.Open()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Toggle Console",
		Description: "Show or hide debug console",
		Keybind:     "",
		Category:    "System",
		Handler: func() {
			m.consolePanel.Toggle()
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Select Option",
		Description: "Open selection dialog",
		Keybind:     "",
		Category:    "System",
		Handler: func() {
			m.selectDialog.SetOptions([]SelectOption{
				{Title: "Option 1", Value: "opt1", Description: "First option"},
				{Title: "Option 2", Value: "opt2", Description: "Second option"},
			})
			m.selectDialog.SetOnSelect(func(val string) {
				m.toastManager.ShowInfo("Selected: " + val)
			})
		},
	})

	m.commandPalette.Register(PaletteCommand{
		Name:        "Fleet Panel",
		Description: "Show fleet management panel",
		Keybind:     "",
		Category:    "System",
		Handler: func() {
			m.fleetPanel.Toggle()
		},
	})
}

func fleetTickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(5 * time.Second)
		return fleetTickMsg{}
	}
}

func (m *Model) Init() tea.Cmd {
	m.setTerminalTitle("Freecode")
	return func() tea.Msg {
		time.Sleep(100 * time.Millisecond)
		return initTick{}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()

	case initTick:
		m.handleInit()

	case fleetTickMsg:
		if m.fleetPanel.IsOpen() {
			m.fleetPanel.refresh()
			return m, tea.Batch(fleetTickCmd())
		}
		m.fleetTicking = false

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

	case PermissionRequestMsg:
		m.permissionDialog.SetRequest(msg.Request)

	case QuestionRequestMsg:
		m.questionDialog.SetRequest(msg.Request)

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
	m.permissionDialog.SetWidth(m.width / 2)
	m.questionDialog.SetWidth(m.width / 2)
	m.selectDialog.SetWidth(m.width / 2)
	m.statusDialog.SetWidth(m.width / 2)
	m.exportDialog.SetWidth(m.width / 2)
	m.mcpDialog.SetWidth(m.width / 2)
	m.toolDialog.SetWidth(m.width / 2)
	m.consolePanel.SetWidth(m.width - 4)
	m.consolePanel.SetHeight(m.height - 4)
	m.autocompleteDialog.SetWidth(m.width / 2)
	m.fleetPanel.SetWidth(m.width / 2)
	m.fleetPanel.SetHeight(m.height / 2)
	m.setupDialog.SetWidth(m.width / 2)
}

func (m *Model) handleSetupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.setupDialog.IsOpen() {
		return m, nil
	}

	switch msg.Type {
	case tea.KeyUp, tea.KeyShiftTab:
		m.setupDialog.MoveUp()
	case tea.KeyDown, tea.KeyTab:
		m.setupDialog.MoveDown()
	case tea.KeyEnter:
		prevStep := m.setupDialog.GetStep()
		m.setupDialog.Next()
		if m.setupDialog.GetStep() == SetupStepModel && prevStep == SetupStepProvider {
			m.setupDialog.SetModels(m.getModelsForProvider(m.setupDialog.GetSelectedProviderID()))
		}
		if m.setupDialog.GetStep() == SetupStepDone {
			providerID, modelID, apiKey := m.setupDialog.GetSelection()
			m.saveSetupConfig(providerID, modelID, apiKey)
			m.statusBar.SetProvider(providerID)
			m.statusBar.SetModel(modelID)
			m.route = RouteHome
		}
	case tea.KeyEsc:
		m.setupDialog.Prev()
		if !m.setupDialog.IsOpen() {
			m.quitting = true
			return m, tea.Quit
		}
	case tea.KeyBackspace:
		m.setupDialog.BackspaceAPIKey()
	default:
		if msg.Runes != nil {
			for _, r := range msg.Runes {
				m.setupDialog.AppendToAPIKey(r)
			}
		}
	}

	return m, nil
}

func (m *Model) getModelsForProvider(providerID string) []string {
	defaultModels := map[string][]string{
		"ollama":       {"llama3.2", "mistral", "codellama", "qwen2.5-coder"},
		"lmstudio":     {"llama-3.2-3b-instruct", "mistral-7b-instruct", "codellama-7b"},
		"minimax":      {"MiniMax-Text-01", "abab6.5s-chat"},
		"openai":       {"gpt-4o", "gpt-4o-mini", "gpt-4-turbo"},
		"anthropic":    {"claude-sonnet-4-20250514", "claude-3-5-sonnet-20241022", "claude-3-5-haiku-20241022"},
	}
	if models, ok := defaultModels[providerID]; ok {
		return models
	}
	return []string{"gpt-4o", "claude-sonnet-4-20250514"}
}

func (m *Model) saveSetupConfig(providerID, modelID, apiKey string) {
	paths := config.PathsGet()

	cfg := config.DefaultConfig()
	cfg.Session.Dir = paths.SessionDir()

	baseURL := m.inferBaseURL(providerID)

	cfg.Models = map[string]config.ModelConfig{
		"default": {
			Provider: providerID,
			Name:     modelID,
		},
	}

	cfg.Providers = map[string]config.ProviderConfig{
		providerID: {
			APIKey:  apiKey,
			BaseURL: baseURL,
		},
	}

	config.SaveConfig(paths.ConfigFile("config.yaml"), cfg)
}

func (m *Model) inferBaseURL(providerID string) string {
	guesses := map[string]string{
		"ollama":       "http://localhost:11434",
		"lmstudio":     "http://localhost:1234",
		"ollama-cloud": "https://ollama.cloud",
		"minimax":      "https://api.minimax.chat/v1",
		"minimax-cn":    "https://api.minimax.chat/v1",
	}
	if url, ok := guesses[providerID]; ok {
		return url
	}
	return ""
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.helpDialog.IsOpen() {
		handled := m.helpDialog.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

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

	if m.permissionDialog.IsVisible() {
		return m.handlePermissionKey(msg)
	}

	if m.questionDialog.IsVisible() {
		return m.handleQuestionKey(msg)
	}

	if m.selectDialog.IsVisible() {
		handled := m.selectDialog.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.statusDialog.IsOpen() {
		handled := m.statusDialog.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.exportDialog.IsOpen() {
		handled := m.exportDialog.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.mcpDialog.IsOpen() {
		handled := m.mcpDialog.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.toolDialog.IsOpen() {
		handled := m.toolDialog.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.consolePanel.IsOpen() {
		handled := m.consolePanel.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.autocompleteDialog.IsVisible() {
		handled := m.autocompleteDialog.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.fleetPanel.IsOpen() {
		handled := m.fleetPanel.HandleKey(msg.String())
		if handled {
			return m, nil
		}
	}

	if m.route == RouteSetup {
		return m.handleSetupKey(msg)
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

	case "ctrl+shift+a":
		m.animationManager.Toggle()
		m.toastManager.ShowInfo("Animation: " + m.animationManager.Level().String())

	case "?":
		m.helpDialog.Toggle()

	case "ctrl+h":
		m.route = RouteHome
		m.setTerminalTitle("Freecode")

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
				if selected.Title != "" && selected.Title != "New Session" {
					title := selected.Title
					if len(title) > 40 {
						title = title[:37] + "..."
					}
					m.setTerminalTitle("FC | " + title)
				} else {
					m.setTerminalTitle("Freecode")
				}
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

	if m.fleetPanel.IsOpen() && !m.fleetTicking {
		m.fleetTicking = true
		return m, fleetTickCmd()
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
	case RouteSetup:
		return m.renderSetup()
	default:
		return m.renderHome()
	}
}

func (m *Model) renderToast() string {
	return m.toastManager.Render()
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
		paletteView = m.commandPalette.RenderCentered(m.width, m.height)
	}

	return s.String() + paletteView
}

func (m *Model) renderSetup() string {
	dialog := m.setupDialog.Render()
	padding := (m.width - 70) / 2
	if padding < 0 {
		padding = 0
	}
	lines := strings.Split(dialog, "\n")
	var s strings.Builder
	startY := (m.height - 20) / 2
	if startY < 0 {
		startY = 0
	}
	for i := 0; i < startY; i++ {
		s.WriteString("\n")
	}
	for _, line := range lines {
		s.WriteString(strings.Repeat(" ", padding))
		s.WriteString(line)
		s.WriteString("\n")
	}
	return s.String()
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
	toast := m.renderToast()
	help := m.helpDialog.Render()
	permission := m.permissionDialog.Render()
	question := m.questionDialog.Render()
	selectView := m.selectDialog.Render()
	statusView := m.statusDialog.Render()
	exportView := m.exportDialog.Render()
	mcpView := m.mcpDialog.Render()
	toolView := m.toolDialog.Render()
	consoleView := m.consolePanel.Render()
	autocompleteView := m.autocompleteDialog.Render()
	fleetView := m.fleetPanel.Render()

	paletteView := ""
	if m.commandPalette.IsOpen() {
		paletteView = "\n" + m.commandPalette.Render()
	}

	return tabBar + "\n" + content + "\n" + input + "\n" + status + toast + help + permission + question + selectView + statusView + exportView + mcpView + toolView + consoleView + autocompleteView + fleetView + paletteView
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

	if m.currentSession == nil {
		agentName := m.cliArgs.Agent
		if agentName == "" {
			agentName = "sisyphus"
		}
		model := m.cliArgs.Model
		if model == "" {
			model = "claude-opus-4-7"
		}
		sess, _ := m.engine.SessionManager().CreateSession("New Session", model, agentName)
		m.currentSession = sess
		m.tabs[m.activeTabIdx].SessionID = sess.ID
		m.setTerminalTitle("Freecode")
	} else if m.currentSession.Title != "" && m.currentSession.Title != "New Session" {
		title := m.currentSession.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		m.setTerminalTitle("FC | " + title)
	}

	m.engine.SessionManager().AddMessage(m.currentSession.ID, "user", content)

	agentName := m.cliArgs.Agent
	if agentName == "" {
		agentName = "sisyphus"
	}
	model := m.cliArgs.Model
	if model == "" {
		model = "claude-opus-4-7"
	}

	resp, err := m.engine.Run(context.Background(), agent.Request{
		SessionID: m.currentSession.ID,
		AgentName: agentName,
		Model:     model,
		Message: agent.Message{
			Role:    "user",
			Content: content,
		},
	})
	if err != nil {
		m.addAssistantMessage("Error: "+err.Error(), nil)
		return
	}
	var parts []MessagePart
	for _, p := range resp.Message.Parts {
		parts = append(parts, MessagePart{
			Type:    p.Type,
			Content: p.Content,
			Tool:    p.Tool,
		})
	}
	m.addAssistantMessage(resp.Message.Content, parts)
}

func (m *Model) addAssistantMessage(content string, parts []MessagePart) {
	msg := Message{
		ID:        uuid.New().String(),
		Role:      "assistant",
		Content:   content,
		Timestamp: time.Now(),
		Parts:     parts,
	}
	m.messageList.AddMessage(msg)
	m.messageList.ScrollToBottom()
	m.soundManager.Play(SoundEventMessageReceived)
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

type initTick struct{}

type fleetTickMsg struct{}

type PermissionRequestMsg struct {
	Request *PermissionRequest
}

type PermissionResponseMsg struct {
	Response *PermissionResponse
}

type ShowPermissionMsg struct {
	Request *PermissionRequest
}

type QuestionRequestMsg struct {
	Request *QuestionRequest
}

type ShowQuestionMsg struct {
	Request *QuestionRequest
}

func (m *Model) handleInit() {
	m.loadSessions()

	if m.route == RouteSetup {
		m.populateSetupProviders()
		return
	}

	if m.cliArgs.Prompt != "" && !m.promptSubmitted {
		m.promptSubmitted = true
		m.inputArea.SetValue(m.cliArgs.Prompt)
		m.route = RouteSession
		m.setTerminalTitle("Freecode")
		m.addUserMessage(m.cliArgs.Prompt)
		return
	}

	if m.cliArgs.SessionID != "" {
		m.route = RouteSession
		m.loadSessionMessages(m.cliArgs.SessionID)
		return
	}

	if m.cliArgs.Continue && len(m.sessions) > 0 {
		m.route = RouteSession
		m.loadSessionMessages(m.sessions[0].ID)
	}
}

func (m *Model) populateSetupProviders() {
	providers := []ProviderInfo{
		{ID: "ollama", Name: "Ollama (Local)"},
		{ID: "lmstudio", Name: "LM Studio (Local)"},
		{ID: "minimax", Name: "Minimax"},
		{ID: "openai", Name: "OpenAI"},
		{ID: "anthropic", Name: "Anthropic"},
	}
	m.setupDialog.SetProviders(providers)
}

func (m *Model) setTerminalTitle(title string) {
	fmt.Print("\x1b]2;" + title + "\x1b\\")
}

func (m *Model) handlePermissionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	stage := m.permissionDialog.GetStage()

	switch msg.String() {
	case "escape":
		if stage == "always" {
			m.permissionDialog.SetStage("permission")
		} else if stage == "reject" {
			m.permissionDialog.SetStage("permission")
		} else {
			m.permissionDialog.Clear()
		}
		return m, nil

	case "enter":
		if stage == "permission" {
			m.permissionDialog.SetStage("always")
		} else if stage == "always" {
			m.permissionDialog.Clear()
			m.toastManager.ShowInfo("Permission remembered for this session")
		} else if stage == "reject" {
			m.permissionDialog.Clear()
			m.toastManager.ShowWarning("Permission denied")
		}
		return m, nil

	case "left", "h":
		if stage == "permission" {
		}
		return m, nil

	case "right", "l":
		if stage == "permission" {
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) handleQuestionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.questionDialog.IsEditing() {
		switch msg.String() {
		case "escape":
			m.questionDialog.SetEditing(false)
			return m, nil
		case "enter":
			m.questionDialog.SetEditing(false)
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "escape":
		m.questionDialog.Clear()
		return m, nil

	case "enter":
		if m.questionDialog.IsConfirm() {
			answers := m.questionDialog.Submit()
			m.questionDialog.Clear()
			m.toastManager.ShowInfo("Question answered")
			_ = answers
		} else {
			m.questionDialog.ToggleCurrentOption()
		}
		return m, nil

	case "left", "h":
		m.questionDialog.PrevTab()
		return m, nil

	case "right", "l":
		m.questionDialog.NextTab()
		return m, nil

	case "up", "k":
		m.questionDialog.PrevOption()
		return m, nil

	case "down", "j":
		m.questionDialog.NextOption()
		return m, nil

	case "shift+tab":
		m.questionDialog.PrevTab()
		return m, nil

	case "tab":
		m.questionDialog.NextTab()
		return m, nil
	}

	return m, nil
}

func (m *Model) ShowPermissionRequest(req *PermissionRequest) {
	m.permissionDialog.SetRequest(req)
}

func (m *Model) loadSessions() {
	store := session.NewStore(filepath.Join(m.configDir(), "sessions"))
	sessions, err := store.ListSessions()
	if err != nil || sessions == nil {
		return
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})

	m.sessions = sessions

	items := make([]SidebarItem, 0, len(sessions))
	for _, sess := range sessions {
		item := SidebarItem{
			ID:        sess.ID,
			Title:     sess.Title,
			Timestamp: formatTimestamp(sess.UpdatedAt),
		}
		if item.Title == "" {
			item.Title = "Untitled session"
		}
		items = append(items, item)
	}

	m.sidebar.SetItems(items)
}

func (m *Model) loadSessionMessages(sessionID string) {
	store := session.NewStore(filepath.Join(m.configDir(), "sessions"))
	sess, err := store.LoadSession(sessionID)
	if err != nil || sess == nil {
		return
	}

	for _, msg := range sess.Messages {
		m.addUserMessage(msg.Content)
	}
}

func (m *Model) configDir() string {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = os.Getenv("HOME")
	}
	if homeDir == "" {
		return ".freecode"
	}
	return filepath.Join(homeDir, ".config", "freecode")
}

func formatTimestamp(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
	return t.Format("Jan 2")
}

func (m *Model) getAvailableModels() []string {
	var models []string

	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		models = append(models, "anthropic/claude-opus-4-7", "anthropic/claude-sonnet-4-5", "anthropic/claude-haiku-4-7")
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		models = append(models, "openai/gpt-4o", "openai/gpt-4o-mini", "openai/gpt-4-turbo", "openai/o1-preview", "openai/o1-mini")
	}
	if os.Getenv("MINIMAX_API_KEY") != "" {
		models = append(models, "minimax/minimax-coding-plan")
	}
	if os.Getenv("GROQ_API_KEY") != "" {
		models = append(models, "groq/llama-3.3-70b-versatile", "groq/mixtral-8x7b-32768")
	}
	if os.Getenv("PERPLEXITY_API_KEY") != "" {
		models = append(models, "perplexity/llama-3.1-sonar-large-128k-online", "perplexity/llama-3.1-sonar-huge-128k-online")
	}
	if os.Getenv("GOOGLE_API_KEY") != "" {
		models = append(models, "google/gemini-2.5-pro", "google/gemini-2.5-flash")
	}
	if os.Getenv("DEEPSEEK_API_KEY") != "" {
		models = append(models, "deepseek/deepseek-coder", "deepseek/deepseek-chat")
	}
	if os.Getenv("OLLAMA_BASE_URL") != "" || os.Getenv("OLLAMA_API_KEY") != "" {
		models = append(models, "ollama/llama3.2", "ollama/codellama", "ollama/mistral")
	}
	if os.Getenv("OPENROUTER_API_KEY") != "" {
		models = append(models, "openrouter/claude-3.5-sonnet", "openrouter/gpt-4", "openrouter/gemini-2.0-flash")
	}

	if len(models) == 0 {
		models = []string{
			"anthropic/claude-opus-4-7",
			"openai/gpt-4o",
			"minimax/minimax-coding-plan",
			"groq/llama-3.3-70b-versatile",
			"deepseek/deepseek-coder",
			"ollama/llama3.2",
		}
	}

	return models
}
