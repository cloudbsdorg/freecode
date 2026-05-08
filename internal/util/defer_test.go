package util

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestDefer_Dispose(t *testing.T) {
	called := atomic.Int32{}

	d := Defer(func() {
		called.Add(1)
	})

	d.Dispose()

	if called.Load() != 1 {
		t.Errorf("expected fn to be called once, got %d", called.Load())
	}
}

func TestDefer_Dispose_Idempotent(t *testing.T) {
	called := atomic.Int32{}

	d := Defer(func() {
		called.Add(1)
	})

	d.Dispose()
	d.Dispose()
	d.Dispose()

	if called.Load() != 1 {
		t.Errorf("expected fn to be called exactly once, got %d", called.Load())
	}
}

func TestDefer_DisposeAsync(t *testing.T) {
	called := atomic.Int32{}

	d := Defer(func() {
		called.Add(1)
	})

	err := d.DisposeAsync()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if called.Load() != 1 {
		t.Errorf("expected fn to be called once, got %d", called.Load())
	}
}

func TestDefer_DisposeAsync_Idempotent(t *testing.T) {
	called := atomic.Int32{}

	d := Defer(func() {
		called.Add(1)
	})

	_ = d.DisposeAsync()
	_ = d.DisposeAsync()
	_ = d.DisposeAsync()

	if called.Load() != 1 {
		t.Errorf("expected fn to be called exactly once, got %d", called.Load())
	}
}

func TestDefer_DisposeAfterDisposeAsync(t *testing.T) {
	called := atomic.Int32{}

	d := Defer(func() {
		called.Add(1)
	})

	_ = d.DisposeAsync()
	d.Dispose()

	if called.Load() != 1 {
		t.Errorf("expected fn to be called exactly once, got %d", called.Load())
	}
}

func TestDefer_DisposeAsyncAfterDispose(t *testing.T) {
	called := atomic.Int32{}

	d := Defer(func() {
		called.Add(1)
	})

	d.Dispose()
	err := d.DisposeAsync()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if called.Load() != 1 {
		t.Errorf("expected fn to be called exactly once, got %d", called.Load())
	}
}

func TestDefer_ConcurrentAccess(t *testing.T) {
	called := atomic.Int32{}

	d := Defer(func() {
		called.Add(1)
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			d.Dispose()
		}()
		go func() {
			defer wg.Done()
			_ = d.DisposeAsync()
		}()
	}
	wg.Wait()

	if called.Load() != 1 {
		t.Errorf("expected fn to be called exactly once, got %d", called.Load())
	}
}
