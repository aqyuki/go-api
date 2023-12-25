package account

import (
	"context"
	"fmt"
)

// AccountApp provides functions to manage and control accounts.
// This structure implements IAccountService interface.
type AccountApp struct {
	repository Repository
	encoder    PasswordEncoder
}

// SignIn verifies the account and returns a account.
func (s *AccountApp) SignIn(ctx context.Context, id string, password string) (*Account, error) {
	hash, err := s.encoder.Encode(password)
	if err != nil {
		return nil, fmt.Errorf("failed to sign in: %w", err)
	}

	account, err := s.repository.FetchAccountWithPassword(ctx, id, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign in: %w", err)
	}
	return account, nil
}

// SignUp creates a new account and returns a account.
func (s *AccountApp) SignUp(ctx context.Context, id string, password string, name string, bio string) (*Account, error) {
	hash,err := s.encoder.Encode(password)
	if err != nil {
		return nil, fmt.Errorf("failed to sign up: %w", err)
	}
	account := Account{
		ID:           id,
		Name:         name,
		Bio:          bio,
		PasswordHash: hash,
	}
	if err := s.repository.CreateAccount(ctx, &account); err != nil {
		return nil, fmt.Errorf("failed to sign up: %w", err)
	}
	return &account, nil
}

// FetchAccountInfo returns account information.
func (s *AccountApp) FetchAccountInfo(ctx context.Context, id string) (*Account, error) {
	account, err := s.repository.FetchAccount(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account info: %w", err)
	}
	return account, nil
}

// NewAccountApp creates a new AccountApp instance.
func NewAccountApp(r Repository, e PasswordEncoder) *AccountApp {
	return &AccountApp{
		repository: r,
		encoder:    e,
	}
}
