package template

import (
	"strings"
	"testing"

	"github.com/freecode/freecode/internal/renderer"
)

type mockRenderer struct{}

func (m mockRenderer) RenderBox(x, y, w, h int, bgColor string) string {
	return ""
}
func (m mockRenderer) RenderText(text string, x, y int, fgColor string) string {
	return text
}
func (m mockRenderer) RenderBorder(x, y, w, h int, fgColor string) string {
	return ""
}
func (m mockRenderer) RenderSelected(text string, x, y, w int, fg, bg string) string {
	return text
}
func (m mockRenderer) Width() int  { return 80 }
func (m mockRenderer) Height() int { return 24 }

func TestParseBasic(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantType ElementType
	}{
		{"window", "<window>", TypeWindow},
		{"vbox", "<vbox>", TypeVBox},
		{"hbox", "<hbox>", TypeHBox},
		{"text", "<text>", TypeText},
		{"list", "<list>", TypeList},
		{"button", "<button>", TypeButton},
		{"input", "<input>", TypeInput},
		{"spacer", "<spacer>", TypeSpacer},
		{"divider", "<divider>", TypeDivider},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := Parse(tt.src)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.src, err)
			}
			if tmpl.Root.Type != tt.wantType {
				t.Errorf("Parse(%q) type = %v, want %v", tt.src, tmpl.Root.Type, tt.wantType)
			}
		})
	}
}

func TestParseAttributes(t *testing.T) {
	src := `<window title="Setup" width="60" height="20" center="true" padding="2">
		<vbox gap="1">
			<text value="Hello" color="#ffffff" bold="true" />
			<spacer height="1" flex="2" />
			<button id="next" label="Next" primary="true" />
		</vbox>
	</window>`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	if tmpl.Root.Attributes["title"] != "Setup" {
		t.Errorf("title = %q, want %q", tmpl.Root.Attributes["title"], "Setup")
	}
	if tmpl.Root.Attributes["width"] != "60" {
		t.Errorf("width = %q, want %q", tmpl.Root.Attributes["width"], "60")
	}
	if tmpl.Root.Attributes["height"] != "20" {
		t.Errorf("height = %q, want %q", tmpl.Root.Attributes["height"], "20")
	}
	if tmpl.Root.Attributes["center"] != "true" {
		t.Errorf("center = %q, want %q", tmpl.Root.Attributes["center"], "true")
	}

	if len(tmpl.Root.Children) != 1 {
		t.Fatalf("Children len = %d, want 1", len(tmpl.Root.Children))
	}

	vbox := tmpl.Root.Children[0]
	if vbox.Type != TypeVBox {
		t.Errorf("child type = %v, want %v", vbox.Type, TypeVBox)
	}
	if vbox.Attributes["gap"] != "1" {
		t.Errorf("vbox gap = %q, want %q", vbox.Attributes["gap"], "1")
	}

	if len(vbox.Children) != 3 {
		t.Fatalf("vbox.Children len = %d, want 3", len(vbox.Children))
	}

	text := vbox.Children[0]
	if text.Type != TypeText {
		t.Errorf("first child type = %v, want %v", text.Type, TypeText)
	}
	if text.Attributes["value"] != "Hello" {
		t.Errorf("text value = %q, want %q", text.Attributes["value"], "Hello")
	}
	if text.Attributes["color"] != "#ffffff" {
		t.Errorf("text color = %q, want %q", text.Attributes["color"], "#ffffff")
	}
	if text.Attributes["bold"] != "true" {
		t.Errorf("text bold = %q, want %q", text.Attributes["bold"], "true")
	}

	spacer := vbox.Children[1]
	if spacer.Type != TypeSpacer {
		t.Errorf("second child type = %v, want %v", spacer.Type, TypeSpacer)
	}
	if spacer.Attributes["height"] != "1" {
		t.Errorf("spacer height = %q, want %q", spacer.Attributes["height"], "1")
	}
	if spacer.Attributes["flex"] != "2" {
		t.Errorf("spacer flex = %q, want %q", spacer.Attributes["flex"], "2")
	}

	button := vbox.Children[2]
	if button.Type != TypeButton {
		t.Errorf("third child type = %v, want %v", button.Type, TypeButton)
	}
	if button.ID != "next" {
		t.Errorf("button ID = %q, want %q", button.ID, "next")
	}
	if button.Attributes["label"] != "Next" {
		t.Errorf("button label = %q, want %q", button.Attributes["label"], "Next")
	}
	if button.Attributes["primary"] != "true" {
		t.Errorf("button primary = %q, want %q", button.Attributes["primary"], "true")
	}
}

