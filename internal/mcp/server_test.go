package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient()

	if c == nil {
		t.Fatal("NewClient() returned nil")
	}

	if c.servers == nil {
		t.Error("Client.servers is nil")
	}
}

func TestClientAddServer(t *testing.T) {
	c := NewClient()

	c.AddServer("test-id", "Test Server", "http://localhost:8080", "test-auth")

	if len(c.servers) != 1 {
		t.Errorf("len(Client.servers) = %d, want 1", len(c.servers))
	}

	srv, ok := c.servers["test-id"]
	if !ok {
		t.Error("Server 'test-id' not found")
	}

	if srv.Name != "Test Server" {
		t.Errorf("Server.Name = %q, want %q", srv.Name, "Test Server")
	}

	if srv.URL != "http://localhost:8080" {
		t.Errorf("Server.URL = %q, want %q", srv.URL, "http://localhost:8080")
	}
}

func TestClientRemoveServer(t *testing.T) {
	c := NewClient()

	c.AddServer("test-id", "Test Server", "http://localhost:8080", "")
	c.RemoveServer("test-id")

	if len(c.servers) != 0 {
		t.Errorf("len(Client.servers) = %d, want 0", len(c.servers))
	}
}

func TestClientListTools(t *testing.T) {
	c := NewClient()

	c.AddServer("test-id", "Test Server", "http://localhost:8080", "")

	tools, err := c.ListTools()
	if err != nil {
		t.Errorf("ListTools() error = %v", err)
	}

	if len(tools) != 0 {
		t.Errorf("len(tools) = %d, want 0", len(tools))
	}
}

func TestClientCallToolServerNotFound(t *testing.T) {
	c := NewClient()

	_, err := c.CallTool(context.Background(), "nonexistent", "test-tool", nil)
	if err == nil {
		t.Error("CallTool() should error for nonexistent server")
	}
}

func TestNewServer(t *testing.T) {
	s := NewServer(8080)

	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	if s.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", s.Port, 8080)
	}

	if s.Tools == nil {
		t.Error("Server.Tools is nil")
	}

	if s.Handlers == nil {
		t.Error("Server.Handlers is nil")
	}
}

func TestServerRegisterTool(t *testing.T) {
	s := NewServer(8080)

	s.RegisterTool("test-tool", "A test tool", func(args map[string]interface{}) (interface{}, error) {
		return "result", nil
	})

	if len(s.Tools) != 1 {
		t.Errorf("len(Server.Tools) = %d, want 1", len(s.Tools))
	}

	if s.Tools[0].Name != "test-tool" {
		t.Errorf("Tool.Name = %q, want %q", s.Tools[0].Name, "test-tool")
	}

	if _, ok := s.Handlers["test-tool"]; !ok {
		t.Error("Handler 'test-tool' not found")
	}
}

func TestServerHandleInitialize(t *testing.T) {
	s := NewServer(8080)

	req := ServerRequest{
		JSONRPC: "2.0",
		Method:  "initialize",
		Params:  nil,
		ID:      1,
	}

	resp := s.processRequest(&req)

	if resp.JSONRPC != "2.0" {
		t.Errorf("resp.JSONRPC = %q, want %q", resp.JSONRPC, "2.0")
	}

	if resp.ID != 1 {
		t.Errorf("resp.ID = %v, want %v", resp.ID, 1)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("resp.Result is not a map")
	}

	if result["protocolVersion"] != "2024-11-05" {
		t.Errorf("protocolVersion = %v, want %v", result["protocolVersion"], "2024-11-05")
	}
}

func TestServerHandleToolsList(t *testing.T) {
	s := NewServer(8080)
	s.RegisterTool("tool1", "Tool 1", nil)
	s.RegisterTool("tool2", "Tool 2", nil)

	req := ServerRequest{
		JSONRPC: "2.0",
		Method:  "tools/list",
		Params:  nil,
		ID:      1,
	}

	resp := s.processRequest(&req)

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("resp.Result is not a map")
	}

	tools, ok := result["tools"].([]Tool)
	if !ok {
		t.Fatal("result[tools] is not a []Tool")
	}

	if len(tools) != 2 {
		t.Errorf("len(tools) = %d, want 2", len(tools))
	}
}

func TestServerHandleToolsCall(t *testing.T) {
	s := NewServer(8080)
	s.RegisterTool("echo", "Echo tool", func(args map[string]interface{}) (interface{}, error) {
		return args, nil
	})

	params := `{"name":"echo","arguments":{"msg":"hello"}}`
	req := ServerRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  json.RawMessage(params),
		ID:      1,
	}

	resp := s.processRequest(&req)

	if resp.Error != nil {
		t.Errorf("resp.Error = %v, want nil", resp.Error)
	}
}

func TestServerHandleUnknownMethod(t *testing.T) {
	s := NewServer(8080)

	req := ServerRequest{
		JSONRPC: "2.0",
		Method:  "unknown",
		ID:      1,
	}

	resp := s.processRequest(&req)

	if resp.Error == nil {
		t.Error("resp.Error = nil, want error for unknown method")
	}

	if resp.Error.Code != -32601 {
		t.Errorf("resp.Error.Code = %d, want %d", resp.Error.Code, -32601)
	}
}

func TestServerHTTPHandler(t *testing.T) {
	s := NewServer(8080)
	s.RegisterTool("test", "Test", func(args map[string]interface{}) (interface{}, error) {
		return "ok", nil
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Handle(w, r)
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"jsonrpc":"2.0","method":"tools/list","id":1}`))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusOK)
	}
}
