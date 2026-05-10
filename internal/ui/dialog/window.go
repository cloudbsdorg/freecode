package dialog

import "github.com/charmbracelet/lipgloss"

type Window struct {
	X      int
	Y      int
	Width  int
	Height int
}

func NewWindow() *Window {
	return &Window{
		X:      0,
		Y:      0,
		Width:  60,
		Height: 20,
	}
}

func (w *Window) SetPosition(x, y int) {
	w.X = x
	w.Y = y
}

func (w *Window) SetSize(width, height int) {
	w.Width = width
	w.Height = height
}

func (w *Window) CenterIn(containerW, containerH int) {
	w.X = (containerW - w.Width) / 2
	if w.X < 0 {
		w.X = 0
	}
	w.Y = (containerH - w.Height) / 2
	if w.Y < 0 {
		w.Y = 0
	}
}

func (w *Window) AlignLeft(containerH int) {
	w.Y = (containerH - w.Height) / 2
	if w.Y < 0 {
		w.Y = 0
	}
}

func (w *Window) AlignRight(containerW, containerH int) {
	w.X = containerW - w.Width
	w.Y = (containerH - w.Height) / 2
	if w.Y < 0 {
		w.Y = 0
	}
}

func (w *Window) AlignTop(containerW int) {
	w.X = (containerW - w.Width) / 2
	if w.X < 0 {
		w.X = 0
	}
	w.Y = 0
}

func (w *Window) AlignBottom(containerW, containerH int) {
	w.X = (containerW - w.Width) / 2
	if w.X < 0 {
		w.X = 0
	}
	w.Y = containerH - w.Height
}

func (w *Window) Style() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(w.Width).
		Height(w.Height)
}

func (w *Window) Render(content string) string {
	return w.Style().Render(content)
}

func (w *Window) RenderWithBackground(content string, bgColor string) string {
	return lipgloss.NewStyle().
		Width(w.Width).
		Height(w.Height).
		Background(lipgloss.Color(bgColor)).
		Render(content)
}

func (w *Window) RenderCentered(containerW, containerH int, content string) string {
	w.CenterIn(containerW, containerH)
	return lipgloss.NewStyle().
		Width(w.Width).
		MarginTop(w.Y).
		MarginLeft(w.X).
		Render(content)
}

func (w *Window) RenderCenteredWithBackground(containerW, containerH int, content, bgColor string) string {
	w.CenterIn(containerW, containerH)
	return lipgloss.NewStyle().
		Width(w.Width).
		Background(lipgloss.Color(bgColor)).
		MarginTop(w.Y).
		MarginLeft(w.X).
		Render(content)
}

func RenderBox(lines []string, colors Colors, width, height int) string {
	content := joinLines(lines)
	return lipgloss.NewStyle().
		Background(lipgloss.Color(colors.Background)).
		Width(width).
		Height(height).
		Render(content)
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func RenderBoxCentered(lines []string, colors Colors, width, height, containerW, containerH int) string {
	content := joinLines(lines)
	x := (containerW - width) / 2
	y := (containerH - height) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return lipgloss.NewStyle().
		Width(width).
		Background(lipgloss.Color(colors.Background)).
		MarginTop(y).
		MarginLeft(x).
		Render(content)
}
