package ui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/freecode/freecode/internal/style"
)

type InputHandler struct{}

func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

type CursorMode int

const (
	CursorModeBlock CursorMode = iota
	CursorModeUnderline
	CursorModeBar
)

type InputState struct {
	Value       string
	Cursor      int
	CursorMode  CursorMode
	Multiline   bool
	Placeholder string
	Focused     bool
}

func NewInputState() *InputState {
	return &InputState{
		Cursor:     0,
		CursorMode: CursorModeBlock,
	}
}

func (s *InputState) Insert(char rune) {
	if s.Cursor >= len(s.Value) {
		s.Value += string(char)
	} else {
		s.Value = s.Value[:s.Cursor] + string(char) + s.Value[s.Cursor:]
	}
	s.Cursor++
}

func (s *InputState) Delete() {
	if s.Cursor >= len(s.Value) || s.Cursor < 0 {
		return
	}
	s.Value = s.Value[:s.Cursor] + s.Value[s.Cursor+1:]
}

func (s *InputState) MoveLeft() {
	if s.Cursor > 0 {
		s.Cursor--
	}
}

func (s *InputState) MoveRight() {
	if s.Cursor < len(s.Value) {
		s.Cursor++
	}
}

func (s *InputState) MoveToStart() {
	s.Cursor = 0
}

func (s *InputState) MoveToEnd() {
	s.Cursor = len(s.Value)
}

func (s *InputState) DeleteWord() {
	if s.Cursor == 0 {
		return
	}

	start := s.Cursor
	for start > 0 && s.Value[start-1] == ' ' {
		start--
	}
	for start > 0 && s.Value[start-1] != ' ' {
		start--
	}

	s.Value = s.Value[:start] + s.Value[s.Cursor:]
	s.Cursor = start
}

func (h *InputHandler) Handle(msg tea.KeyMsg, state *InputState) {
	switch msg.Type {
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			state.Insert(r)
		}
	case tea.KeyLeft:
		state.MoveLeft()
	case tea.KeyRight:
		state.MoveRight()
	case tea.KeyHome:
		state.MoveToStart()
	case tea.KeyEnd:
		state.MoveToEnd()
	case tea.KeyBackspace:
		if state.Cursor > 0 {
			state.Cursor--
			state.Delete()
		}
	case tea.KeyDelete:
		state.Delete()
	case tea.KeyCtrlW:
		state.DeleteWord()
	}
}

func (h *InputHandler) Render(state *InputState) string {
	if !state.Focused && state.Value == "" {
		return style.NewStyle().
			Foreground(style.Color("#808080")).
			Render(state.Placeholder)
	}

	before := state.Value[:state.Cursor]
	after := state.Value[state.Cursor:]

	cursor := style.NewStyle().
		Reverse(true).
		Render(" ")

	if state.CursorMode == CursorModeUnderline {
		cursor = style.NewStyle().
			Underline(true).
			Render(" ")
	} else if state.CursorMode == CursorModeBar {
		cursor = style.NewStyle().
			Foreground(style.Color("#FFFFFF")).
			Background(style.Color("#007ACC")).
			Render(" ")
	}

	return before + cursor + after
}
