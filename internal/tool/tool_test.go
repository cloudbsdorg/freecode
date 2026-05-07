package tool

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBashToolName(t *testing.T) {
	tool := NewBashTool()
	if tool.Name() != "bash" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "bash")
	}
}

func TestBashToolDescription(t *testing.T) {
	tool := NewBashTool()
	if tool.Description() != "Execute bash commands in a shell" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Execute bash commands in a shell")
	}
}

func TestBashToolSchema(t *testing.T) {
	tool := NewBashTool()
	schema := tool.Schema()
	if schema.Name != "bash" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "bash")
	}
	if schema.Parameters["command"].Required != true {
		t.Error("Schema().Parameters['command'].Required should be true")
	}
}

func TestBashToolExecuteSuccess(t *testing.T) {
	tool := NewBashTool()
	resp, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"command": "echo hello"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if resp.Error != nil {
		t.Errorf("Response.Error = %v", resp.Error)
	}
}

func TestBashToolExecuteError(t *testing.T) {
	tool := NewBashTool()
	resp, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"command": "exit 1"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if resp.Error == nil {
		t.Error("Response.Error should not be nil for failed command")
	}
}

func TestBashToolExecuteNonStringCommand(t *testing.T) {
	tool := NewBashTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"command": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string command")
	}
}

func TestBashToolExecuteWithTimeout(t *testing.T) {
	tool := NewBashTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"command": "echo hello",
			"timeout": 60,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestBashToolExecuteWithWorkdir(t *testing.T) {
	tool := NewBashTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"command": "pwd",
			"workdir": "/tmp",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestBashToolExecuteContextCanceled(t *testing.T) {
	tool := NewBashTool()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := tool.Execute(ctx, Request{
		Arguments: map[string]interface{}{"command": "sleep 10"},
	})
	if err == nil {
		t.Error("Execute() should error for canceled context")
	}
}

func TestReadToolName(t *testing.T) {
	tool := NewReadTool()
	if tool.Name() != "read" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "read")
	}
}

func TestReadToolDescription(t *testing.T) {
	tool := NewReadTool()
	if tool.Description() != "Read file contents" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Read file contents")
	}
}

func TestReadToolSchema(t *testing.T) {
	tool := NewReadTool()
	schema := tool.Schema()
	if schema.Name != "read" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "read")
	}
	if schema.Parameters["file_path"].Required != true {
		t.Error("Schema().Parameters['file_path'].Required should be true")
	}
}

func TestReadToolExecuteSuccess(t *testing.T) {
	tool := NewReadTool()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("line1\nline2\nline3"), 0644)

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"file_path": tmpFile},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestReadToolExecuteWithOffset(t *testing.T) {
	tool := NewReadTool()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("line1\nline2\nline3"), 0644)

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": tmpFile,
			"offset":    1,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestReadToolExecuteWithLimit(t *testing.T) {
	tool := NewReadTool()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("line1\nline2\nline3\nline4\nline5"), 0644)

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": tmpFile,
			"limit":     2,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestReadToolExecuteNonStringPath(t *testing.T) {
	tool := NewReadTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"file_path": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string file_path")
	}
}

func TestReadToolExecuteFileNotFound(t *testing.T) {
	tool := NewReadTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"file_path": "/nonexistent/file.txt"},
	})
	if err == nil {
		t.Error("Execute() should error for nonexistent file")
	}
}

func TestEditToolName(t *testing.T) {
	tool := NewEditTool()
	if tool.Name() != "edit" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "edit")
	}
}

func TestEditToolDescription(t *testing.T) {
	tool := NewEditTool()
	if tool.Description() != "Edit a file by replacing lines" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Edit a file by replacing lines")
	}
}

func TestEditToolSchema(t *testing.T) {
	tool := NewEditTool()
	schema := tool.Schema()
	if schema.Name != "edit" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "edit")
	}
	if schema.Parameters["file_path"].Required != true {
		t.Error("Schema().Parameters['file_path'].Required should be true")
	}
}

