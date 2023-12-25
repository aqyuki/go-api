package account

import (
	"context"
	"errors"
	"sync"
)

type MockRepo struct {
	repo map[string]Account
	m    sync.Mutex
}

func (r *MockRepo) FetchAccount(ctx context.Context, id string) (*Account, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if account, ok := r.repo[id]; ok {
		return &account, nil
	}
	return nil, errors.New("account not found")
}

func (r *MockRepo) FetchAccountWithPassword(ctx context.Context, id string, passwordHash HashedPassword) (*Account, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if account, ok := r.repo[id]; ok {
		if account.PasswordHash == passwordHash {
			return &account, nil
		}
	}
	return nil, errors.New("account not found")
}

func (r *MockRepo) CreateAccount(ctx context.Context, account *Account) error {
	r.m.Lock()
	defer r.m.Unlock()
	r.repo[account.ID] = *account
	return nil
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		repo: map[string]Account{},
	}
}

type MockEncoder struct{}

func (e *MockEncoder) Encode(password string) (HashedPassword, error) {
	return HashedPassword(password), nil
}
