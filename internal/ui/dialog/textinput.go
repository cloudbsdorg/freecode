package dialog

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TextInputOption func(*TextInput)

func TextInputWithMaxLen(max int) TextInputOption {
	return func(t *TextInput) {
		t.MaxLen = max
	}
}

func TextInputWithHidden(hidden bool) TextInputOption {
	return func(t *TextInput) {
		t.Hidden = hidden
	}
}

func TextInputWithPlaceholder(placeholder string) TextInputOption {
	return func(t *TextInput) {
		t.Placeholder = placeholder
	}
}

func TextInputWithColors(colors Colors) TextInputOption {
	return func(t *TextInput) {
		t.Colors = colors
	}
}

func TextInputWithOnChange(fn func(string)) TextInputOption {
	return func(t *TextInput) {
		t.OnChange = fn
	}
}

func TextInputWithOnSubmit(fn func(string)) TextInputOption {
	return func(t *TextInput) {
		t.OnSubmit = fn
	}
}

func TextInputWithOnCancel(fn func()) TextInputOption {
	return func(t *TextInput) {
		t.OnCancel = fn
	}
}

type TextInput struct {
	Value       string
	Placeholder string
	Colors      Colors
	MaxLen      int
	Hidden      bool
	Cursor      int
	Focused     bool
	OnChange    func(string)
	OnSubmit    func(string)
	OnCancel    func()
}

func NewTextInput(opts ...TextInputOption) *TextInput {
	t := &TextInput{
		Value:       "",
		Placeholder: "",
		Colors:      DefaultColors,
		MaxLen:      0,
		Hidden:      false,
		Cursor:      0,
		Focused:     true,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *TextInput) SetValue(value string) {
	if t.MaxLen > 0 && len(value) > t.MaxLen {
		value = value[:t.MaxLen]
	}
	t.Value = value
	if t.Cursor > len(t.Value) {
		t.Cursor = len(t.Value)
	}
	if t.OnChange != nil {
		t.OnChange(t.Value)
	}
}

func (t *TextInput) GetValue() string {
	return t.Value
}

func (t *TextInput) Clear() {
	t.Value = ""
	t.Cursor = 0
	if t.OnChange != nil {
		t.OnChange(t.Value)
	}
}

func (t *TextInput) Append(ch rune) {
	if t.MaxLen > 0 && len(t.Value) >= t.MaxLen {
		return
	}
	if t.Cursor == len(t.Value) {
		t.Value += string(ch)
	} else {
		t.Value = t.Value[:t.Cursor] + string(ch) + t.Value[t.Cursor:]
	}
	t.Cursor++
	if t.OnChange != nil {
		t.OnChange(t.Value)
	}
}

func (t *TextInput) Backspace() {
	if t.Cursor == 0 {
		return
	}
	t.Value = t.Value[:t.Cursor-1] + t.Value[t.Cursor:]
	t.Cursor--
	if t.OnChange != nil {
		t.OnChange(t.Value)
	}
}

func (t *TextInput) Delete() {
	if t.Cursor >= len(t.Value) {
		return
	}
	t.Value = t.Value[:t.Cursor] + t.Value[t.Cursor+1:]
	if t.OnChange != nil {
		t.OnChange(t.Value)
	}
}

func (t *TextInput) MoveLeft() {
	if t.Cursor > 0 {
		t.Cursor--
	}
}

func (t *TextInput) MoveRight() {
	if t.Cursor < len(t.Value) {
		t.Cursor++
	}
}

func (t *TextInput) MoveToStart() {
	t.Cursor = 0
}

func (t *TextInput) MoveToEnd() {
	t.Cursor = len(t.Value)
}

func (t *TextInput) Submit() {
	if t.OnSubmit != nil {
		t.OnSubmit(t.Value)
	}
}

func (t *TextInput) Cancel() {
	if t.OnCancel != nil {
		t.OnCancel()
	}
}

func (t *TextInput) RenderDisplay() string {
	display := t.Value
	if t.Hidden {
		display = strings.Repeat("•", len(t.Value))
	}
	return display
}

func (t *TextInput) Render() string {
	return t.RenderWithPrefix("")
}

func (t *TextInput) RenderWithPrefix(prefix string) string {
	display := t.RenderDisplay()

	if len(display) > 50 {
		headLen := 23
		tailLen := 24
		if headLen+tailLen > len(display) {
			headLen = len(display) - tailLen
		}
		display = display[:headLen] + "..." + display[len(display)-tailLen:]
	}

	prefixStr := ""
	if prefix != "" {
		prefixStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Text)).
			Render(prefix + " ")
	}

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Colors.Text)).
		Background(lipgloss.Color(t.Colors.BackgroundAlt)).
		Padding(0, 1)

	if t.Cursor >= len(t.Value) {
		return prefixStr + inputStyle.Render(display + "_")
	}

	before := display[:t.Cursor]
	after := display[t.Cursor:]
	return prefixStr + lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Colors.Text)).
		Background(lipgloss.Color(t.Colors.BackgroundAlt)).
		Padding(0, 1).
		Render(before + "_" + after)
}

func (t *TextInput) RenderLabeled(label string) string {
	prefixStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Colors.TextMuted)).
		Render(label + " ")

	display := t.RenderDisplay()

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Colors.Text)).
		Background(lipgloss.Color(t.Colors.BackgroundAlt)).
		Padding(0, 1)

	if t.Cursor >= len(t.Value) {
		return prefixStr + inputStyle.Render(display + "_")
	}

	before := display[:t.Cursor]
	after := display[t.Cursor:]
	return prefixStr + lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Colors.Text)).
		Background(lipgloss.Color(t.Colors.BackgroundAlt)).
		Padding(0, 1).
		Render(before + "_" + after)
}
