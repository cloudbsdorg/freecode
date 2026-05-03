package file

import (
	"context"
	"os"
)

type Event struct {
	Path   string
	Type   EventType
	Exists bool
}

type EventType string

const (
	EventCreate EventType = "create"
	EventModify EventType = "modify"
	EventDelete EventType = "delete"
)

type Watcher interface {
	Watch(ctx context.Context, path string) error
	Unwatch(path string) error
	Events() <-chan Event
}

type watcher struct {
	events chan Event
}

func NewWatcher() Watcher {
	return &watcher{events: make(chan Event, 100)}
}

func (w *watcher) Watch(ctx context.Context, path string) error {
	go func() {
		stat, err := os.Stat(path)
		w.events <- Event{Path: path, Exists: err == nil}
		if stat != nil {
			_ = stat.IsDir()
		}
	}()
	return nil
}

func (w *watcher) Unwatch(path string) error {
	return nil
}

func (w *watcher) Events() <-chan Event {
	return w.events
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
