package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbletea"
)

var InputContainerStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#1E1E1E")).
	Foreground(lipgloss.Color("#E0E0E0")).
	BorderForeground(lipgloss.Color("#3D3D3D")).
	BorderStyle(lipgloss.NormalBorder()).
	Padding(0, 1)

var InputStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF"))

var PlaceholderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#606060"))

var PromptStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#4EC9B0")).
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
		placeholder: "Type a message...",
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
	display := in.value
	placeholder := in.placeholder

	if display == "" && !in.focused {
		placeholder = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#606060")).
			Italic(true).
			Render(placeholder)
		display = placeholder
	}

	before := ""
	after := ""
	if in.cursorPos <= len(in.value) {
		before = in.value[:in.cursorPos]
		after = in.value[in.cursorPos:]
	} else {
		before = in.value
		after = ""
	}

	var cursor string
	if in.focused {
		cursor = lipgloss.NewStyle().
			Background(lipgloss.Color("#007ACC")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" ")
	} else {
		cursor = " "
	}

	promptStr := PromptStyle.Render(in.prompt)
	inputBefore := InputStyle.Render(before)
	inputAfter := InputStyle.Render(after)

	line := promptStr + inputBefore + cursor + inputAfter

	availWidth := in.width - 4
	if len(line) > availWidth && availWidth > 0 {
		prefixLen := lipgloss.Width(promptStr)
		inputLen := lipgloss.Width(inputBefore) + 1 + lipgloss.Width(inputAfter)
		if prefixLen + inputLen > availWidth {
			maxInput := availWidth - prefixLen - 1
			if maxInput > 0 {
				if lipgloss.Width(inputBefore) > maxInput {
					cutoff := lipgloss.Width(inputBefore) - maxInput
					if cutoff < lipgloss.Width(inputBefore) {
						inputBefore = inputBefore[cutoff:]
					}
				}
			}
		}
		line = promptStr + inputBefore + cursor + inputAfter
	}

	return InputContainerStyle.Render(line)
}
