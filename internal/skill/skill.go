package skill

import (
	"context"
	"os"
	"path/filepath"
)

type Skill struct {
	Name        string
	Description string
	Path        string
	Category    string
}

type Discovery interface {
	Discover(ctx context.Context, paths []string) ([]*Skill, error)
	GetSkill(name string) (*Skill, error)
}

type fileDiscovery struct {
	skills map[string]*Skill
}

func NewFileDiscovery() Discovery {
	return &fileDiscovery{skills: make(map[string]*Skill)}
}

func (d *fileDiscovery) Discover(ctx context.Context, paths []string) ([]*Skill, error) {
	var result []*Skill
	for _, basePath := range paths {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			skillPath := filepath.Join(basePath, entry.Name())
			skillFile := filepath.Join(skillPath, "SKILL.md")
			if _, err := os.Stat(skillFile); err == nil {
				skill := &Skill{
					Name: entry.Name(),
					Path: skillPath,
				}
				d.skills[skill.Name] = skill
				result = append(result, skill)
			}
		}
	}
	return result, nil
}

func (d *fileDiscovery) GetSkill(name string) (*Skill, error) {
	if s, ok := d.skills[name]; ok {
		return s, nil
	}
	return nil, nil
}
