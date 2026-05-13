package style

import (
	"github.com/charmbracelet/lipgloss"
)

type LipglossStyle struct {
	s lipgloss.Style
}

func NewLipglossStyle() LipglossStyle {
	return LipglossStyle{s: lipgloss.NewStyle()}
}

func (l LipglossStyle) Render(text string) string {
	return l.s.Render(text)
}

func (l LipglossStyle) Foreground(color lipgloss.Color) Style {
	l.s = l.s.Foreground(color)
	return l
}

func (l LipglossStyle) Background(color lipgloss.Color) Style {
	l.s = l.s.Background(color)
	return l
}

func (l LipglossStyle) Bold(_ ...bool) Style {
	l.s = l.s.Bold(true)
	return l
}

func (l LipglossStyle) Italic(_ ...bool) Style {
	l.s = l.s.Italic(true)
	return l
}

func (l LipglossStyle) Reverse(_ ...bool) Style {
	l.s = l.s.Reverse(true)
	return l
}

func (l LipglossStyle) Underline(_ ...bool) Style {
	l.s = l.s.Underline(true)
	return l
}

func (l LipglossStyle) Width(w int) Style {
	l.s = l.s.Width(w)
	return l
}

func (l LipglossStyle) Height(h int) Style {
	l.s = l.s.Height(h)
	return l
}

func (l LipglossStyle) Padding(values ...int) Style {
	switch len(values) {
	case 1:
		l.s = l.s.Padding(values[0])
	case 2:
		l.s = l.s.Padding(values[0], values[1])
	case 3:
		l.s = l.s.Padding(values[0], values[1], values[2])
	case 4:
		l.s = l.s.Padding(values[0], values[1], values[2], values[3])
	}
	return l
}

func (l LipglossStyle) Margin(values ...int) Style {
	switch len(values) {
	case 1:
		l.s = l.s.Margin(values[0])
	case 2:
		l.s = l.s.Margin(values[0], values[1])
	case 3:
		l.s = l.s.Margin(values[0], values[1], values[2])
	case 4:
		l.s = l.s.Margin(values[0], values[1], values[2], values[3])
	}
	return l
}

func (l LipglossStyle) MarginTop(v int) Style {
	l.s = l.s.MarginTop(v)
	return l
}

func (l LipglossStyle) MarginLeft(v int) Style {
	l.s = l.s.MarginLeft(v)
	return l
}

func (l LipglossStyle) BorderStyle(b lipgloss.Border) Style {
	l.s = l.s.BorderStyle(b)
	return l
}

func (l LipglossStyle) BorderForeground(color lipgloss.Color) Style {
	l.s = l.s.BorderForeground(color)
	return l
}

func Color(s string) lipgloss.Color {
	return lipgloss.Color(s)
}

func NormalBorder() lipgloss.Border {
	return lipgloss.NormalBorder()
}

func RoundedBorder() lipgloss.Border {
	return lipgloss.RoundedBorder()
}

func HiddenBorder() lipgloss.Border {
	return lipgloss.HiddenBorder()
}

func Normal() lipgloss.Border {
	return lipgloss.NormalBorder()
}

func Rounded() lipgloss.Border {
	return lipgloss.RoundedBorder()
}

func Hidden() lipgloss.Border {
	return lipgloss.HiddenBorder()
}
