package tool

import (
	"context"
	"fmt"
)

type SkillTool struct{}

func NewSkillTool() *SkillTool {
	return &SkillTool{}
}

func (t *SkillTool) Name() string {
	return "skill"
}

func (t *SkillTool) Description() string {
	return "List and invoke skills"
}

func (t *SkillTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "skill",
		Description: "List and invoke skills",
		Parameters: map[string]Parameter{
			"action": {
				Type:        "string",
				Description: "Action: list, invoke",
				Required:    true,
				Enum:        []string{"list", "invoke"},
			},
			"name": {
				Type:        "string",
				Description: "Skill name",
			},
			"args": {
				Type:        "object",
				Description: "Skill arguments",
			},
		},
	}
}

func (t *SkillTool) Execute(ctx context.Context, req Request) (*Response, error) {
	action, ok := req.Arguments["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action must be a string")
	}

	switch action {
	case "list":
		skills := []string{
			"git-master: Atomic commits, rebase, branch management",
			"playwright: Browser automation",
			"frontend-ui-ux: Design-first UI/UX development",
			"review-work: 5 parallel subagent code review",
			"ai-slop-remover: Remove AI-generated code smells",
		}
		return &Response{Result: fmt.Sprintf("Available skills:\n%s", join(skills, "\n"))}, nil
	case "invoke":
		name, _ := req.Arguments["name"].(string)
		return &Response{Result: fmt.Sprintf("Invoking skill: %s (placeholder)", name)}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func join(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
