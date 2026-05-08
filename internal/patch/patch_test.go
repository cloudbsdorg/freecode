package patch

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestApply_EmptyPath(t *testing.T) {
	patch := Patch{Path: ""}
	err := Apply(context.Background(), patch)
	if err != ErrInvalidDiff {
		t.Errorf("Expected ErrInvalidDiff, got %v", err)
	}
}

func TestApply_FileNotFound(t *testing.T) {
	patch := Patch{Path: "/nonexistent/file.txt"}
	err := Apply(context.Background(), patch)
	if err != ErrFileNotFound {
		t.Errorf("Expected ErrFileNotFound, got %v", err)
	}
}

func TestApply_Success(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(file, []byte("line1\nline2\nline3\n"), 0644)

	patch := Patch{
		Path:     file,
		OldLines: []string{"line1", "line2"},
		NewLines: []string{"line1", "modified", "line2"},
	}

	err := Apply(context.Background(), patch)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, _ := os.ReadFile(file)
	expected := "line1\nmodified\nline2\nline3\n"
	if string(content) != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, string(content))
	}
}

func TestApply_Mismatch(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(file, []byte("line1\nline2\nline3\n"), 0644)

	patch := Patch{
		Path:     file,
		OldLines: []string{"different", "content"},
		NewLines: []string{"new", "content"},
	}

	err := Apply(context.Background(), patch)
	if err != ErrMismatch {
		t.Errorf("Expected ErrMismatch, got %v", err)
	}
}

func TestApply_InsertAtBeginning(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(file, []byte("line2\nline3\n"), 0644)

	patch := Patch{
		Path:     file,
		OldLines: []string{},
		NewLines: []string{"line1"},
	}

	err := Apply(context.Background(), patch)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, _ := os.ReadFile(file)
	if string(content) != "line1\nline2\nline3\n" {
		t.Errorf("Unexpected content: %s", string(content))
	}
}

func TestCreate(t *testing.T) {
	patch, err := Create(context.Background(), "old\ncontent", "new\ncontent")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if len(patch.OldLines) != 2 {
		t.Errorf("Expected 2 old lines, got %d", len(patch.OldLines))
	}
	if len(patch.NewLines) != 2 {
		t.Errorf("Expected 2 new lines, got %d", len(patch.NewLines))
	}
}

func TestParse_Simple(t *testing.T) {
	diff := `--- a/file.txt
+++ b/file.txt
@@ -1,3 +1,4 @@
 line1
-line2
+new line
+added line
 line3`

	patches, err := Parse(context.Background(), diff)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(patches) != 1 {
		t.Fatalf("Expected 1 patch, got %d", len(patches))
	}

	if patches[0].Path != "file.txt" {
		t.Errorf("Expected path 'file.txt', got '%s'", patches[0].Path)
	}
}

func TestParse_MultipleFiles(t *testing.T) {
	diff := `--- a/file1.txt
+++ b/file1.txt
@@ -1,2 +1,2 @@
 old1
-old2
+new2
--- a/file2.txt
+++ b/file2.txt
@@ -1,1 +1,1 @@
 f2old
+f2new`

	patches, err := Parse(context.Background(), diff)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(patches) != 2 {
		t.Fatalf("Expected 2 patches, got %d", len(patches))
	}
}

