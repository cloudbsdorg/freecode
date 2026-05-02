package shell

import (
	"testing"
)

func TestStartPTY(t *testing.T) {
	pty, err := StartPTY("bash", "-c", "echo hello")
	if err != nil {
		t.Fatalf("StartPTY() error = %v", err)
	}
	if pty == nil {
		t.Fatal("StartPTY() returned nil")
	}
}

func TestPTYWrite(t *testing.T) {
	pty, _ := StartPTY("bash")
	n, err := pty.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 0 {
		t.Errorf("Write() returned %d, want 0", n)
	}
}

func TestPTYRead(t *testing.T) {
	pty, _ := StartPTY("bash")
	buf := make([]byte, 1024)
	n, err := pty.Read(buf)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if n != 0 {
		t.Errorf("Read() returned %d, want 0", n)
	}
}

func TestPTYClose(t *testing.T) {
	pty, _ := StartPTY("bash")
	err := pty.Close()
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestPTYResize(t *testing.T) {
	pty, _ := StartPTY("bash")
	err := pty.Resize(24, 80)
	if err != nil {
		t.Fatalf("Resize() error = %v", err)
	}
}