package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type TmuxSession struct {
	Name      string
	Windows   int
	CreatedAt time.Time
	Attached  bool
}

type TmuxManager struct {
	mu       sync.RWMutex
	sessions map[string]*TmuxSession
	prefix   string
}

func NewTmuxManager() *TmuxManager {
	return &TmuxManager{
		sessions: make(map[string]*TmuxSession),
		prefix:   "freecode-",
	}
}

func (m *TmuxManager) getPrefix() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.prefix
}

func (m *TmuxManager) SetPrefix(prefix string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.prefix = prefix
}

func (m *TmuxManager) sessionName(name string) string {
	return m.prefix + name
}

func (m *TmuxManager) HasTmux() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

func (m *TmuxManager) CreateSession(name string) (*TmuxSession, error) {
	if !m.HasTmux() {
		return nil, fmt.Errorf("tmux not found")
	}

	sessionName := m.sessionName(name)

	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to create tmux session: %w", err)
	}

	session := &TmuxSession{
		Name:      name,
		Windows:   1,
		CreatedAt: time.Now(),
		Attached:  false,
	}

	m.mu.Lock()
	m.sessions[name] = session
	m.mu.Unlock()

	return session, nil
}

func (m *TmuxManager) DeleteSession(name string) error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	sessionName := m.sessionName(name)

	cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill tmux session: %w", err)
	}

	m.mu.Lock()
	delete(m.sessions, name)
	m.mu.Unlock()

	return nil
}

func (m *TmuxManager) ListSessions() ([]*TmuxSession, error) {
	if !m.HasTmux() {
		return nil, fmt.Errorf("tmux not found")
	}

	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}|#{session_windows}|#{session_attached}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list tmux sessions: %w", err)
	}

	var sessions []*TmuxSession
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}

		sessionName := strings.TrimPrefix(parts[0], m.prefix)
		if !strings.HasPrefix(parts[0], m.prefix) {
			continue
		}

		var windows int
		fmt.Sscanf(parts[1], "%d", &windows)

		sessions = append(sessions, &TmuxSession{
			Name:     sessionName,
			Windows:  windows,
			Attached: parts[2] == "1",
		})
	}

	return sessions, nil
}

func (m *TmuxManager) AttachSession(name string) error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	sessionName := m.sessionName(name)

	cmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	m.mu.Lock()
	if sess, ok := m.sessions[name]; ok {
		sess.Attached = true
	}
	m.mu.Unlock()

	return cmd.Run()
}

func (m *TmuxManager) SendKeys(sessionName, pane string, keys string) error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	target := m.sessionName(sessionName)
	if pane != "" {
		target = target + ":" + pane
	}

	cmd := exec.Command("tmux", "send-keys", "-t", target, keys, "Enter")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send keys: %w", err)
	}

	return nil
}

func (m *TmuxManager) CapturePane(sessionName, pane string) (string, error) {
	if !m.HasTmux() {
		return "", fmt.Errorf("tmux not found")
	}

	target := m.sessionName(sessionName)
	if pane != "" {
		target = target + ":" + pane
	}

	cmd := exec.Command("tmux", "capture-pane", "-t", target, "-p")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to capture pane: %w", err)
	}

	return string(output), nil
}

func (m *TmuxManager) CreateWindow(sessionName, windowName string) error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	sessionName = m.sessionName(sessionName)

	cmd := exec.Command("tmux", "new-window", "-t", sessionName, "-n", windowName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}

	m.mu.Lock()
	if sess, ok := m.sessions[sessionName]; ok {
		sess.Windows++
	}
	m.mu.Unlock()

	return nil
}

func (m *TmuxManager) SelectLayout(sessionName, layout string) error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	cmd := exec.Command("tmux", "select-layout", "-t", m.sessionName(sessionName), layout)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to select layout: %w", err)
	}

	return nil
}

func (m *TmuxManager) SplitWindow(sessionName, pane string) error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	target := m.sessionName(sessionName)
	if pane != "" {
		target = target + "." + pane
	}

	cmd := exec.Command("tmux", "split-window", "-t", target)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to split window: %w", err)
	}

	return nil
}

func (m *TmuxManager) GetSession(name string) (*TmuxSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sess, ok := m.sessions[name]
	return sess, ok
}

func (m *TmuxManager) RefreshSession(name string) error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	sessionName := m.sessionName(name)

	cmd := exec.Command("tmux", "list-windows", "-t", sessionName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("session not found: %s", name)
	}

	cmd = exec.Command("tmux", "list-panes", "-t", sessionName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("session not found: %s", name)
	}

	m.mu.Lock()
	if sess, ok := m.sessions[name]; ok {
		output, _ := exec.Command("tmux", "list-windows", "-t", sessionName).Output()
		sess.Windows = strings.Count(string(output), "\n") + 1
	}
	m.mu.Unlock()

	return nil
}

func (m *TmuxManager) RunInPane(ctx context.Context, sessionName, pane string, command string) (string, error) {
	if !m.HasTmux() {
		return "", fmt.Errorf("tmux not found")
	}

	target := m.sessionName(sessionName)
	if pane != "" {
		target = target + ":" + pane
	}

	cmd := exec.CommandContext(ctx, "tmux", "send-keys", "-t", target, command, "Enter")

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	time.Sleep(100 * time.Millisecond)

	output, err := m.CapturePane(sessionName, pane)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (m *TmuxManager) KillAll() error {
	if !m.HasTmux() {
		return fmt.Errorf("tmux not found")
	}

	m.mu.RLock()
	names := make([]string, 0, len(m.sessions))
	for name := range m.sessions {
		names = append(names, name)
	}
	m.mu.RUnlock()

	for _, name := range names {
		if err := m.DeleteSession(name); err != nil {
			continue
		}
	}

	return nil
}
