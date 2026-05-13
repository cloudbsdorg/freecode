package ui

import (
	"strings"

	"github.com/freecode/freecode/internal/style"
)

type PermissionRequest struct {
	ID         string
	SessionID  string
	Permission string
	Metadata   map[string]interface{}
	Tool       *ToolInfo
	Patterns   []string
	Always     []string
}

type ToolInfo struct {
	MessageID string
	CallID    string
	Input     map[string]interface{}
}

type PermissionResponse struct {
	RequestID string
	Reply     string
	Workspace string
	Message   string
}

type PermissionDialogState struct {
	Request   *PermissionRequest
	Stage     string
	Expanded  bool
	RejectMsg string
}

type PermissionDialog struct {
	state PermissionDialogState
	width int
}

func NewPermissionDialog() *PermissionDialog {
	return &PermissionDialog{
		state: PermissionDialogState{
			Stage: "permission",
		},
		width: 60,
	}
}

func (p *PermissionDialog) SetWidth(w int) {
	p.width = w
}

func (p *PermissionDialog) SetRequest(req *PermissionRequest) {
	p.state.Request = req
	p.state.Stage = "permission"
	p.state.Expanded = false
	p.state.RejectMsg = ""
}

func (p *PermissionDialog) GetRequest() *PermissionRequest {
	return p.state.Request
}

func (p *PermissionDialog) SetStage(stage string) {
	p.state.Stage = stage
}

func (p *PermissionDialog) GetStage() string {
	return p.state.Stage
}

func (p *PermissionDialog) SetRejectMsg(msg string) {
	p.state.RejectMsg = msg
}

func (p *PermissionDialog) GetRejectMsg() string {
	return p.state.RejectMsg
}

func (p *PermissionDialog) IsVisible() bool {
	return p.state.Request != nil
}

func (p *PermissionDialog) ToggleExpanded() {
	p.state.Expanded = !p.state.Expanded
}

func (p *PermissionDialog) IsExpanded() bool {
	return p.state.Expanded
}

func (p *PermissionDialog) Clear() {
	p.state.Request = nil
	p.state.Stage = "permission"
	p.state.Expanded = false
	p.state.RejectMsg = ""
}

func (p *PermissionDialog) Render() string {
	if !p.IsVisible() {
		return ""
	}

	req := p.state.Request
	if req == nil {
		return ""
	}

	dialogStyle := style.NewStyle().
		Background(style.Color("#1E1E1E")).
		BorderStyle(style.HiddenBorder()).
		Width(p.width)

	stage := p.state.Stage

	if stage == "always" {
		return dialogStyle.Render(p.renderAlwaysStage())
	}

	if stage == "reject" {
		return dialogStyle.Render(p.renderRejectStage())
	}

	return dialogStyle.Render(p.renderPermissionStage())
}

func (p *PermissionDialog) renderPermissionStage() string {
	var lines []string

	lines = append(lines, p.renderHeader("Permission required"))
	lines = append(lines, "")
	lines = append(lines, p.renderPermissionInfo()...)
	lines = append(lines, "")
	lines = append(lines, p.renderOptions([]string{"Allow once", "Allow always", "Reject"}, 0)...)

	return strings.Join(lines, "\n")
}

func (p *PermissionDialog) renderAlwaysStage() string {
	req := p.state.Request
	var lines []string

	lines = append(lines, p.renderHeader("Always allow"))
	lines = append(lines, "")

	if len(req.Always) == 1 && req.Always[0] == "*" {
		lines = append(lines, p.renderMuted("This will allow "+req.Permission+" until Freecode is restarted."))
	} else {
		lines = append(lines, p.renderMuted("This will allow the following patterns until Freecode is restarted"))
		lines = append(lines, "")
		for _, pattern := range req.Always {
			lines = append(lines, "  - "+pattern)
		}
	}

	lines = append(lines, "")
	lines = append(lines, p.renderOptions([]string{"Confirm", "Cancel"}, 1)...)

	return strings.Join(lines, "\n")
}

func (p *PermissionDialog) renderRejectStage() string {
	var lines []string

	lines = append(lines, p.renderHeader("Reject permission"))
	lines = append(lines, "")
	lines = append(lines, p.renderMuted("Tell Freecode what to do differently"))
	lines = append(lines, "")
	lines = append(lines, p.renderInput(p.state.RejectMsg))
	lines = append(lines, "")
	lines = append(lines, p.renderOptions([]string{"Confirm", "Cancel"}, 1)...)

	return strings.Join(lines, "\n")
}

func (p *PermissionDialog) renderHeader(title string) string {
	warningStyle := style.NewStyle().Foreground(style.Color("#FFCC00"))
	textStyle := style.NewStyle().Foreground(style.Color("#E0E0E0"))
	return warningStyle.Render("△") + " " + textStyle.Render(title)
}