func TestParseVariable(t *testing.T) {
	src := `<vbox>
		<text value="${title}" />
		<list items="${items}" />
	</vbox>`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	if len(tmpl.Root.Children) != 2 {
		t.Fatalf("Children len = %d, want 2", len(tmpl.Root.Children))
	}

	text := tmpl.Root.Children[0]
	if text.Type != TypeText {
		t.Errorf("first child type = %v, want %v", text.Type, TypeText)
	}
	if text.Attributes["value"] != "${title}" {
		t.Errorf("text value attr = %q, want %q", text.Attributes["value"], "${title}")
	}
	if len(text.AttributeVars["value"]) != 1 {
		t.Errorf("text AttributeVars[value] len = %d, want 1", len(text.AttributeVars["value"]))
	}
	if text.AttributeVars["value"][0] != "title" {
		t.Errorf("text AttributeVars[value][0] = %q, want %q", text.AttributeVars["value"][0], "title")
	}

	list := tmpl.Root.Children[1]
	if list.Type != TypeList {
		t.Errorf("second child type = %v, want %v", list.Type, TypeList)
	}
	if list.Attributes["items"] != "${items}" {
		t.Errorf("list items attr = %q, want %q", list.Attributes["items"], "${items}")
	}
	if len(list.AttributeVars["items"]) != 1 {
		t.Errorf("list AttributeVars[items] len = %d, want 1", len(list.AttributeVars["items"]))
	}
	if list.AttributeVars["items"][0] != "items" {
		t.Errorf("list AttributeVars[items][0] = %q, want %q", list.AttributeVars["items"][0], "items")
	}
}

func TestParseVariableProperty(t *testing.T) {
	src := `<text value="${item.name}" />`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	if tmpl.Root.Attributes["value"] != "${item.name}" {
		t.Errorf("value attr = %q, want %q", tmpl.Root.Attributes["value"], "${item.name}")
	}
	if len(tmpl.Root.AttributeVars["value"]) != 1 {
		t.Errorf("AttributeVars[value] len = %d, want 1", len(tmpl.Root.AttributeVars["value"]))
	}
	if tmpl.Root.AttributeVars["value"][0] != "item.name" {
		t.Errorf("AttributeVars[value][0] = %q, want %q", tmpl.Root.AttributeVars["value"][0], "item.name")
	}
}

func TestParseContent(t *testing.T) {
	src := `<text>Hello World</text>`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	if tmpl.Root.Content != "Hello World" {
		t.Errorf("content = %q, want %q", tmpl.Root.Content, "Hello World")
	}
}

func TestParseGrid(t *testing.T) {
	src := `<grid cols="3" rows="2" gap="2">
		<text value="A" />
		<text value="B" />
		<text value="C" />
		<text value="D" />
		<text value="E" />
		<text value="F" />
	</grid>`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	if tmpl.Root.Type != TypeGrid {
		t.Errorf("type = %v, want %v", tmpl.Root.Type, TypeGrid)
	}
	if tmpl.Root.Attributes["cols"] != "3" {
		t.Errorf("cols = %q, want %q", tmpl.Root.Attributes["cols"], "3")
	}
	if tmpl.Root.Attributes["rows"] != "2" {
		t.Errorf("rows = %q, want %q", tmpl.Root.Attributes["rows"], "2")
	}
	if tmpl.Root.Attributes["gap"] != "2" {
		t.Errorf("gap = %q, want %q", tmpl.Root.Attributes["gap"], "2")
	}
	if len(tmpl.Root.Children) != 6 {
		t.Errorf("Children len = %d, want 6", len(tmpl.Root.Children))
	}
}

