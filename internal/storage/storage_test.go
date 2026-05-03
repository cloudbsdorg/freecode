package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	type Data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err = s.Write([]string{"test"}, Data{Name: "Alice", Age: 30})
	if err != nil {
		t.Errorf("Write error: %v", err)
	}

	var result Data
	err = s.Read([]string{"test"}, &result)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	if result.Name != "Alice" || result.Age != 30 {
		t.Errorf("expected Alice/30, got %s/%d", result.Name, result.Age)
	}

	err = s.Update([]string{"test"}, func(draft any) {
		if m, ok := draft.(map[string]any); ok {
			if age, ok := m["age"].(float64); ok {
				m["age"] = age + 1
			}
		}
	})
	if err != nil {
		t.Errorf("Update error: %v", err)
	}

	var updated Data
	err = s.Read([]string{"test"}, &updated)
	if err != nil {
		t.Errorf("Read after update error: %v", err)
	}
	if updated.Age != 31 {
		t.Errorf("expected age 31, got %d", updated.Age)
	}

	err = s.Remove([]string{"test"})
	if err != nil {
		t.Errorf("Remove error: %v", err)
	}

	var missing Data
	err = s.Read([]string{"test"}, &missing)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	s.Write([]string{"a", "file1"}, map[string]any{"x": 1})
	s.Write([]string{"a", "file2"}, map[string]any{"x": 2})
	s.Write([]string{"b", "file3"}, map[string]any{"x": 3})

	entries, err := s.List([]string{"a"})
	if err != nil {
		t.Errorf("List error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	entries, err = s.List([]string{"b"})
	if err != nil {
		t.Errorf("List error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestNotFound(t *testing.T) {
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	err = s.Read([]string{"nonexistent"}, &result)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
