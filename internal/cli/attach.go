package cli

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var attachCmd = &cobra.Command{
	Use:   "attach <url>",
	Short: "Attach to a running freecode server",
	Args:  cobra.ExactArgs(1),
	RunE:  runAttach,
}

var (
	attachDir      string
	attachContinue bool
	attachSession  string
	attachFork     bool
	attachPassword string
)

func init() {
	attachCmd.Flags().StringVar(&attachDir, "dir", "", "Directory to run in")
	attachCmd.Flags().BoolVarP(&attachContinue, "continue", "c", false, "Continue the last session")
	attachCmd.Flags().StringVarP(&attachSession, "session", "s", "", "Session ID to continue")
	attachCmd.Flags().BoolVar(&attachFork, "fork", false, "Fork the session when continuing")
	attachCmd.Flags().StringVarP(&attachPassword, "password", "p", "", "Server password")
	rootCmd.AddCommand(attachCmd)
}

func runAttach(cmd *cobra.Command, args []string) error {
	serverURL := args[0]

	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Host == "" {
		parsedURL.Host = "localhost:18792"
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}

	targetURL := parsedURL.String()
	fmt.Printf("Attaching to: %s\n", targetURL)

	if attachDir != "" {
		if err := os.Chdir(attachDir); err != nil {
			return fmt.Errorf("failed to change directory: %w", err)
		}
		fmt.Printf("Working directory: %s\n", attachDir)
	}

	headers := make(http.Header)
	if attachPassword != "" {
		headers.Set("Authorization", "Basic "+basicAuth("opencode", attachPassword))
	} else if envPassword := os.Getenv("FREECODE_SERVER_PASSWORD"); envPassword != "" {
		headers.Set("Authorization", "Basic "+basicAuth("opencode", envPassword))
	}

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", targetURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v[0])
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Warning: Could not reach server - %v\n", err)
		fmt.Println("Starting TUI anyway (may fail if server is not running)...")
	}

	if resp != nil {
		resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Printf("Warning: Server returned status %d\n", resp.StatusCode)
		}
	}

	fmt.Println("\nNote: TUI attach mode requires the server to have an active session")
	fmt.Printf("URL: %s\n", targetURL)

	if attachSession != "" {
		fmt.Printf("Session: %s\n", attachSession)
	}
	if attachContinue {
		fmt.Println("Mode: continue last session")
	}

	return startLocalTUI(targetURL, headers)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + btoa(auth)
}

func btoa(s string) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := ""
	for i := 0; i < len(s); i += 3 {
		var n uint32
		n |= uint32(s[i]) << 16
		if i+1 < len(s) {
			n |= uint32(s[i+1]) << 8
		}
		if i+2 < len(s) {
			n |= uint32(s[i+2])
		}
		result += string(alphabet[n>>18&63]) + string(alphabet[n>>12&63])
		if i+1 < len(s) {
			result += string(alphabet[n>>6&63])
		}
		if i+2 < len(s) {
			result += string(alphabet[n&63])
		}
	}
	return result
}

func startLocalTUI(serverURL string, headers http.Header) error {
	fmt.Println("\nStarting TUI...")
	fmt.Println("Note: Full TUI attach functionality requires server-side session support")

	p := os.Getenv("FREECODE_TUI_URL")
	if p == "" {
		p = serverURL
	}
	fmt.Printf("TUI would connect to: %s\n", p)
	return nil
}
