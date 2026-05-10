package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/freecode/freecode/internal/tool"
)

type ToolToggleHandler func(name string, enabled bool)

type ToolDialog struct {
	width    int
	isOpen  bool
	tools   []tool.ToolInfo
	selected int
	onToggle ToolToggleHandler
}

func NewToolDialog() *ToolDialog {
	return &ToolDialog{
		width:    70,
		isOpen:  false,
		tools:   []tool.ToolInfo{},
		selected: 0,
	}
}

func (m *ToolDialog) SetToggleHandler(h ToolToggleHandler) {
	m.onToggle = h
}

func (m *ToolDialog) SetWidth(w int) {
	m.width = w
}

func (m *ToolDialog) SetTools(tools []tool.ToolInfo) {
	m.tools = tools
	if m.selected >= len(tools) {
		m.selected = 0
	}
}

func (m *ToolDialog) Open() {
	m.isOpen = true
	m.selected = 0
}

func (m *ToolDialog) Close() {
	m.isOpen = false
}

func (m *ToolDialog) Toggle() {
	m.isOpen = !m.isOpen
}

func (m *ToolDialog) IsOpen() bool {
	return m.isOpen
}

func (m *ToolDialog) HandleKey(msg string) bool {
	if !m.isOpen {
		return false
	}

	switch msg {
	case "escape":
		m.Close()
		return true
	case "enter":
		m.toggleSelected()
		return true
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return true
	case "down", "j":
		if m.selected < len(m.tools)-1 {
			m.selected++
		}
		return true
	}
	return false
}

func (m *ToolDialog) toggleSelected() {
	if m.selected >= 0 && m.selected < len(m.tools) {
		m.tools[m.selected].Enabled = !m.tools[m.selected].Enabled
		if m.onToggle != nil {
			m.onToggle(m.tools[m.selected].Name, m.tools[m.selected].Enabled)
		}
	}
}

func (m *ToolDialog) Render() string {
	if !m.isOpen {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E1E")).
		Border(lipgloss.HiddenBorder()).
		Width(m.width)

	return dialogStyle.Render(m.renderContent())
}

func (m *ToolDialog) renderContent() string {
	var lines []string

	lines = append(lines, m.renderHeader())
	lines = append(lines, "")
	lines = append(lines, m.renderTools()...)
	lines = append(lines, "")
	lines = append(lines, m.renderHints())

	return strings.Join(lines, "\n")
}

func (m *ToolDialog) renderHeader() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Bold(true).
		Render("Tools Configuration")
}

func (m *ToolDialog) renderTools() []string {
	var lines []string

	if len(m.tools) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		lines = append(lines, "  "+emptyStyle.Render("No tools registered"))
		return lines
	}

	for i, t := range m.tools {
		lines = append(lines, m.renderTool(t, i))
	}

	return lines
}

func (m *ToolDialog) renderTool(t tool.ToolInfo, idx int) string {
	selected := idx == m.selected
	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	}

	var statusStyle lipgloss.Style
	status := "○"
	if t.Enabled {
		status = "●"
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))
	} else {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	}

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))

	line := prefix + " " + statusStyle.Render(status) + " " + nameStyle.Render(t.Name)
	if t.Description != "" {
		line += " " + descStyle.Render("— "+t.Description)
	}

	return line
}

func (m *ToolDialog) renderHints() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	return hintStyle.Render("↑↓ select  enter toggle  esc close")
}

func (m *ToolDialog) GetTools() []tool.ToolInfo {
	return m.tools
}