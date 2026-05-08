package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()

	if Exists(tmpDir) != true {
		t.Errorf("Exists(%q) = false, want true", tmpDir)
	}

	nonExistent := filepath.Join(tmpDir, "doesnotexist")
	if Exists(nonExistent) != false {
		t.Errorf("Exists(%q) = true, want false", nonExistent)
	}

	file := filepath.Join(tmpDir, "testfile")
	if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if Exists(file) != true {
		t.Errorf("Exists(%q) = false, want true", file)
	}
}

func TestIsDir(t *testing.T) {
	tmpDir := t.TempDir()

	if !IsDir(tmpDir) {
		t.Errorf("IsDir(%q) = false, want true", tmpDir)
	}

	nonExistent := filepath.Join(tmpDir, "doesnotexist")
	if IsDir(nonExistent) != false {
		t.Errorf("IsDir(%q) = true, want false", nonExistent)
	}

	file := filepath.Join(tmpDir, "testfile")
	if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if IsDir(file) != false {
		t.Errorf("IsDir(%q) = true, want false", file)
	}
}

func TestStat(t *testing.T) {
	tmpDir := t.TempDir()

	info, err := Stat(tmpDir)
	if err != nil {
		t.Errorf("Stat(%q) error = %v, want nil", tmpDir, err)
	}
	if info == nil {
		t.Errorf("Stat(%q) = nil, want FileInfo", tmpDir)
	}
	if !info.IsDir() {
		t.Errorf("Stat(%q).IsDir() = false, want true", tmpDir)
	}

	nonExistent := filepath.Join(tmpDir, "doesnotexist")
	info, err = Stat(nonExistent)
	if err != nil {
		t.Errorf("Stat(%q) error = %v, want nil for non-existent path", nonExistent, err)
	}
	if info != nil {
		t.Errorf("Stat(%q) = %v, want nil for non-existent path", nonExistent, info)
	}

	file := filepath.Join(tmpDir, "testfile")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	info, err = Stat(file)
	if err != nil {
		t.Errorf("Stat(%q) error = %v, want nil", file, err)
	}
	if info == nil {
		t.Errorf("Stat(%q) = nil, want FileInfo", file)
	}
	if info.Size() != 5 {
		t.Errorf("Stat(%q).Size() = %d, want 5", file, info.Size())
	}
}

func TestReadText(t *testing.T) {
	tmpDir := t.TempDir()

	file := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	got, err := ReadText(file)
	if err != nil {
		t.Errorf("ReadText(%q) error = %v, want nil", file, err)
	}
	if got != content {
		t.Errorf("ReadText(%q) = %q, want %q", file, got, content)
	}

	nonExistent := filepath.Join(tmpDir, "doesnotexist")
	_, err = ReadText(nonExistent)
	if err == nil {
		t.Errorf("ReadText(%q) error = nil, want error", nonExistent)
	}
}

func TestReadBytes(t *testing.T) {
	tmpDir := t.TempDir()

	file := filepath.Join(tmpDir, "test.bin")
	content := []byte{0x01, 0x02, 0x03, 0xff}
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	got, err := ReadBytes(file)
	if err != nil {
		t.Errorf("ReadBytes(%q) error = %v, want nil", file, err)
	}
	if string(got) != string(content) {
		t.Errorf("ReadBytes(%q) = %v, want %v", file, got, content)
	}

	nonExistent := filepath.Join(tmpDir, "doesnotexist")
	_, err = ReadBytes(nonExistent)
	if err == nil {
		t.Errorf("ReadBytes(%q) error = nil, want error", nonExistent)
	}
}

func TestReadJson(t *testing.T) {
	tmpDir := t.TempDir()

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person := Person{Name: "Alice", Age: 30}
	data, _ := json.Marshal(person)
	file := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(file, data, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	var got Person
	got, err := ReadJson[Person](file)
	if err != nil {
		t.Errorf("ReadJson(%q) error = %v, want nil", file, err)
	}
	if got.Name != person.Name || got.Age != person.Age {
		t.Errorf("ReadJson(%q) = %+v, want %+v", file, got, person)
	}

	type Nested struct {
		Items []int          `json:"items"`
		Map   map[string]int `json:"map"`
	}
	nested := Nested{Items: []int{1, 2, 3}, Map: map[string]int{"a": 1}}
	nestedData, _ := json.Marshal(nested)
	nestedFile := filepath.Join(tmpDir, "nested.json")
	if err := os.WriteFile(nestedFile, nestedData, 0644); err != nil {
		t.Fatalf("failed to create nested test file: %v", err)
	}

	var gotNested Nested
	gotNested, err = ReadJson[Nested](nestedFile)
	if err != nil {
		t.Errorf("ReadJson(%q) error = %v, want nil", nestedFile, err)
	}
	if len(gotNested.Items) != 3 || gotNested.Map["a"] != 1 {
		t.Errorf("ReadJson(%q) = %+v, want %+v", nestedFile, gotNested, nested)
	}

	invalidFile := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidFile, []byte("not json"), 0644); err != nil {
		t.Fatalf("failed to create invalid json file: %v", err)
	}
	_, err = ReadJson[Person](invalidFile)
	if err == nil {
		t.Errorf("ReadJson(%q) error = nil, want error for invalid JSON", invalidFile)
	}
}

