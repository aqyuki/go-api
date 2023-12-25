package server

// SignInRequest is a request for sign-in
type SignInRequest struct {
	ID       string `json:"user_id"`
	Password string `json:"password"`
}

// SignUpRequest is a request for sign-up
type SignUpRequest struct {
	ID       string `json:"user_id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
}
