package sync

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Version represents a single version number for a session
type Version int64

// VersionVector tracks versions across multiple nodes/clients
type VersionVector map[string]Version

// ConflictResolution defines strategies for resolving conflicts
type ConflictResolution int

const (
	// LastWriteWins uses the most recently updated version (default)
	LastWriteWins ConflictResolution = iota
	// Merge combines changes from concurrent modifications
	Merge
	// ClientWins client version always wins over server
	ClientWins
	// ServerWins server version always wins over client
	ServerWins
)

// ConflictError represents a conflict detected during sync
type ConflictError struct {
	SessionID  string
	ClientVersion VersionVector
	ServerVersion VersionVector
	ClientData   map[string]any
	ServerData   map[string]any
}

func (e *ConflictError) Error() string {
	return "conflict detected for session " + e.SessionID
}

// Session represents a syncable session with version tracking
type Session struct {
	ID      string
	Data    map[string]any
	Updated time.Time
	Version Version       // Monotonically increasing version
	Vector  VersionVector // Version vector for distributed tracking
}

// Store interface for session storage
type Store interface {
	Get(ctx context.Context, id string) (*Session, error)
	Set(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Session, error)
}

// SyncStore extends Store with sync and conflict resolution capabilities
type SyncStore interface {
	Store
	// Sync attempts to sync a session, returning conflict error if detected and strategy is Merge
	Sync(ctx context.Context, session *Session, resolution ConflictResolution) (*Session, error)
}

// MemoryStore implements Store with version tracking
type MemoryStore struct {
	mu sync.RWMutex
	m  map[string]*Session
}

// NewMemoryStore creates a new MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{m: make(map[string]*Session)}
}

// internalSession wraps Session with version metadata for conflict detection
type internalSession struct {
	*Session
	vector VersionVector
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
	if session.Version == 0 {
		session.Version = 1
	} else {
		session.Version++
	}
	if session.Vector == nil {
		session.Vector = make(VersionVector)
	}
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

// Sync attempts to sync a session with conflict detection and resolution
func (s *MemoryStore) Sync(ctx context.Context, session *Session, resolution ConflictResolution) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session == nil || session.ID == "" {
		return nil, errors.New("invalid session")
	}

	existing, ok := s.m[session.ID]

	// New session - just insert
	if !ok {
		session.Updated = time.Now()
		session.Version = 1
		if session.Vector == nil {
			session.Vector = make(VersionVector)
		}
		s.m[session.ID] = session
		return session, nil
	}

	// Check for conflict using version vector
	if hasConflict(existing.Vector, session.Vector) {
		if resolution == Merge {
			return s.resolveMerge(existing, session)
		}
		if resolution == ClientWins {
			session.Version = existing.Version + 1
			session.Updated = time.Now()
			s.m[session.ID] = session
			return session, nil
		}
		if resolution == ServerWins {
			return existing, &ConflictError{
				SessionID:     session.ID,
				ClientVersion: session.Vector,
				ServerVersion: existing.Vector,
				ClientData:    session.Data,
				ServerData:    existing.Data,
			}
		}
		// LastWriteWins - compare Updated timestamps
		if session.Updated.After(existing.Updated) {
			session.Version = existing.Version + 1
			session.Updated = time.Now()
			s.m[session.ID] = session
			return session, nil
		}
		return existing, nil
	}

	// No conflict - update
	session.Version = existing.Version + 1
	session.Updated = time.Now()
	s.m[session.ID] = session
	return session, nil
}

// hasConflict detects if two version vectors indicate concurrent modifications
// A conflict occurs when client has changes that server doesn't know about
// (client version > server version for some key they both know)
func hasConflict(client, server VersionVector) bool {
	if client == nil || server == nil {
		return false
	}

	for k, clientV := range client {
		serverV := server[k]
		if clientV > serverV {
			return true
		}
	}

	return false
}

// resolveMerge combines changes from client and server
func (s *MemoryStore) resolveMerge(existing, client *Session) (*Session, error) {
	merged := &Session{
		ID:      existing.ID,
		Data:    make(map[string]any),
		Updated: time.Now(),
		Version: max(existing.Version, client.Version) + 1,
		Vector:  make(VersionVector),
	}

	// Copy server data first
	for k, v := range existing.Data {
		merged.Data[k] = v
	}

	// Merge client data - client wins on key conflicts
	for k, v := range client.Data {
		merged.Data[k] = v
	}

	// Merge vectors - take max of each component
	allKeys := make(map[string]struct{})
	for k := range existing.Vector {
		allKeys[k] = struct{}{}
	}
	for k := range client.Vector {
		allKeys[k] = struct{}{}
	}
	for k := range allKeys {
		merged.Vector[k] = max(existing.Vector[k], client.Vector[k])
	}

	s.m[merged.ID] = merged
	return merged, nil
}

func max(a, b Version) Version {
	if a > b {
		return a
	}
	return b
}

// GetVersion returns the current version of a session
func (s *MemoryStore) GetVersion(ctx context.Context, id string) (Version, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if session, ok := s.m[id]; ok {
		return session.Version, nil
	}
	return 0, nil
}

// GetVector returns the version vector of a session
func (s *MemoryStore) GetVector(ctx context.Context, id string) (VersionVector, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if session, ok := s.m[id]; ok {
		return session.Vector, nil
	}
	return nil, nil
}