package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Manage GitHub agent and integrations",
}

var githubPRsCmd = &cobra.Command{
	Use:   "prs [owner/repo]",
	Short: "List recent PRs",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  runGithubPRs,
}

var githubReviewCmd = &cobra.Command{
	Use:   "review <pr-number>",
	Short: "Review a PR",
	Args:  cobra.ExactArgs(1),
	RunE:  runGithubReview,
}

var githubCommentCmd = &cobra.Command{
	Use:   "comment <pr-number> <comment>",
	Short: "Comment on a PR",
	Args:  cobra.ExactArgs(2),
	RunE:  runGithubComment,
}

func init() {
	githubCmd.AddCommand(githubPRsCmd)
	githubCmd.AddCommand(githubReviewCmd)
	githubCmd.AddCommand(githubCommentCmd)
	rootCmd.AddCommand(githubCmd)
}

func runGithubPRs(cmd *cobra.Command, args []string) error {
	repo := ""
	if len(args) > 0 {
		repo = args[0]
	}

	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI not found: %w", err)
	}

	var ghArgs []string
	if repo != "" {
		ghArgs = []string{"pr", "list", "--repo", repo, "--limit", "20"}
	} else {
		ghArgs = []string{"pr", "list", "--limit", "20"}
	}

	ghCmd := exec.Command("gh", ghArgs...)
	ghCmd.Stdout = os.Stdout
	ghCmd.Stderr = os.Stderr
	return ghCmd.Run()
}

func runGithubReview(cmd *cobra.Command, args []string) error {
	prNumber := args[0]

	reviewCmd := exec.Command("gh", "pr", "review", prNumber)
	reviewCmd.Stdout = os.Stdout
	reviewCmd.Stderr = os.Stderr
	reviewCmd.Stdin = os.Stdin
	return reviewCmd.Run()
}

func runGithubComment(cmd *cobra.Command, args []string) error {
	prNumber := args[0]
	comment := args[1]

	ghCmd := exec.Command("gh", "pr", "comment", prNumber, "--body", comment)
	ghCmd.Stdout = os.Stdout
	ghCmd.Stderr = os.Stderr
	return ghCmd.Run()
}

func ghGraphQL(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}

	body, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	return result, nil
}

func parseGitHubURL(input string) (owner, repo string, err error) {
	input = strings.TrimPrefix(input, "https://github.com/")
	input = strings.TrimPrefix(input, "git@github.com:")
	parts := strings.Split(input, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format")
	}
	return parts[0], parts[1], nil
}

func isGitHubURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Host == "github.com" || u.Host == "www.github.com")
}
