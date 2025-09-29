package responses

// ResponsesSimpleMessageResponse is a generic response for returning a simple message.
type ResponsesSimpleMessageResponse struct {
	Message string `json:"message" example:"Action completed successfully"`
}

// ResponsesLoginResponse represents the response after a successful login or registration verification.
type ResponsesLoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}