package util

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewLock(t *testing.T) {
	l := NewLock()
	if l == nil {
		t.Fatal("NewLock returned nil")
	}
	if l.readers == nil {
		t.Error("readers map is nil")
	}
	if l.writers == nil {
		t.Error("writers map is nil")
	}
}

func TestReadLock(t *testing.T) {
	l := NewLock()
	unlock, err := l.Read("test-key")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if unlock == nil {
		t.Fatal("Read returned nil unlock")
	}
	unlock()
}

func TestWriteLock(t *testing.T) {
	l := NewLock()
	unlock, err := l.Write("test-key")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if unlock == nil {
		t.Fatal("Write returned nil unlock")
	}
	unlock()
}

func TestMultipleReadersSameKey(t *testing.T) {
	l := NewLock()
	var wg sync.WaitGroup
	concurrentReaders := 0
	var mu sync.Mutex
	maxConcurrent := 0

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unlock, _ := l.Read("shared-key")
			mu.Lock()
			concurrentReaders++
			if concurrentReaders > maxConcurrent {
				maxConcurrent = concurrentReaders
			}
			mu.Unlock()
			time.Sleep(20 * time.Millisecond)
			mu.Lock()
			concurrentReaders--
			mu.Unlock()
			unlock()
		}()
	}
	wg.Wait()

	if maxConcurrent < 2 {
		t.Errorf("expected multiple concurrent readers, got max %d", maxConcurrent)
	}
}

func TestWriteExcludesReaders(t *testing.T) {
	l := NewLock()
	writerActive := false
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		unlock, _ := l.Write("test-key")
		writerActive = true
		time.Sleep(50 * time.Millisecond)
		writerActive = false
		unlock()
	}()

	time.Sleep(10 * time.Millisecond)

	wg.Add(1)
	var readerSawWriter int32
	go func() {
		defer wg.Done()
		unlock, _ := l.Read("test-key")
		if writerActive {
			atomic.StoreInt32(&readerSawWriter, 1)
		}
		unlock()
	}()

	wg.Wait()

	if atomic.LoadInt32(&readerSawWriter) == 1 {
		t.Error("reader acquired lock while writer was active")
	}
}

func TestWriteExcludesWriters(t *testing.T) {
	l := NewLock()
	var wg sync.WaitGroup
	firstWriterDone := false

	wg.Add(2)
	go func() {
		defer wg.Done()
		unlock, _ := l.Write("test-key")
		time.Sleep(50 * time.Millisecond)
		firstWriterDone = true
		unlock()
	}()

	time.Sleep(10 * time.Millisecond)

	go func() {
		defer wg.Done()
		unlock, _ := l.Write("test-key")
		if !firstWriterDone {
			t.Error("second writer acquired lock while first writer was active")
		}
		unlock()
	}()

	wg.Wait()
}

func TestDifferentKeysIndependent(t *testing.T) {
	l := NewLock()
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		unlock, _ := l.Write("key-a")
		time.Sleep(50 * time.Millisecond)
		unlock()
	}()
	go func() {
		defer wg.Done()
		unlock, _ := l.Read("key-b")
		time.Sleep(10 * time.Millisecond)
		unlock()
	}()

	wg.Wait()
}

func TestReaderUnlockSignal(t *testing.T) {
	l := NewLock()
	unlock1, _ := l.Read("key")
	unlock2, _ := l.Read("key")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		unlock, _ := l.Write("key")
		unlock()
	}()

	time.Sleep(10 * time.Millisecond)
	unlock1()
	unlock2()

	wg.Wait()
}

func TestWriterPriorityBlocksReaders(t *testing.T) {
	l := NewLock()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		unlock, _ := l.Write("priority-key")
		time.Sleep(50 * time.Millisecond)
		unlock()
	}()

	time.Sleep(10 * time.Millisecond)

	var readerGotLock int32
	wg.Add(1)
	go func() {
		defer wg.Done()
		unlock, _ := l.Read("priority-key")
		atomic.StoreInt32(&readerGotLock, 1)
		unlock()
	}()

	time.Sleep(20 * time.Millisecond)

	if atomic.LoadInt32(&readerGotLock) == 1 {
		t.Error("reader acquired lock while writer was waiting (writer priority violated)")
	}

	wg.Wait()
}

func TestReadAfterWriteComplete(t *testing.T) {
	l := NewLock()
	var wg sync.WaitGroup

	unlock, _ := l.Write("key")
	unlock()

	wg.Add(1)
	go func() {
		defer wg.Done()
		unlock, err := l.Read("key")
		if err != nil {
			t.Errorf("Read failed after write released: %v", err)
		}
		unlock()
	}()

	wg.Wait()
}