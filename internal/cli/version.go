package cli

import (
	"fmt"

	"github.com/freecode/freecode/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print freecode version",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.Get()
		fmt.Printf("freecode version %s\n", info.Version)
		fmt.Printf("Platform: %s\n", info.Platform)
		fmt.Printf("Go: %s\n", info.GoVersion)
	},
}
