package cli

import (
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agents",
	Long:  `List, start, and manage freecode agents.`,
}

var (
	agentList   bool
	agentStart  string
	agentStop   string
	agentStatus string
)

func init() {
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentStartCmd)
	agentCmd.AddCommand(agentStopCmd)
	agentCmd.AddCommand(agentStatusCmd)
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var agentStartCmd = &cobra.Command{
	Use:   "start [agent-name]",
	Short: "Start an agent",
	RunE:  runAgentStart,
}

var agentStopCmd = &cobra.Command{
	Use:   "stop [agent-name]",
	Short: "Stop an agent",
	RunE:  runAgentStop,
}

var agentStatusCmd = &cobra.Command{
	Use:   "status [agent-name]",
	Short: "Get agent status",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func runAgentStart(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runAgentStop(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}
