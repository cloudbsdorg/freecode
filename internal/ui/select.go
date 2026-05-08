package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type SelectOption struct {
	Title       string
	Value       string
	Description string
	Category    string
}

type SelectDialog struct {
	options   []SelectOption
	filtered  []SelectOption
	selected  int
	filter   string
	width    int
	isOpen   bool
	onSelect func(string)
}

func NewSelectDialog() *SelectDialog {
	return &SelectDialog{
		options:  []SelectOption{},
		filtered: []SelectOption{},
		selected: 0,
		filter:   "",
		width:    60,
		isOpen:   false,
	}
}

func (s *SelectDialog) SetWidth(w int) {
	s.width = w
}

func (s *SelectDialog) SetOptions(opts []SelectOption) {
	s.options = opts
	s.filtered = opts
	s.selected = 0
	s.filter = ""
	s.isOpen = true
}

func (s *SelectDialog) SetOnSelect(fn func(string)) {
	s.onSelect = fn
}

func (s *SelectDialog) IsVisible() bool {
	return s.isOpen
}

func (s *SelectDialog) Clear() {
	s.isOpen = false
	s.filter = ""
	s.options = []SelectOption{}
	s.filtered = []SelectOption{}
	s.selected = 0
}

func (s *SelectDialog) SetFilter(filter string) {
	s.filter = filter
	s.applyFilter()
}

func (s *SelectDialog) applyFilter() {
	if s.filter == "" {
		s.filtered = s.options
		return
	}
	needle := strings.ToLower(s.filter)
	var result []SelectOption
	for _, opt := range s.options {
		if strings.Contains(strings.ToLower(opt.Title), needle) {
			result = append(result, opt)
		} else if opt.Category != "" && strings.Contains(strings.ToLower(opt.Category), needle) {
			result = append(result, opt)
		} else if strings.Contains(strings.ToLower(opt.Description), needle) {
			result = append(result, opt)
		}
	}
	s.filtered = result
	if s.selected >= len(s.filtered) {
		s.selected = 0
	}
}

func (s *SelectDialog) Next() {
	if len(s.filtered) == 0 {
		return
	}
	s.selected = (s.selected + 1) % len(s.filtered)
}

func (s *SelectDialog) Prev() {
	if len(s.filtered) == 0 {
		return
	}
	s.selected = (s.selected - 1 + len(s.filtered)) % len(s.filtered)
}

func (s *SelectDialog) GetSelected() *SelectOption {
	if s.selected < 0 || s.selected >= len(s.filtered) {
		return nil
	}
	return &s.filtered[s.selected]
}

func (s *SelectDialog) HandleKey(msg string) bool {
	if !s.isOpen {
		return false
	}

	switch msg {
	case "escape":
		s.Clear()
		return true
	case "enter":
		sel := s.GetSelected()
		if sel != nil && s.onSelect != nil {
			s.onSelect(sel.Value)
		}
		s.Clear()
		return true
	case "up", "k":
		s.Prev()
		return true
	case "down", "j":
		s.Next()
		return true
	}
	return false
}

func (s *SelectDialog) Render() string {
	if !s.isOpen {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E1E")).
		Border(lipgloss.HiddenBorder()).
		Width(s.width)

	return dialogStyle.Render(s.renderContent())
}

func (s *SelectDialog) renderContent() string {
	var lines []string

	lines = append(lines, s.renderHeader())
	lines = append(lines, "")
	lines = append(lines, s.renderFilter())
	lines = append(lines, "")
	lines = append(lines, s.renderOptions()...)
	lines = append(lines, "")
	lines = append(lines, s.renderHints())

	return strings.Join(lines, "\n")
}

func (s *SelectDialog) renderHeader() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Render("Select an option")
}

func (s *SelectDialog) renderFilter() string {
	filterStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#3C3C3C")).
		Foreground(lipgloss.Color("#E0E0E0")).
		Padding(0, 1)
	return filterStyle.Render("Filter: " + s.filter + "_")
}

func (s *SelectDialog) renderOptions() []string {
	var lines []string

	if len(s.filtered) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Render("  No matches"))
		return lines
	}

	groups := make(map[string][]SelectOption)
	hasGroups := false

	for _, opt := range s.filtered {
		if opt.Category != "" {
			hasGroups = true
			groups[opt.Category] = append(groups[opt.Category], opt)
		} else {
			groups[""] = append(groups[""], opt)
		}
	}

	if hasGroups {
		for category, opts := range groups {
			if category != "" {
				catStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#007ACC")).
					Bold(true)
				lines = append(lines, catStyle.Render(category))
			}
			for _, opt := range opts {
				lines = append(lines, s.renderOption(opt))
			}
			if category != "" {
				lines = append(lines, "")
			}
		}
	} else {
		for _, opt := range s.filtered {
			lines = append(lines, s.renderOption(opt))
		}
	}

	return lines
}

func (s *SelectDialog) renderOption(opt SelectOption) string {
	idx := -1
	for i, f := range s.filtered {
		if f.Value == opt.Value {
			idx = i
			break
		}
	}

	selected := idx == s.selected

	var prefix string
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	} else {
		prefix = " "
	}

	var titleStyle lipgloss.Style
	if selected {
		titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	} else {
		titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D0D0D0"))
	}

	line := prefix + " " + titleStyle.Render(opt.Title)

	if opt.Description != "" {
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		line += " " + descStyle.Render("— "+opt.Description)
	}

	return line
}

func (s *SelectDialog) renderHints() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	return hintStyle.Render("↑↓ navigate  enter select  esc close")
}