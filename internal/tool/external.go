package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"
)

type ExternalToolPlugin interface {
	Name() string
	Description() string
	Schema() ToolSchema
	Execute(ctx interface{}, req Request) (*Response, error)
}

type externalLoader struct {
	mu      sync.RWMutex
	loaded  map[string]*plugin.Plugin
	factories map[string]func() Tool
}

var extLoader = &externalLoader{
	loaded:   make(map[string]*plugin.Plugin),
	factories: make(map[string]func() Tool),
}

func RegisterExternal(name string, factory func() Tool) {
	mu.Lock()
	defer mu.Unlock()
	extLoader.factories[name] = factory
}

func LoadExternalTools(toolsDir string) error {
	if toolsDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		toolsDir = filepath.Join(homeDir, ".config", "freecode", "tools")
	}

	if _, err := os.Stat(toolsDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(toolsDir)
	if err != nil {
		return fmt.Errorf("failed to read tools directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".so" {
			continue
		}

		pluginPath := filepath.Join(toolsDir, entry.Name())
		if err := loadPlugin(pluginPath); err != nil {
			fmt.Printf("warning: failed to load plugin %s: %v\n", pluginPath, err)
			continue
		}
	}

	return nil
}

func loadPlugin(path string) error {
	extLoader.mu.Lock()
	defer extLoader.mu.Unlock()

	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	symbol, err := p.Lookup("ToolPlugin")
	if err != nil {
		return fmt.Errorf("plugin missing ToolPlugin symbol: %w", err)
	}

	toolPlugin, ok := symbol.(ExternalToolPlugin)
	if !ok {
		return fmt.Errorf("ToolPlugin symbol does not implement ExternalToolPlugin interface")
	}

	name := toolPlugin.Name()
	extLoader.loaded[path] = p
	extLoader.factories[name] = func() Tool {
		return &pluginWrapper{
			name:        toolPlugin.Name(),
			description: toolPlugin.Description(),
			schema:      toolPlugin.Schema(),
			plugin:      toolPlugin,
		}
	}

	mu.Lock()
	tools[name] = extLoader.factories[name]
	enabled[name] = true
	mu.Unlock()

	return nil
}

type pluginWrapper struct {
	name        string
	description string
	schema      ToolSchema
	plugin      ExternalToolPlugin
}

func (w *pluginWrapper) Name() string {
	return w.name
}

func (w *pluginWrapper) Description() string {
	return w.description
}

func (w *pluginWrapper) Schema() ToolSchema {
	return w.schema
}

func (w *pluginWrapper) Execute(ctx context.Context, req Request) (*Response, error) {
	return w.plugin.Execute(ctx, req)
}