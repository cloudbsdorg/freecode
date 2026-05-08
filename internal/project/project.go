package project

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Project represents a detected software project
type Project struct {
	ID       string
	Name     string
	Path     string
	Remote   string
	VCS      string
	Created  int64
	Modified int64
}

// Detector interface for project detection
type Detector interface {
	Detect(ctx context.Context, dir string) (*Project, error)
}

type projectDetector struct{}

// NewDetector creates a new project detector
func NewDetector() *projectDetector {
	return &projectDetector{}
}

func (d *projectDetector) Detect(ctx context.Context, dir string) (*Project, error) {
	root := DetectProjectRoot(dir)
	if root == "" {
		return nil, nil
	}

	if proj := detectGit(root); proj != nil {
		return proj, nil
	}
	if proj := detectNPM(root); proj != nil {
		return proj, nil
	}
	if proj := detectGoMod(root); proj != nil {
		return proj, nil
	}
	if proj := detectCargo(root); proj != nil {
		return proj, nil
	}
	if proj := detectPython(root); proj != nil {
		return proj, nil
	}

	return nil, nil
}

func DetectProjectRoot(dir string) string {
	if dir == "" {
		return ""
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}

	markers := []string{".git", "package.json", "go.mod", "Cargo.toml", "requirements.txt", "pyproject.toml"}

	current := absDir
	for {
		for _, marker := range markers {
			if _, err := os.Stat(filepath.Join(current, marker)); err == nil {
				return current
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func detectGit(dir string) *Project {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return &Project{
			ID:   dir,
			Path: dir,
			VCS:  "git",
		}
	}
	return nil
}

func detectNPM(dir string) *Project {
	pkgFile := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(pkgFile)
	if err != nil {
		return nil
	}

	var pkg struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}

	name := pkg.Name
	if name == "" {
		name = filepath.Base(dir)
	}

	return &Project{
		ID:   dir,
		Name: name,
		Path: dir,
		VCS:  "npm",
	}
}

func detectGoMod(dir string) *Project {
	modFile := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(modFile)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			name := strings.TrimPrefix(line, "module ")
			return &Project{
				ID:   dir,
				Name: name,
				Path: dir,
				VCS:  "go",
			}
		}
	}

	return nil
}

func detectCargo(dir string) *Project {
	cargoFile := filepath.Join(dir, "Cargo.toml")
	data, err := os.ReadFile(cargoFile)
	if err != nil {
		return nil
	}

	name := parseCargoName(string(data))
	if name == "" {
		name = filepath.Base(dir)
	}

	return &Project{
		ID:   dir,
		Name: name,
		Path: dir,
		VCS:  "cargo",
	}
}

func parseCargoName(content string) string {
	lines := strings.Split(content, "\n")
	inPackage := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "[package]" {
			inPackage = true
			continue
		}

		if strings.HasPrefix(line, "[") && line != "[package]" {
			inPackage = false
			continue
		}

		if inPackage && strings.HasPrefix(line, "name") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				name := strings.Trim(parts[1], " \t\"")
				return name
			}
		}
	}

	return ""
}

func detectPython(dir string) *Project {
	pyprojectFile := filepath.Join(dir, "pyproject.toml")
	if data, err := os.ReadFile(pyprojectFile); err == nil {
		name := parsePyprojectName(string(data))
		if name != "" {
			return &Project{
				ID:   dir,
				Name: name,
				Path: dir,
				VCS:  "python",
			}
		}
	}

	reqFile := filepath.Join(dir, "requirements.txt")
	if _, err := os.Stat(reqFile); err == nil {
		return &Project{
			ID:   dir,
			Name: filepath.Base(dir),
			Path: dir,
			VCS:  "python",
		}
	}

	return nil
}

func parsePyprojectName(content string) string {
	lines := strings.Split(content, "\n")
	inProject := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "[project]" {
			inProject = true
			continue
		}

		if strings.HasPrefix(line, "[") && line != "[project]" {
			inProject = false
			continue
		}

		if inProject && strings.HasPrefix(line, "name") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				name := strings.Trim(parts[1], " \t\"")
				return name
			}
		}
	}

	return ""
}

func List(ctx context.Context, dirs []string) ([]*Project, error) {
	var projects []*Project
	for _, dir := range dirs {
		detector := NewDetector()
		if proj, err := detector.Detect(ctx, dir); err == nil && proj != nil {
			projects = append(projects, proj)
		}
	}
	return projects, nil
}