func TestParseNestedContainers(t *testing.T) {
	src := `<vbox>
		<hbox>
			<text value="Left" />
			<spacer width="1" />
			<text value="Right" />
		</hbox>
		<divider char="=" />
		<text value="Bottom" />
	</vbox>`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	if len(tmpl.Root.Children) != 3 {
		t.Fatalf("Children len = %d, want 3", len(tmpl.Root.Children))
	}

	hbox := tmpl.Root.Children[0]
	if hbox.Type != TypeHBox {
		t.Errorf("first child type = %v, want %v", hbox.Type, TypeHBox)
	}
	if len(hbox.Children) != 3 {
		t.Errorf("hbox Children len = %d, want 3", len(hbox.Children))
	}

	divider := tmpl.Root.Children[1]
	if divider.Type != TypeDivider {
		t.Errorf("second child type = %v, want %v", divider.Type, TypeDivider)
	}
	if divider.Attributes["char"] != "=" {
		t.Errorf("divider char = %q, want %q", divider.Attributes["char"], "=")
	}
}

func TestParseSelfClosing(t *testing.T) {
	src := `<vbox>
		<text value="Title" />
		<spacer height="2" />
		<text value="Body" />
	</vbox>`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	if len(tmpl.Root.Children) != 3 {
		t.Fatalf("Children len = %d, want 3", len(tmpl.Root.Children))
	}

	spacer := tmpl.Root.Children[1]
	if spacer.Type != TypeSpacer {
		t.Errorf("middle child type = %v, want %v", spacer.Type, TypeSpacer)
	}
}

func TestTemplateSetVariables(t *testing.T) {
	src := `<text value="${name}" />`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	tmpl.Set("name", "Freecode")
	tmpl.Set("count", 42)

	if tmpl.Vars["name"] != "Freecode" {
		t.Errorf("name = %v, want %q", tmpl.Vars["name"], "Freecode")
	}
	if tmpl.Vars["count"] != 42 {
		t.Errorf("count = %v, want %d", tmpl.Vars["count"], 42)
	}
}

func TestTemplateSetMap(t *testing.T) {
	src := `<vbox>
		<text value="${title}" />
		<text value="${count}" />
	</vbox>`

	tmpl, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}

	tmpl.SetMap(map[string]interface{}{
		"title": "Setup",
		"count": 10,
	})

	if len(tmpl.Vars) != 2 {
		t.Errorf("Vars len = %d, want 2", len(tmpl.Vars))
	}
	if tmpl.Vars["title"] != "Setup" {
		t.Errorf("title = %v, want %q", tmpl.Vars["title"], "Setup")
	}
}

