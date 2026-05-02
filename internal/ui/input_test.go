package ui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
)

func TestNewInputHandler(t *testing.T) {
	h := NewInputHandler()
	if h == nil {
		t.Fatal("NewInputHandler() returned nil")
	}
}

func TestNewInputState(t *testing.T) {
	s := NewInputState()
	if s == nil {
		t.Fatal("NewInputState() returned nil")
	}
	if s.Cursor != 0 {
		t.Errorf("Cursor = %d, want 0", s.Cursor)
	}
	if s.CursorMode != CursorModeBlock {
		t.Errorf("CursorMode = %v, want CursorModeBlock", s.CursorMode)
	}
}

func TestInputStateInsert(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 0}

	s.Insert('x')
	if s.Value != "xhello" {
		t.Errorf("Value = %q, want %q", s.Value, "xhello")
	}
	if s.Cursor != 1 {
		t.Errorf("Cursor = %d, want 1", s.Cursor)
	}
}

func TestInputStateInsertAtMiddle(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 3}

	s.Insert('x')
	if s.Value != "helxlo" {
		t.Errorf("Value = %q, want %q", s.Value, "helxlo")
	}
}

func TestInputStateDelete(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 4}

	s.Delete()
	if s.Value != "hell" {
		t.Errorf("Value = %q, want %q", s.Value, "hell")
	}
}

func TestInputStateDeleteAtStart(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 0}

	s.Delete()
	if s.Value != "ello" {
		t.Errorf("Value = %q, want %q", s.Value, "ello")
	}
}

func TestInputStateMoveLeft(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 3}

	s.MoveLeft()
	if s.Cursor != 2 {
		t.Errorf("Cursor = %d, want 2", s.Cursor)
	}
}

func TestInputStateMoveLeftAtStart(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 0}

	s.MoveLeft()
	if s.Cursor != 0 {
		t.Errorf("Cursor = %d, want 0", s.Cursor)
	}
}

func TestInputStateMoveRight(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 3}

	s.MoveRight()
	if s.Cursor != 4 {
		t.Errorf("Cursor = %d, want 4", s.Cursor)
	}
}

func TestInputStateMoveRightAtEnd(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 5}

	s.MoveRight()
	if s.Cursor != 5 {
		t.Errorf("Cursor = %d, want 5", s.Cursor)
	}
}

func TestInputStateMoveToStart(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 5}

	s.MoveToStart()
	if s.Cursor != 0 {
		t.Errorf("Cursor = %d, want 0", s.Cursor)
	}
}

func TestInputStateMoveToEnd(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 0}

	s.MoveToEnd()
	if s.Cursor != 5 {
		t.Errorf("Cursor = %d, want 5", s.Cursor)
	}
}

func TestInputStateDeleteWord(t *testing.T) {
	s := &InputState{Value: "hello world", Cursor: 11}

	s.DeleteWord()
	if s.Value != "hello " {
		t.Errorf("Value = %q, want %q", s.Value, "hello ")
	}
}

func TestInputStateDeleteWordAtStart(t *testing.T) {
	s := &InputState{Value: "hello", Cursor: 0}

	s.DeleteWord()
	if s.Value != "hello" {
		t.Errorf("Value = %q, want %q", s.Value, "hello")
	}
}

func TestInputHandlerHandle(t *testing.T) {
	h := NewInputHandler()
	state := &InputState{Value: "hello", Cursor: 5}

	h.Handle(tea.KeyMsg{Type: tea.KeyBackspace}, state)
	if state.Cursor != 4 {
		t.Errorf("Cursor = %d, want 4", state.Cursor)
	}
}

func TestInputHandlerHandleCtrlW(t *testing.T) {
	h := NewInputHandler()
	state := &InputState{Value: "hello world", Cursor: 11}

	h.Handle(tea.KeyMsg{Type: tea.KeyCtrlW}, state)
	if state.Value != "hello " {
		t.Errorf("Value = %q, want %q", state.Value, "hello ")
	}
}

func TestInputHandlerRender(t *testing.T) {
	h := NewInputHandler()
	state := &InputState{Value: "test", Cursor: 2, Focused: true}

	rendered := h.Render(state)
	if rendered == "" {
		t.Error("Render() returned empty string")
	}
}

func TestInputHandlerRenderPlaceholder(t *testing.T) {
	h := NewInputHandler()
	state := &InputState{Placeholder: "Enter text...", Focused: false}

	rendered := h.Render(state)
	if rendered == "" {
		t.Error("Render() returned empty string for placeholder")
	}
}