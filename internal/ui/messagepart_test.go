package ui

import (
	"strings"
	"testing"
)

func TestRenderThinkingBlock(t *testing.T) {
	result := RenderThinkingBlock("分析中...", false, 80)
	if !strings.Contains(result, "🤔 Thinking") {
		t.Error("should contain thinking prefix")
	}
}

func TestRenderThinkingBlockCollapsed(t *testing.T) {
	result := RenderThinkingBlock("分析中...", true, 80)
	if !strings.Contains(result, "···") {
		t.Error("collapsed should show ellipsis")
	}
	if !strings.Contains(result, "click to expand") {
		t.Error("collapsed should show expand hint")
	}
}

func TestIsDiffContentTrue(t *testing.T) {
	content := "+++ b/file.go\n- removed line\n+ added line"
	if !isDiffContent(content) {
		t.Error("content with +++ and --- should be diff")
	}
}

func TestIsDiffContentFalse(t *testing.T) {
	content := "This is just regular text"
	if isDiffContent(content) {
		t.Error("regular text should not be diff")
	}
}

func TestIsDiffContentWithPluses(t *testing.T) {
	content := " + added line\n - removed line"
	if !isDiffContent(content) {
		t.Error("content with + and - prefixes should be diff")
	}
}

func TestIsDiffContentWithMinuses(t *testing.T) {
	content := " - removed line\n + added line"
	if !isDiffContent(content) {
		t.Error("content with - and + prefixes should be diff")
	}
}

func TestRenderTextPartWithDiff(t *testing.T) {
	content := "```diff\n+++ b/file.go\n- removed\n+ added\n```"
	result := RenderTextPart(content, 80)
	if !strings.Contains(result, "diff") {
		t.Error("should contain diff marker")
	}
}

func TestRenderTextPartWithCode(t *testing.T) {
	content := "```go\nfunc main() {}\n```"
	result := RenderTextPart(content, 80)
	if !strings.Contains(result, "go") {
		t.Error("should contain language marker")
	}
}

func TestRenderTextPartWithInlineCode(t *testing.T) {
	content := "Use `fmt.Println()` to print"
	result := RenderTextPart(content, 80)
	if !strings.Contains(result, "fmt.Println") {
		t.Error("should render inline code")
	}
}

func TestRenderReasoningPart(t *testing.T) {
	result := RenderReasoningPart("分析中...", 80)
	if !strings.Contains(result, "🤔 Thinking") {
		t.Error("should contain thinking prefix")
	}
}
