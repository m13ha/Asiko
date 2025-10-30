package responses

// SimpleMessage is a generic response for returning a simple message.
type SimpleMessage struct {
	Message string `json:"message" example:"Action completed successfully"`
} // @SimpleResponse
