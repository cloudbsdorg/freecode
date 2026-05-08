package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/freecode/freecode/internal/config"
	"github.com/google/uuid"
)

type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	tabs     map[string]*Tab
	config   *config.Config
}

type Session struct {
	ID        string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Model     string
	Agent     string
	TabID     string
	Messages  []Message
	Metadata  map[string]interface{}
}

type Tab struct {
	ID            string
	Name          string
	CreatedAt     time.Time
	Sessions      []string
	ActiveSession string
}

type Message struct {
	ID        string
	Role      string
	Content   string
	Timestamp time.Time
	Parts     []MessagePart
}

type MessagePart struct {
	Type    string
	Content string
	Tool    string
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		tabs:     make(map[string]*Tab),
		config:   cfg,
	}
}

func (m *Manager) CreateSession(title, model, agent string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New().String()
	sess := &Session{
		ID:        id,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Model:     model,
		Agent:     agent,
		Messages:  make([]Message, 0),
		Metadata:  make(map[string]interface{}),
	}

	m.sessions[id] = sess
	return sess, nil
}

func (m *Manager) GetSession(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sess, ok := m.sessions[id]
	return sess, ok
}

func (m *Manager) ListSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sessions := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		sessions = append(sessions, s)
	}
	return sessions
}

func (m *Manager) DeleteSession(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.sessions[id]; !ok {
		return fmt.Errorf("session not found: %s", id)
	}
	delete(m.sessions, id)
	return nil
}

func (m *Manager) AddMessage(sessionID, role, content string) (*Message, error) {
	return m.AddMessageWithParts(sessionID, role, content, nil)
}

func (m *Manager) AddMessageWithParts(sessionID, role, content string, parts []MessagePart) (*Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sess, ok := m.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	msg := Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
		Parts:     parts,
	}

	sess.Messages = append(sess.Messages, msg)
	sess.UpdatedAt = time.Now()

	return &msg, nil
}

func (m *Manager) CreateTab(name string) (*Tab, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New().String()
	tab := &Tab{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		Sessions:  make([]string, 0),
	}

	m.tabs[id] = tab
	return tab, nil
}

func (m *Manager) GetTab(id string) (*Tab, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tab, ok := m.tabs[id]
	return tab, ok
}

func (m *Manager) ListTabs() []*Tab {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tabs := make([]*Tab, 0, len(m.tabs))
	for _, t := range m.tabs {
		tabs = append(tabs, t)
	}
	return tabs
}

func (m *Manager) CloseTab(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tabs[id]; !ok {
		return fmt.Errorf("tab not found: %s", id)
	}
	delete(m.tabs, id)
	return nil
}
