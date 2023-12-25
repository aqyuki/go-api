package password

import (
	"testing"

	"github.com/aqyuki/jwt-demo/account"
	"github.com/stretchr/testify/assert"
)

var (
	_ account.PasswordEncoder = (*SHA256Encoder)(nil)
)

func TestEncode(t *testing.T) {
	e := &SHA256Encoder{}

	pat1 := "bar"
	pat2 := "foo"

	hash1, err := e.Encode(pat1)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash1)

	hash2, err := e.Encode(pat2)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash2)

	hash3, err := e.Encode(pat1)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash3)

	assert.NotEqual(t, hash1, hash2)
	assert.Equal(t, hash1, hash3)
}
