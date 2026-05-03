package project

import (
	"testing"
)

func TestDetector(t *testing.T) {
	d := NewDetector()
	if d == nil {
		t.Error("expected detector")
	}
}
