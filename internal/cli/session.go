package cli

import (
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage sessions",
	Long:  `List, read, search, import, and export sessions.`,
}

var (
	sessionList   bool
	sessionRead   string
	sessionSearch string
	sessionExport string
	sessionImport string
	sessionDelete string
	sessionInfo   string
)

func init() {
	sessionCmd.AddCommand(sessionListCmd)
	sessionCmd.AddCommand(sessionReadCmd)
	sessionCmd.AddCommand(sessionSearchCmd)
	sessionCmd.AddCommand(sessionExportCmd)
	sessionCmd.AddCommand(sessionImportCmd)
	sessionCmd.AddCommand(sessionDeleteCmd)
	sessionCmd.AddCommand(sessionInfoCmd)
}

var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sessions",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var sessionReadCmd = &cobra.Command{
	Use:   "read [session-id]",
	Short: "Read session messages",
	RunE:  runSessionRead,
}

var sessionSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search session history",
	RunE:  runSessionSearch,
}

var sessionExportCmd = &cobra.Command{
	Use:   "export [session-id]",
	Short: "Export session to file",
	RunE:  runSessionExport,
}

var sessionImportCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import session from file",
	RunE:  runSessionImport,
}

var sessionDeleteCmd = &cobra.Command{
	Use:   "delete [session-id]",
	Short: "Delete a session",
	RunE:  runSessionDelete,
}

var sessionInfoCmd = &cobra.Command{
	Use:   "info [session-id]",
	Short: "Get session info",
	RunE:  runSessionInfo,
}

func runSessionRead(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runSessionSearch(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runSessionExport(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runSessionImport(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runSessionDelete(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}

func runSessionInfo(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}
