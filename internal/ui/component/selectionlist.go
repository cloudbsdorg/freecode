package component

import "github.com/freecode/freecode/internal/renderer"

type Item struct {
	ID          string
	Title       string
	Description string
	Category    string
	Footer      string
	Value       string
}

type SelectionList[R renderer.Renderer] struct {
	Component[R]
	Items       []Item
	Filtered    []Item
	Selected    int
	ScrollOff  int
	Filter     string
	Colors     SelectionColors
}

type SelectionColors struct {
	Background   string
	Foreground   string
	SelectedBg   string
	SelectedFg   string
	MutedColor   string
	BorderColor  string
}

func NewSelectionList[R renderer.Renderer](width, height int, colors SelectionColors) *SelectionList[R] {
	return &SelectionList[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  height,
			Visible: true,
		},
		Items:    []Item{},
		Filtered: []Item{},
		Selected: 0,
		ScrollOff: 0,
		Filter:   "",
		Colors:   colors,
	}
}

func (s *SelectionList[R]) SetItems(items []Item) {
	s.Items = items
	s.ApplyFilter()
	s.Selected = 0
	s.ScrollOff = 0
}

func (s *SelectionList[R]) ApplyFilter() {
	if s.Filter == "" {
		s.Filtered = s.Items
		return
	}
	lower := toLower(s.Filter)
	s.Filtered = []Item{}
	for _, item := range s.Items {
		if contains(toLower(item.Title), lower) ||
			contains(toLower(item.Description), lower) ||
			contains(toLower(item.Category), lower) {
			s.Filtered = append(s.Filtered, item)
		}
	}
}

func (s *SelectionList[R]) MoveUp() {
	if s.Selected > 0 {
		s.Selected--
		s.autoScroll()
	}
}

func (s *SelectionList[R]) MoveDown() {
	if s.Selected < len(s.Filtered)-1 {
		s.Selected++
		s.autoScroll()
	}
}

func (s *SelectionList[R]) autoScroll() {
	visibleH := s.Height - 2
	if visibleH < 1 {
		visibleH = 10
	}
	if s.Selected < s.ScrollOff {
		s.ScrollOff = s.Selected
	}
	if s.Selected >= s.ScrollOff+visibleH {
		s.ScrollOff = s.Selected - visibleH + 1
	}
}

func (s *SelectionList[R]) GetSelected() *Item {
	if s.Selected < 0 || s.Selected >= len(s.Filtered) {
		return nil
	}
	return &s.Filtered[s.Selected]
}

func (s *SelectionList[R]) SetFilter(filter string) {
	s.Filter = filter
	s.ApplyFilter()
	if s.Selected >= len(s.Filtered) {
		s.Selected = len(s.Filtered) - 1
		if s.Selected < 0 {
			s.Selected = 0
		}
	}
}

func (s *SelectionList[R]) Render(r R) string {
	if !s.Visible || len(s.Filtered) == 0 {
		return r.RenderBox(s.X, s.Y, s.Width, s.Height, s.Colors.Background)
	}

	lines := []string{}
	start := s.ScrollOff
	end := start + s.Height - 2
	if end > len(s.Filtered) {
		end = len(s.Filtered)
	}

	for i := start; i < end; i++ {
		item := s.Filtered[i]
		prefix := "  "
		if i == s.Selected {
			prefix = "> "
		}
		text := prefix + item.Title
		if item.Description != "" {
			text += " - " + item.Description
		}
		if i == s.Selected {
			lines = append(lines, r.RenderSelected(text, s.X+1, s.Y+1+(i-start), s.Width-2, s.Colors.SelectedFg, s.Colors.SelectedBg))
		} else {
			lines = append(lines, r.RenderText(text, s.X+1, s.Y+1+(i-start), s.Colors.Foreground))
		}
	}

	border := r.RenderBorder(s.X, s.Y, s.Width, s.Height, s.Colors.BorderColor)

	for i := len(lines); i < s.Height-2; i++ {
		lines = append(lines, r.RenderText("", s.X+1, s.Y+1+i, s.Colors.Foreground))
	}

	result := ""
	for _, line := range lines {
		result += line + "\n"
	}
	return border + result
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
