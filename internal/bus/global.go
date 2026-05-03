package bus

import "sync"

type GlobalEventHandler func(event GlobalEvent)

type globalHandlerEntry struct {
	id      int
	handler GlobalEventHandler
}

type GlobalBus struct {
	mu       sync.RWMutex
	handlers []globalHandlerEntry
	nextID   int
}

var (
	globalBus     *GlobalBus
	globalBusOnce sync.Once
)

func GetGlobalBus() *GlobalBus {
	globalBusOnce.Do(func() {
		globalBus = &GlobalBus{}
	})
	return globalBus
}

func (b *GlobalBus) Emit(event GlobalEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers {
		h.handler(event)
	}
}

func (b *GlobalBus) On(handler GlobalEventHandler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.nextID
	b.nextID++
	b.handlers = append(b.handlers, globalHandlerEntry{id: id, handler: handler})
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		for i, h := range b.handlers {
			if h.id == id {
				b.handlers = append(b.handlers[:i], b.handlers[i+1:]...)
				return
			}
		}
	}
}
