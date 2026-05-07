package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbletea"
)

var PaletteContainerStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#1E1E1E")).
	BorderStyle(lipgloss.HiddenBorder())

var PaletteHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#007ACC")).
	Bold(true).
	Padding(0, 1)

var PaletteItemStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#E0E0E0"))

var PaletteSelectedStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#007ACC")).
	Foreground(lipgloss.Color("#FFFFFF"))

var PaletteDescriptionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#808080"))

type PaletteCommand struct {
	Name        string
	Description string
	Keybind     string
	Category    string
	Handler     func()
}

type CommandPalette struct {
	commands    []PaletteCommand
	filtered    []PaletteCommand
	query       string
	selectedIdx int
	width       int
	height      int
	isOpen      bool
}

func NewCommandPalette() *CommandPalette {
	return &CommandPalette{
		commands:    make([]PaletteCommand, 0),
		filtered:    make([]PaletteCommand, 0),
		query:       "",
		selectedIdx: 0,
		width:       60,
		height:      15,
		isOpen:      false,
	}
}

func (c *CommandPalette) Register(cmd PaletteCommand) {
	c.commands = append(c.commands, cmd)
}

func (c *CommandPalette) Open() {
	c.isOpen = true
	c.query = ""
	c.filtered = c.commands
	c.selectedIdx = 0
}

func (c *CommandPalette) Close() {
	c.isOpen = false
	c.query = ""
}

func (c *CommandPalette) Toggle() {
	if c.isOpen {
		c.Close()
	} else {
		c.Open()
	}
}

func (c *CommandPalette) IsOpen() bool {
	return c.isOpen
}

func (c *CommandPalette) HandleKey(msg tea.KeyMsg) bool {
	if !c.isOpen {
		return false
	}

	switch msg.Type {
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			c.query += string(r)
		}
		c.filterCommands()
		c.selectedIdx = 0

	case tea.KeyBackspace:
		if len(c.query) > 0 {
			c.query = c.query[:len(c.query)-1]
			c.filterCommands()
			c.selectedIdx = 0
		}

	case tea.KeyUp:
		if c.selectedIdx > 0 {
			c.selectedIdx--
		}

	case tea.KeyDown:
		if c.selectedIdx < len(c.filtered)-1 {
			c.selectedIdx++
		}

	case tea.KeyEnter:
		if c.selectedIdx >= 0 && c.selectedIdx < len(c.filtered) {
			cmd := c.filtered[c.selectedIdx]
			cmd.Handler()
			c.Close()
			return true
		}
	}
	return false
}

func (c *CommandPalette) filterCommands() {
	c.filtered = make([]PaletteCommand, 0)
	query := strings.ToLower(c.query)

	for _, cmd := range c.commands {
		nameLower := strings.ToLower(cmd.Name)
		descLower := strings.ToLower(cmd.Description)
		catLower := strings.ToLower(cmd.Category)

		if strings.Contains(nameLower, query) ||
			strings.Contains(descLower, query) ||
			strings.Contains(catLower, query) {
			c.filtered = append(c.filtered, cmd)
		}
	}
}

func (c *CommandPalette) SetWidth(w int) {
	c.width = w
}

func (c *CommandPalette) SetHeight(h int) {
	c.height = h
}

func (c *CommandPalette) Render() string {
	if !c.isOpen {
		return ""
	}

	header := PaletteHeaderStyle.Render("Command Palette")
	if c.query != "" {
		header += " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render(c.query)
	}

	lines := []string{header, ""}

	maxItems := c.height - 4
	if maxItems < 1 {
		maxItems = 10
	}

	displayItems := c.filtered
	if len(displayItems) > maxItems {
		displayItems = displayItems[:maxItems]
	}

	categories := make(map[string][]PaletteCommand)
	for _, cmd := range displayItems {
		cat := cmd.Category
		if cat == "" {
			cat = "Other"
		}
		categories[cat] = append(categories[cat], cmd)
	}

	idx := 0
	for cat, cmds := range categories {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#606060")).Render(cat))
		for _, cmd := range cmds {
			itemStr := "  " + cmd.Name
			if cmd.Description != "" {
				itemStr += " - " + cmd.Description
			}
			if cmd.Keybind != "" {
				itemStr += " (" + cmd.Keybind + ")"
			}

			if idx == c.selectedIdx {
				lines = append(lines, PaletteSelectedStyle.Render(itemStr))
			} else {
				lines = append(lines, PaletteItemStyle.Render(itemStr))
			}
			idx++
		}
		lines = append(lines, "")
	}

	if len(c.filtered) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render("  No commands found"))
	}

	result := strings.Join(lines, "\n")

	if len(result) > c.width {
		lines := strings.Split(result, "\n")
		for i, line := range lines {
			if len(line) > c.width {
				lines[i] = line[:c.width-3] + "..."
			}
		}
		result = strings.Join(lines, "\n")
	}

	return PaletteContainerStyle.Render(result)
}
