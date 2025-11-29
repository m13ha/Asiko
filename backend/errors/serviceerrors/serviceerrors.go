package serviceerrors

import (
	"github.com/m13ha/asiko/errors"
)

// ValidationError returns a validation error.
func ValidationError(message string) error {
	return errors.NewAppError(errors.CodeValidationFailed, "validation", 400, message, nil)
}

// ConflictError returns a conflict error.
func ConflictError(message string) error {
	return errors.NewAppError(errors.CodeConflict, "conflict", 409, message, nil)
}

// NotFoundError returns a not found error.
func NotFoundError(message string) error {
	return errors.NewAppError(errors.CodeResourceNotFound, "not_found", 404, message, nil)
}

// ForbiddenError returns a forbidden error.
func ForbiddenError(message string) error {
	return errors.NewAppError(errors.CodeForbidden, "forbidden", 403, message, nil)
}

// UnauthorizedError returns an unauthorized error.
func UnauthorizedError(message string) error {
	return errors.NewAppError(errors.CodeUnauthorized, "unauthorized", 401, message, nil)
}

// InternalError returns an internal server error.
func InternalError(message string) error {
	return errors.NewAppError(errors.CodeInternalError, "internal", 500, message, nil)
}

// UserError returns a user error (for general user-related errors).
func UserError(message string) error {
	return errors.NewAppError("USER_ERROR", "user", 400, message, nil)
}

// PreconditionFailedError returns a precondition failed error.
func PreconditionFailedError(message string) error {
	return errors.NewAppError(errors.CodePreconditionFailed, "precondition", 400, message, nil)
}

// BookingCapacityExceededError returns a booking capacity exceeded error.
func BookingCapacityExceededError(message string) error {
	return errors.NewAppError(errors.CodeBookingCapacityExceeded, "conflict", 409, message, nil)
}

// BookingSlotUnavailableError returns a booking slot unavailable error.
func BookingSlotUnavailableError(message string) error {
	return errors.NewAppError(errors.CodeBookingSlotUnavailable, "conflict", 409, message, nil)
}

// EmailAlreadyRegisteredError returns an email already registered error.
func EmailAlreadyRegisteredError(message string) error {
	return errors.NewAppError(errors.CodeEmailAlreadyRegistered, "conflict", 409, message, nil)
}

// VerificationExpiredError returns a verification expired error.
func VerificationExpiredError(message string) error {
	return errors.NewAppError(errors.CodeVerificationExpired, "validation", 400, message, nil)
}

// InvalidVerificationCodeError returns an invalid verification code error.
func InvalidVerificationCodeError(message string) error {
	return errors.NewAppError(errors.CodeInvalidVerificationCode, "validation", 400, message, nil)
}

// LoginInvalidCredentialsError returns a login invalid credentials error.
func LoginInvalidCredentialsError(message string) error {
	return errors.NewAppError(errors.CodeLoginInvalidCredentials, "unauthorized", 401, message, nil)
}

// UserPendingVerificationError returns a user pending verification error.
func UserPendingVerificationError(message string) error {
	return errors.NewAppError(errors.CodeUserPendingVerification, "precondition", 202, message, nil)
}

// WrapError wraps an existing error with a service error.
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return &errors.AppError{
		Code:    "WRAPPED_ERROR",
		Message: message,
		Kind:    "internal",
		HTTP:    500,
		Cause:   err,
	}
}

// FromError extracts a *AppError from a generic error.
func FromError(err error) *errors.AppError {
	return errors.FromAppError(err)
}
