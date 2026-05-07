package permission

import (
	"context"
	"sync"
)

type Action string

const (
	ActionRead    Action = "read"
	ActionWrite   Action = "write"
	ActionDelete  Action = "delete"
	ActionExecute Action = "execute"
)

type Permission struct {
	Resource string
	Actions  []Action
}

type Checker interface {
	Check(ctx context.Context, subject string, permission Permission) (bool, error)
	Grant(ctx context.Context, subject string, permission Permission) error
	Revoke(ctx context.Context, subject string, permission Permission) error
}

type memoryChecker struct {
	mu          sync.RWMutex
	permissions map[string][]Permission
}

func NewMemoryChecker() *memoryChecker {
	return &memoryChecker{permissions: make(map[string][]Permission)}
}

func (c *memoryChecker) Check(ctx context.Context, subject string, permission Permission) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	perms, ok := c.permissions[subject]
	if !ok {
		return false, nil
	}
	for _, p := range perms {
		if p.Resource == permission.Resource {
			for _, a := range p.Actions {
				if a == permission.Actions[0] {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (c *memoryChecker) Grant(ctx context.Context, subject string, permission Permission) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.permissions[subject] = append(c.permissions[subject], permission)
	return nil
}

func (c *memoryChecker) Revoke(ctx context.Context, subject string, permission Permission) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	perms := c.permissions[subject]
	for i, p := range perms {
		if p.Resource == permission.Resource {
			c.permissions[subject] = append(perms[:i], perms[i+1:]...)
			return nil
		}
	}
	return nil
}
