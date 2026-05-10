package component

import "github.com/freecode/freecode/internal/renderer"

type Dialog[R renderer.Renderer] struct {
	Component[R]
	Title     string
	Content   string
	Colors    DialogColors
}

type DialogColors struct {
	Background   string
	Foreground   string
	BorderColor string
	TitleColor  string
}

func NewDialog[R renderer.Renderer](width, height int, title string, colors DialogColors) *Dialog[R] {
	return &Dialog[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  height,
			Visible: true,
		},
		Title:  title,
		Colors: colors,
	}
}

func (d *Dialog[R]) SetTitle(title string) {
	d.Title = title
}

func (d *Dialog[R]) SetContent(content string) {
	d.Content = content
}

func (d *Dialog[R]) Render(r R) string {
	if !d.Visible {
		return ""
	}

	lines := d.Title + "\n\n" + d.Content
	return r.RenderBorder(d.X, d.Y, d.Width, d.Height, d.Colors.BorderColor) +
		r.RenderText(lines, d.X+1, d.Y+1, d.Colors.Foreground)
}
