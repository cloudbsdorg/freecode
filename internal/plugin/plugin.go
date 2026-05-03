package plugin

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrPluginNotFound = errors.New("plugin not found")
	ErrPluginDisabled = errors.New("plugin disabled")
)

type Plugin interface {
	Name() string
	Init(ctx context.Context) error
	Close() error
}

type Registry interface {
	Register(plugin Plugin) error
	Unregister(name string) error
	Get(name string) (Plugin, error)
	List() []string
	Enable(name string) error
	Disable(name string) error
}

type memoryRegistry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
	enabled map[string]bool
}

func NewMemoryRegistry() Registry {
	return &memoryRegistry{
		plugins: make(map[string]Plugin),
		enabled: make(map[string]bool),
	}
}

func (r *memoryRegistry) Register(plugin Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[plugin.Name()] = plugin
	r.enabled[plugin.Name()] = true
	return nil
}

func (r *memoryRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.plugins, name)
	delete(r.enabled, name)
	return nil
}

func (r *memoryRegistry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	plugin, ok := r.plugins[name]
	if !ok {
		return nil, ErrPluginNotFound
	}
	if !r.enabled[name] {
		return nil, ErrPluginDisabled
	}
	return plugin, nil
}

func (r *memoryRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var names []string
	for name := range r.plugins {
		names = append(names, name)
	}
	return names
}

func (r *memoryRegistry) Enable(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.plugins[name]; !ok {
		return ErrPluginNotFound
	}
	r.enabled[name] = true
	return nil
}

func (r *memoryRegistry) Disable(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.plugins[name]; !ok {
		return ErrPluginNotFound
	}
	r.enabled[name] = false
	return nil
}
