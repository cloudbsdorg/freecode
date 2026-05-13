package style

import "github.com/charmbracelet/lipgloss"

type Style interface {
	Render(text string) string
	Foreground(color lipgloss.Color) Style
	Background(color lipgloss.Color) Style
	Bold(v ...bool) Style
	Italic(v ...bool) Style
	Reverse(v ...bool) Style
	Underline(v ...bool) Style
	Width(w int) Style
	Height(h int) Style
	Padding(values ...int) Style
	Margin(values ...int) Style
	MarginTop(v int) Style
	MarginLeft(v int) Style
	BorderStyle(b lipgloss.Border) Style
	BorderForeground(color lipgloss.Color) Style
}
