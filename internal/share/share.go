package share

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Share struct {
	ID      string
	Content string
	URL     string
}

type Publisher interface {
	Publish(ctx context.Context, content string) (*Share, error)
	Get(ctx context.Context, id string) (*Share, error)
}

type pastebinPublisher struct {
	apiURL  string
	apiKey  string
	expiry  string
}

func NewPastebinPublisher(apiURL, apiKey string) *pastebinPublisher {
	return &pastebinPublisher{
		apiURL: apiURL,
		apiKey: apiKey,
		expiry: "1week",
	}
}

type localPublisher struct {
	shares map[string]*Share
}

func NewLocalPublisher() *localPublisher {
	return &localPublisher{
		shares: make(map[string]*Share),
	}
}

func (p *localPublisher) Publish(ctx context.Context, content string) (*Share, error) {
	id := uuid.New().String()[:8]
	share := &Share{
		ID:      id,
		Content: content,
		URL:     fmt.Sprintf("local://%s", id),
	}
	p.shares[id] = share
	return share, nil
}

func (p *localPublisher) Get(ctx context.Context, id string) (*Share, error) {
	if share, ok := p.shares[id]; ok {
		return share, nil
	}
	return nil, fmt.Errorf("share not found: %s", id)
}

type httpPublisher struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPPublisher(baseURL string) *httpPublisher {
	return &httpPublisher{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *httpPublisher) Publish(ctx context.Context, content string) (*Share, error) {
	data := map[string]string{
		"content": content,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/api/share", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to publish: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("publish failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Share{
		ID:      result.ID,
		Content: content,
		URL:     result.URL,
	}, nil
}

func (p *httpPublisher) Get(ctx context.Context, id string) (*Share, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/api/share/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get share: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("share not found: %s", id)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result Share
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

type multiPublisher struct {
	publishers []Publisher
}

func NewMultiPublisher(publishers ...Publisher) *multiPublisher {
	return &multiPublisher{publishers: publishers}
}

func (p *multiPublisher) Publish(ctx context.Context, content string) (*Share, error) {
	var lastErr error
	for _, pub := range p.publishers {
		share, err := pub.Publish(ctx, content)
		if err == nil {
			return share, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all publishers failed, last error: %w", lastErr)
}

func (p *multiPublisher) Get(ctx context.Context, id string) (*Share, error) {
	var lastErr error
	for _, pub := range p.publishers {
		share, err := pub.Get(ctx, id)
		if err == nil {
			return share, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all publishers failed, last error: %w", lastErr)
}

type PublisherFunc struct {
	publishFn func(ctx context.Context, content string) (*Share, error)
	getFn     func(ctx context.Context, id string) (*Share, error)
}

func (f *PublisherFunc) Publish(ctx context.Context, content string) (*Share, error) {
	return f.publishFn(ctx, content)
}

func (f *PublisherFunc) Get(ctx context.Context, id string) (*Share, error) {
	return f.getFn(ctx, id)
}

func NewAnonymousPublisher() *PublisherFunc {
	return &PublisherFunc{
		publishFn: func(ctx context.Context, content string) (*Share, error) {
			id := uuid.New().String()[:8]
			return &Share{
				ID:      id,
				Content: content,
				URL:     fmt.Sprintf("anon://%s", id),
			}, nil
		},
		getFn: func(ctx context.Context, id string) (*Share, error) {
			return nil, fmt.Errorf("anonymous shares cannot be retrieved: %s", id)
		},
	}
}