package responses

// SimpleMessageResponse is a generic response for returning a simple message.
type SimpleMessageResponse struct {
	Message string `json:"message" example:"Action completed successfully"`
}