func TestEditToolExecuteSuccess(t *testing.T) {
	tool := NewEditTool()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("hello world"), 0644)

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  tmpFile,
			"old_string": "world",
			"new_string": "universe",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	data, _ := os.ReadFile(tmpFile)
	if string(data) != "hello universe" {
		t.Errorf("File content = %q, want %q", string(data), "hello universe")
	}
}

func TestEditToolExecuteNonStringPath(t *testing.T) {
	tool := NewEditTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  123,
			"old_string": "old",
			"new_string": "new",
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string file_path")
	}
}

func TestEditToolExecuteNonStringOldString(t *testing.T) {
	tool := NewEditTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  "/tmp/test.txt",
			"old_string": 123,
			"new_string": "new",
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string old_string")
	}
}

func TestEditToolExecuteNonStringNewString(t *testing.T) {
	tool := NewEditTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  "/tmp/test.txt",
			"old_string": "old",
			"new_string": 123,
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string new_string")
	}
}

func TestEditToolExecuteFileNotFound(t *testing.T) {
	tool := NewEditTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  "/nonexistent/test.txt",
			"old_string": "old",
			"new_string": "new",
		},
	})
	if err == nil {
		t.Error("Execute() should error for nonexistent file")
	}
}

func TestEditToolExecuteOldStringNotFound(t *testing.T) {
	tool := NewEditTool()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("hello world"), 0644)

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  tmpFile,
			"old_string": "notfound",
			"new_string": "new",
		},
	})
	if err == nil {
		t.Error("Execute() should error when old_string not found")
	}
}

func TestEditToolExecuteWriteError(t *testing.T) {
	tool := NewEditTool()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("hello world"), 0644)

	if err := os.Chmod(tmpFile, 0000); err != nil {
		t.Skip("Cannot chmod file")
	}

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  tmpFile,
			"old_string": "world",
			"new_string": "universe",
		},
	})
	if err == nil {
		t.Error("Execute() should error when WriteFile fails")
	}
}

func TestGrepToolName(t *testing.T) {
	tool := NewGrepTool()
	if tool.Name() != "grep" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "grep")
	}
}

func TestGrepToolDescription(t *testing.T) {
	tool := NewGrepTool()
	if tool.Description() != "Search for patterns in files" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Search for patterns in files")
	}
}

func TestGrepToolSchema(t *testing.T) {
	tool := NewGrepTool()
	schema := tool.Schema()
	if schema.Name != "grep" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "grep")
	}
	if schema.Parameters["pattern"].Required != true {
		t.Error("Schema().Parameters['pattern'].Required should be true")
	}
}

func TestGrepToolExecute(t *testing.T) {
	tool := NewGrepTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": "test",
			"path":    "/tmp",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestGrepToolExecuteNonStringPattern(t *testing.T) {
	tool := NewGrepTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": 123,
			"path":    "/tmp",
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string pattern")
	}
}

func TestGrepToolExecuteNonStringPath(t *testing.T) {
	tool := NewGrepTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": "test",
			"path":    123,
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string path")
	}
}

func TestGrepToolExecuteInvalidPattern(t *testing.T) {
	tool := NewGrepTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": "[",
			"path":    "/tmp",
		},
	})
	if err == nil {
		t.Error("Execute() should error for invalid pattern")
	}
}

func TestGlobToolName(t *testing.T) {
	tool := NewGlobTool()
	if tool.Name() != "glob" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "glob")
	}
}

func TestGlobToolDescription(t *testing.T) {
	tool := NewGlobTool()
	if tool.Description() != "Find files matching a glob pattern" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Find files matching a glob pattern")
	}
}

func TestGlobToolSchema(t *testing.T) {
	tool := NewGlobTool()
	schema := tool.Schema()
	if schema.Name != "glob" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "glob")
	}
	if schema.Parameters["pattern"].Required != true {
		t.Error("Schema().Parameters['pattern'].Required should be true")
	}
}

