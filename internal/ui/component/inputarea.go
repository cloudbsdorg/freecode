package component

import "github.com/freecode/freecode/internal/renderer"

type InputArea[R renderer.Renderer] struct {
	Component[R]
	Lines     []string
	CursorX   int
	CursorY   int
	Prompt    string
	Colors    InputAreaColors
}

type InputAreaColors struct {
	Background    string
	Foreground    string
	CursorColor   string
	MutedColor    string
}

func NewInputArea[R renderer.Renderer](width, height int, colors InputAreaColors) *InputArea[R] {
	return &InputArea[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  height,
			Visible: true,
		},
		Lines:   []string{""},
		CursorX: 0,
		CursorY: 0,
		Prompt:  "> ",
		Colors:  colors,
	}
}

func (in *InputArea[R]) SetValue(value string) {
	in.Lines = []string{""}
	for _, line := range splitLines(value) {
		in.Lines = append(in.Lines, line)
	}
	if len(in.Lines) == 0 {
		in.Lines = []string{""}
	}
	in.CursorX = 0
	in.CursorY = len(in.Lines) - 1
}

func (in *InputArea[R]) GetValue() string {
	result := ""
	for _, line := range in.Lines {
		if result != "" {
			result += "\n"
		}
		result += line
	}
	return result
}

func (in *InputArea[R]) Append(ch rune) {
	if in.CursorY >= len(in.Lines) {
		in.Lines = append(in.Lines, "")
	}
	line := in.Lines[in.CursorY]
	if in.CursorX >= len(line) {
		in.Lines[in.CursorY] += string(ch)
	} else {
		in.Lines[in.CursorY] = line[:in.CursorX] + string(ch) + line[in.CursorX:]
	}
	in.CursorX++
}

func (in *InputArea[R]) Backspace() {
	if in.CursorY >= len(in.Lines) || in.CursorX == 0 && in.CursorY == 0 {
		return
	}
	if in.CursorX > 0 {
		line := in.Lines[in.CursorY]
		in.Lines[in.CursorY] = line[:in.CursorX-1] + line[in.CursorX:]
		in.CursorX--
	} else if in.CursorY > 0 {
		in.CursorY--
		in.CursorX = len(in.Lines[in.CursorY])
	}
}

func (in *InputArea[R]) NewLine() {
	if in.CursorY < len(in.Lines)-1 {
		after := in.Lines[in.CursorY][in.CursorX:]
		in.Lines[in.CursorY] = in.Lines[in.CursorY][:in.CursorX]
		in.Lines = append(in.Lines[:in.CursorY+1], append([]string{after}, in.Lines[in.CursorY+1:]...)...)
	} else {
		in.Lines = append(in.Lines, "")
	}
	in.CursorY++
	in.CursorX = 0
}

func (in *InputArea[R]) MoveLeft() {
	if in.CursorX > 0 {
		in.CursorX--
	} else if in.CursorY > 0 {
		in.CursorY--
		in.CursorX = len(in.Lines[in.CursorY])
	}
}

func (in *InputArea[R]) MoveRight() {
	if in.CursorX < len(in.Lines[in.CursorY]) {
		in.CursorX++
	} else if in.CursorY < len(in.Lines)-1 {
		in.CursorY++
		in.CursorX = 0
	}
}

func (in *InputArea[R]) MoveUp() {
	if in.CursorY > 0 {
		in.CursorY--
		if in.CursorX > len(in.Lines[in.CursorY]) {
			in.CursorX = len(in.Lines[in.CursorY])
		}
	}
}

func (in *InputArea[R]) MoveDown() {
	if in.CursorY < len(in.Lines)-1 {
		in.CursorY++
		if in.CursorX > len(in.Lines[in.CursorY]) {
			in.CursorX = len(in.Lines[in.CursorY])
		}
	}
}

func (in *InputArea[R]) Clear() {
	in.Lines = []string{""}
	in.CursorX = 0
	in.CursorY = 0
}

func (in *InputArea[R]) Render(r R) string {
	if !in.Visible {
		return ""
	}

	lines := []string{}
	for i, line := range in.Lines {
		display := in.Prompt + line
		if i == in.CursorY {
			if in.CursorX >= len(line) {
				display += "_"
			} else {
				display = display[:len(in.Prompt)+in.CursorX] + "_" + display[len(in.Prompt)+in.CursorX+1:]
			}
		}
		lines = append(lines, r.RenderText(display, in.X, in.Y+i, in.Colors.Foreground))
	}

	for i := len(lines); i < in.Height; i++ {
		lines = append(lines, r.RenderText(in.Prompt, in.X, in.Y+i, in.Colors.MutedColor))
	}

	result := ""
	for _, line := range lines {
		result += line + "\n"
	}
	return r.RenderBox(in.X, in.Y, in.Width, in.Height, in.Colors.Background) + result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	lines = append(lines, s[start:])
	return lines
}
