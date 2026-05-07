package session

import (
	"testing"
)

func TestNewCompactor(t *testing.T) {
	c := NewCompactor(1000, 10)

	if c.thresholdTokens != 1000 {
		t.Errorf("thresholdTokens = %d, want 1000", c.thresholdTokens)
	}

	if c.keepLastN != 10 {
		t.Errorf("keepLastN = %d, want 10", c.keepLastN)
	}
}

func TestCompactorCompactNoOpBelowThreshold(t *testing.T) {
	c := NewCompactor(1000, 5)
	sess := &Session{
		Messages: []Message{
			{Role: "user", Content: "short message"},
		},
	}

	err := c.Compact(sess)
	if err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	if len(sess.Messages) != 1 {
		t.Errorf("len(Messages) = %d, want 1", len(sess.Messages))
	}
}

func TestCompactorCompactNoOpBelowKeepLastN(t *testing.T) {
	c := NewCompactor(1, 10)
	sess := &Session{
		Messages: []Message{
			{Role: "user", Content: "msg1"},
			{Role: "user", Content: "msg2"},
			{Role: "user", Content: "msg3"},
		},
	}

	err := c.Compact(sess)
	if err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	if len(sess.Messages) != 3 {
		t.Errorf("len(Messages) = %d, want 3", len(sess.Messages))
	}
}

func TestCompactorCompact(t *testing.T) {
	c := NewCompactor(1, 2)
	sess := &Session{
		Messages: []Message{
			{Role: "system", Content: "system prompt"},
			{Role: "user", Content: "message 1"},
			{Role: "assistant", Content: "message 2"},
			{Role: "user", Content: "message 3"},
			{Role: "assistant", Content: "message 4"},
		},
	}

	err := c.Compact(sess)
	if err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	if len(sess.Messages) != 3 {
		t.Errorf("len(Messages) = %d, want 3 (system + 2 keep)", len(sess.Messages))
	}

	if sess.Messages[0].Role != "system" {
		t.Errorf("Messages[0].Role = %q, want %q", sess.Messages[0].Role, "system")
	}
}

func TestCompactorExtractSystemMessages(t *testing.T) {
	c := NewCompactor(1000, 10)
	messages := []Message{
		{Role: "system", Content: "sys1"},
		{Role: "user", Content: "user1"},
		{Role: "system", Content: "sys2"},
		{Role: "assistant", Content: "asst1"},
	}

	system := c.extractSystemMessages(messages)

	if len(system) != 2 {
		t.Errorf("len(system) = %d, want 2", len(system))
	}

	if system[0].Content != "sys1" {
		t.Errorf("system[0].Content = %q, want %q", system[0].Content, "sys1")
	}

	if system[1].Content != "sys2" {
		t.Errorf("system[1].Content = %q, want %q", system[1].Content, "sys2")
	}
}

func TestCompactorEstimateTokens(t *testing.T) {
	c := NewCompactor(1000, 10)
	messages := []Message{
		{Role: "user", Content: "one two three four five"},
		{Role: "assistant", Content: "six seven eight nine ten"},
	}

	tokens := c.estimateTokens(messages)

	if tokens == 0 {
		t.Error("estimateTokens() returned 0")
	}
}

func TestCompactorEstimateTokensEmpty(t *testing.T) {
	c := NewCompactor(1000, 10)
	messages := []Message{}

	tokens := c.estimateTokens(messages)

	if tokens != 0 {
		t.Errorf("estimateTokens() = %d, want 0", tokens)
	}
}

func TestNewHistory(t *testing.T) {
	h := NewHistory(100)

	if h.maxSize != 100 {
		t.Errorf("maxSize = %d, want 100", h.maxSize)
	}
}

func TestHistoryAddEmpty(t *testing.T) {
	h := NewHistory(10)
	sessions := []*Session{}

	result := h.Add(sessions, Message{Role: "user", Content: "test"})

	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}
}

func TestHistoryAdd(t *testing.T) {
	h := NewHistory(3)
	sessions := []*Session{
		{Messages: []Message{}},
	}

	result := h.Add(sessions, Message{Role: "user", Content: "msg1"})

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}

	if len(result[0].Messages) != 1 {
		t.Errorf("len(result[0].Messages) = %d, want 1", len(result[0].Messages))
	}
}

func TestHistoryAddTruncates(t *testing.T) {
	h := NewHistory(2)
	sessions := []*Session{
		{Messages: []Message{
			{Role: "user", Content: "old1"},
			{Role: "user", Content: "old2"},
		}},
	}

	result := h.Add(sessions, Message{Role: "user", Content: "new"})

	if len(result[0].Messages) != 2 {
		t.Errorf("len(result[0].Messages) = %d, want 2", len(result[0].Messages))
	}
}

func TestHistorySearch(t *testing.T) {
	h := NewHistory(100)
	sessions := []*Session{
		{
			Messages: []Message{
				{Role: "user", Content: "Hello world"},
				{Role: "assistant", Content: "Hi there"},
			},
		},
		{
			Messages: []Message{
				{Role: "user", Content: "Goodbye world"},
			},
		},
	}

	results := h.Search(sessions, "world")

	if len(results) != 2 {
		t.Errorf("len(results) = %d, want 2", len(results))
	}
}

func TestHistorySearchNoMatch(t *testing.T) {
	h := NewHistory(100)
	sessions := []*Session{
		{
			Messages: []Message{
				{Role: "user", Content: "Hello world"},
			},
		},
	}

	results := h.Search(sessions, "nonexistent")

	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestHistorySearchCaseInsensitive(t *testing.T) {
	h := NewHistory(100)
	sessions := []*Session{
		{
			Messages: []Message{
				{Role: "user", Content: "Hello WORLD"},
			},
		},
	}

	results := h.Search(sessions, "world")

	if len(results) != 1 {
		t.Errorf("len(results) = %d, want 1", len(results))
	}
}
