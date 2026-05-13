package renderer

import (
	"strings"

	"github.com/freecode/freecode/internal/style"
)

type BubbleRenderer struct {
	width  int
	height int
}

func NewBubbleRenderer(w, h int) *BubbleRenderer {
	return &BubbleRenderer{width: w, height: h}
}

func (b *BubbleRenderer) RenderBox(x, y, w, h int, bgColor string) string {
	lines := make([]string, h)
	for i := range lines {
		lines[i] = strings.Repeat(" ", w)
	}
	content := strings.Join(lines, "\n")
	s := style.NewStyle().
		Width(w).
		Height(h).
		Background(style.Color(bgColor)).
		MarginTop(y).
		MarginLeft(x)
	return s.Render(content)
}

func (b *BubbleRenderer) RenderText(text string, x, y int, fgColor string) string {
	s := style.NewStyle().
		Foreground(style.Color(fgColor)).
		MarginTop(y).
		MarginLeft(x)
	return s.Render(text)
}

func (b *BubbleRenderer) RenderBorder(x, y, w, h int, fgColor string) string {
	if w < 2 || h < 2 {
		return ""
	}
	s := style.NewStyle().
		Width(w).
		Height(h).
		BorderStyle(style.Rounded()).
		BorderForeground(style.Color(fgColor)).
		MarginTop(y).
		MarginLeft(x)
	return s.Render("")
}

func (b *BubbleRenderer) RenderSelected(text string, x, y, w int, fg, bg string) string {
	if len(text) > w {
		text = text[:w-3] + "..."
	}
	s := style.NewStyle().
		Width(w).
		Foreground(style.Color(fg)).
		Background(style.Color(bg)).
		MarginTop(y).
		MarginLeft(x)
	return s.Render(text)
}

func (b *BubbleRenderer) Width() int  { return b.width }
func (b *BubbleRenderer) Height() int { return b.height }
