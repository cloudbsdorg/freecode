package cli

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var acpCmd = &cobra.Command{
	Use:   "acp",
	Short: "Start ACP (Agent Client Protocol) server",
	RunE:  runACP,
}

var (
	acpPort    int
	acpHost    string
	acpWorkDir string
)

func init() {
	acpCmd.Flags().IntVar(&acpPort, "port", 18792, "Server port")
	acpCmd.Flags().StringVar(&acpHost, "hostname", "127.0.0.1", "Bind host")
	acpCmd.Flags().StringVar(&acpWorkDir, "cwd", "", "Working directory")
	acpCmd.MarkFlagDirname("cwd")
	rootCmd.AddCommand(acpCmd)
}

func runACP(cmd *cobra.Command, args []string) error {
	addr := acpHost + ":" + strconv.Itoa(acpPort)

	fmt.Println("Starting ACP server...")
	fmt.Printf("Listening on %s\n", addr)

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	server := &http.Server{
		Addr:    addr,
		Handler: acpHandler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			shutdownCancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down ACP server...")
	shutdownCtx, shutdownCancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	return server.Shutdown(shutdownCtx)
}

func acpHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","protocol":"acp"}`))
	})

	mux.HandleFunc("/api/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"type":"connected"}`))
	})

	return mux
}

type acpClient struct {
	conn    net.Conn
	readCh  chan []byte
	writeCh chan []byte
}

func (c *acpClient) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			close(c.readCh)
			return
		}
		c.readCh <- buf[:n]
	}
}

func (c *acpClient) writeLoop() {
	for data := range c.writeCh {
		c.conn.Write(data)
	}
	c.conn.Close()
}
