package component

import "github.com/freecode/freecode/internal/renderer"

type TextInput[R renderer.Renderer] struct {
	Component[R]
	Value       string
	Placeholder string
	Cursor      int
	Hidden      bool
	MaxLen      int
	Colors      TextInputColors
}

type TextInputColors struct {
	Background    string
	Foreground    string
	CursorColor   string
}

func NewTextInput[R renderer.Renderer](width int, colors TextInputColors) *TextInput[R] {
	return &TextInput[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  1,
			Visible: true,
		},
		Colors: colors,
	}
}

func (t *TextInput[R]) SetValue(value string) {
	if t.MaxLen > 0 && len(value) > t.MaxLen {
		value = value[:t.MaxLen]
	}
	t.Value = value
	if t.Cursor > len(t.Value) {
		t.Cursor = len(t.Value)
	}
}

func (t *TextInput[R]) GetValue() string {
	return t.Value
}

func (t *TextInput[R]) Append(ch rune) {
	if t.MaxLen > 0 && len(t.Value) >= t.MaxLen {
		return
	}
	if t.Cursor == len(t.Value) {
		t.Value += string(ch)
	} else {
		t.Value = t.Value[:t.Cursor] + string(ch) + t.Value[t.Cursor:]
	}
	t.Cursor++
}

func (t *TextInput[R]) Backspace() {
	if t.Cursor == 0 {
		return
	}
	t.Value = t.Value[:t.Cursor-1] + t.Value[t.Cursor:]
	t.Cursor--
}

func (t *TextInput[R]) MoveLeft() {
	if t.Cursor > 0 {
		t.Cursor--
	}
}

func (t *TextInput[R]) MoveRight() {
	if t.Cursor < len(t.Value) {
		t.Cursor++
	}
}

func (t *TextInput[R]) Clear() {
	t.Value = ""
	t.Cursor = 0
}

func (t *TextInput[R]) Render(r R) string {
	if !t.Visible {
		return ""
	}

	display := t.Value
	if t.Hidden {
		display = ""
		for range t.Value {
			display += "*"
		}
	}

	prefix := ""
	if t.Placeholder != "" && t.Value == "" {
		display = t.Placeholder
		prefix = "> "
	} else {
		prefix = "> "
	}

	if t.Cursor == len(t.Value) {
		display = display + "_"
	} else {
		display = display[:t.Cursor] + "_" + display[t.Cursor:]
	}

	return r.RenderSelected(prefix+display, t.X, t.Y, t.Width, t.Colors.Foreground, t.Colors.Background)
}
