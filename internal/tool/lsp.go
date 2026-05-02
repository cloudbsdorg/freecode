package tool

import (
	"context"
	"fmt"
)

type LSPTool struct{}

func NewLSPTool() *LSPTool {
	return &LSPTool{}
}

func (t *LSPTool) Name() string {
	return "lsp"
}

func (t *LSPTool) Description() string {
	return "Language Server Protocol operations"
}

func (t *LSPTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "lsp",
		Description: "Language Server Protocol operations",
		Parameters: map[string]Parameter{
			"action": {
				Type:        "string",
				Description: "Action: start, stop, goto, hover, completions",
				Required:    true,
				Enum:       []string{"start", "stop", "goto", "hover", "completions"},
			},
			"language": {
				Type:        "string",
				Description: "Programming language",
			},
			"file": {
				Type:        "string",
				Description: "File path",
			},
			"line": {
				Type:        "integer",
				Description: "Line number",
			},
			"character": {
				Type:        "integer",
				Description: "Character position",
			},
		},
	}
}

func (t *LSPTool) Execute(ctx context.Context, req Request) (*Response, error) {
	action, ok := req.Arguments["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action must be a string")
	}

	switch action {
	case "start":
		lang, _ := req.Arguments["language"].(string)
		return &Response{Result: fmt.Sprintf("Starting LSP server for: %s", lang)}, nil
	case "stop":
		return &Response{Result: "Stopping LSP server"}, nil
	case "goto":
		line, _ := req.Arguments["line"].(int)
		char, _ := req.Arguments["character"].(int)
		return &Response{Result: fmt.Sprintf("Go to: line=%d char=%d", line, char)}, nil
	case "hover":
		line, _ := req.Arguments["line"].(int)
		char, _ := req.Arguments["character"].(int)
		return &Response{Result: fmt.Sprintf("Hover info at: line=%d char=%d", line, char)}, nil
	case "completions":
		return &Response{Result: "Completions: (placeholder)"}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
