package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

type Client struct {
	conn   *jsonrpc2.Conn
	proc   *exec.Cmd
	root   string
	mu     sync.RWMutex
	files  map[string]*textDocument
}

type textDocument struct {
	uri     string
	version int
	text    string
}

type ServerCapabilities struct {
	TextDocumentSync         any       `json:"textDocumentSync,omitempty"`
	HoverProvider            bool      `json:"hoverProvider,omitempty"`
	DefinitionProvider       bool      `json:"definitionProvider,omitempty"`
	ReferencesProvider       bool      `json:"referencesProvider,omitempty"`
	ImplementationProvider   bool      `json:"implementationProvider,omitempty"`
	DocumentSymbolProvider   bool      `json:"documentSymbolProvider,omitempty"`
	WorkspaceSymbolProvider  bool      `json:"workspaceSymbolProvider,omitempty"`
	CompletionProvider       *struct{} `json:"completionProvider,omitempty"`
	DiagnosticProvider       any       `json:"diagnosticProvider,omitempty"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
	Source   string `json:"source,omitempty"`
	Code     any    `json:"code,omitempty"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      uint32 `json:"line"`
	Character uint32 `json:"character"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	InsertText    string `json:"insertText,omitempty"`
	Kind          int    `json:"kind,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

type InitializeParams struct {
	RootURI        string `json:"rootUri,omitempty"`
	ProcessID      any    `json:"processId,omitempty"`
	WorkspaceRoots []struct {
		Name string `json:"name"`
		URI  string `json:"uri"`
	} `json:"workspaceFolders,omitempty"`
	Capabilities any `json:"capabilities"`
}

func filePathToURI(path string) string {
	return "file://" + filepath.ToSlash(path)
}

func uriToFilePath(uri string) string {
	if len(uri) > 8 && uri[:7] == "file://" {
		return uri[7:]
	}
	return uri
}

func NewClient() *Client {
	return &Client{files: make(map[string]*textDocument)}
}

func (c *Client) Connect(ctx context.Context, server string, root string) error {
	cmd := exec.CommandContext(ctx, server)
	cmd.Dir = root

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	c.proc = cmd
	c.root = root

	conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(stdin), jsonrpc2.NewPlainObjectStream(stdout), c.handler())
	c.conn = conn

	return nil
}

func (c *Client) handler() jsonrpc2.Handler {
	return jsonrpc2.HandlerFunc(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		switch req.Method {
		case "window/workDoneProgress/create":
			return nil, nil
		case "workspace/configuration":
			return []any{}, nil
		case "client/registerCapability":
			return nil, nil
		case "client/unregisterCapability":
			return nil, nil
		case "workspace/workspaceFolders":
			return []any{{Name: "workspace", URI: filePathToURI(c.root)}}, nil
		case "workspace/diagnostic/refresh":
			return nil, nil
		case "textDocument/publishDiagnostics":
			// handled via callback if registered
			return nil, nil
		default:
			return nil, &jsonrpc2.Error{Code: jsonrpc2.MethodNotFound, Message: fmt.Sprintf("method %q not handled", req.Method)}
		}
	})
}

func (c *Client) Initialize(ctx context.Context) error {
	params := InitializeParams{
		RootURI:    filePathToURI(c.root),
		ProcessID:  os.Getpid(),
		Capabilities: map[string]any{
			"window": map[string]any{
				"workDoneProgress": true,
			},
			"workspace": map[string]any{
				"configuration":     true,
				"didChangeWatchedFiles": map[string]any{"dynamicRegistration": true},
			},
			"textDocument": map[string]any{
				"synchronization": map[string]any{
					"didOpen": true,
					"didChange": true,
				},
				"diagnostic": map[string]any{
					"dynamicRegistration": true,
				},
				"publishDiagnostics": map[string]any{
					"versionSupport": false,
				},
			},
		},
	}

	var result InitializeResult
	err := c.conn.Call(ctx, "initialize", params, &result)
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	err = c.conn.Notify(ctx, "initialized", map[string]any{})
	if err != nil {
		return fmt.Errorf("initialized notification: %w", err)
	}

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	var result any
	return c.conn.Call(ctx, "shutdown", nil, &result)
}

