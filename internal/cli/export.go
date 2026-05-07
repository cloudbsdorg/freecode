package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [sessionID]",
	Short: "Export session data as JSON",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  runExport,
}

var exportOutput string

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path")
	rootCmd.AddCommand(exportCmd)
}

type exportedSession struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Messages  []exportedMessage `json:"messages,omitempty"`
}

type exportedMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func runExport(cmd *cobra.Command, args []string) error {
	sessionID := ""
	if len(args) > 0 {
		sessionID = args[0]
	}

	home, _ := os.UserHomeDir()
	sessionsDir := filepath.Join(home, ".local", "share", "freecode", "sessions")

	if sessionID == "" {
		sessionID = getMostRecentSession(sessionsDir)
		if sessionID == "" {
			return fmt.Errorf("no session ID provided and no sessions found")
		}
	}

	sessionFile := filepath.Join(sessionsDir, sessionID+".json")
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("session %s not found", sessionID)
		}
		return fmt.Errorf("failed to read session: %w", err)
	}

	var session exportedSession
	if err := json.Unmarshal(data, &session); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	output := os.Stdout
	if exportOutput != "" {
		output, err = os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	}

	enc := json.NewEncoder(output)
	enc.SetIndent("", "  ")
	if err := enc.Encode(session); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Exported session %s\n", sessionID)
	return nil
}

func getMostRecentSession(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		return ""
	}

	var newest string
	var newestTime time.Time
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		info, _ := e.Info()
		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newest = strings.TrimSuffix(e.Name(), ".json")
		}
	}
	return newest
}
