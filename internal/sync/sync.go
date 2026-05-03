package sync

import (
	"context"
	"sync"
	"time"
)

type Session struct {
	ID      string
	Data    map[string]any
	Updated time.Time
}

type Store interface {
	Get(ctx context.Context, id string) (*Session, error)
	Set(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Session, error)
}

type MemoryStore struct {
	mu  sync.RWMutex
	m   map[string]*Session
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{m: make(map[string]*Session)}
}

func (s *MemoryStore) Get(ctx context.Context, id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if session, ok := s.m[id]; ok {
		return session, nil
	}
	return nil, nil
}

func (s *MemoryStore) Set(ctx context.Context, session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	session.Updated = time.Now()
	s.m[session.ID] = session
	return nil
}

func (s *MemoryStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, id)
	return nil
}

func (s *MemoryStore) List(ctx context.Context) ([]*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var sessions []*Session
	for _, session := range s.m {
		sessions = append(sessions, session)
	}
	return sessions, nil
}
