package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var TabBarStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#2D2D2D")).
	Foreground(lipgloss.Color("#FFFFFF")).
	Padding(0, 1)

var ActiveTabStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#007ACC")).
	Foreground(lipgloss.Color("#FFFFFF")).
	Bold(true).
	Padding(0, 1)

var InactiveTabStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#3D3D3D")).
	Foreground(lipgloss.Color("#808080")).
	Padding(0, 1)

type TabBarComponent struct {
	tabs      []TabItem
	activeIdx int
	width     int
}

type TabItem struct {
	ID    string
	Name  string
	Dirty bool
}

func NewTabBar() *TabBarComponent {
	return &TabBarComponent{
		tabs:      make([]TabItem, 0),
		activeIdx: 0,
	}
}

func (t *TabBarComponent) AddTab(id, name string) {
	t.tabs = append(t.tabs, TabItem{ID: id, Name: name})
	t.activeIdx = len(t.tabs) - 1
}

func (t *TabBarComponent) CloseTab(idx int) bool {
	if idx < 0 || idx >= len(t.tabs) {
		return false
	}
	t.tabs = append(t.tabs[:idx], t.tabs[idx+1:]...)
	if t.activeIdx >= len(t.tabs) {
		t.activeIdx = len(t.tabs) - 1
	}
	if t.activeIdx < 0 {
		t.activeIdx = 0
	}
	return true
}

func (t *TabBarComponent) SetActive(idx int) {
	if idx >= 0 && idx < len(t.tabs) {
		t.activeIdx = idx
	}
}

func (t *TabBarComponent) NextTab() {
	if len(t.tabs) > 0 {
		t.activeIdx = (t.activeIdx + 1) % len(t.tabs)
	}
}

func (t *TabBarComponent) PrevTab() {
	if len(t.tabs) > 0 {
		t.activeIdx = (t.activeIdx - 1 + len(t.tabs)) % len(t.tabs)
	}
}

func (t *TabBarComponent) ActiveTab() *TabItem {
	if t.activeIdx >= 0 && t.activeIdx < len(t.tabs) {
		return &t.tabs[t.activeIdx]
	}
	return nil
}

func (t *TabBarComponent) SetWidth(w int) {
	t.width = w
}

func (t *TabBarComponent) Render() string {
	if len(t.tabs) == 0 {
		return ""
	}

	const (
		tabPrefix    = "["
		tabSuffix    = "]"
		newTabButton = "[+]"
		separator    = ""
	)

	fixedWidth := len(tabPrefix) + len(tabSuffix) + len(newTabButton) + len(separator)
	availableWidth := t.width - fixedWidth

	if availableWidth < 20 {
		availableWidth = 80
	}

	tabStr := ""
	for i, tab := range t.tabs {
		name := tab.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		var itemStr string
		if i == t.activeIdx {
			itemStr = ActiveTabStyle.Render(fmt.Sprintf("%s%s%s", tabPrefix, name, tabSuffix))
		} else {
			itemStr = InactiveTabStyle.Render(fmt.Sprintf("%s%s%s", tabPrefix, name, tabSuffix))
		}
		tabStr += itemStr + separator
	}

	newTabStr := InactiveTabStyle.Render(newTabButton)

	totalLen := len(tabStr) + len(newTabStr)
	if totalLen > t.width && t.width > 0 {
		maxTabLen := (t.width - len(newTabStr)) / (len(tabPrefix) + len(tabSuffix) + 1)
		truncated := ""
		for i := 0; i < len(t.tabs) && i < maxTabLen; i++ {
			tab := t.tabs[i]
			name := tab.Name
			if len(name) > 15 {
				name = name[:12] + "..."
			}
			if i == t.activeIdx {
				truncated += ActiveTabStyle.Render(fmt.Sprintf("%s%s%s", tabPrefix, name, tabSuffix)) + separator
			} else {
				truncated += InactiveTabStyle.Render(fmt.Sprintf("%s%s%s", tabPrefix, name, tabSuffix)) + separator
			}
		}
		tabStr = truncated
	}

	return TabBarStyle.Render(tabStr + newTabStr)
}
