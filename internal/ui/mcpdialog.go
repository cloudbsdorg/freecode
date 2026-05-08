package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type MCPServer struct {
	Name    string
	URL     string
	Enabled bool
}

type MCPDialog struct {
	width    int
	isOpen  bool
	servers []MCPServer
	selected int
}

func NewMCPDialog() *MCPDialog {
	return &MCPDialog{
		width:    60,
		isOpen:  false,
		servers: []MCPServer{},
		selected: 0,
	}
}

func (m *MCPDialog) SetWidth(w int) {
	m.width = w
}

func (m *MCPDialog) SetServers(servers []MCPServer) {
	m.servers = servers
	if m.selected >= len(servers) {
		m.selected = 0
	}
}

func (m *MCPDialog) Open() {
	m.isOpen = true
	m.selected = 0
}

func (m *MCPDialog) Close() {
	m.isOpen = false
}

func (m *MCPDialog) Toggle() {
	m.isOpen = !m.isOpen
}

func (m *MCPDialog) IsOpen() bool {
	return m.isOpen
}

func (m *MCPDialog) HandleKey(msg string) bool {
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
		if m.selected < len(m.servers)-1 {
			m.selected++
		}
		return true
	}
	return false
}

func (m *MCPDialog) toggleSelected() {
	if m.selected >= 0 && m.selected < len(m.servers) {
		m.servers[m.selected].Enabled = !m.servers[m.selected].Enabled
	}
}

func (m *MCPDialog) Render() string {
	if !m.isOpen {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E1E")).
		Border(lipgloss.HiddenBorder()).
		Width(m.width)

	return dialogStyle.Render(m.renderContent())
}

func (m *MCPDialog) renderContent() string {
	var lines []string

	lines = append(lines, m.renderHeader())
	lines = append(lines, "")
	lines = append(lines, m.renderServers()...)
	lines = append(lines, "")
	lines = append(lines, m.renderHints())

	return strings.Join(lines, "\n")
}

func (m *MCPDialog) renderHeader() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Bold(true).
		Render("MCP Servers")
}

func (m *MCPDialog) renderServers() []string {
	var lines []string

	if len(m.servers) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		lines = append(lines, "  "+emptyStyle.Render("No MCP servers configured"))
		lines = append(lines, "")
		lines = append(lines, "  "+emptyStyle.Render("Configure servers in ~/.config/freecode/mcp.json"))
		return lines
	}

	for i, server := range m.servers {
		lines = append(lines, m.renderServer(server, i))
	}

	return lines
}

func (m *MCPDialog) renderServer(server MCPServer, idx int) string {
	selected := idx == m.selected
	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	}

	var statusStyle lipgloss.Style
	status := "○"
	if server.Enabled {
		status = "●"
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))
	} else {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	}

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))

	line := prefix + " " + statusStyle.Render(status) + " " + nameStyle.Render(server.Name)
	if server.URL != "" {
		line += " " + urlStyle.Render("("+server.URL+")")
	}

	return line
}

func (m *MCPDialog) renderHints() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	return hintStyle.Render("↑↓ select  enter toggle  esc close")
}