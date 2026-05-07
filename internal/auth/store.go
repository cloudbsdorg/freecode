package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type AuthType string

const (
	AuthTypeAPI       AuthType = "api"
	AuthTypeOAuth     AuthType = "oauth"
	AuthTypeWellKnown AuthType = "wellknown"
)

type Info struct {
	Type     AuthType       `json:"type"`
	Key      string         `json:"key,omitempty"`
	Access   string         `json:"access,omitempty"`
	Refresh  string         `json:"refresh,omitempty"`
	Expires  int64          `json:"expires,omitempty"`
	Token    string         `json:"token,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Store struct {
	path string
	data map[string]Info
	mu   sync.RWMutex
}

var (
	globalStore *Store
	once        sync.Once
)

func DefaultPath() string {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = os.TempDir()
	}
	return filepath.Join(homeDir, ".config", "freecode", "auth.json")
}

func Default() *Store {
	once.Do(func() {
		path := DefaultPath()
		globalStore = NewStore(path)
	})
	return globalStore
}

func NewStore(path string) *Store {
	s := &Store{
		path: path,
		data: make(map[string]Info),
	}
	s.load()
	return s
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read auth file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &s.data); err != nil {
		return fmt.Errorf("failed to parse auth file: %w", err)
	}

	return nil
}

func (s *Store) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create auth directory: %w", err)
	}

	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal auth data: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %w", err)
	}

	return nil
}

func (s *Store) Get(provider string) (Info, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.data[provider]
	return info, ok
}

func (s *Store) Set(provider string, info Info) error {
	s.mu.Lock()
	s.data[provider] = info
	s.mu.Unlock()
	return s.save()
}

func (s *Store) Remove(provider string) error {
	s.mu.Lock()
	delete(s.data, provider)
	s.mu.Unlock()
	return s.save()
}

func (s *Store) All() map[string]Info {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]Info, len(s.data))
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

func (s *Store) Clear() error {
	s.mu.Lock()
	s.data = make(map[string]Info)
	s.mu.Unlock()
	return s.save()
}
