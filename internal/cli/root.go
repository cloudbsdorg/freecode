package cli

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	yolo    bool
)

var rootCmd = &cobra.Command{
	Use:   "freecode",
	Short: "Freecode - AI coding assistant",
	Long: `Freecode is a platform-independent AI coding assistant.

Built with Go for FreeBSD 16, Linux, macOS, and IllumOS.`,
	SilenceErrors: true,
	SilenceUsage:   true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/freecode/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&yolo, "yolo", false, "skip all confirmations")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.AddCommand(tabCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(versionCmd)
}
