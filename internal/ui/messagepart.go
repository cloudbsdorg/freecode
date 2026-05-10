package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	ToolStatusPending   = "pending"
	ToolStatusRunning   = "running"
	ToolStatusCompleted = "completed"
	ToolStatusError     = "error"
)

var (
	TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0"))

	ReasoningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Italic(true)

	ThinkingPrefixStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6A9955")).
		Bold(true)

	ThinkingContentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CDCFE"))

	ToolCallStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#DCDCAA"))

	ToolRunningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4EC9B0"))

	ToolErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F14C4C"))

	ToolSuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4EC9B0"))

	ToolArgsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CDCFE"))

	InlineCodeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CE9178"))

	ContextGroupStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C586C0"))
)

func RenderPart(part MessagePart, width int) string {
	switch part.Type {
	case "text":
		return RenderTextPart(part.Content, width)
	case "reasoning":
		return RenderReasoningPart(part.Content, width)
	case "tool":
		return RenderToolPart(part.Tool, part.Content, width)
	case "tool_call":
		return RenderToolCallPart(part.Tool, part.Content, width)
	case "tool_result":
		return RenderToolResultPart(part.Tool, part.Content, width)
	default:
		return TextStyle.Render(part.Content)
	}
}

func RenderTextPart(content string, width int) string {
	if width <= 0 {
		width = 80
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	var result strings.Builder

	codeBlocks := parseCodeBlocks(content)
	for i, block := range codeBlocks {
		if i > 0 {
			result.WriteString("\n")
		}
		if block.isCode {
			result.WriteString(renderCodeBlock(block.content, block.lang, width))
		} else {
			result.WriteString(renderTextWithInlineCode(block.content, width))
		}
	}

	return result.String()
}

type codeBlock struct {
	isCode bool
	content string
	lang    string
}

func parseCodeBlocks(content string) []codeBlock {
	var blocks []codeBlock
	lines := strings.Split(content, "\n")
	var current strings.Builder
	inCode := false
	codeLang := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if !inCode {
				if current.Len() > 0 {
					blocks = append(blocks, codeBlock{isCode: false, content: current.String(), lang: ""})
					current.Reset()
				}
				inCode = true
				codeLang = strings.TrimPrefix(strings.TrimSpace(line), "```")
				if codeLang == "" {
					codeLang = "plaintext"
				}
			} else {
				blocks = append(blocks, codeBlock{isCode: true, content: current.String(), lang: codeLang})
				current.Reset()
				inCode = false
				codeLang = ""
			}
		} else {
			if current.Len() > 0 {
				current.WriteString("\n")
			}
			current.WriteString(line)
		}
	}

	if current.Len() > 0 {
		blocks = append(blocks, codeBlock{isCode: inCode, content: current.String(), lang: codeLang})
	}

	return blocks
}

func renderCodeBlock(code, lang string, width int) string {
	lines := strings.Split(code, "\n")
	var result strings.Builder

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3A3A3A"))

	result.WriteString(borderStyle.Render("┌─"))
	result.WriteString(InlineCodeStyle.Render(lang))
	result.WriteString(borderStyle.Render("─"))
	for i := 0; i < width-20 && i < len(code)/20+len(lang); i++ {
		result.WriteString(borderStyle.Render("─"))
	}
	result.WriteString(borderStyle.Render("┐"))
	result.WriteString("\n")

	for _, line := range lines {
		result.WriteString(borderStyle.Render("│ "))
		if len(line) > width-4 {
			line = line[:width-7] + " ..."
		}
		result.WriteString(InlineCodeStyle.Render(line))
		result.WriteString("\n")
	}

	result.WriteString(borderStyle.Render("└─"))
	for i := 0; i < width-4; i++ {
		result.WriteString(borderStyle.Render("─"))
	}
	result.WriteString(borderStyle.Render("┘"))

	return result.String()
}

func renderTextWithInlineCode(text string, width int) string {
	var result strings.Builder

	lines := strings.Split(text, "\n")
	for lineIdx, line := range lines {
		if lineIdx > 0 {
			result.WriteString("\n")
		}
		rendered := renderInlineCode(line)
		if len(rendered) > width-4 {
			rendered = wordWrapText(rendered, width-4)
		}
		result.WriteString(TextStyle.Render(rendered))
	}

	return result.String()
}

func wordWrapText(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	var result strings.Builder
	words := strings.Fields(text)
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

	return result.String()
}

