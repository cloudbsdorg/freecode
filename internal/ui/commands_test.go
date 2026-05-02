package ui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
)

func TestNewCommandHandler(t *testing.T) {
	h := NewCommandHandler()
	if h == nil {
		t.Fatal("NewCommandHandler() returned nil")
	}
	if h.commands == nil {
		t.Error("CommandHandler.commands is nil")
	}
}

func TestCommandHandlerRegister(t *testing.T) {
	h := NewCommandHandler()
	cmd := Command{
		Name:        "test",
		Description: "Test command",
		Shortcut:    "ctrl+t",
	}
	h.Register(cmd)

	if len(h.commands) != 1 {
		t.Errorf("len(commands) = %d, want 1", len(h.commands))
	}
}

func TestCommandHandlerList(t *testing.T) {
	h := NewCommandHandler()
	cmd := Command{Name: "test"}
	h.Register(cmd)

	list := h.List()
	if len(list) != 1 {
		t.Errorf("len(List()) = %d, want 1", len(list))
	}
}

func TestCommandHandlerSearch(t *testing.T) {
	h := NewCommandHandler()
	h.Register(Command{Name: "test", Description: "A test command"})
	h.Register(Command{Name: "foo", Description: "Something else"})

	results := h.Search("test")
	if len(results) != 1 {
		t.Errorf("len(Search('test')) = %d, want 1", len(results))
	}

	results = h.Search("something")
	if len(results) != 1 {
		t.Errorf("len(Search('something')) = %d, want 1", len(results))
	}

	results = h.Search("notfound")
	if len(results) != 0 {
		t.Errorf("len(Search('notfound')) = %d, want 0", len(results))
	}
}

func TestCommandHandlerSearchCaseInsensitive(t *testing.T) {
	h := NewCommandHandler()
	h.Register(Command{Name: "TEST", Description: "uppercase"})

	results := h.Search("test")
	if len(results) != 1 {
		t.Errorf("len(Search('test')) = %d, want 1", len(results))
	}
}

func TestCommandHandlerHandle(t *testing.T) {
	h := NewCommandHandler()
	m := &Model{}

	cmd := h.Handle(tea.KeyMsg{}, m)
	if cmd != nil {
		t.Error("Handle() should return nil for unknown key")
	}
}