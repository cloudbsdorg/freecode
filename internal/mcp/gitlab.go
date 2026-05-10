package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type GitLabClient struct {
	httpClient *http.Client
	token     string
	baseURL   string
}

func NewGitLabClient() *GitLabClient {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		token = os.Getenv("GITLAB_PERSONAL_ACCESS_TOKEN")
	}

	baseURL := "https://gitlab.com/api/v4"
	if override := os.Getenv("GITLAB_API_URL"); override != "" {
		baseURL = strings.TrimSuffix(override, "/")
	}

	return &GitLabClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		baseURL:   baseURL,
	}
}

func (c *GitLabClient) IsConfigured() bool {
	return c.token != ""
}

type GitLabProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	Private           bool   `json:"visibility" json:"visibility"`
	WebURL            string `json:"web_url"`
	StarCount         int    `json:"star_count"`
	ForksCount        int    `json:"forks_count"`
}

type GitLabIssue struct {
	ID          int      `json:"id"`
	IID         int      `json:"iid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	State       string   `json:"state"`
	Author      struct {
		Username string `json:"username"`
	} `json:"author"`
	Labels     []string `json:"labels"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	WebURL     string   `json:"web_url"`
	Weight     int      `json:"weight"`
}

type GitLabMergeRequest struct {
	ID          int    `json:"id"`
	IID         int    `json:"iid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Author      struct {
		Username string `json:"username"`
	} `json:"author"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Additions   int    `json:"changes_count"`
	DiffRefs     struct {
		BaseSHA  string `json:"base_sha"`
		HeadSHA  string `json:"head_sha"`
		StartSHA string `json:"start_sha"`
	} `json:"diff_refs"`
	WebURL     string `json:"web_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GitLabSearchResult struct {
	Projects []GitLabProject `json:"projects,omitempty"`
	Issues   []GitLabIssue   `json:"issues,omitempty"`
	MergeRequests []GitLabMergeRequest `json:"merge_requests,omitempty"`
}

func (c *GitLabClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("GitLab token not configured. Set GITLAB_TOKEN or GITLAB_PERSONAL_ACCESS_TOKEN environment variable")
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

	req.Header.Set("PRIVATE-TOKEN", c.token)
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
		return nil, fmt.Errorf("GitLab API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

func (c *GitLabClient) GetProject(projectID string) (*GitLabProject, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/projects/%s", url.PathEscape(projectID))

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result GitLabProject
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitLabClient) ListIssues(projectID string, state string) ([]GitLabIssue, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/projects/%s/issues", url.PathEscape(projectID))

	params := "?per_page=20"
	if state != "" {
		params += "&state=" + state
	}
	path += params

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result []GitLabIssue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *GitLabClient) GetIssue(projectID string, issueIID int) (*GitLabIssue, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/projects/%s/issues/%d", url.PathEscape(projectID), issueIID)

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result GitLabIssue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitLabClient) CreateIssue(projectID, title, description string, labels []string) (*GitLabIssue, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/projects/%s/issues", url.PathEscape(projectID))

	reqBody := map[string]interface{}{
		"title":       title,
		"description": description,
	}
	if len(labels) > 0 {
		reqBody["labels"] = strings.Join(labels, ",")
	}

	data, err := c.doRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var result GitLabIssue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitLabClient) ListMergeRequests(projectID string, state string) ([]GitLabMergeRequest, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/projects/%s/merge_requests", url.PathEscape(projectID))

	params := "?per_page=20"
	if state != "" {
		params += "&state=" + state
	}
	path += params

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result []GitLabMergeRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *GitLabClient) GetMergeRequest(projectID string, mrIID int) (*GitLabMergeRequest, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/projects/%s/merge_requests/%d", url.PathEscape(projectID), mrIID)

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result GitLabMergeRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitLabClient) CreateMergeRequest(projectID, title, description, sourceBranch, targetBranch string) (*GitLabMergeRequest, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/projects/%s/merge_requests", url.PathEscape(projectID))

	reqBody := map[string]interface{}{
		"title":         title,
		"description":   description,
		"source_branch": sourceBranch,
		"target_branch": targetBranch,
	}

	data, err := c.doRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var result GitLabMergeRequest
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *GitLabClient) SearchProjects(query string) ([]GitLabProject, error) {
	ctx := context.Background()
	path := "/search/projects?q=" + url.QueryEscape(query) + "&per_page=10"

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result []GitLabProject
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *GitLabClient) SearchIssues(projectID, query string) ([]GitLabIssue, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/search/issues?scope=issues&search=%s&per_page=10", url.QueryEscape(query))
	if projectID != "" {
		path += "&project_id=" + url.PathEscape(projectID)
	}

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result []GitLabIssue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *GitLabClient) GetFileContent(projectID, path, ref string) (string, error) {
	ctx := context.Background()
	apiPath := fmt.Sprintf("/projects/%s/repository/files/%s", url.PathEscape(projectID), url.PathEscape(path))
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

func (c *GitLabClient) GitLabExecute(ctx context.Context, input GitLabToolInput) (string, error) {
	if !c.IsConfigured() {
		return "", fmt.Errorf("GitLab not configured: set GITLAB_TOKEN or GITLAB_PERSONAL_ACCESS_TOKEN environment variable")
	}

	switch input.Operation {
	case "get_project":
		project, err := c.GetProject(input.ProjectID)
		if err != nil {
			return "", err
		}
		return formatJSON(project)

	case "list_issues":
		issues, err := c.ListIssues(input.ProjectID, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(issues)

	case "get_issue":
		issue, err := c.GetIssue(input.ProjectID, input.IssueIID)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "create_issue":
		issue, err := c.CreateIssue(input.ProjectID, input.Title, input.Description, input.Labels)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "list_mrs":
		mrs, err := c.ListMergeRequests(input.ProjectID, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(mrs)

	case "get_mr":
		mr, err := c.GetMergeRequest(input.ProjectID, input.IssueIID)
		if err != nil {
			return "", err
		}
		return formatJSON(mr)

	case "create_mr":
		mr, err := c.CreateMergeRequest(input.ProjectID, input.Title, input.Description, input.SourceBranch, input.TargetBranch)
		if err != nil {
			return "", err
		}
		return formatJSON(mr)

	case "search_projects":
		projects, err := c.SearchProjects(input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(projects)

	case "search_issues":
		issues, err := c.SearchIssues(input.ProjectID, input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(issues)

	case "get_file":
		content, err := c.GetFileContent(input.ProjectID, input.Path, input.Ref)
		if err != nil {
			return "", err
		}
		return content, nil

	default:
		return "", fmt.Errorf("unknown operation: %s", input.Operation)
	}
}

type GitLabTool struct {
	client *GitLabClient
}

func NewGitLabTool() *GitLabTool {
	return &GitLabTool{
		client: NewGitLabClient(),
	}
}

func (t *GitLabTool) Name() string {
	return "gitlab"
}

func (t *GitLabTool) Description() string {
	return "GitLab API tools - search projects, manage issues and MRs"
}

func (t *GitLabTool) IsConfigured() bool {
	return t.client.IsConfigured()
}

type GitLabToolInput struct {
	Operation   string   `json:"operation"`
	ProjectID   string   `json:"project_id"`
	IssueIID    int     `json:"issue_iid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	SourceBranch string  `json:"source_branch"`
	TargetBranch string  `json:"target_branch"`
	Query       string   `json:"query"`
	Path        string   `json:"path"`
	Ref         string   `json:"ref"`
	Labels      []string `json:"labels"`
	State       string   `json:"state"`
}

