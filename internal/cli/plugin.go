package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:     "plugin <module>",
	Aliases: []string{"plug"},
	Short:   "Install a freecode plugin",
	Long: `Install a plugin from npm and update configuration.
The module should be an npm package name or URL.`,
	Args: cobra.ExactArgs(1),
	RunE: runPlugin,
}

var (
	pluginGlobal bool
	pluginForce  bool
)

func init() {
	pluginCmd.Flags().BoolVarP(&pluginGlobal, "global", "g", false, "Install globally")
	pluginCmd.Flags().BoolVarP(&pluginForce, "force", "f", false, "Replace existing plugin version")
	rootCmd.AddCommand(pluginCmd)
}

func runPlugin(cmd *cobra.Command, args []string) error {
	module := args[0]

	fmt.Printf("Installing plugin: %s\n", module)

	if os.Getenv("FREECODE_PLUGINS_ENABLED") != "true" {
		fmt.Println("Warning: Plugin system is experimental")
	}

	scope := "local"
	if pluginGlobal {
		scope = "global"
	}
	fmt.Printf("Scope: %s\n", scope)

	fmt.Println("\nNote: Plugin installation requires npm/node.js")
	fmt.Printf("To install manually: npm install -g %s\n", module)

	execPath, err := exec.LookPath("npm")
	if err != nil {
		fmt.Println("\nnpm not found. Please install Node.js to use plugins.")
		return nil
	}
	fmt.Printf("Found npm at: %s\n", execPath)

	fmt.Printf("\nRunning: npm install -g %s\n", module)

	installCmd := exec.Command("npm", "install", "-g", module)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Stdin = os.Stdin

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("npm install failed: %w", err)
	}

	fmt.Println("\nPlugin installed successfully!")
	fmt.Println("Restart freecode to load the plugin.")
	return nil
}
