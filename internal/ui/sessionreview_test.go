package ui

import (
	"testing"
	"time"

	"github.com/freecode/freecode/internal/session"
)

func TestSessionReviewDialogOpenClose(t *testing.T) {
	d := NewSessionReviewDialog()
	if d.IsOpen() {
		t.Error("dialog should start closed")
	}
	d.Open()
	if !d.IsOpen() {
		t.Error("dialog should be open after Open()")
	}
	d.Close()
	if d.IsOpen() {
		t.Error("dialog should be closed after Close()")
	}
}

func TestSessionReviewDialogNavigation(t *testing.T) {
	d := NewSessionReviewDialog()
	d.SetSessions([]*session.Session{
		{ID: "1", Title: "Session 1"},
		{ID: "2", Title: "Session 2"},
		{ID: "3", Title: "Session 3"},
	})

	if d.selectedIdx != 0 {
		t.Errorf("initial selection = %d, want 0", d.selectedIdx)
	}

	d.Next()
	if d.selectedIdx != 1 {
		t.Errorf("after Next, selection = %d, want 1", d.selectedIdx)
	}

	d.Prev()
	if d.selectedIdx != 0 {
		t.Errorf("after Prev, selection = %d, want 0", d.selectedIdx)
	}

	d.Next()
	d.Next()
	if d.selectedIdx != 2 {
		t.Errorf("after 2x Next, selection = %d, want 2", d.selectedIdx)
	}
	d.Next()
	if d.selectedIdx != 2 {
		t.Error("selection should not exceed list bounds")
	}
}

func TestSessionReviewDialogBoundaryNavigation(t *testing.T) {
	d := NewSessionReviewDialog()
	d.SetSessions([]*session.Session{
		{ID: "1", Title: "Only Session"},
	})

	d.Prev()
	if d.selectedIdx != 0 {
		t.Error("selection should not go below 0")
	}
}

func TestSessionReviewDialogKeyHandling(t *testing.T) {
	d := NewSessionReviewDialog()
	d.SetSessions([]*session.Session{
		{ID: "1", Title: "Session 1"},
		{ID: "2", Title: "Session 2"},
	})
	d.Open()

	if d.HandleKey("down") != true {
		t.Error("down key should be handled")
	}
	if d.selectedIdx != 1 {
		t.Errorf("selection = %d, want 1", d.selectedIdx)
	}

	if d.HandleKey("up") != true {
		t.Error("up key should be handled")
	}
	if d.selectedIdx != 0 {
		t.Errorf("selection = %d, want 0", d.selectedIdx)
	}

	if d.HandleKey("escape") != true {
		t.Error("escape key should close dialog")
	}
	if d.IsOpen() {
		t.Error("dialog should be closed after escape")
	}
}

func TestSessionReviewDialogGetSelectedSession(t *testing.T) {
	sessions := []*session.Session{
		{ID: "1", Title: "First"},
		{ID: "2", Title: "Second"},
	}
	d := NewSessionReviewDialog()
	d.SetSessions(sessions)

	sel := d.GetSelectedSession()
	if sel.ID != "1" {
		t.Errorf("selected = %s, want 1", sel.ID)
	}

	d.Next()
	sel = d.GetSelectedSession()
	if sel.ID != "2" {
		t.Errorf("selected = %s, want 2", sel.ID)
	}
}

func TestSessionReviewDialogRender(t *testing.T) {
	d := NewSessionReviewDialog()
	d.SetSessions([]*session.Session{
		{
			ID:        "1",
			Title:     "Test Session",
			Model:     "claude-3",
			Agent:     "sisyphus",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Messages:  []session.Message{{Role: "user", Content: "hello"}},
		},
	})
	d.Open()

	output := d.Render()
	if output == "" {
		t.Error("Render() should return non-empty string")
	}
	if !contains(output, "Session History") {
		t.Error("Render() should contain header")
	}
	if !contains(output, "Test Session") {
		t.Error("Render() should contain session title")
	}
}

func TestSessionReviewDialogRenderClosed(t *testing.T) {
	d := NewSessionReviewDialog()
	output := d.Render()
	if output != "" {
		t.Errorf("Render() on closed dialog should return empty string, got %q", output)
	}
}

func TestSessionReviewDialogRenderEmpty(t *testing.T) {
	d := NewSessionReviewDialog()
	d.Open()
	d.SetSessions([]*session.Session{})

	output := d.Render()
	if !contains(output, "No sessions") {
		t.Error("Render() should show 'No sessions' when list is empty")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
