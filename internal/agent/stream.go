package agent

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type StreamHandler struct {
	mu      sync.RWMutex
	outputs map[string]chan Message
}

func NewStreamHandler() *StreamHandler {
	return &StreamHandler{
		outputs: make(map[string]chan Message),
	}
}

func (sh *StreamHandler) CreateStream(sessionID string) chan Message {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	ch := make(chan Message, 100)
	sh.outputs[sessionID] = ch
	return ch
}

func (sh *StreamHandler) CloseStream(sessionID string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if ch, ok := sh.outputs[sessionID]; ok {
		close(ch)
		delete(sh.outputs, sessionID)
	}
}

func (sh *StreamHandler) Send(sessionID string, msg Message) error {
	sh.mu.RLock()
	ch, ok := sh.outputs[sessionID]
	sh.mu.RUnlock()

	if !ok {
		return fmt.Errorf("stream not found for session: %s", sessionID)
	}

	select {
	case ch <- msg:
		return nil
	default:
		return fmt.Errorf("stream buffer full")
	}
}

type StreamResponse struct {
	SessionID string
	Delta     string
	Final     bool
}

func (sh *StreamHandler) Stream(ctx context.Context, sessionID string, w io.Writer) error {
	sh.mu.RLock()
	ch, ok := sh.outputs[sessionID]
	sh.mu.RUnlock()

	if !ok {
		return fmt.Errorf("stream not found: %s", sessionID)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", msg.Content); err != nil {
				return err
			}
		}
	}
}
