package project

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDetector(t *testing.T) {
	d := NewDetector()
	if d == nil {
		t.Error("expected detector")
	}
}

func TestDetectProjectRoot(t *testing.T) {
	tmpDir := t.TempDir()

	root := DetectProjectRoot(tmpDir)
	if root != "" {
		t.Errorf("expected empty root for non-project dir, got %s", root)
	}

	gitDir := filepath.Join(tmpDir, ".git")
	os.Mkdir(gitDir, 0755)

	root = DetectProjectRoot(tmpDir)
	if root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, root)
	}
}

func TestDetectGit(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.Mkdir(gitDir, 0755)

	proj := detectGit(tmpDir)
	if proj == nil {
		t.Fatal("expected git project")
	}
	if proj.VCS != "git" {
		t.Errorf("expected VCS 'git', got %s", proj.VCS)
	}
	if proj.Path != tmpDir {
		t.Errorf("expected Path %s, got %s", tmpDir, proj.Path)
	}
}

func TestDetectNPM(t *testing.T) {
	tmpDir := t.TempDir()
	pkgFile := filepath.Join(tmpDir, "package.json")
	os.WriteFile(pkgFile, []byte(`{"name": "my-npm-project"}`), 0644)

	proj := detectNPM(tmpDir)
	if proj == nil {
		t.Fatal("expected npm project")
	}
	if proj.VCS != "npm" {
		t.Errorf("expected VCS 'npm', got %s", proj.VCS)
	}
	if proj.Name != "my-npm-project" {
		t.Errorf("expected Name 'my-npm-project', got %s", proj.Name)
	}
}

func TestDetectNPMWithoutName(t *testing.T) {
	tmpDir := t.TempDir()
	pkgFile := filepath.Join(tmpDir, "package.json")
	os.WriteFile(pkgFile, []byte(`{}`), 0644)

	proj := detectNPM(tmpDir)
	if proj == nil {
		t.Fatal("expected npm project")
	}
	if proj.Name != filepath.Base(tmpDir) {
		t.Errorf("expected Name '%s', got %s", filepath.Base(tmpDir), proj.Name)
	}
}

func TestDetectGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	modFile := filepath.Join(tmpDir, "go.mod")
	os.WriteFile(modFile, []byte("module my-go-module\n\ngo 1.21"), 0644)

	proj := detectGoMod(tmpDir)
	if proj == nil {
		t.Fatal("expected go module")
	}
	if proj.VCS != "go" {
		t.Errorf("expected VCS 'go', got %s", proj.VCS)
	}
	if proj.Name != "my-go-module" {
		t.Errorf("expected Name 'my-go-module', got %s", proj.Name)
	}
}

func TestDetectCargo(t *testing.T) {
	tmpDir := t.TempDir()
	cargoFile := filepath.Join(tmpDir, "Cargo.toml")
	os.WriteFile(cargoFile, []byte(`[package]
name = "my-cargo-project"
version = "0.1.0"
`), 0644)

	proj := detectCargo(tmpDir)
	if proj == nil {
		t.Fatal("expected cargo project")
	}
	if proj.VCS != "cargo" {
		t.Errorf("expected VCS 'cargo', got %s", proj.VCS)
	}
	if proj.Name != "my-cargo-project" {
		t.Errorf("expected Name 'my-cargo-project', got %s", proj.Name)
	}
}

func TestDetectPythonPyproject(t *testing.T) {
	tmpDir := t.TempDir()
	pyprojectFile := filepath.Join(tmpDir, "pyproject.toml")
	os.WriteFile(pyprojectFile, []byte(`[project]
name = "my-python-project"
version = "0.1.0"
`), 0644)

	proj := detectPython(tmpDir)
	if proj == nil {
		t.Fatal("expected python project")
	}
	if proj.VCS != "python" {
		t.Errorf("expected VCS 'python', got %s", proj.VCS)
	}
	if proj.Name != "my-python-project" {
		t.Errorf("expected Name 'my-python-project', got %s", proj.Name)
	}
}

func TestDetectPythonRequirements(t *testing.T) {
	tmpDir := t.TempDir()
	reqFile := filepath.Join(tmpDir, "requirements.txt")
	os.WriteFile(reqFile, []byte("requests>=2.28.0\n"), 0644)

	proj := detectPython(tmpDir)
	if proj == nil {
		t.Fatal("expected python project")
	}
	if proj.VCS != "python" {
		t.Errorf("expected VCS 'python', got %s", proj.VCS)
	}
	if proj.Name != filepath.Base(tmpDir) {
		t.Errorf("expected Name '%s', got %s", filepath.Base(tmpDir), proj.Name)
	}
}

func TestDetectWalkUp(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub", "path")
	os.MkdirAll(subDir, 0755)

	gitDir := filepath.Join(tmpDir, ".git")
	os.Mkdir(gitDir, 0755)

	root := DetectProjectRoot(subDir)
	if root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, root)
	}
}

func TestDetect(t *testing.T) {
	tmpDir := t.TempDir()
	modFile := filepath.Join(tmpDir, "go.mod")
	os.WriteFile(modFile, []byte("module test-module\n\ngo 1.21"), 0644)

	d := NewDetector()
	proj, err := d.Detect(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if proj == nil {
		t.Fatal("expected project")
	}
	if proj.Name != "test-module" {
		t.Errorf("expected Name 'test-module', got %s", proj.Name)
	}
}

func TestList(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	modFile1 := filepath.Join(tmpDir1, "go.mod")
	os.WriteFile(modFile1, []byte("module test-module-1\n\ngo 1.21"), 0644)

	pkgFile2 := filepath.Join(tmpDir2, "package.json")
	os.WriteFile(pkgFile2, []byte(`{"name": "test-module-2"}`), 0644)

	projects, err := List(context.Background(), []string{tmpDir1, tmpDir2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
}