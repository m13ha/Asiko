package errors

// UserError is a custom error type for user-facing errors
// Only the Message is returned to the client
// Internal errors are logged but not exposed

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

func NewUserError(msg string) error {
	return &UserError{Message: msg}
}
