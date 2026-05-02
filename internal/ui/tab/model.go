package tab

import (
	"sync"
)

type TabState struct {
	mu          sync.RWMutex
	tabs        []Tab
	active      int
	width       int
	height      int
}

type Tab struct {
	ID             string
	Name           string
	Sessions       []string
	ActiveSession  string
	SplitVertical bool
	SplitRatio    float64
}

func New() *TabState {
	return &TabState{
		tabs:   make([]Tab, 0),
		active: 0,
	}
}

func (m *TabState) AddTab(name string) Tab {
	m.mu.Lock()
	defer m.mu.Unlock()

	tab := Tab{
		ID:       generateID(),
		Name:     name,
		Sessions: make([]string, 0),
	}
	m.tabs = append(m.tabs, tab)
	return tab
}

func (m *TabState) CloseTab(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, t := range m.tabs {
		if t.ID == id {
			m.tabs = append(m.tabs[:i], m.tabs[i+1:]...)
			if m.active >= len(m.tabs) {
				m.active = len(m.tabs) - 1
			}
			return true
		}
	}
	return false
}

func (m *TabState) SetActive(idx int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if idx >= 0 && idx < len(m.tabs) {
		m.active = idx
	}
}

func (m *TabState) GetActive() Tab {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.active >= 0 && m.active < len(m.tabs) {
		return m.tabs[m.active]
	}
	return Tab{}
}

func (m *TabState) List() []Tab {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tabs
}

func (m *TabState) AddSession(tabID, sessionID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, t := range m.tabs {
		if t.ID == tabID {
			m.tabs[i].Sessions = append(m.tabs[i].Sessions, sessionID)
			m.tabs[i].ActiveSession = sessionID
			return true
		}
	}
	return false
}

func (m *TabState) RemoveSession(tabID, sessionID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, t := range m.tabs {
		if t.ID == tabID {
			for j, s := range t.Sessions {
				if s == sessionID {
					m.tabs[i].Sessions = append(m.tabs[i].Sessions[:j], m.tabs[i].Sessions[j+1:]...)
					return true
				}
			}
		}
	}
	return false
}

func (m *TabState) SetSplit(tabID string, vertical bool, ratio float64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, t := range m.tabs {
		if t.ID == tabID {
			m.tabs[i].SplitVertical = vertical
			m.tabs[i].SplitRatio = ratio
			return true
		}
	}
	return false
}

func generateID() string {
	return "tab-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
