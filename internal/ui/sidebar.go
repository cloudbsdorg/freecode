package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var SidebarStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#252525")).
	Foreground(lipgloss.Color("#E0E0E0")).
	Width(42).
	Padding(1, 2)

var SidebarHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#007ACC")).
	Bold(true)

var SidebarItemStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#E0E0E0"))

var SidebarSelectedStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#007ACC")).
	Foreground(lipgloss.Color("#FFFFFF")).
	Bold(true)

var SidebarTimestampStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#606060"))

var SidebarDirtyIndicator = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFCC00")).
	Bold(true)

type SidebarItem struct {
	ID        string
	Title     string
	Timestamp string
	Dirty     bool
	Selected  bool
}

type Sidebar struct {
	items      []SidebarItem
	selectedIdx int
	width      int
	height     int
	isOpen     bool
}

func NewSidebar() *Sidebar {
	return &Sidebar{
		items:      make([]SidebarItem, 0),
		selectedIdx: 0,
		width:      42,
		height:     20,
		isOpen:     true,
	}
}

func (s *Sidebar) AddItem(item SidebarItem) {
	s.items = append(s.items, item)
}

func (s *Sidebar) SetItems(items []SidebarItem) {
	s.items = items
	if s.selectedIdx >= len(s.items) {
		s.selectedIdx = len(s.items) - 1
	}
	if s.selectedIdx < 0 {
		s.selectedIdx = 0
	}
}

func (s *Sidebar) Clear() {
	s.items = make([]SidebarItem, 0)
	s.selectedIdx = 0
}

func (s *Sidebar) Open() {
	s.isOpen = true
}

func (s *Sidebar) Close() {
	s.isOpen = false
}

func (s *Sidebar) Toggle() {
	s.isOpen = !s.isOpen
}

func (s *Sidebar) IsOpen() bool {
	return s.isOpen
}

func (s *Sidebar) SetWidth(w int) {
	s.width = w
}

func (s *Sidebar) SetHeight(h int) {
	s.height = h
}

func (s *Sidebar) SelectNext() {
	if s.selectedIdx < len(s.items)-1 {
		s.selectedIdx++
	}
}

func (s *Sidebar) SelectPrev() {
	if s.selectedIdx > 0 {
		s.selectedIdx--
	}
}

func (s *Sidebar) SelectedItem() *SidebarItem {
	if s.selectedIdx >= 0 && s.selectedIdx < len(s.items) {
		return &s.items[s.selectedIdx]
	}
	return nil
}

func (s *Sidebar) SetSelected(idx int) {
	if idx >= 0 && idx < len(s.items) {
		s.selectedIdx = idx
	}
}

func (s *Sidebar) Width() int {
	return s.width
}

func (s *Sidebar) Render() string {
	if !s.isOpen {
		return ""
	}

	lines := []string{}

	header := SidebarHeaderStyle.Render("Sessions")
	lines = append(lines, header, "")

	maxItems := s.height - 4
	if maxItems < 1 {
		maxItems = 10
	}

	displayItems := s.items
	if len(displayItems) > maxItems {
		start := 0
		if s.selectedIdx >= maxItems {
			start = s.selectedIdx - maxItems + 1
		}
		end := start + maxItems
		if end > len(displayItems) {
			end = len(displayItems)
		}
		displayItems = displayItems[start:end]
	}

	for i, item := range displayItems {
		actualIdx := i
		for j, stored := range s.items {
			if stored.ID == item.ID {
				actualIdx = j
				break
			}
		}

		dirty := ""
		if item.Dirty {
			dirty = " " + SidebarDirtyIndicator.Render("●")
		}

		title := item.Title
		if len(title) > s.width-10 {
			title = title[:s.width-13] + "..."
		}

		itemStr := fmt.Sprintf("  %s%s", title, dirty)
		if item.Timestamp != "" {
			itemStr += "\n    " + SidebarTimestampStyle.Render(item.Timestamp)
		}

		if actualIdx == s.selectedIdx {
			lines = append(lines, SidebarSelectedStyle.Render(itemStr))
		} else {
			lines = append(lines, SidebarItemStyle.Render(itemStr))
		}
	}

	if len(s.items) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#606060")).Render("  No sessions"))
	}

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("#606060")).Render("  Ctrl+T: New | Enter: Select")
	lines = append(lines, "", footer)

	result := strings.Join(lines, "\n")
	return SidebarStyle.Render(result)
}
