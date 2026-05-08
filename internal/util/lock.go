package util

import (
	"sync"
)

type Lock struct {
	mu             sync.Mutex
	readers        map[string]int
	waitingWriters map[string]int
	waitingReaders map[string]int
	writers        map[string]bool
	keyConds       map[string]*sync.Cond
}

type RUnlock func()
type WUnlock func()

func NewLock() *Lock {
	return &Lock{
		readers:        make(map[string]int),
		waitingWriters: make(map[string]int),
		waitingReaders: make(map[string]int),
		writers:        make(map[string]bool),
		keyConds:       make(map[string]*sync.Cond),
	}
}

func (l *Lock) Read(key string) (RUnlock, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.writers[key] && l.waitingWriters[key] == 0 {
		l.readers[key]++
		return func() {
			l.mu.Lock()
			defer l.mu.Unlock()
			l.readers[key]--
			if l.readers[key] == 0 && !l.writers[key] {
				l.getKeyCond(key).Signal()
			}
		}, nil
	}

	l.waitingReaders[key]++
	for l.writers[key] || l.waitingWriters[key] > 0 {
		l.getKeyCond(key).Wait()
	}
	l.waitingReaders[key]--
	l.readers[key]++

	return func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.readers[key]--
		if l.readers[key] == 0 && !l.writers[key] {
			l.getKeyCond(key).Signal()
		}
	}, nil
}

func (l *Lock) Write(key string) (WUnlock, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.writers[key] && l.readers[key] == 0 {
		l.writers[key] = true
		return func() {
			l.mu.Lock()
			defer l.mu.Unlock()
			l.writers[key] = false
			l.getKeyCond(key).Broadcast()
		}, nil
	}

	l.waitingWriters[key]++
	for l.writers[key] || l.readers[key] > 0 {
		l.getKeyCond(key).Wait()
	}
	l.waitingWriters[key]--
	l.writers[key] = true

	return func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.writers[key] = false
		l.getKeyCond(key).Broadcast()
	}, nil
}

func (l *Lock) getKeyCond(key string) *sync.Cond {
	if _, ok := l.keyConds[key]; !ok {
		l.keyConds[key] = sync.NewCond(&l.mu)
	}
	return l.keyConds[key]
}