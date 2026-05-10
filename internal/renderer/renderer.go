package renderer

type Renderer interface {
	RenderBox(x, y, w, h int, bgColor string) string
	RenderText(text string, x, y int, fgColor string) string
	RenderBorder(x, y, w, h int, fgColor string) string
	RenderSelected(text string, x, y, w int, fg, bg string) string
	Width() int
	Height() int
}
