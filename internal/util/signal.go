package util

import (
	"sync"
)

// Trigger provides a one-shot signal mechanism.
// It can be triggered once (first trigger wins, subsequent are no-ops)
// and waited on via a channel that closes when triggered.
type Trigger struct {
	trigger func()
	wait    func() <-chan struct{}
}

type signalState struct {
	mu        sync.Mutex
	triggered bool
	ch        chan struct{}
}

// NewSignal creates a new Trigger.
// The returned Trigger can be triggered once and waited on via its Wait() method.
// After the first trigger, subsequent triggers are no-ops.
func NewSignal() Trigger {
	s := &signalState{
		ch: make(chan struct{}),
	}

	return Trigger{
		trigger: func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			if !s.triggered {
				s.triggered = true
				close(s.ch)
			}
		},
		wait: func() <-chan struct{} {
			return s.ch
		},
	}
}

// Trigger signals the trigger. It is safe to call from multiple goroutines.
// The first call triggers the signal; subsequent calls are no-ops.
func (t Trigger) Trigger() {
	t.trigger()
}

// Wait returns a channel that is closed when the trigger is signaled.
// The returned channel is already closed if the trigger was already fired.
func (t Trigger) Wait() <-chan struct{} {
	return t.wait()
}
