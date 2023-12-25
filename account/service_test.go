package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_ Service = new(AccountApp)
)

func TestAccountService(t *testing.T) {
	t.Parallel()

	t.Run("should be success to sign up - sign in - fetch account info", func(t *testing.T) {
		t.Parallel()

		service := NewAccountApp(NewMockRepo(), &MockEncoder{})
		ctx := context.Background()

		data := struct {
			id       string
			password string
			name     string
			bio      string
		}{
			id:       "foo",
			password: "bar",
			name:     "baz",
			bio:      "qux",
		}

		// sign up
		account, err := service.SignUp(ctx, data.id, data.password, data.name, data.bio)
		assert.NoError(t, err, "should return no error but received %v", err)
		if assert.NotNil(t, account, "should return account but received nil") {
			assert.EqualValues(t, data.id, account.ID, "should return id %v but received %v", data.id, account.ID)
			assert.EqualValues(t, data.name, account.Name, "should return name %v but received %v", data.name, account.Name)
			assert.EqualValues(t, data.bio, account.Bio, "should return bio %v but received %v", data.bio, account.Bio)
		} else {
			t.Fatal("failed to sign up test")
		}

		// sign in
		account, err = service.SignIn(ctx, data.id, data.password)
		assert.NoError(t, err, "should return no error but received %v", err)
		if assert.NotNil(t, account, "should return account but received nil") {
			assert.EqualValues(t, data.id, account.ID, "should return id %v but received %v", data.id, account.ID)
			assert.EqualValues(t, data.name, account.Name, "should return name %v but received %v", data.name, account.Name)
			assert.EqualValues(t, data.bio, account.Bio, "should return bio %v but received %v", data.bio, account.Bio)
		} else {
			t.Fatal("failed to sign in test")
		}

		// fetch account info
		account, err = service.FetchAccountInfo(ctx, data.id)
		assert.NoError(t, err, "should return no error but received %v", err)
		if assert.NotNil(t, account, "should return account but received nil") {
			assert.EqualValues(t, data.id, account.ID, "should return id %v but received %v", data.id, account.ID)
			assert.EqualValues(t, data.name, account.Name, "should return name %v but received %v", data.name, account.Name)
			assert.EqualValues(t, data.bio, account.Bio, "should return bio %v but received %v", data.bio, account.Bio)
		} else {
			t.Fatal("failed to fetch account info test")
		}
	})
}
