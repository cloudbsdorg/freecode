package effect

import (
	"context"
	"testing"
)

func TestRegistry(t *testing.T) {
	Register("test", func() Effect {
		return &testEffect{}
	})

	factory, ok := Get("test")
	if !ok {
		t.Error("expected to get test effect")
	}
	if factory() == nil {
		t.Error("expected effect")
	}
}

type testEffect struct{}

func (e *testEffect) Run(ctx context.Context) error {
	return nil
}
