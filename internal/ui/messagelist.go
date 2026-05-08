package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var MessageContainerStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#1E1E1E")).
	Foreground(lipgloss.Color("#E0E0E0")).
	Padding(1, 2)

var UserMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#4EC9B0")).
	Bold(true)

var AssistantMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#DCDCAA"))

var SystemMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#808080")).
	Italic(true)

var ToolMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#CE9178")).
	Background(lipgloss.Color("#2D2D2D")).
	Padding(0, 1)

var TimestampStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#606060"))

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
	style := UserMessageStyle

	switch msg.Role {
	case "user":
		roleLabel = "You"
		style = UserMessageStyle
	case "assistant":
		roleLabel = "Assistant"
		style = AssistantMessageStyle
	case "system":
		roleLabel = "System"
		style = SystemMessageStyle
	case "tool":
		roleLabel = "Tool"
		style = ToolMessageStyle
	}

	var timestampStr string
	if m.showTimestamps {
		timestampStr = " " + TimestampStyle.Render(msg.Timestamp.Format("15:04"))
	}

	var result strings.Builder
	result.WriteString(style.Render(roleLabel) + ": ")

	if len(msg.Parts) > 0 {
		for _, part := range msg.Parts {
			switch part.Type {
			case "text":
				content := part.Content
				if len(content) > m.width-15 {
					content = content[:m.width-18] + "..."
				}
				result.WriteString(lipgloss.NewStyle().Render(content))
			case "reasoning":
				content := part.Content
				if len(content) > m.width-20 {
					content = content[:m.width-23] + "..."
				}
				reasoningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Italic(true)
				result.WriteString(reasoningStyle.Render("[Thinking: " + content + "]"))
			case "tool":
				toolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CE9178"))
				if part.Tool != "" {
					result.WriteString(toolStyle.Render("[Tool: " + part.Tool + "]"))
				} else {
					result.WriteString(toolStyle.Render("[Tool]"))
				}
			default:
				content := part.Content
				if len(content) > m.width-15 {
					content = content[:m.width-18] + "..."
				}
				result.WriteString(lipgloss.NewStyle().Render(content))
			}
		}
	} else {
		content := msg.Content
		if len(content) > m.width-10 {
			content = content[:m.width-13] + "..."
		}
		result.WriteString(lipgloss.NewStyle().Render(content))
	}

	result.WriteString(timestampStr)
	result.WriteString("\n")
	return result.String()
}
