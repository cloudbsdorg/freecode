package id

import (
	"testing"
)

func TestNew(t *testing.T) {
	id1 := New()
	id2 := New()
	if id1 == id2 {
		t.Error("expected unique IDs")
	}
	if len(id1) == 0 {
		t.Error("expected non-empty ID")
	}
}

func TestRandom(t *testing.T) {
	id1 := Random()
	id2 := Random()
	if id1 == id2 {
		t.Error("expected unique random IDs")
	}
	if len(id1) < 10 {
		t.Error("expected reasonable length ID")
	}
}

func TestShort(t *testing.T) {
	id := Short()
	if len(id) < 5 {
		t.Error("expected short ID")
	}
}
