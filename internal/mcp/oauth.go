package mcp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURI  string
	Scopes       []string
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type OAuthState struct {
	ServerID    string
	Provider    string
	ReturnURL   string
	CreatedAt   time.Time
	CallbackURL string
}

type OAuthHandler struct {
	mu          sync.RWMutex
	configs     map[string]*OAuthConfig
	tokens      map[string]*OAuthToken
	states      map[string]*OAuthState
	port        int
	callbackURL string
}

func NewOAuthHandler(port int) *OAuthHandler {
	return &OAuthHandler{
		configs:     make(map[string]*OAuthConfig),
		tokens:      make(map[string]*OAuthToken),
		states:      make(map[string]*OAuthState),
		port:        port,
		callbackURL: fmt.Sprintf("http://localhost:%d/callback", port),
	}
}

func (h *OAuthHandler) RegisterProvider(providerID string, config *OAuthConfig) {
	h.mu.Lock()
	defer h.mu.Unlock()
	config.RedirectURI = h.callbackURL
	h.configs[providerID] = config
}

func (h *OAuthHandler) GetAuthURL(providerID, serverID, returnURL string) (string, string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	config, ok := h.configs[providerID]
	if !ok {
		return "", "", fmt.Errorf("provider not registered: %s", providerID)
	}

	state := generateState()
	h.states[state] = &OAuthState{
		ServerID:    serverID,
		Provider:    providerID,
		ReturnURL:   returnURL,
		CreatedAt:   time.Now(),
		CallbackURL: h.callbackURL,
	}

	params := url.Values{}
	params.Set("client_id", config.ClientID)
	params.Set("redirect_uri", config.RedirectURI)
	params.Set("response_type", "code")
	params.Set("state", state)
	if len(config.Scopes) > 0 {
		params.Set("scope", strings.Join(config.Scopes, " "))
	}

	authURL := config.AuthURL + "?" + params.Encode()
	return authURL, state, nil
}

func (h *OAuthHandler) HandleCallback(ctx context.Context, code, state string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	stateData, ok := h.states[state]
	if !ok {
		return "", fmt.Errorf("invalid state parameter")
	}

	if time.Since(stateData.CreatedAt) > 10*time.Minute {
		delete(h.states, state)
		return "", fmt.Errorf("state expired")
	}

	config, ok := h.configs[stateData.Provider]
	if !ok {
		return "", fmt.Errorf("provider not found: %s", stateData.Provider)
	}

	token, err := h.exchangeCode(ctx, config, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	h.tokens[stateData.ServerID] = token
	delete(h.states, state)

	return stateData.ReturnURL, nil
}

func (h *OAuthHandler) exchangeCode(ctx context.Context, config *OAuthConfig, code string) (*OAuthToken, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var token OAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	if token.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second).Unix()
	}

	return &token, nil
}

func (h *OAuthHandler) RefreshToken(ctx context.Context, serverID string) error {
	h.mu.Lock()
	token, ok := h.tokens[serverID]
	if !ok {
		h.mu.Unlock()
		return fmt.Errorf("no token found for server: %s", serverID)
	}

	if token.RefreshToken == "" {
		h.mu.Unlock()
		return fmt.Errorf("no refresh token available")
	}

	config, ok := h.configs[tokenToProvider(token)]
	if !ok {
		h.mu.Unlock()
		return fmt.Errorf("provider config not found")
	}
	h.mu.Unlock()

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("refresh_token", token.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("refresh failed: %s", string(body))
	}

	var newToken OAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&newToken); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.tokens[serverID] = &newToken

	return nil
}

func (h *OAuthHandler) GetToken(serverID string) (*OAuthToken, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	token, ok := h.tokens[serverID]
	if !ok {
		return nil, fmt.Errorf("no token found for server: %s", serverID)
	}

	return token, nil
}

func (h *OAuthHandler) IsTokenExpired(serverID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	token, ok := h.tokens[serverID]
	if !ok {
		return true
	}

	if token.ExpiresAt == 0 {
		return false
	}

	return time.Now().Unix() >= token.ExpiresAt-300
}

func (h *OAuthHandler) RemoveToken(serverID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.tokens, serverID)
}

func (h *OAuthHandler) SaveTokensToFile(serverID, path string) error {
	h.mu.RLock()
	token, ok := h.tokens[serverID]
	h.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no token to save for server: %s", serverID)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (h *OAuthHandler) LoadTokensFromFile(serverID, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var token OAuthToken
	if err := json.Unmarshal(data, &token); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.tokens[serverID] = &token

	return nil
}

func (h *OAuthHandler) StartCallbackServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if code == "" || state == "" {
			http.Error(w, "Missing code or state", http.StatusBadRequest)
			return
		}

		returnURL, err := h.HandleCallback(ctx, code, state)
		if err != nil {
			http.Error(w, fmt.Sprintf("OAuth error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`<html><body><p>Authentication successful! You can close this window.</p><script>window.close();</script></body></html>`)))
		go func() {
			<-ctx.Done()
		}()
		_ = returnURL
	})

	addr := fmt.Sprintf(":%d", h.port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	return server.ListenAndServe()
}

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func tokenToProvider(token *OAuthToken) string {
	return ""
}
