package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall freecode and remove related files",
	Long: `Remove freecode binary and optionally configuration/data directories.
Use --keep-config to preserve configuration and --keep-data to preserve session data.`,
	RunE: runUninstall,
}

var (
	uninstallKeepConfig bool
	uninstallKeepData   bool
	uninstallDryRun     bool
	uninstallForce      bool
)

func init() {
	uninstallCmd.Flags().BoolVarP(&uninstallKeepConfig, "keep-config", "c", false, "Keep configuration files")
	uninstallCmd.Flags().BoolVarP(&uninstallKeepData, "keep-data", "d", false, "Keep session data and snapshots")
	uninstallCmd.Flags().BoolVar(&uninstallDryRun, "dry-run", false, "Show what would be removed without removing")
	uninstallCmd.Flags().BoolVarP(&uninstallForce, "force", "f", false, "Skip confirmation prompts")
	rootCmd.AddCommand(uninstallCmd)
}

type removalTarget struct {
	path  string
	label string
	keep  bool
}

func runUninstall(cmd *cobra.Command, args []string) error {
	fmt.Println()
	printUninstallLogo()
	fmt.Println()

	targets := collectRemovalTargets()

	if !uninstallDryRun {
		fmt.Println("The following will be removed:")
	} else {
		fmt.Println("Dry run - the following would be removed:")
	}

	hasErrors := false
	for _, t := range targets {
		exists, _ := dirExists(t.path)
		if !exists {
			continue
		}

		size, _ := getDirSize(t.path)
		status := ""
		prefix := "✓"
		if t.keep {
			prefix = "○"
			status = " (keeping)"
		}
		fmt.Printf("  %s %s: %s %s%s\n", prefix, t.label, shortenPath(t.path), formatSize(size), status)
	}

	if uninstallDryRun {
		fmt.Println("\nNo changes made (dry run)")
		return nil
	}

	if !uninstallForce {
		fmt.Print("\nAre you sure you want to uninstall? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	fmt.Println("\nRemoving files...")

	for _, t := range targets {
		if t.keep {
			fmt.Printf("  Skipping %s (--keep-%s)\n", t.label, strings.ToLower(t.label))
			continue
		}

		exists, _ := dirExists(t.path)
		if !exists {
			continue
		}

		if err := os.RemoveAll(t.path); err != nil {
			fmt.Printf("  Failed to remove %s: %v\n", t.label, err)
			hasErrors = true
			continue
		}
		fmt.Printf("  Removed %s\n", t.label)
	}

	if hasErrors {
		fmt.Println("\nSome operations failed.")
	} else {
		fmt.Println("\nThank you for using Freecode!")
	}

	return nil
}

func collectRemovalTargets() []removalTarget {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "freecode")
	dataDir := filepath.Join(homeDir, ".local", "share", "freecode")
	cacheDir := filepath.Join(os.TempDir(), "freecode-cache")

	return []removalTarget{
		{path: configDir, label: "Config", keep: uninstallKeepConfig},
		{path: dataDir, label: "Data", keep: uninstallKeepData},
		{path: cacheDir, label: "Cache", keep: false},
	}
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.1f GB", float64(bytes)/(1024*1024*1024))
}

func shortenPath(p string) string {
	home, _ := os.UserHomeDir()
	if home != "" && strings.HasPrefix(p, home) {
		return strings.Replace(p, home, "~", 1)
	}
	return p
}

func printUninstallLogo() {
	logo := `██╗     ██╗ ██████╗ ██████╗ ███████╗██████╗ ███╗   ██╗██╗███╗   ██╗ ██████╗
██║     ██║██╔════╝██╔═══██╗██╔════╝██╔══██╗████╗  ██║██║████╗  ██║██╔════╝
██║     ██║██║     ██║   ██║█████╗  ██████╔╝██╔██╗ ██║██║██╔██╗ ██║██║  ███╗
██║     ██║██║     ██║   ██║██╔══╝  ██╔══██╗██║╚██╗██║██║██║╚██╗██║██║   ██║
███████╗██║╚██████╗╚██████╔╝███████╗██║  ██║██║ ╚████║██║██║ ╚████║╚██████╔╝
╚══════╝╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝╚═╝  ╚═══╝ ╚═════╝`
	fmt.Println(logo)
}
