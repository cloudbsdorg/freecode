package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var StatusBarStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#2D2D2D")).
	Foreground(lipgloss.Color("#808080")).
	Height(1).
	Width(100).
	Padding(0, 1)

var StatusBarActiveStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#4EC9B0"))

var StatusBarWarningStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFCC00"))

var StatusBarErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#F44747"))

type StatusBar struct {
	width       int
	model       string
	agent       string
	provider    string
	yoloEnabled bool
	tabCount    int
	connectionStatus string
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		width:            80,
		model:           "claude",
		agent:           "sisyphus",
		provider:        "anthropic",
		yoloEnabled:     false,
		tabCount:        1,
		connectionStatus: "connected",
	}
}

func (s *StatusBar) SetWidth(w int) {
	s.width = w
}

func (s *StatusBar) SetModel(model string) {
	s.model = model
}

func (s *StatusBar) SetAgent(agent string) {
	s.agent = agent
}

func (s *StatusBar) SetProvider(provider string) {
	s.provider = provider
}

func (s *StatusBar) SetYOLO(enabled bool) {
	s.yoloEnabled = enabled
}

func (s *StatusBar) SetTabCount(count int) {
	s.tabCount = count
}

func (s *StatusBar) SetConnectionStatus(status string) {
	s.connectionStatus = status
}

func (s *StatusBar) Render() string {
	parts := []string{}

	parts = append(parts, s.connectionIndicator())
	parts = append(parts, s.yoloIndicator())
	parts = append(parts, s.modelIndicator())
	parts = append(parts, s.agentIndicator())
	parts = append(parts, s.tabIndicator())

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += "  "
		}
		result += part
	}

	totalLen := len(result)
	if totalLen > s.width && s.width > 0 {
		parts = parts[:4]
		result = ""
		for i, part := range parts {
			if i > 0 {
				result += "  "
			}
			result += part
		}
	}

	return StatusBarStyle.Render(result)
}

func (s *StatusBar) connectionIndicator() string {
	status := s.connectionStatus
	var style lipgloss.Style
	switch status {
	case "connected":
		style = StatusBarActiveStyle
	case "connecting":
		style = StatusBarWarningStyle
	default:
		style = StatusBarErrorStyle
	}
	return style.Render(fmt.Sprintf("● %s", status))
}

func (s *StatusBar) yoloIndicator() string {
	if s.yoloEnabled {
		return StatusBarWarningStyle.Render("YOLO: ON")
	}
	return "YOLO: OFF"
}

func (s *StatusBar) modelIndicator() string {
	return fmt.Sprintf("Model: %s", s.model)
}

func (s *StatusBar) agentIndicator() string {
	return fmt.Sprintf("Agent: %s", s.agent)
}

func (s *StatusBar) tabIndicator() string {
	return fmt.Sprintf("Tabs: %d", s.tabCount)
}
