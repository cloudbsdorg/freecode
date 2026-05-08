package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type AutocompleteItem struct {
	Label       string
	Value       string
	Description string
	Category    string
}

type AutocompleteDialog struct {
	visible   bool
	items     []AutocompleteItem
	filtered  []AutocompleteItem
	selected  int
	input    string
	width    int
	onSelect func(string)
	onClose  func()
}

func NewAutocompleteDialog() *AutocompleteDialog {
	return &AutocompleteDialog{
		visible:  false,
		items:    []AutocompleteItem{},
		filtered: []AutocompleteItem{},
		selected: 0,
		input:   "",
		width:   50,
	}
}

func (a *AutocompleteDialog) SetWidth(w int) {
	a.width = w
}

func (a *AutocompleteDialog) SetItems(items []AutocompleteItem) {
	a.items = items
	a.filtered = items
	a.selected = 0
}

func (a *AutocompleteDialog) Show(input string) {
	a.visible = true
	a.input = input
	a.applyFilter()
}

func (a *AutocompleteDialog) Hide() {
	a.visible = false
}

func (a *AutocompleteDialog) IsVisible() bool {
	return a.visible
}

func (a *AutocompleteDialog) SetOnSelect(fn func(string)) {
	a.onSelect = fn
}

func (a *AutocompleteDialog) SetOnClose(fn func()) {
	a.onClose = fn
}

func (a *AutocompleteDialog) applyFilter() {
	if a.input == "" {
		a.filtered = a.items
		return
	}

	needle := strings.ToLower(a.input)
	var result []AutocompleteItem
	for _, item := range a.items {
		if strings.Contains(strings.ToLower(item.Label), needle) {
			result = append(result, item)
		} else if strings.Contains(strings.ToLower(item.Description), needle) {
			result = append(result, item)
		}
	}
	a.filtered = result
	if a.selected >= len(a.filtered) {
		a.selected = 0
	}
}

func (a *AutocompleteDialog) Next() {
	if len(a.filtered) == 0 {
		return
	}
	a.selected = (a.selected + 1) % len(a.filtered)
}

func (a *AutocompleteDialog) Prev() {
	if len(a.filtered) == 0 {
		return
	}
	a.selected = (a.selected - 1 + len(a.filtered)) % len(a.filtered)
}

func (a *AutocompleteDialog) Select() *AutocompleteItem {
	if a.selected < 0 || a.selected >= len(a.filtered) {
		return nil
	}
	item := &a.filtered[a.selected]
	if a.onSelect != nil {
		a.onSelect(item.Value)
	}
	a.Hide()
	return item
}

func (a *AutocompleteDialog) HandleKey(msg string) bool {
	if !a.visible {
		return false
	}

	switch msg {
	case "escape":
		a.Hide()
		if a.onClose != nil {
			a.onClose()
		}
		return true
	case "enter":
		a.Select()
		return true
	case "up", "k":
		a.Prev()
		return true
	case "down", "j":
		a.Next()
		return true
	}
	return false
}

func (a *AutocompleteDialog) Render() string {
	if !a.visible || len(a.filtered) == 0 {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#2D2D2D")).
		Border(lipgloss.HiddenBorder()).
		Width(a.width)

	return dialogStyle.Render(a.renderContent())
}

func (a *AutocompleteDialog) renderContent() string {
	var lines []string

	for i, item := range a.filtered {
		lines = append(lines, a.renderItem(item, i == a.selected))
	}

	return strings.Join(lines, "\n")
}

func (a *AutocompleteDialog) renderItem(item AutocompleteItem, selected bool) string {
	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	}

	var labelStyle lipgloss.Style
	if selected {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	} else {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D0D0D0"))
	}

	line := prefix + " " + labelStyle.Render(item.Label)

	if item.Category != "" {
		catStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#007ACC"))
		line += " " + catStyle.Render("["+item.Category+"]")
	}

	if item.Description != "" {
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		line += " — " + descStyle.Render(item.Description)
	}

	return line
}