package pty

import (
	"testing"
)

func TestTerminal(t *testing.T) {
	term := NewTerminal(80, 24)
	cols, rows := term.Size()
	if cols != 80 || rows != 24 {
		t.Errorf("expected 80x24, got %dx%d", cols, rows)
	}

	term.SetSize(100, 40)
	cols, rows = term.Size()
	if cols != 100 || rows != 40 {
		t.Errorf("expected 100x40, got %dx%d", cols, rows)
	}
}
