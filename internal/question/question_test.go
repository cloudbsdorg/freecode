package question

import (
	"testing"
	"time"
)

func TestNewFlow(t *testing.T) {
	flow := NewFlow("test-flow")
	if flow == nil {
		t.Fatal("NewFlow() returned nil")
	}
	if flow.id != "test-flow" {
		t.Errorf("id = %q, want %q", flow.id, "test-flow")
	}
}

func TestFlowAddQuestion(t *testing.T) {
	flow := NewFlow("test")
	q := &Question{
		ID:      "q1",
		Text:    "What is your name?",
		Options: []string{"Alice", "Bob"},
	}

	flow.AddQuestion(q)

	flow.mu.RLock()
	if len(flow.questions) != 1 {
		t.Errorf("len(questions) = %d, want 1", len(flow.questions))
	}
	flow.mu.RUnlock()
}

func TestFlowNext(t *testing.T) {
	flow := NewFlow("test")
	q := &Question{ID: "q1", Text: "Question 1"}
	flow.AddQuestion(q)

	flow.mu.Lock()
	flow.current = 1
	flow.mu.Unlock()

	next, err := flow.Next()
	if err == nil {
		t.Error("Next() should return error when no more questions")
	}
	if next != nil {
		t.Error("Next() should return nil when no more questions")
	}
}

func TestFlowNextQuestion(t *testing.T) {
	flow := NewFlow("test")
	q := &Question{ID: "q1", Text: "Question 1"}
	flow.AddQuestion(q)

	flow.mu.Lock()
	flow.current = 0
	flow.mu.Unlock()

	next, err := flow.Next()
	if err != nil {
		t.Fatalf("Next() error = %v", err)
	}
	if next == nil {
		t.Fatal("Next() returned nil")
	}
	if next.ID != "q1" {
		t.Errorf("next.ID = %q, want %q", next.ID, "q1")
	}
}

func TestFlowAnswer(t *testing.T) {
	flow := NewFlow("test")
	q := &Question{ID: "q1", Text: "Question 1"}
	flow.AddQuestion(q)

	flow.mu.Lock()
	flow.current = 1
	flow.mu.Unlock()

	flow.Answer("q1", []string{"answer"}, "answer")

	flow.mu.RLock()
	answer := flow.answers["q1"]
	flow.mu.RUnlock()

	if answer == nil {
		t.Fatal("Answer not found")
	}
	if answer.QuestionID != "q1" {
		t.Errorf("QuestionID = %q, want %q", answer.QuestionID, "q1")
	}
}

func TestQuestion(t *testing.T) {
	q := Question{
		ID:          "q1",
		Text:        "What is your name?",
		Options:     []string{"Alice", "Bob"},
		MultiSelect: false,
		Required:    true,
		Timeout:     30 * time.Second,
	}

	if q.ID != "q1" {
		t.Errorf("ID = %q, want %q", q.ID, "q1")
	}
	if q.Text != "What is your name?" {
		t.Errorf("Text = %q, want %q", q.Text, "What is your name?")
	}
	if len(q.Options) != 2 {
		t.Errorf("len(Options) = %d, want 2", len(q.Options))
	}
	if !q.Required {
		t.Error("Required should be true")
	}
	if q.MultiSelect {
		t.Error("MultiSelect should be false")
	}
}

func TestAnswer(t *testing.T) {
	now := time.Now()
	a := Answer{
		QuestionID: "q1",
		Selected:   []string{"Alice"},
		Text:      "Alice",
		Timestamp: now,
	}

	if a.QuestionID != "q1" {
		t.Errorf("QuestionID = %q, want %q", a.QuestionID, "q1")
	}
	if len(a.Selected) != 1 {
		t.Errorf("len(Selected) = %d, want 1", len(a.Selected))
	}
	if a.Selected[0] != "Alice" {
		t.Errorf("Selected[0] = %q, want %q", a.Selected[0], "Alice")
	}
}

func TestFlowCurrentIndex(t *testing.T) {
	flow := NewFlow("test")
	flow.AddQuestion(&Question{ID: "q1"})
	flow.AddQuestion(&Question{ID: "q2"})

	flow.mu.RLock()
	current := flow.current
	flow.mu.RUnlock()

	if current != 0 {
		t.Errorf("current = %d, want 0", current)
	}
}