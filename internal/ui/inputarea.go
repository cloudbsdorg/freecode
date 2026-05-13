package ui

import (
	"github.com/freecode/freecode/internal/style"
	"github.com/charmbracelet/bubbletea"
)

var InputContainerStyle = style.NewStyle().
	Background(style.Color("#1E1E1E")).
	Foreground(style.Color("#E0E0E0")).
	BorderForeground(style.Color("#3D3D3D")).
	BorderStyle(style.NormalBorder()).
	Padding(0, 1)

var InputStyle = style.NewStyle().
	Foreground(style.Color("#FFFFFF"))

var PlaceholderStyle = style.NewStyle().
	Foreground(style.Color("#606060"))

var PromptStyle = style.NewStyle().
	Foreground(style.Color("#4EC9B0")).
	Bold(true)

type InputArea struct {
	value       string
	cursorPos   int
	prompt      string
	placeholder string
	focused     bool
	width       int
	height      int
	history     []string
	historyIdx  int
}

func NewInputArea() *InputArea {
	return &InputArea{
		value:       "",
		cursorPos:   0,
		prompt:      "> ",
		placeholder: "Ask anything... (e.g., 'Fix a bug', 'Explain this code')",
		focused:     true,
		width:       80,
		height:      3,
		history:     make([]string, 0),
		historyIdx:  -1,
	}
}

func (in *InputArea) SetValue(value string) {
	in.value = value
	in.cursorPos = len(value)
}

func (in *InputArea) GetValue() string {
	return in.value
}

func (in *InputArea) SetPrompt(prompt string) {
	in.prompt = prompt
}

func (in *InputArea) SetPlaceholder(placeholder string) {
	in.placeholder = placeholder
}

func (in *InputArea) SetWidth(w int) {
	in.width = w
}

func (in *InputArea) SetHeight(h int) {
	in.height = h
}

func (in *InputArea) Focus() {
	in.focused = true
}

func (in *InputArea) Blur() {
	in.focused = false
}

func (in *InputArea) HandleKey(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			if in.cursorPos >= len(in.value) {
				in.value += string(r)
			} else {
				in.value = in.value[:in.cursorPos] + string(r) + in.value[in.cursorPos:]
			}
			in.cursorPos++
		}
		in.historyIdx = -1

	case tea.KeyLeft:
		if in.cursorPos > 0 {
			in.cursorPos--
		}

	case tea.KeyRight:
		if in.cursorPos < len(in.value) {
			in.cursorPos++
		}

	case tea.KeyHome:
		in.cursorPos = 0

	case tea.KeyEnd:
		in.cursorPos = len(in.value)

	case tea.KeyBackspace:
		if in.cursorPos > 0 {
			in.value = in.value[:in.cursorPos-1] + in.value[in.cursorPos:]
			in.cursorPos--
		}

	case tea.KeyDelete:
		if in.cursorPos < len(in.value) {
			in.value = in.value[:in.cursorPos] + in.value[in.cursorPos+1:]
		}

	case tea.KeyUp:
		if len(in.history) > 0 {
			if in.historyIdx < len(in.history)-1 {
				in.historyIdx++
				in.value = in.history[len(in.history)-1-in.historyIdx]
				in.cursorPos = len(in.value)
			}
		}

	case tea.KeyDown:
		if in.historyIdx > 0 {
			in.historyIdx--
			in.value = in.history[len(in.history)-1-in.historyIdx]
			in.cursorPos = len(in.value)
		} else if in.historyIdx == 0 {
			in.historyIdx = -1
			in.value = ""
			in.cursorPos = 0
		}

	case tea.KeyEnter:
		if in.value != "" {
			in.history = append(in.history, in.value)
			in.historyIdx = -1
		}
	}
}

func (in *InputArea) Submit() string {
	value := in.value
	in.value = ""
	in.cursorPos = 0
	return value
}

func (in *InputArea) Render() string {
	promptStr := PromptStyle.Render(in.prompt)
	promptLen := len(in.prompt)

	display := in.value
	if display == "" {
		placeholder := style.NewStyle().
			Foreground(style.Color("#606060")).
			Italic(true).
			Render(in.placeholder)
		display = placeholder
	}

	before := display[:in.cursorPos]
	after := display[in.cursorPos:]

	inputBefore := InputStyle.Render(before)
	inputAfter := InputStyle.Render(after)

	var cursor string
	if in.focused {
		cursor = style.NewStyle().
			Background(style.Color("#007ACC")).
			Foreground(style.Color("#FFFFFF")).
			Render(" ")
	} else {
		cursor = " "
	}

	maxWidth := in.width - 2
	totalLen := promptLen + len(before) + 1 + len(after)

	if totalLen > maxWidth && maxWidth > promptLen {
		availInput := maxWidth - promptLen - 1
		if availInput < 0 {
			availInput = 0
		}
		if len(before) > availInput {
			before = before[len(before)-availInput:]
			if len(before) > 3 {
				before = "..." + before[3:]
			}
			after = ""
		} else if len(before)+len(after) > availInput {
			availAfter := availInput - len(before)
			if availAfter < 0 {
				availAfter = 0
			}
			if availAfter < len(after) && availAfter > 3 {
				after = after[:availAfter-3] + "..."
			} else if availAfter <= 3 {
				after = ""
			}
		}
		inputBefore = InputStyle.Render(before)
		inputAfter = InputStyle.Render(after)
	}

	line := promptStr + inputBefore + cursor + inputAfter

	style := InputContainerStyle.Width(in.width)
	return style.Render(line)
}
