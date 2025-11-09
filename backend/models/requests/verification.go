package requests

// VerificationRequest represents the request body for verifying an email address.
type VerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

// ResendVerificationRequest represents the payload to resend a verification code.
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}
