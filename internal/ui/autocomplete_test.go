package ui

import (
	"path/filepath"
	"testing"
)

func TestFrecencyStoreRecord(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)

	store.Record("test prompt")
	store.Save()

	suggestions := store.GetSuggestions("test", 10)
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Text != "test prompt" {
		t.Errorf("Expected 'test prompt', got %q", suggestions[0].Text)
	}
	if suggestions[0].UseCount != 1 {
		t.Errorf("Expected UseCount 1, got %d", suggestions[0].UseCount)
	}
}

func TestFrecencyStoreMultipleRecords(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)

	store.Record("prompt one")
	store.Record("prompt two")
	store.Record("prompt one")
	store.Save()

	suggestions := store.GetSuggestions("", 10)
	if len(suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(suggestions))
	}

	for _, s := range suggestions {
		if s.Text == "prompt one" && s.UseCount != 2 {
			t.Errorf("Expected 'prompt one' to have UseCount 2, got %d", s.UseCount)
		}
	}
}

func TestFrecencyStoreFuzzyMatch(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)

	store.Record("Fix bug in authentication")
	store.Record("Explain this code")
	store.Record("Refactor the parser")

	suggestions := store.GetSuggestions("fix", 10)
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion for 'fix', got %d", len(suggestions))
	}

	suggestions = store.GetSuggestions("auth", 10)
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion for 'auth', got %d", len(suggestions))
	}

	suggestions = store.GetSuggestions("xyz", 10)
	if len(suggestions) != 0 {
		t.Errorf("Expected 0 suggestions for 'xyz', got %d", len(suggestions))
	}
}

func TestFrecencyStoreFuzzyMatchSubsequence(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)

	store.Record("Fix authentication bug")

	suggestions := store.GetSuggestions("fab", 10)
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion for 'fab' (f=abug), got %d", len(suggestions))
	}
}

func TestFrecencyStoreScoreCalculation(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)

	store.Record("recent prompt")

	suggestions := store.GetSuggestions("", 10)
	if len(suggestions) != 1 {
		t.Fatal("Expected 1 suggestion")
	}

	initialScore := suggestions[0].Score

	store.Record("recent prompt")
	store.Save()

	suggestions = store.GetSuggestions("", 10)
	if suggestions[0].Score <= initialScore {
		t.Errorf("Expected score to increase after second record, got %f <= %f", suggestions[0].Score, initialScore)
	}
}

func TestFrecencyStoreClear(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)

	store.Record("prompt one")
	store.Record("prompt two")
	store.Save()

	store.Clear()

	suggestions := store.GetSuggestions("", 10)
	if len(suggestions) != 0 {
		t.Errorf("Expected 0 suggestions after clear, got %d", len(suggestions))
	}
}

func TestFrecencyStorePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")

	store1 := NewFrecencyStore(storePath)
	store1.Record("persistent prompt")
	store1.Save()

	store2 := NewFrecencyStore(storePath)
	suggestions := store2.GetSuggestions("persistent", 10)
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion after reload, got %d", len(suggestions))
	}
}

func TestFrecencyStoreMaxItems(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)
	store.maxItems = 5

	for i := 0; i < 10; i++ {
		store.Record("prompt number " + string(rune('a'+i)))
	}

	suggestions := store.GetSuggestions("", 10)
	if len(suggestions) > 5 {
		t.Errorf("Expected max 5 suggestions, got %d", len(suggestions))
	}
}

func TestFrecencyStoreEmptyPrompt(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "autocomplete.json")
	store := NewFrecencyStore(storePath)

	store.Record("")
	store.Record("   ")
	store.Record("valid prompt")

	suggestions := store.GetSuggestions("", 10)
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 valid suggestion, got %d", len(suggestions))
	}
}

func TestAutocompleteDialogShowHide(t *testing.T) {
	dialog := NewAutocompleteDialog()

	if dialog.IsVisible() {
		t.Error("Expected dialog to be initially hidden")
	}

	dialog.Show("test")
	if !dialog.IsVisible() {
		t.Error("Expected dialog to be visible after Show()")
	}

	dialog.Hide()
	if dialog.IsVisible() {
		t.Error("Expected dialog to be hidden after Hide()")
	}
}

func TestAutocompleteDialogNavigation(t *testing.T) {
	dialog := NewAutocompleteDialog()

	suggestions := []*Suggestion{
		{Text: "first", Score: 3.0, UseCount: 3},
		{Text: "second", Score: 2.0, UseCount: 2},
		{Text: "third", Score: 1.0, UseCount: 1},
	}
	dialog.SetItems(suggestions)
	dialog.Show("")

	dialog.Next()
	dialog.Next()
	dialog.Next()

	selected := dialog.GetSelectedText()
	if selected != "first" {
		t.Errorf("Expected 'first' after cycling through, got %q", selected)
	}

	dialog.Prev()
	selected = dialog.GetSelectedText()
	if selected != "third" {
		t.Errorf("Expected 'third' after Prev(), got %q", selected)
	}
}

func TestAutocompleteDialogComplete(t *testing.T) {
	dialog := NewAutocompleteDialog()

	suggestions := []*Suggestion{
		{Text: "Fix authentication bug", Score: 2.0, UseCount: 2},
	}
	dialog.SetItems(suggestions)
	dialog.Show("")

	dialog.SetOnComplete(func(text string) {
		if text != "Fix authentication bug" {
			t.Errorf("Expected 'Fix authentication bug', got %q", text)
		}
	})

	completed := dialog.Complete()
	if completed != "Fix authentication bug" {
		t.Errorf("Expected 'Fix authentication bug', got %q", completed)
	}
}

func TestHighlightMatchNotFound(t *testing.T) {
	label := "Fix authentication bug"
	partial := "xyz"

	result := HighlightMatch(label, partial)
	if result != label {
		t.Error("Expected no highlight when partial not found")
	}
}

func TestHighlightMatchNoPartial(t *testing.T) {
	label := "Fix authentication bug"
	result := HighlightMatch(label, "")
	if result != label {
		t.Error("Expected no modification when partial is empty")
	}
}