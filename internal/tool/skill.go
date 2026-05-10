package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/freecode/freecode/internal/skill"
)

type SkillTool struct {
	homeDir string
}

func init() {
	Register("skill", func() Tool { return &SkillTool{} })
}

func NewSkillTool() *SkillTool {
	home, _ := os.UserHomeDir()
	return &SkillTool{homeDir: home}
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
				Description: "Skill name to invoke",
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
		return t.listSkills()
	case "invoke":
		name, ok := req.Arguments["name"].(string)
		if !ok {
			return nil, fmt.Errorf("name must be a string")
		}
		return t.invokeSkill(name)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (t *SkillTool) listSkills() (*Response, error) {
	skills, err := skill.Discover(t.homeDir)
	if err != nil {
		return &Response{
			Result: fmt.Sprintf("Error discovering skills: %v", err),
		}, nil
	}

	if len(skills) == 0 {
		return &Response{
			Result: "No skills found. Create skill files at:\n" +
				"  ~/.config/freecode/skills/<name>/SKILL.md\n" +
				"  ~/.config/claude/skills/<name>/SKILL.md\n" +
				"  ~/.agents/skills/<name>/SKILL.md\n\n" +
				"Skill files should have YAML frontmatter with 'name' and 'description'.",
		}, nil
	}

	var lines []string
	lines = append(lines, "Available skills:")
	for _, s := range skills {
		lines = append(lines, fmt.Sprintf("  %s: %s", s.Name, s.Description))
	}

	return &Response{Result: join(lines, "\n")}, nil
}

func (t *SkillTool) invokeSkill(name string) (*Response, error) {
	skills, err := skill.Discover(t.homeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to discover skills: %w", err)
	}

	s := skill.GetSkill(skills, name)
	if s == nil {
		available := []string{}
		for _, sk := range skills {
			available = append(available, sk.Name)
		}
		return nil, fmt.Errorf("skill not found: %s. Available: %v", name, available)
	}

	baseDir := filepath.Dir(s.Location)

	output := fmt.Sprintf(`<skill_content name="%s">
# Skill: %s

%s

Base directory for this skill: %s
Relative paths in this skill (e.g., scripts/, reference/) are relative to this base directory.

</skill_content>`, s.Name, s.Name, s.Content, baseDir)

	return &Response{Result: output}, nil
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