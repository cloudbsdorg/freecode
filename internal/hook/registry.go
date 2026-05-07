package hook

import (
	"context"
	"fmt"
	"sync"
)

type Registry struct {
	mu                sync.RWMutex
	sessionHooks      map[string][]SessionHook
	toolHooks         map[string][]ToolHook
	transformHooks    []TransformHook
	continuationHooks []ContinuationHook
	ralphHooks        []RalphHook
	skillHooks        map[string][]SkillHook
}

type SessionHook func(ctx context.Context, evt SessionEvent) error
type ToolHook func(ctx context.Context, evt ToolEvent) (error, bool)
type TransformHook func(msg *Message) (*Message, error)
type ContinuationHook func(ctx context.Context, session *SessionData) (*ContinueSignal, error)
type RalphHook func(ctx context.Context, input string) (string, error)
type SkillHook func(ctx context.Context, skill string, args map[string]interface{}) error

type SessionEvent struct {
	Type      string
	SessionID string
	Data      map[string]interface{}
}

type ToolEvent struct {
	Type      string
	ToolName  string
	SessionID string
	Arguments map[string]interface{}
	Result    interface{}
	Error     error
}

type Message struct {
	Role    string
	Content string
}

type SessionData struct {
	ID       string
	Messages []Message
}

type ContinueSignal struct {
	Continue bool
	Reason   string
}

func NewRegistry() *Registry {
	return &Registry{
		sessionHooks:      make(map[string][]SessionHook),
		toolHooks:         make(map[string][]ToolHook),
		transformHooks:    make([]TransformHook, 0),
		continuationHooks: make([]ContinuationHook, 0),
		ralphHooks:        make([]RalphHook, 0),
		skillHooks:        make(map[string][]SkillHook),
	}
}

func (r *Registry) RegisterSessionHook(eventType string, hook SessionHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessionHooks[eventType] = append(r.sessionHooks[eventType], hook)
}

func (r *Registry) RegisterToolHook(eventType string, hook ToolHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.toolHooks[eventType] = append(r.toolHooks[eventType], hook)
}

func (r *Registry) RegisterTransformHook(hook TransformHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transformHooks = append(r.transformHooks, hook)
}

func (r *Registry) RegisterContinuationHook(hook ContinuationHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.continuationHooks = append(r.continuationHooks, hook)
}

func (r *Registry) RegisterRalphHook(hook RalphHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ralphHooks = append(r.ralphHooks, hook)
}

func (r *Registry) RegisterSkillHook(skill string, hook SkillHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skillHooks[skill] = append(r.skillHooks[skill], hook)
}

func (r *Registry) EmitSessionEvent(ctx context.Context, eventType string, sessionID string, data map[string]interface{}) error {
	r.mu.RLock()
	hooks := r.sessionHooks[eventType]
	r.mu.RUnlock()

	if len(hooks) == 0 {
		return nil
	}

	evt := SessionEvent{
		Type:      eventType,
		SessionID: sessionID,
		Data:      data,
	}

	for _, hook := range hooks {
		if err := hook(ctx, evt); err != nil {
			return fmt.Errorf("session hook error (%s): %w", eventType, err)
		}
	}
	return nil
}

func (r *Registry) EmitToolEvent(ctx context.Context, eventType, toolName, sessionID string, args map[string]interface{}) (error, bool) {
	r.mu.RLock()
	hooks := r.toolHooks[eventType]
	r.mu.RUnlock()

	if len(hooks) == 0 {
		return nil, false
	}

	evt := ToolEvent{
		Type:      eventType,
		ToolName:  toolName,
		SessionID: sessionID,
		Arguments: args,
	}

	for _, hook := range hooks {
		if err, handled := hook(ctx, evt); handled {
			return err, true
		}
	}
	return nil, false
}

func (r *Registry) ApplyTransformHooks(msg *Message) (*Message, error) {
	r.mu.RLock()
	hooks := r.transformHooks
	r.mu.RUnlock()

	result := msg
	var err error
	for _, hook := range hooks {
		result, err = hook(result)
		if err != nil {
			return nil, fmt.Errorf("transform hook error: %w", err)
		}
	}
	return result, nil
}

func (r *Registry) CheckContinuation(ctx context.Context, session *SessionData) (*ContinueSignal, error) {
	r.mu.RLock()
	hooks := r.continuationHooks
	r.mu.RUnlock()

	for _, hook := range hooks {
		signal, err := hook(ctx, session)
		if err != nil {
			return nil, fmt.Errorf("continuation hook error: %w", err)
		}
		if signal != nil && !signal.Continue {
			return signal, nil
		}
	}
	return &ContinueSignal{Continue: true}, nil
}

func (r *Registry) ApplyRalphHooks(ctx context.Context, input string) (string, error) {
	r.mu.RLock()
	hooks := r.ralphHooks
	r.mu.RUnlock()

	result := input
	var err error
	for _, hook := range hooks {
		result, err = hook(ctx, result)
		if err != nil {
			return "", fmt.Errorf("ralph hook error: %w", err)
		}
	}
	return result, nil
}
