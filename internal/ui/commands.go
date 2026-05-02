package ui

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
)

type Command struct {
	Name        string
	Description string
	Shortcut    string
	Handler     func(*Model) tea.Cmd
}

type CommandHandler struct {
	commands []Command
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands: make([]Command, 0),
	}
}

func (h *CommandHandler) Register(cmd Command) {
	h.commands = append(h.commands, cmd)
}

func (h *CommandHandler) Handle(msg tea.KeyMsg, m *Model) tea.Cmd {
	switch msg.String() {
	case "ctrl+p":
		return m.OpenCommandPalette()
	case "ctrl+b":
		return m.ToggleSidebar()
	case "ctrl+f":
		return m.ToggleFleetPanel()
	}

	return nil
}

func (h *CommandHandler) Search(query string) []Command {
	query = strings.ToLower(query)
	var results []Command

	for _, cmd := range h.commands {
		if strings.Contains(strings.ToLower(cmd.Name), query) ||
			strings.Contains(strings.ToLower(cmd.Description), query) {
			results = append(results, cmd)
		}
	}

	return results
}

func (h *CommandHandler) List() []Command {
	return h.commands
}