func renderInlineCode(line string) string {
	var result strings.Builder
	inCode := false
	var codeContent strings.Builder

	for _, char := range line {
		if char == '`' {
			if inCode {
				result.WriteString(InlineCodeStyle.Render(codeContent.String()))
				codeContent.Reset()
				inCode = false
			} else {
				inCode = true
			}
		} else if inCode {
			codeContent.WriteRune(char)
		} else {
			result.WriteRune(char)
		}
	}

	if codeContent.Len() > 0 && inCode {
		result.WriteString(InlineCodeStyle.Render(codeContent.String()))
	} else if codeContent.Len() > 0 {
		result.WriteString(codeContent.String())
	}

	return result.String()
}

func RenderReasoningPart(content string, width int) string {
	var result strings.Builder

	result.WriteString("\n")
	result.WriteString(ThinkingPrefixStyle.Render("🤔 Thinking"))
	result.WriteString("\n")

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		indent := "  "
		if len(line) > width-6 {
			line = line[:width-9] + "..."
		}
		result.WriteString(indent)
		result.WriteString(ThinkingContentStyle.Render(line))
	}

	result.WriteString("\n")
	return result.String()
}

type MessageToolInfo struct {
	Name      string
	Status    string
	Title     string
	Subtitle  string
	Args      map[string]interface{}
	Result    string
	Error     string
	Collapsed bool
}

func GetToolInfo(toolName string, input map[string]interface{}) MessageToolInfo {
	info := MessageToolInfo{
		Name:   toolName,
		Status: ToolStatusPending,
		Args:   input,
	}

	if desc, ok := input["description"].(string); ok {
		info.Title = desc
	} else if query, ok := input["query"].(string); ok {
		info.Title = query
	} else if url, ok := input["url"].(string); ok {
		info.Title = url
	} else if path, ok := input["filePath"].(string); ok {
		info.Title = path
	} else if name, ok := input["name"].(string); ok {
		info.Title = name
	}

	if pattern, ok := input["pattern"].(string); ok {
		info.Subtitle = pattern
	} else if path, ok := input["path"].(string); ok {
		info.Subtitle = path
	}

	return info
}

func RenderToolPart(toolName string, content string, width int) string {
	var result strings.Builder

	result.WriteString("\n")
	result.WriteString(ToolCallStyle.Render("🔧 " + strings.ToUpper(toolName)))

	if content != "" {
		result.WriteString("\n")
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if i > 0 {
				result.WriteString("\n")
			}
			indent := "  "
			if len(line) > width-6 {
				line = line[:width-9] + "..."
			}
			result.WriteString(indent)
			result.WriteString(ToolArgsStyle.Render(line))
		}
	}

	result.WriteString("\n")
	return result.String()
}

func RenderToolCallPart(toolName string, input string, width int) string {
	var result strings.Builder

	status := ToolStatusPending
	if strings.Contains(input, `"status":"running"`) || strings.Contains(input, `"status": "running"`) {
		status = ToolStatusRunning
	} else if strings.Contains(input, `"status":"completed"`) || strings.Contains(input, `"status": "completed"`) {
		status = ToolStatusCompleted
	} else if strings.Contains(input, `"status":"error"`) || strings.Contains(input, `"status": "error"`) {
		status = ToolStatusError
	}

	icon := "🔧"
	switch status {
	case ToolStatusRunning:
		icon = "⚡"
		result.WriteString(ToolRunningStyle.Render(icon + " " + strings.ToUpper(toolName)))
	case ToolStatusError:
		icon = "❌"
		result.WriteString(ToolErrorStyle.Render(icon + " " + strings.ToUpper(toolName)))
	case ToolStatusCompleted:
		icon = "✅"
		result.WriteString(ToolSuccessStyle.Render(icon + " " + strings.ToUpper(toolName)))
	default:
		result.WriteString(ToolCallStyle.Render(icon + " " + strings.ToUpper(toolName)))
	}

	if input != "" {
		result.WriteString("\n")
		lines := strings.Split(input, "\n")
		for i, line := range lines {
			if i > 0 {
				result.WriteString("\n")
			}
			indent := "  "
			line = highlightJSONFields(line)
			if len(line) > width-6 {
				line = line[:width-9] + "..."
			}
			result.WriteString(indent)
			result.WriteString(ToolArgsStyle.Render(line))
		}
	}

	result.WriteString("\n")
	return result.String()
}

func highlightJSONFields(line string) string {
	line = strings.ReplaceAll(line, `"description"`, `"`+ToolCallStyle.Render("description")+`"`)
	line = strings.ReplaceAll(line, `"query"`, `"`+ToolCallStyle.Render("query")+`"`)
	line = strings.ReplaceAll(line, `"filePath"`, `"`+ToolCallStyle.Render("filePath")+`"`)
	line = strings.ReplaceAll(line, `"path"`, `"`+ToolCallStyle.Render("path")+`"`)
	line = strings.ReplaceAll(line, `"name"`, `"`+ToolCallStyle.Render("name")+`"`)
	return line
}

