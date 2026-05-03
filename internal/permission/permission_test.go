package permission

import (
	"context"
	"testing"
)

func TestChecker(t *testing.T) {
	c := NewMemoryChecker()
	ctx := context.Background()

	perm := Permission{Resource: "file.txt", Actions: []Action{ActionRead}}
	granted, err := c.Check(ctx, "user1", perm)
	if err != nil {
		t.Errorf("Check error: %v", err)
	}
	if granted {
		t.Error("expected not granted initially")
	}

	if err := c.Grant(ctx, "user1", perm); err != nil {
		t.Errorf("Grant error: %v", err)
	}

	granted, _ = c.Check(ctx, "user1", perm)
	if !granted {
		t.Error("expected granted after Grant")
	}

	if err := c.Revoke(ctx, "user1", perm); err != nil {
		t.Errorf("Revoke error: %v", err)
	}

	granted, _ = c.Check(ctx, "user1", perm)
	if granted {
		t.Error("expected not granted after Revoke")
	}
}
