package account

import (
	"bytes"
	"encoding/gob"
)

// ConvertToBinary converts account to binary.
func ConvertToBinary(account *Account) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(account); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ConvertFromBinary converts binary to account.
func ConvertFromBinary(b []byte) (*Account, error) {
	buf := bytes.NewBuffer(b)
	account := &Account{}
	if err := gob.NewDecoder(buf).Decode(account); err != nil {
		return nil, err
	}
	return account, nil
}
