// Package util provides utility functions for deferred execution.
package util

import "sync"

// Deferred provides both synchronous and asynchronous disposal of a function.
// It implements the pattern of deferring cleanup actions with support for
// both sync and async disposal patterns.
type Deferred interface {
	// Dispose synchronously executes the deferred function.
	// This is safe to call multiple times; subsequent calls are no-ops.
	Dispose()

	// DisposeAsync executes the deferred function asynchronously and returns
	// any error that occurs during execution.
	// This is safe to call multiple times; subsequent calls are no-ops.
	// If fn returns a non-nil error, it will be returned by DisposeAsync.
	DisposeAsync() error
}

// deferredFunc wraps a function with optional error handling.
type deferredFunc struct {
	fn        func()
	mu        sync.Mutex
	disposed  bool
}

// Defer creates a new Deferred that will execute the provided function.
// The function is guaranteed to be called exactly once, either on Dispose()
// or DisposeAsync(). Subsequent calls to either method are no-ops.
//
// Usage:
//
//	d := Defer(func() { cleanup() })
//	defer d.Dispose()
//
// Or for async contexts:
//
//	d := Defer(func() error { return cleanup() })
//	go func() { _ = d.DisposeAsync() }()
func Defer(fn func()) Deferred {
	return &deferredFunc{fn: fn}
}

// Dispose synchronously executes the deferred function.
func (d *deferredFunc) Dispose() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.disposed {
		return
	}
	d.disposed = true
	d.fn()
}

// DisposeAsync executes the deferred function asynchronously.
// It returns any error that occurs during function execution.
func (d *deferredFunc) DisposeAsync() error {
	d.mu.Lock()
	if d.disposed {
		d.mu.Unlock()
		return nil
	}
	d.disposed = true
	fn := d.fn
	d.mu.Unlock()

	fn()
	return nil
}
