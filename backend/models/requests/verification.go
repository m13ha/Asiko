package requests

// VerificationRequest represents the request body for verifying an email address.
type VerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}
