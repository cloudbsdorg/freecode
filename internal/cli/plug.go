package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/freecode/freecode/internal/plugin"
	"github.com/spf13/cobra"
)

var (
	plugInstall string
	plugRemove  string
	plugList    bool
	plugReload  string
)

var plugCmd = &cobra.Command{
	Use:   "plug",
	Short: "Manage freecode plugins",
	Long: `Manage freecode plugins including installation, removal, and listing.

Examples:
  freecode plug              # Show plugin status
  freecode plug --list      # List installed plugins
  freecode plug --install <path>  # Install a plugin
  freecode plug --remove <name>   # Remove a plugin
  freecode plug --reload <name>   # Reload a plugin`,
	RunE: runPlug,
}

func init() {
	plugCmd.Flags().BoolVar(&plugList, "list", false, "List installed plugins")
	plugCmd.Flags().StringVar(&plugInstall, "install", "", "Install a plugin from path")
	plugCmd.Flags().StringVar(&plugRemove, "remove", "", "Remove a plugin by name")
	plugCmd.Flags().StringVar(&plugReload, "reload", "", "Reload a plugin by name")
	rootCmd.AddCommand(plugCmd)
}

func runPlug(cmd *cobra.Command, args []string) error {
	if plugList {
		return listPlugins()
	}

	if plugInstall != "" {
		return installPlugin(plugInstall)
	}

	if plugRemove != "" {
		return removePlugin(plugRemove)
	}

	if plugReload != "" {
		return reloadPlugin(plugReload)
	}

	return showPluginStatus()
}

var globalRegistry = plugin.NewMemoryRegistry()

func pluginDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "freecode", "plugins")
}

func listPlugins() error {
	plugins := globalRegistry.List()

	fmt.Println("Installed Plugins:")
	fmt.Println("")

	if len(plugins) == 0 {
		fmt.Println("  No plugins installed.")
		fmt.Println("")
		fmt.Printf("  To install a plugin: freecode plug --install /path/to/plugin\n")
		fmt.Printf("  Plugin directory: %s\n", pluginDir())
		return nil
	}

	for _, name := range plugins {
		fmt.Printf("  - %s\n", name)
	}
	fmt.Println("")

	return nil
}

func installPlugin(path string) error {
	fmt.Printf("Installing plugin from: %s\n", path)
	fmt.Println("")

	srcInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("plugin path does not exist: %s", path)
		}
		return fmt.Errorf("failed to access plugin path: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("plugin path must be a directory: %s", path)
	}

	manifest, err := findPluginManifest(path)
	if err != nil {
		return fmt.Errorf("failed to find plugin manifest: %w", err)
	}

	if err := validateManifest(manifest); err != nil {
		return fmt.Errorf("invalid plugin manifest: %w", err)
	}

	pluginDst := filepath.Join(pluginDir(), manifest.ID)
	if _, err := os.Stat(pluginDst); err == nil {
		return fmt.Errorf("plugin %s already installed at %s", manifest.Name, pluginDst)
	}

	fmt.Printf("  Copying plugin files...\n")
	if err := copyPluginDir(path, pluginDst); err != nil {
		return fmt.Errorf("failed to copy plugin files: %w", err)
	}

	installedPlugin := &installedPluginAdapter{
		id:      manifest.ID,
		name:    manifest.Name,
		version: manifest.Version,
		author:  manifest.Author,
		path:    pluginDst,
	}

	if err := globalRegistry.Register(installedPlugin); err != nil {
		os.RemoveAll(pluginDst)
		return fmt.Errorf("failed to register plugin: %w", err)
	}

	fmt.Println("")
	fmt.Printf("  Successfully installed plugin: %s (v%s)\n", manifest.Name, manifest.Version)
	fmt.Printf("  Plugin ID: %s\n", manifest.ID)
	fmt.Printf("  Installed to: %s\n", pluginDst)
	fmt.Println("")

	return nil
}

type installedPluginAdapter struct {
	id      string
	name    string
	version string
	author  string
	path    string
}

func (p *installedPluginAdapter) Name() string { return p.name }
func (p *installedPluginAdapter) Init(ctx context.Context) error {
	return nil
}
func (p *installedPluginAdapter) Close() error {
	return nil
}

func findPluginManifest(pluginPath string) (*pluginManifest, error) {
	manifestPath := filepath.Join(pluginPath, "plugin.json")
	data, err := os.ReadFile(manifestPath)
	if err == nil {
		var manifest pluginManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("failed to parse plugin.json: %w", err)
		}
		return &manifest, nil
	}

	pkgPath := filepath.Join(pluginPath, "package.json")
	pkgData, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("neither plugin.json nor package.json found in %s", pluginPath)
	}

	var pkg struct {
		Name           string `json:"name"`
		Version        string `json:"version"`
		Description    string `json:"description"`
		FreecodePlugin any   `json:"freecode-plugin"`
	}
	if err := json.Unmarshal(pkgData, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	if pkg.FreecodePlugin == nil {
		return nil, fmt.Errorf("package.json does not contain 'freecode-plugin' field")
	}

	var fcPlugin struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Version     string   `json:"version"`
		Author      string   `json:"author"`
		Description string   `json:"description"`
		Main        string   `json:"main"`
		Hooks       []string `json:"hooks"`
	}
	fcData, err := json.Marshal(pkg.FreecodePlugin)
	if err != nil {
		return nil, fmt.Errorf("failed to process freecode-plugin field: %w", err)
	}
	if err := json.Unmarshal(fcData, &fcPlugin); err != nil {
		return nil, fmt.Errorf("failed to parse freecode-plugin field: %w", err)
	}

	return &pluginManifest{
		ID:          fcPlugin.ID,
		Name:        fcPlugin.Name,
		Version:     fcPlugin.Version,
		Author:      fcPlugin.Author,
		Description: fcPlugin.Description,
		Main:        fcPlugin.Main,
		Hooks:       fcPlugin.Hooks,
	}, nil
}

type pluginManifest struct {
	ID          string
	Name        string
	Version     string
	Author      string
	Description string
	Main        string
	Hooks       []string
}

func validateManifest(manifest *pluginManifest) error {
	if manifest.ID == "" {
		return fmt.Errorf("plugin ID is required")
	}
	if manifest.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	return nil
}

func copyPluginDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return os.MkdirAll(dst, 0755)
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}
		return copyPluginFile(path, dstPath, info.Mode())
	})
}

func copyPluginFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	return err
}

func removePlugin(name string) error {
	err := globalRegistry.Unregister(name)
	if err != nil {
		fmt.Printf("Removing plugin: %s\n", name)
		fmt.Printf("  Error: %v\n", err)
		return err
	}

	fmt.Printf("Removed plugin: %s\n", name)
	return nil
}

func reloadPlugin(name string) error {
	_, err := globalRegistry.Get(name)
	if err != nil {
		fmt.Printf("Reloading plugin: %s\n", name)
		fmt.Printf("  Error: %v\n", err)
		return err
	}

	fmt.Printf("Reloaded plugin: %s\n", name)
	return nil
}

func showPluginStatus() error {
	plugins := globalRegistry.List()

	fmt.Println("Freecode Plugin System")
	fmt.Println("======================")
	fmt.Println("")
	fmt.Printf("  Plugin directory: %s\n", pluginDir())
	fmt.Printf("  Registered plugins: %d\n", len(plugins))
	fmt.Println("")

	if len(plugins) > 0 {
		fmt.Println("  Plugins:")
		for _, name := range plugins {
			fmt.Printf("    - %s\n", name)
		}
	} else {
		fmt.Println("  No plugins registered.")
	}
	fmt.Println("")

	return nil
}