func TestWrite(t *testing.T) {
	tmpDir := t.TempDir()

	file := filepath.Join(tmpDir, "subdir", "test.txt")
	content := []byte("Hello, World!")

	err := Write(file, content, 0644)
	if err != nil {
		t.Errorf("Write(%q, ...) error = %v, want nil", file, err)
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("Write wrote %q, want %q", string(got), string(content))
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("failed to stat written file: %v", err)
	}
	if info.Mode() != 0644 {
		t.Errorf("Write file mode = %v, want 0644", info.Mode())
	}

	overwrite := filepath.Join(tmpDir, "overwrite.txt")
	if err := os.WriteFile(overwrite, []byte("old"), 0644); err != nil {
		t.Fatalf("failed to create overwrite target: %v", err)
	}
	err = Write(overwrite, []byte("new"), 0600)
	if err != nil {
		t.Errorf("Write(%q, ...) error = %v, want nil", overwrite, err)
	}
	info, _ = os.Stat(overwrite)
	if info.Mode() != 0600 {
		t.Errorf("Write file mode = %v, want 0600", info.Mode())
	}
}

func TestWriteJson(t *testing.T) {
	tmpDir := t.TempDir()

	type Config struct {
		Name string `json:"name"`
		Port int    `json:"port"`
	}
	config := Config{Name: "server", Port: 8080}

	file := filepath.Join(tmpDir, "config.json")
	err := WriteJson(file, config, 0644)
	if err != nil {
		t.Errorf("WriteJson(%q, ...) error = %v, want nil", file, err)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	var got Config
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("written file is not valid JSON: %v", err)
	}
	if got.Name != config.Name || got.Port != config.Port {
		t.Errorf("WriteJson wrote %+v, want %+v", got, config)
	}

	if !strings.HasSuffix(string(data), "\n") {
		t.Errorf("WriteJson should end with newline")
	}
}

func TestResolve(t *testing.T) {
	tmpDir := t.TempDir()

	absPath := filepath.Join(tmpDir, "test")
	if err := os.MkdirAll(absPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	resolved := Resolve(absPath)
	if resolved != absPath {
		t.Errorf("Resolve(%q) = %q, want %q", absPath, resolved, absPath)
	}

	parent := filepath.Dir(absPath)
	childPath := filepath.Join(absPath, "subdir", "..", "..", filepath.Base(absPath))
	resolved = Resolve(childPath)
	if resolved != parent && !strings.HasPrefix(resolved, parent) {
		t.Errorf("Resolve(%q) = %q, should resolve to within %q", childPath, resolved, parent)
	}

	symlinkPath := filepath.Join(tmpDir, "symlink")
	if err := os.Symlink(absPath, symlinkPath); err != nil {
		t.Skipf("symlinks not supported, skipping symlink test: %v", err)
	}

	resolved = Resolve(symlinkPath)
	if resolved != absPath {
		t.Errorf("Resolve(symlink) = %q, want %q", resolved, absPath)
	}
}

func TestPathContains(t *testing.T) {
	tmpDir := t.TempDir()

	parent := filepath.Join(tmpDir, "project")
	child := filepath.Join(parent, "src", "main.go")

	if !PathContains(parent, child) {
		t.Errorf("PathContains(%q, %q) = false, want true", parent, child)
	}

	if !PathContains(parent, parent) {
		t.Errorf("PathContains(%q, %q) = false, want true (same path)", parent, parent)
	}

	sibling := filepath.Join(tmpDir, "other")
	if PathContains(parent, sibling) {
		t.Errorf("PathContains(%q, %q) = true, want false (sibling)", parent, sibling)
	}

	grandParent := filepath.Dir(parent)
	if PathContains(parent, grandParent) {
		t.Errorf("PathContains(%q, %q) = true, want false (parent dir)", parent, grandParent)
	}

	unrelated := filepath.Join(tmpDir, "..", "other")
	if PathContains(parent, unrelated) {
		t.Errorf("PathContains(%q, %q) = true, want false (outside)", parent, unrelated)
	}
}

func TestFindUp(t *testing.T) {
	tmpDir := t.TempDir()

	level1 := filepath.Join(tmpDir, "level1")
	level2 := filepath.Join(level1, "level2")
	level3 := filepath.Join(level2, "level3")

	if err := os.MkdirAll(level3, 0755); err != nil {
		t.Fatalf("failed to create directories: %v", err)
	}

	goMod := filepath.Join(level1, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	tests := []struct {
		name    string
		target  string
		start   string
		stop    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "find go.mod from level3",
			target:  "go.mod",
			start:   level3,
			stop:    "",
			wantLen: 1,
		},
		{
			name:    "find go.mod from level3 with stop",
			target:  "go.mod",
			start:   level3,
			stop:    level1,
			wantLen: 0,
		},
		{
			name:    "find non-existent file",
			target:  "nonexistent.txt",
			start:   level3,
			stop:    "",
			wantLen: 0,
		},
		{
			name:    "find root directory marker",
			target:  "go.mod",
			start:   level3,
			stop:    tmpDir,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindUp(tt.target, tt.start, tt.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUp(%q, %q, %q) error = %v, wantErr %v", tt.target, tt.start, tt.stop, err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("FindUp(%q, %q, %q) returned %d items, want %d", tt.target, tt.start, tt.stop, len(got), tt.wantLen)
			}
		})
	}
}