func RenderToolResultPart(toolName string, result string, width int) string {
	var builder strings.Builder

	builder.WriteString("\n")

	isError := strings.Contains(strings.ToLower(result), "error") ||
		strings.HasPrefix(strings.TrimSpace(result), "!")
	if isError {
		builder.WriteString(ToolErrorStyle.Render("❌ Result:"))
	} else {
		builder.WriteString(ToolSuccessStyle.Render("✅ Result:"))
	}

	builder.WriteString("\n")

	lines := strings.Split(result, "\n")
	for i, line := range lines {
		if i > 0 {
			builder.WriteString("\n")
		}
		indent := "  "
		if len(line) > width-6 {
			line = line[:width-9] + "..."
		}
		if isError {
			builder.WriteString(indent)
			builder.WriteString(ToolErrorStyle.Render(line))
		} else {
			builder.WriteString(indent)
			builder.WriteString(ToolSuccessStyle.Render(line))
		}
	}

	builder.WriteString("\n")
	return builder.String()
}

func RenderToolCall(info MessageToolInfo, width int) string {
	var result strings.Builder

	result.WriteString("\n")
	result.WriteString("┌─ ")
	result.WriteString(ToolCallStyle.Render("Tool: " + info.Name))
	result.WriteString("\n")

	statusIcon := "⏳"
	statusStyle := ToolRunningStyle
	switch info.Status {
	case ToolStatusCompleted:
		statusIcon = "✅"
		statusStyle = ToolSuccessStyle
	case ToolStatusError:
		statusIcon = "❌"
		statusStyle = ToolErrorStyle
	case ToolStatusRunning:
		statusIcon = "⚡"
		statusStyle = ToolRunningStyle
	}

	result.WriteString("│ ")
	result.WriteString(statusStyle.Render(statusIcon + " " + info.Status))
	result.WriteString("\n")

	if info.Title != "" {
		result.WriteString("│ ")
		result.WriteString(ToolCallStyle.Render(info.Title))
		result.WriteString("\n")
	}

	if info.Subtitle != "" {
		result.WriteString("│ ")
		result.WriteString(ToolArgsStyle.Render(info.Subtitle))
		result.WriteString("\n")
	}

	if len(info.Args) > 0 {
		result.WriteString("│ Arguments:\n")
		for key, value := range info.Args {
			result.WriteString("│   ")
			result.WriteString(ToolCallStyle.Render(key + ": "))
			result.WriteString(ToolArgsStyle.Render(fmt.Sprintf("%v", value)))
			result.WriteString("\n")
		}
	}

	if info.Error != "" {
		result.WriteString("│ Error: ")
		result.WriteString(ToolErrorStyle.Render(info.Error))
		result.WriteString("\n")
	}

	result.WriteString("└─")
	for i := 0; i < 40 && i < width-2; i++ {
		result.WriteString("─")
	}

	return result.String()
}

func RenderContextGroup(tools []MessageToolInfo, width int) string {
	if len(tools) == 0 {
		return ""
	}

	var result strings.Builder

	result.WriteString("\n")
	result.WriteString(ContextGroupStyle.Render("📚 Context"))
	result.WriteString("\n")

	readCount := 0
	searchCount := 0
	listCount := 0

	for _, tool := range tools {
		switch strings.ToLower(tool.Name) {
		case "read", "glob", "grep":
			readCount++
		case "search", "websearch":
			searchCount++
		case "list", "ls":
			listCount++
		}
	}

	if readCount > 0 {
		result.WriteString("  Read: ")
		result.WriteString(ContextGroupStyle.Render(fmt.Sprintf("%d", readCount)))
		result.WriteString(" | ")
	}
	if searchCount > 0 {
		result.WriteString("  Search: ")
		result.WriteString(ContextGroupStyle.Render(fmt.Sprintf("%d", searchCount)))
		result.WriteString(" | ")
	}
	if listCount > 0 {
		result.WriteString("  List: ")
		result.WriteString(ContextGroupStyle.Render(fmt.Sprintf("%d", listCount)))
	}

	result.WriteString("\n")

	return result.String()
}

func RenderThinkingBlock(content string, collapsed bool, width int) string {
	var result strings.Builder

	if collapsed {
		result.WriteString(ThinkingPrefixStyle.Render("🤔 ···"))
		result.WriteString(" (click to expand)")
	} else {
		result.WriteString(ThinkingPrefixStyle.Render("🤔 Thinking"))
		result.WriteString("\n")

		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if i > 0 {
				result.WriteString("\n")
			}
			indent := "  "
			if len(line) > width-6 {
				line = line[:width-9] + "..."
			}
			result.WriteString(indent)
			result.WriteString(ThinkingContentStyle.Render(line))
		}
	}

	return result.String()
}
