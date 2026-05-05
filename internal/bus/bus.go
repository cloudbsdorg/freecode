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
	mu        sync.RWMutex
	wildcard  []GlobalHandler
	typed     map[string][]Handler
	closed    bool
	nextID    int32
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

type GlobalEvent struct {
	Directory string
	Project   string
	Workspace string
	Payload   Payload
}

var (
	globalState     *State
	globalStateOnce sync.Once
	handlerID       int32
	globalHandlerID int32
)

func getState() *State {
	globalStateOnce.Do(func() {
		globalState = &State{
			typed: make(map[string][]Handler),
		}
	})
	if globalState.closed {
		globalState.mu.Lock()
		globalState.closed = false
		globalState.typed = make(map[string][]Handler)
		globalState.wildcard = nil
		globalState.mu.Unlock()
	}
	return globalState
}

func Publish(ctx context.Context, def Definition, properties any) error {
	return publishTo(getState(), def, properties)
}

func publishTo(s *State, def Definition, properties any) error {
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

func Reset() {
	s := getState()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = false
	s.typed = make(map[string][]Handler)
	s.wildcard = nil
}
