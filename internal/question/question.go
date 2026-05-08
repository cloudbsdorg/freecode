package question

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Question struct {
	ID        string
	Text      string
	Options   []string
	MultiSelect bool
	Required  bool
	Timeout   time.Duration
}

type Answer struct {
	QuestionID string
	Selected  []string
	Text      string
	Timestamp time.Time
}

type Flow struct {
	id       string
	questions []*Question
	current   int
	answers   map[string]*Answer
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewFlow(id string) *Flow {
	ctx, cancel := context.WithCancel(context.Background())
	return &Flow{
		id:        id,
		questions: make([]*Question, 0),
		answers:   make(map[string]*Answer),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (f *Flow) AddQuestion(q *Question) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.questions = append(f.questions, q)
}

func (f *Flow) Next() (*Question, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.current >= len(f.questions) {
		return nil, fmt.Errorf("no more questions")
	}

	q := f.questions[f.current]
	f.current++
	return q, nil
}

func (f *Flow) Answer(questionID string, selected []string, text string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	answer := &Answer{
		QuestionID: questionID,
		Selected:   selected,
		Text:       text,
		Timestamp:  time.Now(),
	}

	f.answers[questionID] = answer
	return nil
}

func (f *Flow) GetAnswer(questionID string) (*Answer, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	answer, ok := f.answers[questionID]
	return answer, ok
}

func (f *Flow) GetAllAnswers() map[string]*Answer {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make(map[string]*Answer, len(f.answers))
	for k, v := range f.answers {
		result[k] = v
	}
	return result
}

func (f *Flow) IsComplete() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, q := range f.questions {
		if q.Required {
			if _, ok := f.answers[q.ID]; !ok {
				return false
			}
		}
	}
	return true
}

func (f *Flow) Progress() (int, int) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.current, len(f.questions)
}

func (f *Flow) Cancel() {
	f.cancel()
}

func (f *Flow) Done() <-chan struct{} {
	return f.ctx.Done()
}

type Manager struct {
	flows   map[string]*Flow
	active  string
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		flows: make(map[string]*Flow),
	}
}

func (m *Manager) CreateFlow(id string) *Flow {
	m.mu.Lock()
	defer m.mu.Unlock()
	flow := NewFlow(id)
	m.flows[id] = flow
	m.active = id
	return flow
}

func (m *Manager) GetFlow(id string) (*Flow, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	flow, ok := m.flows[id]
	return flow, ok
}

func (m *Manager) ActiveFlow() (*Flow, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.active == "" {
		return nil, false
	}
	flow, ok := m.flows[m.active]
	return flow, ok
}

func (m *Manager) SetActive(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active = id
}

func (m *Manager) DeleteFlow(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.flows, id)
	if m.active == id {
		m.active = ""
	}
}

type PromptStyle int

const (
	PromptStyleDefault PromptStyle = iota
	PromptStyleCompact
	PromptStyleDetailed
)

type AskOptions struct {
	Style       PromptStyle
	MultiSelect bool
	Required    bool
	Timeout     time.Duration
	Default     string
}

func Ask(ctx context.Context, text string, opts AskOptions) (string, error) {
	flow := NewFlow("")
	q := &Question{
		ID:          "singleton",
		Text:       text,
		Required:   opts.Required,
		Timeout:    opts.Timeout,
		MultiSelect: opts.MultiSelect,
	}
	flow.AddQuestion(q)

	_, err := flow.Next()
	if err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("timeout waiting for answer")
	}
}

func AskChoice(ctx context.Context, text string, options []string, opts AskOptions) (int, error) {
	if len(options) == 0 {
		return -1, fmt.Errorf("no options provided")
	}

	flow := NewFlow("")
	q := &Question{
		ID:           "choice",
		Text:         text,
		Options:      options,
		MultiSelect:  opts.MultiSelect,
		Required:     opts.Required,
		Timeout:      opts.Timeout,
	}
	flow.AddQuestion(q)

	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	case <-time.After(30 * time.Second):
		return -1, fmt.Errorf("timeout waiting for answer")
	}
}

func AskConfirm(ctx context.Context, text string) (bool, error) {
	flow := NewFlow("")
	q := &Question{
		ID:       "confirm",
		Text:     text,
		Options:  []string{"Yes", "No"},
		Required: true,
	}
	flow.AddQuestion(q)

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case <-time.After(30 * time.Second):
		return false, fmt.Errorf("timeout waiting for answer")
	}
}

func AskMulti(ctx context.Context, text string, options []string) ([]string, error) {
	flow := NewFlow("")
	q := &Question{
		ID:          "multi",
		Text:        text,
		Options:     options,
		MultiSelect: true,
		Required:    true,
	}
	flow.AddQuestion(q)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout waiting for answer")
	}
}

func AskWithDefault(ctx context.Context, text, defaultVal string) (string, error) {
	flow := NewFlow("")
	q := &Question{
		ID:       "default",
		Text:     text,
		Required: false,
	}
	flow.AddQuestion(q)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(30 * time.Second):
		return defaultVal, nil
	}
}