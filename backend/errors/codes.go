package errors

const (
    // Generic
    CodeValidationFailed   = "VALIDATION_FAILED"
    CodeUnauthorized       = "AUTH_UNAUTHORIZED"
    CodeForbidden          = "FORBIDDEN"
    CodeResourceNotFound   = "RESOURCE_NOT_FOUND"
    CodeConflict           = "CONFLICT"
    CodeRateLimited        = "RATE_LIMITED"
    CodePreconditionFailed = "PRECONDITION_FAILED"
    CodeTimeout            = "TIMEOUT"
    CodeCanceled           = "CANCELED"
    CodeExternalError      = "EXTERNAL_ERROR"
    CodeInternalError      = "INTERNAL_ERROR"

    // User/Auth
    CodeUserPendingVerification = "USER_PENDING_VERIFICATION"
    CodeInvalidVerificationCode = "INVALID_VERIFICATION_CODE"
    CodeVerificationExpired     = "VERIFICATION_EXPIRED"
    CodeEmailAlreadyRegistered  = "EMAIL_ALREADY_REGISTERED"
    CodeLoginInvalidCredentials = "LOGIN_INVALID_CREDENTIALS"

    // Booking/Appointment
    CodeAppointmentNotFound   = "APPOINTMENT_NOT_FOUND"
    CodeBookingNotFound       = "BOOKING_NOT_FOUND"
    CodeBookingSlotUnavailable = "BOOKING_SLOT_UNAVAILABLE"
    CodeBookingCapacityExceeded = "BOOKING_CAPACITY_EXCEEDED"
    CodeBanListBlocked        = "BANLIST_BLOCKED"

    // DB
    CodeDBUniqueViolation       = "DB_UNIQUE_VIOLATION"
    CodeDBForeignKeyConstraint  = "DB_FOREIGN_KEY_CONSTRAINT"
    CodeDBSerializationFailure  = "DB_SERIALIZATION_FAILURE"
    CodeDBDeadlockDetected      = "DB_DEADLOCK_DETECTED"
    CodeDBConnectionFailed      = "DB_CONNECTION_FAILED"

    // External
    CodeExternalSendgridForbidden   = "EXTERNAL_SENDGRID_FORBIDDEN"
    CodeExternalSendgridUnavailable = "EXTERNAL_SENDGRID_UNAVAILABLE"
)

