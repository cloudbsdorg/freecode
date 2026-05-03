package bus

import (
	"context"
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
	mu      sync.RWMutex
	wildcard []GlobalHandler
	typed    map[string][]Handler
	closed   bool
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

var globalHandlerID int32
var handlerID int32

type GlobalEvent struct {
	Directory string
	Project   string
	Workspace string
	Payload   any
}

var (
	defaultState *State
	once        sync.Once
)

func getState() *State {
	once.Do(func() {
		defaultState = &State{
			typed: make(map[string][]Handler),
		}
	})
	if defaultState.closed {
		defaultState.mu.Lock()
		defaultState.closed = false
		defaultState.typed = make(map[string][]Handler)
		defaultState.wildcard = nil
		defaultState.mu.Unlock()
	}
	return defaultState
}

func Publish(ctx context.Context, def Definition, properties any) error {
	s := getState()
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	payload := Payload{Type: def.Type, Properties: properties}

	for _, h := range s.typed[def.Type] {
		h.callback(payload)
	}

	for _, h := range s.wildcard {
		h.callback(GlobalEvent{Payload: payload})
	}

	return nil
}

func Subscribe(def Definition, callback EventHandler) func() {
	s := getState()
	s.mu.Lock()
	defer s.mu.Unlock()

	id := atomic.AddInt32(&handlerID, 1)
	h := Handler{id: int(id), eventType: def.Type, callback: callback}
	s.typed[def.Type] = append(s.typed[def.Type], h)

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		handlers := s.typed[def.Type]
		for i, h := range handlers {
			if h.id == int(id) {
				s.typed[def.Type] = append(handlers[:i], handlers[i+1:]...)
				return
			}
		}
	}
}

func SubscribeAll(callback func(event GlobalEvent)) func() {
	s := getState()
	s.mu.Lock()
	defer s.mu.Unlock()

	id := atomic.AddInt32(&globalHandlerID, 1)
	handler := GlobalHandler{id: int(id), callback: callback}
	s.wildcard = append(s.wildcard, handler)

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		for i, h := range s.wildcard {
			if h.id == int(id) {
				s.wildcard = append(s.wildcard[:i], s.wildcard[i+1:]...)
				return
			}
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
}
