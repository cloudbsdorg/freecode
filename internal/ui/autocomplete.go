package ui

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/freecode/freecode/internal/style"
)

const (
	minCharsForAutocomplete = 2
	maxAutocompleteItems    = 5
)

type FrecencyStore struct {
	mu         sync.RWMutex
	suggestions map[string]*Suggestion
	filePath   string
	maxItems   int
}

type Suggestion struct {
	Text       string  `json:"text"`
	Score      float64 `json:"score"`
	UseCount   int     `json:"useCount"`
	LastUsedAt int64   `json:"lastUsedAt"`
}

func NewFrecencyStore(filePath string) *FrecencyStore {
	store := &FrecencyStore{
		suggestions: make(map[string]*Suggestion),
		filePath:    filePath,
		maxItems:    100,
	}
	store.load()
	return store
}

func (f *FrecencyStore) Record(prompt string) {
	if prompt == "" {
		return
	}

	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now().Unix()
	s, exists := f.suggestions[prompt]
	if exists {
		s.UseCount++
		s.LastUsedAt = now
		s.Score = f.calculateScore(s)
	} else {
		f.suggestions[prompt] = &Suggestion{
			Text:       prompt,
			UseCount:   1,
			LastUsedAt: now,
			Score:      1.0,
		}
		if len(f.suggestions) > f.maxItems {
			f.evictLowestScored()
		}
	}
}

