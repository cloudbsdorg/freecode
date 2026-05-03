package command

import (
	"testing"
)

func TestRegistry(t *testing.T) {
	var called bool
	Register("test", func(args []string) error {
		called = true
		return nil
	})

	fn, ok := Get("test")
	if !ok {
		t.Error("expected to get test command")
	}

	if err := fn(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !called {
		t.Error("command was not called")
	}
}
