package util

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithTimeout_Success(t *testing.T) {
	result, err := WithTimeout(func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 42, nil
	}, 100)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("expected 42, got %v", result)
	}
}

func TestWithTimeout_Timeout(t *testing.T) {
	_, err := WithTimeout(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 42, nil
	}, 50)

	if err == nil {
		t.Error("expected timeout error, got nil")
	}
	if err.Error() != "Operation timed out after 50ms" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWithTimeout_PromiseError(t *testing.T) {
	promiseErr := errors.New("promise failed")
	_, err := WithTimeout(func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 0, promiseErr
	}, 100)

	if err != promiseErr {
		t.Errorf("expected promise error, got %v", err)
	}
}

func TestWithTimeout_ZeroMs(t *testing.T) {
	result, err := WithTimeout(func() (int, error) {
		return 99, nil
	}, 0)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != 99 {
		t.Errorf("expected 99, got %v", result)
	}
}

func TestWithTimeout_NegativeMs(t *testing.T) {
	result, err := WithTimeout(func() (int, error) {
		return 77, nil
	}, -50)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != 77 {
		t.Errorf("expected 77, got %v", result)
	}
}

func TestWithTimeout_NoLeak(t *testing.T) {
	done := make(chan struct{})

	go func() {
		WithTimeout(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 1, nil
		}, 10)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Error("goroutine appears to have leaked")
	}
}

func TestWithTimeout_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result, err := WithTimeout(func() (int, error) {
		<-ctx.Done()
		return 0, ctx.Err()
	}, 100)

	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
	if result != 0 {
		t.Errorf("expected zero value, got %v", result)
	}
}