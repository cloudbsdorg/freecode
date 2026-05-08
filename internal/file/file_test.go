package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	w := NewWatcher()
	if w == nil {
		t.Fatal("NewWatcher() returned nil")
	}
}

func TestWatcherEvents(t *testing.T) {
	w := NewWatcher()
	events := w.Events()
	if events == nil {
		t.Error("Events() returned nil channel")
	}
}

func TestWatcherWatch(t *testing.T) {
	w := NewWatcher()
	ctx := context.Background()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if err := w.Watch(ctx, tmpFile); err != nil {
		t.Fatalf("Watch() error = %v", err)
	}

	select {
	case event := <-w.Events():
		if event.Path != tmpFile {
			t.Errorf("event.Path = %q, want %q", event.Path, tmpFile)
		}
		if !event.Exists {
			t.Error("event.Exists = false, want true")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Watch did not send event")
	}
}

func TestWatcherUnwatch(t *testing.T) {
	w := NewWatcher()
	ctx := context.Background()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	w.Watch(ctx, tmpFile)
	if err := w.Unwatch(tmpFile); err != nil {
		t.Errorf("Unwatch() error = %v", err)
	}
}

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("hello world")

	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	data, err := ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("ReadFile() = %q, want %q", string(data), string(content))
	}
}

func TestReadFileNotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/path/to/file")
	if err == nil {
		t.Error("ReadFile() for nonexistent file should return error")
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("hello world")

	if err := WriteFile(tmpFile, content); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("file content = %q, want %q", string(data), string(content))
	}
}

func TestMkdirAll(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "a", "b", "c")

	if err := MkdirAll(path, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to stat created directory: %v", err)
	}
	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	if Exists(tmpFile) {
		t.Error("Exists() = true for nonexistent file")
	}

	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if !Exists(tmpFile) {
		t.Error("Exists() = false for existing file")
	}
}

func TestIsDir(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	if IsDir(tmpFile) {
		t.Error("IsDir() = true for file")
	}

	if !IsDir(tmpDir) {
		t.Error("IsDir() = false for directory")
	}
}