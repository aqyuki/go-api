package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:           "id",
		PasswordHash: "password_hash",
		Name:         "name",
		Bio:          "bio",
	}

	actual, err := ConvertToBinary(account)
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	actualAccount, err := ConvertFromBinary(actual)
	assert.NoError(t, err)
	assert.NotNil(t, actualAccount)

	assert.EqualValues(t, account, actualAccount)
}
