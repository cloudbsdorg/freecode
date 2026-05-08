package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient("https://api.example.com")
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.baseURL != "https://api.example.com" {
		t.Fatalf("Expected baseURL 'https://api.example.com', got '%s'", c.baseURL)
	}
	if c.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
	if c.httpClient.Timeout != 30*time.Second {
		t.Fatalf("Expected default timeout 30s, got %v", c.httpClient.Timeout)
	}
}

func TestWithAPIKey(t *testing.T) {
	c := NewClient("https://api.example.com", WithAPIKey("test-key"))
	if c.apiKey != "test-key" {
		t.Fatalf("Expected apiKey 'test-key', got '%s'", c.apiKey)
	}
}

func TestWithHeader(t *testing.T) {
	c := NewClient("https://api.example.com", WithHeader("X-Custom", "value"))
	if c.headers["X-Custom"] != "value" {
		t.Fatalf("Expected header 'X-Custom: value', got '%s'", c.headers["X-Custom"])
	}
}

func TestWithTimeout(t *testing.T) {
	c := NewClient("https://api.example.com", WithTimeout(60*time.Second))
	if c.httpClient.Timeout != 60*time.Second {
		t.Fatalf("Expected timeout 60s, got %v", c.httpClient.Timeout)
	}
}

func TestSetAPIKey(t *testing.T) {
	c := NewClient("https://api.example.com")
	c.SetAPIKey("new-key")
	if c.apiKey != "new-key" {
		t.Fatalf("Expected apiKey 'new-key', got '%s'", c.apiKey)
	}
}

func TestSetHeader(t *testing.T) {
	c := NewClient("https://api.example.com")
	c.SetHeader("X-Test", "test-value")
	if c.headers["X-Test"] != "test-value" {
		t.Fatalf("Expected header 'X-Test: test-value', got '%s'", c.headers["X-Test"])
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	resp, err := c.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(resp) != `{"result":"ok"}` {
		t.Fatalf("Unexpected response: %s", string(resp))
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["key"] != "value" {
			t.Errorf("Unexpected body: %v", body)
		}
		w.Write([]byte(`{"posted":true}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	resp, err := c.Post(context.Background(), "/test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if string(resp) != `{"posted":true}` {
		t.Fatalf("Unexpected response: %s", string(resp))
	}
}

func TestClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"updated":true}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	resp, err := c.Put(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	if string(resp) != `{"updated":true}` {
		t.Fatalf("Unexpected response: %s", string(resp))
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	resp, err := c.Delete(context.Background(), "/test")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if string(resp) != `{"deleted":true}` {
		t.Fatalf("Unexpected response: %s", string(resp))
	}
}

func TestClient_GetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"name":"test","value":123}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	err := c.GetJSON(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}
	if result.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", result.Name)
	}
	if result.Value != 123 {
		t.Errorf("Expected value 123, got %d", result.Value)
	}
}

func TestClient_PostJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	var result struct {
		Success bool `json:"success"`
	}
	err := c.PostJSON(context.Background(), "/test", map[string]string{"key": "value"}, &result)
	if err != nil {
		t.Fatalf("PostJSON failed: %v", err)
	}
	if !result.Success {
		t.Error("Expected success to be true")
	}
}

func TestClient_GetPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("Expected page=2, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("per_page") != "10" {
			t.Errorf("Expected per_page=10, got %s", r.URL.Query().Get("per_page"))
		}
		w.Write([]byte(`{"data":[1,2,3]}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	var result struct {
		Data []int `json:"data"`
	}
	err := c.GetPage(context.Background(), "/test", 2, 10, &result)
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}
}

func TestClient_Do_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.Get(context.Background(), "/test")
	if err == nil {
		t.Fatal("Expected error for 404 response")
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "not found",
	}
	expected := "api error: status 404, message: not found"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestParseAPIError(t *testing.T) {
	body := []byte(`{"error":"bad request","message":"invalid input","details":{"field":"value"}}`)
	err := ParseAPIError(400, body)
	if err.StatusCode != 400 {
		t.Errorf("Expected status 400, got %d", err.StatusCode)
	}
	if err.Message != "invalid input" {
		t.Errorf("Expected message 'invalid input', got '%s'", err.Message)
	}
	if err.Details["field"] != "value" {
		t.Errorf("Expected details field 'value', got '%v'", err.Details["field"])
	}
}

func TestParseAPIError_Fallback(t *testing.T) {
	body := []byte(`not json at all`)
	err := ParseAPIError(500, body)
	if err.StatusCode != 500 {
		t.Errorf("Expected status 500, got %d", err.StatusCode)
	}
	if err.Message != "not json at all" {
		t.Errorf("Expected message 'not json at all', got '%s'", err.Message)
	}
}

func TestPageResponse(t *testing.T) {
	pr := PageResponse[string]{
		Data:       []string{"a", "b", "c"},
		Page:       1,
		PerPage:    10,
		Total:      25,
		TotalPages: 3,
	}
	if len(pr.Data) != 3 {
		t.Errorf("Expected 3 items, got %d", len(pr.Data))
	}
	if pr.Total != 25 {
		t.Errorf("Expected total 25, got %d", pr.Total)
	}
	if pr.TotalPages != 3 {
		t.Errorf("Expected 3 pages, got %d", pr.TotalPages)
	}
}

func TestClient_AuthorizationHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected 'Bearer test-token', got '%s'", auth)
		}
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, WithAPIKey("test-token"))
	_, _ = c.Get(context.Background(), "/test")
}

func TestClient_CustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		accept := r.Header.Get("Accept")
		custom := r.Header.Get("X-Custom")
		if contentType != "application/json" {
			t.Errorf("Expected 'application/json', got '%s'", contentType)
		}
		if accept != "application/json" {
			t.Errorf("Expected 'application/json', got '%s'", accept)
		}
		if custom != "custom-value" {
			t.Errorf("Expected 'custom-value', got '%s'", custom)
		}
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, WithHeader("X-Custom", "custom-value"))
	_, _ = c.Get(context.Background(), "/test")
}