package util

import (
	"sync"
)

// AsyncQueue is a thread-safe async queue for Go.
// It supports blocking and non-blocking receive operations.
type AsyncQueue[T any] struct {
	queue    []T
	resolvers []chan T
	mu       sync.Mutex
}

// NewAsyncQueue creates a new AsyncQueue.
func NewAsyncQueue[T any]() *AsyncQueue[T] {
	return &AsyncQueue[T]{}
}

// Push adds an item to the queue. If there's a waiting receiver,
// it will be delivered immediately; otherwise, the item is queued.
func (q *AsyncQueue[T]) Push(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.resolvers) > 0 {
		resolver := q.resolvers[0]
		q.resolvers = q.resolvers[1:]
		resolver <- item
		return
	}

	q.queue = append(q.queue, item)
}

// Next returns the next item from the queue, blocking until an item is available.
func (q *AsyncQueue[T]) Next() T {
	q.mu.Lock()

	if len(q.queue) > 0 {
		item := q.queue[0]
		q.queue = q.queue[1:]
		q.mu.Unlock()
		return item
	}

	resolver := make(chan T)
	q.resolvers = append(q.resolvers, resolver)
	q.mu.Unlock()

	return <-resolver
}

// TryNext attempts to return the next item from the queue without blocking.
// It returns the item and true if available, or the zero value and false otherwise.
func (q *AsyncQueue[T]) TryNext() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) > 0 {
		item := q.queue[0]
		q.queue = q.queue[1:]
		return item, true
	}

	var zero T
	return zero, false
}

// Len returns the current number of items in the queue (non-blocking receivers not counted).
func (q *AsyncQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}

// Work runs concurrent workers processing items from the slice.
// It spawns `concurrency` goroutines that process items concurrently.
// Returns the first error encountered, or nil if all items processed successfully.
func Work[T any](concurrency int, items []T, fn func(T) error) error {
	if concurrency <= 0 {
		concurrency = 1
	}

	if len(items) == 0 {
		return nil
	}

	pending := make([]T, len(items))
	copy(pending, items)

	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				item, ok := func() (T, bool) {
					mu.Lock()
					defer mu.Unlock()
					if len(pending) > 0 {
						item := pending[len(pending)-1]
						pending = pending[:len(pending)-1]
						return item, true
					}
					var zero T
					return zero, false
				}()
				if !ok {
					return
				}
				if err := fn(item); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	return nil
}

func WorkWithContext[T any](concurrency int, items []T, fn func(T) error) error {
	if concurrency <= 0 {
		concurrency = 1
	}

	if len(items) == 0 {
		return nil
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, item := range items {
		item := item
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			if err := fn(item); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}
