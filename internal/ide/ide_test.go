package ide

import (
	"testing"
)

func TestStubIDE(t *testing.T) {
	ide := NewStubIDE()
	if ide == nil {
		t.Error("expected IDE")
	}
}
