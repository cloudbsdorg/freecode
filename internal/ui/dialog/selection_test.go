package dialog

import (
	"testing"
)

func TestColorsDark(t *testing.T) {
	c := Dark
	if c.Background != "#1E1E1E" {
		t.Errorf("expected #1E1E1E, got %s", c.Background)
	}
	if c.Primary != "#22C55E" {
		t.Errorf("expected #22C55E, got %s", c.Primary)
	}
}

func TestColorsLight(t *testing.T) {
	c := Light
	if c.Background != "#FFFFFF" {
		t.Errorf("expected #FFFFFF, got %s", c.Background)
	}
	if c.Primary != "#0066CC" {
		t.Errorf("expected #0066CC, got %s", c.Primary)
	}
}

func TestDefaultColors(t *testing.T) {
	c := DefaultColors
	if c.Background != "#1E1E1E" {
		t.Errorf("expected dark background, got %s", c.Background)
	}
}

func TestItem(t *testing.T) {
	item := Item{
		ID:          "1",
		Title:       "Test",
		Description: "desc",
		Category:    "cat",
		Disabled:    false,
		Footer:      "foot",
		Value:       "val",
	}
	if item.ID != "1" {
		t.Errorf("expected ID=1, got %s", item.ID)
	}
}

func TestNewSelectionList(t *testing.T) {
	s := NewSelectionList()
	if s.Width != 60 {
		t.Errorf("expected Width=60, got %d", s.Width)
	}
	if s.Height != 20 {
		t.Errorf("expected Height=20, got %d", s.Height)
	}
	if s.Selected != 0 {
		t.Errorf("expected Selected=0, got %d", s.Selected)
	}
	if s.Colors != DefaultColors {
		t.Error("expected DefaultColors")
	}
	if s.Title != "Select" {
		t.Errorf("expected Title=Select, got %s", s.Title)
	}
}

func TestSelectionListSetItems(t *testing.T) {
	s := NewSelectionList()
	items := []Item{
		{ID: "1", Title: "One"},
		{ID: "2", Title: "Two"},
	}
	s.SetItems(items)
	if len(s.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(s.Items))
	}
	if len(s.Filtered) != 2 {
		t.Errorf("expected 2 filtered items, got %d", len(s.Filtered))
	}
}

func TestSelectionListSetFilter(t *testing.T) {
	s := NewSelectionList()
	s.Items = []Item{
		{ID: "1", Title: "Apple"},
		{ID: "2", Title: "Banana"},
		{ID: "3", Title: "Cherry"},
	}
	s.ApplyFilter()

	filterCalled := false
	s.OnFilter = func(f string) {
		filterCalled = true
	}

	s.SetFilter("ap")
	if len(s.Filtered) != 1 {
		t.Errorf("expected 1 filtered item, got %d", len(s.Filtered))
	}
	if !filterCalled {
		t.Error("expected OnFilter to be called")
	}
}

func TestSelectionListApplyFilterNoFilter(t *testing.T) {
	s := NewSelectionList()
	s.SkipFilter = true
	s.Items = []Item{
		{ID: "1", Title: "One", Disabled: false},
		{ID: "2", Title: "Two", Disabled: true},
	}
	s.ApplyFilter()
	if len(s.Filtered) != 1 {
		t.Errorf("expected 1 filtered (non-disabled), got %d", len(s.Filtered))
	}
}

func TestSelectionListApplyFilterWithFilter(t *testing.T) {
	s := NewSelectionList()
	s.Items = []Item{
		{ID: "1", Title: "Apple", Description: "fruit"},
		{ID: "2", Title: "Banana", Category: "fruit"},
		{ID: "3", Title: "Cherry"},
	}
	s.SetFilter("fruit")

	if len(s.Filtered) != 2 {
		t.Errorf("expected 2 filtered items, got %d", len(s.Filtered))
	}
}

func TestSelectionListNext(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One"},
		{ID: "2", Title: "Two"},
		{ID: "3", Title: "Three"},
	}
	s.Selected = 0
	s.Height = 10

	s.Next()
	if s.Selected != 1 {
		t.Errorf("expected Selected=1, got %d", s.Selected)
	}

	s.Next()
	if s.Selected != 2 {
		t.Errorf("expected Selected=2, got %d", s.Selected)
	}

	s.Next()
	if s.Selected != 0 {
		t.Errorf("expected wrap to Selected=0, got %d", s.Selected)
	}
}

