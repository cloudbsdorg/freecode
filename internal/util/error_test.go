package util

import (
	"errors"
	"fmt"
	"testing"
)

type customError struct {
	msg   string
	code  int
	cause error
}

func (e *customError) Error() string { return e.msg }
func (e *customError) Cause() error { return e.cause }

type causerError struct {
	msg   string
	cause error
}

func (e *causerError) Error() string { return e.msg }
func (e *causerError) Cause() error  { return e.cause }

type errorWithData struct {
	Message string
	Code    int
	Details string
}

func (e *errorWithData) Error() string { return e.Message }

func TestErrorFormat(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantEmpty bool
		wantContains string
	}{
		{
			name:     "nil error returns empty",
			err:      nil,
			wantEmpty: true,
		},
		{
			name:     "standard error",
			err:      errors.New("test error"),
			wantContains: "test error",
		},
		{
			name:     "formatted error with %v",
			err:      fmt.Errorf("wrapped: %w", errors.New("inner")),
			wantContains: "wrapped",
		},
		{
			name:     "custom error with cause",
			err:      &customError{msg: "outer", cause: errors.New("inner")},
			wantContains: "outer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorFormat(tt.err)
			if tt.wantEmpty {
				if got != "" {
					t.Errorf("ErrorFormat() = %q, want empty string", got)
				}
				return
			}
			if tt.wantContains != "" && !contains(got, tt.wantContains) {
				t.Errorf("ErrorFormat() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantEmpty bool
		want     string
		wantContains string
	}{
		{
			name:     "nil error returns empty",
			err:      nil,
			wantEmpty: true,
		},
		{
			name:     "standard error",
			err:      errors.New("test error message"),
			want:     "test error message",
		},
		{
			name:     "fmt.Errorf",
			err:      fmt.Errorf("wrapped: %w", errors.New("inner error")),
			wantContains: "inner error",
		},
		{
			name:     "custom error with Message method",
			err:      &errorWithData{Message: "custom message", Code: 42},
			want:     "custom message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorMessage(tt.err)
			if tt.wantEmpty {
				if got != "" {
					t.Errorf("ErrorMessage() = %q, want empty string", got)
				}
				return
			}
			if tt.want != "" && got != tt.want {
				t.Errorf("ErrorMessage() = %q, want %q", got, tt.want)
			}
			if tt.wantContains != "" && !contains(got, tt.wantContains) {
				t.Errorf("ErrorMessage() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestErrorData(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		got := ErrorData(nil)
		if got != nil {
			t.Errorf("ErrorData(nil) = %v, want nil", got)
		}
	})

	t.Run("standard error", func(t *testing.T) {
		err := errors.New("test error")
		got := ErrorData(err)
		if got == nil {
			t.Fatal("ErrorData() returned nil")
		}
		if got["type"] != "error" && got["type"] != "Error" && got["type"] != "errorString" {
			t.Errorf("ErrorData() type = %q, want 'error', 'Error', or 'errorString'", got["type"])
		}
		if got["message"] != "test error" {
			t.Errorf("ErrorData() message = %q, want %q", got["message"], "test error")
		}
		if got["formatted"] == "" {
			t.Error("ErrorData() formatted should not be empty")
		}
	})

	t.Run("wrapped error", func(t *testing.T) {
		inner := errors.New("inner error")
		err := fmt.Errorf("wrapper: %w", inner)
		got := ErrorData(err)
		if got == nil {
			t.Fatal("ErrorData() returned nil")
		}
		if got["message"] == "" {
			t.Error("ErrorData() message should not be empty for wrapped error")
		}
	})

	t.Run("custom error with cause", func(t *testing.T) {
		inner := errors.New("inner cause")
		err := &customError{msg: "outer message", cause: inner}
		got := ErrorData(err)
		if got == nil {
			t.Fatal("ErrorData() returned nil")
		}
		if got["type"] != "customError" {
			t.Errorf("ErrorData() type = %q, want %q", got["type"], "customError")
		}
		if got["message"] != "outer message" {
			t.Errorf("ErrorData() message = %q, want %q", got["message"], "outer message")
		}
	})

	t.Run("error with cause interface", func(t *testing.T) {
		inner := errors.New("inner cause")
		err := &causerError{msg: "outer", cause: inner}
		got := ErrorData(err)
		if got == nil {
			t.Fatal("ErrorData() returned nil")
		}
		if got["cause"] == nil {
			t.Error("ErrorData() cause should not be nil")
		}
	})

	t.Run("error with stack", func(t *testing.T) {
		err := errors.New("error with stack")
		got := ErrorData(err)
		if got == nil {
			t.Fatal("ErrorData() returned nil")
		}
		if got["stack"] != nil {
			t.Logf("ErrorData() has stack trace: %v", got["stack"])
		}
	})
}

func TestErrorData_ExtractsAllFields(t *testing.T) {
	err := &errorWithData{
		Message: "test message",
		Code:    100,
		Details: "some details",
	}
	got := ErrorData(err)
	if got == nil {
		t.Fatal("ErrorData() returned nil")
	}

	if got["Message"] != "test message" {
		t.Errorf("ErrorData() Message = %q, want %q", got["Message"], "test message")
	}
	if got["Code"] != 100 {
		t.Errorf("ErrorData() Code = %v, want %d", got["Code"], 100)
	}
	if got["Details"] != "some details" {
		t.Errorf("ErrorData() Details = %q, want %q", got["Details"], "some details")
	}
}

func TestErrorFormat_Wrapping(t *testing.T) {
	t.Run("unwrapping chain", func(t *testing.T) {
		err1 := errors.New("level 1")
		err2 := fmt.Errorf("level 2: %w", err1)
		err3 := fmt.Errorf("level 3: %w", err2)

		if !errors.Is(err3, err1) {
			t.Error("errors.Is(err3, err1) = false, want true")
		}
		if !errors.Is(err3, err2) {
			t.Error("errors.Is(err3, err2) = false, want true")
		}
	})

	t.Run("error as", func(t *testing.T) {
		err := &customError{msg: "test", code: 42}
		var ce *customError
		if !errors.As(err, &ce) {
			t.Error("errors.As(err, &ce) = false, want true")
		}
		if ce.code != 42 {
			t.Errorf("ce.code = %d, want %d", ce.code, 42)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type mockError struct {
	name    string
	message string
	stack   []uintptr
}

func (e *mockError) Error() string {
	return e.name + ": " + e.message
}

var _ error = (*mockError)(nil)

func TestErrorFormat_Empty(t *testing.T) {
	result := ErrorFormat(nil)
	if result != "" {
		t.Errorf("ErrorFormat(nil) = %q, want empty string", result)
	}
}

func TestErrorMessage_Empty(t *testing.T) {
	result := ErrorMessage(nil)
	if result != "" {
		t.Errorf("ErrorMessage(nil) = %q, want empty string", result)
	}
}

func TestErrorData_Empty(t *testing.T) {
	result := ErrorData(nil)
	if result != nil {
		t.Errorf("ErrorData(nil) = %v, want nil", result)
	}
}

func TestErrorMessage_WithCause(t *testing.T) {
	inner := errors.New("inner error")
	outer := fmt.Errorf("outer error: %w", inner)

	msg := ErrorMessage(outer)
	if msg == "" {
		t.Error("ErrorMessage() returned empty string for wrapped error")
	}
	if !contains(msg, "outer error") {
		t.Errorf("ErrorMessage() = %q, want to contain 'outer error'", msg)
	}
}

func TestErrorData_FormattedIncludesError(t *testing.T) {
	err := errors.New("test error")
	data := ErrorData(err)
	if data == nil {
		t.Fatal("ErrorData() returned nil")
	}
	formatted, ok := data["formatted"].(string)
	if !ok {
		t.Fatal("formatted is not a string")
	}
	if !contains(formatted, "test error") {
		t.Errorf("formatted = %q, want to contain 'test error'", formatted)
	}
}

type stringableError struct {
	val string
}

func (e *stringableError) Error() string {
	return e.val
}

func (e *stringableError) String() string {
	return "stringable: " + e.val
}

func TestErrorFormat_Stringable(t *testing.T) {
	err := &stringableError{val: "test value"}
	result := ErrorFormat(err)
	if !contains(result, "test value") {
		t.Errorf("ErrorFormat() = %q, want to contain 'test value'", result)
	}
}

func TestGetErrorType(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil", nil, "error"},
		{"standard", errors.New("test"), "errorString"},
		{"custom", &customError{msg: "test"}, "customError"},
		{"fmt", fmt.Errorf("test"), "errorString"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getErrorType(tt.err)
			if got != tt.want {
				t.Errorf("getErrorType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractCause(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if extractCause(nil) != nil {
			t.Error("extractCause(nil) should return nil")
		}
	})

	t.Run("custom error with cause", func(t *testing.T) {
		inner := errors.New("inner")
		outer := &customError{msg: "outer", cause: inner}
		got := extractCause(outer)
		if got != inner {
			t.Errorf("extractCause() = %v, want %v", got, inner)
		}
	})

	t.Run("causer error", func(t *testing.T) {
		inner := errors.New("inner")
		outer := &causerError{msg: "outer", cause: inner}
		got := extractCause(outer)
		if got != inner {
			t.Errorf("extractCause() = %v, want %v", got, inner)
		}
	})
}

func TestExtractCustomData(t *testing.T) {
	t.Run("struct fields", func(t *testing.T) {
		err := &errorWithData{Message: "msg", Code: 42, Details: "details"}
		data := make(map[string]interface{})
		extractCustomData(err, data)

		if data["Message"] != "msg" {
			t.Errorf("Message = %q, want %q", data["Message"], "msg")
		}
		if data["Code"] != 42 {
			t.Errorf("Code = %v, want %d", data["Code"], 42)
		}
		if data["Details"] != "details" {
			t.Errorf("Details = %q, want %q", data["Details"], "details")
		}
	})
}

func TestErrorData_AllTypes(t *testing.T) {
	err := &errorWithData{
		Message: "test error",
		Code:    42,
		Details: "details",
	}

	data := ErrorData(err)
	if data == nil {
		t.Fatal("ErrorData() returned nil")
	}

	if data["Message"] != "test error" {
		t.Errorf("Message = %v, want %q", data["Message"], "test error")
	}
	if data["Code"] != 42 {
		t.Errorf("Code = %v, want %d", data["Code"], 42)
	}
	if data["Details"] != "details" {
		t.Errorf("Details = %v, want %q", data["Details"], "details")
	}
}

func TestGetStackTrace(t *testing.T) {
	err := errors.New("test error")
	stack := getStackTrace(err)
	t.Logf("Stack trace: %s", stack)
}

type nestedError struct {
	msg   string
	cause error
}

func (e *nestedError) Error() string { return e.msg }
func (e *nestedError) Cause() error  { return e.cause }

func TestErrorData_NestedCauses(t *testing.T) {
	inner := &nestedError{msg: "inner most"}
	middle := &nestedError{msg: "middle", cause: inner}
	outer := &nestedError{msg: "outer", cause: middle}

	data := ErrorData(outer)
	if data == nil {
		t.Fatal("ErrorData() returned nil")
	}

	if data["cause"] == nil {
		t.Error("Expected cause to be extracted")
	}
}

func BenchmarkErrorFormat(b *testing.B) {
	err := errors.New("benchmark error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ErrorFormat(err)
	}
}

func BenchmarkErrorMessage(b *testing.B) {
	err := errors.New("benchmark error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ErrorMessage(err)
	}
}

func BenchmarkErrorData(b *testing.B) {
	err := errors.New("benchmark error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ErrorData(err)
	}
}
