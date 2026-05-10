package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/freecode/freecode/internal/agent"
	"github.com/freecode/freecode/internal/config"
	"github.com/freecode/freecode/internal/mcp"
	"github.com/freecode/freecode/internal/server"
	"github.com/freecode/freecode/internal/tool"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a headless freecode server",
	Long:  `Start the freecode API and MCP server without the TUI.`,
	RunE:  runServe,
}

var (
	servePort int
	serveMCP  int
	serveHost string
	serveMDNS bool
)

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 18792, "API server port")
	serveCmd.Flags().IntVar(&serveMCP, "mcp-port", 18793, "MCP server port")
	serveCmd.Flags().StringVar(&serveHost, "hostname", "127.0.0.1", "Bind host")
	serveCmd.Flags().BoolVar(&serveMDNS, "mdns", false, "Enable mDNS discovery")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg := config.DefaultConfig()
	eng := agent.NewEngine(cfg)

	mcpServer := mcp.NewServer(serveMCP)
	for _, t := range eng.ListTools() {
		schema := t.Schema()
		mcpServer.RegisterTool(schema.Name, schema.Description, func(args map[string]interface{}) (interface{}, error) {
			req := tool.Request{
				Name:      schema.Name,
				Arguments: args,
			}
			t, exists := eng.GetTool(schema.Name)
			if !exists {
				return nil, fmt.Errorf("tool not found: %s", schema.Name)
			}
			result, err := t.Execute(context.Background(), req)
			if err != nil {
				return nil, err
			}
			if result != nil && result.Error != nil {
				return nil, result.Error
			}
			return result.Result, nil
		})
	}

	addr := serveHost + ":" + strconv.Itoa(servePort)
	srv := server.New(addr)
	srv.SetupRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := mcpServer.Start(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		}
	}()

	go func() {
		if err := srv.Start(ctx); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "API server error: %v\n", err)
			cancel()
		}
	}()

	host := serveHost
	if serveMDNS {
		host = "0.0.0.0"
	}

	fmt.Printf("Freecode server running:\n")
	fmt.Printf("  API:  http://%s:%d\n", host, servePort)
	fmt.Printf("  MCP:  http://%s:%d (MCP protocol)\n", host, serveMCP)
	fmt.Printf("\nPress Ctrl+C to stop\n")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil
}
