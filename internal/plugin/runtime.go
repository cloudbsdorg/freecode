package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	ErrPluginNotLoaded   = errors.New("plugin not loaded")
	ErrPluginLoadFailed  = errors.New("plugin load failed")
	ErrPluginUnloadFailed = errors.New("plugin unload failed")
	ErrManifestNotFound  = errors.New("plugin manifest not found")
	ErrInvalidManifest   = errors.New("invalid plugin manifest")
	ErrPluginInitFailed  = errors.New("plugin initialization failed")
)

// Request represents a plugin execution request.
type Request struct {
	Name      string
	Arguments map[string]interface{}
	SessionID string
}

// Response represents a plugin execution response.
type Response struct {
	Result any
	Error  error
}

// ExecutablePlugin is the interface that plugin implementations must satisfy.
// This extends the basic Plugin interface with execution capabilities.
type ExecutablePlugin interface {
	Plugin
	Execute(ctx context.Context, req Request) (*Response, error)
	Shutdown() error
}

// PluginInfo contains metadata about a loaded plugin.
type PluginInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Author    string            `json:"author"`
	Description string          `json:"description"`
	Path      string            `json:"path"`
	Hooks     []string          `json:"hooks"`
	APIKeys   map[string]string `json:"api_keys,omitempty"`
	LoadedAt  time.Time         `json:"loaded_at"`
	State     PluginState       `json:"state"`
}

type PluginState string

const (
	PluginStateLoading    PluginState = "loading"
	PluginStateLoaded     PluginState = "loaded"
	PluginStateError      PluginState = "error"
	PluginStateUnloading  PluginState = "unloading"
	PluginStateUnloaded   PluginState = "unloaded"
)

// PluginManifest represents the plugin.json manifest file.
type PluginManifest struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Description string            `json:"description"`
	Main        string            `json:"main"`
	Hooks       []string          `json:"hooks"`
	APIKeys     map[string]string `json:"api_keys,omitempty"`
	Dependencies []string         `json:"dependencies,omitempty"`
}

// PluginRuntime manages the plugin lifecycle including loading, unloading,
// and hot reload capabilities.
type PluginRuntime struct {
	mu       sync.RWMutex
	plugins  map[string]*loadedPlugin
	watcher  *fsnotify.Watcher
	pluginsDir string
	hookRegistry interface {
		Register(name string, fn interface{}) error
	}
	closed     chan struct{}
}

// loadedPlugin holds a plugin instance along with its metadata.
type loadedPlugin struct {
	info    *PluginInfo
	plugin  ExecutablePlugin
	symlink string
}

// NewPluginRuntime creates a new plugin runtime instance.
func NewPluginRuntime(pluginsDir string) (*PluginRuntime, error) {
	if pluginsDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		pluginsDir = filepath.Join(homeDir, ".config", "freecode", "plugins")
	}

	// Ensure plugins directory exists
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugins directory: %w", err)
	}

	runtime := &PluginRuntime{
		plugins:    make(map[string]*loadedPlugin),
		pluginsDir: pluginsDir,
		closed:     make(chan struct{}),
	}

	// Start file watcher for hot reload
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}
	runtime.watcher = watcher

	// Watch the plugins directory
	if err := runtime.watcher.Add(pluginsDir); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to watch plugins directory: %w", err)
	}

	// Start hot reload handler
	go runtime.handleFileEvents()

	return runtime, nil
}

// handleFileEvents processes file system events for hot reload.
func (r *PluginRuntime) handleFileEvents() {
	for {
		select {
		case <-r.closed:
			return
		case event, ok := <-r.watcher.Events:
			if !ok {
				return
			}
			r.handleEvent(event)
		case err, ok := <-r.watcher.Errors:
			if !ok {
				return
			}
			// Log error but continue watching
			fmt.Printf("plugin watcher error: %v\n", err)
		}
	}
}

// handleEvent processes a single file system event.
func (r *PluginRuntime) handleEvent(event fsnotify.Event) {
	// Only care about plugin manifest changes
	if filepath.Ext(event.Name) != ".json" {
		return
	}

	// Extract plugin ID from path
	pluginID := filepath.Base(filepath.Dir(event.Name))

	r.mu.RLock()
	lp, exists := r.plugins[pluginID]
	r.mu.RUnlock()

	if !exists {
		return
	}

	switch {
	case event.Op&fsnotify.Write == fsnotify.Write:
		// Reload plugin on manifest change
		if err := r.reloadPlugin(pluginID, lp); err != nil {
			fmt.Printf("failed to reload plugin %s: %v\n", pluginID, err)
		}
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		// Plugin directory removed, unload it
		if err := r.Unload(pluginID); err != nil {
			fmt.Printf("failed to unload removed plugin %s: %v\n", pluginID, err)
		}
	}
}