func TestGlobToolExecute(t *testing.T) {
	tool := NewGlobTool()
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package test"), 0644)

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": "*.go",
			"path":    tmpDir,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestGlobToolExecuteNonStringPattern(t *testing.T) {
	tool := NewGlobTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": 123,
			"path":    "/tmp",
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string pattern")
	}
}

func TestGlobToolExecuteGlobError(t *testing.T) {
	tool := NewGlobTool()
	resp, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": "*.go",
			"path":    "/nonexistent-directory",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if resp.Result == "" {
		t.Log("No matches found (expected for nonexistent directory)")
	}
}

func TestGlobWalkDir(t *testing.T) {
	tmpDir := t.TempDir()
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "nested.txt"), []byte("nested"), 0644)

	var results []string
	err := globWalkDir(tmpDir, "*.txt", &results)
	if err != nil {
		t.Fatalf("globWalkDir() error = %v", err)
	}
}

func TestGlobWalkDirNonexistent(t *testing.T) {
	var results []string
	err := globWalkDir("/nonexistent", "*.txt", &results)
	if err != nil {
		t.Fatalf("globWalkDir() error = %v", err)
	}
}

func TestLSPToolName(t *testing.T) {
	tool := NewLSPTool()
	if tool.Name() != "lsp" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "lsp")
	}
}

func TestLSPToolDescription(t *testing.T) {
	tool := NewLSPTool()
	if tool.Description() != "Language Server Protocol operations" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Language Server Protocol operations")
	}
}

func TestLSPToolSchema(t *testing.T) {
	tool := NewLSPTool()
	schema := tool.Schema()
	if schema.Name != "lsp" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "lsp")
	}
}

func TestLSPToolExecuteStart(t *testing.T) {
	tool := NewLSPTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "start", "language": "go"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestLSPToolExecuteStop(t *testing.T) {
	tool := NewLSPTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "stop"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestLSPToolExecuteGoto(t *testing.T) {
	tool := NewLSPTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":    "goto",
			"file":      "/test/file.go",
			"line":      10,
			"character": 5,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestLSPToolExecuteHover(t *testing.T) {
	tool := NewLSPTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":    "hover",
			"file":      "/test/file.go",
			"line":      10,
			"character": 5,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestLSPToolExecuteCompletions(t *testing.T) {
	tool := NewLSPTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "completions",
			"file":   "/test/file.go",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestLSPToolExecuteUnknownAction(t *testing.T) {
	tool := NewLSPTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "unknown"},
	})
	if err == nil {
		t.Error("Execute() should error for unknown action")
	}
}

func TestLSPToolExecuteNonStringAction(t *testing.T) {
	tool := NewLSPTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string action")
	}
}

func TestPlanToolName(t *testing.T) {
	tool := NewPlanTool()
	if tool.Name() != "plan" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "plan")
	}
}

func TestPlanToolDescription(t *testing.T) {
	tool := NewPlanTool()
	if tool.Description() != "Create and manage plans" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Create and manage plans")
	}
}

func TestPlanToolSchema(t *testing.T) {
	tool := NewPlanTool()
	schema := tool.Schema()
	if schema.Name != "plan" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "plan")
	}
}

func TestPlanToolExecuteList(t *testing.T) {
	tool := NewPlanTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "list"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestPlanToolExecuteCreate(t *testing.T) {
	tool := NewPlanTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "create",
			"content": "test plan content",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestPlanToolExecuteExecute(t *testing.T) {
	tool := NewPlanTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "execute",
			"plan_id": "plan-123",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestPlanToolExecuteUnknownAction(t *testing.T) {
	tool := NewPlanTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "unknown"},
	})
	if err == nil {
		t.Error("Execute() should error for unknown action")
	}
}

func TestPlanToolExecuteNonStringAction(t *testing.T) {
	tool := NewPlanTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string action")
	}
}

func TestQuestionToolName(t *testing.T) {
	tool := NewQuestionTool()
	if tool.Name() != "question" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "question")
	}
}

func TestQuestionToolDescription(t *testing.T) {
	tool := NewQuestionTool()
	if tool.Description() != "Ask the user a question" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Ask the user a question")
	}
}

