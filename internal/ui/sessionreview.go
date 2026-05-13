package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/freecode/freecode/internal/session"
	"github.com/freecode/freecode/internal/ui/dialog"
)

type SessionReviewDialog struct {
	sessions     []*session.Session
	selectedIdx  int
	width        int
	height       int
	isOpen       bool
	colors       dialog.Colors
}

func NewSessionReviewDialog() *SessionReviewDialog {
	return &SessionReviewDialog{
		selectedIdx: 0,
		width:       60,
		height:      20,
		isOpen:      false,
		colors:      dialog.Dark,
	}
}

func (s *SessionReviewDialog) SetSessions(sessions []*session.Session) {
	s.sessions = sessions
	if s.selectedIdx >= len(sessions) {
		s.selectedIdx = 0
	}
}

func (s *SessionReviewDialog) Open() {
	s.isOpen = true
	s.selectedIdx = 0
}

func (s *SessionReviewDialog) Close() {
	s.isOpen = false
}

func (s *SessionReviewDialog) IsOpen() bool {
	return s.isOpen
}

func (s *SessionReviewDialog) SetWidth(w int) {
	s.width = w
}

func (s *SessionReviewDialog) SetHeight(h int) {
	s.height = h
}

func (s *SessionReviewDialog) Prev() {
	if s.selectedIdx > 0 {
		s.selectedIdx--
	}
}

func (s *SessionReviewDialog) Next() {
	if s.selectedIdx < len(s.sessions)-1 {
		s.selectedIdx++
	}
}

func (s *SessionReviewDialog) GetSelectedSession() *session.Session {
	if s.selectedIdx >= 0 && s.selectedIdx < len(s.sessions) {
		return s.sessions[s.selectedIdx]
	}
	return nil
}

func (s *SessionReviewDialog) HandleKey(key string) bool {
	if !s.isOpen {
		return false
	}

	switch key {
	case "q", "escape":
		s.Close()
		return true
	case "j", "down":
		s.Next()
		return true
	case "k", "up":
		s.Prev()
		return true
	}
	return false
}

func (s *SessionReviewDialog) Render() string {
	if !s.isOpen {
		return ""
	}

	var lines []string
	lines = append(lines, dialog.Header("Session History", s.colors))
	lines = append(lines, "")

	if len(s.sessions) == 0 {
		lines = append(lines, "  No sessions found")
		lines = append(lines, "")
		lines = append(lines, "  Press ESC to close")
	} else {
		lines = append(lines, s.renderSessionList()...)
		lines = append(lines, "")
		lines = append(lines, s.renderSessionPreview())
		lines = append(lines, "")
		lines = append(lines, "  ↑↓ Navigate | ENTER Fork | ESC Close")
	}

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(s.width).
		Height(s.height).
		Background(lipgloss.Color(s.colors.Background)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(s.colors.Border)).
		Padding(1).
		Render(content)
}

func (s *SessionReviewDialog) renderSessionList() []string {
	var lines []string
	count := 10
	if len(s.sessions) < count {
		count = len(s.sessions)
	}

	for i := 0; i < count; i++ {
		sess := s.sessions[i]
		prefix := "  "
		if i == s.selectedIdx {
			prefix = "▶ "
		}

		title := sess.Title
		if title == "" {
			title = "Untitled"
		}
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		date := sess.UpdatedAt.Format("01-02 15:04")
		msgCount := len(sess.Messages)

		line := prefix + title + " (" + date + ", " + itoa(msgCount) + " msgs)"
		if i == s.selectedIdx {
			lines = append(lines, dialog.Selected(line, s.colors))
		} else {
			lines = append(lines, dialog.Muted(line, s.colors))
		}
	}
	return lines
}

func (s *SessionReviewDialog) renderSessionPreview() string {
	sess := s.GetSelectedSession()
	if sess == nil {
		return ""
	}

	var preview []string
	preview = append(preview, "  ── Session Preview ──")
	preview = append(preview, "  Title: "+sess.Title)
	preview = append(preview, "  Model: "+sess.Model)
	preview = append(preview, "  Agent: "+sess.Agent)
	preview = append(preview, "  Created: "+sess.CreatedAt.Format(time.RFC1123))
	preview = append(preview, "  Updated: "+sess.UpdatedAt.Format(time.RFC1123))
	preview = append(preview, "  Messages: "+itoa(len(sess.Messages)))

	if len(sess.Messages) > 0 {
		lastMsg := sess.Messages[len(sess.Messages)-1]
		content := lastMsg.Content
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		content = strings.ReplaceAll(content, "\n", " ")
		preview = append(preview, "  Last: "+content)
	}

	return strings.Join(preview, "\n")
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var result []byte
	for i > 0 {
		result = append([]byte{byte('0' + i%10)}, result...)
		i /= 10
	}
	return string(result)
}
