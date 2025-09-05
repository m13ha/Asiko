package responses

// LoginResponse represents the data returned upon a successful login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
