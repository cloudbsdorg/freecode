package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import session data from JSON file or URL",
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

type importSession struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	CreatedAt string          `json:"created_at"`
	Messages  []importMessage `json:"messages,omitempty"`
}

type importMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func runImport(cmd *cobra.Command, args []string) error {
	input := args[0]

	var data []byte
	var err error

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		data, err = fetchFromURL(input)
		if err != nil {
			return fmt.Errorf("failed to fetch URL: %w", err)
		}
	} else {
		data, err = os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	}

	var session importSession
	if err := json.Unmarshal(data, &session); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	home, _ := os.UserHomeDir()
	sessionsDir := filepath.Join(home, ".local", "share", "freecode", "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	sessionID := session.ID
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	sessionFile := filepath.Join(sessionsDir, sessionID+".json")
	if err := os.WriteFile(sessionFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	fmt.Printf("Imported session: %s\n", sessionID)
	return nil
}

func fetchFromURL(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func generateSessionID() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}