func TestParseFile(t *testing.T) {
	tmpDir := t.TempDir()
	diffFile := filepath.Join(tmpDir, "diff.txt")
	os.WriteFile(diffFile, []byte(`--- a/file.txt
+++ b/file.txt
@@ -1,2 +1,2 @@
 old
-old2
+new2`), 0644)

	patches, err := ParseFile(context.Background(), diffFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(patches) != 1 {
		t.Fatalf("Expected 1 patch, got %d", len(patches))
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile(context.Background(), "/nonexistent/diff.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestApplyToFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create original file
	originalFile := filepath.Join(tmpDir, "original.txt")
	os.WriteFile(originalFile, []byte("line1\nline2\nline3\n"), 0644)

	// Create diff file
	diffFile := filepath.Join(tmpDir, "diff.txt")
	os.WriteFile(diffFile, []byte(`--- a/original.txt
+++ b/original.txt
@@ -1,3 +1,3 @@
 line1
-line2
+modified line2
 line3`), 0644)

	err := ApplyToFile(context.Background(), diffFile, tmpDir)
	if err != nil {
		t.Fatalf("ApplyToFile failed: %v", err)
	}

	content, _ := os.ReadFile(originalFile)
	expected := "line1\nmodified line2\nline3\n"
	if string(content) != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, string(content))
	}
}

func TestCreateDiff(t *testing.T) {
	oldContent := "line1\nline2\nline3\n"
	newContent := "line1\nmodified\nline3\n"

	diff := CreateDiff("file.txt", oldContent, newContent)

	if !contains(diff, "--- file.txt") {
		t.Error("Diff should contain '--- file.txt'")
	}
	if !contains(diff, "+++ file.txt") {
		t.Error("Diff should contain '+++ file.txt'")
	}
	if !contains(diff, "-line2") {
		t.Error("Diff should contain removed line '-line2'")
	}
	if !contains(diff, "+modified") {
		t.Error("Diff should contain added line '+modified'")
	}
}

func TestCreateDiff_EmptyOld(t *testing.T) {
	diff := CreateDiff("new.txt", "", "new content\n")
	if !contains(diff, "--- new.txt") {
		t.Error("Diff should contain '--- new.txt'")
	}
	if !contains(diff, "+++ new.txt") {
		t.Error("Diff should contain '+++ new.txt'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestFindSequence(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e"}
	seq := []string{"b", "c"}

	idx := findSequence(lines, seq)
	if idx != 1 {
		t.Errorf("Expected index 1, got %d", idx)
	}
}

func TestFindSequence_NotFound(t *testing.T) {
	lines := []string{"a", "b", "c"}
	seq := []string{"x", "y"}

	idx := findSequence(lines, seq)
	if idx != -1 {
		t.Errorf("Expected -1, got %d", idx)
	}
}

func TestFindSequence_Empty(t *testing.T) {
	lines := []string{"a", "b", "c"}

	idx := findSequence(lines, []string{})
	if idx != 0 {
		t.Errorf("Expected 0 for empty sequence, got %d", idx)
	}
}

func TestFindSequence_SingleElement(t *testing.T) {
	lines := []string{"a", "b", "c"}

	idx := findSequence(lines, []string{"b"})
	if idx != 1 {
		t.Errorf("Expected 1, got %d", idx)
	}
}

func TestApplyPatchToLines_NoOldLines(t *testing.T) {
	original := []string{"a", "b", "c"}
	new := []string{"x", "y"}

	result, ok := applyPatchToLines(original, []string{}, new)
	if !ok {
		t.Error("Expected ok=true")
	}
	if len(result) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(result))
	}
	if result[0] != "x" || result[1] != "y" {
		t.Error("Expected new lines at beginning")
	}
}

func TestApplyPatchToLines_MiddleReplacement(t *testing.T) {
	original := []string{"a", "b", "c", "d"}
	old := []string{"b", "c"}
	new := []string{"x", "y"}

	result, ok := applyPatchToLines(original, old, new)
	if !ok {
		t.Error("Expected ok=true")
	}
	if len(result) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(result))
	}
	expected := []string{"a", "x", "y", "d"}
	for i, e := range expected {
		if result[i] != e {
			t.Errorf("Line %d: expected '%s', got '%s'", i, e, result[i])
		}
	}
}

func TestApplyPatchToLines_NotFound(t *testing.T) {
	original := []string{"a", "b", "c"}
	old := []string{"x", "y"}
	new := []string{"p", "q"}

	result, ok := applyPatchToLines(original, old, new)
	if ok {
		t.Error("Expected ok=false")
	}
	if result != nil {
		t.Error("Expected nil result when not found")
	}
}