func TestSelectionListNextEmpty(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{}
	s.Selected = 0

	s.Next()
	if s.Selected != 0 {
		t.Errorf("expected Selected=0, got %d", s.Selected)
	}
}

func TestSelectionListPrev(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One"},
		{ID: "2", Title: "Two"},
	}
	s.Selected = 1
	s.Height = 10

	s.Prev()
	if s.Selected != 0 {
		t.Errorf("expected Selected=0, got %d", s.Selected)
	}

	s.Prev()
	if s.Selected != 1 {
		t.Errorf("expected wrap to Selected=1, got %d", s.Selected)
	}
}

func TestSelectionListPrevEmpty(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{}
	s.Selected = 0

	s.Prev()
	if s.Selected != 0 {
		t.Errorf("expected Selected=0, got %d", s.Selected)
	}
}

func TestSelectionListMoveUp(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One"},
		{ID: "2", Title: "Two"},
	}
	s.Selected = 1
	s.Height = 10

	s.MoveUp()
	if s.Selected != 0 {
		t.Errorf("expected Selected=0, got %d", s.Selected)
	}

	s.MoveUp()
	if s.Selected != 0 {
		t.Errorf("expected Selected=0 at boundary, got %d", s.Selected)
	}
}

func TestSelectionListMoveDown(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One"},
		{ID: "2", Title: "Two"},
	}
	s.Selected = 0
	s.Height = 10

	s.MoveDown()
	if s.Selected != 1 {
		t.Errorf("expected Selected=1, got %d", s.Selected)
	}

	s.MoveDown()
	if s.Selected != 1 {
		t.Errorf("expected Selected=1 at boundary, got %d", s.Selected)
	}
}

func TestSelectionListAutoScroll(t *testing.T) {
	s := NewSelectionList()
	s.Height = 10
	s.ScrollOffset = 0
	s.Selected = 0

	s.autoScroll()
	if s.ScrollOffset != 0 {
		t.Errorf("expected ScrollOffset=0, got %d", s.ScrollOffset)
	}

	s.Selected = 8
	s.autoScroll()
	if s.ScrollOffset != 3 {
		t.Errorf("expected ScrollOffset=3, got %d", s.ScrollOffset)
	}
}

func TestSelectionListAutoScrollZeroHeight(t *testing.T) {
	s := NewSelectionList()
	s.Height = 0
	s.autoScroll()
}

func TestSelectionListNotifyMove(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}}
	s.Selected = 0

	called := false
	s.OnMove = func(item Item) {
		called = true
	}

	s.notifyMove()
	if !called {
		t.Error("expected OnMove to be called")
	}
}

func TestSelectionListNotifyMoveOutOfBounds(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}}
	s.Selected = 10

	called := false
	s.OnMove = func(item Item) {
		called = true
	}

	s.notifyMove()
	if called {
		t.Error("expected OnMove NOT to be called")
	}
}

func TestSelectionListSelect(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One"},
		{ID: "2", Title: "Two"},
	}

	if !s.Select(1) {
		t.Error("expected Select(1) to succeed")
	}
	if s.Selected != 1 {
		t.Errorf("expected Selected=1, got %d", s.Selected)
	}

	if s.Select(5) {
		t.Error("expected Select(5) to fail")
	}
}

func TestSelectionListGetSelected(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One"},
		{ID: "2", Title: "Two"},
	}
	s.Selected = 1

	item := s.GetSelected()
	if item == nil {
		t.Fatal("expected non-nil item")
	}
	if item.ID != "2" {
		t.Errorf("expected ID=2, got %s", item.ID)
	}

	s.Selected = 10
	if s.GetSelected() != nil {
		t.Error("expected nil for out-of-bounds")
	}
}

func TestSelectionListSelectByValue(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One", Value: "val1"},
		{ID: "2", Title: "Two", Value: "val2"},
	}

	if !s.SelectByValue("val2") {
		t.Error("expected SelectByValue to succeed")
	}
	if s.Selected != 1 {
		t.Errorf("expected Selected=1, got %d", s.Selected)
	}

	if s.SelectByValue("val3") {
		t.Error("expected SelectByValue to fail")
	}
}

func TestSelectionListConfirm(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}}
	s.Selected = 0

	called := false
	s.OnSelect = func(item Item) {
		called = true
	}

	if !s.Confirm() {
		t.Error("expected Confirm to succeed")
	}
	if !called {
		t.Error("expected OnSelect to be called")
	}
}

