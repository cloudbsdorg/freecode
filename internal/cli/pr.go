package cli

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr <number>",
	Short: "Fetch and checkout a GitHub PR branch, then run freecode",
	Args:  cobra.ExactArgs(1),
	RunE:  runPR,
}

func init() {
	rootCmd.AddCommand(prCmd)
}

func runPR(cmd *cobra.Command, args []string) error {
	prNumber := args[0]

	fmt.Printf("Fetching and checking out PR #%s...\n", prNumber)

	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI not found. Please install GitHub CLI: https://cli.github.com")
	}

	localBranch := "pr/" + prNumber

	checkoutCmd := exec.Command("gh", "pr", "checkout", prNumber, "--branch", localBranch, "--force")
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr
	checkoutCmd.Stdin = os.Stdin

	if err := checkoutCmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout PR #%s: %w", prNumber, err)
	}

	fmt.Printf("Successfully checked out PR #%s as branch '%s'\n", prNumber, localBranch)

	prViewCmd := exec.Command("gh", "pr", "view", prNumber, "--json", "body", "--jq", ".body")
	body, err := prViewCmd.Output()
	if err == nil {
		bodyStr := strings.TrimSpace(string(body))
		if sessionURL := extractSessionURL(bodyStr); sessionURL != "" {
			fmt.Printf("Found freecode session: %s\n", sessionURL)
			fmt.Println("Importing session...")
			importCmd := exec.Command("freecode", "import", sessionURL)
			importCmd.Stdout = os.Stdout
			importCmd.Stderr = os.Stderr
			importCmd.Run()
		}
	}

	fmt.Println("\nStarting freecode...")
	runCmd := exec.Command("freecode")
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin
	return runCmd.Run()
}

func extractSessionURL(body string) string {
	patterns := []string{
		"https://opncd.ai/s/",
		"https://freecode.ai/s/",
		"https://localhost:",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(body, pattern); idx != -1 {
			end := idx + len(pattern)
			for end < len(body) && isSessionChar(body[end]) {
				end++
			}
			return body[idx:end]
		}
	}
	return ""
}

func isSessionChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
}

func isValidURL(u string) bool {
	_, err := url.Parse(u)
	return err == nil
}
