package env

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	e := New()
	if e == nil {
		t.Fatal("New() returned nil")
	}
}

func TestEnvGetSet(t *testing.T) {
	e := New()

	e.Set("TEST_KEY", "test_value")
	if got := e.Get("TEST_KEY"); got != "test_value" {
		t.Errorf("Get(TEST_KEY) = %q, want %q", got, "test_value")
	}
}

func TestEnvGetUnset(t *testing.T) {
	e := New()

	e.Set("TEST_KEY", "test_value")
	e.Unset("TEST_KEY")
	if got := e.Get("TEST_KEY"); got != "" {
		t.Errorf("Get(TEST_KEY) after Unset = %q, want empty", got)
	}
}

func TestEnvGetFallbackToOS(t *testing.T) {
	e := New()

	os.Setenv("TEST_OS_KEY", "os_value")
	defer os.Unsetenv("TEST_OS_KEY")

	if got := e.Get("TEST_OS_KEY"); got != "os_value" {
		t.Errorf("Get(TEST_OS_KEY) = %q, want %q (from OS)", got, "os_value")
	}
}

func TestEnvOverrideOS(t *testing.T) {
	e := New()

	os.Setenv("TEST_KEY", "os_value")
	defer os.Unsetenv("TEST_KEY")

	e.Set("TEST_KEY", "env_value")
	if got := e.Get("TEST_KEY"); got != "env_value" {
		t.Errorf("Get(TEST_KEY) = %q, want %q (env overrides OS)", got, "env_value")
	}
}

func TestEnvAll(t *testing.T) {
	e := New()

	e.Set("KEY1", "value1")
	e.Set("KEY2", "value2")

	all := e.All()
	if len(all) != 2 {
		t.Errorf("All() returned %d keys, want 2", len(all))
	}
	if all["KEY1"] != "value1" {
		t.Errorf("All()[KEY1] = %q, want %q", all["KEY1"], "value1")
	}
	if all["KEY2"] != "value2" {
		t.Errorf("All()[KEY2] = %q, want %q", all["KEY2"], "value2")
	}
}

func TestEnvExpand(t *testing.T) {
	e := New()

	e.Set("FOO", "bar")
	e.Set("BAZ", "qux")

	result := e.Expand("$FOO and $BAZ")
	if result == "" {
		t.Error("Expand() returned empty string")
	}
}

func TestEnvUnsetNonExistent(t *testing.T) {
	e := New()
	e.Unset("NONEXISTENT")
}

func TestEnvConcurrentAccess(t *testing.T) {
	e := New()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			e.Set("KEY", "value")
			e.Get("KEY")
			e.All()
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}