package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	style := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Background(lipgloss.Color(bgColor)).
		MarginTop(y).
		MarginLeft(x)
	return style.Render(content)
}

func (b *BubbleRenderer) RenderText(text string, x, y int, fgColor string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(fgColor)).
		MarginTop(y).
		MarginLeft(x)
	return style.Render(text)
}

func (b *BubbleRenderer) RenderBorder(x, y, w, h int, fgColor string) string {
	if w < 2 || h < 2 {
		return ""
	}
	style := lipgloss.NewStyle().
		Width(w).
		Height(h).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(fgColor)).
		MarginTop(y).
		MarginLeft(x)
	return style.Render("")
}

func (b *BubbleRenderer) RenderSelected(text string, x, y, w int, fg, bg string) string {
	if len(text) > w {
		text = text[:w-3] + "..."
	}
	style := lipgloss.NewStyle().
		Width(w).
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg)).
		MarginTop(y).
		MarginLeft(x)
	return style.Render(text)
}

func (b *BubbleRenderer) Width() int  { return b.width }
func (b *BubbleRenderer) Height() int { return b.height }
