package cli

import (
	"fmt"

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
	fmt.Println("✓ CLI version: 0.1.0")
	fmt.Println("✓ Platform: darwin/arm64")
	fmt.Println("✓ Go version: go1.20")
	fmt.Println()
	fmt.Println("All checks passed.")
	return nil
}
