package errors

import (
	"errors"
	"fmt"
)

// FromAppError extracts an *AppError from a generic error.
// If the error is already an *AppError, it is returned directly.
// Otherwise the error is wrapped as an internal AppError.
func FromAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	// Not an AppError – wrap as internal error
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: err.Error(),
		Kind:    "internal",
		HTTP:    500,
		Cause:   err,
	}
}

// Ensure fmt is used for compilation (fmt imported for potential future formatting).
var _ = fmt.Sprintf
