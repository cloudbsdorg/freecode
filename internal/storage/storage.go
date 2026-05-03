package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage interface {
	Remove(key []string) error
	Read(key []string, dest any) error
	Update(key []string, fn func(draft any)) error
	Write(key []string, content any) error
	List(prefix []string) ([][]string, error)
}

type fileStorage struct {
	mu   sync.RWMutex
	dir  string
	locks map[string]*sync.Mutex
}

func New(dir string) (Storage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &fileStorage{
		dir:  dir,
		locks: make(map[string]*sync.Mutex),
	}, nil
}

func (s *fileStorage) lock(key []string) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()
	path := s.filePath(key)
	if s.locks[path] == nil {
		s.locks[path] = &sync.Mutex{}
	}
	return s.locks[path]
}

func (s *fileStorage) filePath(key []string) string {
	return filepath.Join(s.dir, filepath.Join(key...)) + ".json"
}

func (s *fileStorage) Remove(key []string) error {
	mu := s.lock(key)
	mu.Lock()
	defer mu.Unlock()

	path := s.filePath(key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *fileStorage) Read(key []string, dest any) error {
	mu := s.lock(key)
	mu.Lock()
	defer mu.Unlock()

	path := s.filePath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal(data, dest)
}

func (s *fileStorage) Update(key []string, fn func(draft any)) error {
	mu := s.lock(key)
	mu.Lock()
	defer mu.Unlock()

	path := s.filePath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}

	var draft any
	if err := json.Unmarshal(data, &draft); err != nil {
		return err
	}

	fn(draft)

	out, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}

func (s *fileStorage) Write(key []string, content any) error {
	mu := s.lock(key)
	mu.Lock()
	defer mu.Unlock()

	path := s.filePath(key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}

func (s *fileStorage) List(prefix []string) ([][]string, error) {
	dir := filepath.Join(s.dir, filepath.Join(prefix...))

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return [][]string{}, nil
		}
		return nil, err
	}

	var result [][]string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) > 5 && name[len(name)-5:] == ".json" {
			result = append(result, append(prefix, name[:len(name)-5]))
		}
	}
	return result, nil
}
