package tab

import (
	"testing"
)

func TestNewSplitState(t *testing.T) {
	s := NewSplitState()
	if s == nil {
		t.Fatal("NewSplitState() returned nil")
	}
	if s.splits == nil {
		t.Error("splits is nil")
	}
}

func TestSplitStateCreate(t *testing.T) {
	s := NewSplitState()
	spl := s.Create("split-1", SplitVertical, 0.5)

	if spl.ID != "split-1" {
		t.Errorf("ID = %q, want %q", spl.ID, "split-1")
	}
	if spl.Direction != SplitVertical {
		t.Errorf("Direction = %v, want SplitVertical", spl.Direction)
	}
	if spl.Ratio != 0.5 {
		t.Errorf("Ratio = %f, want 0.5", spl.Ratio)
	}
}

func TestSplitStateGet(t *testing.T) {
	s := NewSplitState()
	s.Create("split-1", SplitVertical, 0.5)

	spl, ok := s.Get("split-1")
	if !ok {
		t.Error("Get() returned false")
	}
	if spl.ID != "split-1" {
		t.Errorf("ID = %q, want %q", spl.ID, "split-1")
	}
}

func TestSplitStateGetNotFound(t *testing.T) {
	s := NewSplitState()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for nonexistent")
	}
}

func TestSplitStateRemove(t *testing.T) {
	s := NewSplitState()
	s.Create("split-1", SplitVertical, 0.5)

	s.Remove("split-1")
	_, ok := s.Get("split-1")
	if ok {
		t.Error("Get() should return false after Remove()")
	}
}

func TestSplitStateSetRatio(t *testing.T) {
	s := NewSplitState()
	s.Create("split-1", SplitVertical, 0.5)

	result := s.SetRatio("split-1", 0.75)
	if !result {
		t.Error("SetRatio() returned false")
	}

	spl, _ := s.Get("split-1")
	if spl.Ratio != 0.75 {
		t.Errorf("Ratio = %f, want 0.75", spl.Ratio)
	}
}

func TestSplitStateSetRatioNotFound(t *testing.T) {
	s := NewSplitState()
	result := s.SetRatio("nonexistent", 0.5)
	if result {
		t.Error("SetRatio() should return false for nonexistent")
	}
}

func TestSplitStateToggle(t *testing.T) {
	s := NewSplitState()
	s.Create("split-1", SplitVertical, 0.5)

	s.Toggle("split-1")
	spl, _ := s.Get("split-1")
	if spl.Direction != SplitHorizontal {
		t.Errorf("Direction = %v, want SplitHorizontal", spl.Direction)
	}
}

func TestSplitStateToggleNotFound(t *testing.T) {
	s := NewSplitState()
	s.Toggle("nonexistent")
}

func TestSplitStateList(t *testing.T) {
	s := NewSplitState()
	s.Create("split-1", SplitVertical, 0.5)
	s.Create("split-2", SplitHorizontal, 0.5)

	splits := s.List()
	if len(splits) != 2 {
		t.Errorf("len(List()) = %d, want 2", len(splits))
	}
}

func TestSplitStateListEmpty(t *testing.T) {
	s := NewSplitState()
	splits := s.List()
	if len(splits) != 0 {
		t.Errorf("len(List()) = %d, want 0", len(splits))
	}
}