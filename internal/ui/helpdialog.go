package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type HelpDialog struct {
	width     int
	height   int
	isOpen   bool
	selected int
}

func NewHelpDialog() *HelpDialog {
	return &HelpDialog{
		width:     60,
		height:    20,
		isOpen:    false,
		selected:  0,
	}
}

func (h *HelpDialog) Open() {
	h.isOpen = true
	h.selected = 0
}

func (h *HelpDialog) Close() {
	h.isOpen = false
}

func (h *HelpDialog) IsOpen() bool {
	return h.isOpen
}

func (h *HelpDialog) Toggle() {
	if h.isOpen {
		h.Close()
	} else {
		h.Open()
	}
}

func (h *HelpDialog) HandleKey(key string) bool {
	if !h.isOpen {
		return false
	}

	switch key {
	case "q", "escape", "enter":
		h.Close()
		return true
	case "j", "down":
		h.selected++
		if h.selected >= h.itemCount() {
			h.selected = 0
		}
		return true
	case "k", "up":
		h.selected--
		if h.selected < 0 {
			h.selected = h.itemCount() - 1
		}
		return true
	}
	return false
}

func (h *HelpDialog) itemCount() int {
	return 12
}

func (h *HelpDialog) SetWidth(w int) {
	h.width = w
}

func (h *HelpDialog) SetHeight(height int) {
	h.height = height
}

func (h *HelpDialog) Render() string {
	if !h.isOpen {
		return ""
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E1E")).
		Border(lipgloss.HiddenBorder()).
		Width(h.width).
		Height(h.height)

	helpItems := []struct {
		key      string
		action   string
		category string
	}{
		{"Ctrl+P", "Command Palette", "General"},
		{"Ctrl+H", "Go to Home", "Navigation"},
		{"Ctrl+Q", "Quit", "General"},
		{"Ctrl+T", "New Tab", "Tab"},
		{"Ctrl+W", "Close Tab", "Tab"},
		{"Ctrl+B", "Toggle Sidebar", "View"},
		{"Ctrl+Y", "Toggle YOLO Mode", "General"},
		{"Enter", "Submit / Start Session", "Navigation"},
		{"Tab", "Next Tab", "Tab"},
		{"Shift+Tab", "Previous Tab", "Tab"},
		{"j/k", "Navigate", "Navigation"},
		{"g/G", "Scroll to Top/Bottom", "Navigation"},
	}

	var lines []string
	lines = append(lines, headerStyle.Render("Help"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#606060")).Render("Keyboard Shortcuts"))
	lines = append(lines, "")

	for i, item := range helpItems {
		catStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))

		line := fmt.Sprintf("  %-10s %-25s %s", item.key, item.action, catStyle.Render(item.category))
		if i == h.selected {
			line = lipgloss.NewStyle().
				Background(lipgloss.Color("#007ACC")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Render(line)
		}
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#606060")).Render("Press q, esc, or enter to close"))

	content := strings.Join(lines, "\n")
	return dialogStyle.Render(content)
}