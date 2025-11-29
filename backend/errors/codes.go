package errors

// ---------------------------------------------------------------------
// Generic error codes
// ---------------------------------------------------------------------

const (
	// Validation
	CodeValidationFailed = "VALIDATION_FAILED"
	CodeBadRequest       = "BAD_REQUEST"

	// Authentication / Authorization
	CodeUnauthorized = "AUTH_UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"

	// Resource
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"

	// Conflict & Rate limiting
	CodeConflict    = "CONFLICT"
	CodeRateLimited = "RATE_LIMITED"

	// Precondition & Timeout
	CodePreconditionFailed = "PRECONDITION_FAILED"
	CodeTimeout            = "TIMEOUT"

	// Cancellation & External
	CodeCanceled      = "CANCELED"
	CodeExternalError = "EXTERNAL_ERROR"

	// Internal server
	CodeInternalError = "INTERNAL_ERROR"

	// -----------------------------------------------------------------
	// User / Auth specific codes
	// -----------------------------------------------------------------
	CodeUserPendingVerification = "USER_PENDING_VERIFICATION"
	CodeInvalidVerificationCode = "INVALID_VERIFICATION_CODE"
	CodeVerificationExpired     = "VERIFICATION_EXPIRED"
	CodeEmailAlreadyRegistered  = "EMAIL_ALREADY_REGISTERED"
	CodeLoginInvalidCredentials = "LOGIN_INVALID_CREDENTIALS"

	// -----------------------------------------------------------------
	// Booking / Appointment codes
	// -----------------------------------------------------------------
	CodeAppointmentNotFound     = "APPOINTMENT_NOT_FOUND"
	CodeBookingNotFound         = "BOOKING_NOT_FOUND"
	CodeBookingSlotUnavailable  = "BOOKING_SLOT_UNAVAILABLE"
	CodeBookingCapacityExceeded = "BOOKING_CAPACITY_EXCEEDED"
	CodeBanListBlocked          = "BANLIST_BLOCKED"

	// -----------------------------------------------------------------
	// Database error codes
	// -----------------------------------------------------------------
	CodeDBUniqueViolation      = "DB_UNIQUE_VIOLATION"
	CodeDBForeignKeyConstraint = "DB_FOREIGN_KEY_CONSTRAINT"
	CodeDBSerializationFailure = "DB_SERIALIZATION_FAILURE"
	CodeDBDeadlockDetected     = "DB_DEADLOCK_DETECTED"
	CodeDBConnectionFailed     = "DB_CONNECTION_FAILED"

	// -----------------------------------------------------------------
	// External service error codes
	// -----------------------------------------------------------------
	CodeExternalSendgridForbidden   = "EXTERNAL_SENDGRID_FORBIDDEN"
	CodeExternalSendgridUnavailable = "EXTERNAL_SENDGRID_UNAVAILABLE"

	// -----------------------------------------------------------------
	// Repository layer codes
	// -----------------------------------------------------------------
	CodeRepoValidationError = "REPO_VALIDATION_ERROR"
	CodeRepoConflictError   = "REPO_CONFLICT_ERROR"
	CodeRepoNotFoundError   = "REPO_NOT_FOUND_ERROR"
	CodeRepoInternalError   = "REPO_INTERNAL_ERROR"
)
