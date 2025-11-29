package errors

import "fmt"

// AppError is a unified error type used across repository, service and API layers.
// It implements the error interface and provides helpers for error unwrapping
// and comparison via the standard errors.Is / errors.As mechanisms.
type AppError struct {
	Code    string // machine‑readable error identifier (e.g. "VALIDATION_ERROR")
	Message string // human‑readable description
	Kind    string // logical grouping such as "validation", "conflict", "internal"
	HTTP    int    // HTTP status code to be returned to the client
	Cause   error  // optional wrapped error for stack traces
}

// Error satisfies the built‑in error interface.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause for errors.Is / errors.As.
func (e *AppError) Unwrap() error { return e.Cause }

// Is reports whether the target error has the same Code.
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}

// NewAppError creates a new AppError instance.
func NewAppError(code, kind string, http int, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Kind:    kind,
		HTTP:    http,
		Cause:   cause,
	}
}
