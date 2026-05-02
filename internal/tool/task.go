package tool

import (
	"context"
	"fmt"
)

type TaskTool struct{}

func NewTaskTool() *TaskTool {
	return &TaskTool{}
}

func (t *TaskTool) Name() string {
	return "task"
}

func (t *TaskTool) Description() string {
	return "Create and manage tasks"
}

func (t *TaskTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "task",
		Description: "Create and manage tasks",
		Parameters: map[string]Parameter{
			"action": {
				Type:        "string",
				Description: "Action: list, create, update, delete",
				Required:    true,
				Enum:       []string{"list", "create", "update", "delete"},
			},
			"title": {
				Type:        "string",
				Description: "Task title",
			},
			"description": {
				Type:        "string",
				Description: "Task description",
			},
			"task_id": {
				Type:        "string",
				Description: "Task ID",
			},
			"status": {
				Type:        "string",
				Description: "Task status",
				Enum:       []string{"pending", "in_progress", "completed", "blocked"},
			},
		},
	}
}

func (t *TaskTool) Execute(ctx context.Context, req Request) (*Response, error) {
	action, ok := req.Arguments["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action must be a string")
	}

	switch action {
	case "list":
		return &Response{Result: "Task list: (placeholder)"}, nil
	case "create":
		title, _ := req.Arguments["title"].(string)
		return &Response{Result: fmt.Sprintf("Created task: %s", title)}, nil
	case "update":
		taskID, _ := req.Arguments["task_id"].(string)
		return &Response{Result: fmt.Sprintf("Updated task: %s", taskID)}, nil
	case "delete":
		taskID, _ := req.Arguments["task_id"].(string)
		return &Response{Result: fmt.Sprintf("Deleted task: %s", taskID)}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
