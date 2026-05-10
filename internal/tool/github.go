package tool

import (
	"context"

	"github.com/freecode/freecode/internal/mcp"
)

type GitHubTool struct {
	client *mcp.GitHubClient
}

func init() {
	Register("github", func() Tool { return &GitHubTool{client: mcp.NewGitHubClient()} })
}

func NewGitHubTool() *GitHubTool {
	return &GitHubTool{
		client: mcp.NewGitHubClient(),
	}
}

func (t *GitHubTool) Name() string {
	return "github"
}

func (t *GitHubTool) Description() string {
	return "GitHub API - search repos, manage issues and pull requests. Requires GITHUB_TOKEN or GITHUB_PAT env var."
}

func (t *GitHubTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "github",
		Description: "GitHub API - search repos, manage issues and pull requests. Requires GITHUB_TOKEN or GITHUB_PAT env var.",
		Parameters: map[string]Parameter{
			"operation": {
				Type:        "string",
				Description: "Operation: get_repo, list_issues, get_issue, create_issue, list_prs, get_pr, create_pr, search_repos, search_code, get_file",
				Required:    true,
				Enum: []string{
					"get_repo", "list_issues", "get_issue", "create_issue",
					"list_prs", "get_pr", "create_pr",
					"search_repos", "search_code", "get_file",
				},
			},
			"owner": {
				Type:        "string",
				Description: "Repository owner (user or organization)",
				Required:    false,
			},
			"repo": {
				Type:        "string",
				Description: "Repository name",
				Required:    false,
			},
			"number": {
				Type:        "integer",
				Description: "Issue or PR number",
				Required:    false,
			},
			"title": {
				Type:        "string",
				Description: "Title for create_issue or create_pr",
				Required:    false,
			},
			"body": {
				Type:        "string",
				Description: "Body/description for create_issue or create_pr",
				Required:    false,
			},
			"head": {
				Type:        "string",
				Description: "Source branch for create_pr",
				Required:    false,
			},
			"base": {
				Type:        "string",
				Description: "Target branch for create_pr",
				Required:    false,
			},
			"query": {
				Type:        "string",
				Description: "Search query for search_repos or search_code",
				Required:    false,
			},
			"path": {
				Type:        "string",
				Description: "File path for get_file",
				Required:    false,
			},
			"ref": {
				Type:        "string",
				Description: "Git ref (branch, tag, commit) for get_file",
				Required:    false,
			},
			"state": {
				Type:        "string",
				Description: "State filter: open, closed, all (for list_issues, list_prs)",
				Required:    false,
			},
		},
	}
}

func (t *GitHubTool) Execute(ctx context.Context, req Request) (*Response, error) {
	if !t.client.IsConfigured() {
		return &Response{
			Result: "GitHub not configured: set GITHUB_TOKEN or GITHUB_PAT environment variable",
		}, nil
	}

	input := mcp.GitHubToolInput{
		Operation: getStringArg(req.Arguments, "operation"),
		Owner:     getStringArg(req.Arguments, "owner"),
		Repo:      getStringArg(req.Arguments, "repo"),
		Title:     getStringArg(req.Arguments, "title"),
		Body:      getStringArg(req.Arguments, "body"),
		Head:      getStringArg(req.Arguments, "head"),
		Base:      getStringArg(req.Arguments, "base"),
		Query:     getStringArg(req.Arguments, "query"),
		Path:      getStringArg(req.Arguments, "path"),
		Ref:       getStringArg(req.Arguments, "ref"),
		State:     getStringArg(req.Arguments, "state"),
	}

	if n, ok := req.Arguments["number"].(float64); ok {
		input.Number = int(n)
	}
	if labels, ok := req.Arguments["labels"].([]any); ok {
		input.Labels = make([]string, len(labels))
		for i, l := range labels {
			if s, ok := l.(string); ok {
				input.Labels[i] = s
			}
		}
	}

	result, err := t.client.GitHubExecute(ctx, input)
	if err != nil {
		return &Response{Error: err}, nil
	}
	return &Response{Result: result}, nil
}

func getStringArg(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}