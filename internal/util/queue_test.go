package util

import (
	"sync"
	"testing"
	"time"
)

func TestAsyncQueue_Basic(t *testing.T) {
	q := NewAsyncQueue[int]()

	if q.Len() != 0 {
		t.Errorf("expected len 0, got %d", q.Len())
	}

	q.Push(1)
	if q.Len() != 1 {
		t.Errorf("expected len 1, got %d", q.Len())
	}

	val, ok := q.TryNext()
	if !ok || val != 1 {
		t.Errorf("expected 1, got %d, ok=%v", val, ok)
	}

	if q.Len() != 0 {
		t.Errorf("expected len 0, got %d", q.Len())
	}
}

func TestAsyncQueue_TryNext_Empty(t *testing.T) {
	q := NewAsyncQueue[int]()

	val, ok := q.TryNext()
	if ok {
		t.Errorf("expected ok=false for empty queue, got ok=%v val=%d", ok, val)
	}
}

func TestAsyncQueue_Next_Blocks(t *testing.T) {
	q := NewAsyncQueue[int]()

	var received int
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		received = q.Next()
	}()

	time.Sleep(10 * time.Millisecond)

	if received != 0 {
		t.Errorf("expected 0 (blocked), got %d", received)
	}

	q.Push(42)

	wg.Wait()
	if received != 42 {
		t.Errorf("expected 42, got %d", received)
	}
}

func TestAsyncQueue_PushWhileBlocked(t *testing.T) {
	q := NewAsyncQueue[int]()

	var received int
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		received = q.Next()
	}()

	q.Push(10)
	q.Push(20)

	wg.Wait()
	if received != 10 {
		t.Errorf("expected first value 10, got %d", received)
	}
}

func TestAsyncQueue_FIFO(t *testing.T) {
	q := NewAsyncQueue[int]()

	q.Push(1)
	q.Push(2)
	q.Push(3)

	if v := q.Next(); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	if v := q.Next(); v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
	if v := q.Next(); v != 3 {
		t.Errorf("expected 3, got %d", v)
	}
}

func TestAsyncQueue_Concurrent(t *testing.T) {
	q := NewAsyncQueue[int]()

	const n = 1000
	var received int
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			q.Push(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			val := q.Next()
			mu.Lock()
			received += val
			mu.Unlock()
		}
	}()

	wg.Wait()

	expected := n*(n-1) / 2
	if received != expected {
		t.Errorf("expected sum %d, got %d", expected, received)
	}
}

func TestWork_Success(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	var sum int
	var mu sync.Mutex

	err := Work(3, items, func(item int) error {
		mu.Lock()
		sum += item
		mu.Unlock()
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := 15
	if sum != expected {
		t.Errorf("expected sum %d, got %d", expected, sum)
	}
}

func TestWork_Error(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	err := Work(3, items, func(item int) error {
		if item == 3 {
			return &testError{value: 3}
		}
		return nil
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	testErr, ok := err.(*testError)
	if !ok {
		t.Errorf("expected *testError, got %T", err)
	} else if testErr.value != 3 {
		t.Errorf("expected error value 3, got %d", testErr.value)
	}
}

type testError struct {
	value int
}

func (e *testError) Error() string {
	return "test error"
}

func TestWork_EmptyItems(t *testing.T) {
	items := []int{}
	var called bool

	err := Work(3, items, func(item int) error {
		called = true
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if called {
		t.Errorf("fn should not be called for empty items")
	}
}

func TestWork_ZeroConcurrency(t *testing.T) {
	items := []int{1, 2, 3}
	var count int

	err := Work(0, items, func(item int) error {
		count++
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 calls, got %d", count)
	}
}

func TestWork_SingleConcurrency(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	var order []int
	var mu sync.Mutex

	err := Work(1, items, func(item int) error {
		mu.Lock()
		order = append(order, item)
		mu.Unlock()
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(order) != 5 {
		t.Errorf("expected 5 items, got %d", len(order))
	}
}

func TestWork_LargeConcurrency(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	var count int
	var mu sync.Mutex

	err := Work(100, items, func(item int) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 5 {
		t.Errorf("expected 5, got %d", count)
	}
}

func TestWorkWithContext_Success(t *testing.T) {
	items := []int{1, 2, 3}

	err := WorkWithContext(2, items, func(item int) error {
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkWithContext_Error(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	err := WorkWithContext(2, items, func(item int) error {
		if item == 2 {
			return &testError{value: 2}
		}
		return nil
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