func (f *FrecencyStore) GetSuggestions(partial string, limit int) []*Suggestion {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if limit <= 0 {
		limit = 5
	}

	var results []*Suggestion
	needle := strings.ToLower(strings.TrimSpace(partial))

	for _, s := range f.suggestions {
		if partial == "" || f.fuzzyMatch(s.Text, needle) {
			results = append(results, &Suggestion{
				Text:       s.Text,
				Score:      s.Score,
				UseCount:   s.UseCount,
				LastUsedAt: s.LastUsedAt,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func (f *FrecencyStore) fuzzyMatch(haystack, needle string) bool {
	if needle == "" {
		return true
	}

	haystack = strings.ToLower(haystack)

	if strings.Contains(haystack, needle) {
		return true
	}

	needleRunes := []rune(needle)
	haystackRunes := []rune(haystack)

	j := 0
	for i := 0; i < len(haystackRunes) && j < len(needleRunes); i++ {
		if haystackRunes[i] == needleRunes[j] {
			j++
		}
	}
	return j == len(needleRunes)
}

func (f *FrecencyStore) calculateScore(s *Suggestion) float64 {
	now := time.Now().Unix()
	hoursSinceLastUse := float64(now-s.LastUsedAt) / 3600.0
	return float64(s.UseCount) * (1 + math.Log(hoursSinceLastUse+1))
}

func (f *FrecencyStore) evictLowestScored() {
	var lowest *Suggestion
	var lowestKey string
	for k, s := range f.suggestions {
		if lowest == nil || s.Score < lowest.Score {
			lowest = s
			lowestKey = k
		}
	}
	if lowestKey != "" {
		delete(f.suggestions, lowestKey)
	}
}

func (f *FrecencyStore) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.suggestions = make(map[string]*Suggestion)
}

func (f *FrecencyStore) Save() error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	data, err := json.MarshalIndent(f.suggestions, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(f.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(f.filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

func (f *FrecencyStore) load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := os.ReadFile(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var suggestions map[string]*Suggestion
	if err := json.Unmarshal(data, &suggestions); err != nil {
		return err
	}

	f.suggestions = suggestions
	return nil
}

type AutocompleteItem struct {
	Label       string
	Value       string
	Description string
	Category    string
	Score       float64
}

type AutocompleteDialog struct {
	visible    bool
	items      []AutocompleteItem
	filtered   []AutocompleteItem
	selected   int
	input      string
	width      int
	onSelect   func(string)
	onClose    func()
	onComplete func(string)
}

func NewAutocompleteDialog() *AutocompleteDialog {
	return &AutocompleteDialog{
		visible:  false,
		items:    []AutocompleteItem{},
		filtered: []AutocompleteItem{},
		selected: 0,
		input:   "",
		width:   50,
	}
}

func (a *AutocompleteDialog) SetWidth(w int) {
	a.width = w
}

func (a *AutocompleteDialog) SetItems(suggestions []*Suggestion) {
	a.items = make([]AutocompleteItem, len(suggestions))
	for i, s := range suggestions {
		a.items[i] = AutocompleteItem{
			Label:       s.Text,
			Value:       s.Text,
			Description: fmt.Sprintf("Used %d time(s)", s.UseCount),
			Category:    "history",
			Score:       s.Score,
		}
	}
	a.filtered = a.items
	a.selected = 0
}

func (a *AutocompleteDialog) Show(input string) {
	a.visible = true
	a.input = input
	a.applyFilter()
}

func (a *AutocompleteDialog) Hide() {
	a.visible = false
}

func (a *AutocompleteDialog) IsVisible() bool {
	return a.visible
}

func (a *AutocompleteDialog) SetOnSelect(fn func(string)) {
	a.onSelect = fn
}

func (a *AutocompleteDialog) SetOnClose(fn func()) {
	a.onClose = fn
}

func (a *AutocompleteDialog) SetOnComplete(fn func(string)) {
	a.onComplete = fn
}

func (a *AutocompleteDialog) GetSelectedText() string {
	if a.selected < 0 || a.selected >= len(a.filtered) {
		return ""
	}
	return a.filtered[a.selected].Label
}

func (a *AutocompleteDialog) applyFilter() {
	if a.input == "" {
		a.filtered = a.items
		a.selected = 0
		return
	}

	needle := strings.ToLower(a.input)
	var result []AutocompleteItem
	for _, item := range a.items {
		if a.fuzzyMatch(item.Label, needle) {
			result = append(result, item)
		} else if strings.Contains(strings.ToLower(item.Description), needle) {
			result = append(result, item)
		}
	}
	a.filtered = result
	if a.selected >= len(a.filtered) {
		a.selected = 0
	}
}

func (a *AutocompleteDialog) fuzzyMatch(haystack, needle string) bool {
	if needle == "" {
		return true
	}

	haystack = strings.ToLower(haystack)

	if strings.Contains(haystack, needle) {
		return true
	}

	needleRunes := []rune(needle)
	haystackRunes := []rune(haystack)

	j := 0
	for i := 0; i < len(haystackRunes) && j < len(needleRunes); i++ {
		if haystackRunes[i] == needleRunes[j] {
			j++
		}
	}
	return j == len(needleRunes)
}

func (a *AutocompleteDialog) Next() {
	if len(a.filtered) == 0 {
		return
	}
	a.selected = (a.selected + 1) % len(a.filtered)
}

func (a *AutocompleteDialog) Prev() {
	if len(a.filtered) == 0 {
		return
	}
	a.selected = (a.selected - 1 + len(a.filtered)) % len(a.filtered)
}

func (a *AutocompleteDialog) Select() *AutocompleteItem {
	if a.selected < 0 || a.selected >= len(a.filtered) {
		return nil
	}
	item := &a.filtered[a.selected]
	if a.onSelect != nil {
		a.onSelect(item.Value)
	}
	a.Hide()
	return item
}

func (a *AutocompleteDialog) Complete() string {
	if a.selected < 0 || a.selected >= len(a.filtered) {
		return ""
	}
	item := &a.filtered[a.selected]
	if a.onComplete != nil {
		a.onComplete(item.Value)
	}
	a.Hide()
	return item.Label
}

func (a *AutocompleteDialog) HandleKey(msg string) bool {
	if !a.visible {
		return false
	}

	switch msg {
	case "escape":
		a.Hide()
		if a.onClose != nil {
			a.onClose()
		}
		return true
	case "enter":
		a.Select()
		return true
	case "tab":
		a.Complete()
		return true
	case "up", "k":
		a.Prev()
		return true
	case "down", "j":
		a.Next()
		return true
	}
	return false
}

func (a *AutocompleteDialog) Render() string {
	if !a.visible || len(a.filtered) == 0 {
		return ""
	}

	dialogStyle := style.NewStyle().
		Background(style.Color("#2D2D2D")).
		BorderStyle(style.HiddenBorder()).
		Width(a.width)

	return dialogStyle.Render(a.renderContent())
}

func (a *AutocompleteDialog) renderContent() string {
	var lines []string

	for i, item := range a.filtered {
		lines = append(lines, a.renderItem(item, i == a.selected))
	}

	return strings.Join(lines, "\n")
}

func (a *AutocompleteDialog) renderItem(item AutocompleteItem, selected bool) string {
	prefix := "  "
	if selected {
		prefix = style.NewStyle().
			Foreground(style.Color("#FFCC00")).
			Render("▶")
	}

	var labelStyle style.Style
	if selected {
		labelStyle = style.NewStyle().Foreground(style.Color("#FFFFFF"))
	} else {
		labelStyle = style.NewStyle().Foreground(style.Color("#D0D0D0"))
	}

	line := prefix + " " + labelStyle.Render(item.Label)

	if item.Category != "" {
		catStyle := style.NewStyle().Foreground(style.Color("#007ACC"))
		line += " " + catStyle.Render("["+item.Category+"]")
	}

	if item.Description != "" {
		descStyle := style.NewStyle().Foreground(style.Color("#808080"))
		line += " — " + descStyle.Render(item.Description)
	}

	return line
}

func HighlightMatch(label, partial string) string {
	if partial == "" {
		return label
	}

	lowerLabel := strings.ToLower(label)
	lowerPartial := strings.ToLower(partial)

	idx := strings.Index(lowerLabel, lowerPartial)
	if idx == -1 {
		return label
	}

	before := label[:idx]
	match := label[idx : idx+len(partial)]
	after := label[idx+len(partial):]

	highlightStyle := style.NewStyle().
		Foreground(style.Color("#FFCC00")).
		Bold(true)

	return before + highlightStyle.Render(match) + after
}
