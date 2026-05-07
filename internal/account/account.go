package account

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrRepoError       = errors.New("account repository error")
)

type AccountID string
type OrgID string
type AccessToken string
type RefreshToken string

type Info struct {
	ID          AccountID
	Email       string
	URL         string
	ActiveOrgID OrgID
}

type Org struct {
	ID   OrgID
	Name string
}

type Login struct {
	Code     string
	UserCode string
	URL      string
	Server   string
	Expiry   time.Duration
	Interval time.Duration
}

type PollResult interface{ isPollResult() }
type PollSuccess struct{ Email string }
type PollPending struct{}
type PollSlow struct{}
type PollExpired struct{}
type PollDenied struct{}
type PollError struct{ Cause error }

func (PollSuccess) isPollResult() {}
func (PollPending) isPollResult() {}
func (PollSlow) isPollResult()    {}
func (PollExpired) isPollResult() {}
func (PollDenied) isPollResult()  {}
func (PollError) isPollResult()   {}

type AccountRepo interface {
	Active(ctx context.Context) (*Info, error)
	List(ctx context.Context) ([]*Info, error)
	Remove(ctx context.Context, id AccountID) error
	Use(ctx context.Context, id AccountID, orgID OrgID) error
	PersistToken(ctx context.Context, id AccountID, access, refresh string, expiry time.Time) error
	PersistAccount(ctx context.Context, info *Info, access, refresh string) error
}

type memoryAccountRepo struct {
	mu       sync.RWMutex
	active   *Info
	accounts map[AccountID]*Info
}

func NewMemoryRepo() AccountRepo {
	return &memoryAccountRepo{
		accounts: make(map[AccountID]*Info),
	}
}

func (r *memoryAccountRepo) Active(ctx context.Context) (*Info, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.active == nil {
		return nil, ErrAccountNotFound
	}
	return r.active, nil
}

func (r *memoryAccountRepo) List(ctx context.Context) ([]*Info, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Info
	for _, acc := range r.accounts {
		result = append(result, acc)
	}
	return result, nil
}

func (r *memoryAccountRepo) Remove(ctx context.Context, id AccountID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.active != nil && r.active.ID == id {
		r.active = nil
	}
	delete(r.accounts, id)
	return nil
}

func (r *memoryAccountRepo) Use(ctx context.Context, id AccountID, orgID OrgID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	acc, ok := r.accounts[id]
	if !ok {
		return ErrAccountNotFound
	}
	acc.ActiveOrgID = orgID
	r.active = acc
	return nil
}

func (r *memoryAccountRepo) PersistToken(ctx context.Context, id AccountID, access, refresh string, expiry time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	acc, ok := r.accounts[id]
	if !ok {
		return ErrAccountNotFound
	}
	_ = access
	_ = refresh
	_ = expiry
	_ = acc
	return nil
}

func (r *memoryAccountRepo) PersistAccount(ctx context.Context, info *Info, access, refresh string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[info.ID] = info
	r.active = info
	_ = access
	_ = refresh
	return nil
}
