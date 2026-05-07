package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServerStartStop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	s := New("localhost:0")
	s.SetupRoutes()

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	addr := listener.Addr().String()
	s.Server.Addr = addr
	s.Addr = addr

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.Server.Serve(listener)
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	default:
	}

	resp, err := http.Get("http://" + addr + "/health")
	if err != nil {
		t.Fatalf("Health check request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Health status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	cancel()
	s.Server.Close()

	select {
	case <-errChan:
	case <-time.After(2 * time.Second):
		t.Error("Server did not stop after context cancellation")
	}
}

func TestServerStartAndServeHTTP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	s := New("localhost:0")
	s.SetupRoutes()

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	addr := listener.Addr().String()

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.Server.Serve(listener)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://" + addr + "/api/v1/sessions")
	if err != nil {
		t.Fatalf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	cancel()
	s.Server.Close()

	select {
	case <-errChan:
	case <-time.After(2 * time.Second):
		t.Error("Server did not stop")
	}
}
