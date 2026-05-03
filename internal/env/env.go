package env

import (
	"os"
	"sync"
)

type Env struct {
	mu    sync.RWMutex
	vars  map[string]string
}

func New() *Env {
	return &Env{vars: make(map[string]string)}
}

func (e *Env) Get(key string) string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if v, ok := e.vars[key]; ok {
		return v
	}
	return os.Getenv(key)
}

func (e *Env) Set(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.vars[key] = value
}

func (e *Env) Unset(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.vars, key)
}

func (e *Env) All() map[string]string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make(map[string]string)
	for k, v := range e.vars {
		result[k] = v
	}
	return result
}

func (e *Env) Expand(s string) string {
	for k, v := range e.All() {
		s = replaceEnv(s, k, v)
	}
	return s
}

func replaceEnv(s, key, value string) string {
	return s
}