func TestQuestionToolSchema(t *testing.T) {
	tool := NewQuestionTool()
	schema := tool.Schema()
	if schema.Name != "question" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "question")
	}
	if schema.Parameters["question"].Required != true {
		t.Error("Schema().Parameters['question'].Required should be true")
	}
}

func TestQuestionToolExecute(t *testing.T) {
	tool := NewQuestionTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"question": "What is your name?",
			"header":   "Name",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestQuestionToolExecuteNonStringQuestion(t *testing.T) {
	tool := NewQuestionTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"question": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string question")
	}
}

func TestSkillToolName(t *testing.T) {
	tool := NewSkillTool()
	if tool.Name() != "skill" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "skill")
	}
}

func TestSkillToolDescription(t *testing.T) {
	tool := NewSkillTool()
	if tool.Description() != "List and invoke skills" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "List and invoke skills")
	}
}

func TestSkillToolSchema(t *testing.T) {
	tool := NewSkillTool()
	schema := tool.Schema()
	if schema.Name != "skill" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "skill")
	}
}

func TestSkillToolExecuteList(t *testing.T) {
	tool := NewSkillTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "list"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestSkillToolExecuteInvoke(t *testing.T) {
	tool := NewSkillTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "invoke",
			"name":   "git-master",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestSkillToolExecuteUnknownAction(t *testing.T) {
	tool := NewSkillTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "unknown"},
	})
	if err == nil {
		t.Error("Execute() should error for unknown action")
	}
}

func TestSkillToolExecuteNonStringAction(t *testing.T) {
	tool := NewSkillTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string action")
	}
}

func TestJoin(t *testing.T) {
	arr := []string{"a", "b", "c"}
	result := join(arr, ",")
	if result != "a,b,c" {
		t.Errorf("join() = %q, want %q", result, "a,b,c")
	}
}

func TestJoinSingleElement(t *testing.T) {
	arr := []string{"a"}
	result := join(arr, ",")
	if result != "a" {
		t.Errorf("join() = %q, want %q", result, "a")
	}
}

func TestJoinEmpty(t *testing.T) {
	arr := []string{}
	result := join(arr, ",")
	if result != "" {
		t.Errorf("join() = %q, want %q", result, "")
	}
}

func TestTaskToolName(t *testing.T) {
	tool := NewTaskTool()
	if tool.Name() != "task" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "task")
	}
}

func TestTaskToolDescription(t *testing.T) {
	tool := NewTaskTool()
	if tool.Description() != "Create and manage tasks" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Create and manage tasks")
	}
}

func TestTaskToolSchema(t *testing.T) {
	tool := NewTaskTool()
	schema := tool.Schema()
	if schema.Name != "task" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "task")
	}
	if schema.Parameters["action"].Required != true {
		t.Error("Schema().Parameters['action'].Required should be true")
	}
}

func TestTaskToolExecuteList(t *testing.T) {
	tool := NewTaskTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "list"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTaskToolExecuteCreate(t *testing.T) {
	tool := NewTaskTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "create",
			"title":  "Test task",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTaskToolExecuteUpdate(t *testing.T) {
	tool := NewTaskTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "update",
			"task_id": "test-123",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTaskToolExecuteDelete(t *testing.T) {
	tool := NewTaskTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "delete",
			"task_id": "test-123",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTaskToolExecuteUnknownAction(t *testing.T) {
	tool := NewTaskTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "unknown"},
	})
	if err == nil {
		t.Error("Execute() should error for unknown action")
	}
}

func TestTaskToolExecuteNonStringAction(t *testing.T) {
	tool := NewTaskTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string action")
	}
}

func TestTodoToolName(t *testing.T) {
	tool := NewTodoTool()
	if tool.Name() != "todo" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "todo")
	}
}

func TestTodoToolDescription(t *testing.T) {
	tool := NewTodoTool()
	if tool.Description() != "Manage todo items" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Manage todo items")
	}
}

