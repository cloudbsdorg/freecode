package git

import (
	"testing"
)

func TestRepository(t *testing.T) {
	r, err := Open("/tmp")
	if err != nil {
		t.Skip("git not available")
	}
	if r == nil {
		t.Error("expected repository")
	}
}