func TestInterpolate(t *testing.T) {
	vars := map[string]interface{}{
		"name":  "Freecode",
		"count": 42,
		"items": "a,b,c",
	}

	tests := []struct {
		input string
		want  string
	}{
		{"Hello ${name}", "Hello Freecode"},
		{"Count: ${count}", "Count: 42"},
		{"${name} has ${count} items", "Freecode has 42 items"},
		{"Items: ${items}", "Items: a,b,c"},
		{"No substitution ${missing}", "No substitution ${missing}"},
	}

	for _, tt := range tests {
		got := interpolate(tt.input, vars)
		if got != tt.want {
			t.Errorf("interpolate(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestEngineRenderWindow(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<window title="Test" width="40" height="10" padding="1">
		<text value="Hello World" />
	</window>`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}
}

func TestEngineRenderMultiLineText(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<vbox>
		<text value="Line 1\nLine 2\nLine 3" />
	</vbox>`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}

	if !strings.Contains(result, "Line 1") {
		t.Error("result should contain 'Line 1'")
	}
	if !strings.Contains(result, "Line 2") {
		t.Error("result should contain 'Line 2'")
	}
	if !strings.Contains(result, "Line 3") {
		t.Error("result should contain 'Line 3'")
	}
}

func TestEngineRenderBanner(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	banner := ` _____     _
| ____|___| |__   ___  _ __
|  _| / _ \ '_ \ / _ \| '__|
| |__|  __/ |_) | (_) | |
|_____\___|_.__/ \___/|_|
`

	src := `<vbox>
		<text value="` + banner + `" />
	</vbox>`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}

	if !strings.Contains(result, " _____") {
		t.Error("result should contain banner characters")
	}
}

func TestEngineRenderVBox(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<vbox gap="1">
		<text value="Line 1" />
		<text value="Line 2" />
		<text value="Line 3" />
	</vbox>`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}
}

func TestEngineRenderHBox(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<hbox gap="2">
		<text value="A" />
		<text value="B" />
		<text value="C" />
	</hbox>`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}
}

func TestEngineRenderButton(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<button label="Click Me" />`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if !strings.Contains(result, "Click Me") {
		t.Errorf("result = %q, expected to contain 'Click Me'", result)
	}
}

func TestEngineRenderInput(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<input placeholder="Enter text..." />`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if !strings.Contains(result, "Enter text...") {
		t.Errorf("result = %q, expected to contain 'Enter text...'", result)
	}
}

func TestEngineRenderProgress(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<progress value="50" width="20" />`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}
}

func TestEngineRenderDivider(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<divider char="-" />`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}
}

