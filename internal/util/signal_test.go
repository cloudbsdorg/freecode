package util

import (
	"sync"
	"testing"
	"time"
)

func TestNewSignal_Basic(t *testing.T) {
	s := NewSignal()
	ch := s.Wait()

	select {
	case <-ch:
		t.Fatal("channel should not be closed before trigger")
	default:
	}

	s.Trigger()

	select {
	case <-ch:
	default:
		t.Fatal("channel should be closed after trigger")
	}
}

func TestNewSignal_OneShot(t *testing.T) {
	s := NewSignal()

	s.Trigger()
	s.Trigger()
	s.Trigger()

	ch := s.Wait()
	select {
	case <-ch:
	default:
		t.Fatal("channel should be closed after trigger")
	}
}

func TestNewSignal_MultipleWaiters(t *testing.T) {
	s := NewSignal()

	const n = 5
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			<-s.Wait()
		}()
	}

	time.Sleep(10 * time.Millisecond)

	s.Trigger()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for waiters")
	}
}

func TestNewSignal_WaitAfterTrigger(t *testing.T) {
	s := NewSignal()

	s.Trigger()

	ch := s.Wait()

	select {
	case <-ch:
	default:
		t.Fatal("channel should be immediately closed when wait is called after trigger")
	}
}

func TestNewSignal_TriggerFromDifferentGoroutine(t *testing.T) {
	s := NewSignal()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.Trigger()
	}()

	go func() {
		defer wg.Done()
		<-s.Wait()
	}()

	if !waitWithTimeout(&wg, time.Second) {
		t.Fatal("timed out")
	}
}

func TestNewSignal_ConcurrentTriggers(t *testing.T) {
	s := NewSignal()

	const n = 10
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.Trigger()
		}()
	}

	if !waitWithTimeout(&wg, time.Second) {
		t.Fatal("timed out")
	}

	ch := s.Wait()
	select {
	case <-ch:
	default:
		t.Fatal("channel should be closed")
	}
}

func waitWithTimeout(wg *sync.WaitGroup, d time.Duration) bool {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(d):
		return false
	}
}
