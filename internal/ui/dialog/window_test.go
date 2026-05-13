package dialog

import (
	"testing"
)

func TestNewWindow(t *testing.T) {
	w := NewWindow()
	if w.X != 0 {
		t.Errorf("expected X=0, got %d", w.X)
	}
	if w.Y != 0 {
		t.Errorf("expected Y=0, got %d", w.Y)
	}
	if w.Width != 60 {
		t.Errorf("expected Width=60, got %d", w.Width)
	}
	if w.Height != 20 {
		t.Errorf("expected Height=20, got %d", w.Height)
	}
}

func TestWindowSetPosition(t *testing.T) {
	w := NewWindow()
	w.SetPosition(10, 20)
	if w.X != 10 {
		t.Errorf("expected X=10, got %d", w.X)
	}
	if w.Y != 20 {
		t.Errorf("expected Y=20, got %d", w.Y)
	}
}

func TestWindowSetSize(t *testing.T) {
	w := NewWindow()
	w.SetSize(80, 30)
	if w.Width != 80 {
		t.Errorf("expected Width=80, got %d", w.Width)
	}
	if w.Height != 30 {
		t.Errorf("expected Height=30, got %d", w.Height)
	}
}

func TestWindowCenterIn(t *testing.T) {
	tests := []struct {
		name       string
		containerW int
		containerH int
		w          *Window
		wantX      int
		wantY      int
	}{
		{
			name:       "centered fits",
			containerW: 100,
			containerH: 50,
			w:          &Window{Width: 60, Height: 20},
			wantX:      20,
			wantY:      15,
		},
		{
			name:       "centered overflow left",
			containerW: 40,
			containerH: 50,
			w:          &Window{Width: 60, Height: 20},
			wantX:      0,
			wantY:      15,
		},
		{
			name:       "centered overflow top",
			containerW: 100,
			containerH: 10,
			w:          &Window{Width: 60, Height: 20},
			wantX:      20,
			wantY:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.CenterIn(tt.containerW, tt.containerH)
			if tt.w.X != tt.wantX {
				t.Errorf("X: expected %d, got %d", tt.wantX, tt.w.X)
			}
			if tt.w.Y != tt.wantY {
				t.Errorf("Y: expected %d, got %d", tt.wantY, tt.w.Y)
			}
		})
	}
}

func TestWindowAlignLeft(t *testing.T) {
	w := &Window{Width: 60, Height: 20}
	w.AlignLeft(50)
	if w.X != 0 {
		t.Errorf("expected X=0, got %d", w.X)
	}
	if w.Y != 15 {
		t.Errorf("expected Y=15, got %d", w.Y)
	}

	// Overflow top
	w.AlignLeft(10)
	if w.Y != 0 {
		t.Errorf("expected Y=0, got %d", w.Y)
	}
}

func TestWindowAlignRight(t *testing.T) {
	w := &Window{Width: 60, Height: 20}
	w.AlignRight(100, 50)
	if w.X != 40 {
		t.Errorf("expected X=40, got %d", w.X)
	}
	if w.Y != 15 {
		t.Errorf("expected Y=15, got %d", w.Y)
	}

	// Overflow top
	w.AlignRight(100, 10)
	if w.Y != 0 {
		t.Errorf("expected Y=0, got %d", w.Y)
	}
}

func TestWindowAlignTop(t *testing.T) {
	w := &Window{Width: 60, Height: 20}
	w.AlignTop(100)
	if w.X != 20 {
		t.Errorf("expected X=20, got %d", w.X)
	}
	if w.Y != 0 {
		t.Errorf("expected Y=0, got %d", w.Y)
	}

	// Overflow left
	w.AlignTop(40)
	if w.X != 0 {
		t.Errorf("expected X=0, got %d", w.X)
	}
}

func TestWindowAlignBottom(t *testing.T) {
	w := &Window{Width: 60, Height: 20}
	w.AlignBottom(100, 50)
	if w.X != 20 {
		t.Errorf("expected X=20, got %d", w.X)
	}
	if w.Y != 30 {
		t.Errorf("expected Y=30, got %d", w.Y)
	}

	// Overflow left
	w.AlignBottom(40, 50)
	if w.X != 0 {
		t.Errorf("expected X=0, got %d", w.X)
	}
}

func TestWindowStyle(t *testing.T) {
	w := &Window{Width: 60, Height: 20}
	style := w.Style()
	_ = style
}

func TestWindowRender(t *testing.T) {
	w := &Window{Width: 10, Height: 3}
	result := w.Render("hello")
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestWindowRenderWithBackground(t *testing.T) {
	w := &Window{Width: 10, Height: 3}
	result := w.RenderWithBackground("hello", "#FF0000")
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestWindowRenderCentered(t *testing.T) {
	w := &Window{Width: 10, Height: 3}
	result := w.RenderCentered(100, 50, "hello")
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestWindowRenderCenteredWithBackground(t *testing.T) {
	w := &Window{Width: 10, Height: 3}
	result := w.RenderCenteredWithBackground(100, 50, "hello", "#FF0000")
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestRenderBox(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	colors := Dark
	result := RenderBox(lines, colors, 20, 10)
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestJoinLines(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected string
	}{
		{
			name:     "empty",
			lines:    []string{},
			expected: "",
		},
		{
			name:     "single",
			lines:    []string{"hello"},
			expected: "hello",
		},
		{
			name:     "multiple",
			lines:    []string{"hello", "world"},
			expected: "hello\nworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinLines(tt.lines)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRenderBoxCentered(t *testing.T) {
	lines := []string{"line1", "line2"}
	colors := Dark
	result := RenderBoxCentered(lines, colors, 20, 5, 100, 50)
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestRenderBoxCenteredOverflow(t *testing.T) {
	lines := []string{"line1"}
	colors := Dark
	// Container smaller than box
	result := RenderBoxCentered(lines, colors, 200, 100, 50, 30)
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestRenderBackdrop(t *testing.T) {
	result := RenderBackdrop(80, 24)
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestRenderWithBackdrop(t *testing.T) {
	backdrop := RenderBackdrop(80, 24)
	content := "dialog content"
	result := RenderWithBackdrop(backdrop, content)
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestRenderBackdropCentered(t *testing.T) {
	result := RenderBackdropCentered(80, 24, "content")
	if result == "" {
		t.Error("expected non-empty render result")
	}
}

func TestRenderDialogWithBackdrop(t *testing.T) {
	result := RenderDialogWithBackdrop("content", Dark, 80, 24)
	if result == "" {
		t.Error("expected non-empty render result")
	}
}