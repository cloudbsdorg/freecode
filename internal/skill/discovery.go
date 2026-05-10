package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type SkillInfo struct {
	Name        string
	Description string
	Location    string
	Content     string
}

type Frontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

var DefaultSkillDirs = []string{
	".config/freecode/skills",
	".config/claude/skills",
	".agents/skills",
}

func Discover(homeDir string) ([]SkillInfo, error) {
	var skills []SkillInfo

	for _, dir := range DefaultSkillDirs {
		skillDir := filepath.Join(homeDir, dir)
		if err := scanDir(skillDir, &skills); err != nil {
			continue
		}
	}

	return skills, nil
}

func DiscoverWithPaths(homeDir string, extraPaths []string) ([]SkillInfo, error) {
	skills, err := Discover(homeDir)
	if err != nil {
		return nil, err
	}

	for _, p := range extraPaths {
		expanded := p
		if strings.HasPrefix(p, "~/") {
			expanded = filepath.Join(homeDir, p[2:])
		}
		if err := scanDir(expanded, &skills); err != nil {
			continue
		}
	}

	return skills, nil
}

func scanDir(dir string, skills *[]SkillInfo) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			continue
		}

		skill, err := LoadSkill(skillPath)
		if err != nil {
			continue
		}

		*skills = append(*skills, *skill)
	}

	return nil
}

func LoadSkill(path string) (*SkillInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill file: %w", err)
	}

	content := string(data)
	frontmatter, body, err := parseFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	if frontmatter.Name == "" {
		name := filepath.Base(filepath.Dir(path))
		frontmatter.Name = name
	}

	if frontmatter.Description == "" {
		frontmatter.Description = "No description"
	}

	return &SkillInfo{
		Name:        frontmatter.Name,
		Description: frontmatter.Description,
		Location:    path,
		Content:     strings.TrimSpace(body),
	}, nil
}

func parseFrontmatter(content string) (*Frontmatter, string, error) {
	if !strings.HasPrefix(content, "---") {
		return &Frontmatter{}, content, nil
	}

	parts := strings.SplitN(content[4:], "---", 2)
	if len(parts) < 2 {
		return &Frontmatter{}, content, nil
	}

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(parts[0]), &fm); err != nil {
		return &Frontmatter{}, content, nil
	}

	return &fm, strings.TrimSpace(parts[1]), nil
}

func GetSkill(skills []SkillInfo, name string) *SkillInfo {
	for i := range skills {
		if skills[i].Name == name {
			return &skills[i]
		}
	}
	return nil
}