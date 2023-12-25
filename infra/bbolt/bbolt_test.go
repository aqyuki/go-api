package bbolt

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/aqyuki/go-api/account"
	"github.com/stretchr/testify/assert"
)

var (
	_ account.Repository = (*BBoltAccountRepository)(nil)
)

func TestBBolt(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "test.db")
	db, err := NewAccountRepository(path)
	if !assert.NoError(t, err, "failed to create db") {
		t.Fatal(err)
	}

	// mock data
	data := account.Account{
		ID:           "test",
		PasswordHash: "test",
		Name:         "test",
		Bio:          "test",
	}

	// Create account
	err = db.CreateAccount(context.Background(), &data)
	if !assert.NoError(t, err, "failed to create account") {
		t.Fatal(err)
	}

	// Fetch account with password
	got, err := db.FetchAccountWithPassword(context.Background(), data.ID, data.PasswordHash)
	if !assert.NoError(t, err, "failed to fetch account with password") {
		t.Fatal(err)
	}
	assert.EqualValues(t, &data, got, "account does not match")

	// Fetch account
	got, err = db.FetchAccount(context.Background(), data.ID)
	if !assert.NoError(t, err, "failed to fetch account") {
		t.Fatal(err)
	}
	assert.EqualValues(t, &data, got, "account does not match")

	err = db.Close()
	assert.NoError(t, err, "failed to close db")
}
