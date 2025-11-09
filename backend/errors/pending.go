package errors

// PendingVerificationError indicates that a request was accepted but requires
// email verification to complete. Handlers should typically map this to
// HTTP 202 Accepted.
type PendingVerificationError struct {
    Message string
}

func (e *PendingVerificationError) Error() string {
    return e.Message
}

func NewPendingVerificationError(msg string) error {
    return &PendingVerificationError{Message: msg}
}

