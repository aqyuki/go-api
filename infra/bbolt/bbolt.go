package bbolt

import (
	"context"
	"fmt"

	"github.com/aqyuki/jwt-demo/account"
	bolt "go.etcd.io/bbolt"
)

const (
	accountBucketName = "account"
)

// BBoltAccountRepository is a repository for account.
// BBoltAccountRepository implements AccountRepository interface.
type BBoltAccountRepository struct {
	db *bolt.DB
}

func (r *BBoltAccountRepository) FetchAccountWithPassword(ctx context.Context, id string, passwordHash account.HashedPassword) (*account.Account, error) {
	var got *account.Account

	r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(accountBucketName))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", accountBucketName)
		}

		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("account %s does not exist", id)
		}

		a, err := account.ConvertFromBinary(v)
		if err != nil {
			return fmt.Errorf("failed to convert from binary: %w", err)
		}

		if a.PasswordHash != passwordHash {
			return fmt.Errorf("password does not match")
		}
		got = a
		return nil
	})
	return got, nil
}

func (r *BBoltAccountRepository) FetchAccount(ctx context.Context, id string) (*account.Account, error) {
	var got *account.Account

	r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(accountBucketName))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", accountBucketName)
		}

		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("account %s does not exist", id)
		}

		a, err := account.ConvertFromBinary(v)
		if err != nil {
			return fmt.Errorf("failed to convert from binary: %w", err)
		}
		got = a
		return nil
	})
	return got, nil
}

func (r *BBoltAccountRepository) CreateAccount(ctx context.Context, ac *account.Account) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(accountBucketName))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", accountBucketName)
		}

		v, err := account.ConvertToBinary(ac)
		if err != nil {
			return fmt.Errorf("failed to convert to binary: %w", err)
		}

		if err := b.Put([]byte(ac.ID), v); err != nil {
			return fmt.Errorf("failed to put an account: %w", err)
		}
		return nil
	})
}

func (r *BBoltAccountRepository) Close() error {
	return r.db.Close()
}

// NewAccountRepository returns a new BBoltAccountRepository.
// BBoltAccountRepository implements AccountRepository interface.
// path is a path to a bbolt database file.
// Recommend to use a file with .db extension.
func NewAccountRepository(path string) (*BBoltAccountRepository, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open a database file: %w", err)
	}

	// try to initialize a database
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(accountBucketName)); err != nil {
			return fmt.Errorf("failed to create a bucket: %w", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize a database: %w", err)
	}

	return &BBoltAccountRepository{
		db: db,
	}, nil
}