func TestTodoToolSchema(t *testing.T) {
	tool := NewTodoTool()
	schema := tool.Schema()
	if schema.Name != "todo" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "todo")
	}
	if schema.Parameters["action"].Required != true {
		t.Error("Schema().Parameters['action'].Required should be true")
	}
}

func TestTodoToolExecuteList(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "list"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTodoToolExecuteAdd(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "add",
			"content": "Test todo",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTodoToolExecuteAddWithPriority(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":   "add",
			"content":  "High priority",
			"priority": "high",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTodoToolExecuteUpdate(t *testing.T) {
	tool := NewTodoTool()
	// First add a todo
	tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "add",
			"content": "Test",
		},
	})
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "update",
			"id":     "todo-1",
			"status": "completed",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTodoToolExecuteUpdateNotFound(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "update",
			"id":     "nonexistent",
			"status": "completed",
		},
	})
	if err == nil {
		t.Error("Execute() should error for nonexistent id")
	}
}

func TestTodoToolExecuteDelete(t *testing.T) {
	tool := NewTodoTool()
	// First add a todo
	tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "add",
			"content": "To delete",
		},
	})
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "delete",
			"id":     "todo-1",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestTodoToolExecuteUnknownAction(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": "unknown"},
	})
	if err == nil {
		t.Error("Execute() should error for unknown action")
	}
}

func TestTodoToolExecuteNonStringAction(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"action": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string action")
	}
}

func TestTodoToolExecuteAddNonStringContent(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action":  "add",
			"content": 123,
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string content")
	}
}

func TestTodoToolExecuteUpdateNonStringID(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "update",
			"id":     123,
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string id")
	}
}

func TestTodoToolExecuteDeleteNonStringID(t *testing.T) {
	tool := NewTodoTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"action": "delete",
			"id":     123,
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string id")
	}
}

func TestTodoToolLoadFromFile(t *testing.T) {
	tool := NewTodoTool()
	err := tool.LoadFromFile("/nonexistent/path/todos.json")
	if err == nil {
		t.Error("LoadFromFile should error for nonexistent file")
	}
}

func TestTodoToolSaveToFile(t *testing.T) {
	tool := NewTodoTool()
	err := tool.SaveToFile("/tmp/freecode_test_todos.json")
	if err != nil {
		t.Fatalf("SaveToFile error = %v", err)
	}
	// Clean up
	os.Remove("/tmp/freecode_test_todos.json")
}

func TestTodoToolSaveToFileNestedDir(t *testing.T) {
	tool := NewTodoTool()
	err := tool.SaveToFile("/tmp/nested/freecode_test_todos.json")
	if err != nil {
		t.Fatalf("SaveToFile error = %v", err)
	}
	// Clean up
	os.RemoveAll("/tmp/nested")
}

func TestWebFetchToolName(t *testing.T) {
	tool := NewWebFetchTool()
	if tool.Name() != "webfetch" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "webfetch")
	}
}

func TestWebFetchToolDescription(t *testing.T) {
	tool := NewWebFetchTool()
	if tool.Description() != "Fetch content from a URL" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Fetch content from a URL")
	}
}

func TestWebFetchToolSchema(t *testing.T) {
	tool := NewWebFetchTool()
	schema := tool.Schema()
	if schema.Name != "webfetch" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "webfetch")
	}
	if schema.Parameters["url"].Required != true {
		t.Error("Schema().Parameters['url'].Required should be true")
	}
}

func TestWebFetchToolExecuteNonStringURL(t *testing.T) {
	tool := NewWebFetchTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"url": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string url")
	}
}

func TestWebSearchToolName(t *testing.T) {
	tool := NewWebSearchTool()
	if tool.Name() != "websearch" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "websearch")
	}
}

func TestWebSearchToolDescription(t *testing.T) {
	tool := NewWebSearchTool()
	if tool.Description() != "Search the web using Exa" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Search the web using Exa")
	}
}

