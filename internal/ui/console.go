package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelDebug LogLevel = "debug"
)

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

type ConsolePanel struct {
	width     int
	height    int
	isOpen    bool
	logs      []LogEntry
	maxLogs   int
}

func NewConsolePanel() *ConsolePanel {
	return &ConsolePanel{
		width:   80,
		height:  20,
		isOpen:  false,
		logs:    []LogEntry{},
		maxLogs: 100,
	}
}

func (c *ConsolePanel) SetWidth(w int) {
	c.width = w
}

func (c *ConsolePanel) SetHeight(h int) {
	c.height = h
}

func (c *ConsolePanel) Toggle() {
	c.isOpen = !c.isOpen
}

func (c *ConsolePanel) Open() {
	c.isOpen = true
}

func (c *ConsolePanel) Close() {
	c.isOpen = false
}

func (c *ConsolePanel) IsOpen() bool {
	return c.isOpen
}

func (c *ConsolePanel) HandleKey(msg string) bool {
	if !c.isOpen {
		return false
	}
	if msg == "escape" || msg == "c" {
		c.Close()
		return true
	}
	return false
}

func (c *ConsolePanel) Log(level LogLevel, format string, args ...interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   strings.TrimSpace(fmt.Sprintf(format, args...)),
	}
	c.logs = append(c.logs, entry)
	if len(c.logs) > c.maxLogs {
		c.logs = c.logs[len(c.logs)-c.maxLogs:]
	}
}

func (c *ConsolePanel) Info(format string, args ...interface{}) {
	c.Log(LogLevelInfo, format, args...)
}

func (c *ConsolePanel) Warn(format string, args ...interface{}) {
	c.Log(LogLevelWarn, format, args...)
}

func (c *ConsolePanel) Error(format string, args ...interface{}) {
	c.Log(LogLevelError, format, args...)
}

func (c *ConsolePanel) Debug(format string, args ...interface{}) {
	c.Log(LogLevelDebug, format, args...)
}

func (c *ConsolePanel) Clear() {
	c.logs = []LogEntry{}
}

func (c *ConsolePanel) Render() string {
	if !c.isOpen {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#0D0D0D")).
		Border(lipgloss.HiddenBorder()).
		Width(c.width).
		Height(c.height)

	return dialogStyle.Render(c.renderContent())
}

func (c *ConsolePanel) renderContent() string {
	var lines []string

	lines = append(lines, c.renderHeader())
	lines = append(lines, c.renderLogs()...)

	return strings.Join(lines, "\n")
}

func (c *ConsolePanel) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(0, 1)

	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	return headerStyle.Render("Console") + " " + hintStyle.Render("esc to close")
}

func (c *ConsolePanel) renderLogs() []string {
	var lines []string

	if len(c.logs) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		lines = append(lines, emptyStyle.Render("  (no logs)"))
		return lines
	}

	sort.Slice(c.logs, func(i, j int) bool {
		return c.logs[i].Timestamp.Before(c.logs[j].Timestamp)
	})

	for _, entry := range c.logs {
		lines = append(lines, c.renderEntry(entry))
	}

	if len(lines) > c.height-2 {
		lines = lines[len(lines)-(c.height-2):]
	}

	return lines
}

func (c *ConsolePanel) renderEntry(entry LogEntry) string {
	timeStr := entry.Timestamp.Format("15:04:05.000")
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#606060"))

	var levelStyle lipgloss.Style
	var levelStr string
	switch entry.Level {
	case LogLevelInfo:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))
		levelStr = "INFO"
	case LogLevelWarn:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DCDCAA"))
		levelStr = "WARN"
	case LogLevelError:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F44747"))
		levelStr = "ERROR"
	case LogLevelDebug:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		levelStr = "DEBUG"
	}

	msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))

	return timeStyle.Render(timeStr) + " " + levelStyle.Render(levelStr) + " " + msgStyle.Render(entry.Message)
}