func TestFindUpMultiple(t *testing.T) {
	tmpDir := t.TempDir()

	project := filepath.Join(tmpDir, "project")
	subdir := filepath.Join(project, "subdir")

	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("failed to create directories: %v", err)
	}

	readme := filepath.Join(project, "README.md")
	gitignore := filepath.Join(project, ".gitignore")

	if err := os.WriteFile(readme, []byte("# Project"), 0644); err != nil {
		t.Fatalf("failed to create README: %v", err)
	}
	if err := os.WriteFile(gitignore, []byte("node_modules"), 0644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	got, err := FindUpMulti([]string{"README.md", ".gitignore"}, subdir, "", false)
	if err != nil {
		t.Errorf("FindUpMulti error = %v, want nil", err)
	}
	if len(got) != 2 {
		t.Errorf("FindUpMulti returned %d items, want 2", len(got))
	}
}

func TestWindowsPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("WindowsPath is for non-Windows platforms")
	}

	tests := []struct {
		name string
		p    string
		want string
	}{
		{
			name: "git bash lowercase",
			p:    "/c/Users/test",
			want: "C:/Users/test",
		},
		{
			name: "git bash uppercase",
			p:    "/C/Users/test",
			want: "C:/Users/test",
		},
		{
			name: "cygwin path",
			p:    "/cygdrive/c/Users/test",
			want: "C:/Users/test",
		},
		{
			name: "wsl path",
			p:    "/mnt/c/Users/test",
			want: "C:/Users/test",
		},
		{
			name: "unix path unchanged",
			p:    "/home/user/test",
			want: "/home/user/test",
		},
		{
			name: "relative path unchanged",
			p:    "relative/path",
			want: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WindowsPath(tt.p)
			if got != tt.want {
				t.Errorf("WindowsPath(%q) = %q, want %q", tt.p, got, tt.want)
			}
		})
	}
}

func TestWriteStream(t *testing.T) {
	tmpDir := t.TempDir()

	file := filepath.Join(tmpDir, "subdir", "stream.txt")
	content := "Hello from stream!"

	err := WriteStream(file, strings.NewReader(content), 0644)
	if err != nil {
		t.Errorf("WriteStream error = %v, want nil", err)
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(got) != content {
		t.Errorf("WriteStream wrote %q, want %q", string(got), content)
	}
}

func TestSize(t *testing.T) {
	tmpDir := t.TempDir()

	file := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	size, err := Size(file)
	if err != nil {
		t.Errorf("Size(%q) error = %v, want nil", file, err)
	}
	if size != int64(len(content)) {
		t.Errorf("Size(%q) = %d, want %d", file, size, len(content))
	}

	nonExistent := filepath.Join(tmpDir, "doesnotexist")
	size, err = Size(nonExistent)
	if err != nil {
		t.Errorf("Size(%q) error = %v, want nil for non-existent", nonExistent, err)
	}
	if size != 0 {
		t.Errorf("Size(%q) = %d, want 0 for non-existent", nonExistent, size)
	}
}

func TestCopy(t *testing.T) {
	tmpDir := t.TempDir()

	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "subdir", "dest.txt")
	content := "Copy me!"

	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	err := Copy(dst, src)
	if err != nil {
		t.Errorf("Copy error = %v, want nil", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(got) != content {
		t.Errorf("Copy wrote %q, want %q", string(got), content)
	}

	nonExistent := filepath.Join(tmpDir, "nonexistent.txt")
	err = Copy(dst, nonExistent)
	if err == nil {
		t.Errorf("Copy from non-existent source error = nil, want error")
	}
}

func TestHasDotDotPrefix(t *testing.T) {
	tests := []struct {
		p    string
		want bool
	}{
		{"..", true},
		{"../foo", true},
		{"..foo", false},
		{".", false},
		{"./", false},
		{"foo", false},
		{"foo/..", false},
	}

	for _, tt := range tests {
		got := hasDotDotPrefix(tt.p)
		if got != tt.want {
			t.Errorf("hasDotDotPrefix(%q) = %v, want %v", tt.p, got, tt.want)
		}
	}
}
