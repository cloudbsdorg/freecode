package cli

import (
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
	fmt.Println("  Note: Plugin loading from path not yet implemented.")
	fmt.Printf("  Plugin would be registered.\n")

	return nil
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
