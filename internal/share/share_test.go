package share

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocalPublisher_Publish(t *testing.T) {
	publisher := NewLocalPublisher()
	ctx := context.Background()

	share, err := publisher.Publish(ctx, "test content")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if share.ID == "" {
		t.Error("Expected non-empty ID")
	}
	if share.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", share.Content)
	}
	if share.URL == "" {
		t.Error("Expected non-empty URL")
	}
}

func TestLocalPublisher_Get(t *testing.T) {
	publisher := NewLocalPublisher()
	ctx := context.Background()

	share, err := publisher.Publish(ctx, "test content")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	retrieved, err := publisher.Get(ctx, share.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Content != "test content" {
		t.Errorf("Expected 'test content', got '%s'", retrieved.Content)
	}
}

func TestLocalPublisher_Get_NotFound(t *testing.T) {
	publisher := NewLocalPublisher()
	ctx := context.Background()

	_, err := publisher.Get(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent share")
	}
}

func TestHTTPPublisher_Publish_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["content"] != "test content" {
			t.Errorf("Expected content 'test content', got '%s'", body["content"])
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "abc123", "url": "https://example.com/abc123"})
	}))
	defer server.Close()

	publisher := NewHTTPPublisher(server.URL)
	share, err := publisher.Publish(context.Background(), "test content")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if share.ID != "abc123" {
		t.Errorf("Expected ID 'abc123', got '%s'", share.ID)
	}
	if share.URL != "https://example.com/abc123" {
		t.Errorf("Expected URL 'https://example.com/abc123', got '%s'", share.URL)
	}
}

func TestHTTPPublisher_Publish_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	publisher := NewHTTPPublisher(server.URL)
	_, err := publisher.Publish(context.Background(), "test content")
	if err == nil {
		t.Error("Expected error for failed publish")
	}
}

func TestHTTPPublisher_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(Share{ID: "abc123", Content: "test content"})
	}))
	defer server.Close()

	publisher := NewHTTPPublisher(server.URL)
	share, err := publisher.Get(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if share.ID != "abc123" {
		t.Errorf("Expected ID 'abc123', got '%s'", share.ID)
	}
	if share.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", share.Content)
	}
}

func TestHTTPPublisher_Get_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	publisher := NewHTTPPublisher(server.URL)
	_, err := publisher.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for not found")
	}
}

func TestMultiPublisher_Publish_Success(t *testing.T) {
	localPub := &mockPublisher{err: nil, share: &Share{ID: "first-share", Content: "content", URL: "http://first.com"}}
	httpPub := &mockPublisher{err: nil, share: &Share{ID: "http-share", Content: "content", URL: "http://example.com"}}

	publisher := NewMultiPublisher(localPub, httpPub)
	share, err := publisher.Publish(context.Background(), "test")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if share.ID != "first-share" {
		t.Errorf("Expected ID 'first-share', got '%s'", share.ID)
	}
}

func TestMultiPublisher_Publish_FallbackToSecond(t *testing.T) {
	localPub := &mockPublisher{err: context.DeadlineExceeded, share: nil}
	httpPub := &mockPublisher{err: nil, share: &Share{ID: "fallback-share", Content: "content", URL: "http://example.com"}}

	publisher := NewMultiPublisher(localPub, httpPub)
	share, err := publisher.Publish(context.Background(), "test")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if share.ID != "fallback-share" {
		t.Errorf("Expected ID 'fallback-share', got '%s'", share.ID)
	}
}

func TestMultiPublisher_Publish_AllFail(t *testing.T) {
	failPub := &mockPublisher{err: context.DeadlineExceeded, share: nil}

	publisher := NewMultiPublisher(failPub)
	_, err := publisher.Publish(context.Background(), "test")
	if err == nil {
		t.Error("Expected error when all publishers fail")
	}
}

func TestMultiPublisher_Get_Success(t *testing.T) {
	localPub := &mockPublisher{err: nil, share: &Share{ID: "local-share", Content: "local content", URL: "local://test"}}
	httpPub := &mockPublisher{err: nil, share: &Share{ID: "http-share", Content: "http content", URL: "http://example.com"}}

	publisher := NewMultiPublisher(localPub, httpPub)
	share, err := publisher.Get(context.Background(), "local-share")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if share.Content != "local content" {
		t.Errorf("Expected 'local content', got '%s'", share.Content)
	}
}

func TestMultiPublisher_Get_FallbackToSecond(t *testing.T) {
	localPub := &mockPublisher{err: context.DeadlineExceeded, share: nil}
	httpPub := &mockPublisher{err: nil, share: &Share{ID: "http-share", Content: "http content", URL: "http://example.com"}}

	publisher := NewMultiPublisher(localPub, httpPub)
	share, err := publisher.Get(context.Background(), "http-share")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if share.ID != "http-share" {
		t.Errorf("Expected ID 'http-share', got '%s'", share.ID)
	}
}

func TestAnonymousPublisher_Publish(t *testing.T) {
	publisher := NewAnonymousPublisher()
	ctx := context.Background()

	share, err := publisher.Publish(ctx, "anonymous content")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if share.ID == "" {
		t.Error("Expected non-empty ID")
	}
	if share.Content != "anonymous content" {
		t.Errorf("Expected 'anonymous content', got '%s'", share.Content)
	}
}

func TestAnonymousPublisher_Get_Error(t *testing.T) {
	publisher := NewAnonymousPublisher()
	ctx := context.Background()

	_, err := publisher.Get(ctx, "some-id")
	if err == nil {
		t.Error("Expected error when getting from anonymous publisher")
	}
}

func TestShare_Struct(t *testing.T) {
	share := Share{
		ID:      "test-id",
		Content: "test content",
		URL:     "https://example.com/test",
	}

	if share.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", share.ID)
	}
	if share.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", share.Content)
	}
	if share.URL != "https://example.com/test" {
		t.Errorf("Expected URL 'https://example.com/test', got '%s'", share.URL)
	}
}

func TestPublisherFunc(t *testing.T) {
	var publishCalled, getCalled bool

	publisher := &PublisherFunc{
		publishFn: func(ctx context.Context, content string) (*Share, error) {
			publishCalled = true
			return &Share{ID: "func-id", Content: content, URL: "func://url"}, nil
		},
		getFn: func(ctx context.Context, id string) (*Share, error) {
			getCalled = true
			return &Share{ID: id, Content: "content", URL: "func://" + id}, nil
		},
	}

	share, err := publisher.Publish(context.Background(), "test")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
	if share.ID != "func-id" {
		t.Errorf("Expected 'func-id', got '%s'", share.ID)
	}
	if !publishCalled {
		t.Error("publishFn was not called")
	}

	_, err = publisher.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !getCalled {
		t.Error("getFn was not called")
	}
}

type mockPublisher struct {
	share *Share
	err   error
}

func (m *mockPublisher) Publish(ctx context.Context, content string) (*Share, error) {
	return m.share, m.err
}

func (m *mockPublisher) Get(ctx context.Context, id string) (*Share, error) {
	return m.share, m.err
}