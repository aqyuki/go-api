package server

type Account struct {
	ID           string
	PasswordHash PasswordHash
	Name         string
	Bio          string
}

// ErrorResponse is a response for error
type ErrorResponse struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// AccountInfoResponse is a response for account info
type AccountInfoResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

// AccountVerifyResponse is a response for sign-in and sign-up
type AccountVerifyResponse struct {
	Token string `json:"token"` // Token is a jwt token.
}
