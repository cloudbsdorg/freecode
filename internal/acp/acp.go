package acp

import (
	"context"
)

type Session struct {
	ID    string
	Agent AgentState
}

type AgentState struct {
	Ready     bool
	Paused    bool
	Stopped   bool
	LastError string
}

type ACP interface {
	CreateSession(ctx context.Context) (*Session, error)
	GetSession(ctx context.Context, id string) (*Session, error)
	PauseSession(ctx context.Context, id string) error
	ResumeSession(ctx context.Context, id string) error
	StopSession(ctx context.Context, id string) error
}

type memoryACP struct {
	sessions map[string]*Session
}

func NewMemoryACP() ACP {
	return &memoryACP{sessions: make(map[string]*Session)}
}

func (a *memoryACP) CreateSession(ctx context.Context) (*Session, error) {
	s := &Session{ID: "session-1", Agent: AgentState{Ready: true}}
	a.sessions[s.ID] = s
	return s, nil
}

func (a *memoryACP) GetSession(ctx context.Context, id string) (*Session, error) {
	if s, ok := a.sessions[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (a *memoryACP) PauseSession(ctx context.Context, id string) error {
	if s, ok := a.sessions[id]; ok {
		s.Agent.Paused = true
	}
	return nil
}

func (a *memoryACP) ResumeSession(ctx context.Context, id string) error {
	if s, ok := a.sessions[id]; ok {
		s.Agent.Paused = false
	}
	return nil
}

func (a *memoryACP) StopSession(ctx context.Context, id string) error {
	if s, ok := a.sessions[id]; ok {
		s.Agent.Stopped = true
	}
	return nil
}
