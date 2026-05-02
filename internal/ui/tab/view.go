package tab

import (
	"fmt"
	"strings"
)

func (m *TabState) View() string {
	var s strings.Builder

	s.WriteString(m.tabBar())
	s.WriteString("\n")
	s.WriteString(m.content())
	s.WriteString("\n")
	s.WriteString(m.statusLine())

	return s.String()
}

func (m *TabState) tabBar() string {
	var s strings.Builder
	tabs := m.List()

	for i, tab := range tabs {
		if i == m.active {
			s.WriteString(fmt.Sprintf("[%s]", tab.Name))
		} else {
			s.WriteString(fmt.Sprintf(" %s ", tab.Name))
		}
		s.WriteString(" ")
	}

	s.WriteString("[+]")
	return s.String()
}

func (m *TabState) content() string {
	active := m.GetActive()
	if len(active.Sessions) == 0 {
		return "(no session)"
	}
	return fmt.Sprintf("Session: %s", active.ActiveSession)
}

func (m *TabState) statusLine() string {
	active := m.GetActive()
	parts := []string{
		fmt.Sprintf("Tab: %s", active.Name),
		fmt.Sprintf("Sessions: %d", len(active.Sessions)),
	}

	if active.SplitVertical {
		parts = append(parts, "VSplit")
	} else {
		parts = append(parts, "HSplit")
	}

	return strings.Join(parts, " | ")
}

func (m *TabState) SetSize(width, height int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.width = width
	m.height = height
}
