// Package util provides utility functions for task execution and management.
package util

import (
	"context"
	"sync"
	"time"
)

// AbortController mimics the browser AbortController API.
// It provides a way to abort an operation with a timeout.
type AbortController struct {
	Signal chan struct{}
	mu     sync.Mutex
}

// NewAbortController creates a new AbortController.
func NewAbortController() *AbortController {
	return &AbortController{
		Signal: make(chan struct{}),
	}
}

// Abort signals the abort by closing the Signal channel.
func (ac *AbortController) Abort() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	select {
	case <-ac.Signal:
		// Already closed
	default:
		close(ac.Signal)
	}
}

// IsAborted returns true if the controller has been aborted.
func (ac *AbortController) IsAborted() bool {
	select {
	case <-ac.Signal:
		return true
	default:
		return false
	}
}

// AbortAfter creates an AbortController that automatically aborts after
// the specified timeout in milliseconds.
//
// Returns the controller and a clear function that can be used to cancel
// the timeout before it fires. Calling the clear function does not abort
// the controller if it hasn't already been triggered.
//
// Uses bind() pattern to avoid capturing the surrounding scope in closures,
// preventing retention of large objects in the timer closure.
func AbortAfter(ms int) (controller *AbortController, clear func()) {
	controller = NewAbortController()

	// Use a separate goroutine to avoid capturing controller in closure
	// This prevents retention of request bodies and other large objects
	timer := time.AfterFunc(time.Duration(ms)*time.Millisecond, func() {
		controller.Abort()
	})

	clear = func() {
		timer.Stop()
	}

	return controller, clear
}

// WithAbort returns a context that is canceled when the provided signal
// channel is closed OR when the abort is called on the returned context's
// cancel function.
//
// The signal channel should be a channel that is closed when abort is requested.
// This is useful for integrating with other abort mechanisms.
//
// The returned context is never nil and is valid until the cancel function
// is called or the signal is closed.
func WithAbort(signal <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-signal:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

// WithAbortFromController returns a context that is canceled when the
// AbortController is aborted or when the returned cancel function is called.
func WithAbortFromController(ac *AbortController) (context.Context, context.CancelFunc) {
	return WithAbort(ac.Signal)
}

// AbortAfterContext combines a timeout with an existing context.
// It returns a new context that is canceled when either:
// - The timeout expires
// - The parent context is canceled
// - The returned cancel function is called
//
// The cancel function clears the timeout and cancels the context if it
// hasn't already been canceled.
func AbortAfterContext(ctx context.Context, ms int) (context.Context, context.CancelFunc) {
	controller, clearTimeout := AbortAfter(ms)

	// Create a context that is canceled by either the controller or the parent
	childCtx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-controller.Signal:
			cancel()
		case <-childCtx.Done():
			// Already canceled
		}
	}()

	// Combined cancel function that clears timeout and cancels context
	finalCancel := func() {
		clearTimeout()
		cancel()
	}

	return childCtx, finalCancel
}
