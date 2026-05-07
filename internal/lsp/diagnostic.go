package lsp

import (
	"fmt"
	"sync"
	"time"
)

const (
	DiagnosticSeverityError   = 1
	DiagnosticSeverityWarning = 2
	DiagnosticSeverityInfo    = 3
	DiagnosticSeverityHint    = 4
)

const (
	DiagnosticsDebounceMs = 150
)

type DiagnosticStore struct {
	mu        sync.RWMutex
	push      map[string][]Diagnostic
	pull      map[string][]Diagnostic
	published map[string]*PublishInfo
	pending   map[string]*pendingDiagnostic
	onPublish func(uri string, diagnostics []Diagnostic)
}

type PublishInfo struct {
	At      time.Time
	Version int
}

type pendingDiagnostic struct {
	timer *time.Timer
	uri   string
}

func NewDiagnosticStore(onPublish func(uri string, diagnostics []Diagnostic)) *DiagnosticStore {
	return &DiagnosticStore{
		push:      make(map[string][]Diagnostic),
		pull:      make(map[string][]Diagnostic),
		published: make(map[string]*PublishInfo),
		pending:   make(map[string]*pendingDiagnostic),
		onPublish: onPublish,
	}
}

func (ds *DiagnosticStore) UpdatePushDiagnostics(uri string, diagnostics []Diagnostic) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.push[uri] = diagnostics
	ds.published[uri] = &PublishInfo{
		At:      time.Now(),
		Version: diagnosticsVersion(diagnostics),
	}

	ds.schedulePublish(uri, diagnostics)
}

func (ds *DiagnosticStore) schedulePublish(uri string, diagnostics []Diagnostic) {
	if existing, ok := ds.pending[uri]; ok {
		existing.timer.Stop()
	}

	ds.pending[uri] = &pendingDiagnostic{
		timer: time.AfterFunc(time.Millisecond*DiagnosticsDebounceMs, func() {
			ds.doPublish(uri, diagnostics)
		}),
	}
}

func (ds *DiagnosticStore) doPublish(uri string, diagnostics []Diagnostic) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	delete(ds.pending, uri)

	if ds.onPublish != nil {
		ds.onPublish(uri, diagnostics)
	}
}

func (ds *DiagnosticStore) UpdatePullDiagnostics(uri string, diagnostics []Diagnostic) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.pull[uri] = diagnostics
}

func (ds *DiagnosticStore) GetDiagnostics(uri string) []Diagnostic {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	pushDiags := ds.push[uri]
	pullDiags := ds.pull[uri]

	if len(pushDiags) == 0 && len(pullDiags) == 0 {
		return nil
	}

	if len(pushDiags) == 0 {
		return pullDiags
	}
	if len(pullDiags) == 0 {
		return pushDiags
	}

	return dedupeDiagnostics(append(pushDiags, pullDiags...))
}

func (ds *DiagnosticStore) GetAllDiagnostics() map[string][]Diagnostic {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	result := make(map[string][]Diagnostic)
	for uri, diags := range ds.push {
		result[uri] = diags
	}
	for uri, diags := range ds.pull {
		if _, exists := result[uri]; exists {
			result[uri] = dedupeDiagnostics(append(result[uri], diags...))
		} else {
			result[uri] = diags
		}
	}
	return result
}

func (ds *DiagnosticStore) Clear(uri string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	delete(ds.push, uri)
	delete(ds.pull, uri)
	delete(ds.published, uri)
	if pending, ok := ds.pending[uri]; ok {
		pending.timer.Stop()
		delete(ds.pending, uri)
	}
}

func diagnosticsVersion(diags []Diagnostic) int {
	h := 0
	for _, d := range diags {
		h += int(d.Range.Start.Line) + int(d.Range.Start.Character) + int(d.Range.End.Line) + int(d.Range.End.Character)
		h += len(d.Message)
	}
	return h
}

func dedupeDiagnostics(items []Diagnostic) []Diagnostic {
	seen := make(map[string]bool)
	result := make([]Diagnostic, 0, len(items))

	for _, item := range items {
		key := fmt.Sprintf("%v-%d-%s-%v", item.Range, item.Severity, item.Message, item.Source)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return result
}

func PrettyDiagnostic(diagnostic Diagnostic) string {
	severityMap := map[int]string{
		DiagnosticSeverityError:   "ERROR",
		DiagnosticSeverityWarning: "WARN",
		DiagnosticSeverityInfo:    "INFO",
		DiagnosticSeverityHint:    "HINT",
	}

	severity := severityMap[diagnostic.Severity]
	if severity == "" {
		severity = "ERROR"
	}

	line := diagnostic.Range.Start.Line + 1
	col := diagnostic.Range.Start.Character + 1

	return fmt.Sprintf("%s [%d:%d] %s", severity, line, col, diagnostic.Message)
}