// reloadPlugin reloads a plugin from disk.
func (r *PluginRuntime) reloadPlugin(id string, lp *loadedPlugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Shutdown existing plugin
	if lp.plugin != nil {
		if err := lp.plugin.Shutdown(); err != nil {
			fmt.Printf("warning: error shutting down plugin %s: %v\n", id, err)
		}
	}

	manifest, err := r.readManifest(lp.symlink)
	if err != nil {
		lp.info.State = PluginStateError
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	lp.info.Version = manifest.Version
	lp.info.Description = manifest.Description

	// Attempt to load the new plugin
	newPlugin, err := r.loadPluginFile(manifest, lp.symlink)
	if err != nil {
		lp.info.State = PluginStateError
		return fmt.Errorf("failed to reload plugin: %w", err)
	}

	// Initialize the plugin
	if err := newPlugin.Init(context.Background()); err != nil {
		lp.info.State = PluginStateError
		return fmt.Errorf("failed to initialize reloaded plugin: %w", err)
	}

	lp.plugin = newPlugin
	lp.info.State = PluginStateLoaded
	lp.info.LoadedAt = time.Now()

	return nil
}

// Load loads a plugin from the specified path.
func (r *PluginRuntime) Load(pluginPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Resolve symlinks to get actual plugin directory
	realPath, err := filepath.EvalSymlinks(pluginPath)
	if err != nil {
		realPath = pluginPath
	}

	// Check if already loaded
	for _, lp := range r.plugins {
		if lp.symlink == realPath {
			return fmt.Errorf("plugin already loaded from %s", pluginPath)
		}
	}

	// Read manifest
	manifest, err := r.readManifest(realPath)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrManifestNotFound, err)
	}

	// Check if plugin ID already loaded
	if _, exists := r.plugins[manifest.ID]; exists {
		return fmt.Errorf("plugin with ID %s already loaded", manifest.ID)
	}

	// Load the plugin
	pl, err := r.loadPluginFile(manifest, realPath)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPluginLoadFailed, err)
	}

	// Initialize plugin
	if err := pl.Init(context.Background()); err != nil {
		return fmt.Errorf("%w: %v", ErrPluginInitFailed, err)
	}

	// Create plugin info
	info := &PluginInfo{
		ID:          manifest.ID,
		Name:        manifest.Name,
		Version:     manifest.Version,
		Author:      manifest.Author,
		Description: manifest.Description,
		Path:        pluginPath,
		Hooks:       manifest.Hooks,
		APIKeys:     manifest.APIKeys,
		LoadedAt:    time.Now(),
		State:       PluginStateLoaded,
	}

	// Store loaded plugin
	r.plugins[manifest.ID] = &loadedPlugin{
		info:    info,
		plugin:  pl,
		symlink: realPath,
	}

	// Watch plugin directory for changes
	if err := r.watcher.Add(realPath); err != nil {
		fmt.Printf("warning: failed to watch plugin directory %s: %v\n", realPath, err)
	}

	return nil
}

// loadPluginFile loads a plugin from its manifest and main file.
func (r *PluginRuntime) loadPluginFile(manifest *PluginManifest, pluginDir string) (ExecutablePlugin, error) {
	// For Go plugins, we use the Go plugin package
	// The main file should be a Go plugin (.so file)
	mainPath := filepath.Join(pluginDir, manifest.Main)
	if manifest.Main == "" {
		mainPath = filepath.Join(pluginDir, "plugin.so")
	}

	// Check if main file exists
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		// Try as a directory plugin (more common for Go)
		return nil, fmt.Errorf("plugin main file not found: %s", mainPath)
	}

	// Open the plugin
	p, err := plugin.Open(mainPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up the plugin symbol
	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin symbol not found: %w", err)
	}

	// Type assert to our plugin interface
	pl, ok := symPlugin.(ExecutablePlugin)
	if !ok {
		return nil, fmt.Errorf("plugin does not implement ExecutablePlugin interface")
	}

	return pl, nil
}

