package mcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type GitHubClient struct {
	httpClient *http.Client
	token      string
	baseURL   string
}

func NewGitHubClient() *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_PAT")
	}

	baseURL := "https://api.github.com"
	if override := os.Getenv("GITHUB_API_URL"); override != "" {
		baseURL = strings.TrimSuffix(override, "/")
	}

	return &GitHubClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		baseURL:   baseURL,
	}
}

func (c *GitHubClient) IsConfigured() bool {
	return c.token != ""
}

type GitHubRepo struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Private  bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
}

type GitHubIssue struct {
	Number      int    `json:"number"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	State      string `json:"state"`
	Author     string `json:"user"`
	Labels     []string `json:"labels"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Comments   int    `json:"comments"`
	URL        string `json:"html_url"`
}

type GitHubPullRequest struct {
	Number      int    `json:"number"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	State      string `json:"state"`
	Author     string `json:"user"`
	Head       string `json:"head"`
	Base       string `json:"base"`
	Additions  int    `json:"additions"`
	Deletions  int    `json:"deletions"`
	Commits    int    `json:"commits"`
	URL        string `json:"html_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type GitHubSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []GitHubSearchItem `json:"items"`
}

type GitHubSearchItem struct {
	FullName      string `json:"full_name"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Owner         string `json:"owner"`
	Private       bool   `json:"private"`
	Stars         int    `json:"stargazers_count"`
	Forks         int    `json:"forks_count"`
	Language      string `json:"language"`
	URL           string `json:"html_url"`
}

func (c *GitHubClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("GitHub token not configured. Set GITHUB_TOKEN or GITHUB_PAT environment variable")
	}

	reqURL := c.baseURL + path
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

func (c *GitHubClient) GetRepo(owner, repo string) (*GitHubRepo, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/repos/%s/%s", owner, repo)

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result GitHubRepo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitHubClient) ListIssues(owner, repo string, state string) ([]GitHubIssue, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/repos/%s/%s/issues", owner, repo)

	if state == "" {
		state = "open"
	}
	path += "?state=" + state + "&per_page=20"

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result []GitHubIssue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *GitHubClient) GetIssue(owner, repo string, number int) (*GitHubIssue, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number)

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result GitHubIssue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitHubClient) CreateIssue(owner, repo, title, body string, labels []string) (*GitHubIssue, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/repos/%s/%s/issues", owner, repo)

	reqBody := map[string]interface{}{
		"title": title,
		"body":  body,
	}
	if len(labels) > 0 {
		reqBody["labels"] = labels
	}

	data, err := c.doRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var result GitHubIssue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitHubClient) ListPullRequests(owner, repo string, state string) ([]GitHubPullRequest, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)

	if state == "" {
		state = "open"
	}
	path += "?state=" + state + "&per_page=20"

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result []GitHubPullRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *GitHubClient) GetPullRequest(owner, repo string, number int) (*GitHubPullRequest, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number)

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result GitHubPullRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitHubClient) CreatePullRequest(owner, repo, title, body, head, base string) (*GitHubPullRequest, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)

	reqBody := map[string]interface{}{
		"title": title,
		"body":  body,
		"head":  head,
		"base":  base,
	}

	data, err := c.doRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var result GitHubPullRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitHubClient) SearchRepositories(query string) ([]GitHubSearchItem, error) {
	ctx := context.Background()
	path := "/search/repositories?q=" + url.QueryEscape(query) + "&per_page=10"

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result GitHubSearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Items, nil
}

func (c *GitHubClient) SearchCode(query string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	path := "/search/code?q=" + url.QueryEscape(query) + "&per_page=10"

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		TotalCount int                   `json:"total_count"`
		Items      []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Items, nil
}

func (c *GitHubClient) GetFileContent(owner, repo, path, ref string) (string, error) {
	ctx := context.Background()
	apiPath := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	if ref != "" {
		apiPath += "?ref=" + ref
	}

	data, err := c.doRequest(ctx, "GET", apiPath, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Content string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Encoding == "base64" {
		content, err := decodeBase64(result.Content)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 content: %w", err)
		}
		return content, nil
	}
	return result.Content, nil
}

func (c *GitHubClient) GitHubExecute(ctx context.Context, input GitHubToolInput) (string, error) {
	if !c.IsConfigured() {
		return "", fmt.Errorf("GitHub not configured: set GITHUB_TOKEN or GITHUB_PAT environment variable")
	}

	switch input.Operation {
	case "get_repo":
		repo, err := c.GetRepo(input.Owner, input.Repo)
		if err != nil {
			return "", err
		}
		return formatJSON(repo)

	case "list_issues":
		issues, err := c.ListIssues(input.Owner, input.Repo, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(issues)

	case "get_issue":
		issue, err := c.GetIssue(input.Owner, input.Repo, input.Number)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "create_issue":
		issue, err := c.CreateIssue(input.Owner, input.Repo, input.Title, input.Body, input.Labels)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "list_prs":
		prs, err := c.ListPullRequests(input.Owner, input.Repo, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(prs)

	case "get_pr":
		pr, err := c.GetPullRequest(input.Owner, input.Repo, input.Number)
		if err != nil {
			return "", err
		}
		return formatJSON(pr)

	case "create_pr":
		pr, err := c.CreatePullRequest(input.Owner, input.Repo, input.Title, input.Body, input.Head, input.Base)
		if err != nil {
			return "", err
		}
		return formatJSON(pr)

	case "search_repos":
		repos, err := c.SearchRepositories(input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(repos)

	case "search_code":
		results, err := c.SearchCode(input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(results)

	case "get_file":
		content, err := c.GetFileContent(input.Owner, input.Repo, input.Path, input.Ref)
		if err != nil {
			return "", err
		}
		return content, nil

	default:
		return "", fmt.Errorf("unknown operation: %s", input.Operation)
	}
}

func decodeBase64(s string) (string, error) {
	s = strings.TrimSpace(s)
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(dst, []byte(s))
	if err != nil {
		return "", err
	}
	return string(dst[:n]), nil
}

type GitHubTool struct {
	client *GitHubClient
}

func NewGitHubTool() *GitHubTool {
	return &GitHubTool{
		client: NewGitHubClient(),
	}
}

func (t *GitHubTool) Name() string {
	return "github"
}

func (t *GitHubTool) Description() string {
	return "GitHub API tools - search repos, manage issues and PRs"
}

func (t *GitHubTool) IsConfigured() bool {
	return t.client.IsConfigured()
}

type GitHubToolInput struct {
	Operation string                 `json:"operation"`
	Owner     string                `json:"owner"`
	Repo      string                `json:"repo"`
	Number    int                   `json:"number"`
	Title     string                `json:"title"`
	Body      string                `json:"body"`
	Head      string                `json:"head"`
	Base      string                `json:"base"`
	Query     string                `json:"query"`
	Path      string                `json:"path"`
	Ref       string                `json:"ref"`
	Labels    []string              `json:"labels"`
	State     string                `json:"state"`
}

func (t *GitHubTool) Execute(ctx context.Context, input GitHubToolInput) (string, error) {
	if !t.client.IsConfigured() {
		return "", fmt.Errorf("GitHub not configured: set GITHUB_TOKEN or GITHUB_PAT environment variable")
	}

	switch input.Operation {
	case "get_repo":
		repo, err := t.client.GetRepo(input.Owner, input.Repo)
		if err != nil {
			return "", err
		}
		return formatJSON(repo)

	case "list_issues":
		issues, err := t.client.ListIssues(input.Owner, input.Repo, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(issues)

	case "get_issue":
		issue, err := t.client.GetIssue(input.Owner, input.Repo, input.Number)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "create_issue":
		issue, err := t.client.CreateIssue(input.Owner, input.Repo, input.Title, input.Body, input.Labels)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "list_prs":
		prs, err := t.client.ListPullRequests(input.Owner, input.Repo, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(prs)

	case "get_pr":
		pr, err := t.client.GetPullRequest(input.Owner, input.Repo, input.Number)
		if err != nil {
			return "", err
		}
		return formatJSON(pr)

	case "create_pr":
		pr, err := t.client.CreatePullRequest(input.Owner, input.Repo, input.Title, input.Body, input.Head, input.Base)
		if err != nil {
			return "", err
		}
		return formatJSON(pr)

	case "search_repos":
		repos, err := t.client.SearchRepositories(input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(repos)

	case "search_code":
		results, err := t.client.SearchCode(input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(results)

	case "get_file":
		content, err := t.client.GetFileContent(input.Owner, input.Repo, input.Path, input.Ref)
		if err != nil {
			return "", err
		}
		return content, nil

	default:
		return "", fmt.Errorf("unknown operation: %s", input.Operation)
	}
}

func formatJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}