package dialog

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Colors struct {
	Background      string
	BackgroundAlt  string
	Text           string
	TextMuted      string
	Primary        string
	PrimaryBg      string
	PrimaryFg      string
	Secondary      string
	Success        string
	Warning        string
	Error          string
	SelectedBg     string
	SelectedFg     string
	Border         string
	Muted          string
}

var Dark = Colors{
	Background:      "#1E1E1E",
	BackgroundAlt:    "#2D2D2D",
	Text:            "#E0E0E0",
	TextMuted:       "#808080",
	Primary:         "#22C55E",
	PrimaryBg:       "#007ACC",
	PrimaryFg:       "#FFFFFF",
	Secondary:       "#007ACC",
	Success:         "#4EC9B0",
	Warning:         "#FFCC00",
	Error:           "#F44747",
	SelectedBg:      "#007ACC",
	SelectedFg:      "#FFFFFF",
	Border:          "#3C3C3C",
	Muted:           "#808080",
}

var Light = Colors{
	Background:      "#FFFFFF",
	BackgroundAlt:    "#F5F5F5",
	Text:            "#333333",
	TextMuted:       "#666666",
	Primary:         "#0066CC",
	PrimaryBg:       "#0066CC",
	PrimaryFg:       "#FFFFFF",
	Secondary:       "#0066CC",
	Success:         "#00875A",
	Warning:         "#FFAB00",
	Error:           "#DE350B",
	SelectedBg:      "#0066CC",
	SelectedFg:      "#FFFFFF",
	Border:          "#CCCCCC",
	Muted:           "#666666",
}

var DefaultColors = Dark

type Item struct {
	ID          string
	Title       string
	Description string
	Category    string
	Disabled    bool
	Footer      string
	Value       interface{}
}

type ItemRenderer func(item Item, selected, current bool, colors Colors) string
type FilterRenderer func(filter string, colors Colors) string
type FooterRenderer func(colors Colors) string

type SelectionList struct {
	Width          int
	Height         int
	Items          []Item
	Filtered       []Item
	Selected       int
	ScrollOffset   int
	Filter         string
	Colors         Colors
	Title          string
	SkipFilter     bool
	Flat           bool
	Current        interface{}
	ItemRenderer   ItemRenderer
	FilterRenderer FilterRenderer
	FooterRenderer FooterRenderer
	OnSelect       func(Item)
	OnMove         func(Item)
	OnFilter       func(string)
}

func NewSelectionList(opts ...func(*SelectionList)) *SelectionList {
	s := &SelectionList{
		Width:      60,
		Height:     20,
		Items:      []Item{},
		Filtered:   []Item{},
		Selected:   0,
		Filter:     "",
		Colors:     DefaultColors,
		Title:      "Select",
		SkipFilter: false,
		Flat:       false,
		ItemRenderer: func(item Item, selected, current bool, colors Colors) string {
			prefix := "  "
			if selected {
				prefix = SelectedPrefix("▶", colors)
			}
			titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
			if selected {
				titleStyle = lipgloss.NewStyle().
					Background(lipgloss.Color(colors.SelectedBg)).
					Foreground(lipgloss.Color(colors.SelectedFg))
			} else if current {
				titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Primary))
			}
			line := prefix + " " + titleStyle.Render(item.Title)
			if item.Description != "" {
				line += " " + Muted("— "+item.Description, colors)
			}
			if item.Footer != "" {
				line += " " + Muted(item.Footer, colors)
			}
			return line
		},
		FilterRenderer: func(filter string, colors Colors) string {
			return Input("Filter: "+filter+"_", colors)
		},
		FooterRenderer: func(colors Colors) string {
			return Muted("↑↓ navigate  enter select  esc close", colors)
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *SelectionList) SetItems(items []Item) {
	s.Items = items
	s.Selected = 0
	s.ScrollOffset = 0
	s.ApplyFilter()
}

func (s *SelectionList) SetFilter(filter string) {
	s.Filter = filter
	s.ApplyFilter()
	if s.OnFilter != nil {
		s.OnFilter(filter)
	}
}

