package account

import "context"

// Service provides functions to manage and control accounts.
type Service interface {
	// SignIn verifies the account and returns a jwt token.
	SignIn(ctx context.Context, id string, password string) (*Account, error)
	// SignUp creates a new account and returns a jwt token.
	SignUp(ctx context.Context, id string, password string, name string, bio string) (*Account, error)
	// FetchAccountInfo returns account information.
	FetchAccountInfo(ctx context.Context, id string) (*Account, error)
}

// Repository provides functions to access database for accounts.
type Repository interface {
	// FetchAccountWithPassword finds an account matching id and password, and returns it.
	// If there is not account matching id and password, it returns an error.
	FetchAccountWithPassword(ctx context.Context, id string, passwordHash HashedPassword) (*Account, error)
	// FetchAccount finds an account matching id, and returns it.
	// If there is not account matching id, it returns an error.
	FetchAccount(ctx context.Context, id string) (*Account, error)
	// CreateAccount creates a new account.
	CreateAccount(ctx context.Context, account *Account) error
}

// Account is a type to hold account data.
type Account struct {
	ID           string
	Name         string
	Bio          string
	PasswordHash HashedPassword
}

// HashedPassword is a type to hold hashed password
type HashedPassword string

// PasswordEncoder provides function to encode password to hashed password
type PasswordEncoder interface {
	// Encode converts password to hashed password
	Encode(password string) (HashedPassword,error)
}
