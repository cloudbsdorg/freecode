package style

type GlyphStyle struct{}

func NewGlyphStyle() GlyphStyle {
	return GlyphStyle{}
}

func (g GlyphStyle) Render(text string) string {
	return text
}

func (g GlyphStyle) Foreground(color string) Style {
	return g
}

func (g GlyphStyle) Background(color string) Style {
	return g
}

func (g GlyphStyle) Bold(v ...bool) Style {
	return g
}

func (g GlyphStyle) Italic(v ...bool) Style {
	return g
}

func (g GlyphStyle) Width(w int) Style {
	return g
}

func (g GlyphStyle) Height(h int) Style {
	return g
}

func (g GlyphStyle) Padding(values ...int) Style {
	return g
}

func (g GlyphStyle) Margin(values ...int) Style {
	return g
}

func (g GlyphStyle) MarginTop(v int) Style {
	return g
}

func (g GlyphStyle) MarginLeft(v int) Style {
	return g
}

func (g GlyphStyle) BorderStyle(b Border) Style {
	return g
}

func (g GlyphStyle) BorderForeground(color string) Style {
	return g
}
