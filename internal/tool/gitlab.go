package tool

import (
	"context"

	"github.com/freecode/freecode/internal/mcp"
)

type GitLabTool struct {
	client *mcp.GitLabClient
}

func init() {
	Register("gitlab", func() Tool { return &GitLabTool{client: mcp.NewGitLabClient()} })
}

func NewGitLabTool() *GitLabTool {
	return &GitLabTool{
		client: mcp.NewGitLabClient(),
	}
}

func (t *GitLabTool) Name() string {
	return "gitlab"
}

func (t *GitLabTool) Description() string {
	return "GitLab API - search projects, manage issues and merge requests. Requires GITLAB_TOKEN or GITLAB_PERSONAL_ACCESS_TOKEN env var."
}

func (t *GitLabTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "gitlab",
		Description: "GitLab API - search projects, manage issues and merge requests. Requires GITLAB_TOKEN or GITLAB_PERSONAL_ACCESS_TOKEN env var.",
		Parameters: map[string]Parameter{
			"operation": {
				Type:        "string",
				Description: "Operation: get_project, list_issues, get_issue, create_issue, list_mrs, get_mr, create_mr, search_projects, search_issues, get_file",
				Required:    true,
				Enum: []string{
					"get_project", "list_issues", "get_issue", "create_issue",
					"list_mrs", "get_mr", "create_mr",
					"search_projects", "search_issues", "get_file",
				},
			},
			"project_id": {
				Type:        "string",
				Description: "Project ID or path with namespace (e.g., 'group/project')",
				Required:    false,
			},
			"issue_iid": {
				Type:        "integer",
				Description: "Issue or MR IID number",
				Required:    false,
			},
			"title": {
				Type:        "string",
				Description: "Title for create_issue or create_mr",
				Required:    false,
			},
			"description": {
				Type:        "string",
				Description: "Description for create_issue or create_mr",
				Required:    false,
			},
			"source_branch": {
				Type:        "string",
				Description: "Source branch for create_mr",
				Required:    false,
			},
			"target_branch": {
				Type:        "string",
				Description: "Target branch for create_mr",
				Required:    false,
			},
			"query": {
				Type:        "string",
				Description: "Search query for search_projects or search_issues",
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
				Description: "State filter: opened, closed, all (for list_issues, list_mrs)",
				Required:    false,
			},
			"labels": {
				Type:        "array",
				Description: "Labels for create_issue",
				Required:    false,
				Items: &Parameter{
					Type: "string",
				},
			},
		},
	}
}

func (t *GitLabTool) Execute(ctx context.Context, req Request) (*Response, error) {
	if !t.client.IsConfigured() {
		return &Response{
			Result: "GitLab not configured: set GITLAB_TOKEN or GITLAB_PERSONAL_ACCESS_TOKEN environment variable",
		}, nil
	}

	input := mcp.GitLabToolInput{
		Operation:    getStringArg(req.Arguments, "operation"),
		ProjectID:    getStringArg(req.Arguments, "project_id"),
		Title:        getStringArg(req.Arguments, "title"),
		Description:  getStringArg(req.Arguments, "description"),
		SourceBranch: getStringArg(req.Arguments, "source_branch"),
		TargetBranch: getStringArg(req.Arguments, "target_branch"),
		Query:        getStringArg(req.Arguments, "query"),
		Path:         getStringArg(req.Arguments, "path"),
		Ref:          getStringArg(req.Arguments, "ref"),
		State:        getStringArg(req.Arguments, "state"),
	}

	if iid, ok := req.Arguments["issue_iid"].(float64); ok {
		input.IssueIID = int(iid)
	}
	if labels, ok := req.Arguments["labels"].([]any); ok {
		input.Labels = make([]string, len(labels))
		for i, l := range labels {
			if s, ok := l.(string); ok {
				input.Labels[i] = s
			}
		}
	}

	result, err := t.client.GitLabExecute(ctx, input)
	if err != nil {
		return &Response{Error: err}, nil
	}
	return &Response{Result: result}, nil
}