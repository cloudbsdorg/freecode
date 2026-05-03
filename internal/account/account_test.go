package account

import (
	"context"
	"testing"
)

func TestMemoryRepo(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()

	info := &Info{
		ID:    "acc-1",
		Email: "test@example.com",
		URL:   "https://example.com",
	}

	err := repo.PersistAccount(ctx, info, "access", "refresh")
	if err != nil {
		t.Errorf("PersistAccount error: %v", err)
	}

	active, err := repo.Active(ctx)
	if err != nil {
		t.Errorf("Active error: %v", err)
	}
	if active == nil {
		t.Error("expected active account")
	}
	if active.Email != "test@example.com" {
		t.Errorf("expected test@example.com, got %s", active.Email)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Errorf("List error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 account, got %d", len(list))
	}

	err = repo.Remove(ctx, "acc-1")
	if err != nil {
		t.Errorf("Remove error: %v", err)
	}

	active, err = repo.Active(ctx)
	if active != nil {
		t.Error("expected nil active after remove")
	}
}
