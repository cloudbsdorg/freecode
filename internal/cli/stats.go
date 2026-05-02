package cli

import (
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics",
	Long:  `Display freecode usage statistics.`,
}

var (
	statsSessions bool
	statsTools    bool
	statsAgents   bool
	statsAll      bool
)

var statsRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run and show stats",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	statsCmd.Flags().BoolVar(&statsSessions, "sessions", false, "Show session statistics")
	statsCmd.Flags().BoolVar(&statsTools, "tools", false, "Show tool usage statistics")
	statsCmd.Flags().BoolVar(&statsAgents, "agents", false, "Show agent statistics")
	statsCmd.Flags().BoolVar(&statsAll, "all", false, "Show all statistics")
	statsCmd.AddCommand(statsRunCmd)
}
