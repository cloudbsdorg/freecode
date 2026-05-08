package lsp

import (
	"testing"
)

func TestServerCapabilities(t *testing.T) {
	caps := ServerCapabilities{
		TextDocumentSync:        1,
		HoverProvider:          true,
		DefinitionProvider:      true,
		ReferencesProvider:      false,
		ImplementationProvider:  false,
		DocumentSymbolProvider:  true,
		WorkspaceSymbolProvider: true,
	}

	if !caps.HoverProvider {
		t.Error("HoverProvider should be true")
	}
	if caps.ReferencesProvider {
		t.Error("ReferencesProvider should be false")
	}
}

func TestTextDocument(t *testing.T) {
	td := &textDocument{
		uri:     "file:///test.go",
		version: 1,
		text:    "package main",
	}

	if td.uri != "file:///test.go" {
		t.Errorf("uri = %q, want %q", td.uri, "file:///test.go")
	}
	if td.version != 1 {
		t.Errorf("version = %d, want %d", td.version, 1)
	}
	if td.text != "package main" {
		t.Errorf("text = %q, want %q", td.text, "package main")
	}
}

func TestDiagnosticStore(t *testing.T) {
	store := NewDiagnosticStore(nil)
	if store == nil {
		t.Fatal("NewDiagnosticStore() returned nil")
	}
}

func TestDiagnostic(t *testing.T) {
	d := Diagnostic{
		Range: Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 1, Character: 10}},
		Severity: 1,
		Message: "test error",
	}

	if d.Message != "test error" {
		t.Errorf("Message = %q, want %q", d.Message, "test error")
	}
}

func TestPosition(t *testing.T) {
	p := Position{
		Line:      10,
		Character: 5,
	}

	if p.Line != 10 {
		t.Errorf("Line = %d, want %d", p.Line, 10)
	}
	if p.Character != 5 {
		t.Errorf("Character = %d, want %d", p.Character, 5)
	}
}

func TestRange(t *testing.T) {
	r := Range{
		Start: Position{Line: 1, Character: 0},
		End:   Position{Line: 1, Character: 10},
	}

	if r.Start.Line != 1 {
		t.Errorf("Start.Line = %d, want %d", r.Start.Line, 1)
	}
	if r.End.Character != 10 {
		t.Errorf("End.Character = %d, want %d", r.End.Character, 10)
	}
}