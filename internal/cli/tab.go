package cli

import (
	"github.com/spf13/cobra"
)

var tabCmd = &cobra.Command{
	Use:   "tab",
	Short: "Manage tabs",
	Long:  `Create, close, and manage session tabs.`,
}

var (
	tabNew    bool
	tabClose  string
	tabList   bool
	tabMove   string
	tabRename string
)

func init() {
	tabCmd.AddCommand(tabNewCmd)
	tabCmd.AddCommand(tabCloseCmd)
	tabCmd.AddCommand(tabListCmd)
	tabCmd.AddCommand(tabMoveCmd)
	tabCmd.AddCommand(tabRenameCmd)
}

var tabNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new tab",
	RunE:  runTabNew,
}

var tabCloseCmd = &cobra.Command{
	Use:   "close [tab-id]",
	Short: "Close a tab",
	RunE:  runTabClose,
}

var tabListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tabs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var tabMoveCmd = &cobra.Command{
	Use:   "move [session-id] [tab-id]",
	Short: "Move session to tab",
	RunE:  runTabMove,
}

var tabRenameCmd = &cobra.Command{
	Use:   "rename [tab-id] [name]",
	Short: "Rename a tab",
	RunE:  runTabRename,
}

func runTabNew(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runTabClose(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runTabMove(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runTabRename(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}