func TestSelectionListConfirmNoCallback(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}}
	s.Selected = 0
	s.OnSelect = nil

	if s.Confirm() {
		t.Error("expected Confirm to fail without callback")
	}
}

func TestSelectionListIsItemCurrent(t *testing.T) {
	s := NewSelectionList()
	s.Current = "val1"

	item := Item{ID: "1", Value: "val1"}
	if !s.IsItemCurrent(item) {
		t.Error("expected IsItemCurrent to return true")
	}

	item.Value = "val2"
	if s.IsItemCurrent(item) {
		t.Error("expected IsItemCurrent to return false")
	}
}

func TestSelectionListIsItemCurrentNilCurrent(t *testing.T) {
	s := NewSelectionList()
	s.Current = nil

	item := Item{ID: "1"}
	if s.IsItemCurrent(item) {
		t.Error("expected false with nil Current")
	}
}

func TestSelectionListRenderList(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}}
	s.Height = 10

	lines := s.RenderList()
	if len(lines) == 0 {
		t.Error("expected non-empty lines")
	}
}

func TestSelectionListRenderListWithFilter(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}}
	s.Filter = "filter"
	s.Height = 10

	lines := s.RenderList()
	if len(lines) == 0 {
		t.Error("expected non-empty lines")
	}
}

func TestSelectionListRenderListEmpty(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{}
	s.Height = 10

	lines := s.RenderList()
	if len(lines) != 1 {
		t.Errorf("expected 1 line (no matches), got %d", len(lines))
	}
}

func TestSelectionListRenderVisibleItems(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}, {ID: "2", Title: "Two"}}
	s.Height = 10

	lines := s.renderVisibleItems()
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestSelectionListRenderVisibleItemsScrolled(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}, {ID: "2", Title: "Two"}, {ID: "3", Title: "Three"}}
	s.Height = 10
	s.ScrollOffset = 1

	lines := s.renderVisibleItems()
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestSelectionListRender(t *testing.T) {
	s := NewSelectionList()
	s.Items = []Item{{ID: "1", Title: "One"}}
	s.ApplyFilter()
	s.Height = 10

	result := s.Render()
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestSelectionListRenderItemsFlat(t *testing.T) {
	s := NewSelectionList()
	s.Flat = true
	s.Filtered = []Item{{ID: "1", Title: "One"}}

	lines := s.renderItems()
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
}

func TestSelectionListRenderItemsGrouped(t *testing.T) {
	s := NewSelectionList()
	s.Flat = false
	s.Filtered = []Item{
		{ID: "1", Title: "One", Category: "Fruit"},
		{ID: "2", Title: "Two", Category: "Fruit"},
	}

	lines := s.renderItems()
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 items), got %d", len(lines))
	}
}

func TestSelectionListHasCategories(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{
		{ID: "1", Title: "One", Category: ""},
		{ID: "2", Title: "Two", Category: ""},
	}

	if s.HasCategories() {
		t.Error("expected false for no categories")
	}

	s.Filtered[0].Category = "Fruit"
	if !s.HasCategories() {
		t.Error("expected true with category")
	}
}

func TestSelectionListIsSelected(t *testing.T) {
	s := NewSelectionList()
	s.Filtered = []Item{{ID: "1", Title: "One"}, {ID: "2", Title: "Two"}}
	s.Selected = 1

	item := Item{ID: "1"}
	if s.IsSelected(item) {
		t.Error("expected IsSelected false for ID=1 when Selected=1")
	}

	item.ID = "2"
	if !s.IsSelected(item) {
		t.Error("expected IsSelected true for ID=2 when Selected=1")
	}
}

func TestStyle(t *testing.T) {
	s := Style(Dark)
	_ = s
}

func TestStyled(t *testing.T) {
	s := Styled(Dark, 50, 20)
	_ = s
}

func TestHeader(t *testing.T) {
	result := Header("Title", Dark)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestMuted(t *testing.T) {
	result := Muted("text", Dark)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestTextStyled(t *testing.T) {
	result := TextStyled("text", "#FF0000")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestErrorText(t *testing.T) {
	result := ErrorText("error", Dark)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestSuccessText(t *testing.T) {
	result := SuccessText("success", Dark)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestInput(t *testing.T) {
	result := Input("text", Dark)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestSelected(t *testing.T) {
	result := Selected("text", Dark)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestSelectedPrefix(t *testing.T) {
	result := SelectedPrefix("▶", Dark)
	if result == "" {
		t.Error("expected non-empty result")
	}
}