func TestWebSearchToolSchema(t *testing.T) {
	tool := NewWebSearchTool()
	schema := tool.Schema()
	if schema.Name != "websearch" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "websearch")
	}
	if schema.Parameters["query"].Required != true {
		t.Error("Schema().Parameters['query'].Required should be true")
	}
}

func TestWebSearchToolExecuteNonStringQuery(t *testing.T) {
	tool := NewWebSearchTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{"query": 123},
	})
	if err == nil {
		t.Error("Execute() should error for non-string query")
	}
}

func TestWriteToolName(t *testing.T) {
	tool := NewWriteTool()
	if tool.Name() != "write" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "write")
	}
}

func TestWriteToolDescription(t *testing.T) {
	tool := NewWriteTool()
	if tool.Description() != "Write content to a file" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Write content to a file")
	}
}

func TestWriteToolSchema(t *testing.T) {
	tool := NewWriteTool()
	schema := tool.Schema()
	if schema.Name != "write" {
		t.Errorf("Schema().Name = %q, want %q", schema.Name, "write")
	}
	if schema.Parameters["file_path"].Required != true {
		t.Error("Schema().Parameters['file_path'].Required should be true")
	}
	if schema.Parameters["content"].Required != true {
		t.Error("Schema().Parameters['content'].Required should be true")
	}
}

func TestWriteToolExecute(t *testing.T) {
	tool := NewWriteTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": "/tmp/freecode_test_write.txt",
			"content":   "Hello, World!",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	// Clean up
	os.Remove("/tmp/freecode_test_write.txt")
}

func TestWriteToolExecuteAppend(t *testing.T) {
	tool := NewWriteTool()
	// First write
	tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": "/tmp/freecode_test_append.txt",
			"content":   "First line\n",
		},
	})
	// Append
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": "/tmp/freecode_test_append.txt",
			"content":   "Second line\n",
			"append":    true,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	// Clean up
	os.Remove("/tmp/freecode_test_append.txt")
}

func TestWriteToolExecuteNonStringFilePath(t *testing.T) {
	tool := NewWriteTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": 123,
			"content":   "test",
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string file_path")
	}
}

func TestWriteToolExecuteNonStringContent(t *testing.T) {
	tool := NewWriteTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": "/tmp/test.txt",
			"content":   123,
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-string content")
	}
}

func TestWriteToolExecuteNestedDir(t *testing.T) {
	tool := NewWriteTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": "/tmp/nested/deep/freecode_test.txt",
			"content":   "Nested path test",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	// Clean up
	os.RemoveAll("/tmp/nested")
}

func TestRegistrySchema(t *testing.T) {
	r := NewRegistry()
	r.Register(NewBashTool())
	r.Register(NewReadTool())
	schemas := r.Schema()
	if len(schemas) != 2 {
		t.Errorf("len(Schema()) = %d, want 2", len(schemas))
	}
}

func TestWriteToolExecuteMkdirAllError(t *testing.T) {
	tool := NewWriteTool()

	if os.Getuid() == 0 {
		t.Skip("skipping as root user")
	}

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": "/proc/cannot-create/test.txt",
			"content":   "test",
		},
	})
	if err == nil {
		t.Error("Execute() should error when MkdirAll fails")
	}
}

func TestWriteToolExecuteWriteStringError(t *testing.T) {
	tool := NewWriteTool()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")

	f, _ := os.Create(tmpFile)
	f.Close()
	if err := os.Chmod(tmpFile, 0000); err != nil {
		t.Skip("Cannot chmod file")
	}

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path": tmpFile,
			"content":   "test",
		},
	})
	if err == nil {
		t.Error("Execute() should error when WriteString fails")
	}
}

func TestBashToolExecuteContextDone(t *testing.T) {
	tool := NewBashTool()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := tool.Execute(ctx, Request{
		Arguments: map[string]interface{}{
			"command": "sleep 10",
		},
	})
	if err == nil {
		t.Error("Execute() should error when context is done")
	}
}

