package tool

import (
	"context"
	"fmt"
)

type PlanTool struct{}

func init() {
	Register("plan", func() Tool { return &PlanTool{} })
}

func NewPlanTool() *PlanTool {
	return &PlanTool{}
}

func (t *PlanTool) Name() string {
	return "plan"
}

func (t *PlanTool) Description() string {
	return "Create and manage plans"
}

func (t *PlanTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "plan",
		Description: "Create and manage plans",
		Parameters: map[string]Parameter{
			"action": {
				Type:        "string",
				Description: "Action: create, update, execute, list",
				Required:    true,
				Enum:        []string{"create", "update", "execute", "list"},
			},
			"plan_id": {
				Type:        "string",
				Description: "Plan ID",
			},
			"content": {
				Type:        "string",
				Description: "Plan content",
			},
		},
	}
}

func (t *PlanTool) Execute(ctx context.Context, req Request) (*Response, error) {
	action, ok := req.Arguments["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action must be a string")
	}

	switch action {
	case "list":
		return &Response{Result: "Plans:\n- Task breakdown plan\n- Implementation plan"}, nil
	case "create":
		content, _ := req.Arguments["content"].(string)
		return &Response{Result: fmt.Sprintf("Created plan: %s", content)}, nil
	case "execute":
		planID, _ := req.Arguments["plan_id"].(string)
		return &Response{Result: fmt.Sprintf("Executing plan: %s", planID)}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
