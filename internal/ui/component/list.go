package component

import "github.com/freecode/freecode/internal/renderer"

type List[R renderer.Renderer] struct {
	Component[R]
	Items     []string
	Colors    ListColors
}

type ListColors struct {
	Background   string
	Foreground   string
	BulletColor  string
	MutedColor   string
}

func NewList[R renderer.Renderer](width, height int, colors ListColors) *List[R] {
	return &List[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  height,
			Visible: true,
		},
		Items:  []string{},
		Colors: colors,
	}
}

func (l *List[R]) SetItems(items []string) {
	l.Items = items
}

func (l *List[R]) Add(item string) {
	l.Items = append(l.Items, item)
}

func (l *List[R]) Remove(idx int) {
	if idx < 0 || idx >= len(l.Items) {
		return
	}
	l.Items = append(l.Items[:idx], l.Items[idx+1:]...)
}

func (l *List[R]) Clear() {
	l.Items = []string{}
}

func (l *List[R]) Render(r R) string {
	if !l.Visible {
		return ""
	}

	lines := []string{}

	for i := 0; i < len(l.Items) && i < l.Height-2; i++ {
		bullet := "•"
		if l.Colors.BulletColor != "" {
			bullet = r.RenderText("•", l.X+1, l.Y+1+i, l.Colors.BulletColor)
		}
		item := truncate(l.Items[i], l.Width-len(bullet)-3)
		lines = append(lines, bullet+" "+item)
	}

	for i := len(lines); i < l.Height-2; i++ {
		lines = append(lines, "")
	}

	result := ""
	for _, line := range lines {
		result += r.RenderText(line, l.X+1, l.Y+1+len(result), l.Colors.Foreground) + "\n"
	}

	return r.RenderBox(l.X, l.Y, l.Width, l.Height, l.Colors.Background) + result
}
