package component

import "github.com/freecode/freecode/internal/renderer"

type Window[R renderer.Renderer] struct {
	Component[R]
	Title  string
	Colors WindowColors
}

type WindowColors struct {
	Background   string
	BorderColor  string
	TitleColor   string
}

func NewWindow[R renderer.Renderer](width, height int, title string, colors WindowColors) *Window[R] {
	return &Window[R]{
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

func (w *Window[R]) CenterIn(containerW, containerH int) {
	w.X = (containerW - w.Width) / 2
	if w.X < 0 {
		w.X = 0
	}
	w.Y = (containerH - w.Height) / 2
	if w.Y < 0 {
		w.Y = 0
	}
}

func (w *Window[R]) AlignLeft(containerH int) {
	w.Y = (containerH - w.Height) / 2
	if w.Y < 0 {
		w.Y = 0
	}
}

func (w *Window[R]) AlignRight(containerW, containerH int) {
	w.X = containerW - w.Width
	w.Y = (containerH - w.Height) / 2
	if w.Y < 0 {
		w.Y = 0
	}
}

func (w *Window[R]) AlignTop(containerW int) {
	w.X = (containerW - w.Width) / 2
	if w.X < 0 {
		w.X = 0
	}
	w.Y = 0
}

func (w *Window[R]) AlignBottom(containerW, containerH int) {
	w.X = (containerW - w.Width) / 2
	if w.X < 0 {
		w.X = 0
	}
	w.Y = containerH - w.Height
}

func (w *Window[R]) RenderContent(content string, r R) string {
	if !w.Visible {
		return ""
	}
	border := r.RenderBorder(w.X, w.Y, w.Width, w.Height, w.Colors.BorderColor)
	text := r.RenderText(w.Title, w.X+1, w.Y, w.Colors.TitleColor)
	text += r.RenderText(content, w.X+1, w.Y+1, w.Colors.Background)
	return border + text
}
