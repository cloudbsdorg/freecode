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
			return true
		}

	case tea.KeyDown:
		if c.selectedIdx < len(c.filtered)-1 {
			c.selectedIdx++
			return true
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

	startIdx := 0
	endIdx := maxItems
	if endIdx > len(c.filtered) {
		endIdx = len(c.filtered)
	}
	if startIdx > c.selectedIdx {
		startIdx = c.selectedIdx
		endIdx = startIdx + maxItems
		if endIdx > len(c.filtered) {
			endIdx = len(c.filtered)
		}
	}
	if c.selectedIdx >= endIdx {
		endIdx = c.selectedIdx + 1
		startIdx = endIdx - maxItems
		if startIdx < 0 {
			startIdx = 0
		}
	}

	displayItems := c.filtered[startIdx:endIdx]

	seenCats := make(map[string]bool)
	for i, cmd := range displayItems {
		cat := cmd.Category
		if cat == "" {
			cat = "Other"
		}
		if !seenCats[cat] {
			lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#606060")).Render(cat))
			seenCats[cat] = true
		}

		itemStr := "  " + cmd.Name
		if cmd.Description != "" {
			itemStr += " - " + cmd.Description
		}
		if cmd.Keybind != "" {
			itemStr += " (" + cmd.Keybind + ")"
		}

		if startIdx+i == c.selectedIdx {
			lines = append(lines, PaletteSelectedStyle.Render(itemStr))
		} else {
			lines = append(lines, PaletteItemStyle.Render(itemStr))
		}
	}

	if len(c.filtered) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render("  No commands found"))
	}

	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Width(c.width).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3B3B3B")).
		Background(lipgloss.Color("#1E1E1E")).
		MarginTop((c.height - len(lines) - 2) / 2).
		Padding(1).
		Render(content)
}

func (c *CommandPalette) RenderCentered(termWidth, termHeight int) string {
	if !c.isOpen {
		return ""
	}

	oldWidth := c.width
	oldHeight := c.height
	c.width = termWidth
	c.height = termHeight
	content := c.Render()
	c.width = oldWidth
	c.height = oldHeight
	return content
}
