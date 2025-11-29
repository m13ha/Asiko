package repoerrors

import "github.com/m13ha/asiko/errors"

// ValidationError returns a validation error for the repository layer.
func ValidationError(message string) error {
	return errors.NewAppError(errors.CodeRepoValidationError, "validation", 400, message, nil)
}

// ConflictError returns a conflict error for the repository layer.
func ConflictError(message string) error {
	return errors.NewAppError(errors.CodeRepoConflictError, "conflict", 409, message, nil)
}

// NotFoundError returns a not found error for the repository layer.
func NotFoundError(message string) error {
	return errors.NewAppError(errors.CodeRepoNotFoundError, "not_found", 404, message, nil)
}

// InternalError returns an internal error for the repository layer.
func InternalError(message string) error {
	return errors.NewAppError(errors.CodeRepoInternalError, "internal", 500, message, nil)
}

// FromError extracts a *AppError from a generic error.
func FromError(err error) *errors.AppError {
	return errors.FromAppError(err)
}
