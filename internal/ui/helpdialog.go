package ui

import (
	"fmt"
	"strings"

	"github.com/freecode/freecode/internal/style"

	"github.com/freecode/freecode/internal/ui/dialog"
)

type HelpDialog struct {
	width     int
	height    int
	isOpen    bool
	selected  int
	itemCount int
	colors    dialog.Colors
}

func NewHelpDialog() *HelpDialog {
	return &HelpDialog{
		width:     60,
		height:    20,
		isOpen:    false,
		selected:  0,
		itemCount: 12,
		colors:    dialog.Dark,
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
		if h.selected >= h.itemCount {
			h.selected = 0
		}
		return true
	case "k", "up":
		h.selected--
		if h.selected < 0 {
			h.selected = h.itemCount - 1
		}
		return true
	}
	return false
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
	lines = append(lines, dialog.Header("Help", h.colors))
	lines = append(lines, "")
	lines = append(lines, dialog.Muted("Keyboard Shortcuts", h.colors))
	lines = append(lines, "")

	for i, item := range helpItems {
		line := fmt.Sprintf("  %-10s %-25s %s", item.key, item.action, dialog.Muted(item.category, h.colors))
		if i == h.selected {
			line = dialog.Selected(line, h.colors)
		}
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, dialog.Muted("Press q, esc, or enter to close", h.colors))

	content := strings.Join(lines, "\n")
	return style.NewStyle().
		Background(style.Color(h.colors.Background)).
		Width(h.width).
		Height(h.height).
		Render(content)
}
