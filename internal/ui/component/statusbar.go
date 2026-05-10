package component

import "github.com/freecode/freecode/internal/renderer"

type StatusBarContent struct {
	LeftText  string
	RightText string
	Model     string
	Agent     string
	Yolo      bool
}

type StatusBar[R renderer.Renderer] struct {
	Component[R]
	Content StatusBarContent
	Colors  StatusBarColors
}

type StatusBarColors struct {
	Background string
	Foreground string
	Active    string
	Warning   string
	Error     string
}

func NewStatusBar[R renderer.Renderer](width int, colors StatusBarColors) *StatusBar[R] {
	return &StatusBar[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  1,
			Visible: true,
		},
		Content: StatusBarContent{
			LeftText:  "",
			RightText: "",
			Model:    "none",
			Agent:    "sisyphus",
			Yolo:     false,
		},
		Colors: colors,
	}
}

func (s *StatusBar[R]) SetContent(content StatusBarContent) {
	s.Content = content
}

func (s *StatusBar[R]) Render(r R) string {
	if !s.Visible {
		return ""
	}

	box := r.RenderBox(s.X, s.Y, s.Width, 1, s.Colors.Background)

	left := s.Content.LeftText
	if s.Content.Model != "" {
		left += " | " + s.Content.Model
	}
	if s.Content.Agent != "" {
		left += " | " + s.Content.Agent
	}
	if s.Content.Yolo {
		left += " | YOLO"
	}

	right := ""
	if s.Content.RightText != "" {
		right = s.Content.RightText
	}

	text := r.RenderText(left, s.X+1, s.Y, s.Colors.Foreground)
	if right != "" {
		text += r.RenderText(right, s.X+s.Width-len(right)-1, s.Y, s.Colors.Foreground)
	}

	return box + text
}