func TestEngineRenderGrid(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<grid cols="2">
		<text value="A" />
		<text value="B" />
		<text value="C" />
		<text value="D" />
	</grid>`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}
}

func TestEngineComplexLayout(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	src := `<window title="Setup" width="60" height="20" center="true">
		<vbox padding="2" gap="1">
			<text value="Welcome to Freecode Setup" bold="true" />
			<divider />
			<vbox gap="0">
				<text value="Provider:" />
				<input placeholder="Select a provider..." />
			</vbox>
			<spacer flex="1" />
			<hbox gap="2">
				<spacer flex="1" />
				<button id="cancel" label="Cancel" />
				<button id="next" label="Next" primary="true" />
			</hbox>
		</vbox>
	</window>`

	result, err := engine.ParseAndRender(src, 80, 24, mockRenderer{})
	if err != nil {
		t.Fatalf("ParseAndRender error = %v", err)
	}

	if result == "" {
		t.Error("result is empty, expected rendered content")
	}

	if !strings.Contains(result, "Welcome to Freecode Setup") {
		t.Errorf("result missing title text")
	}
	if !strings.Contains(result, "Cancel") {
		t.Errorf("result missing Cancel button")
	}
	if !strings.Contains(result, "Next") {
		t.Errorf("result missing Next button")
	}
	if !strings.Contains(result, "Provider:") {
		t.Errorf("result missing Provider label")
	}
}

func TestEngineGetComponent(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	node := &ComponentNode{
		ID:     "mybutton",
		Type:   TypeButton,
		Width:  10,
		Height: 1,
	}
	engine.RegisterComponent("mybutton", node)

	btn := engine.GetComponent("mybutton")
	if btn == nil {
		t.Error("GetComponent returned nil, expected button")
	}
	if btn != nil && btn.ID != "mybutton" {
		t.Errorf("button ID = %q, want %q", btn.ID, "mybutton")
	}

	missing := engine.GetComponent("nonexistent")
	if missing != nil {
		t.Error("GetComponent should return nil for nonexistent component")
	}
}

func TestEngineRegisterComponent(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	node := &ComponentNode{
		ID:     "custom",
		Type:   TypeText,
		Width:  10,
		Height: 1,
	}

	engine.RegisterComponent("custom", node)

	got := engine.GetComponent("custom")
	if got == nil {
		t.Error("GetComponent returned nil, expected custom component")
	}
	if got != nil && got.ID != "custom" {
		t.Errorf("component ID = %q, want %q", got.ID, "custom")
	}
}

func TestGetInt(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	attrs := map[string]string{
		"width":  "100",
		"height": "50",
		"empty":  "",
	}

	if got := engine.getInt(attrs, "width", 0); got != 100 {
		t.Errorf("getInt width = %d, want 100", got)
	}
	if got := engine.getInt(attrs, "height", 0); got != 50 {
		t.Errorf("getInt height = %d, want 50", got)
	}
	if got := engine.getInt(attrs, "missing", 42); got != 42 {
		t.Errorf("getInt missing = %d, want 42", got)
	}
	if got := engine.getInt(attrs, "empty", 99); got != 99 {
		t.Errorf("getInt empty = %d, want 99", got)
	}
}

func TestGetBool(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	attrs := map[string]string{
		"true1":  "true",
		"true2":  "yes",
		"true3":  "1",
		"false1": "false",
		"false2": "no",
		"false3": "0",
		"empty":  "",
	}

	if !engine.getBool(attrs, "true1", false) {
		t.Error("getBool true1 should be true")
	}
	if !engine.getBool(attrs, "true2", false) {
		t.Error("getBool true2 should be true")
	}
	if !engine.getBool(attrs, "true3", false) {
		t.Error("getBool true3 should be true")
	}
	if engine.getBool(attrs, "false1", true) {
		t.Error("getBool false1 should be false")
	}
	if engine.getBool(attrs, "false2", true) {
		t.Error("getBool false2 should be false")
	}
	if engine.getBool(attrs, "false3", true) {
		t.Error("getBool false3 should be false")
	}
	if engine.getBool(attrs, "empty", true) {
		t.Error("getBool empty should return default")
	}
	if !engine.getBool(attrs, "missing", true) {
		t.Error("getBool missing should return default")
	}
}

func TestGetColor(t *testing.T) {
	engine := NewEngine[mockRenderer]()

	attrs := map[string]string{
		"fg":     "#ffffff",
		"bg":     "#000000",
		"border": "red",
		"empty":  "",
	}

	if got := engine.getColor(attrs, "fg", ""); got != "#ffffff" {
		t.Errorf("getColor fg = %q, want %q", got, "#ffffff")
	}
	if got := engine.getColor(attrs, "bg", ""); got != "#000000" {
		t.Errorf("getColor bg = %q, want %q", got, "#000000")
	}
	if got := engine.getColor(attrs, "missing", "default"); got != "default" {
		t.Errorf("getColor missing = %q, want %q", got, "default")
	}
}

func TestAllElementTypes(t *testing.T) {
	types := []struct {
		name     string
		elemType ElementType
	}{
		{"window", TypeWindow},
		{"vbox", TypeVBox},
		{"hbox", TypeHBox},
		{"grid", TypeGrid},
		{"text", TypeText},
		{"list", TypeList},
		{"button", TypeButton},
		{"input", TypeInput},
		{"spacer", TypeSpacer},
		{"divider", TypeDivider},
		{"progress", TypeProgress},
	}

	for _, tt := range types {
		t.Run(tt.name, func(t *testing.T) {
			src := "<" + tt.name + ">"
			tmpl, err := Parse(src)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", src, err)
			}
			if tmpl.Root.Type != tt.elemType {
				t.Errorf("Parse(%q) type = %v, want %v", src, tmpl.Root.Type, tt.elemType)
			}
		})
	}
}

var _ renderer.Renderer = mockRenderer{}
