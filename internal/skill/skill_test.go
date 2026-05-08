package skill

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileDiscovery(t *testing.T) {
	d := NewFileDiscovery()
	if d == nil {
		t.Fatal("NewFileDiscovery() returned nil")
	}
}

func TestFileDiscoveryDiscover(t *testing.T) {
	d := NewFileDiscovery()
	ctx := context.Background()

	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("Failed to create skill dir: %v", err)
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillFile, []byte("# Test Skill"), 0644); err != nil {
		t.Fatalf("Failed to create SKILL.md: %v", err)
	}

	skills, err := d.Discover(ctx, []string{tmpDir})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(skills) != 1 {
		t.Errorf("Discover() returned %d skills, want 1", len(skills))
	}
}

func TestFileDiscoveryGetSkill(t *testing.T) {
	d := NewFileDiscovery()
	ctx := context.Background()

	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("Failed to create skill dir: %v", err)
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillFile, []byte("# My Skill"), 0644); err != nil {
		t.Fatalf("Failed to create SKILL.md: %v", err)
	}

	d.Discover(ctx, []string{tmpDir})

	skill, err := d.GetSkill("my-skill")
	if err != nil {
		t.Fatalf("GetSkill() error = %v", err)
	}
	if skill == nil {
		t.Fatal("GetSkill() returned nil")
	}
	if skill.Name != "my-skill" {
		t.Errorf("skill.Name = %q, want %q", skill.Name, "my-skill")
	}

	nonexistent, err := d.GetSkill("nonexistent")
	if err != nil {
		t.Fatalf("GetSkill() for nonexistent returned error: %v", err)
	}
	if nonexistent != nil {
		t.Error("GetSkill() for nonexistent should return nil")
	}
}

func TestFileDiscoveryDiscoverNoSkills(t *testing.T) {
	d := NewFileDiscovery()
	ctx := context.Background()

	tmpDir := t.TempDir()

	skills, err := d.Discover(ctx, []string{tmpDir})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("Discover() returned %d skills, want 0", len(skills))
	}
}

func TestFileDiscoveryDiscoverNonexistentPath(t *testing.T) {
	d := NewFileDiscovery()
	ctx := context.Background()

	skills, err := d.Discover(ctx, []string{"/nonexistent/path"})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("Discover() returned %d skills, want 0", len(skills))
	}
}

func TestSkill(t *testing.T) {
	skill := &Skill{
		Name:        "test-skill",
		Description: "A test skill",
		Path:        "/path/to/skill",
		Category:    "testing",
	}

	if skill.Name != "test-skill" {
		t.Errorf("Name = %q, want %q", skill.Name, "test-skill")
	}
	if skill.Description != "A test skill" {
		t.Errorf("Description = %q, want %q", skill.Description, "A test skill")
	}
	if skill.Path != "/path/to/skill" {
		t.Errorf("Path = %q, want %q", skill.Path, "/path/to/skill")
	}
	if skill.Category != "testing" {
		t.Errorf("Category = %q, want %q", skill.Category, "testing")
	}
}