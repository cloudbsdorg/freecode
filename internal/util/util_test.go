package util

import (
	"testing"
)

func TestRandomID(t *testing.T) {
	id1 := RandomID()
	id2 := RandomID()
	if id1 == id2 {
		t.Error("expected unique random IDs")
	}
	if len(id1) < 20 {
		t.Error("expected long ID")
	}
}

func TestTrimSpace(t *testing.T) {
	if TrimSpace("  hello  ") != "hello" {
		t.Error("expected trimmed string")
	}
}

func TestContains(t *testing.T) {
	if !Contains("hello world", "world") {
		t.Error("expected to contain")
	}
	if Contains("hello", "xyz") {
		t.Error("expected not to contain")
	}
}

func TestToLower(t *testing.T) {
	if ToLower("HELLO") != "hello" {
		t.Error("expected lowercase")
	}
}

func TestToUpper(t *testing.T) {
	if ToUpper("hello") != "HELLO" {
		t.Error("expected uppercase")
	}
}

func TestHash(t *testing.T) {
	h1 := Hash("test")
	h2 := Hash("test")
	if h1 != h2 {
		t.Error("expected same hash for same string")
	}
	h3 := Hash("other")
	if h1 == h3 {
		t.Error("expected different hash for different string")
	}
}
