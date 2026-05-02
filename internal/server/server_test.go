package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := New("localhost:18792")
	if s == nil {
		t.Fatal("New() returned nil")
	}
	if s.Addr != "localhost:18792" {
		t.Errorf("Addr = %q, want %q", s.Addr, "localhost:18792")
	}
}

func TestServerHandle(t *testing.T) {
	s := New("localhost:0")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	s.Handle("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestServerHandleFunc(t *testing.T) {
	s := New("localhost:0")
	s.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handled"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestServerMount(t *testing.T) {
	s := New("localhost:0")
	mux := http.NewServeMux()
	mux.HandleFunc("/sub", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("mounted"))
	})
	s.Mount("/mount/", mux)

	req := httptest.NewRequest("GET", "/mount/sub", nil)
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)
}

func TestServerURL(t *testing.T) {
	s := New("localhost:18792")
	url := s.URL()
	if url != "http://localhost:18792" {
		t.Errorf("URL() = %q, want %q", url, "http://localhost:18792")
	}
}

func TestServerStop(t *testing.T) {
	s := New("localhost:0")
	err := s.Stop()
	if err != nil {
		t.Errorf("Stop() error = %v", err)
	}
}

func TestHealthHandler(t *testing.T) {
	h := NewHealthHandler()
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("status = %q, want %q", resp["status"], "ok")
	}
}

func TestSetupRoutes(t *testing.T) {
	s := New("localhost:0")
	s.SetupRoutes()

	tests := []struct {
		path       string
		method     string
		wantStatus int
		wantBody   string
	}{
		{"/health", "GET", http.StatusOK, `"status":"ok"`},
		{"/api/v1/sessions", "GET", http.StatusOK, `"sessions":[]`},
		{"/api/v1/agents", "GET", http.StatusOK, `"agents":[]`},
		{"/api/v1/tools", "GET", http.StatusOK, `"tools":[]`},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			s.Router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Status = %d, want %d", rr.Code, tt.wantStatus)
			}
			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("Body = %q, want to contain %q", rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestSessionHandlerGet(t *testing.T) {
	t.Skip("requires non-nil Manager with CreateSession implemented")
}

func TestSessionHandlerPostInvalidJSON(t *testing.T) {
	h := NewSessionHandler(nil)
	req := httptest.NewRequest("POST", "/sessions", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestAgentHandlerGet(t *testing.T) {
	h := NewAgentHandler(nil, nil)
	req := httptest.NewRequest("GET", "/agents", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
}

func TestAgentHandlerPostInvalidJSON(t *testing.T) {
	h := NewAgentHandler(nil, nil)
	req := httptest.NewRequest("POST", "/agents", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestAgentHandlerPostWithEngine(t *testing.T) {
	t.Skip("requires non-nil Engine")
}