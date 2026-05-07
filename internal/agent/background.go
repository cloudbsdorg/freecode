package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type BackgroundTask struct {
	ID          string
	Name        string
	Type        string
	Status      TaskStatus
	Result      interface{}
	Error       error
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Progress    int
	CancelFn    context.CancelFunc
}

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

type BackgroundTaskHandler interface {
	HandleTask(task *BackgroundTask)
}

type BackgroundManager struct {
	mu      sync.RWMutex
	tasks   map[string]*BackgroundTask
	workers int
	handler BackgroundTaskHandler
}

func NewBackgroundManager(workers int, handler BackgroundTaskHandler) *BackgroundManager {
	if workers <= 0 {
		workers = 5
	}
	return &BackgroundManager{
		tasks:   make(map[string]*BackgroundTask),
		workers: workers,
		handler: handler,
	}
}

func (m *BackgroundManager) Submit(task *BackgroundTask) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if task.ID == "" {
		return fmt.Errorf("task ID is required")
	}

	if task.Name == "" {
		return fmt.Errorf("task name is required")
	}

	task.Status = TaskStatusPending
	task.CreatedAt = time.Now()
	m.tasks[task.ID] = task

	return nil
}

func (m *BackgroundManager) SubmitWithCancel(task *BackgroundTask) (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	task.CancelFn = cancel
	if err := m.Submit(task); err != nil {
		cancel()
		return nil, nil, err
	}

	return ctx, func() {
		m.Cancel(task.ID)
		cancel()
	}, nil
}

func (m *BackgroundManager) Run(ctx context.Context, task *BackgroundTask) {
	m.mu.Lock()
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	m.tasks[task.ID] = task
	m.mu.Unlock()

	defer func() {
		if task.Status == TaskStatusRunning {
			m.Complete(task.ID, nil)
		}
	}()

	if m.handler != nil {
		m.handler.HandleTask(task)
	}
}

func (m *BackgroundManager) RunAsync(ctx context.Context, task *BackgroundTask) {
	go m.Run(ctx, task)
}

func (m *BackgroundManager) Get(id string) (*BackgroundTask, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	task, ok := m.tasks[id]
	return task, ok
}

func (m *BackgroundManager) List() []*BackgroundTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]*BackgroundTask, 0, len(m.tasks))
	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (m *BackgroundManager) ListByStatus(status TaskStatus) []*BackgroundTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tasks []*BackgroundTask
	for _, t := range m.tasks {
		if t.Status == status {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (m *BackgroundManager) ListByType(taskType string) []*BackgroundTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tasks []*BackgroundTask
	for _, t := range m.tasks {
		if t.Type == taskType {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (m *BackgroundManager) Complete(id string, result interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return
	}

	task.Status = TaskStatusCompleted
	task.Result = result
	now := time.Now()
	task.CompletedAt = &now
}

func (m *BackgroundManager) Fail(id string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return
	}

	task.Status = TaskStatusFailed
	task.Error = err
	now := time.Now()
	task.CompletedAt = &now
}

func (m *BackgroundManager) Cancel(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return
	}

	if task.CancelFn != nil {
		task.CancelFn()
	}
	task.Status = TaskStatusCancelled
	now := time.Now()
	task.CompletedAt = &now
}

func (m *BackgroundManager) UpdateProgress(id string, progress int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return
	}

	task.Progress = progress
}

func (m *BackgroundManager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tasks, id)
}

func (m *BackgroundManager) ClearCompleted() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, task := range m.tasks {
		if task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed || task.Status == TaskStatusCancelled {
			delete(m.tasks, id)
		}
	}
}

func (m *BackgroundManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.tasks)
}

func (m *BackgroundManager) CountByStatus(status TaskStatus) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, t := range m.tasks {
		if t.Status == status {
			count++
		}
	}
	return count
}

type BackgroundWorkerPool struct {
	mu      sync.RWMutex
	manager *BackgroundManager
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewBackgroundWorkerPool(workers int, handler BackgroundTaskHandler) *BackgroundWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &BackgroundWorkerPool{
		manager: NewBackgroundManager(workers, handler),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (p *BackgroundWorkerPool) Submit(task *BackgroundTask) error {
	return p.manager.Submit(task)
}

func (p *BackgroundWorkerPool) Run(task *BackgroundTask) {
	p.manager.Run(p.ctx, task)
}

func (p *BackgroundWorkerPool) RunAsync(task *BackgroundTask) {
	p.manager.RunAsync(p.ctx, task)
}

func (p *BackgroundWorkerPool) Get(id string) (*BackgroundTask, bool) {
	return p.manager.Get(id)
}

func (p *BackgroundWorkerPool) List() []*BackgroundTask {
	return p.manager.List()
}

func (p *BackgroundWorkerPool) Shutdown() {
	p.cancel()
}

func (p *BackgroundWorkerPool) Wait() {
	<-p.ctx.Done()
}
