package ui

import (
	"strings"
	"time"

	"github.com/freecode/freecode/internal/style"
)

var MessageContainerStyle = style.NewStyle().
	Background(style.Color("#1E1E1E")).
	Foreground(style.Color("#E0E0E0")).
	Padding(1, 2)

var UserMessageStyle = style.NewStyle().
	Foreground(style.Color("#4EC9B0")).
	Bold(true)

var AssistantMessageStyle = style.NewStyle().
	Foreground(style.Color("#DCDCAA"))

var SystemMessageStyle = style.NewStyle().
	Foreground(style.Color("#808080")).
	Italic(true)

var ToolMessageStyle = style.NewStyle().
	Foreground(style.Color("#CE9178")).
	Background(style.Color("#2D2D2D")).
	Padding(0, 1)

var TimestampStyle = style.NewStyle().
	Foreground(style.Color("#606060"))

type Message struct {
	ID        string
	Role      string
	Content   string
	Timestamp time.Time
	Parts     []MessagePart
}

type MessagePart struct {
	Type    string
	Content string
	Tool    string
}

type MessageList struct {
	messages []Message
	width    int
	height   int
	scrollPos int
	showTimestamps bool
}

func NewMessageList() *MessageList {
	return &MessageList{
		messages:     make([]Message, 0),
		width:        80,
		height:       20,
		scrollPos:    0,
		showTimestamps: false,
	}
}

func (m *MessageList) AddMessage(msg Message) {
	m.messages = append(m.messages, msg)
}

func (m *MessageList) SetMessages(msgs []Message) {
	m.messages = msgs
}

func (m *MessageList) Clear() {
	m.messages = make([]Message, 0)
}

func (m *MessageList) RemoveLastAssistant() {
	for i := len(m.messages) - 1; i >= 0; i-- {
		if m.messages[i].Role == "assistant" {
			m.messages = append(m.messages[:i], m.messages[i+1:]...)
			return
		}
	}
}

func (m *MessageList) GetMessages() []Message {
	return m.messages
}

func (m *MessageList) SetWidth(w int) {
	m.width = w
}

func (m *MessageList) SetHeight(h int) {
	m.height = h
}

func (m *MessageList) ScrollUp() {
	if m.scrollPos > 0 {
		m.scrollPos--
	}
}

func (m *MessageList) ScrollDown() {
	maxScroll := len(m.messages) - m.height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollPos < maxScroll {
		m.scrollPos++
	}
}

func (m *MessageList) ScrollToBottom() {
	m.scrollPos = 0
}

func (m *MessageList) scrollToTop() {
	maxScroll := len(m.messages) - m.height
	if maxScroll < 0 {
		maxScroll = 0
	}
	m.scrollPos = maxScroll
}

func (m *MessageList) ToggleTimestamps() {
	m.showTimestamps = !m.showTimestamps
}

func (m *MessageList) Render() string {
	if len(m.messages) == 0 {
		return MessageContainerStyle.Render("\n\n  No messages yet. Start a conversation!\n\n")
	}

	result := ""
	visibleMessages := m.messages
	if len(visibleMessages) > m.height {
		start := len(visibleMessages) - m.height
		if m.scrollPos > 0 && m.scrollPos < start {
			start = m.scrollPos
		}
		visibleMessages = visibleMessages[start:]
	}

	for _, msg := range visibleMessages {
		result += m.renderMessage(msg)
	}

	return MessageContainerStyle.Render(result)
}

func (m *MessageList) renderMessage(msg Message) string {
	roleLabel := ""
	labelStyle := UserMessageStyle

	switch msg.Role {
	case "user":
		roleLabel = "You"
		labelStyle = UserMessageStyle
	case "assistant":
		roleLabel = "Assistant"
		labelStyle = AssistantMessageStyle
	case "system":
		roleLabel = "System"
		labelStyle = SystemMessageStyle
	case "tool":
		roleLabel = "Tool"
		labelStyle = ToolMessageStyle
	}

	var timestampStr string
	if m.showTimestamps {
		timestampStr = " " + TimestampStyle.Render(msg.Timestamp.Format("15:04"))
	}

	var result strings.Builder
	result.WriteString(labelStyle.Render(roleLabel) + ": ")

	if len(msg.Parts) > 0 {
		for _, part := range msg.Parts {
			partContent := RenderPart(part, m.width-10)
			result.WriteString(partContent)
		}
	} else {
		content := wordWrap(msg.Content, m.width-10)
		result.WriteString(style.NewStyle().Render(content))
	}

	result.WriteString(timestampStr)
	result.WriteString("\n")
	return result.String()
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for lineIdx, line := range lines {
		if lineIdx > 0 {
			result.WriteString("\n")
		}

		if len(line) <= width {
			result.WriteString(line)
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}

		currentLine := ""
		for _, word := range words {
			if len(currentLine)+len(word)+1 <= width {
				if currentLine != "" {
					currentLine += " "
				}
				currentLine += word
			} else {
				if currentLine != "" {
					result.WriteString(currentLine)
					result.WriteString("\n")
				}
				for len(word) > width {
					result.WriteString(word[:width])
					result.WriteString("\n")
					word = word[width:]
				}
				currentLine = word
			}
		}
		if currentLine != "" {
			result.WriteString(currentLine)
		}
	}

	return result.String()
}