func (c *Client) DidOpen(ctx context.Context, path, text string) error {
	uri := filePathToURI(path)
	c.mu.Lock()
	c.files[path] = &textDocument{uri: uri, version: 0, text: text}
	c.mu.Unlock()

	return c.conn.Notify(ctx, "textDocument/didOpen", map[string]any{}{
		"textDocument": map[string]any{
			"uri":      uri,
			"languageId": languageID(filepath.Ext(path)),
			"version":  0,
			"text":     text,
		},
	})
}

func (c *Client) DidChange(ctx context.Context, path, text string) error {
	c.mu.Lock()
	doc := c.files[path]
	if doc == nil {
		doc = &textDocument{uri: filePathToURI(path), version: 0, text: text}
		c.files[path] = doc
	}
	doc.version++
	doc.text = text
	c.mu.Unlock()

	return c.conn.Notify(ctx, "textDocument/didChange", map[string]any{}{
		"textDocument": map[string]any{
			"uri":    doc.uri,
			"version": doc.version,
		},
		"contentChanges": []map[string]any{{"text": text}},
	})
}

func (c *Client) Hover(ctx context.Context, path string, line, character uint32) (*HoverResult, error) {
	c.mu.RLock()
	doc := c.files[path]
	c.mu.RUnlock()
	if doc == nil {
		return nil, nil
	}

	var result HoverResult
	err := c.conn.Call(ctx, "textDocument/hover", map[string]any{}{
		"textDocument": map[string]any{"uri": doc.uri},
		"position":     map[string]any{"line": line, "character": character},
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type HoverResult struct {
	Contents any `json:"contents"`
	Range    *Range `json:"range,omitempty"`
}

func (c *Client) Definition(ctx context.Context, path string, line, character uint32) ([]Location, error) {
	c.mu.RLock()
	doc := c.files[path]
	c.mu.RUnlock()
	if doc == nil {
		return nil, nil
	}

	var result []Location
	err := c.conn.Call(ctx, "textDocument/definition", map[string]any{}{
		"textDocument": map[string]any{"uri": doc.uri},
		"position":     map[string]any{"line": line, "character": character},
	}, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Location struct {
	URI  string `json:"uri"`
	Range Range `json:"range"`
}

func (c *Client) References(ctx context.Context, path string, line, character uint32) ([]Location, error) {
	c.mu.RLock()
	doc := c.files[path]
	c.mu.RUnlock()
	if doc == nil {
		return nil, nil
	}

	var result []Location
	err := c.conn.Call(ctx, "textDocument/references", map[string]any{}{
		"textDocument": map[string]any{"uri": doc.uri},
		"position":     map[string]any{"line": line, "character": character},
		"context":      map[string]any{"includeDeclaration": true},
	}, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Completion(ctx context.Context, path string, line, character uint32) ([]CompletionItem, error) {
	c.mu.RLock()
	doc := c.files[path]
	c.mu.RUnlock()
	if doc == nil {
		return nil, nil
	}

	var result CompletionList
	err := c.conn.Call(ctx, "textDocument/completion", map[string]any{}{
		"textDocument": map[string]any{"uri": doc.uri},
		"position":     map[string]any{"line": line, "character": character},
	}, &result)
	if err != nil {
		return nil, err
	}
	if result.IsIncomplete {
		return nil, nil
	}
	return result.Items, nil
}

type CompletionList struct {
	Items      []CompletionItem `json:"items"`
	IsIncomplete bool            `json:"isIncomplete"`
}

func (c *Client) DocumentSymbol(ctx context.Context, path string) ([]DocumentSymbol, error) {
	c.mu.RLock()
	doc := c.files[path]
	c.mu.RUnlock()
	if doc == nil {
		return nil, nil
	}

	var result []DocumentSymbol
	err := c.conn.Call(ctx, "textDocument/documentSymbol", map[string]any{}{
		"textDocument": map[string]any{"uri": doc.uri},
	}, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type DocumentSymbol struct {
	Name        string  `json:"name"`
	Detail      string  `json:"detail,omitempty"`
	Kind        int     `json:"kind"`
	Range       Range   `json:"range"`
	SelectionRange Range `json:"selectionRange"`
	Children    []DocumentSymbol `json:"children,omitempty"`
}

func (c *Client) WorkspaceSymbol(ctx context.Context, query string) ([]SymbolInformation, error) {
	var result []SymbolInformation
	err := c.conn.Call(ctx, "workspace/symbol", map[string]any{
		"query": query,
	}, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type SymbolInformation struct {
	Name       string   `json:"name"`
	Kind       int      `json:"kind"`
	Location   Location `json:"location"`
	ContainerName string `json:"containerName,omitempty"`
}

func languageID(ext string) string {
	switch ext {
	case ".go":
		return "go"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".py":
		return "python"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".cs":
		return "csharp"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".swift":
		return "swift"
	case ".kt", ".kts":
		return "kotlin"
	case ".scala":
		return "scala"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".scss", ".sass":
		return "scss"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".xml":
		return "xml"
	case ".md":
		return "markdown"
	case ".sql":
		return "sql"
	case ".sh", ".bash":
		return "shellscript"
	case ".dockerfile":
		return "dockerfile"
	default:
		return "plaintext"
	}
}

func (c *Client) Close() error {
	if c.conn != nil {
		c.conn.Close()
	}
	if c.proc != nil && c.proc.Process != nil {
		c.proc.Process.Kill()
	}
	return nil
}

// --- Server side ---

type Server struct {
	conn    *jsonrpc2.Conn
	stdin   chan json.RawMessage
	stdout  chan json.RawMessage
}

func NewServer() *Server {
	return &Server{
		stdin:  make(chan json.RawMessage, 100),
		stdout: make(chan json.RawMessage, 100),
	}
}

func (s *Server) Start(ctx context.Context) error {
	conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(s.stdout), jsonrpc2.NewPlainObjectStream(s.stdin), s.handler())
	s.conn = conn
	return nil
}

func (s *Server) handler() jsonrpc2.Handler {
	return jsonrpc2.HandlerFunc(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		switch req.Method {
		case "initialize":
			var params InitializeParams
			if err := json.Unmarshal(*req.Params, &params); err != nil {
				return nil, err
			}
			return InitializeResult{
				Capabilities: ServerCapabilities{
					TextDocumentSync:        2,
					HoverProvider:           true,
					DefinitionProvider:      true,
					ReferencesProvider:     true,
					DocumentSymbolProvider: true,
					WorkspaceSymbolProvider: true,
					CompletionProvider:     &struct{}{},
				},
			}, nil
		case "shutdown":
			return nil, nil
		case "textDocument/didOpen":
			return nil, nil
		case "textDocument/didChange":
			return nil, nil
		case "textDocument/hover":
			return nil, nil
		case "textDocument/definition":
			return nil, nil
		case "textDocument/references":
			return nil, nil
		case "textDocument/completion":
			return CompletionList{Items: []CompletionItem{}}, nil
		case "textDocument/documentSymbol":
			return []DocumentSymbol{}, nil
		case "workspace/symbol":
			return []SymbolInformation{}, nil
		default:
			return nil, &jsonrpc2.Error{Code: jsonrpc2.MethodNotFound, Message: fmt.Sprintf("method %q not handled", req.Method)}
		}
	})
}

func (s *Server) Stop(ctx context.Context) error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *Server) Send(ctx context.Context, msg json.RawMessage) error {
	s.stdin <- msg
	return nil
}

func (s *Server) Receive(ctx context.Context) (json.RawMessage, error) {
	select {
	case msg := <-s.stdout:
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
