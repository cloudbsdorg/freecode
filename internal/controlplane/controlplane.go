package controlplane

import (
	"context"
	"sync"
)

type Agent struct {
	ID       string
	Endpoint string
	Status   string
}

type ControlPlane interface {
	Register(ctx context.Context, agent Agent) error
	Unregister(ctx context.Context, id string) error
	ListAgents(ctx context.Context) ([]Agent, error)
	GetAgent(ctx context.Context, id string) (*Agent, error)
}

type memoryControlPlane struct {
	mu     sync.RWMutex
	agents map[string]Agent
}

func NewMemoryControlPlane() ControlPlane {
	return &memoryControlPlane{agents: make(map[string]Agent)}
}

func (cp *memoryControlPlane) Register(ctx context.Context, agent Agent) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	agent.Status = "online"
	cp.agents[agent.ID] = agent
	return nil
}

func (cp *memoryControlPlane) Unregister(ctx context.Context, id string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	delete(cp.agents, id)
	return nil
}

func (cp *memoryControlPlane) ListAgents(ctx context.Context) ([]Agent, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	var result []Agent
	for _, a := range cp.agents {
		result = append(result, a)
	}
	return result, nil
}

func (cp *memoryControlPlane) GetAgent(ctx context.Context, id string) (*Agent, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	if a, ok := cp.agents[id]; ok {
		return &a, nil
	}
	return nil, nil
}
