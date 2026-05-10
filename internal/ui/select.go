package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/freecode/freecode/internal/ui/dialog"
)

type SelectOption struct {
	Title       string
	Value       string
	Description string
	Category    string
}

type SelectDialog struct {
	list      *dialog.SelectionList
	width     int
	isOpen    bool
	onSelect  func(string)
	colors    dialog.Colors
}

func NewSelectDialog() *SelectDialog {
	s := &SelectDialog{
		width:    60,
		isOpen:   false,
		colors:   dialog.Dark,
	}

	s.list = dialog.NewSelectionList(
		func(d *dialog.SelectionList) {
			d.Title = "Select an option"
			d.Width = 60
			d.Colors = s.colors
			d.OnSelect = func(item dialog.Item) {
				if s.onSelect != nil {
					if v, ok := item.Value.(string); ok {
						s.onSelect(v)
					} else {
						s.onSelect(item.ID)
					}
				}
			}
		},
	)

	return s
}

func (s *SelectDialog) SetWidth(w int) {
	s.width = w
}

func (s *SelectDialog) SetOptions(opts []SelectOption) {
	items := make([]dialog.Item, len(opts))
	for i, opt := range opts {
		items[i] = dialog.Item{
			ID:          opt.Value,
			Title:       opt.Title,
			Description: opt.Description,
			Category:    opt.Category,
			Value:       opt.Value,
		}
	}
	s.list.SetItems(items)
	s.list.Width = s.width
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
	s.list.SetItems([]dialog.Item{})
}

func (s *SelectDialog) SetFilter(filter string) {
	s.list.SetFilter(filter)
}

func (s *SelectDialog) Next() {
	s.list.Next()
}

func (s *SelectDialog) Prev() {
	s.list.Prev()
}

func (s *SelectDialog) GetSelected() *SelectOption {
	item := s.list.GetSelected()
	if item == nil {
		return nil
	}
	return &SelectOption{
		Title:       item.Title,
		Value:       item.ID,
		Description: item.Description,
		Category:    item.Category,
	}
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
		s.list.Confirm()
		s.Clear()
		return true
	case "up", "k":
		s.list.Prev()
		return true
	case "down", "j":
		s.list.Next()
		return true
	}
	return false
}

func (s *SelectDialog) Render() string {
	if !s.isOpen {
		return ""
	}

	s.list.Width = s.width
	lines := s.list.RenderList()

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(s.width).
		Height(s.list.Height).
		Background(lipgloss.Color(s.colors.Background)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(s.colors.Border)).
		Padding(1).
		Render(content)
}

func (s *SelectDialog) RenderWithTitle(title string) string {
	if !s.isOpen {
		return ""
	}

	s.list.Width = s.width

	oldRenderer := s.list.ItemRenderer
	s.list.ItemRenderer = func(item dialog.Item, selected, current bool, colors dialog.Colors) string {
		prefix := " "
		if selected {
			prefix = dialog.TextStyled("▶", colors.Primary)
		}
		titleStr := dialog.TextStyled(item.Title, colors.Text)
		line := prefix + " " + titleStr
		if item.Description != "" {
			line += " " + dialog.Muted("— "+item.Description, colors)
		}
		return line
	}
	defer func() { s.list.ItemRenderer = oldRenderer }()

	lines := s.list.RenderList()
	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(s.width).
		Height(s.list.Height).
		Background(lipgloss.Color(s.colors.Background)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(s.colors.Border)).
		Padding(1).
		Render(content)
}
