package bus

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
)

type Definition struct {
	Type       string
	Properties any
}

type Payload struct {
	Type       string
	Properties any
}

type EventHandler func(payload Payload)

type State struct {
	mu       sync.RWMutex
	wildcard []GlobalHandler
	typed    map[string][]Handler
	patterns []PatternHandler
	closed   bool
	nextID   int32
}

type Handler struct {
	id        int
	eventType string
	callback  EventHandler
}

type GlobalHandler struct {
	id       int
	callback func(GlobalEvent)
}

type PatternHandler struct {
	id       int
	pattern  string
	callback func(GlobalEvent)
}

type GlobalEvent struct {
	SessionID string
	Directory string
	Project   string
	Workspace string
	Payload   Payload
}

var (
	globalState       *State
	globalStateOnce    sync.Once
	handlerID          int32
	globalHandlerID    int32
	patternHandlerID   int32
)

func getState() *State {
	globalStateOnce.Do(func() {
		globalState = &State{
			typed:    make(map[string][]Handler),
			patterns: make([]PatternHandler, 0),
		}
	})
	if globalState.closed {
		globalState.mu.Lock()
		globalState.closed = false
		globalState.typed = make(map[string][]Handler)
		globalState.wildcard = nil
		globalState.patterns = make([]PatternHandler, 0)
		globalState.mu.Unlock()
	}
	return globalState
}

func Publish(ctx context.Context, def Definition, properties any) error {
	return publishTo(getState(), "", def, properties)
}

func PublishWithSession(ctx context.Context, sessionID string, def Definition, properties any) error {
	return publishTo(getState(), sessionID, def, properties)
}

func publishTo(s *State, sessionID string, def Definition, properties any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	payload := Payload{Type: def.Type, Properties: properties}
	globalEvent := GlobalEvent{
		SessionID: sessionID,
		Payload:   payload,
	}

	for _, h := range s.typed[def.Type] {
		h.callback(payload)
	}

	for _, h := range s.wildcard {
		h.callback(globalEvent)
	}

	for _, h := range s.patterns {
		if matchPattern(def.Type, h.pattern) {
			h.callback(globalEvent)
		}
	}

	return nil
}

func matchPattern(topic, pattern string) bool {
	topicParts := strings.Split(topic, ".")
	patternParts := strings.Split(pattern, ".")

	if len(topicParts) != len(patternParts) {
		return false
	}

	for i := range topicParts {
		if patternParts[i] == "*" {
			continue
		}
		if topicParts[i] != patternParts[i] {
			return false
		}
	}
	return true
}

func Subscribe(def Definition, callback EventHandler) func() {
	return subscribeTo(getState(), def, callback)
}

func subscribeTo(s *State, def Definition, callback EventHandler) func() {
	id := atomic.AddInt32(&handlerID, 1)

	s.mu.Lock()
	defer s.mu.Unlock()

	h := Handler{id: int(id), eventType: def.Type, callback: callback}
	s.typed[def.Type] = append(s.typed[def.Type], h)

	return func() {
		unsubscribe(s, def.Type, int(id))
	}
}

func unsubscribe(s *State, eventType string, id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	handlers := s.typed[eventType]
	for i, h := range handlers {
		if h.id == id {
			s.typed[eventType] = append(handlers[:i], handlers[i+1:]...)
			return
		}
	}
}

func SubscribeAll(callback func(event GlobalEvent)) func() {
	return subscribeAllTo(getState(), callback)
}

func subscribeAllTo(s *State, callback func(event GlobalEvent)) func() {
	id := atomic.AddInt32(&globalHandlerID, 1)

	s.mu.Lock()
	defer s.mu.Unlock()

	handler := GlobalHandler{id: int(id), callback: callback}
	s.wildcard = append(s.wildcard, handler)

	return func() {
		unsubscribeAll(s, int(id))
	}
}

func unsubscribeAll(s *State, id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, h := range s.wildcard {
		if h.id == id {
			s.wildcard = append(s.wildcard[:i], s.wildcard[i+1:]...)
			return
		}
	}
}

func SubscribePattern(pattern string, callback func(GlobalEvent)) func() {
	return subscribePatternTo(getState(), pattern, callback)
}

func subscribePatternTo(s *State, pattern string, callback func(GlobalEvent)) func() {
	id := atomic.AddInt32(&patternHandlerID, 1)

	s.mu.Lock()
	defer s.mu.Unlock()

	h := PatternHandler{id: int(id), pattern: pattern, callback: callback}
	s.patterns = append(s.patterns, h)

	return func() {
		unsubscribePattern(s, int(id))
	}
}

func unsubscribePattern(s *State, id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, h := range s.patterns {
		if h.id == id {
			s.patterns = append(s.patterns[:i], s.patterns[i+1:]...)
			return
		}
	}
}

func Define(eventType string, properties any) Definition {
	return Definition{Type: eventType, Properties: properties}
}

func Close() {
	s := getState()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	s.typed = make(map[string][]Handler)
	s.wildcard = nil
	s.patterns = make([]PatternHandler, 0)
}

func Reset() {
	s := getState()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = false
	s.typed = make(map[string][]Handler)
	s.wildcard = nil
	s.patterns = make([]PatternHandler, 0)
}