func (s *SelectionList) ApplyFilter() {
	if s.SkipFilter || s.Filter == "" {
		s.Filtered = []Item{}
		for _, item := range s.Items {
			if !item.Disabled {
				s.Filtered = append(s.Filtered, item)
			}
		}
	} else {
		needle := strings.ToLower(s.Filter)
		s.Filtered = []Item{}
		for _, item := range s.Items {
			if item.Disabled {
				continue
			}
			if strings.Contains(strings.ToLower(item.Title), needle) {
				s.Filtered = append(s.Filtered, item)
			} else if item.Category != "" && strings.Contains(strings.ToLower(item.Category), needle) {
				s.Filtered = append(s.Filtered, item)
			} else if strings.Contains(strings.ToLower(item.Description), needle) {
				s.Filtered = append(s.Filtered, item)
			}
		}
	}
	if s.Selected >= len(s.Filtered) {
		s.Selected = 0
	}
}

func (s *SelectionList) Next() {
	if len(s.Filtered) == 0 {
		return
	}
	s.Selected = (s.Selected + 1) % len(s.Filtered)
	s.notifyMove()
}

func (s *SelectionList) Prev() {
	if len(s.Filtered) == 0 {
		return
	}
	s.Selected = (s.Selected - 1 + len(s.Filtered)) % len(s.Filtered)
	s.notifyMove()
}

func (s *SelectionList) MoveUp() {
	if s.Selected > 0 {
		s.Selected--
		s.autoScroll()
		s.notifyMove()
	}
}

func (s *SelectionList) MoveDown() {
	if s.Selected < len(s.Filtered)-1 {
		s.Selected++
		s.autoScroll()
		s.notifyMove()
	}
}

func (s *SelectionList) autoScroll() {
	if s.Height <= 0 {
		return
	}
	visibleHeight := s.Height - 4
	if visibleHeight < 1 {
		visibleHeight = 10
	}
	if s.Selected < s.ScrollOffset {
		s.ScrollOffset = s.Selected
	}
	if s.Selected >= s.ScrollOffset+visibleHeight {
		s.ScrollOffset = s.Selected - visibleHeight + 1
	}
}

func (s *SelectionList) notifyMove() {
	if s.OnMove != nil && s.Selected >= 0 && s.Selected < len(s.Filtered) {
		s.OnMove(s.Filtered[s.Selected])
	}
}

func (s *SelectionList) Select(idx int) bool {
	if idx >= 0 && idx < len(s.Filtered) {
		s.Selected = idx
		s.notifyMove()
		return true
	}
	return false
}

func (s *SelectionList) GetSelected() *Item {
	if s.Selected < 0 || s.Selected >= len(s.Filtered) {
		return nil
	}
	return &s.Filtered[s.Selected]
}

func (s *SelectionList) SelectByValue(value interface{}) bool {
	for i, item := range s.Filtered {
		if item.Value == value {
			s.Selected = i
			s.notifyMove()
			return true
		}
	}
	return false
}

func (s *SelectionList) Confirm() bool {
	item := s.GetSelected()
	if item != nil && s.OnSelect != nil {
		s.OnSelect(*item)
		return true
	}
	return false
}

func (s *SelectionList) IsItemCurrent(item Item) bool {
	if s.Current == nil {
		return false
	}
	return item.Value == s.Current || item.ID == s.Current
}

func (s *SelectionList) RenderList() []string {
	var lines []string

	if !s.SkipFilter && s.FilterRenderer != nil && s.Filter != "" {
		lines = append(lines, s.FilterRenderer(s.Filter, s.Colors))
		lines = append(lines, "")
	}

	if len(s.Filtered) == 0 {
		lines = append(lines, Muted("  No matches", s.Colors))
	} else {
		lines = append(lines, s.renderVisibleItems()...)
	}

	return lines
}

