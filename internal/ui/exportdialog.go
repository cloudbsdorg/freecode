package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ExportFormat string

const (
	ExportFormatMarkdown ExportFormat = "markdown"
	ExportFormatJSON    ExportFormat = "json"
	ExportFormatText    ExportFormat = "text"
)

type ExportOptions struct {
	Format           ExportFormat
	IncludeTimestamps bool
	IncludeMetadata  bool
}

type ExportDialog struct {
	width         int
	isOpen        bool
	format        ExportFormat
	includeTime   bool
	includeMeta   bool
	selected      int
	sessionTitle  string
	messages      []string
}

func NewExportDialog() *ExportDialog {
	return &ExportDialog{
		width:         60,
		isOpen:        false,
		format:        ExportFormatMarkdown,
		includeTime:   true,
		includeMeta:    true,
		selected:      0,
		sessionTitle:  "Session",
		messages:      []string{},
	}
}

func (e *ExportDialog) SetWidth(w int) {
	e.width = w
}

func (e *ExportDialog) SetSession(title string, messages []string) {
	e.sessionTitle = title
	e.messages = messages
}

func (e *ExportDialog) Open() {
	e.isOpen = true
	e.selected = 0
}

func (e *ExportDialog) Close() {
	e.isOpen = false
}

func (e *ExportDialog) IsOpen() bool {
	return e.isOpen
}

func (e *ExportDialog) HandleKey(msg string) bool {
	if !e.isOpen {
		return false
	}

	switch msg {
	case "escape":
		e.Close()
		return true
	case "enter":
		e.export()
		return true
	case "up", "k":
		if e.selected > 0 {
			e.selected--
		}
		return true
	case "down", "j":
		if e.selected < 3 {
			e.selected++
		}
		return true
	case "left", "h":
		e.toggleOption(-1)
		return true
	case "right", "l":
		e.toggleOption(1)
		return true
	}
	return false
}

func (e *ExportDialog) toggleOption(dir int) {
	switch e.selected {
	case 0:
		if dir > 0 {
			if e.format == ExportFormatMarkdown {
				e.format = ExportFormatJSON
			} else if e.format == ExportFormatJSON {
				e.format = ExportFormatText
			} else {
				e.format = ExportFormatMarkdown
			}
		} else {
			if e.format == ExportFormatMarkdown {
				e.format = ExportFormatText
			} else if e.format == ExportFormatText {
				e.format = ExportFormatJSON
			} else {
				e.format = ExportFormatMarkdown
			}
		}
	case 1:
		e.includeTime = !e.includeTime
	case 2:
		e.includeMeta = !e.includeMeta
	}
}

func (e *ExportDialog) export() {
	e.Close()
}

func (e *ExportDialog) GetOptions() ExportOptions {
	return ExportOptions{
		Format:           e.format,
		IncludeTimestamps: e.includeTime,
		IncludeMetadata:  e.includeMeta,
	}
}

func (e *ExportDialog) Render() string {
	if !e.isOpen {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E1E")).
		Border(lipgloss.HiddenBorder()).
		Width(e.width)

	return dialogStyle.Render(e.renderContent())
}

func (e *ExportDialog) renderContent() string {
	var lines []string

	lines = append(lines, e.renderHeader())
	lines = append(lines, "")
	lines = append(lines, e.renderOption("Format", string(e.format), 0))
	lines = append(lines, e.renderToggle("Timestamps", e.includeTime, 1))
	lines = append(lines, e.renderToggle("Metadata", e.includeMeta, 2))
	lines = append(lines, "")
	lines = append(lines, e.renderPreview()...)
	lines = append(lines, "")
	lines = append(lines, e.renderHints())

	return strings.Join(lines, "\n")
}

func (e *ExportDialog) renderHeader() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Bold(true).
		Render("Export Session")
}

func (e *ExportDialog) renderOption(title, value string, idx int) string {
	selected := e.selected == idx
	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	}

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))

	return prefix + " " + titleStyle.Render(title+":") + " " + valueStyle.Render(value)
}

func (e *ExportDialog) renderToggle(title string, enabled bool, idx int) string {
	selected := e.selected == idx
	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	}

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	checkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))
	disabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F44747"))

	check := "✗"
	valueStyle := disabledStyle
	if enabled {
		check = "✓"
		valueStyle = checkStyle
	}

	return prefix + " " + titleStyle.Render(title+":") + " " + valueStyle.Render(check)
}

func (e *ExportDialog) renderPreview() []string {
	var lines []string
	previewStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	lines = append(lines, "  "+previewStyle.Render("Preview:"))

	switch e.format {
	case ExportFormatMarkdown:
		lines = append(lines, "  # "+e.sessionTitle)
		if e.includeTime {
			lines = append(lines, "  *timestamp*")
		}
		lines = append(lines, "  ## Messages")
	case ExportFormatJSON:
		lines = append(lines, `  {"title": "...", "messages": []}`)
	case ExportFormatText:
		lines = append(lines, "  Session: "+e.sessionTitle)
	}

	return lines
}

func (e *ExportDialog) renderHints() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	return hintStyle.Render("↑↓ select  ←→ change  enter export  esc cancel")
}