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

	"github.com/freecode/freecode/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a headless freecode server",
	Long:  `Start the freecode API server without the TUI.`,
	RunE:  runServe,
}

var (
	servePort    int
	serveHost    string
	serveMDNS    bool
	serveMDNSDomain string
)

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 18792, "Server port (0 = random)")
	serveCmd.Flags().StringVar(&serveHost, "hostname", "127.0.0.1", "Bind host")
	serveCmd.Flags().BoolVar(&serveMDNS, "mdns", false, "Enable mDNS discovery")
	serveCmd.Flags().StringVar(&serveMDNSDomain, "mdns-domain", "freecode.local", "mDNS domain name")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	addr := serveHost + ":" + strconv.Itoa(servePort)

	srv := server.New(addr)
	srv.SetupRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := srv.Start(ctx); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			cancel()
		}
	}()

	host := serveHost
	if serveMDNS {
		host = "0.0.0.0"
	}

	fmt.Printf("Freecode server listening on http://%s:%d\n", host, servePort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil
}