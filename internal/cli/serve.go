package cli

import (
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start freecode server mode",
	Long:  `Start the freecode API server and web UI.`,
	RunE:  runServe,
}

var (
	servePort    int
	serveHost    string
	enableMCP    bool
	enableWebUI  bool
)

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 18792, "API server port")
	serveCmd.Flags().StringVar(&serveHost, "host", "127.0.0.1", "Bind host (localhost only)")
	serveCmd.Flags().BoolVar(&enableMCP, "enable-mcp", true, "Enable MCP server")
	serveCmd.Flags().BoolVar(&enableWebUI, "enable-webui", true, "Enable web UI")
}

func runServe(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}
