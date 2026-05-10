package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ToolResultComponent struct {
	width      int
	toolName   string
	status     string
	content    string
	error      string
	truncated  bool
	maxLines   int
}

func NewToolResultComponent() *ToolResultComponent {
	return &ToolResultComponent{
		width:    80,
		maxLines: 50,
		status:   ToolStatusCompleted,
	}
}

func (t *ToolResultComponent) SetWidth(w int) {
	t.width = w
}

func (t *ToolResultComponent) SetToolName(name string) {
	t.toolName = name
}

func (t *ToolResultComponent) SetStatus(status string) {
	t.status = status
}

func (t *ToolResultComponent) SetContent(content string) {
	if len(content) > t.maxLines*200 {
		t.truncated = true
		content = content[:t.maxLines*200]
	}
	t.content = content
}

func (t *ToolResultComponent) SetError(err string) {
	t.error = err
	t.status = ToolStatusError
}

func (t *ToolResultComponent) Render() string {
	var result strings.Builder

	result.WriteString("\n")

	isError := t.status == ToolStatusError || t.error != ""

	if isError {
		result.WriteString(ToolErrorStyle.Render("❌ " + t.toolName + " result"))
	} else {
		result.WriteString(ToolSuccessStyle.Render("✅ " + t.toolName + " result"))
	}

	result.WriteString("\n")

	if isError && t.error != "" {
		result.WriteString(ToolErrorStyle.Render(t.error))
	} else {
		result.WriteString(t.renderContent())
	}

	if t.truncated {
		result.WriteString("\n")
		result.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Render("... (output truncated)"))
	}

	return result.String()
}

func (t *ToolResultComponent) renderContent() string {
	var result strings.Builder

	lines := strings.Split(t.content, "\n")
	lineCount := 0

	for _, line := range lines {
		if lineCount >= t.maxLines {
			break
		}

		if len(line) > t.width-6 {
			line = line[:t.width-9] + "..."
		}

		result.WriteString("  ")
		result.WriteString(ToolSuccessStyle.Render(line))
		result.WriteString("\n")
		lineCount++
	}

	return result.String()
}

type ToolResultView struct {
	toolName   string
	status     string
	result     string
	error      string
	showHeader bool
	width      int
}

func NewToolResultView() *ToolResultView {
	return &ToolResultView{
		width:      80,
		showHeader: true,
		status:    ToolStatusCompleted,
	}
}

func (t *ToolResultView) SetToolName(name string) {
	t.toolName = name
}

func (t *ToolResultView) SetStatus(status string) {
	t.status = status
}

func (t *ToolResultView) SetResult(result string) {
	t.result = result
}

func (t *ToolResultView) SetError(err string) {
	t.error = err
	t.status = ToolStatusError
}

func (t *ToolResultView) SetShowHeader(show bool) {
	t.showHeader = show
}

func (t *ToolResultView) SetWidth(w int) {
	t.width = w
}

func (t *ToolResultView) Render() string {
	var result strings.Builder

	if t.showHeader {
		isError := t.status == ToolStatusError || t.error != ""

		result.WriteString("\n")

		if isError {
			result.WriteString(ToolErrorStyle.Render("┌─ " + t.toolName + " error"))
			result.WriteString("\n")
		} else {
			result.WriteString(ToolSuccessStyle.Render("┌─ " + t.toolName + " result"))
			result.WriteString("\n")
		}
	}

	if t.error != "" {
		result.WriteString("│ ")
		result.WriteString(ToolErrorStyle.Render(t.error))
		result.WriteString("\n")
	} else if t.result != "" {
		lines := strings.Split(t.result, "\n")
		for _, line := range lines {
			if len(line) > t.width-6 {
				line = line[:t.width-9] + "..."
			}
			result.WriteString("│ ")
			result.WriteString(ToolSuccessStyle.Render(line))
			result.WriteString("\n")
		}
	}

	if t.showHeader {
		result.WriteString("└─")
		for i := 0; i < t.width-4 && i < 40; i++ {
			result.WriteString("─")
		}
	}

	return result.String()
}