// readManifest reads and parses a plugin.json file.
func (r *PluginRuntime) readManifest(pluginDir string) (*PluginManifest, error) {
	manifestPath := filepath.Join(pluginDir, "plugin.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Validate required fields
	if manifest.ID == "" {
		return nil, ErrInvalidManifest
	}
	if manifest.Name == "" {
		return nil, ErrInvalidManifest
	}

	return &manifest, nil
}

// Unload unloads a plugin by its ID.
func (r *PluginRuntime) Unload(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	lp, exists := r.plugins[id]
	if !exists {
		return fmt.Errorf("%w: %s", ErrPluginNotFound, id)
	}

	lp.info.State = PluginStateUnloading

	if lp.plugin != nil {
		if err := lp.plugin.Shutdown(); err != nil {
			fmt.Printf("warning: error shutting down plugin %s: %v\n", id, err)
		}
	}

	delete(r.plugins, id)

	// Stop watching the plugin directory
	if lp.symlink != "" {
		r.watcher.Remove(lp.symlink)
	}

	return nil
}

// List returns information about all loaded plugins.
func (r *PluginRuntime) List() []*PluginInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]*PluginInfo, 0, len(r.plugins))
	for _, lp := range r.plugins {
		infos = append(infos, lp.info)
	}
	return infos
}

// Get returns a plugin by its ID.
func (r *PluginRuntime) Get(id string) (ExecutablePlugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lp, exists := r.plugins[id]
	if !exists {
		return nil, ErrPluginNotFound
	}
	return lp.plugin, nil
}

// GetInfo returns plugin info by ID.
func (r *PluginRuntime) GetInfo(id string) (*PluginInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lp, exists := r.plugins[id]
	if !exists {
		return nil, ErrPluginNotFound
	}
	return lp.info, nil
}

// Discover discovers all plugins in the plugins directory.
func (r *PluginRuntime) Discover() ([]*PluginManifest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var manifests []*PluginManifest

	entries, err := os.ReadDir(r.pluginsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginDir := filepath.Join(r.pluginsDir, entry.Name())
		manifest, err := r.readManifest(pluginDir)
		if err != nil {
			// Skip plugins with invalid manifests
			continue
		}
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

// LoadAll discovers and loads all valid plugins from the plugins directory.
func (r *PluginRuntime) LoadAll() error {
	manifests, err := r.Discover()
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		pluginDir := filepath.Join(r.pluginsDir, manifest.ID)
		if err := r.Load(pluginDir); err != nil {
			// Log error but continue loading other plugins
			fmt.Printf("failed to load plugin %s: %v\n", manifest.ID, err)
		}
	}

	return nil
}

// Execute runs a plugin by ID with the given request.
func (r *PluginRuntime) Execute(ctx context.Context, id string, req Request) (*Response, error) {
	r.mu.RLock()
	lp, exists := r.plugins[id]
	r.mu.RUnlock()

	if !exists {
		return nil, ErrPluginNotFound
	}

	if lp.info.State != PluginStateLoaded {
		return nil, fmt.Errorf("plugin %s is not in loaded state: %s", id, lp.info.State)
	}

	return lp.plugin.Execute(ctx, req)
}

// Close shuts down the plugin runtime and unloads all plugins.
func (r *PluginRuntime) Close() error {
	close(r.closed)

	r.mu.Lock()
	defer r.mu.Unlock()

	// Unload all plugins
	for id := range r.plugins {
		lp := r.plugins[id]
		lp.info.State = PluginStateUnloading

		if lp.plugin != nil {
			if err := lp.plugin.Shutdown(); err != nil {
				fmt.Printf("warning: error shutting down plugin %s: %v\n", id, err)
			}
		}

		if lp.symlink != "" {
			r.watcher.Remove(lp.symlink)
		}
	}

	r.plugins = make(map[string]*loadedPlugin)

	return r.watcher.Close()
}

// RegisterHook registers a hook function for a plugin.
func (r *PluginRuntime) RegisterHook(pluginID, hookName string, hookFn interface{}) error {
	r.mu.RLock()
	_, exists := r.plugins[pluginID]
	r.mu.RUnlock()

	if !exists {
		return ErrPluginNotFound
	}

	// Hook registration would integrate with the hook registry
	// This is a placeholder for the hook integration
	return nil
}
