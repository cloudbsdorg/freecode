package lsp

import (
	"context"
	"os/exec"
	"path/filepath"
	"sync"
)

type ServerConfig struct {
	Command string
	Args    []string
}

var DefaultServers = map[string]ServerConfig{
	"go":         {Command: "gopls", Args: []string{"--stdio"}},
	"typescript": {Command: "typescript-language-server", Args: []string{"--stdio"}},
	"javascript": {Command: "typescript-language-server", Args: []string{"--stdio"}},
	"python":     {Command: "pyright-langserver", Args: []string{"--stdio"}},
	"rust":       {Command: "rust-analyzer", Args: []string{}},
	"csharp":     {Command: "csharp-ls", Args: []string{}},
	"c":          {Command: "clangd", Args: []string{}},
	"cpp":        {Command: "clangd", Args: []string{}},
	"java":       {Command: "jdtls", Args: []string{}},
	"ruby":       {Command: "ruby-lsp", Args: []string{}},
	"php":        {Command: "intelephense", Args: []string{"--stdio"}},
	"swift":      {Command: "sourcekit-lsp", Args: []string{}},
	"kotlin":     {Command: "kotlin-language-server", Args: []string{}},
	"lua":        {Command: "lua-language-server", Args: []string{}},
	"svelte":     {Command: "svelte-language-server", Args: []string{"--stdio"}},
	"vue":        {Command: "vue-language-server", Args: []string{"--stdio"}},
	"html":       {Command: "vscode-html-language-server", Args: []string{"--stdio"}},
	"css":        {Command: "vscode-css-language-server", Args: []string{"--stdio"}},
	"json":       {Command: "vscode-json-language-server", Args: []string{"--stdio"}},
	"yaml":       {Command: "yaml-language-server", Args: []string{"--stdio"}},
	"xml":        {Command: "lemminx", Args: []string{}},
	"markdown":    {Command: "markdown-language-server", Args: []string{"--stdio"}},
	"dockerfile": {Command: "docker-langserver", Args: []string{"--stdio"}},
	"shellscript": {Command: "bash-language-server", Args: []string{"--stdio"}},
}

type ServerManager struct {
	mu      sync.RWMutex
	servers map[string]*Client
	root    string
}

func NewServerManager(root string) *ServerManager {
	return &ServerManager{
		servers: make(map[string]*Client),
		root:    root,
	}
}

func (sm *ServerManager) GetServer(language string) (*Client, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	client, ok := sm.servers[language]
	return client, ok
}

func (sm *ServerManager) GetOrCreateServer(ctx context.Context, language string) (*Client, error) {
	sm.mu.RLock()
	client, exists := sm.servers[language]
	sm.mu.RUnlock()

	if exists {
		return client, nil
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if client, exists = sm.servers[language]; exists {
		return client, nil
	}

	config, ok := DefaultServers[language]
	if !ok {
		return nil, &ServerNotFoundError{Language: language}
	}

	if !isCommandAvailable(config.Command) {
		return nil, &ServerNotFoundError{Language: language}
	}

	client = NewClient()
	cmd := config.Command

	if err := client.Connect(ctx, cmd, sm.root); err != nil {
		return nil, err
	}

	if err := client.Initialize(ctx); err != nil {
		client.Close()
		return nil, err
	}

	sm.servers[language] = client
	return client, nil
}

func (sm *ServerManager) CloseServer(language string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	client, ok := sm.servers[language]
	if !ok {
		return nil
	}

	err := client.Close()
	delete(sm.servers, language)
	return err
}

func (sm *ServerManager) CloseAll() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var lastErr error
	for lang, client := range sm.servers {
		if err := client.Close(); err != nil {
			lastErr = err
		}
		delete(sm.servers, lang)
	}
	return lastErr
}

type ServerNotFoundError struct {
	Language string
}

func (e *ServerNotFoundError) Error() string {
	return "LSP server not found for language: " + e.Language
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func DetectLanguage(filePath string) string {
	ext := filepath.Ext(filePath)
	lang := LanguageByExtension(ext)
	return lang
}
