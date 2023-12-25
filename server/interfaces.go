package server

import "context"

type PasswordHash string

// AccountRepository is an interface for account repository
type AccountRepository interface {
	// FetchAccount fetches account by id
	FetchAccount(ctx context.Context, id string) (*Account, error)

	// FetchAccountWithPassword fetches account by id and password hash
	FetchAccountWithPassword(ctx context.Context, id string, hash PasswordHash) (*Account, error)

	// RegisterAccount registers account
	RegisterAccount(ctx context.Context, account *Account) error
}

// PasswordEncoder is an interface for password encoder
type PasswordEncoder interface {
	// Encode converts password to hash
	Encode(password string) PasswordHash
}
