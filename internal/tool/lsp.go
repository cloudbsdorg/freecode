package tool

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/freecode/freecode/internal/lsp"
)

type LSPTool struct {
	manager *lsp.ServerManager
}

func init() {
	Register("lsp", func() Tool { return NewLSPTool() })
}

func NewLSPTool() *LSPTool {
	return &LSPTool{
		manager: lsp.NewServerManager(""),
	}
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
		Description: "LSP operations for hover, goto definition, find references, and diagnostics",
		Parameters: map[string]Parameter{
			"action": {
				Type:        "string",
				Description: "Action: hover, definition, references, completion, diagnostics, start",
				Required:    true,
				Enum:        []string{"hover", "definition", "references", "completion", "diagnostics", "start"},
			},
			"file": {
				Type:        "string",
				Description: "File path",
				Required:    true,
			},
			"line": {
				Type:        "integer",
				Description: "Line number (1-based)",
			},
			"character": {
				Type:        "integer",
				Description: "Character position (1-based)",
			},
			"language": {
				Type:        "string",
				Description: "Language ID (e.g., go, python, typescript)",
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
		return t.executeStart(ctx, req.Arguments)
	case "stop":
		return t.executeStop(ctx, req.Arguments)
	case "hover", "definition", "references", "completion", "diagnostics":
		file, ok := req.Arguments["file"].(string)
		if !ok {
			return nil, fmt.Errorf("file must be a string")
		}
		return t.executeFileAction(ctx, action, file, req.Arguments)
	case "goto":
		file, ok := req.Arguments["file"].(string)
		if !ok {
			return nil, fmt.Errorf("file must be a string")
		}
		return t.executeDefinition(ctx, file, req.Arguments)
	case "completions":
		file, ok := req.Arguments["file"].(string)
		if !ok {
			return nil, fmt.Errorf("file must be a string")
		}
		return t.executeCompletion(ctx, file, req.Arguments)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (t *LSPTool) executeStart(ctx context.Context, args map[string]interface{}) (*Response, error) {
	lang, _ := args["language"].(string)
	if lang == "" {
		return nil, fmt.Errorf("language is required for start action")
	}

	dir := ""
	if root, ok := args["root"].(string); ok {
		dir = root
	}
	if dir == "" {
		dir = "."
	}

	t.manager = lsp.NewServerManager(dir)

	_, err := t.manager.GetOrCreateServer(ctx, lang)
	if err != nil {
		return &Response{Result: fmt.Sprintf("Failed to start LSP server for %s: %v", lang, err)}, nil
	}

	return &Response{Result: fmt.Sprintf("LSP server started for %s (root: %s)", lang, dir)}, nil
}

func (t *LSPTool) executeStop(ctx context.Context, args map[string]interface{}) (*Response, error) {
	if t.manager == nil {
		return &Response{Result: "No LSP server running"}, nil
	}

	if err := t.manager.CloseAll(); err != nil {
		return &Response{Result: fmt.Sprintf("Error stopping LSP server: %v", err)}, nil
	}

	return &Response{Result: "LSP server stopped"}, nil
}

func (t *LSPTool) executeFileAction(ctx context.Context, action string, file string, args map[string]interface{}) (*Response, error) {
	switch action {
	case "hover":
		return t.executeHover(ctx, file, args)
	case "definition":
		return t.executeDefinition(ctx, file, args)
	case "references":
		return t.executeReferences(ctx, file, args)
	case "completion":
		return t.executeCompletion(ctx, file, args)
	case "diagnostics":
		return t.executeDiagnostics(ctx, file, args)
	default:
		return nil, fmt.Errorf("unknown file action: %s", action)
	}
}

func (t *LSPTool) executeHover(ctx context.Context, file string, args map[string]interface{}) (*Response, error) {
	if t.manager == nil {
		return nil, fmt.Errorf("no LSP server started. Use 'start' action first")
	}

	line, char := getLineChar(args)
	lang := getLanguage(file, args)

	client, ok := t.manager.GetServer(lang)
	if !ok {
		return &Response{Result: "No LSP server for language: " + lang}, nil
	}

	text, err := readFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err := client.DidOpen(ctx, file, text); err != nil {
		return nil, fmt.Errorf("didOpen: %w", err)
	}

	result, err := client.Hover(ctx, file, uint32(line-1), uint32(char-1))
	if err != nil {
		return nil, fmt.Errorf("hover: %w", err)
	}

	if result == nil {
		return &Response{Result: "No hover information available"}, nil
	}

	return &Response{Result: fmt.Sprintf("%v", result.Contents)}, nil
}

func (t *LSPTool) executeDefinition(ctx context.Context, file string, args map[string]interface{}) (*Response, error) {
	if t.manager == nil {
		return nil, fmt.Errorf("no LSP server started. Use 'start' action first")
	}

	line, char := getLineChar(args)
	lang := getLanguage(file, args)

	client, ok := t.manager.GetServer(lang)
	if !ok {
		return &Response{Result: "No LSP server for language: " + lang}, nil
	}

	text, err := readFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err := client.DidOpen(ctx, file, text); err != nil {
		return nil, fmt.Errorf("didOpen: %w", err)
	}

	locations, err := client.Definition(ctx, file, uint32(line-1), uint32(char-1))
	if err != nil {
		return nil, fmt.Errorf("definition: %w", err)
	}

	if len(locations) == 0 {
		return &Response{Result: "No definition found"}, nil
	}

	var results []string
	for _, loc := range locations {
		results = append(results, fmt.Sprintf("%s:%d:%d", lsp.UriToFilePath(loc.URI), loc.Range.Start.Line+1, loc.Range.Start.Character+1))
	}

	return &Response{Result: strings.Join(results, "\n")}, nil
}

func (t *LSPTool) executeReferences(ctx context.Context, file string, args map[string]interface{}) (*Response, error) {
	if t.manager == nil {
		return nil, fmt.Errorf("no LSP server started. Use 'start' action first")
	}

	line, char := getLineChar(args)
	lang := getLanguage(file, args)

	client, ok := t.manager.GetServer(lang)
	if !ok {
		return &Response{Result: "No LSP server for language: " + lang}, nil
	}

	text, err := readFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err := client.DidOpen(ctx, file, text); err != nil {
		return nil, fmt.Errorf("didOpen: %w", err)
	}

	locations, err := client.References(ctx, file, uint32(line-1), uint32(char-1))
	if err != nil {
		return nil, fmt.Errorf("references: %w", err)
	}

	if len(locations) == 0 {
		return &Response{Result: "No references found"}, nil
	}

	var results []string
	for _, loc := range locations {
		results = append(results, fmt.Sprintf("%s:%d:%d", lsp.UriToFilePath(loc.URI), loc.Range.Start.Line+1, loc.Range.Start.Character+1))
	}

	return &Response{Result: strings.Join(results, "\n")}, nil
}

func (t *LSPTool) executeCompletion(ctx context.Context, file string, args map[string]interface{}) (*Response, error) {
	if t.manager == nil {
		return nil, fmt.Errorf("no LSP server started. Use 'start' action first")
	}

	line, char := getLineChar(args)
	lang := getLanguage(file, args)

	client, ok := t.manager.GetServer(lang)
	if !ok {
		return &Response{Result: "No LSP server for language: " + lang}, nil
	}

	text, err := readFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err := client.DidOpen(ctx, file, text); err != nil {
		return nil, fmt.Errorf("didOpen: %w", err)
	}

	items, err := client.Completion(ctx, file, uint32(line-1), uint32(char-1))
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}

	if len(items) == 0 {
		return &Response{Result: "No completions available"}, nil
	}

	var results []string
	for _, item := range items {
		results = append(results, item.Label)
	}

	return &Response{Result: strings.Join(results, "\n")}, nil
}

func (t *LSPTool) executeDiagnostics(ctx context.Context, file string, args map[string]interface{}) (*Response, error) {
	if t.manager == nil {
		return nil, fmt.Errorf("no LSP server started. Use 'start' action first")
	}

	lang := getLanguage(file, args)

	client, ok := t.manager.GetServer(lang)
	if !ok {
		return &Response{Result: "No LSP server for language: " + lang}, nil
	}

	text, err := readFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err := client.DidOpen(ctx, file, text); err != nil {
		return nil, fmt.Errorf("didOpen: %w", err)
	}

	uri := lsp.FilePathToURI(file)
	diagnostics := client.GetDiagnostics(uri)

	if len(diagnostics) == 0 {
		return &Response{Result: "No diagnostics"}, nil
	}

	var results []string
	for _, d := range diagnostics {
		results = append(results, lsp.PrettyDiagnostic(d))
	}

	return &Response{Result: strings.Join(results, "\n")}, nil
}

func getLineChar(args map[string]interface{}) (line int, char int) {
	line, _ = args["line"].(int)
	if line == 0 {
		line = 1
	}
	char, _ = args["character"].(int)
	if char == 0 {
		char = 1
	}
	return line, char
}

func getLanguage(file string, args map[string]interface{}) string {
	if lang, ok := args["language"].(string); ok && lang != "" {
		return lang
	}
	return lsp.DetectLanguage(file)
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
