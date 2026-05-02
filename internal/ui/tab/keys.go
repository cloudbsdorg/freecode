package tab

import (
	"github.com/charmbracelet/bubbletea"
)

type KeyHandler struct{}

func NewKeyHandler() *KeyHandler {
	return &KeyHandler{}
}

func (h *KeyHandler) Handle(msg tea.KeyMsg, m *TabState) (*TabState, tea.Cmd) {
	switch msg.String() {
	case "ctrl+t":
		name := "New Tab"
		tab := m.AddTab(name)
		m.SetActive(len(m.tabs) - 1)
		return m, func() tea.Msg {
			return NewTabCreatedMsg{Tab: tab}
		}

	case "ctrl+w":
		active := m.GetActive()
		if len(m.tabs) > 1 {
			m.CloseTab(active.ID)
		}
		return m, nil

	case "tab":
		current := m.active
		m.SetActive((current + 1) % len(m.tabs))
		return m, nil

	case "shift+tab":
		current := m.active
		m.SetActive((current - 1 + len(m.tabs)) % len(m.tabs))
		return m, nil

	case "ctrl+shift+v":
		active := m.GetActive()
		m.SetSplit(active.ID, true, 0.5)
		return m, nil

	case "ctrl+shift+h":
		active := m.GetActive()
		m.SetSplit(active.ID, false, 0.5)
		return m, nil
	}

	return m, nil
}

type NewTabCreatedMsg struct {
	Tab Tab
}
