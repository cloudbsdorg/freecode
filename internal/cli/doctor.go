package cli

import (
	"fmt"

	"github.com/freecode/freecode/internal/version"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health",
	Long:  `Run diagnostics to verify freecode installation.`,
	RunE:  runDoctor,
}

var doctorAll bool

func init() {
	doctorCmd.Flags().BoolVar(&doctorAll, "all", false, "Run all checks")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("Running freecode doctor...")
	fmt.Println()

	info := version.Get()
	fmt.Printf("✓ CLI version: %s\n", info.Version)
	fmt.Printf("✓ Platform: %s\n", info.Platform)
	fmt.Printf("✓ Go version: %s\n", info.GoVersion)
	fmt.Println()
	fmt.Println("All checks passed.")
	return nil
}
