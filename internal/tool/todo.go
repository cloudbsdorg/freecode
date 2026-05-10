package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type TodoTool struct {
	mu    sync.RWMutex
	todos map[string]TodoItem
}

type TodoItem struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Status   string                 `json:"status"`
	Priority string                 `json:"priority"`
	Metadata map[string]interface{} `json:"metadata"`
}

func init() {
	Register("todo", func() Tool { return &TodoTool{todos: make(map[string]TodoItem)} })
}

func NewTodoTool() *TodoTool {
	return &TodoTool{
		todos: make(map[string]TodoItem),
	}
}

func (t *TodoTool) Name() string {
	return "todo"
}

func (t *TodoTool) Description() string {
	return "Manage todo items"
}

func (t *TodoTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "todo",
		Description: "Manage todo items",
		Parameters: map[string]Parameter{
			"action": {
				Type:        "string",
				Description: "Action to perform: list, add, update, delete",
				Required:    true,
				Enum:        []string{"list", "add", "update", "delete"},
			},
			"content": {
				Type:        "string",
				Description: "Todo content",
			},
			"id": {
				Type:        "string",
				Description: "Todo ID",
			},
			"status": {
				Type:        "string",
				Description: "Todo status",
				Enum:        []string{"pending", "in_progress", "completed", "cancelled"},
			},
			"priority": {
				Type:        "string",
				Description: "Todo priority",
				Enum:        []string{"high", "medium", "low"},
			},
		},
	}
}

func (t *TodoTool) Execute(ctx context.Context, req Request) (*Response, error) {
	action, ok := req.Arguments["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action must be a string")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	switch action {
	case "list":
		todos := make([]TodoItem, 0, len(t.todos))
		for _, todo := range t.todos {
			todos = append(todos, todo)
		}
		data, _ := json.Marshal(todos)
		return &Response{Result: string(data)}, nil

	case "add":
		content, ok := req.Arguments["content"].(string)
		if !ok {
			return nil, fmt.Errorf("content must be a string")
		}
		id := fmt.Sprintf("todo-%d", len(t.todos)+1)
		priority := "medium"
		if p, ok := req.Arguments["priority"].(string); ok {
			priority = p
		}
		t.todos[id] = TodoItem{
			ID:       id,
			Content:  content,
			Status:   "pending",
			Priority: priority,
		}
		return &Response{Result: fmt.Sprintf("Added: %s", id)}, nil

	case "update":
		id, ok := req.Arguments["id"].(string)
		if !ok {
			return nil, fmt.Errorf("id must be a string")
		}
		todo, exists := t.todos[id]
		if !exists {
			return nil, fmt.Errorf("todo not found: %s", id)
		}
		if status, ok := req.Arguments["status"].(string); ok {
			todo.Status = status
		}
		if priority, ok := req.Arguments["priority"].(string); ok {
			todo.Priority = priority
		}
		t.todos[id] = todo
		return &Response{Result: fmt.Sprintf("Updated: %s", id)}, nil

	case "delete":
		id, ok := req.Arguments["id"].(string)
		if !ok {
			return nil, fmt.Errorf("id must be a string")
		}
		delete(t.todos, id)
		return &Response{Result: fmt.Sprintf("Deleted: %s", id)}, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (t *TodoTool) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &t.todos)
}

func (t *TodoTool) SaveToFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t.todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
