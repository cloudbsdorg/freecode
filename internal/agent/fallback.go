package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ModelFallback struct {
	mu         sync.RWMutex
	primary    string
	fallbacks  []string
	currentIdx int
	maxRetries int
	retryDelay time.Duration
	attempts   map[string]int
	disabled   map[string]bool
}

type FallbackConfig struct {
	Primary    string
	Fallbacks  []string
	MaxRetries int
	RetryDelay time.Duration
}

func NewModelFallback(cfg FallbackConfig) *ModelFallback {
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryDelay == 0 {
		cfg.RetryDelay = 5 * time.Second
	}

	return &ModelFallback{
		primary:    cfg.Primary,
		fallbacks:  cfg.Fallbacks,
		maxRetries: cfg.MaxRetries,
		retryDelay: cfg.RetryDelay,
		attempts:   make(map[string]int),
		disabled:   make(map[string]bool),
	}
}

func (f *ModelFallback) GetPrimary() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.primary
}

func (f *ModelFallback) GetCurrent() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.currentIdx == 0 {
		return f.primary
	}

	if f.currentIdx-1 < len(f.fallbacks) {
		return f.fallbacks[f.currentIdx-1]
	}

	return f.primary
}

func (f *ModelFallback) GetAll() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	all := make([]string, 0, 1+len(f.fallbacks))
	all = append(all, f.primary)
	all = append(all, f.fallbacks...)
	return all
}

func (f *ModelFallback) RecordFailure(model string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.attempts[model]++
	currentAttempts := f.attempts[model]

	if currentAttempts >= f.maxRetries {
		f.disabled[model] = true
	}
}

func (f *ModelFallback) RecordSuccess(model string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.attempts[model] = 0
	f.disabled[model] = false
	f.currentIdx = 0
}

func (f *ModelFallback) GetNextFallback() (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := f.currentIdx; i < len(f.fallbacks); i++ {
		candidate := f.fallbacks[i]
		if !f.disabled[candidate] && f.attempts[candidate] < f.maxRetries {
			f.currentIdx = i + 1
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no available fallback models")
}

func (f *ModelFallback) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.currentIdx = 0
	for k := range f.attempts {
		f.attempts[k] = 0
	}
	for k := range f.disabled {
		f.disabled[k] = false
	}
}

func (f *ModelFallback) IsModelDisabled(model string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.disabled[model]
}

func (f *ModelFallback) GetAttempts(model string) int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.attempts[model]
}

func (f *ModelFallback) ShouldRetry(model string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.attempts[model] < f.maxRetries && !f.disabled[model]
}

type FallbackHandler struct {
	mu        sync.RWMutex
	fallbacks map[string]*ModelFallback
	defaultFB *ModelFallback
}

func NewFallbackHandler() *FallbackHandler {
	return &FallbackHandler{
		fallbacks: make(map[string]*ModelFallback),
	}
}

func (h *FallbackHandler) Register(agentName string, fb *ModelFallback) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fallbacks[agentName] = fb
}

func (h *FallbackHandler) SetDefault(fb *ModelFallback) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.defaultFB = fb
}

func (h *FallbackHandler) Get(agentName string) *ModelFallback {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if fb, ok := h.fallbacks[agentName]; ok {
		return fb
	}

	return h.defaultFB
}

func (h *FallbackHandler) GetForAgent(agentName string) (*ModelFallback, error) {
	h.mu.RLock()
	fb, ok := h.fallbacks[agentName]
	h.mu.RUnlock()

	if !ok {
		fb = h.defaultFB
	}

	if fb == nil {
		return nil, fmt.Errorf("no fallback configured for agent: %s", agentName)
	}

	return fb, nil
}

func (h *FallbackHandler) Remove(agentName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.fallbacks, agentName)
}

type FallbackExecutor struct {
	fallbacks *FallbackHandler
	maxDepth  int
}

func NewFallbackExecutor(fallbacks *FallbackHandler) *FallbackExecutor {
	return &FallbackExecutor{
		fallbacks: fallbacks,
		maxDepth:  5,
	}
}

type FallbackableTask func(ctx context.Context, model string) error

func (e *FallbackExecutor) ExecuteWithFallback(ctx context.Context, agentName string, task FallbackableTask) error {
	fb, err := e.fallbacks.GetForAgent(agentName)
	if err != nil {
		return err
	}

	models := append([]string{fb.GetPrimary()}, fb.fallbacks...)
	visited := make(map[string]bool)
	depth := 0

	for {
		if depth >= e.maxDepth {
			return fmt.Errorf("max fallback depth reached")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		model := models[depth]
		if visited[model] {
			depth++
			continue
		}
		visited[model] = true

		err := task(ctx, model)
		if err == nil {
			fb.RecordSuccess(model)
			return nil
		}

		fb.RecordFailure(model)

		if depth >= len(models)-1 {
			return fmt.Errorf("all models failed, last error: %w", err)
		}

		time.Sleep(fb.retryDelay)
		depth++
	}
}

func (e *FallbackExecutor) ExecuteWithFallbackAsync(ctx context.Context, agentName string, task FallbackableTask, resultCh chan<- error) {
	go func() {
		resultCh <- e.ExecuteWithFallback(ctx, agentName, task)
	}()
}
