package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	dir string
}

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) SaveSession(sess *Session) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("failed to create session dir: %w", err)
	}

	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	path := filepath.Join(s.dir, sess.ID+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write session: %w", err)
	}

	return nil
}

func (s *Store) LoadSession(id string) (*Session, error) {
	path := filepath.Join(s.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read session: %w", err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &sess, nil
}

func (s *Store) DeleteSession(id string) error {
	path := filepath.Join(s.dir, id+".json")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (s *Store) ListSessions() ([]*Session, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read session dir: %w", err)
	}

	sessions := make([]*Session, 0)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		id := entry.Name()[:len(entry.Name())-5]
		sess, err := s.LoadSession(id)
		if err != nil {
			continue
		}
		sessions = append(sessions, sess)
	}

	return sessions, nil
}

// Tab management methods

func (s *Store) tabsPath() string {
	return filepath.Join(s.dir, "tabs.json")
}

func (s *Store) SaveTabs(tabs []*Tab) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("failed to create session dir: %w", err)
	}

	data, err := json.MarshalIndent(tabs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tabs: %w", err)
	}

	if err := os.WriteFile(s.tabsPath(), data, 0644); err != nil {
		return fmt.Errorf("failed to write tabs: %w", err)
	}

	return nil
}

func (s *Store) LoadTabs() ([]*Tab, error) {
	data, err := os.ReadFile(s.tabsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return []*Tab{}, nil
		}
		return nil, fmt.Errorf("failed to read tabs: %w", err)
	}

	var tabs []*Tab
	if err := json.Unmarshal(data, &tabs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tabs: %w", err)
	}

	return tabs, nil
}

func (s *Store) SaveTab(tab *Tab) error {
	tabs, err := s.LoadTabs()
	if err != nil {
		return err
	}

	// Find and update or append
	found := false
	for i, t := range tabs {
		if t.ID == tab.ID {
			tabs[i] = tab
			found = true
			break
		}
	}
	if !found {
		tabs = append(tabs, tab)
	}

	return s.SaveTabs(tabs)
}

func (s *Store) DeleteTab(id string) error {
	tabs, err := s.LoadTabs()
	if err != nil {
		return err
	}

	for i, t := range tabs {
		if t.ID == id {
			tabs = append(tabs[:i], tabs[i+1:]...)
			return s.SaveTabs(tabs)
		}
	}

	return fmt.Errorf("tab not found: %s", id)
}