func (s *SelectionList) renderVisibleItems() []string {
	if s.Height <= 0 || s.Height >= len(s.Filtered)+4 {
		return s.renderItems()
	}

	visibleHeight := s.Height - 4
	start := s.ScrollOffset
	end := start + visibleHeight
	if end > len(s.Filtered) {
		end = len(s.Filtered)
	}
	if start >= len(s.Filtered) {
		start = 0
		end = visibleHeight
		if end > len(s.Filtered) {
			end = len(s.Filtered)
		}
	}

	var lines []string
	for i := start; i < end; i++ {
		item := s.Filtered[i]
		lines = append(lines, s.ItemRenderer(item, s.IsSelected(item), s.IsItemCurrent(item), s.Colors))
	}
	return lines
}

func (s *SelectionList) Render() string {
	var lines []string

	lines = append(lines, Header(s.Title, s.Colors))
	lines = append(lines, "")

	if !s.SkipFilter && s.FilterRenderer != nil {
		lines = append(lines, s.FilterRenderer(s.Filter, s.Colors))
		lines = append(lines, "")
	}

	if len(s.Filtered) == 0 {
		lines = append(lines, Muted("  No matches", s.Colors))
	} else {
		lines = append(lines, s.renderItems()...)
	}

	lines = append(lines, "")
	if s.FooterRenderer != nil {
		lines = append(lines, s.FooterRenderer(s.Colors))
	}

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Colors.Background)).
		Width(s.Width).
		Height(s.Height).
		Render(content)
}

func (s *SelectionList) renderItems() []string {
	if s.Flat || !s.HasCategories() {
		return s.renderFlatItems()
	}
	return s.renderGroupedItems()
}

func (s *SelectionList) HasCategories() bool {
	for _, item := range s.Filtered {
		if item.Category != "" {
			return true
		}
	}
	return false
}

func (s *SelectionList) renderFlatItems() []string {
	var lines []string
	for _, item := range s.Filtered {
		lines = append(lines, s.ItemRenderer(item, s.IsSelected(item), s.IsItemCurrent(item), s.Colors))
	}
	return lines
}

func (s *SelectionList) renderGroupedItems() []string {
	var lines []string
	categoryGroups := make(map[string][]Item)

	for _, item := range s.Filtered {
		cat := item.Category
		if cat == "" {
			cat = "Other"
		}
		categoryGroups[cat] = append(categoryGroups[cat], item)
	}

	first := true
	for category, items := range categoryGroups {
		if !first {
			lines = append(lines, "")
		}
		if category != "Other" {
			lines = append(lines, TextStyled("  "+category, s.Colors.Primary))
		}
		for _, item := range items {
			lines = append(lines, s.ItemRenderer(item, s.IsSelected(item), s.IsItemCurrent(item), s.Colors))
		}
		first = false
	}
	return lines
}

func (s *SelectionList) IsSelected(item Item) bool {
	idx := -1
	for i, f := range s.Filtered {
		if f.ID == item.ID {
			idx = i
			break
		}
	}
	return idx == s.Selected
}

func Style(colors Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(colors.Background))
}

func Styled(colors Colors, width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(colors.Background)).
		Width(width).
		Height(height)
}

func Header(text string, colors Colors) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.Primary)).
		Padding(0, 1).
		Render(text)
}

func Muted(text string, colors Colors) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.TextMuted)).
		Render(text)
}

func TextStyled(text string, color string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(text)
}

func ErrorText(text string, colors Colors) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Error)).
		Render(text)
}

func SuccessText(text string, colors Colors) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Success)).
		Render(text)
}

func Input(text string, colors Colors) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Background(lipgloss.Color(colors.BackgroundAlt)).
		Padding(0, 1).
		Render(text)
}

func Selected(text string, colors Colors) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(colors.SelectedBg)).
		Foreground(lipgloss.Color(colors.SelectedFg)).
		Padding(0, 1).
		Render(text)
}

func SelectedPrefix(prefix string, colors Colors) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(colors.SelectedBg)).
		Foreground(lipgloss.Color(colors.SelectedFg)).
		Render(prefix)
}
