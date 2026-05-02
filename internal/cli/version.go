package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print freecode version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("freecode version 0.1.0")
		fmt.Println("Platform: darwin/arm64")
		fmt.Println("Go: go1.20")
	},
}
