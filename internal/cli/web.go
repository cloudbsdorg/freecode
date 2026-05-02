package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/freecode/freecode/internal/server"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start freecode server and open web interface",
	Long: `Start the freecode API server and open the web interface in your browser.

The web interface provides:
- API status and health monitoring
- REST API endpoints for sessions, agents, and tools
- Simple web UI for basic operations

Note: The full interactive experience is in the TUI application.`,
	RunE: runWeb,
}

var (
	webPort        int
	webHost        string
	webNoBrowser   bool
	webDisableUI   bool
)

func init() {
	webCmd.Flags().IntVar(&webPort, "port", 18791, "Web UI port")
	webCmd.Flags().StringVar(&webHost, "hostname", "127.0.0.1", "Bind host (use 0.0.0.0 for network access)")
	webCmd.Flags().BoolVar(&webNoBrowser, "no-browser", false, "Don't open browser automatically")
	webCmd.Flags().BoolVar(&webDisableUI, "no-ui", false, "Disable web UI (API only)")
	rootCmd.AddCommand(webCmd)
}

func runWeb(cmd *cobra.Command, args []string) error {
	addr := webHost + ":" + strconv.Itoa(webPort)

	srv := server.New(addr)
	srv.SetupWebRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := srv.Start(ctx); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			cancel()
		}
	}()

	time.Sleep(100 * time.Millisecond)

	url := fmt.Sprintf("http://%s:%d", webHost, webPort)

	fmt.Println()
	fmt.Println("  Freecode Web Interface")
	fmt.Println("  =====================")
	fmt.Println()
	fmt.Printf("  Local:   %s\n", url)
	fmt.Printf("  Health:  %s/health\n", url)
	fmt.Printf("  API:     %s/api/v1\n", url)

	if webHost == "0.0.0.0" {
		fmt.Println()
		fmt.Println("  Network access enabled (0.0.0.0)")
		fmt.Println("  Other devices can access via the IP address above")
	}

	fmt.Println()

	if !webNoBrowser && !webDisableUI {
		if err := openBrowser(url); err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: Could not open browser: %v\n", err)
		} else {
			fmt.Println("  Browser opened automatically")
		}
	}

	fmt.Println()
	fmt.Println("  Press Ctrl+C to stop the server")
	fmt.Println()

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

func openBrowser(url string) error {
	var err error

	switch {
	case isCommandAvailable("xdg-open"):
		err = exec.Command("xdg-open", url).Run()
	case isCommandAvailable("open"):
		err = exec.Command("open", url).Run()
	case isCommandAvailable("start"):
		err = exec.Command("cmd", "/c", "start", url).Run()
	default:
		return fmt.Errorf("no browser command found")
	}

	return err
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
