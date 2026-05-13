package ui

import (
	"strings"

	"github.com/freecode/freecode/internal/style"
)

type StatusDialog struct {
	width   int
	isOpen  bool
	model   string
	agent   string
	provider string
}

func NewStatusDialog() *StatusDialog {
	return &StatusDialog{
		width:   60,
		isOpen:  false,
		model:   "claude-opus-4-7",
		agent:   "sisyphus",
		provider: "anthropic",
	}
}

func (s *StatusDialog) SetWidth(w int) {
	s.width = w
}

func (s *StatusDialog) SetInfo(model, agent, provider string) {
	s.model = model
	s.agent = agent
	s.provider = provider
}

func (s *StatusDialog) Open() {
	s.isOpen = true
}

func (s *StatusDialog) Close() {
	s.isOpen = false
}

func (s *StatusDialog) Toggle() {
	s.isOpen = !s.isOpen
}

func (s *StatusDialog) IsOpen() bool {
	return s.isOpen
}

func (s *StatusDialog) HandleKey(msg string) bool {
	if !s.isOpen {
		return false
	}
	if msg == "escape" || msg == "enter" || msg == "q" {
		s.Close()
		return true
	}
	return false
}

func (s *StatusDialog) Render() string {
	if !s.isOpen {
		return ""
	}

	dialogStyle := style.NewStyle().
		Background(style.Color("#1E1E1E")).
		BorderStyle(style.HiddenBorder()).
		Width(s.width)

	return dialogStyle.Render(s.renderContent())
}

func (s *StatusDialog) renderContent() string {
	var lines []string

	lines = append(lines, s.renderHeader())
	lines = append(lines, "")
	lines = append(lines, s.renderSection("Model", s.model))
	lines = append(lines, s.renderSection("Provider", s.provider))
	lines = append(lines, s.renderSection("Agent", s.agent))
	lines = append(lines, "")
	lines = append(lines, s.renderHints())

	return strings.Join(lines, "\n")
}

func (s *StatusDialog) renderHeader() string {
	headerStyle := style.NewStyle().
		Foreground(style.Color("#E0E0E0")).
		Bold(true)
	return headerStyle.Render("Status")
}

func (s *StatusDialog) renderSection(title, value string) string {
	titleStyle := style.NewStyle().Foreground(style.Color("#808080"))
	valueStyle := style.NewStyle().Foreground(style.Color("#4EC9B0"))
	return "  " + titleStyle.Render(title+":") + " " + valueStyle.Render(value)
}

func (s *StatusDialog) renderHints() string {
	hintStyle := style.NewStyle().Foreground(style.Color("#808080"))
	return hintStyle.Render("press esc or enter to close")
}