func (t *GitLabTool) Execute(ctx context.Context, input GitLabToolInput) (string, error) {
	if !t.client.IsConfigured() {
		return "", fmt.Errorf("GitLab not configured: set GITLAB_TOKEN or GITLAB_PERSONAL_ACCESS_TOKEN environment variable")
	}

	switch input.Operation {
	case "get_project":
		project, err := t.client.GetProject(input.ProjectID)
		if err != nil {
			return "", err
		}
		return formatJSON(project)

	case "list_issues":
		issues, err := t.client.ListIssues(input.ProjectID, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(issues)

	case "get_issue":
		issue, err := t.client.GetIssue(input.ProjectID, input.IssueIID)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "create_issue":
		issue, err := t.client.CreateIssue(input.ProjectID, input.Title, input.Description, input.Labels)
		if err != nil {
			return "", err
		}
		return formatJSON(issue)

	case "list_mrs":
		mrs, err := t.client.ListMergeRequests(input.ProjectID, input.State)
		if err != nil {
			return "", err
		}
		return formatJSON(mrs)

	case "get_mr":
		mr, err := t.client.GetMergeRequest(input.ProjectID, input.IssueIID)
		if err != nil {
			return "", err
		}
		return formatJSON(mr)

	case "create_mr":
		mr, err := t.client.CreateMergeRequest(input.ProjectID, input.Title, input.Description, input.SourceBranch, input.TargetBranch)
		if err != nil {
			return "", err
		}
		return formatJSON(mr)

	case "search_projects":
		projects, err := t.client.SearchProjects(input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(projects)

	case "search_issues":
		issues, err := t.client.SearchIssues(input.ProjectID, input.Query)
		if err != nil {
			return "", err
		}
		return formatJSON(issues)

	case "get_file":
		content, err := t.client.GetFileContent(input.ProjectID, input.Path, input.Ref)
		if err != nil {
			return "", err
		}
		return content, nil

	default:
		return "", fmt.Errorf("unknown operation: %s", input.Operation)
	}
}