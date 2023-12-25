package password

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/aqyuki/jwt-demo/account"
)

type SHA256Encoder struct{}

// Encode converts password to hashed password
func (e *SHA256Encoder) Encode(password string) (account.HashedPassword, error) {
	b := []byte(password)

	hash := sha256.New()
	if _, err := hash.Write(b); err != nil {
		return "", fmt.Errorf("failed to encode password: %w", err)
	}

	hashed := hex.EncodeToString(hash.Sum(nil))
	return account.HashedPassword(hashed), nil
}
