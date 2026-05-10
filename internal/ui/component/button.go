package component

import "github.com/freecode/freecode/internal/renderer"

type Button[R renderer.Renderer] struct {
	Component[R]
	Text      string
	Selected  bool
	Pressed   bool
	OnClick   func()
	Colors    ButtonColors
}

type ButtonColors struct {
	Background    string
	Foreground    string
	SelectedBg    string
	SelectedFg    string
	PressedBg     string
	PressedFg     string
}

func NewButton[R renderer.Renderer](text string, x, y int, colors ButtonColors) *Button[R] {
	return &Button[R]{
		Component: Component[R]{
			X:       x,
			Y:       y,
			Width:   len(text) + 2,
			Height:  1,
			Visible: true,
		},
		Text:    text,
		Pressed: false,
		Colors:  colors,
	}
}

func (b *Button[R]) SetText(text string) {
	b.Text = text
	b.Width = len(text) + 2
}

func (b *Button[R]) Select() {
	b.Selected = true
}

func (b *Button[R]) Press() {
	b.Pressed = true
	if b.OnClick != nil {
		b.OnClick()
	}
}

func (b *Button[R]) Release() {
	b.Pressed = false
}

func (b *Button[R]) Render(r R) string {
	if !b.Visible {
		return ""
	}

	var fg, bg string
	if b.Pressed {
		fg = b.Colors.PressedFg
		bg = b.Colors.PressedBg
	} else if b.Selected {
		fg = b.Colors.SelectedFg
		bg = b.Colors.SelectedBg
	} else {
		fg = b.Colors.Foreground
		bg = b.Colors.Background
	}

	return r.RenderSelected(" "+b.Text+" ", b.X, b.Y, b.Width, fg, bg)
}
