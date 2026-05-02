package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/freecode/freecode/internal/auth"
)

func TestProvidersCommand(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"providers"})
	if err != nil || cmd == rootCmd {
		t.Error("providers command not found")
	}
	if len(cmd.Commands()) == 0 {
		t.Error("providers should have subcommands")
	}
}

func TestProvidersListCmd(t *testing.T) {
	cmd, _, err := providersCmd.Find([]string{"list"})
	if err != nil || cmd == providersCmd {
		t.Error("list subcommand not found")
	}
}

func TestProvidersLoginCmd(t *testing.T) {
	cmd, _, err := providersCmd.Find([]string{"login"})
	if err != nil || cmd == providersCmd {
		t.Error("login subcommand not found")
	}
}

func TestProvidersLogoutCmd(t *testing.T) {
	cmd, _, err := providersCmd.Find([]string{"logout"})
	if err != nil || cmd == providersCmd {
		t.Error("logout subcommand not found")
	}
}

func TestRunProvidersListEmpty(t *testing.T) {
	buf := &bytes.Buffer{}
	providersListCmd.SetOut(buf)
	providersListCmd.SetErr(buf)

	err := runProvidersList(providersListCmd, []string{})
	if err != nil {
		t.Errorf("runProvidersList() error = %v", err)
	}
}

func TestProvidersLoginProviderFlag(t *testing.T) {
	providersLoginCmd.Flags().Set("provider", "openai")
	defer providersLoginCmd.Flags().Set("provider", "")

	if providersLoginProvider != "openai" {
		t.Errorf("providersLoginProvider = %q, want %q", providersLoginProvider, "openai")
	}
}

func TestSupportedProviders(t *testing.T) {
	providers := supportedProviders()
	if len(providers) == 0 {
		t.Error("supportedProviders() returned empty slice")
	}

	for _, p := range providers {
		if !isValidProvider(p) {
			t.Errorf("isValidProvider(%q) = false, want true", p)
		}
	}

	if isValidProvider("nonexistent") {
		t.Error("isValidProvider(\"nonexistent\") = true, want false")
	}
}

func TestGetProviderDisplayName(t *testing.T) {
	tests := []struct {
		id   string
		name string
	}{
		{"openai", "OpenAI"},
		{"anthropic", "Anthropic"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			name := getProviderDisplayName(tt.id)
			if name != tt.name {
				t.Errorf("getProviderDisplayName(%q) = %q, want %q", tt.id, name, tt.name)
			}
		})
	}
}

func TestGetProviderModel(t *testing.T) {
	tests := []struct {
		id     string
		model  string
	}{
		{"openai", "gpt-4o"},
		{"anthropic", "claude-sonnet-4-5"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			model := getProviderModel(tt.id)
			if model != tt.model {
				t.Errorf("getProviderModel(%q) = %q, want %q", tt.id, model, tt.model)
			}
		})
	}
}

func TestAuthStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "freecode-auth-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	authPath := tmpDir + "/auth.json"
	store := auth.NewStore(authPath)

	info := auth.Info{
		Type: auth.AuthTypeAPI,
		Key:  "test-api-key",
	}

	err = store.Set("test-provider", info)
	if err != nil {
		t.Fatalf("store.Set() error = %v", err)
	}

	retrieved, ok := store.Get("test-provider")
	if !ok {
		t.Fatal("store.Get() = false, want true")
	}
	if retrieved.Type != auth.AuthTypeAPI {
		t.Errorf("retrieved.Type = %q, want %q", retrieved.Type, auth.AuthTypeAPI)
	}
	if retrieved.Key != "test-api-key" {
		t.Errorf("retrieved.Key = %q, want %q", retrieved.Key, "test-api-key")
	}

	all := store.All()
	if len(all) != 1 {
		t.Errorf("len(store.All()) = %d, want 1", len(all))
	}

	err = store.Remove("test-provider")
	if err != nil {
		t.Fatalf("store.Remove() error = %v", err)
	}

	_, ok = store.Get("test-provider")
	if ok {
		t.Error("store.Get() after Remove() = true, want false")
	}
}

func TestAuthStoreClear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "freecode-auth-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	authPath := tmpDir + "/auth.json"
	store := auth.NewStore(authPath)

	info := auth.Info{Type: auth.AuthTypeAPI, Key: "test-key"}
	store.Set("provider1", info)
	store.Set("provider2", info)

	err = store.Clear()
	if err != nil {
		t.Fatalf("store.Clear() error = %v", err)
	}

	all := store.All()
	if len(all) != 0 {
		t.Errorf("len(store.All()) after Clear() = %d, want 0", len(all))
	}
}