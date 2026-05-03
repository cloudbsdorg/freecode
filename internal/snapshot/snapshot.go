package snapshot

import (
	"context"
	"fmt"
	"time"
)

type Snapshot struct {
	ID        string
	CreatedAt time.Time
	Data      map[string]any
}

type Store interface {
	Create(ctx context.Context, data map[string]any) (*Snapshot, error)
	Get(ctx context.Context, id string) (*Snapshot, error)
	List(ctx context.Context) ([]*Snapshot, error)
	Delete(ctx context.Context, id string) error
}

type memoryStore struct {
	snapshots map[string]*Snapshot
}

func NewMemoryStore() *memoryStore {
	return &memoryStore{snapshots: make(map[string]*Snapshot)}
}

func (s *memoryStore) Create(ctx context.Context, data map[string]any) (*Snapshot, error) {
	snap := &Snapshot{
		ID:        fmt.Sprintf("snap-%d", time.Now().UnixNano()),
		CreatedAt: time.Now(),
		Data:      data,
	}
	s.snapshots[snap.ID] = snap
	return snap, nil
}

func (s *memoryStore) Get(ctx context.Context, id string) (*Snapshot, error) {
	if snap, ok := s.snapshots[id]; ok {
		return snap, nil
	}
	return nil, nil
}

func (s *memoryStore) List(ctx context.Context) ([]*Snapshot, error) {
	var snaps []*Snapshot
	for _, snap := range s.snapshots {
		snaps = append(snaps, snap)
	}
	return snaps, nil
}

func (s *memoryStore) Delete(ctx context.Context, id string) error {
	delete(s.snapshots, id)
	return nil
}

var _ = context.Background
