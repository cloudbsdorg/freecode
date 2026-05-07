package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	s := NewStore("/tmp/test")

	if s.dir != "/tmp/test" {
		t.Errorf("dir = %q, want %q", s.dir, "/tmp/test")
	}
}

func TestStoreSaveAndLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)
	sess := &Session{
		ID:    "test-session",
		Title: "Test Session",
	}

	err = s.SaveSession(sess)
	if err != nil {
		t.Fatalf("SaveSession() error = %v", err)
	}

	loaded, err := s.LoadSession("test-session")
	if err != nil {
		t.Fatalf("LoadSession() error = %v", err)
	}

	if loaded.ID != sess.ID {
		t.Errorf("loaded.ID = %q, want %q", loaded.ID, sess.ID)
	}

	if loaded.Title != sess.Title {
		t.Errorf("loaded.Title = %q, want %q", loaded.Title, sess.Title)
	}
}

func TestStoreLoadNonexistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)

	_, err = s.LoadSession("nonexistent")
	if err == nil {
		t.Error("LoadSession() should error for nonexistent session")
	}
}

func TestStoreDeleteSession(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)
	sess := &Session{ID: "delete-me"}

	if err := s.SaveSession(sess); err != nil {
		t.Fatalf("SaveSession() error = %v", err)
	}

	err = s.DeleteSession("delete-me")
	if err != nil {
		t.Fatalf("DeleteSession() error = %v", err)
	}

	_, err = s.LoadSession("delete-me")
	if err == nil {
		t.Error("LoadSession() should error after DeleteSession()")
	}
}

func TestStoreDeleteNonexistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)

	err = s.DeleteSession("nonexistent")
	if err == nil {
		t.Error("DeleteSession() should error for nonexistent session")
	}
}

func TestStoreListSessions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)

	s.SaveSession(&Session{ID: "sess1", Title: "Session 1"})
	s.SaveSession(&Session{ID: "sess2", Title: "Session 2"})
	s.SaveSession(&Session{ID: "sess3", Title: "Session 3"})

	sessions, err := s.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("len(sessions) = %d, want 3", len(sessions))
	}
}

func TestStoreListSessionsEmpty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)

	sessions, err := s.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("len(sessions) = %d, want 0", len(sessions))
	}
}

func TestStoreListSessionsNonexistentDir(t *testing.T) {
	s := NewStore("/nonexistent/path")

	sessions, err := s.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}

	if sessions != nil {
		t.Errorf("ListSessions() = %v, want nil", sessions)
	}
}

func TestStoreSaveSessionCreatesDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	nestedDir := filepath.Join(tmpDir, "nested", "path")
	s := NewStore(nestedDir)

	err = s.SaveSession(&Session{ID: "test"})
	if err != nil {
		t.Fatalf("SaveSession() error = %v", err)
	}

	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("SaveSession() did not create nested directory")
	}
}

func TestStoreSaveSessionInvalidJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)

	sess := &Session{
		ID:       "test",
		Metadata: map[string]interface{}{"unserializable": make(chan int)},
	}

	err = s.SaveSession(sess)
	if err == nil {
		t.Error("SaveSession() should error for unserializable data")
	}
}

func TestStoreSaveSessionWriteFileFails(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := NewStore(tmpDir)

	preExistingFile := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(preExistingFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Chmod(preExistingFile, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	err = s.SaveSession(&Session{ID: "test"})
	if err == nil {
		t.Error("SaveSession() should error when WriteFile fails on read-only file")
	}
}

func TestStoreSaveSessionMkdirAllFails(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping as root user")
	}

	s := NewStore("/proc/fake-path-that-cannot-be-created")

	err := s.SaveSession(&Session{ID: "test"})
	if err == nil {
		t.Error("SaveSession() should error when MkdirAll fails")
	}
}

func TestStoreLoadSessionCorruptJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	corruptFile := filepath.Join(tmpDir, "corrupt.json")
	if err := os.WriteFile(corruptFile, []byte("not valid json{{{"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	s := NewStore(tmpDir)

	_, err = s.LoadSession("corrupt")
	if err == nil {
		t.Error("LoadSession() should error for corrupt JSON")
	}
}

func TestStoreListSessionsWithSubdirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(subdir, "nested.json"), []byte(`{"ID":"nested"}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	s := NewStore(tmpDir)

	sessions, err := s.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("len(sessions) = %d, want 0 (subdirectory files ignored)", len(sessions))
	}
}

func TestStoreListSessionsReadDirError(t *testing.T) {
	f, err := os.CreateTemp("", "not-a-dir")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	s := NewStore(f.Name())

	_, err = s.ListSessions()
	if err == nil {
		t.Error("ListSessions() should error when path is not a directory")
	}
}

func TestStoreListSessionsWithLoadError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := os.WriteFile(filepath.Join(tmpDir, "load-error.json"), []byte("invalid json"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	s := NewStore(tmpDir)

	sessions, err := s.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("len(sessions) = %d, want 0 (load errors skipped)", len(sessions))
	}
}