func (p *PermissionDialog) renderPermissionInfo() []string {
	req := p.state.Request
	var lines []string

	icon := p.getPermissionIcon()
	title := p.getPermissionTitle()
	muted := p.renderMuted(title)

	lines = append(lines, "  "+icon+" "+muted)

	switch req.Permission {
	case "edit":
		filepath := p.getMetadataString("filepath")
		if filepath != "" {
			lines = append(lines, p.renderMuted("  Path: "+filepath))
		}
		diff := p.getMetadataString("diff")
		if diff != "" {
			lines = append(lines, "")
			lines = append(lines, p.renderMuted("  Diff:"))
			lines = append(lines, p.renderDiff(diff))
		}
	case "shell":
		command := p.getToolInputString("command")
		if command != "" {
			lines = append(lines, "  $ "+command)
		}
		description := p.getToolInputString("description")
		if description != "" {
			lines = append(lines, p.renderMuted("  "+description))
		}
	case "read":
		filePath := p.getToolInputString("filePath")
		if filePath != "" {
			lines = append(lines, p.renderMuted("  Path: "+filePath))
		}
	case "glob":
		pattern := p.getToolInputString("pattern")
		if pattern != "" {
			lines = append(lines, p.renderMuted("  Pattern: "+pattern))
		}
	case "grep":
		pattern := p.getToolInputString("pattern")
		if pattern != "" {
			lines = append(lines, p.renderMuted("  Pattern: "+pattern))
		}
	case "list":
		path := p.getToolInputString("path")
		if path != "" {
			lines = append(lines, p.renderMuted("  Path: "+path))
		}
	case "task":
		subagentType := p.getToolInputString("subagent_type")
		description := p.getToolInputString("description")
		lines = append(lines, p.renderMuted("  "+subagentType+" Task"))
		if description != "" {
			lines = append(lines, "  ◎ "+description)
		}
	case "webfetch":
		url := p.getToolInputString("url")
		if url != "" {
			lines = append(lines, p.renderMuted("  URL: "+url))
		}
	case "websearch":
		query := p.getToolInputString("query")
		if query != "" {
			lines = append(lines, p.renderMuted("  Query: "+query))
		}
	}

	return lines
}

func (p *PermissionDialog) renderOptions(options []string, selected int) []string {
	var lines []string
	var optionStr string
	for i, opt := range options {
		if i > 0 {
			optionStr += "  "
		}
		if i == selected {
			optionStr += "[" + opt + "]"
		} else {
			optionStr += opt
		}
	}
	lines = append(lines, optionStr)
	lines = append(lines, p.renderHint("←→ select  enter confirm  esc cancel"))
	return lines
}

func (p *PermissionDialog) renderInput(value string) string {
	if value == "" {
		value = ""
	}
	return style.NewStyle().
		Background(style.Color("#3C3C3C")).
		Foreground(style.Color("#E0E0E0")).
		Padding(0, 1).
		Render(value + "_")
}

func (p *PermissionDialog) renderDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimPrefix(line, "+")
		isAdded := strings.HasPrefix(line, "+")
		isRemoved := strings.HasPrefix(line, "-")

		if isAdded {
			result = append(result, style.NewStyle().
				Foreground(style.Color("#4EC9B0")).
				Render("+ "+trimmed))
		} else if isRemoved {
			result = append(result, style.NewStyle().
				Foreground(style.Color("#F44747")).
				Render("- "+trimmed))
		} else {
			result = append(result, "  "+line)
		}
	}
	return strings.Join(result, "\n")
}

func (p *PermissionDialog) renderMuted(text string) string {
	return style.NewStyle().
		Foreground(style.Color("#808080")).
		Render(text)
}

func (p *PermissionDialog) renderHint(text string) string {
	return style.NewStyle().
		Foreground(style.Color("#808080")).
		Render(text)
}

func (p *PermissionDialog) getPermissionIcon() string {
	req := p.state.Request
	switch req.Permission {
	case "edit":
		return "→"
	case "read":
		return "→"
	case "shell":
		return "#"
	case "glob":
		return "✱"
	case "grep":
		return "✱"
	case "list":
		return "→"
	case "task":
		return "#"
	case "webfetch":
		return "%"
	case "websearch":
		return "◈"
	case "external_directory":
		return "←"
	default:
		return "⚙"
	}
}

func (p *PermissionDialog) getPermissionTitle() string {
	req := p.state.Request
	switch req.Permission {
	case "edit":
		filepath := p.getMetadataString("filepath")
		return "Edit " + filepath
	case "read":
		filePath := p.getToolInputString("filePath")
		return "Read " + filePath
	case "shell":
		description := p.getToolInputString("description")
		if description == "" {
			return "Shell command"
		}
		return description
	case "glob":
		pattern := p.getToolInputString("pattern")
		return "Glob \"" + pattern + "\""
	case "grep":
		pattern := p.getToolInputString("pattern")
		return "Grep \"" + pattern + "\""
	case "list":
		path := p.getToolInputString("path")
		return "List " + path
	case "task":
		subagentType := p.getToolInputString("subagent_type")
		return subagentType + " Task"
	case "webfetch":
		url := p.getToolInputString("url")
		return "WebFetch " + url
	case "websearch":
		query := p.getToolInputString("query")
		return "Exa Web Search \"" + query + "\""
	case "external_directory":
		filepath := p.getMetadataString("filepath")
		if filepath != "" {
			return "Access external directory " + filepath
		}
		parentDir := p.getMetadataString("parentDir")
		if parentDir != "" {
			return "Access external directory " + parentDir
		}
		return "Access external directory"
	default:
		return "Call tool " + req.Permission
	}
}

func (p *PermissionDialog) getMetadataString(key string) string {
	req := p.state.Request
	if req.Metadata == nil {
		return ""
	}
	if v, ok := req.Metadata[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (p *PermissionDialog) getToolInputString(key string) string {
	req := p.state.Request
	if req.Tool == nil || req.Tool.Input == nil {
		return ""
	}
	if v, ok := req.Tool.Input[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}