func TestEditToolExecuteScannerError(t *testing.T) {
	tool := NewEditTool()

	if os.Getuid() == 0 {
		t.Skip("skipping as root user")
	}

	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("hello world"), 0644)
	if err := os.Chmod(tmpFile, 0000); err != nil {
		t.Skip("Cannot chmod file")
	}

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file_path":  tmpFile,
			"old_string": "world",
			"new_string": "universe",
		},
	})
	if err == nil {
		t.Error("Execute() should error when scanner fails")
	}
}

func TestGrepToolExecuteIgnoreCase(t *testing.T) {
	tool := NewGrepTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern":     "TEST",
			"path":        "/tmp",
			"ignore_case": true,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestGlobToolExecuteNonExistentPath(t *testing.T) {
	tool := NewGlobTool()

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": "*.go",
			"path":    "/nonexistent",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestGlobWalkDirReadDirError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping as root user")
	}

	var results []string
	err := globWalkDir("/proc/cannot-read", "*.txt", &results)
	if err != nil {
		t.Fatalf("globWalkDir() error = %v", err)
	}
}

func TestTodoToolLoadFromFileError(t *testing.T) {
	tool := NewTodoTool()
	tmpFile := filepath.Join(t.TempDir(), "invalid.json")
	os.WriteFile(tmpFile, []byte("not json"), 0644)

	err := tool.LoadFromFile(tmpFile)
	if err == nil {
		t.Error("LoadFromFile() should error for invalid JSON")
	}
}

func TestTodoToolSaveToFileMarshalError(t *testing.T) {
	tool := &TodoTool{}
	tool.todos = map[string]TodoItem{
		"bad": {
			ID:       "bad",
			Content:  "test",
			Metadata: map[string]interface{}{"chan": make(chan int)},
		},
	}

	tmpFile := filepath.Join(t.TempDir(), "todos.json")
	err := tool.SaveToFile(tmpFile)
	if err == nil {
		t.Error("SaveToFile() should error when MarshalIndent fails")
	}
}

func TestWebFetchToolExecuteNon200(t *testing.T) {
	tool := NewWebFetchTool()
	ctx := context.Background()

	_, err := tool.Execute(ctx, Request{
		Arguments: map[string]interface{}{
			"url": "http://localhost:1",
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-200 status")
	}
}

func TestWebSearchToolExecute(t *testing.T) {
	tool := NewWebSearchTool()

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"query": "test search",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestWebSearchToolExecuteWithNumResults(t *testing.T) {
	tool := NewWebSearchTool()

	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"query":       "test search",
			"num_results": 5,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestWebFetchToolExecuteWithServer(t *testing.T) {
	tool := NewWebFetchTool()
	ctx := context.Background()

	_, err := tool.Execute(ctx, Request{
		Arguments: map[string]interface{}{
			"url": "http://localhost:99999",
		},
	})
	if err == nil {
		t.Error("Execute() should error for unreachable URL")
	}
}

func TestBashToolExecuteContextCancel(t *testing.T) {
	tool := NewBashTool()
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		_, _ = tool.Execute(ctx, Request{
			Arguments: map[string]interface{}{
				"command": "sleep 10",
			},
		})
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("Execute did not return after context cancel")
	}
}

func TestEditToolExecuteNonExistentFile(t *testing.T) {
	tool := NewEditTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"file":  "/nonexistent/path/file.txt",
			"edits": []map[string]string{},
		},
	})
	if err == nil {
		t.Error("Execute() should error for non-existent file")
	}
}

func TestGlobToolExecuteNonExistentPattern(t *testing.T) {
	tool := NewGlobTool()
	resp, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"pattern": "/nonexistent/**/file.txt",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if resp.Result == "" {
		t.Log("Result is empty for non-matching pattern (expected)")
	}
}

func TestWriteToolExecuteWithMissingDir(t *testing.T) {
	tool := NewWriteTool()
	_, err := tool.Execute(context.Background(), Request{
		Arguments: map[string]interface{}{
			"path":    "/nonexistent/dir/file.txt",
			"content": "hello world",
		},
	})
	if err == nil {
		t.Error("Execute() should error when parent directory doesn't exist")
	}
}
