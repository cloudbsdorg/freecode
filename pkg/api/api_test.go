package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler()
	if h == nil {
		t.Fatal("NewHandler() returned nil")
	}
	if h.tools == nil {
		t.Error("tools map not initialized")
	}
}

func TestHandlerRegisterTool(t *testing.T) {
	h := NewHandler()
	h.RegisterTool("test", "Test tool", func(args map[string]interface{}) (interface{}, error) {
		return "result", nil
	})

	if len(h.tools) != 1 {
		t.Errorf("tools count = %d, want %d", len(h.tools), 1)
	}
	if _, ok := h.tools["test"]; !ok {
		t.Error("tool 'test' not registered")
	}
}

func TestHandlerServeHTTPGetMethod(t *testing.T) {
	h := NewHandler()
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandlerServeHTTPInvalidJSON(t *testing.T) {
	h := NewHandler()
	req := httptest.NewRequest("POST", "/", strings.NewReader("invalid json"))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandlerInitialize(t *testing.T) {
	h := NewHandler()
	body := `{"method":"initialize","id":1}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, "2.0")
	}
	if fmt.Sprintf("%v", resp.ID) != "1" {
		t.Errorf("ID = %v, want 1", resp.ID)
	}
}

func TestHandlerToolsList(t *testing.T) {
	h := NewHandler()
	h.RegisterTool("tool1", "Tool 1", nil)
	h.RegisterTool("tool2", "Tool 2", nil)

	body := `{"method":"tools/list","id":2}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	result := resp.Result.(map[string]interface{})
	tools := result["tools"].([]interface{})
	if len(tools) != 2 {
		t.Errorf("tools count = %d, want %d", len(tools), 2)
	}
}

func TestHandlerUnknownMethod(t *testing.T) {
	h := NewHandler()
	body := `{"method":"unknown","id":3}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.Error == nil {
		t.Fatal("Error should not be nil for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("Error.Code = %d, want %d", resp.Error.Code, -32601)
	}
}

func TestResponseJSON(t *testing.T) {
	resp := Response{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"test": "value",
		},
		ID: 1,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", decoded.JSONRPC, "2.0")
	}
}

func TestErrorJSON(t *testing.T) {
	apiErr := Error{
		Code:    -32600,
		Message: "Invalid request",
	}

	data, err := json.Marshal(apiErr)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded Error
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.Code != -32600 {
		t.Errorf("Code = %d, want %d", decoded.Code, -32600)
	}
}

func TestToolJSON(t *testing.T) {
	tool := Tool{
		Name:        "test",
		Description: "Test tool",
		InputSchema: map[string]interface{}{"type": "object"},
	}

	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded Tool
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.Name != "test" {
		t.Errorf("Name = %q, want %q", decoded.Name, "test")
	}
}
