package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	cmdDryRun  bool
	cmdVerbose bool
)

var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Execute freecode internal commands",
	Long: `Execute freecode internal commands for debugging and development.

Examples:
  freecode cmd list          # List available commands
  freecode cmd exec <name>  # Execute a specific command`,
	RunE: runCmdCommand,
}

func init() {
	cmdCmd.Flags().BoolVar(&cmdDryRun, "dry-run", false, "Show what would be executed without running")
	cmdCmd.Flags().BoolVar(&cmdVerbose, "verbose", false, "Verbose output")
}

func runCmdCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return listCommands()
	}

	subcmd := args[0]
	switch subcmd {
	case "list":
		return listCommands()
	case "exec":
		if len(args) < 2 {
			return fmt.Errorf("exec requires a command name")
		}
		return execCommand(args[1], args[2:])
	default:
		return fmt.Errorf("unknown subcommand: %s", subcmd)
	}
}

func listCommands() error {
	commands := []struct {
		name        string
		description string
	}{
		{"health", "Check system health"},
		{"cache-clear", "Clear internal caches"},
		{"cache-stats", "Show cache statistics"},
		{"session-list", "List active sessions"},
		{"session-cleanup", "Clean up old sessions"},
		{"config-reload", "Reload configuration"},
		{"config-validate", "Validate configuration files"},
		{"tool-list", "List registered tools"},
		{"hook-list", "List registered hooks"},
		{"agent-list", "List available agents"},
	}

	fmt.Println("Available Commands:")
	fmt.Println("")
	for _, c := range commands {
		fmt.Printf("  %-20s %s\n", c.name, c.description)
	}

	return nil
}

func execCommand(name string, args []string) error {
	if cmdDryRun {
		fmt.Printf("Would execute: freecode %s %v\n", name, args)
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable: %w", err)
	}

	execArgs := append([]string{name}, args...)
	execCmd := exec.Command(exe, execArgs...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = os.Environ()

	if cmdVerbose {
		fmt.Printf("Executing: %s %v\n", exe, execArgs)
	}

	return execCmd.Run()
}
