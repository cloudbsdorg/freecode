package component

import "github.com/freecode/freecode/internal/renderer"

type Message struct {
	Role      string
	Content   string
	Timestamp string
}

type MessageList[R renderer.Renderer] struct {
	Component[R]
	Messages    []Message
	ScrollOff  int
	Colors     MessageListColors
}

type MessageListColors struct {
	Background    string
	UserBg        string
	UserFg        string
	AssistantBg   string
	AssistantFg   string
	SystemBg     string
	SystemFg     string
	MutedColor   string
}

func NewMessageList[R renderer.Renderer](width, height int, colors MessageListColors) *MessageList[R] {
	return &MessageList[R]{
		Component: Component[R]{
			X:       0,
			Y:       0,
			Width:   width,
			Height:  height,
			Visible: true,
		},
		Messages:   []Message{},
		ScrollOff: 0,
		Colors:    colors,
	}
}

func (m *MessageList[R]) AddMessage(role, content, timestamp string) {
	m.Messages = append(m.Messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: timestamp,
	})
}

func (m *MessageList[R]) ScrollUp() {
	if m.ScrollOff > 0 {
		m.ScrollOff--
	}
}

func (m *MessageList[R]) ScrollDown() {
	maxScroll := len(m.Messages) - m.Height + 2
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.ScrollOff < maxScroll {
		m.ScrollOff++
	}
}

func (m *MessageList[R]) ScrollToTop() {
	m.ScrollOff = 0
}

func (m *MessageList[R]) ScrollToBottom() {
	m.ScrollOff = len(m.Messages) - m.Height + 2
	if m.ScrollOff < 0 {
		m.ScrollOff = 0
	}
}

func (m *MessageList[R]) Clear() {
	m.Messages = []Message{}
	m.ScrollOff = 0
}

func (m *MessageList[R]) Render(r R) string {
	if !m.Visible {
		return ""
	}

	lines := []string{}

	start := m.ScrollOff
	end := start + m.Height - 2
	if end > len(m.Messages) {
		end = len(m.Messages)
	}
	if start >= end {
		start = 0
		end = 0
	}

	for i := start; i < end; i++ {
		msg := m.Messages[i]
		var fg, bg string

		switch msg.Role {
		case "user":
			fg = m.Colors.UserFg
			bg = m.Colors.UserBg
		case "assistant":
			fg = m.Colors.AssistantFg
			bg = m.Colors.AssistantBg
		default:
			fg = m.Colors.SystemFg
			bg = m.Colors.SystemBg
		}

		prefix := msg.Role + ": "
		text := truncate(msg.Content, m.Width-len(prefix)-2)

		if msg.Timestamp != "" {
			text += " (" + msg.Timestamp + ")"
		}

		lines = append(lines, r.RenderSelected(prefix+text, m.X+1, m.Y+1+(i-start), m.Width-2, fg, bg))
	}

	for i := len(lines); i < m.Height-2; i++ {
		lines = append(lines, r.RenderText("", m.X+1, m.Y+1+i, m.Colors.MutedColor))
	}

	result := ""
	for _, line := range lines {
		result += line + "\n"
	}
	return r.RenderBox(m.X, m.Y, m.Width, m.Height, m.Colors.Background) + result
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
